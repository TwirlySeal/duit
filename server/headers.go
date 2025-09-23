package main

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

// todo: support whitespace and handle malformed header content
func parseEncoding(s string) uint8 {
	var alg uint8
	var weight float32
	start := 0

	precedences := map[string]uint8{
		"identity": identity,
		"deflate": deflate,
		"gzip": gz,
		"br": br,
		"zstd": zs,
		"*": zs,
	}

	update := func(name string, weight2 float32) {
		c, found := precedences[name]
		if !found || weight2 == 0 {
			return
		}
		
		if weight2 > weight || (weight2 == weight && c > alg) {
			alg = c
			weight = weight2
		}
	}

	parseUpdate := func(name2 string, s string) {
		result, err := strconv.ParseFloat(s, 32)
		if err == nil {
			weight2 := float32(result)
			update(name2, weight2)
		}
	}

	// Accept-Encoding only contains ASCII
	i := 0
	outer: for {
		if i == len(s) {
			update(s[start:i], 1.0)
			break
		}

		switch s[i] {
		case ',':
			update(s[start:i], 1.0)

			i += 2 // skip ", "
			start = i

		case ';':
			name := s[start:i]
			i += 3 // skip ";q="
			start = i

			for {
				if i == len(s) {
					parseUpdate(name, s[start:i])
					break outer
				} else if s[i] == ',' {
					parseUpdate(name, s[start:i])
					i += 2 // skip ", "
					start = i
					break
				}

				i++
			}
		}
		i++
	}

	return alg
}

// Select the optimal encoder from the 'Accept-Encoding' request header.
// The 'Content-Encoding' header of the response is set to the selected encoding.
// Returns nil for the 'identity' encoding.
func selectEncoding(w http.ResponseWriter, r *http.Request) io.WriteCloser {
	switch parseEncoding(r.Header.Get("Accept-Encoding")) {
		case identity:
			return nil
		case deflate:
			w.Header().Set("Content-Encoding", "deflate")
			encoder, _ := flate.NewWriter(w, flate.DefaultCompression)
			return encoder
		case gz:
			w.Header().Set("Content-Encoding", "gzip")
			return gzip.NewWriter(w)
		case br:
			w.Header().Set("Content-Encoding", "br")
			return brotli.NewWriter(w)
		case zs:
			w.Header().Set("Content-Encoding", "zstd")
			encoder, _ := zstd.NewWriter(w)
			return encoder
		default:
			return nil
	}
}
