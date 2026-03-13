//go:generate mockgen -source=consideration.go -destination=mock_consideration.go -package=ai
package ai

type ConsiderationID string

type Curve interface {
	Apply(x float64) float64
}

type Consideration struct {
	ID ConsiderationID

	InputFunc func(ctx *Context, target Target) float64

	MinValue float64
	MaxValue float64

	ResponseCurve Curve

	ShouldCache bool
}

func (c *Consideration) Input(ctx *Context, target Target) float64 {
	if c.ShouldCache {
		if cached, exists := ctx.GetCachedConsideration(c.ID, target); exists {
			return cached
		}

		computed := c.InputFunc(ctx, target)
		ctx.SetCachedConsideration(c.ID, target, computed)
		return computed
	}

	return c.InputFunc(ctx, target)
}

func (c *Consideration) Min() float64 {
	return c.MinValue
}

func (c *Consideration) Max() float64 {
	return c.MaxValue
}

func (c *Consideration) Curve() Curve {
	return c.ResponseCurve
}
