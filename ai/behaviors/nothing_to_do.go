package behaviors

import (
	"github.com/eliot-louet/go-utility-ai/ai"
	"github.com/eliot-louet/go-utility-ai/ai/considerations"
	"github.com/eliot-louet/go-utility-ai/ai/provider"
)

type NothingToDo struct {
	considerations []*ai.Consideration
	provider       ai.TargetProvider
	action         ai.Action
}

func NewNothingToDo(cCache *ai.ConsiderationCache, action ai.Action) *NothingToDo {
	return &NothingToDo{
		considerations: []*ai.Consideration{
			cCache.Get(considerations.InverseBoredomConsiderationID, considerations.OwnInverseBoredom),
		},
		provider: provider.SelfProvider{},
		action:   action,
	}
}

func (a *NothingToDo) ShouldAddToHistory() bool {
	return false
}

func (a *NothingToDo) ID() ai.BehaviorID {
	return "NothingToDo"
}

func (a *NothingToDo) Name() string { return "Nothing To Do" }

func (a *NothingToDo) Considerations(ctx *ai.Context, target ai.Target) []*ai.Consideration {
	return a.considerations
}

func (a *NothingToDo) Weight(ctx *ai.Context, target ai.Target) float64 { return 1 }

func (a *NothingToDo) Provider(ctx *ai.Context) ai.TargetProvider {
	return a.provider
}

func (a *NothingToDo) Action(ctx *ai.Context, target ai.Target) ai.Action {
	return a.action
}
