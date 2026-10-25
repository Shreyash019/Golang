package banner

import "testing"

func TestBanner(t *testing.T) {
	const want = "orders-api listening on :8080"
	if got := Banner(); got != want {
		t.Errorf("Banner() = %q, want %q", got, want)
	}
}
