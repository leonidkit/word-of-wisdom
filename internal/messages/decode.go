package messages

import (
	"encoding/base64"
	"fmt"

	"github.com/leonidkit/word-of-wisdom/internal/models/messages"
)

type Decoder struct{}

func (d Decoder) Decode(msg []byte) (messages.Message, error) {
	msgDecoded := make([]byte, base64.RawStdEncoding.DecodedLen(len(msg)))
	_, err := base64.RawStdEncoding.Decode(msgDecoded, msg)
	if err != nil {
		return nil, fmt.Errorf("msg base64 decoding: %v", err)
	}

	baseMsg := Message{}
	err = baseMsg.UnmarshalJSON(msgDecoded)
	if err != nil {
		return nil, fmt.Errorf("unmarshal message: %v", err)
	}

	value, err := baseMsg.ValueByDiscriminator()
	if err != nil {
		return nil, fmt.Errorf("message value by discriminator: %v", err)
	}

	switch v := value.(type) {
	case ChallengeRequestMessage:
		return messages.NewChallengeRequestMessage(v.Challenge), nil
	case ChallengeResponseMessage:
		return messages.NewChallengeResponseMessage(v.Solution), nil
	case WordOfWisdomRequestMessage:
		return messages.NewQuoteRequestMessage(), nil
	case WordOfWisdomResponseMessage:
		return messages.NewQuoteResponseMessage(v.Quote), nil
	default:
		return nil, fmt.Errorf("unknown object: %v", string(msg))
	}
}
