package proxy_list

import (
	"errors"
	"fmt"
	"github.com/emirpasic/gods/maps"
	"github.com/emirpasic/gods/maps/linkedhashmap"
	"net/url"
)

const (
	PrefixError      = "ev-proxy: "
	InvalidAddr      = PrefixError + "addr %v: %w"
	InfiniteRecovery = -1
	EmptyAddress     = ""
)

var ErrEmptyPool = errors.New(PrefixError + "proxy pool is empty")
var ErrActivePool = errors.New(PrefixError + "proxy pool should be more or equal of ProxyListDTO.ActivePool")

type ProxyList interface {
	GetAddress() (string, error)
	Ban(string) bool
}

type ProxyListDTO struct {
	Addresses     []string
	BulkPool      int
	AddressGetter GetAddress
	ActivePool    int
}

func NewProxyListFromStrings(dto ProxyListDTO) (list ProxyList, errs []error) {
	var addrs = make([]*Address, 0)
	errs = make([]error, 0)

	addrs, errs = getAddressesFromString(dto.Addresses)

	if dto.AddressGetter == nil {
		dto.AddressGetter = GetRandomAddress
	}

	if dto.ActivePool > len(addrs) {
		errs = append(errs, ErrActivePool)
		return nil, errs
	}

	return &proxyList{
		minUsing:      dto.ActivePool,
		bulkPool:      dto.BulkPool,
		indexPool:     0,
		pool:          addrs,
		addressGetter: dto.AddressGetter,
		using:         newMap(),
		banned:        newMap(),
	}, errs
}

func getAddressesFromString(addrStrings []string) (addrs []*Address, errs []error) {
	for _, addr := range addrStrings {
		_, err := url.Parse(addr)
		if err != nil {
			errs = append(errs, fmt.Errorf(InvalidAddr, addr, err))
			continue
		}

		addrs = append(addrs, &Address{
			url: addr,
		})
	}

	return
}

func setMapFromList(addrs []*Address, m MapAddress) MapAddress {
	if m == nil {
		m = newMap()
	}

	for _, addr := range addrs {
		m.Put(addr.url, addr)
	}

	return m
}

func newMap() MapAddress {
	return linkedhashmap.New()
}

type MapAddress maps.Map // linkedhashmap.Map

// TODO add strategy struct for changing behavior
type proxyList struct {
	bulkPool            int // count of new getting proxies
	indexPool           int
	pool                []*Address
	using               MapAddress
	minUsing            int // minimal count of using address in one time
	banned              MapAddress
	banRecovering       int
	requestNewAddresses func() []*Address
	addressGetter       GetAddress
}

func (p *proxyList) GetAddress() (string, error) {
	if p.needMore() {
		if p.shouldGetNewAddresses() {
			p.getNewAddresses()
		}

		if p.hasUnusedInPoll() {
			var nextBulkPool = p.indexPool + p.bulkPool

			if p.bulkPool == 0 || len(p.pool) < nextBulkPool {
				nextBulkPool = len(p.pool)
			}

			usingExtending := linkedhashmap.New()
			setMapFromList(p.pool[p.indexPool:nextBulkPool], usingExtending)
			p.using = mergeAddress(p.using, usingExtending)
			p.indexPool = nextBulkPool
		} else if p.canRecoveryBan() {
			if p.hasAttempts() {
				p.banRecovering--
			}
			p.using = mergeAddress(p.using, p.banned)
		}
	}

	if p.using.Size() == 0 {
		return EmptyAddress, ErrEmptyPool
	}

	addr := p.addressGetter(p.using, p.using.Keys())
	addr.used++

	return addr.url, nil
}

func (p *proxyList) needMore() bool {
	return p.using.Size() < p.minUsing || p.using.Size() == 0
}

func (p *proxyList) shouldGetNewAddresses() bool {
	return len(p.pool) <= p.indexPool && p.requestNewAddresses != nil
}

func (p *proxyList) getNewAddresses() {
	p.pool = append(p.pool, p.requestNewAddresses()...)
}

func (p *proxyList) hasUnusedInPoll() bool {
	return len(p.pool) > p.indexPool
}

func (p *proxyList) canRecoveryBan() bool {
	return p.hasAttempts() || p.banRecovering == InfiniteRecovery
}

func (p *proxyList) hasAttempts() bool {
	return p.banRecovering > 0
}

func (p *proxyList) Ban(addrKey string) bool {
	if _, hasKey := p.using.Get(addrKey); !hasKey {
		return false
	}

	addr, _ := p.using.Get(addrKey)
	p.banned.Put(addrKey, addr)
	p.using.Remove(addrKey)
	addr.(*Address).ban = true

	return true
}

func mergeAddress(addrsSource MapAddress, addrsExt MapAddress) MapAddress {
	if addrsSource == nil || addrsSource.Size() == 0 {
		return addrsExt
	}

	for key := range addrsExt.Keys() {
		addrExt, _ := addrsExt.Get(key)
		addrsSource.Put(key, addrExt)
	}

	return addrsSource
}
