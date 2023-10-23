package challenger_test

import (
	"errors"

	"github.com/golang/mock/gomock"

	challenger "github.com/leonidkit/word-of-wisdom/internal/services/server-challenger"
	challengermocks "github.com/leonidkit/word-of-wisdom/internal/services/server-challenger/mocks"
	"github.com/leonidkit/word-of-wisdom/internal/testingh"
	"github.com/leonidkit/word-of-wisdom/pkg/hashcash"
)

type ChallengerSuite struct {
	testingh.ContextSuite

	domain string

	ctrl *gomock.Controller

	challengesRepo *challengermocks.MockchallengesRepository
	challenger     *challenger.Challenger
}

func (cs *ChallengerSuite) SetupTest() {
	cs.ctrl = gomock.NewController(cs.T())
	cs.challengesRepo = challengermocks.NewMockchallengesRepository(cs.ctrl)
	cs.challenger = challenger.New(cs.challengesRepo)
	cs.domain = "domain"
}

func (cs *ChallengerSuite) TestNewChallenge_Success() {
	// Action.
	challenge, err := cs.challenger.NewChallenge(cs.Ctx, cs.domain)

	// Assert.
	cs.Require().NoError(err)
	cs.Require().Empty(challenge)
}

func (cs *ChallengerSuite) TestNewChallenge_GenerateHeaderError() {
	// Arrange.
	cs.challenger.HashcashNewHeader = func(_ int, _ string) (*hashcash.HashcashHeader, error) {
		return nil, errors.New("unexpected")
	}

	// Action.
	c, err := cs.challenger.NewChallenge(cs.Ctx, cs.domain)

	// Assert.
	cs.Require().Error(err)
	cs.Require().Empty(c)
}

func (cs *ChallengerSuite) TestNewChallenge_RepositoryError() {
	// Arrange.
	cs.challengesRepo.EXPECT().AddChallenge(cs.Ctx, gomock.Any(), gomock.Any()).Return(errors.New("unexpected"))

	// Action.
	c, err := cs.challenger.NewChallenge(cs.Ctx, cs.domain)

	// Arrange.
	cs.Require().Error(err)
	cs.Require().Empty(c)
}

func (cs *ChallengerSuite) TestValidate_Success() {
	// Arrange.
	cs.challengesRepo.EXPECT().AddChallenge(cs.Ctx, gomock.Any(), gomock.Any()).Return(nil)

	// Action.
	c, err := cs.challenger.NewChallenge(cs.Ctx, cs.domain)

	// Arrange.
	cs.Require().NoError(err)
	cs.Require().NotEmpty(c)
}

func (cs *ChallengerSuite) TestValidate_ParseError() {
	// Arrange.
	cs.challenger.HashcashParseHeader = func(header string) (*hashcash.HashcashHeader, error) {
		return nil, errors.New("unexpected")
	}

	// Action.
	ok, err := cs.challenger.Validate(cs.Ctx, "challenge")

	// Arrange.
	cs.Require().Error(err)
	cs.Require().False(ok)
}

func (cs *ChallengerSuite) TestValidate_UnknownChallenge() {
	// Arrange.
	cs.challenger.HashcashParseHeader = func(header string) (*hashcash.HashcashHeader, error) {
		return &hashcash.HashcashHeader{}, nil
	}
	cs.challengesRepo.EXPECT().IsChallengeExists(cs.Ctx, gomock.Any()).Return(false, nil)

	// Action.
	ok, err := cs.challenger.Validate(cs.Ctx, "challenge")

	// Arrange.
	cs.Require().Error(err)
	cs.Require().False(ok)
}

func (cs *ChallengerSuite) TestValidate_RepositoryError() {
	// Arrange.
	cs.challenger.HashcashParseHeader = func(header string) (*hashcash.HashcashHeader, error) {
		return &hashcash.HashcashHeader{}, nil
	}
	cs.challengesRepo.EXPECT().IsChallengeExists(cs.Ctx, gomock.Any()).Return(true, nil)
	cs.challengesRepo.EXPECT().DeleteChallenge(cs.Ctx, gomock.Any()).Return(errors.New("unexpected"))

	// Action.
	ok, err := cs.challenger.Validate(cs.Ctx, "challenge")

	// Arrange.
	cs.Require().Error(err)
	cs.Require().False(ok)
}

func (cs *ChallengerSuite) TestValidate_VerifyError() {
	// Arrange.
	cs.challenger.HashcashParseHeader = func(header string) (*hashcash.HashcashHeader, error) {
		return &hashcash.HashcashHeader{}, nil
	}
	cs.challenger.HashcashVerify = func(header *hashcash.HashcashHeader) bool {
		return true
	}
	cs.challengesRepo.EXPECT().IsChallengeExists(cs.Ctx, gomock.Any()).Return(true, nil)
	cs.challengesRepo.EXPECT().DeleteChallenge(cs.Ctx, gomock.Any()).Return(nil)

	// Action.
	ok, err := cs.challenger.Validate(cs.Ctx, "challenge")

	// Arrange.
	cs.Require().NoError(err)
	cs.Require().True(ok)
}
