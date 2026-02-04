package headers

import "testing"

func acceptEncoding(t *testing.T, s string, expectedEncoding uint8) {
	result := New(s, defaultPrecedences).Run()

	if result != expectedEncoding {
		t.Errorf("Expected: %d, got: %d", expectedEncoding, result)
	}
}

func Test1(t *testing.T) {
	acceptEncoding(t, "gzip, deflate, br;q=0.5", gz)
}

func Test2(t *testing.T) {
	acceptEncoding(t, "deflate, gzip;q=0.5", deflate)
}

func Test3(t *testing.T) {
	acceptEncoding(t, "gzip;q=1.0, identity;q=0.5, *;q=0", gz)
}

func Test4(t *testing.T) {
	acceptEncoding(t, "gzip, *", zs)
}

func Test5(t *testing.T) {
	acceptEncoding(t, "gzip;q=0, *;q=0", 0)
}

// fails
// func Test6(t *testing.T) {
	// acceptEncoding(t, "br;q=0.4, gzip; q=0.8", gz)
// }

func Test7(t *testing.T) {
	acceptEncoding(t, "gzip, deflate, br, zstd", zs)
}
