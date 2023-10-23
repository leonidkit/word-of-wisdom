package challenger

import (
	"context"
	"fmt"

	"github.com/leonidkit/word-of-wisdom/pkg/hashcash"
)

var defaultChallengeComplexity = 10

//go:generate mockgen -source=$GOFILE -destination=./mocks/challenger_mocks.gen.go -package=challenger_repo_mocks

type challengesRepository interface {
	IsChallengeExists(ctx context.Context, key string) (bool, error)
	DeleteChallenge(ctx context.Context, key string) error
	AddChallenge(ctx context.Context, key string, header *hashcash.HashcashHeader) error
}

type Challenger struct {
	complexity     int
	challengesRepo challengesRepository

	HashcashNewHeader   func(bitsAmount int, resource string) (*hashcash.HashcashHeader, error)
	HashcashParseHeader func(header string) (*hashcash.HashcashHeader, error)
	HashcashVerify      func(header *hashcash.HashcashHeader) bool
}

func WithComplexity(complexity int) func(*Challenger) {
	return func(c *Challenger) {
		c.complexity = complexity
	}
}

func New(challengesRepo challengesRepository, opts ...func(*Challenger)) *Challenger {
	c := &Challenger{
		complexity:     defaultChallengeComplexity,
		challengesRepo: challengesRepo,

		HashcashNewHeader:   hashcash.NewHeader,
		HashcashParseHeader: hashcash.ParseHeader,
		HashcashVerify:      hashcash.Verify,
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

func (c Challenger) NewChallenge(ctx context.Context, domain string) (string, error) {
	hh, err := c.HashcashNewHeader(c.complexity, domain)
	if err != nil {
		return "", fmt.Errorf("`server-challenger service` hashcash new header: %v", err)
	}
	challengeString := hh.String()

	err = c.challengesRepo.AddChallenge(ctx, hh.Rand, hh)
	if err != nil {
		return "", fmt.Errorf("`server-challenger service` add challenge to repo: %v", err)
	}
	return challengeString, nil
}

func (c Challenger) Validate(ctx context.Context, challenge string) (bool, error) {
	hh, err := c.HashcashParseHeader(challenge)
	if err != nil {
		return false, fmt.Errorf("`server-challenger service` hashcash parse header: %v", err)
	}

	exists, err := c.challengesRepo.IsChallengeExists(ctx, hh.Rand)
	if err != nil {
		return false, fmt.Errorf("`server-challenger service` challenge exists: %v", err)
	}
	if !exists {
		return false, fmt.Errorf("`server-challenger service` challenge not found")
	}

	err = c.challengesRepo.DeleteChallenge(ctx, hh.Rand)
	if err != nil {
		return false, fmt.Errorf("`server-challenger service` delete challenge")
	}

	return c.HashcashVerify(hh), nil
}
