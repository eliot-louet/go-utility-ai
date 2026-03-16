package ai

import "sort"

type BehaviorPackage struct {
	Name      string
	Behaviors []Behavior

	ConditionFunc func(ctx *Context) bool
}

func NewBehaviorPackage(name string, behaviors []Behavior, conditionFunc func(ctx *Context) bool) *BehaviorPackage {
	pkg := &BehaviorPackage{
		Name:          name,
		Behaviors:     behaviors,
		ConditionFunc: conditionFunc,
	}

	// Sort behaviors by max score in descending order for efficient decision making
	// This way, we can quickly skip packages that can't outperform the current behavior
	// and also prioritize higher scoring behaviors within the package
	sort.Slice(pkg.Behaviors, func(i, j int) bool {
		return pkg.Behaviors[i].MaxScore() > pkg.Behaviors[j].MaxScore()
	})

	return pkg
}
