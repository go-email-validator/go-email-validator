package proxy

import (
	"errors"
	"math/rand"
	"time"
)

const PrefixError = "ev-proxy: "

var ErrEmptyPool = errors.New(PrefixError + "proxy pool is empty")
var ErrNotUsingAddress = errors.New(PrefixError + "Address is not used")

type Address struct {
	url  string
	used uint
	ban  bool
}

type GetAddress func(map[string]*Address, []string) *Address

type ProxyList interface {
	GetAddress() (string, error)
}

type ProxyListDTO struct {
	Addresses     []string
	BulkPool      int
	RandomAddress GetAddress
	ActivePool    int
}

func NewProxyListFromStrings(dto ProxyListDTO) (list ProxyList, errs map[string]error) {
	var addrs = make([]*Address, 0)
	errs = make(map[string]error, 0)

	for _, addr := range dto.Addresses {
		addrs = append(addrs, &Address{
			url: addr,
		})
	}

	if dto.RandomAddress == nil {
		dto.RandomAddress = GetRandomAddress
	}

	return &proxyList{
		activePool:    dto.ActivePool,
		bulkPool:      dto.BulkPool,
		indexPool:     0,
		pool:          addrs,
		randomAddress: dto.RandomAddress,
	}, errs
}

func GetRandomAddress(Addrs map[string]*Address, addrs []string) *Address {
	rand.Seed(time.Now().UnixNano())
	return Addrs[addrs[rand.Intn(len(addrs))]]
}

func CreateCircleAddress(i int) GetAddress {
	return func(m map[string]*Address, addrs []string) *Address {
		if i >= len(addrs) {
			i = 0
		}
		i++
		return m[addrs[i-1]]
	}
}

type MapAddress map[string]*Address

type proxyList struct {
	activePool          int
	bulkPool            int // count of new getting proxies
	indexPool           int
	pool                []*Address
	using               MapAddress
	usingKeys           []string
	ban                 MapAddress
	banRecovering       int
	requestNewAddresses func() []*Address
	randomAddress       GetAddress
}

func (p *proxyList) GetAddress() (string, error) {
	if len(p.using) <= p.activePool {
		poolLen := len(p.pool)
		if poolLen <= p.indexPool && p.requestNewAddresses != nil {
			p.pool = append(p.pool, p.requestNewAddresses()...)
		}

		hasInPoll := poolLen > p.indexPool
		infiniteBanReuse := p.banRecovering == -1
		hasBanRecoveryAttempt := p.banRecovering > 0
		updateKeys := false
		if hasInPoll {
			var nextBulkPool int
			if p.bulkPool == 0 {
				nextBulkPool = len(p.pool)
			} else {
				nextBulkPool = p.indexPool + p.bulkPool
			}

			if poolLen > nextBulkPool {
				nextBulkPool = poolLen
			}

			using := make(MapAddress, nextBulkPool-p.indexPool)
			for _, addr := range p.pool[p.indexPool:nextBulkPool] {
				using[addr.url] = addr
			}
			p.using = mergeAddress(p.using, using)
			p.indexPool = nextBulkPool

			updateKeys = true
		} else if infiniteBanReuse || hasBanRecoveryAttempt {
			if hasBanRecoveryAttempt {
				p.banRecovering--
			}
			p.using = mergeAddress(p.using, p.ban)
			updateKeys = true
		}

		if updateKeys {
			p.usingKeys = make([]string, len(p.using))
			i := 0
			for _, addr := range p.using {
				p.usingKeys[i] = addr.url
				i++
			}
		}
	}

	if len(p.using) == 0 {
		return "", ErrEmptyPool
	}

	addr := p.randomAddress(p.using, p.usingKeys)
	addr.used++

	return addr.url, nil
}

func (p *proxyList) Ban(addrKey string) bool {
	if _, hasKey := p.using[addrKey]; !hasKey {
		return false
	}

	p.ban[addrKey] = p.using[addrKey]
	delete(p.using, addrKey)
	p.ban[addrKey].ban = true

	return true
}

func (p *proxyList) GetNew(banAddress string) (string, error) {
	if !p.Ban(banAddress) {
		return "", ErrNotUsingAddress
	}

	return p.GetAddress()
}

func mergeAddress(addrsSource MapAddress, addrsExt MapAddress) MapAddress {
	if addrsSource == nil || len(addrsSource) == 0 {
		return addrsExt
	}

	for key, addrExt := range addrsExt {
		addrsSource[key] = addrExt
	}

	return addrsSource
}
