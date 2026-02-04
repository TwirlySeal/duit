package headers

import "testing"

func acceptEncoding(t *testing.T, s string, expectedEncoding uint8) {
	result := New(s, defaultPrecedences).Run()

	if result != expectedEncoding {
		t.Errorf("Expected: %d, got: %d", expectedEncoding, result)
	}
}

func TestBrotliLowQuality(t *testing.T) {
	acceptEncoding(t, "gzip, deflate, br;q=0.5", gz)
}

func TestGzipLowQuality(t *testing.T) {
	acceptEncoding(t, "deflate, gzip;q=0.5", deflate)
}

func TestRejectElse(t *testing.T) {
	acceptEncoding(t, "gzip;q=1.0, identity;q=0.5, *;q=0", gz)
}

func TestGzipElse(t *testing.T) {
	acceptEncoding(t, "gzip, *", zs)
}

func TestNoValidEncoding(t *testing.T) {
	acceptEncoding(t, "gzip;q=0, *;q=0", 0)
}

func TestNoSpaces(t *testing.T) {
	acceptEncoding(t, "gzip,br", br)
}

func TestSpaceAfterSemicolon(t *testing.T) {
 acceptEncoding(t, "br;q=0.4, gzip; q=0.8", gz)
}

func TestBrowserValue(t *testing.T) {
	acceptEncoding(t, "gzip, deflate, br, zstd", zs)
}

func TestQualityTooHigh(t *testing.T) {
	acceptEncoding(t, "gzip;q=3, br", br)
}

func TestQualityTooManyDecimalPlaces(t *testing.T) {
	acceptEncoding(t, "zstd;=0.1111", 0)
}
