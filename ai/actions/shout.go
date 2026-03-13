package actions

import "github.com/eliot-louet/go-utility-ai/ai"

type ShoutAction struct {
}

func NewShoutAction() *ShoutAction {
	return &ShoutAction{}
}

func (a *ShoutAction) ID() ai.ActionID {
	return "ShoutAction"
}

func (a *ShoutAction) Start(ctx *ai.Context, target ai.Target) {
}

func (a *ShoutAction) Update(ctx *ai.Context, target ai.Target) ai.ActionStatus {
	ctx.Self.OffsetValue("boredom", -40)

	println("Shouting out of boredom! %s", target.(string))
	return ai.Success
}

func (a *ShoutAction) Cancel(ctx *ai.Context, target ai.Target) {
}

func (a *ShoutAction) ShouldAddToHistory() bool {
	return true
}
