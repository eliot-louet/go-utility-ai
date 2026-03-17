package ai

type ActionCache struct {
	cache map[ActionID]Action
}

func NewActionCache() *ActionCache {
	return &ActionCache{
		cache: make(map[ActionID]Action),
	}
}

func (cc *ActionCache) register(action Action) {
	cc.cache[action.ID()] = action
}

func (cc *ActionCache) Get(
	id ActionID,
	makeConsFunc func() Action,
) Action {
	if action, exists := cc.cache[id]; exists {
		return action
	}

	action := makeConsFunc()

	if action.ID() != id {
		panic("Action ID mismatch: expected " + string(id) + ", got " + string(action.ID()))
	}

	cc.register(action)
	return action
}
