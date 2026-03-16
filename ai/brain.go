package ai

type Brain struct {
	Context            *Context
	BehaviorPackages   []*BehaviorPackage
	Current            RunningBehavior
	InterruptThreshold float64

	StateCache map[string]interface{}
}

func (b *Brain) shouldInterrupt(decision Decision) bool {
	return decision.Score > b.Current.Score*b.InterruptThreshold &&
		!(b.Current.Behavior.Name() == decision.Behavior.Name() &&
			b.Current.Target == decision.Target)
}

// Interrupt the currently running action if the new decision has a significantly higher score
func (b *Brain) interruptCurrent(ctx *Context) {
	b.Current.Action.Cancel(ctx, b.Current.Target)

	b.Current.Running = false
	clear(b.Current.State)
}

// Update the currently running action and return whether it is still running
func (b *Brain) updateCurrent(ctx *Context) bool {
	status := b.Current.Action.Update(ctx, b.Current.Target)

	if status == Running {
		return true
	}

	b.finishCurrent(ctx, b.Current.Action)
	return false
}

func (b *Brain) finishCurrent(ctx *Context, action Action) {
	b.Current.Running = false
	clear(b.Current.State)

	if action.ShouldAddToHistory() {
		ctx.AddActionHistory(action, b.Current.Target)
	}

	if b.Current.Behavior.ShouldAddToHistory() {
		ctx.AddBehaviorHistory(b.Current.Behavior, b.Current.Target)
	}
}

func (b *Brain) Update(ctx *Context) {
	// Find the best behavior and target based on the current context
	decision := b.Decide(ctx)

	// If there is a currently running behavior, check if we should interrupt it
	if b.Current.Running {
		// If the new decision is not better than the current one,
		if !b.shouldInterrupt(decision) {
			// Just update the current action and return
			b.updateCurrent(ctx)
			return
		} else {
			// Else, interrupt the current action
			b.interruptCurrent(ctx)
		}
	}

	// If their is no new behavior to run, just return
	if decision.Behavior == nil {
		return
	}

	// Start the new behavior
	b.startDecision(ctx, decision)

	// Clear caches after decision to avoid stale data
	b.clearCaches(ctx)
}

func (b *Brain) startDecision(ctx *Context, decision Decision) {
	clear(b.StateCache)

	action := decision.Behavior.Action(ctx, decision.Target)

	action.Start(ctx, decision.Target)

	b.Current = RunningBehavior{
		Running:  true,
		Behavior: decision.Behavior,
		Target:   decision.Target,
		Score:    decision.Score,
		Action:   action,
	}

	status := action.Update(ctx, decision.Target)

	if status != Running {
		b.finishCurrent(ctx, action)
	}
}

func (b *Brain) clearCaches(ctx *Context) {
	clear(ctx.providerCache)
	clear(ctx.considerationCache)
}

func (b *Brain) boostScoreIfCurrent(behavior Behavior, target Target, score float64) float64 {
	if !b.Current.Running {
		return score
	}

	if behavior.Name() == b.Current.Behavior.Name() && target == b.Current.Target {
		return score * 1.25
	}

	return score
}

func (b *Brain) Decide(ctx *Context) Decision {
	bestScore := 0.0
	var bestBehavior Behavior
	var bestTarget Target

	for _, pkg := range b.BehaviorPackages {
		for _, behavior := range pkg.Behaviors {
			if pkg.ConditionFunc != nil && !pkg.ConditionFunc(ctx) {
				continue
			}

			provider := behavior.Provider(ctx)

			targets := provider.Targets(ctx)

			for _, target := range targets {
				score := EvaluateBehavior(ctx, behavior, target)
				score = b.boostScoreIfCurrent(behavior, target, score)

				// Update the best behavior if this one has a higher score
				if score > bestScore {
					bestScore = score
					bestBehavior = behavior
					bestTarget = target
				}
			}
		}
	}

	return Decision{
		Behavior: bestBehavior,
		Target:   bestTarget,
		Score:    bestScore,
	}
}
