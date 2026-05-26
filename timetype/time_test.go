package timetype

import (
	"testing"
	"time"
)

func TestUnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		input   []byte
		want    string
		wantErr bool
	}{
		{
			name:  "RFC3339 input",
			input: []byte(`"2024-01-15T10:30:00Z"`),
			want:  "2024-01-15T10:30:00Z",
		},
		{
			name:  "datetime-local format",
			input: []byte(`"2024-01-15T10:30"`),
			want:  "2024-01-15T10:30:00Z",
		},
		{
			name:    "invalid input",
			input:   []byte(`"not-a-time"`),
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var ct CustomTime
			err := ct.UnmarshalJSON(tc.input)

			if tc.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			got := ct.UTC().Format(time.RFC3339)
			if got != tc.want {
				t.Errorf("got %q, want %q", got, tc.want)
			}
		})
	}
}

func TestScan(t *testing.T) {
	stringTestWant, _ := time.Parse("2006-01-02 15:04:05", "2006-01-02 15:04:05")
	byteTestWant, _ := time.Parse("2006-01-02 15:04:05", string([]byte("2006-01-02 15:04:05")))
	tests := []struct {
		name    string
		input   any
		want    time.Time
		wantErr bool
	}{
		{
			name:    "string test",
			input:   "2006-01-02 15:04:05",
			want:    stringTestWant.UTC(),
			wantErr: false,
		},
		{
			name:    "byte test",
			input:   []byte("2006-01-02 15:04:05"),
			want:    byteTestWant.UTC(),
			wantErr: false,
		},
		{
			name:    "int test",
			input:   int64(1779828983),
			want:    time.Unix(int64(1779828983), 0).UTC(),
			wantErr: false,
		},
		{
			name:    "nil test",
			input:   nil,
			want:    time.Time{}.UTC(),
			wantErr: false,
		},
		{
			name:    "invalid type",
			input:   true,
			want:    time.Time{}.UTC(),
			wantErr: true,
		},
	}

	for _, tc := range tests {
		var ct CustomTime
		err := ct.Scan(tc.input)

		if tc.wantErr {
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			return
		}
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		_ = ct.Scan(tc.input)
		if ct.Time != tc.want {
			t.Errorf("%s got %q, want %q", tc.name, ct.Time, tc.want)
		}

	}
}
