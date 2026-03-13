package actions

import "github.com/eliot-louet/go-utility-ai/internal/ai"

type WaitAction struct {
}

func NewWaitAction() *WaitAction {
	return &WaitAction{}
}

func (a *WaitAction) ID() ai.ActionID {
	return "WaitAction"
}

func (a *WaitAction) Start(ctx *ai.Context, target ai.Target) {
}

func (a *WaitAction) Update(ctx *ai.Context, target ai.Target) ai.ActionStatus {
	return ai.Success
}

func (a *WaitAction) Cancel(ctx *ai.Context, target ai.Target) {
}

func (a *WaitAction) ShouldAddToHistory() bool {
	return false
}
