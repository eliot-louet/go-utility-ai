package ai

type History[ID comparable, T any] struct {
	Entity   T
	Targets  map[Target]int64
	LastTime int64
}

type ActionHistory = History[ActionID, Action]
type BehaviorHistory = History[BehaviorID, Behavior]

func (h *History[ID, T]) HaveTarget(target Target) bool {
	if h.Targets == nil {
		return false
	}

	_, exists := h.Targets[target]
	return exists
}

func addHistory[ID comparable, T interface {
	ID() ID
	ShouldAddToHistory() bool
}](
	store map[ID]*History[ID, T],
	entity T,
	target Target,
	tick int64,
) {
	// If the entity should not be added to history, return early
	if entity.ShouldAddToHistory() == false {
		return
	}

	// Get the entity's ID
	id := entity.ID()

	// Get or create the history for the entity
	history, exists := store[id]
	if !exists {
		history = &History[ID, T]{
			Entity:  entity,
			Targets: make(map[Target]int64),
		}
		store[id] = history
	}

	history.LastTime = tick

	if target != nil {
		history.Targets[target] = tick
	}
}

func (c *Context) AddActionHistory(action Action, target Target) {
	addHistory(c.ActionHistory, action, target, c.Environment.TimeSinceStart())
}

func (c *Context) AddBehaviorHistory(behavior Behavior, target Target) {
	addHistory(c.BehaviorHistory, behavior, target, c.Environment.TimeSinceStart())
}

func getLastTime[ID comparable, T interface{ ID() ID }](
	store map[ID]*History[ID, T],
	entity T,
	target Target,
) (int64, bool) {

	history, exists := store[entity.ID()]
	if !exists {
		return 0, false
	}

	if target != nil {
		lastTime, ok := history.Targets[target]
		return lastTime, ok
	}

	return history.LastTime, true
}

func (c *Context) GetActionLastTime(action Action, target Target) (int64, bool) {
	return getLastTime(c.ActionHistory, action, target)
}

func (c *Context) GetBehaviorLastTime(behavior Behavior, target Target) (int64, bool) {
	return getLastTime(c.BehaviorHistory, behavior, target)
}
