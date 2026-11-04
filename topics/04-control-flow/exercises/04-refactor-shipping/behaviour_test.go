package shipping

import (
	"slices"
	"testing"
	"time"
)

func TestQuote(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		p       Parcel
		want    int64
		wantErr bool
	}{
		{"India, 500 g", Parcel{WeightGrams: 500, Country: "IN"}, 700, false},
		{"US, exactly 1 kg", Parcel{WeightGrams: 1000, Country: "US"}, 1700, false},
		{"Canada, a started 2nd kg", Parcel{WeightGrams: 1001, Country: "CA"}, 1900, false},
		{"France, express", Parcel{WeightGrams: 2500, Country: "FR", Express: true}, 3600, false},
		{"Netherlands, heaviest allowed", Parcel{WeightGrams: MaxWeightGrams, Country: "NL"}, 7200, false},
		{"too heavy", Parcel{WeightGrams: MaxWeightGrams + 1, Country: "IN"}, 0, true},
		{"zero weight", Parcel{Country: "IN"}, 0, true},
		{"negative weight", Parcel{WeightGrams: -1, Country: "DE"}, 0, true},
		{"unsupported country", Parcel{WeightGrams: 100, Country: "BR"}, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := Quote(tt.p)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Quote(%+v) error = %v, want error: %t", tt.p, err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("Quote(%+v) = %d, want %d", tt.p, got, tt.want)
			}
		})
	}
}

func TestCheapestIndex(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		quotes []int64
		want   int
	}{
		{"empty", nil, -1},
		{"one", []int64{900}, 0},
		{"first tie wins", []int64{500, 300, 300, 900}, 1},
		{"last", []int64{5, 4, 3}, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := CheapestIndex(tt.quotes); got != tt.want {
				t.Errorf("CheapestIndex(%v) = %d, want %d", tt.quotes, got, tt.want)
			}
		})
	}
}

func TestHasRestricted(t *testing.T) {
	t.Parallel()
	restricted := []string{"battery", "perfume"}
	tests := []struct {
		name      string
		shipments [][]string
		want      bool
	}{
		{"none", nil, false},
		{"clean", [][]string{{"book", "pen"}, {"mug"}}, false},
		{"in the last shipment", [][]string{{"book"}, {"mug", "perfume"}}, true},
		{"first item", [][]string{{"battery"}}, true},
		{"empty shipments", [][]string{{}, {}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := HasRestricted(tt.shipments, restricted); got != tt.want {
				t.Errorf("HasRestricted(%v, %v) = %t, want %t", tt.shipments, restricted, got, tt.want)
			}
		})
	}
}

func TestBackoff(t *testing.T) {
	t.Parallel()
	ms := time.Millisecond
	tests := []struct {
		name string
		n    int
		want []time.Duration
	}{
		{"zero retries", 0, nil},
		{"negative", -2, nil},
		{"one", 1, []time.Duration{100 * ms}},
		{"doubling", 3, []time.Duration{100 * ms, 200 * ms, 400 * ms}},
		{"capped at the limit", 6, []time.Duration{100 * ms, 200 * ms, 400 * ms, 800 * ms, 1000 * ms, 1000 * ms}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := Backoff(100*ms, time.Second, tt.n); !slices.Equal(got, tt.want) {
				t.Errorf("Backoff(100ms, 1s, %d) = %v, want %v", tt.n, got, tt.want)
			}
		})
	}
}
