package ai

import (
	"math"
	"sync"
)

const epsilon = 1e-6

// =====================
// EFFECT DEFINITION
// =====================

type StatusEffect struct {
	ID string

	Duration int // -1 = infinite

	MaxStacks        int
	RefreshOnReapply bool

	TagsApplied   []string
	StatModifiers map[string]float64

	RemovalCondition func(*Context) bool

	OnApply  func(*Context, map[string]any)
	OnTick   func(*Context, map[string]any)
	OnRemove func(*Context, map[string]any)
}

// =====================
// RUNTIME INSTANCE
// =====================

type ActiveStatus struct {
	Effect    *StatusEffect
	AppliedAt int64
	Data      map[string]any
}

// =====================
// POOLS
// =====================

var activeStatusPool = sync.Pool{
	New: func() any {
		return &ActiveStatus{
			Data: make(map[string]any),
		}
	},
}

// =====================
// MANAGER
// =====================

type StatusManager struct {
	ActiveStatuses     []*ActiveStatus
	ActiveStatusesByID map[string][]*ActiveStatus

	GlobalStatsModifiers map[string]float64
	GlobalTags           map[string]int
}

func NewStatusManager() *StatusManager {
	return &StatusManager{
		ActiveStatuses:       make([]*ActiveStatus, 0, 16),
		ActiveStatusesByID:   make(map[string][]*ActiveStatus),
		GlobalStatsModifiers: make(map[string]float64),
		GlobalTags:           make(map[string]int),
	}
}

// =====================
// HAS TAG
// =====================
func (sm *StatusManager) HasTag(tag string) bool {
	return sm.GlobalTags[tag] > 0
}

// =====================
// HAS STATUS
// =====================
func (sm *StatusManager) HasStatus(id string) bool {
	return len(sm.ActiveStatusesByID[id]) > 0
}

// =====================
// GET STAT MODIFIER
// =====================
func (sm *StatusManager) GetStatModifier(stat string) float64 {
	return sm.GlobalStatsModifiers[stat]
}

// =====================
// APPLY
// =====================
func (sm *StatusManager) ApplyStatus(effect *StatusEffect, ctx *Context) {
	now := ctx.Environment.TimeSinceStart()
	list := sm.ActiveStatusesByID[effect.ID]

	// --- Stacking rules ---
	if effect.MaxStacks > 0 && len(list) >= effect.MaxStacks {
		if effect.RefreshOnReapply {
			for _, s := range list {
				s.AppliedAt = now
			}
		}
		return
	}

	// --- Get from pool ---
	activeStatus := activeStatusPool.Get().(*ActiveStatus)

	activeStatus.Effect = effect
	activeStatus.AppliedAt = now

	// --- Insert ---
	sm.ActiveStatuses = append(sm.ActiveStatuses, activeStatus)
	sm.ActiveStatusesByID[effect.ID] = append(list, activeStatus)

	// --- Apply tags ---
	for _, tag := range effect.TagsApplied {
		sm.GlobalTags[tag]++
	}

	// --- Apply stats ---
	for stat, mod := range effect.StatModifiers {
		sm.GlobalStatsModifiers[stat] += mod
	}

	// --- Callback ---
	if effect.OnApply != nil {
		effect.OnApply(ctx, activeStatus.Data)
	}
}

// =====================
// TICK
// =====================

func (sm *StatusManager) Tick(ctx *Context) {
	now := ctx.Environment.TimeSinceStart()

	for i := len(sm.ActiveStatuses) - 1; i >= 0; i-- {
		s := sm.ActiveStatuses[i]
		effect := s.Effect

		// --- Removal check FIRST ---
		if (effect.Duration != -1 &&
			now-s.AppliedAt >= int64(effect.Duration)) ||
			(effect.RemovalCondition != nil && effect.RemovalCondition(ctx)) {

			sm.removeAt(i, s, ctx)
			continue
		}

		// --- Tick ---
		if effect.OnTick != nil {
			effect.OnTick(ctx, s.Data)
		}
	}
}

// =====================
// REMOVE (O(1))
// =====================

func (sm *StatusManager) removeAt(index int, s *ActiveStatus, ctx *Context) {
	effect := s.Effect

	// --- Remove tags ---
	for _, tag := range effect.TagsApplied {
		sm.GlobalTags[tag]--
		if sm.GlobalTags[tag] <= 0 {
			delete(sm.GlobalTags, tag)
		}
	}

	// --- Remove stats ---
	for stat, mod := range effect.StatModifiers {
		sm.GlobalStatsModifiers[stat] -= mod
		if math.Abs(sm.GlobalStatsModifiers[stat]) < epsilon {
			delete(sm.GlobalStatsModifiers, stat)
		}
	}

	// --- Callback ---
	if effect.OnRemove != nil {
		effect.OnRemove(ctx, s.Data)
	}

	// --- Remove from ActiveStatuses (swap-remove) ---
	lastIndex := len(sm.ActiveStatuses) - 1
	sm.ActiveStatuses[index] = sm.ActiveStatuses[lastIndex]
	sm.ActiveStatuses = sm.ActiveStatuses[:lastIndex]

	// --- Remove from ID map ---
	list := sm.ActiveStatusesByID[effect.ID]
	for i, item := range list {
		if item == s {
			last := len(list) - 1
			list[i] = list[last]
			sm.ActiveStatusesByID[effect.ID] = list[:last]
			break
		}
	}

	// --- Return to pool ---
	sm.release(s)
}

// =====================
// POOL RELEASE
// =====================

func (sm *StatusManager) release(s *ActiveStatus) {
	s.Effect = nil
	s.AppliedAt = 0

	clear(s.Data)

	activeStatusPool.Put(s)
}
