package messages

import (
	"fmt"

	validatorlib "github.com/go-playground/validator/v10"
)

var (
	_ Message = (*ChallengeRequestMessage)(nil)
	_ Message = (*ChallengeResponseMessage)(nil)
	_ Message = (*QuoteRequestMessage)(nil)
	_ Message = (*QuoteResponseMessage)(nil)
)

var validator = validatorlib.New()

//sumtype:decl
type Message interface {
	messageMarker()
	Validate() error
	Name() string
}

type message struct{}

func (*message) messageMarker() {}

type ChallengeRequestMessage struct {
	message
	Challenge string `validate:"required"`
}

func NewChallengeRequestMessage(challenge string) *ChallengeRequestMessage {
	return &ChallengeRequestMessage{
		Challenge: challenge,
	}
}

func (m *ChallengeRequestMessage) Name() string {
	return "ChallengeRequestMessage"
}

func (m *ChallengeRequestMessage) Validate() error {
	if err := validator.Struct(m); err != nil {
		return fmt.Errorf("validation: %v", err)
	}
	return nil
}

type ChallengeResponseMessage struct {
	message
	Solution string `validate:"required"`
}

func NewChallengeResponseMessage(solution string) *ChallengeResponseMessage {
	return &ChallengeResponseMessage{
		Solution: solution,
	}
}

func (m *ChallengeResponseMessage) Name() string {
	return "ChallengeResponseMessage"
}

func (m *ChallengeResponseMessage) Validate() error {
	if err := validator.Struct(m); err != nil {
		return fmt.Errorf("validation: %v", err)
	}
	return nil
}

type QuoteRequestMessage struct {
	message
}

func NewQuoteRequestMessage() *QuoteRequestMessage {
	return &QuoteRequestMessage{}
}

func (m *QuoteRequestMessage) Name() string {
	return "QuoteRequestMessage"
}

func (m *QuoteRequestMessage) Validate() error {
	return nil
}

type QuoteResponseMessage struct {
	message
	Quote string `validate:"required"`
}

func NewQuoteResponseMessage(quote string) *QuoteResponseMessage {
	return &QuoteResponseMessage{
		Quote: quote,
	}
}

func (m *QuoteResponseMessage) Name() string {
	return "QuoteResponseMessage"
}

func (m *QuoteResponseMessage) Validate() error {
	if err := validator.Struct(m); err != nil {
		return fmt.Errorf("validation: %v", err)
	}
	return nil
}
