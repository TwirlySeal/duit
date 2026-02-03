package main

import (
	"compress/gzip"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/andybalholm/brotli"
	"github.com/evanw/esbuild/pkg/api"
	"github.com/klauspost/compress/zstd"
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/html"
)

type Config struct {
	OutDir string `toml:"out_dir"`

	JsFiles []string `toml:"js_files"`

	CssFiles []string `toml:"css_files"`

	HtmlFiles []string `toml:"html_files"`

	// JS struct {
	// 	EntryPoints []string `toml:"entry_points"`
	// 	Minify bool `toml:"minify"`
	// 	Bundle bool `toml:"bundle"`
	// } `toml:"js"`

	// HTML struct {
	// 	SourceDir string `toml:"source_dir"`
	// }
}

func main() {
	var config Config
	_, err := toml.DecodeFile("build.toml", &config)
	if err != nil {
		log.Fatalf("Failed to load config: %v\n", err)
	}

	result := api.Build(api.BuildOptions{
		EntryPoints: append(config.JsFiles, config.CssFiles...),
		Outdir: config.OutDir,
		Write: true,
		Bundle: true,
		MinifyWhitespace: true,
		MinifyIdentifiers: true,
		MinifySyntax: true,
		// LogLevel: api.LogLevelInfo,
	})

	if len(result.Errors) > 0 {
		os.Exit(1)
	}

	for _, out := range result.OutputFiles {
		compress(out.Path, out.Contents)
	}

	minifier := minify.New()
	minifier.Add("text/html", &html.Minifier{
		KeepDocumentTags: true,
		KeepEndTags: true,
		KeepQuotes: true,
		TemplateDelims: html.GoTemplateDelims,
	})

	for _, path := range config.HtmlFiles {
		source, err := os.Open(path)
		if err != nil {
			panic(err)
		}
		destPath := filepath.Join(
			config.OutDir,
			filepath.Base(path),
		)
		destFile, err := os.Create(destPath)
		if err != nil {
			panic(err)
		}
		minifier.Minify("text/html", destFile, source)
	}
}

func compress(path string, data []byte) {
	write := func(extension string, compressor func(file *os.File) io.WriteCloser) {
		file, err := os.Create(path + extension)
		if err != nil {
			panic(err)
		}
		writer := compressor(file)
		writer.Write(data)
		writer.Close()
	}

	write(".gz", func(file *os.File) io.WriteCloser {
		w, err := gzip.NewWriterLevel(file, gzip.BestCompression)
		if err != nil {
			panic(err)
		}
		return w
	})

	write(".br", func(file *os.File) io.WriteCloser {
		return brotli.NewWriterLevel(file, brotli.BestCompression)
	})

	write(".zst", func(file *os.File) io.WriteCloser {
		writer, err := zstd.NewWriter(file, zstd.WithEncoderLevel(zstd.SpeedBestCompression))
		if err != nil {
			panic(err)
		}
		return writer
	})
}
