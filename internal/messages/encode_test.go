package messages

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/leonidkit/word-of-wisdom/internal/models/messages"
)

func TestAdapter_Adapt(t *testing.T) {
	cases := []struct {
		name string
		msg  messages.Message
		exp  string
	}{
		{
			name: "smoke",
			msg:  messages.NewQuoteResponseMessage("some quote"),
			exp:  "eyJtZXNzYWdlVHlwZSI6IldvcmRPZldpc2RvbVJlc3BvbnNlTWVzc2FnZSIsInF1b3RlIjoic29tZSBxdW90ZSJ9",
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			encoded, err := Encode{}.Encode(tt.msg)
			require.NoError(t, err)
			assert.Equal(t, tt.exp, string(encoded))
		})
	}
}
