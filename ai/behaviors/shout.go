package behaviors

import (
	"github.com/eliot-louet/go-utility-ai/ai"
	"github.com/eliot-louet/go-utility-ai/ai/considerations"
	"github.com/eliot-louet/go-utility-ai/ai/provider"
)

type ShoutOutOfBoredom struct {
	considerations []*ai.Consideration
	provider       ai.TargetProvider
	action         ai.Action
}

func NewShoutOutOfBoredom(cCache *ai.ConsiderationCache, action ai.Action) *ShoutOutOfBoredom {
	return &ShoutOutOfBoredom{
		considerations: []*ai.Consideration{
			cCache.Get(considerations.BoredomConsiderationID, considerations.OwnBoredom),
		},
		provider: provider.ShoutProvider{},
		action:   action,
	}
}

func (a *ShoutOutOfBoredom) ID() ai.BehaviorID {
	return "ShoutOutOfBoredom"
}

func (a *ShoutOutOfBoredom) ShouldAddToHistory() bool {
	return true
}

func (a *ShoutOutOfBoredom) Name() string {
	return "Shout Out Of Boredom"
}

func (a *ShoutOutOfBoredom) Provider(ctx *ai.Context) ai.TargetProvider {
	return a.provider
}

func (a *ShoutOutOfBoredom) Considerations(ctx *ai.Context, target ai.Target) []*ai.Consideration {
	return a.considerations
}

func (a *ShoutOutOfBoredom) Weight(ctx *ai.Context, target ai.Target) float64 {
	return 1.0
}

func (a *ShoutOutOfBoredom) Action(ctx *ai.Context, target ai.Target) ai.Action {
	return a.action
}
