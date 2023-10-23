package messages

import (
	"encoding/base64"
	"fmt"

	"github.com/leonidkit/word-of-wisdom/internal/models/messages"
)

type Encode struct{}

func (e Encode) Encode(msg messages.Message) ([]byte, error) {
	res := Message{}

	switch v := msg.(type) {
	case *messages.ChallengeRequestMessage:
		r := ChallengeRequestMessage{
			Challenge: v.Challenge,
		}
		if err := res.MergeChallengeRequestMessage(r); err != nil {
			return nil, fmt.Errorf("merge challenge request message: %v", err)
		}
	case *messages.ChallengeResponseMessage:
		r := ChallengeResponseMessage{
			Solution: v.Solution,
		}
		if err := res.MergeChallengeResponseMessage(r); err != nil {
			return nil, fmt.Errorf("merge challenge response message: %v", err)
		}
	case *messages.QuoteRequestMessage:
		if err := res.MergeWordOfWisdomRequestMessage(WordOfWisdomRequestMessage{}); err != nil {
			return nil, fmt.Errorf("merge word of wisdom request message: %v", err)
		}
	case *messages.QuoteResponseMessage:
		r := WordOfWisdomResponseMessage{
			Quote: v.Quote,
		}
		if err := res.MergeWordOfWisdomResponseMessage(r); err != nil {
			return nil, fmt.Errorf("merge word of wisdom response message: %v", err)
		}
	default:
		return nil, fmt.Errorf("unknown event: %v", msg)
	}

	jsonData, err := res.MarshalJSON()
	if err != nil {
		return nil, fmt.Errorf("result marshaling: %v", err)
	}

	bytesData := make([]byte, base64.RawStdEncoding.EncodedLen(len(jsonData)))
	base64.RawStdEncoding.Encode(bytesData, jsonData)

	return bytesData, nil
}
