package messages

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/leonidkit/word-of-wisdom/internal/models/messages"
)

func TestDecoder_Decode(t *testing.T) {
	cases := []struct {
		name    string
		objJSON string
		wantEv  messages.Message
	}{
		{
			name:    "smoke",
			objJSON: "eyJxdW90ZSI6ICJzb21lIHF1b3RlIiwgIm1lc3NhZ2VUeXBlIjogIldvcmRPZldpc2RvbVJlc3BvbnNlTWVzc2FnZSJ9Cg",
			wantEv:  messages.NewQuoteResponseMessage("some quote"),
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			decoded, err := Decoder{}.Decode([]byte(tt.objJSON))
			require.NoError(t, err)

			assert.Equal(t, tt.wantEv, decoded)
		})
	}
}
