package challenger

import (
	"context"
	"fmt"

	"github.com/leonidkit/word-of-wisdom/pkg/hashcash"
)

type Challenger struct{}

func New() Challenger {
	return Challenger{}
}

func (c Challenger) Accept(ctx context.Context, challenge string) (string, error) {
	hh, err := hashcash.ParseHeader(challenge)
	if err != nil {
		return "", fmt.Errorf("`client-challenger service` hashcash parse header: %v", err)
	}

	err = hashcash.CalculateCounter(ctx, hh)
	if err != nil {
		return "", fmt.Errorf("`client-challenger service` hashcash calculate counter: %v", err)
	}

	return hh.String(), nil
}
