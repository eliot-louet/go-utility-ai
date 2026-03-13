package provider

import "github.com/eliot-louet/go-utility-ai/ai"

type ShoutProvider struct{}

var shout_target = []ai.Target{
	"I am the strongest!",
	"You can't defeat me!",
	"Feel my wrath!",
	"Is that all you've got?",
}

func (p ShoutProvider) Targets(ctx *ai.Context) []ai.Target {
	return shout_target
}

func (p ShoutProvider) ID() ai.TargetProviderID {
	return "ShoutProvider"
}

func (p ShoutProvider) ShouldCache() bool {
	return false
}
