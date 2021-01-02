package evcache

import (
	"errors"
	"github.com/eko/gocache/cache"
	"github.com/eko/gocache/store"
	mocks "github.com/eko/gocache/test/mocks/cache"
	"github.com/go-email-validator/go-email-validator/pkg/ev/evtests"
	"github.com/golang/mock/gomock"
	"reflect"
	"testing"
)

const key = "key"

var (
	simpleValue = map[int]string{0: "123"}
	simpleError = errors.New("simpleError")
)

func TestMain(m *testing.M) {
	evtests.TestMain(m)
}

func Test_gocacheAdapter_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type fields struct {
		cache func() cache.CacheInterface
	}
	type args struct {
		key interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name: "success get",
			fields: fields{
				cache: func() cache.CacheInterface {
					cacheMock := mocks.NewMockCacheInterface(ctrl)
					cacheMock.EXPECT().Get(key).Return(simpleValue, nil).Times(1)

					return cacheMock
				},
			},
			args: args{
				key: key,
			},
			want:    simpleValue,
			wantErr: false,
		},
		{
			name: "error get",
			fields: fields{
				cache: func() cache.CacheInterface {
					cacheMock := mocks.NewMockCacheInterface(ctrl)
					cacheMock.EXPECT().Get(key).Return(nil, simpleError).Times(1)

					return cacheMock
				},
			},
			args: args{
				key: key,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewCache(tt.fields.cache(), nil)
			got, err := c.Get(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_gocacheAdapter_Set(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	emptyOptions := &store.Options{}

	type fields struct {
		cache  func() cache.CacheInterface
		option *store.Options
	}
	type args struct {
		key    interface{}
		object interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				cache: func() cache.CacheInterface {
					cacheMock := mocks.NewMockCacheInterface(ctrl)
					cacheMock.EXPECT().Set(key, simpleValue, emptyOptions).Return(nil).Times(1)

					return cacheMock
				},
				option: emptyOptions,
			},
			args: args{
				key:    key,
				object: simpleValue,
			},
			wantErr: false,
		},
		{
			name: "error",
			fields: fields{
				cache: func() cache.CacheInterface {
					cacheMock := mocks.NewMockCacheInterface(ctrl)
					cacheMock.EXPECT().Set(key, simpleValue, emptyOptions).Return(simpleError).Times(1)

					return cacheMock
				},
				option: emptyOptions,
			},
			args: args{
				key:    key,
				object: simpleValue,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewCache(tt.fields.cache(), tt.fields.option)
			if err := c.Set(tt.args.key, tt.args.object); (err != nil) != tt.wantErr {
				t.Errorf("Set() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
