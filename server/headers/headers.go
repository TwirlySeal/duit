package headers

import (
	"compress/flate"
	"compress/gzip"
	"io"
	"net/http"
	"strconv"

	"github.com/andybalholm/brotli"
	"github.com/klauspost/compress/zstd"
)

const (
	_ = iota
	identity
	deflate
	gz
	br
	zs
)

// todo: handle variable whitespace and malformed header content
type Parser struct {
	content string
	bestEncoding uint8
	bestWeight float32
	start int
	position int
	precedences map[string]uint8
}

func New(content string, precedences map[string]uint8) *Parser {
	return &Parser{
		content: content,
		precedences: precedences,
	}
}

func (p *Parser) addEncoding(name string, quality float32) {
	precedence, found := p.precedences[name]
	if !found || quality == 0 {
		return
	}

	if quality > p.bestWeight || (quality == p.bestWeight && precedence > p.bestEncoding) {
		p.bestEncoding = precedence
		p.bestWeight = quality
	}
}

func (p *Parser) addParsedEncoding(name string, quality string) {
	num, err := strconv.ParseFloat(quality, 32)
	if err == nil {
		p.addEncoding(name, float32(num))
	}
}

func (p *Parser) next() {
	p.position++
}

func (p *Parser) slice() string {
	return p.content[p.start:p.position]
}

func (p *Parser) currentByte() byte {
	return p.content[p.position]
}

func (p *Parser) isEnd() bool {
	return p.position == len(p.content)
}

func (p *Parser) Run() uint8 {
	outer:
	for {
		if p.isEnd() {
			p.addEncoding(p.slice(), 1.0)
			break
		}

		switch p.currentByte() {
			case ',':
				name := p.slice()
				p.addEncoding(name, 1.0)

				p.position += 2 // skip ", "
				p.start = p.position

			case ';':
				name := p.slice()
				p.position += 3 // skip ";q="
				p.start = p.position

				for {
					if p.isEnd() {
						p.addParsedEncoding(name, p.slice())
						break outer
					} else if p.currentByte() == ',' {
						p.addParsedEncoding(name, p.slice())
						p.position += 2 // skip ", "
						p.start = p.position
						break
					}

					p.next()
				}
		}
		p.next()
	}

	return p.bestEncoding
}

var defaultPrecedences = map[string]uint8{
	"identity": identity,
	"deflate": deflate,
	"gzip": gz,
	"br": br,
	"zstd": zs,
	"*": zs,
}

// Select the optimal encoder from the 'Accept-Encoding' request header.
// The 'Content-Encoding' header of the response is set to the selected encoding.

type compressedResponseWriter struct {
	writer http.ResponseWriter
	encoder io.WriteCloser
}

func (c compressedResponseWriter) Header() http.Header {
	return c.writer.Header()
}

func (c compressedResponseWriter) Write(data []byte) (int, error) {
	return c.encoder.Write(data)
}

func (c compressedResponseWriter) WriteHeader(statusCode int) {
	c.writer.WriteHeader(statusCode)
}

// todo: avoid compression for small responses
func Compress(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		encoding := New(
			r.Header.Get("Accept-Encoding"),
			defaultPrecedences,
		).Run()

		var encodingName string
		var encoder io.WriteCloser
		switch encoding {
			case zs:
				encoder, _ = zstd.NewWriter(w)
				encodingName = "zstd"
			case br:
				encoder = brotli.NewWriter(w)
				encodingName = "br"
			case gz:
				encoder = gzip.NewWriter(w)
				encodingName = "gzip"
			case deflate:
				encoder, _ = flate.NewWriter(w, flate.DefaultCompression)
				encodingName = "deflate"
			case identity:
				// no compression
				next.ServeHTTP(w, r)
				return
			default:
				// no valid encoding
				w.WriteHeader(http.StatusUnsupportedMediaType)
				return
		}

		defer encoder.Close()
		w.Header().Set("Content-Encoding", encodingName)
		next.ServeHTTP(compressedResponseWriter{w, encoder}, r)
	})
}
