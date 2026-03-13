//go:generate mockgen -source=action.go -destination=mock_action.go -package=ai
package ai

type ActionID string

type Action interface {
	ID() ActionID
	Start(ctx *Context, target Target)
	Update(ctx *Context, target Target) ActionStatus
	Cancel(ctx *Context, target Target)
	ShouldAddToHistory() bool
}
