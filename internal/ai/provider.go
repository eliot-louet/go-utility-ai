//go:generate mockgen -source=provider.go -destination=mock_provider.go -package=ai

package ai

type TargetProviderID string

type TargetProvider interface {
	ID() TargetProviderID
	Targets(ctx *Context) []Target
	ShouldCache() bool
}
