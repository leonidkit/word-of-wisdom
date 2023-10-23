package hashcash

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestCalculateCounter(t *testing.T) {
	t.Parallel()

	t.Run("unsupported algorithm", func(t *testing.T) {
		hh := generateHeader()
		hh.Alg = "RND"
		err := CalculateCounter(context.Background(), hh)
		require.Error(t, err)
		require.ErrorIs(t, err, ErrAlgNotSupported)
	})

	t.Run("process interrupted", func(t *testing.T) {
		hh := generateHeader()
		hh.Bits = 100

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		err := CalculateCounter(ctx, hh)
		require.Error(t, err)
		require.ErrorIs(t, err, ErrCalculationInterrupted)
	})

	t.Run("success calculation", func(t *testing.T) {
		hh := generateHeader()

		err := CalculateCounter(context.Background(), hh)
		require.NoError(t, err)
		require.Equal(t, uint64(3450), hh.Counter)
	})
}

func TestVerify(t *testing.T) {
	t.Run("invalid header", func(t *testing.T) {
		hh := generateHeader()
		require.False(t, Verify(hh))
	})

	t.Run("valid header", func(t *testing.T) {
		hh := generateHeader()
		err := CalculateCounter(context.Background(), hh)
		require.NoError(t, err)
		require.True(t, Verify(hh))
	})
}

func TestHashcash_String(t *testing.T) {
	hh := &HashcashHeader{
		Version: 1,
		Bits:    10,
		Timestamp: func() int64 {
			tt, err := time.Parse(time.DateTime, "2023-01-02 15:04:05")
			require.NoError(t, err)
			return tt.Unix()
		}(),
		Alg:      "SHA-256",
		Resource: "resource",
		Rand:     "4PF4B5e0_spEr0b3n0OM4g",
		Counter:  123,
	}

	require.Equal(t, "1:10:1672671845:resource:SHA-256:4PF4B5e0_spEr0b3n0OM4g:MTIz", hh.String())
}

func TestParseHeader(t *testing.T) {
	tests := []struct {
		header  string
		want    *HashcashHeader
		wantErr bool
	}{
		// negative
		{
			header:  "dsff:10:1672671845:resource:SHA-256:4PF4B5e0_spEr0b3n0OM4g:MTIz",
			wantErr: true,
		},
		{
			header:  "1:10:1672671f845:resource:SHA-256:4PF4B5e0_spEr0b3n0OM4g:MTIz",
			wantErr: true,
		},
		{
			header:  "1:10:1672671f845:resource:SHA-256:4PF4B5e0_spEr0b3n0OM4g:sdf",
			wantErr: true,
		},

		// positive
		{
			header: "1:10:1672671845:resource:SHA-256:4PF4B5e0_spEr0b3n0OM4g:MTIz",
			want: &HashcashHeader{
				Version:   1,
				Bits:      10,
				Timestamp: 1672671845,
				Resource:  "resource",
				Alg:       "SHA-256",
				Rand:      "4PF4B5e0_spEr0b3n0OM4g",
				Counter:   123,
			},
		},
	}

	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			res, err := ParseHeader(tt.header)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.want, res)
		})
	}
}

func generateHeader() *HashcashHeader {
	return &HashcashHeader{
		Version:   1,
		Bits:      10,
		Timestamp: 1672671845,
		Resource:  "resource",
		Alg:       AlgorithmSHA256,
		Rand:      "4PF4B5e0_spEr0b3n0OM4g",
	}
}
