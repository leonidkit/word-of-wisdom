package challenges

import (
	"context"
	"sync"

	"github.com/leonidkit/word-of-wisdom/pkg/hashcash"
)

type Repo struct {
	mx         *sync.RWMutex
	challenges map[string]*hashcash.HashcashHeader
}

func New() *Repo {
	return &Repo{
		mx:         new(sync.RWMutex),
		challenges: make(map[string]*hashcash.HashcashHeader),
	}
}

func (r *Repo) IsChallengeExists(ctx context.Context, key string) (bool, error) {
	r.mx.RLock()
	defer r.mx.RUnlock()

	if _, ok := r.challenges[key]; ok {
		return true, nil
	}
	return false, nil
}

func (r *Repo) DeleteChallenge(ctx context.Context, key string) error {
	r.mx.Lock()
	defer r.mx.Unlock()

	delete(r.challenges, key)
	return nil
}

func (r *Repo) AddChallenge(ctx context.Context, key string, header *hashcash.HashcashHeader) error {
	r.mx.Lock()
	defer r.mx.Unlock()

	r.challenges[key] = header
	return nil
}
