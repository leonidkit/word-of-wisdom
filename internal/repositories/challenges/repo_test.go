package challenges_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/leonidkit/word-of-wisdom/internal/repositories/challenges"
	"github.com/leonidkit/word-of-wisdom/pkg/hashcash"
)

func TestRepoSmoke(t *testing.T) {
	repo := challenges.New()
	ctx := context.Background()
	key := "key"

	err := repo.AddChallenge(ctx, key, &hashcash.HashcashHeader{})
	require.NoError(t, err)

	ok, err := repo.IsChallengeExists(ctx, key)
	require.NoError(t, err)
	require.True(t, ok)

	err = repo.DeleteChallenge(ctx, key)
	require.NoError(t, err)

	ok, err = repo.IsChallengeExists(ctx, key)
	require.NoError(t, err)
	require.False(t, ok)
}
