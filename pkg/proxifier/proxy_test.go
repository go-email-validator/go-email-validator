package proxifier

import (
	"errors"
	"fmt"
	"github.com/go-email-validator/go-email-validator/pkg/ev/evtests"
	"reflect"
	"testing"
)

func TestMain(m *testing.M) {
	evtests.TestMain(m)
}

func Test_list_GetAddress(t *testing.T) {
	type fields struct {
		minUsing            int
		bulkPool            int
		indexPool           int
		pool                []*Address
		using               MapAddress
		banned              MapAddress
		banRecovering       int
		requestNewAddresses func() []*Address
		addressGetter       GetAddress
	}

	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr error
	}{
		{
			name: "active pool",
			fields: fields{
				minUsing:      1,
				pool:          getAddrsTest(t, getTestAddrsStr()),
				addressGetter: GetLastAddress,
				using:         newMap(),
				banned:        newMap(),
			},
			want:    addressSecond,
			wantErr: nil,
		},
		{
			name: "active pool with using",
			fields: fields{
				minUsing:      1,
				pool:          getAddrsTest(t, getTestAddrsStr()),
				using:         setMapFromList(getAddrsTest(t, getTestAddrsStr()), newMap()),
				addressGetter: GetLastAddress,
			},
			want:    addressSecond,
			wantErr: nil,
		},
		{
			name: "bulk test less length",
			fields: fields{
				bulkPool:      1,
				pool:          getAddrsTest(t, getTestAddrsStr()),
				addressGetter: GetLastAddress,
				using:         newMap(),
				banned:        newMap(),
			},
			want:    addressFirst,
			wantErr: nil,
		},
		{
			name: "bulk test more length",
			fields: fields{
				bulkPool:      len(getAddrsTest(t, getTestAddrsStr())) + 1,
				pool:          getAddrsTest(t, getTestAddrsStr()),
				addressGetter: GetLastAddress,
				using:         newMap(),
				banned:        newMap(),
			},
			want:    addressSecond,
			wantErr: nil,
		},
		{
			name: "infinite ban recovery",
			fields: fields{
				pool:          nil,
				banRecovering: InfiniteRecovery,
				banned:        setMapFromList(getAddrsTest(t, getTestAddrsStr()), nil),
				addressGetter: GetLastAddress,
				using:         newMap(),
			},
			want:    addressSecond,
			wantErr: nil,
		},
		{
			name: "ban recovery success attempt",
			fields: fields{
				pool:          nil,
				banRecovering: 1,
				banned:        setMapFromList(getAddrsTest(t, getTestAddrsStr()), nil),
				addressGetter: GetFirstAddress,
				using:         newMap(),
			},
			want:    addressFirst,
			wantErr: nil,
		},
		{
			name: "ban recovery fail attempt",
			fields: fields{
				pool:          nil,
				banRecovering: 0,
				banned:        setMapFromList(getAddrsTest(t, getTestAddrsStr()), nil),
				addressGetter: GetLastAddress,
				using:         newMap(),
			},
			want:    EmptyAddress,
			wantErr: ErrEmptyPool,
		},
		{
			name: "requestNewAddresses success",
			fields: fields{
				pool: nil,
				requestNewAddresses: func() []*Address {
					return getAddrsTest(t, getTestAddrsStr())
				},
				addressGetter: GetLastAddress,
				using:         newMap(),
				banned:        newMap(),
			},
			want:    addressSecond,
			wantErr: nil,
		},
		{
			name: "requestNewAddresses empty result",
			fields: fields{
				pool: nil,
				requestNewAddresses: func() []*Address {
					return []*Address{}
				},
				addressGetter: GetLastAddress,
				using:         newMap(),
				banned:        newMap(),
			},
			want:    EmptyAddress,
			wantErr: ErrEmptyPool,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &list{
				minUsing:            tt.fields.minUsing,
				bulkPool:            tt.fields.bulkPool,
				indexPool:           tt.fields.indexPool,
				pool:                tt.fields.pool,
				using:               tt.fields.using,
				banned:              tt.fields.banned,
				banRecovering:       tt.fields.banRecovering,
				requestNewAddresses: tt.fields.requestNewAddresses,
				addressGetter:       tt.fields.addressGetter,
			}

			got, err := p.GetAddress()
			if (err != nil) && !errors.Is(err, tt.wantErr) {
				t.Errorf("GetAddress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetAddress() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewListFromStrings(t *testing.T) {
	type args struct {
		dto ListDTO
	}
	tests := []struct {
		name     string
		args     args
		wantLst  List
		wantErrs []error
	}{
		{
			name: "without Error",
			args: args{
				dto: ListDTO{
					Addresses: getTestAddrsStr(),
					MinUsing:  2,
					BulkPool:  1,
				},
			},
			wantLst: &list{
				minUsing:  2,
				bulkPool:  1,
				indexPool: 0,
				pool:      getAddrsTest(t, getTestAddrsStr()),
				using:     newMap(),
				banned:    newMap(),
			},
			wantErrs: nil,
		},
		{
			name: "with only address Error",
			args: args{
				dto: ListDTO{
					Addresses: []string{addressInvalid, addressInvalid, addressFirst},
					BulkPool:  3,
					MinUsing:  1,
				},
			},
			wantLst: &list{
				minUsing:  1,
				bulkPool:  3,
				indexPool: 0,
				pool:      getAddrsTest(t, []string{addressFirst}),
				using:     newMap(),
				banned:    newMap(),
			},
			wantErrs: append(getAddrErrs([]string{addressInvalid, addressInvalid, addressFirst})),
		},
		{
			name: "with Error",
			args: args{
				dto: ListDTO{
					Addresses: []string{addressInvalid, addressInvalid, addressFirst},
					BulkPool:  3,
					MinUsing:  5,
				},
			},
			wantLst:  nil,
			wantErrs: append(getAddrErrs([]string{addressInvalid, addressInvalid, addressFirst}), ErrNotEnough),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotLst, gotErrs := NewListFromStrings(tt.args.dto)

			if gotLst != nil {
				got := gotLst.(*list)
				// TODO fix comparing of function
				got.requestNewAddresses = nil
				got.addressGetter = nil
			}

			fmt.Printf("%#v", gotErrs)

			if !reflect.DeepEqual(gotLst, tt.wantLst) {
				t.Errorf("NewListFromStrings() gotLst = %v, want %v", gotLst, tt.wantLst)
			}
			if !reflect.DeepEqual(gotErrs, tt.wantErrs) {
				t.Errorf("NewListFromStrings() gotErrs = %v, want %v", gotErrs, tt.wantErrs)
			}
		})
	}
}

type addressValue struct {
	key   interface{}
	value interface{}
}

func mapAddress(values ...addressValue) MapAddress {
	m := newMap()

	for _, value := range values {
		m.Put(value.key, value.value)
	}

	return m
}

func Test_list_Ban(t *testing.T) {
	type fields struct {
		bulkPool            int
		indexPool           int
		pool                []*Address
		using               MapAddress
		minUsing            int
		banned              MapAddress
		banRecovering       int
		requestNewAddresses func() []*Address
		addressGetter       GetAddress
	}
	type args struct {
		addrKey string
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		want     bool
		wantList *list
	}{
		{
			name: "first banned",
			fields: fields{
				using: mapAddress(addressValue{
					key: "key1",
					value: &Address{
						url:  "key1",
						used: 0,
						ban:  false,
					},
				}),
				banned: newMap(),
			},
			args: args{
				addrKey: "key1",
			},
			want: true,
			wantList: &list{
				using: newMap(),
				banned: mapAddress(addressValue{
					key: "key1",
					value: &Address{
						url:  "key1",
						used: 0,
						ban:  true,
					},
				}),
			},
		},
		{
			name: "second banned",
			fields: fields{
				using: mapAddress(addressValue{
					key: "key2",
					value: &Address{
						url:  "key2",
						used: 0,
						ban:  true,
					},
				}),
				banned: mapAddress(addressValue{
					key: "key1",
					value: &Address{
						url:  "key1",
						used: 0,
						ban:  true,
					},
				}),
			},
			args: args{
				addrKey: "key2",
			},
			want: true,
			wantList: &list{
				using: newMap(),
				banned: mapAddress(
					addressValue{
						key: "key1",
						value: &Address{
							url:  "key1",
							used: 0,
							ban:  true,
						},
					},
					addressValue{
						key: "key2",
						value: &Address{
							url:  "key2",
							used: 0,
							ban:  true,
						},
					},
				),
			},
		},
		{
			name: "missing key for ban",
			fields: fields{
				using:  mapAddress(),
				banned: mapAddress(),
			},
			args: args{
				addrKey: "missing key for ban",
			},
			want: false,
			wantList: &list{
				using:  newMap(),
				banned: mapAddress(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &list{
				bulkPool:            tt.fields.bulkPool,
				indexPool:           tt.fields.indexPool,
				pool:                tt.fields.pool,
				using:               tt.fields.using,
				minUsing:            tt.fields.minUsing,
				banned:              tt.fields.banned,
				banRecovering:       tt.fields.banRecovering,
				requestNewAddresses: tt.fields.requestNewAddresses,
				addressGetter:       tt.fields.addressGetter,
			}

			if got := p.Ban(tt.args.addrKey); got != tt.want {
				t.Errorf("Ban() = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(p, tt.wantList) {
				t.Errorf("list = %v, want %v", p, tt.wantList)
			}
		})
	}
}

func Test_mergeAddress(t *testing.T) {
	type args struct {
		addrsSource MapAddress
		addrsExt    MapAddress
	}
	tests := []struct {
		name string
		args args
		want MapAddress
	}{
		{
			name: "addrsSource - nil",
			args: args{
				addrsSource: nil,
				addrsExt:    mapAddress(addressValue{}),
			},
			want: mapAddress(addressValue{}),
		},
		{
			name: "addrsSource - empty",
			args: args{
				addrsSource: mapAddress(),
				addrsExt:    mapAddress(addressValue{}),
			},
			want: mapAddress(addressValue{}),
		},
		{
			name: "addrsSource - extend",
			args: args{
				addrsSource: mapAddress(addressValue{key: "key1", value: "value1"}),
				addrsExt:    mapAddress(addressValue{key: "key2", value: "value2"}),
			},
			want: mapAddress(addressValue{key: "key1", value: "value1"}, addressValue{key: "key2", value: "value2"}),
		},
		{
			name: "addrsSource - rewrite",
			args: args{
				addrsSource: mapAddress(addressValue{key: "key1", value: "value1"}),
				addrsExt:    mapAddress(addressValue{key: "key1", value: "value2"}),
			},
			want: mapAddress(addressValue{key: "key1", value: "value2"}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mergeAddress(tt.args.addrsSource, tt.args.addrsExt); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("mergeAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}
