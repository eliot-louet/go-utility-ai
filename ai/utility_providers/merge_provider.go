package utility_providers

import (
	"strings"

	"github.com/eliot-louet/go-utility-ai/ai"
)

type MergeProvider struct {
	providers       []ai.TargetProvider
	CachedID        ai.TargetProviderID
	MergeDuplicates bool

	buf  []ai.Target
	seen map[ai.Target]struct{}
}

func NewMergeProvider(mergeDuplicates bool, providers ...ai.TargetProvider) *MergeProvider {
	var providerIDs []string
	for _, provider := range providers {
		providerIDs = append(providerIDs, string(provider.ID()))
	}

	return &MergeProvider{
		providers:       providers,
		MergeDuplicates: mergeDuplicates,
		CachedID:        ai.TargetProviderID("merge:" + strings.Join(providerIDs, ",")),
		seen:            make(map[ai.Target]struct{}),
	}
}

func (p *MergeProvider) Targets(ctx *ai.Context) []ai.Target {
	p.buf = p.buf[:0]

	if p.MergeDuplicates {
		clear(p.seen)
	}

	for _, provider := range p.providers {
		for _, target := range provider.Targets(ctx) {

			if p.MergeDuplicates {
				if _, exists := p.seen[target]; exists {
					continue
				}
				p.seen[target] = struct{}{}
			}

			p.buf = append(p.buf, target)
		}
	}

	return p.buf
}

func (p *MergeProvider) ID() ai.TargetProviderID {
	return p.CachedID
}

func (p *MergeProvider) ShouldCache() bool {
	return false
}
