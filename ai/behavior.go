//go:generate mockgen -source=behavior.go -destination=mock_behavior.go -package=ai
package ai

import (
	"math"

	"github.com/eliot-louet/go-utility-ai/ai/math_utils"
)

type ActionStatus int
type BehaviorID string
type Target interface{}

const (
	Running ActionStatus = iota
	Success
	Failure
)

type Behavior interface {
	ID() BehaviorID
	Name() string

	Considerations(ctx *Context, target Target) []*Consideration

	Weight(ctx *Context, target Target) float64

	Provider(ctx *Context) TargetProvider

	Action(ctx *Context, target Target) Action

	ShouldAddToHistory() bool
}

type RunningBehavior struct {
	Action   Action
	Behavior Behavior
	Target   Target
	Score    float64
	State    map[string]interface{} // Private state for this specific action instance
	Running  bool
}

func EvaluateBehavior(ctx *Context, behavior Behavior, target Target) float64 {
	score_product := 1.0
	values := 0.0

	for _, c := range behavior.Considerations(ctx, target) {
		raw := c.Input(ctx, target)

		normalized := math_utils.Normalize(raw, c.Min(), c.Max())

		curved := c.Curve().Apply(normalized)

		// If any consideration returns 0, the whole behavior is not worth it
		if curved == 0 {
			return 0
		}

		score_product *= curved
		values++
	}

	final := math.Pow(score_product, 1.0/values)

	return final * behavior.Weight(ctx, target)
}
