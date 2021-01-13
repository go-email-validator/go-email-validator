package ev

import (
	"errors"
	"github.com/allegro/bigcache"
	"github.com/eko/gocache/marshaler"
	"github.com/eko/gocache/store"
	"github.com/go-email-validator/go-email-validator/pkg/ev/evcache"
	"github.com/go-email-validator/go-email-validator/pkg/ev/evmail"
	"github.com/go-email-validator/go-email-validator/pkg/ev/evsmtp"
	mockevcache "github.com/go-email-validator/go-email-validator/test/mock/ev/evcache"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"net/textproto"
	"reflect"
	"testing"
	"time"
)

func Test_cacheDecorator_Validate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	results := make([]ValidationResult, 0)
	key := validEmail.String()

	type fields struct {
		validator Validator
		cache     func() evcache.Interface
		getKey    CacheKeyGetter
	}
	type args struct {
		email   evmail.Address
		results []ValidationResult
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantResult ValidationResult
	}{
		{
			name: "without cache, error and with set error",
			fields: fields{
				validator: inValidMockValidator,
				cache: func() evcache.Interface {
					cacheMock := mockevcache.NewMockInterface(ctrl)
					cacheMock.EXPECT().Get(key).Return(nil, nil).Times(1)
					cacheMock.EXPECT().Set(key, invalidResult).Return(errorSimple).Times(1)

					return cacheMock
				},
				getKey: EmailCacheKeyGetter,
			},
			args: args{
				email:   validEmail,
				results: results,
			},
			wantResult: invalidResult,
		},
		{
			name: "without cache and with get error",
			fields: fields{
				validator: validMockValidator,
				cache: func() evcache.Interface {
					cacheMock := mockevcache.NewMockInterface(ctrl)
					cacheMock.EXPECT().Get(key).Return(nil, errorSimple).Times(1)
					cacheMock.EXPECT().Set(key, validResult).Return(nil).Times(1)

					return cacheMock
				},
				getKey: nil,
			},
			args: args{
				email:   validEmail,
				results: results,
			},
			wantResult: validResult,
		},
		{
			name: "with cache",
			fields: fields{
				validator: validMockValidator,
				cache: func() evcache.Interface {
					cacheMock := mockevcache.NewMockInterface(ctrl)
					cacheMock.EXPECT().Get(key).Return(&validResult, nil).Times(1)

					return cacheMock
				},
				getKey: EmailCacheKeyGetter,
			},
			args: args{
				email:   validEmail,
				results: nil,
			},
			wantResult: validResult,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewCacheDecorator(tt.fields.validator, tt.fields.cache(), tt.fields.getKey)
			if gotResult := c.Validate(tt.args.email, tt.args.results...); !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("Validate() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func Test_cacheDecorator_GetDeps(t *testing.T) {
	deps := []ValidatorName{OtherValidator}

	type fields struct {
		validator Validator
		cache     evcache.Interface
		getKey    CacheKeyGetter
	}
	tests := []struct {
		name   string
		fields fields
		want   []ValidatorName
	}{
		{
			name: "return deps",
			fields: fields{
				validator: mockValidator{deps: deps},
			},
			want: deps,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &cacheDecorator{
				validator: tt.fields.validator,
				cache:     tt.fields.cache,
				getKey:    tt.fields.getKey,
			}
			if got := c.GetDeps(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetDeps() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEmailCacheKeyGetter(t *testing.T) {
	type args struct {
		email   evmail.Address
		results []ValidationResult
	}
	tests := []struct {
		name string
		args args
		want interface{}
	}{
		{
			name: "success",
			args: args{
				email:   validEmail,
				results: nil,
			},
			want: validEmail.String(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EmailCacheKeyGetter(tt.args.email, tt.args.results...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EmailCacheKeyGetter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDomainCacheKeyGetter(t *testing.T) {
	type args struct {
		email evmail.Address
	}
	tests := []struct {
		name string
		args args
		want interface{}
	}{
		{
			name: "success",
			args: args{
				email: validEmail,
			},
			want: validEmail.Domain(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DomainCacheKeyGetter(tt.args.email); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DomainCacheKeyGetter() = %v, want %v", got, tt.want)
			}
		})
	}
}

type customErr struct{}

func (customErr) Error() string {
	return "customErr"
}

var cacheErrs = []error{
	//error(&customErr{}), TODO find way to marshal and unmarshal all interfaces
	NewDepsError(),
	evsmtp.NewError(1, &textproto.Error{Code: 505, Msg: "msg1"}),
	evsmtp.NewError(1, errors.New("msg2")),
}
var validatorResult = NewResult(true, cacheErrs, cacheErrs, OtherValidator)

func Test_Cache(t *testing.T) {
	bigCacheClient, err := bigcache.NewBigCache(bigcache.DefaultConfig(1 * time.Second))
	require.Nil(t, err)
	bigCacheStore := store.NewBigcache(bigCacheClient, nil)

	marshal := marshaler.New(bigCacheStore)

	cache := evcache.NewCacheMarshaller(marshal, func() interface{} {
		return new(ValidationResult)
	}, nil)

	key := "key"

	err = cache.Set(key, validatorResult)
	require.Nil(t, err)

	gotInterface, err := cache.Get(key)
	require.Nil(t, err)
	got := *gotInterface.(*ValidationResult)
	require.Equal(t, validatorResult, got)
}
