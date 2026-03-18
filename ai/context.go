package ai

type considerationCacheKey struct {
	id     ConsiderationID
	target Target
}

type Context struct {
	Environment        Environment
	Self               Actor
	providerCache      map[TargetProviderID][]Target
	considerationCache map[considerationCacheKey]float64
	ActionHistory      map[ActionID]*ActionHistory
	BehaviorHistory    map[BehaviorID]*BehaviorHistory
}

func MakeContext(environment Environment, self Actor) *Context {
	return &Context{
		Environment:        environment,
		Self:               self,
		providerCache:      make(map[TargetProviderID][]Target),
		considerationCache: make(map[considerationCacheKey]float64),
		ActionHistory:      make(map[ActionID]*ActionHistory),
		BehaviorHistory:    make(map[BehaviorID]*BehaviorHistory),
	}
}

// GetCachedTargets retrieves cached targets by key, or evaluates them using the fallback function if missing.
func (c *Context) GetCachedTargets(key TargetProviderID, fetcher func() []Target) []Target {
	// Return cached targets if they exist
	if cached, exists := c.providerCache[key]; exists {
		return cached
	}

	// Compute, cache, and return targets if not cached
	computed := fetcher()
	c.providerCache[key] = computed

	return computed
}

// GetCachedConsideration retrieves a consideration safely.
func (c *Context) GetCachedConsideration(id ConsiderationID, target Target) (float64, bool) {
	val, exists := c.considerationCache[considerationCacheKey{id: id, target: target}]

	return val, exists
}

// SetCachedConsideration stores a consideration result.
func (c *Context) SetCachedConsideration(id ConsiderationID, target Target, value float64) {
	c.considerationCache[considerationCacheKey{id: id, target: target}] = value
}
