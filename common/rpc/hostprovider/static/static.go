package static

import (
	"github.com/bytom/btmcpool/common/rpc/hostprovider/types"
)

type provider struct {
	hosts   []string
	service string
}

// NewClient creates a new service discovery client with static host list
func NewProvider(service string, hosts []string) types.Provider {
	return &provider{
		hosts:   hosts,
		service: service,
	}
}

// Get is the implementations of provider interface
func (p *provider) Get() ([]string, error) {
	return p.hosts, nil
}

func (p *provider) Ensure() {

}

func (p *provider) Name() string {
	return p.service
}
