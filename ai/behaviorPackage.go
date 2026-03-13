package ai

type BehaviorPackage struct {
	Name      string
	Behaviors []Behavior

	ConditionFunc func(ctx *Context) bool
}
