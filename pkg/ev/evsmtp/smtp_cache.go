package evsmtp

import (
	"fmt"
	"github.com/go-email-validator/go-email-validator/pkg/ev/evcache"
	"github.com/go-email-validator/go-email-validator/pkg/ev/evmail"
	"github.com/go-email-validator/go-email-validator/pkg/log"
	"go.uber.org/zap"
)

// ARandomRCPT is abstract realization of RandomRCPT
type ARandomRCPT struct {
	fn RandomRCPTFunc
}

// Call is calling of RandomRCPTFunc
func (a *ARandomRCPT) Call(sm SendMail, email evmail.Address) []error {
	return a.fn(sm, email)
}

func (a *ARandomRCPT) set(fn RandomRCPTFunc) {
	a.fn = fn
}

func (a *ARandomRCPT) get() RandomRCPTFunc {
	return a.fn
}

// RandomCacheKeyGetter is type of function to get cache key
type RandomCacheKeyGetter func(email evmail.Address) interface{}

// DefaultRandomCacheKeyGetter generates of cache key for RandomRCPT
func DefaultRandomCacheKeyGetter(email evmail.Address) interface{} {
	return email.Domain()
}

// NewCheckerCacheRandomRCPT creates Checker with caching of RandomRCPT calling
func NewCheckerCacheRandomRCPT(checker CheckerWithRandomRCPT, cache evcache.Interface, getKey RandomCacheKeyGetter) Checker {
	if getKey == nil {
		getKey = DefaultRandomCacheKeyGetter
	}

	c := &checkerCacheRandomRCPT{
		CheckerWithRandomRCPT: checker,
		randomRCPT:            &ARandomRCPT{fn: checker.get()},
		cache:                 cache,
		getKey:                getKey,
	}

	c.CheckerWithRandomRCPT.set(c.RandomRCPT)

	return c
}

type checkerCacheRandomRCPT struct {
	CheckerWithRandomRCPT
	randomRCPT RandomRCPT
	cache      evcache.Interface
	getKey     RandomCacheKeyGetter
}

func (c checkerCacheRandomRCPT) RandomRCPT(sm SendMail, email evmail.Address) (errs []error) {
	key := c.getKey(email)
	resultInterface, err := c.cache.Get(key)
	if err == nil && resultInterface != nil {
		errs = *resultInterface.(*[]error)
	} else {
		errs = c.randomRCPT.Call(sm, email)
		if err = c.cache.Set(key, ErrorsToEVSMTPErrors(errs)); err != nil {
			log.Logger().Error(fmt.Sprintf("cache RandomRCPT: %s", err),
				zap.String("email", email.String()),
				zap.String("key", fmt.Sprint(key)),
			)
		}
	}

	return errs
}
