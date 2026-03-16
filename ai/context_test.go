package ai

import (
	"testing"
	// Testify
	"github.com/stretchr/testify/assert"
)

type contextFixture struct {
	context *Context
}

func newContextFixture(t *testing.T) *contextFixture {
	t.Helper()

	return &contextFixture{
		context: &Context{
			providerCache:      make(map[TargetProviderID][]Target),
			considerationCache: make(map[considerationCacheKey]float64),
			ActionHistory:      make(map[ActionID]*ActionHistory),
			BehaviorHistory:    make(map[BehaviorID]*BehaviorHistory),
		},
	}
}

func TestCachedTargets(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name string
		// Arrange
		cacheContent map[TargetProviderID][]Target
		setup        func(t *testing.T, f *contextFixture)

		// Act
		fetchKey TargetProviderID
		fetcher  func() []Target

		// Assert
		expectTargets []Target
	}

	tests := []testCase{
		{
			name: "Cache hit",
			cacheContent: map[TargetProviderID][]Target{
				"provider1": {Target{ID: 1, Type: TargetTypeActor}, Target{ID: 2, Type: TargetTypeObject}},
			},
			fetchKey:      "provider1",
			fetcher:       func() []Target { return []Target{Target{ID: 3, Type: TargetTypeActor}} },
			expectTargets: []Target{Target{ID: 1, Type: TargetTypeActor}, Target{ID: 2, Type: TargetTypeObject}},
		},
		{
			name:         "Cache miss",
			cacheContent: map[TargetProviderID][]Target{},
			fetchKey:     "provider1",
			fetcher: func() []Target {
				return []Target{Target{ID: 1, Type: TargetTypeActor}, Target{ID: 2, Type: TargetTypeObject}}
			},
			expectTargets: []Target{Target{ID: 1, Type: TargetTypeActor}, Target{ID: 2, Type: TargetTypeObject}},
		},
		{
			name: "Cache miss with existing different key",
			cacheContent: map[TargetProviderID][]Target{
				"provider1": {Target{ID: 1, Type: TargetTypeActor}, Target{ID: 2, Type: TargetTypeObject}},
			},
			fetchKey:      "provider2",
			fetcher:       func() []Target { return []Target{Target{ID: 3, Type: TargetTypeActor}} },
			expectTargets: []Target{Target{ID: 3, Type: TargetTypeActor}},
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// Arrange
			f := newContextFixture(t)

			f.context.providerCache = tc.cacheContent

			if tc.setup != nil {
				tc.setup(t, f)
			}

			// Act
			targets := f.context.GetCachedTargets(tc.fetchKey, tc.fetcher)

			// Assert
			assert.Equal(t, tc.expectTargets, targets, "Expected targets to match")
		})
	}
}

func TestConsiderationCache(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name string
		// Arrange
		cacheContent map[considerationCacheKey]float64
		setup        func(t *testing.T, f *contextFixture)

		// Act
		fetchId     ConsiderationID
		fetchTarget Target

		// Assert
		expectFloat  float64
		expectExists bool
	}

	tests := []testCase{
		{
			name: "Cache hit",
			cacheContent: map[considerationCacheKey]float64{
				{id: "cons1", target: Target{ID: 1, Type: TargetTypeActor}}: 0.75,
			},
			fetchId:      "cons1",
			fetchTarget:  Target{ID: 1, Type: TargetTypeActor},
			expectFloat:  0.75,
			expectExists: true,
		},
		{
			name: "Cache miss (missing id and target)",
			cacheContent: map[considerationCacheKey]float64{
				{id: "cons1", target: Target{ID: 1, Type: TargetTypeActor}}: 0.75,
			},
			fetchId:      "cons2",
			fetchTarget:  Target{ID: 2, Type: TargetTypeObject},
			expectFloat:  0.0,
			expectExists: false,
		},
		{
			name: "Cache miss (existing id but different target)",
			cacheContent: map[considerationCacheKey]float64{
				{id: "cons1", target: Target{ID: 1, Type: TargetTypeActor}}: 0.75,
			},
			fetchId:      "cons1",
			fetchTarget:  Target{ID: 2, Type: TargetTypeObject},
			expectFloat:  0.0,
			expectExists: false,
		},
		{
			name: "Cache update",
			cacheContent: map[considerationCacheKey]float64{
				{id: "cons1", target: Target{ID: 1, Type: TargetTypeActor}}: 0.75,
			},
			setup: func(t *testing.T, f *contextFixture) {
				f.context.SetCachedConsideration("cons1", Target{ID: 1, Type: TargetTypeActor}, 0.85)
			},
			fetchId:      "cons1",
			fetchTarget:  Target{ID: 1, Type: TargetTypeActor},
			expectFloat:  0.85,
			expectExists: true,
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// Arrange
			f := newContextFixture(t)

			f.context.considerationCache = tc.cacheContent

			if tc.setup != nil {
				tc.setup(t, f)
			}

			// Act
			val, exists := f.context.GetCachedConsideration(tc.fetchId, tc.fetchTarget)

			// Assert
			assert.Equal(t, tc.expectFloat, val, "Expected consideration value to match")
			assert.Equal(t, tc.expectExists, exists, "Expected existence flag to match")
		})
	}
}
