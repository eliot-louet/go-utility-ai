package utility_providers

import "github.com/eliot-louet/go-utility-ai/ai"

type SelfProvider struct{}

var self_target = []ai.Target{}

func (p SelfProvider) Targets(ctx *ai.Context) []ai.Target {
	return self_target
}

func (p SelfProvider) ID() ai.TargetProviderID {
	return "SelfProvider"
}

func (p SelfProvider) ShouldCache() bool {
	return false
}
