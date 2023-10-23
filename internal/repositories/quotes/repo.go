package quotes

import (
	"context"
	_ "embed"
	"math/rand"
	"strings"
)

//go:embed quotes.txt
var quotes string

type Repo struct{}

func New() *Repo {
	return &Repo{}
}

func (r *Repo) GetRandom(ctx context.Context) (string, error) {
	splitted := strings.Split(quotes, "\n")
	return splitted[rand.Int()%len(splitted)], nil //nolint:gosec
}
