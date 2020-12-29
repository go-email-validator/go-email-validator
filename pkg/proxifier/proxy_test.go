package proxifier

import (
	"errors"
	"testing"
)

const (
	addressFirst  = "addressFirst"
	addressSecond = "addressSecond"
)

func getTwoAddrs(t *testing.T) []*Address {
	twoAddresses, errs := getAddressesFromString([]string{addressFirst, addressSecond})
	if len(errs) > 0 {
		t.Error(errs)
	}

	return twoAddresses
}

func Test_proxyList_GetAddress(t *testing.T) {
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
				pool:          getTwoAddrs(t),
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
				pool:          getTwoAddrs(t),
				using:         setMapFromList(getTwoAddrs(t), newMap()),
				addressGetter: GetLastAddress,
			},
			want:    addressSecond,
			wantErr: nil,
		},
		{
			name: "bulk test less length",
			fields: fields{
				bulkPool:      1,
				pool:          getTwoAddrs(t),
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
				bulkPool:      len(getTwoAddrs(t)) + 1,
				pool:          getTwoAddrs(t),
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
				banned:        setMapFromList(getTwoAddrs(t), nil),
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
				banned:        setMapFromList(getTwoAddrs(t), nil),
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
				banned:        setMapFromList(getTwoAddrs(t), nil),
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
					return getTwoAddrs(t)
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
