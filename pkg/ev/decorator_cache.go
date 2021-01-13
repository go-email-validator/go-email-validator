package ev

import (
	"fmt"
	"github.com/go-email-validator/go-email-validator/pkg/ev/evcache"
	"github.com/go-email-validator/go-email-validator/pkg/ev/evmail"
	"github.com/go-email-validator/go-email-validator/pkg/ev/utils"
	"github.com/go-email-validator/go-email-validator/pkg/log"
	"go.uber.org/zap"
)

// To use complex keys you can use https://github.com/vmihailenco/msgpack/
type CacheKeyGetter func(email evmail.Address, results ...ValidationResult) interface{}

func EmailCacheKeyGetter(email evmail.Address, _ ...ValidationResult) interface{} {
	return email.String()
}

func DomainCacheKeyGetter(email evmail.Address, _ ...ValidationResult) interface{} {
	return email.Domain()
}

func NewCacheDecorator(validator Validator, cache evcache.Interface, getKey CacheKeyGetter) Validator {
	if getKey == nil {
		getKey = EmailCacheKeyGetter
	}

	return &cacheDecorator{
		validator: validator,
		cache:     cache,
		getKey:    getKey,
	}
}

type cacheDecorator struct {
	validator Validator
	cache     evcache.Interface
	getKey    CacheKeyGetter
}

func (c *cacheDecorator) GetDeps() []ValidatorName {
	return c.validator.GetDeps()
}

func (c *cacheDecorator) Validate(email evmail.Address, results ...ValidationResult) (result ValidationResult) {
	key := c.getKey(email, results...)
	resultInterface, err := c.cache.Get(key)
	if err == nil && resultInterface != nil {
		result = *resultInterface.(*ValidationResult)
	} else {
		result = c.validator.Validate(email, results...)
		if err := c.cache.Set(key, result); err != nil {
			log.Logger().Error(fmt.Sprintf("cache decorator %v", err),
				zap.String("validator", utils.StructName(c.validator)),
				zap.String("key", fmt.Sprint(key)),
				zap.String("email", email.String()),
				zap.String("results", fmt.Sprint(results)),
			)
		}
	}

	return result
}
