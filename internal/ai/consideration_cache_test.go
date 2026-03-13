package ai

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type considerationCacheTest struct {
	considerationCache *ConsiderationCache
}

func newConsiderationCacheTest() *considerationCacheTest {
	return &considerationCacheTest{
		considerationCache: NewConsiderationCache(),
	}
}

func TestConsiderationCacheGet(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name string

		// Arrange
		cacheContent map[ConsiderationID]*Consideration
		setup        func(t *testing.T, f *considerationCacheTest)

		// Act
		fetchId     ConsiderationID
		fetchTarget func() *Consideration

		// Assert
		expectCacheContent    map[ConsiderationID]*Consideration
		expectConsiderationId ConsiderationID
		expectPanic           bool
	}

	tests := []testCase{
		{
			name: "Cache hit",
			cacheContent: map[ConsiderationID]*Consideration{
				"cons1": {ID: "cons1"},
			},
			fetchId: "cons1",
			fetchTarget: func() *Consideration {
				return &Consideration{ID: "cons1"}
			},
			expectConsiderationId: "cons1",
			expectPanic:           false,
			expectCacheContent: map[ConsiderationID]*Consideration{
				"cons1": {ID: "cons1"},
			},
		},
		{
			name: "Cache miss",
			cacheContent: map[ConsiderationID]*Consideration{
				"cons1": {ID: "cons1"},
			},
			fetchId: "cons2",
			fetchTarget: func() *Consideration {
				return &Consideration{ID: "cons2"}
			},
			expectConsiderationId: "cons2",
			expectPanic:           false,
			expectCacheContent: map[ConsiderationID]*Consideration{
				"cons1": {ID: "cons1"},
				"cons2": {ID: "cons2"},
			},
		},
		{
			name:         "Panic on ID mismatch",
			cacheContent: map[ConsiderationID]*Consideration{},
			fetchId:      "cons3",
			fetchTarget: func() *Consideration {
				return &Consideration{ID: "wrong_id"}
			},
			expectConsiderationId: "cons3",
			expectPanic:           true,
			expectCacheContent:    map[ConsiderationID]*Consideration{},
		},
	}

	for _, tc := range tests {
		tc := tc // capture range variable
		t.Run(tc.name, func(t *testing.T) {
			f := newConsiderationCacheTest()

			// Arrange
			for _, cons := range tc.cacheContent {
				f.considerationCache.register(cons)
			}

			if tc.setup != nil {
				tc.setup(t, f)
			}

			// Act
			var got *Consideration
			var panicked bool
			func() {
				defer func() {
					if r := recover(); r != nil {
						panicked = true
					}
				}()
				got = f.considerationCache.Get(tc.fetchId, tc.fetchTarget)
			}()

			// Assert
			assert.Equal(t, tc.expectPanic, panicked, "Expected panic state to match")
			if !panicked {
				assert.Equal(t, tc.expectConsiderationId, got.ID, "Expected consideration ID to match")
				assert.Equal(t, tc.expectCacheContent, f.considerationCache.cache, "Expected cache content to match")
			}
		})
	}
}
