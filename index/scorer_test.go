package index

import (
	"math"
	"testing"
)

func TestScore(t *testing.T) {
	tests := []struct {
		name       string
		posting    Posting
		doc        Document
		docListLen int
		totalDocs  int
		want       float64
		expectZero bool
	}{
		{
			name:       "happy path",
			posting:    Posting{Frequency: 2},
			doc:        Document{Length: 4},
			docListLen: 2,
			totalDocs:  10,
			want:       0.5 * math.Log(5),
		},
		{
			name:       "zero document length",
			posting:    Posting{Frequency: 3},
			doc:        Document{Length: 0},
			docListLen: 2,
			totalDocs:  10,
			expectZero: true,
		},
		{
			name:       "rare term higher score then common term",
			posting:    Posting{Frequency: 1},
			doc:        Document{Length: 10},
			docListLen: 1,
			totalDocs:  10,
			want:       (1.0 / 10.0) * math.Log(10),
		},
		{
			name:       "common term lower score",
			posting:    Posting{Frequency: 1},
			doc:        Document{Length: 10},
			docListLen: 5,
			totalDocs:  10,
			want:       (1.0 / 10.0) * math.Log(2),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Score(tt.posting, tt.doc, tt.docListLen, tt.totalDocs)

			if tt.expectZero {
				if got != 0.0 {
					t.Errorf("Score() = %f, want 0.0", got)
				}
				return
			}

			if math.Abs(got-tt.want) > 1e-9 {
				t.Errorf("Score() = %f, want %f", got, tt.want)
			}
		})
	}
}
