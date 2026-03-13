package ai

type ConsiderationCache struct {
	cache map[ConsiderationID]*Consideration
}

func NewConsiderationCache() *ConsiderationCache {
	return &ConsiderationCache{
		cache: make(map[ConsiderationID]*Consideration),
	}
}

func (cc *ConsiderationCache) register(consideration *Consideration) {
	cc.cache[consideration.ID] = consideration
}

func (cc *ConsiderationCache) Get(
	id ConsiderationID,
	makeConsFunc func() *Consideration,
) *Consideration {
	if consideration, exists := cc.cache[id]; exists {
		return consideration
	}

	consideration := makeConsFunc()

	if consideration.ID != id {
		panic("Consideration ID mismatch: expected " + string(id) + ", got " + string(consideration.ID))
	}

	cc.register(consideration)
	return consideration
}
