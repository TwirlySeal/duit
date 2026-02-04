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

// https://httpwg.org/specs/rfc9110.html#field.accept-encoding
type parser struct {
	content string
	bestEncoding uint8
	bestWeight float32
	start int
	position int
	precedences map[string]uint8
}

func New(content string, precedences map[string]uint8) *parser {
	return &parser{
		content: content,
		precedences: precedences,
	}
}

func (p *parser) addEncoding(name string, quality float32) {
	precedence, found := p.precedences[name]
	if !found || quality == 0 {
		return
	}

	if quality > p.bestWeight || (quality == p.bestWeight && precedence > p.bestEncoding) {
		p.bestEncoding = precedence
		p.bestWeight = quality
	}
}

// Encodings with invalid quality values are ignored
func (p *parser) addParsedEncoding(name string, quality string) {
	if len(quality) > 4 {
		// "A sender of qvalue MUST NOT generate more than three digits
		// after the decimal point."
		return
	}
	
	num, err := strconv.ParseFloat(quality, 32)
	if err == nil && num >= 0 && num <= 1 {
		p.addEncoding(name, float32(num))
	}
}

func (p *parser) next() {
	p.position++
}

func (p *parser) slice() string {
	return p.content[p.start:p.position]
}

func (p *parser) currentByte() byte {
	return p.content[p.position]
}

func (p *parser) isEnd() bool {
	return p.position == len(p.content)
}

func (p *parser) setStart() {
	p.start = p.position
}

func (p *parser) skipSpace() {
	for !p.isEnd() && (p.currentByte() == ' ' || p.currentByte() == '\t') {
		p.next()
	}
	p.setStart()
}

// todo: fallback to identity or * if no valid encoding is specified but they are not prohibited
func (p *parser) Run() uint8 {
	if len(p.content) == 0 {
		// "An Accept-Encoding header field with a field value that is empty
		// implies that the user agent does not want any content coding in response."
		return identity
	}

	mainLoop:
	for {
		if p.isEnd() {
			p.addEncoding(p.slice(), 1.0)
			break
		}

		switch p.currentByte() {
			case ' ':
				p.skipSpace()

			case ',':
				p.addEncoding(p.slice(), 1.0)
				p.next()
				p.skipSpace()

			case ';':
				name := p.slice()
				p.next()
				p.skipSpace()
	
				// todo: handle malformed syntax here
				p.position += 2 // skip "q="
				p.start = p.position

				for {
					if p.isEnd() {
						p.addParsedEncoding(name, p.slice())
						break mainLoop
					} else if p.currentByte() == ',' {
						p.addParsedEncoding(name, p.slice())
						p.next()
						p.skipSpace()
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
		acceptEncoding, exists := r.Header["Accept-Encoding"]

		var encoding uint8
		if exists {
			encoding = New(
				acceptEncoding[0],
				defaultPrecedences,
			).Run()
		} else {
			// "If no Accept-Encoding header field is in the request, any
			// content coding is considered acceptable by the user agent."
			encoding = zs
		}

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
				// https://httpwg.org/specs/rfc9110.html#conneg.absent
				w.WriteHeader(http.StatusNotAcceptable)
				return
		}

		defer encoder.Close()
		w.Header().Set("Content-Encoding", encodingName)
		next.ServeHTTP(compressedResponseWriter{w, encoder}, r)
	})
}
