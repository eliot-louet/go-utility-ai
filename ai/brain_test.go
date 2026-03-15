package ai_test

import (
	"testing"

	"github.com/eliot-louet/go-utility-ai/ai"
	gomock "go.uber.org/mock/gomock"
)

type currentData struct {
	Action   actionData
	Behavior behaviorData
	Target   ai.Target
	Score    float64
	State    map[string]interface{}
	Running  bool
}

type behaviorPackageData struct {
	Behaviors []behaviorData
	Cond      func(ctx *ai.Context) bool
}

type behaviorData struct {
	ID                 ai.BehaviorID
	Name               string
	Considerations     []considerationData
	Weight             float64
	Provider           providerData
	Action             actionData
	ShouldAddToHistory bool
}

func (b behaviorData) MakeBehavior(f *brainFixture, ctrl *gomock.Controller) *ai.MockBehavior {
	mockBehavior := ai.NewMockBehavior(ctrl)
	mockBehavior.EXPECT().ID().Return(b.ID).AnyTimes()
	mockBehavior.EXPECT().Name().Return(b.Name).AnyTimes()
	mockBehavior.EXPECT().ShouldAddToHistory().Return(b.ShouldAddToHistory).AnyTimes()
	mockBehavior.EXPECT().Weight(gomock.Any(), gomock.Any()).Return(1.0).AnyTimes()

	provider := ai.NewMockTargetProvider(ctrl)
	provider.EXPECT().ID().Return(b.Provider.ID).AnyTimes()
	provider.EXPECT().ShouldCache().Return(b.Provider.ShouldCache).AnyTimes()
	provider.EXPECT().ForEachTarget(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx *ai.Context, yield func(ai.Target) bool) {
		for _, target := range b.Provider.Targets {
			if !yield(target) {
				return
			}
		}
	}).AnyTimes()

	mockBehavior.EXPECT().Provider(gomock.Any()).Return(provider).AnyTimes()

	considerations := make([]*ai.Consideration, 0, len(b.Considerations))

	for _, consideration := range b.Considerations {
		mockCurve := ai.NewMockCurve(ctrl)
		mockCurve.EXPECT().Apply(gomock.Any()).Return(consideration.Score).AnyTimes()

		consideration := &ai.Consideration{
			ID:            consideration.ID,
			MinValue:      0,
			MaxValue:      100,
			ResponseCurve: mockCurve,
			ShouldCache:   consideration.ShouldCache,
			InputFunc: func(ctx *ai.Context, target ai.Target) float64 {
				return 50.0
			},
		}

		considerations = append(considerations, consideration)
	}

	mockBehavior.EXPECT().ForEachConsideration(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(ctx *ai.Context, target ai.Target, yield func(*ai.Consideration) bool) {
		for _, c := range considerations {
			if !yield(c) {
				return
			}
		}
	}).AnyTimes()

	mockBehavior.EXPECT().Action(gomock.Any(), gomock.Any()).Return(b.Action.MakeAction(f, ctrl)).AnyTimes()

	return mockBehavior
}

type providerData struct {
	ID          ai.TargetProviderID
	Targets     []ai.Target
	ShouldCache bool
}

type actionData struct {
	ID                 ai.ActionID
	UpdateStatus       ai.ActionStatus
	ShouldAddToHistory bool
}

func (a actionData) MakeAction(f *brainFixture, ctrl *gomock.Controller) ai.Action {
	action := ai.NewMockAction(ctrl)
	action.EXPECT().Start(gomock.Any(), gomock.Any()).AnyTimes().Do(func(ctx *ai.Context, target ai.Target) {
		f.actionHistories = append(f.actionHistories, actionHistoryEntry{
			ID:       a.ID,
			Status:   ai.Running,
			Function: "Start",
		})
	})
	action.EXPECT().Update(gomock.Any(), gomock.Any()).Return(a.UpdateStatus).AnyTimes().Do(func(ctx *ai.Context, target ai.Target) {
		f.actionHistories = append(f.actionHistories, actionHistoryEntry{
			ID:       a.ID,
			Status:   a.UpdateStatus,
			Function: "Update",
		})
	})
	action.EXPECT().Cancel(gomock.Any(), gomock.Any()).AnyTimes().Do(func(ctx *ai.Context, target ai.Target) {
		// If the action is canceled, we consider it as not added to history
		f.actionHistories = append(f.actionHistories, actionHistoryEntry{
			ID:       a.ID,
			Status:   ai.Failure,
			Function: "Cancel",
		})
	})
	action.EXPECT().ShouldAddToHistory().Return(a.ShouldAddToHistory).AnyTimes()

	return action
}

type considerationData struct {
	ID          ai.ConsiderationID
	ShouldCache bool
	Score       float64
}

type actionHistoryEntry struct {
	ID       ai.ActionID
	Status   ai.ActionStatus
	Function string
}

type brainFixture struct {
	brain           *ai.Brain
	ctrl            *gomock.Controller
	actionHistories []actionHistoryEntry
}

func newBrainFixture(t *testing.T) *brainFixture {
	t.Helper()

	return &brainFixture{
		brain: &ai.Brain{
			Context:            ai.MakeContext(nil, nil),
			BehaviorPackages:   []*ai.BehaviorPackage{},
			Current:            ai.RunningBehavior{},
			InterruptThreshold: 1.5,
			StateCache:         make(map[string]interface{}),
		},
	}
}

func (f *brainFixture) withBehaviorPackages(ctrl *gomock.Controller, packages ...behaviorPackageData) *brainFixture {
	for _, pkg := range packages {
		behaviorPackage := &ai.BehaviorPackage{
			Behaviors:     []ai.Behavior{},
			ConditionFunc: pkg.Cond,
		}

		for _, behavior := range pkg.Behaviors {
			mockBehavior := behavior.MakeBehavior(f, ctrl)

			behaviorPackage.Behaviors = append(behaviorPackage.Behaviors, mockBehavior)
		}

		f.brain.BehaviorPackages = append(f.brain.BehaviorPackages, behaviorPackage)
	}

	return f
}

func (f *brainFixture) withCurrentAction(ctrl *gomock.Controller, t *testing.T, current currentData) *brainFixture {
	t.Helper()

	f.brain.Current = ai.RunningBehavior{
		Behavior: current.Behavior.MakeBehavior(f, ctrl),
		Action:   current.Action.MakeAction(f, ctrl),
		Target:   current.Target,
		Score:    current.Score,
		State:    current.State,
		Running:  current.Running,
	}

	return f
}

func Current(id string, score float64, actionUpdateStatus ai.ActionStatus) currentData {
	return currentData{
		Action: actionData{
			ID:           ai.ActionID(id + "_action"),
			UpdateStatus: actionUpdateStatus,
		},
		Behavior: Behavior(id, score, actionUpdateStatus),
		Target:   "target",
		Score:    score,
		Running:  true,
		State:    map[string]interface{}{},
	}
}

func Behavior(id string, score float64, actionUpdateStatus ai.ActionStatus) behaviorData {
	return behaviorData{
		ID:   ai.BehaviorID(id),
		Name: id,
		Considerations: []considerationData{
			{
				ID:    ai.ConsiderationID(id + "_c"),
				Score: score,
			},
		},
		Action: actionData{
			ID:           ai.ActionID(id + "_action"),
			UpdateStatus: actionUpdateStatus,
		},
		Provider: providerData{
			ID:      ai.TargetProviderID(id + "_provider"),
			Targets: []ai.Target{"target"},
		},
	}
}

func Package(conditionChecked bool, behaviors ...behaviorData) behaviorPackageData {
	return behaviorPackageData{
		Cond: func(ctx *ai.Context) bool {
			return conditionChecked
		},
		Behaviors: behaviors,
	}
}

func TestBrain_Decide(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name string
		// Arrange
		current  currentData
		packages []behaviorPackageData
		setup    func(t *testing.T, f *brainFixture)

		// Act
		expectBehaviorID ai.BehaviorID
	}

	tests := []testCase{
		{
			name: "Single behavior with one consideration",
			packages: []behaviorPackageData{
				Package(true, Behavior("behavior1", 0.8, ai.Success)),
			},
			expectBehaviorID: "behavior1",
		},
		{
			name: "Single behavior with one considerations but wrong package condition",
			packages: []behaviorPackageData{
				Package(false, Behavior("behavior1", 0.8, ai.Success)),
				Package(true, Behavior("behavior2", 0.9, ai.Success)),
			},
			expectBehaviorID: "behavior2",
		},
		{
			name: "Multiple behaviors with different scores",
			packages: []behaviorPackageData{
				Package(true, Behavior("behavior1", 0.5, ai.Success), Behavior("behavior2", 0.9, ai.Success)),
			},
			expectBehaviorID: "behavior2",
		},
		{
			name:    "Current behavior with lower score than new behavior (but above interrupt threshold)",
			current: Current("currentBehavior", 0.6, ai.Running),
			packages: []behaviorPackageData{
				Package(true,
					Behavior("currentBehavior", 0.6, ai.Running),
					Behavior("newBehavior", 0.9, ai.Success),
				),
			},
			expectBehaviorID: "newBehavior",
		},
		{
			name:    "Current behavior with lower score than new behavior (but below interrupt threshold)",
			current: Current("currentBehavior", 0.6, ai.Running),
			packages: []behaviorPackageData{
				Package(true,
					Behavior("currentBehavior", 0.6, ai.Running),
					Behavior("newBehavior", 0.65, ai.Success),
				),
			},
			expectBehaviorID: "currentBehavior",
		},
		{
			name:    "Current behavior with higher score than new behavior",
			current: Current("currentBehavior", 0.8, ai.Running),
			packages: []behaviorPackageData{
				Package(true,
					Behavior("currentBehavior", 0.8, ai.Running),
					Behavior("newBehavior", 0.6, ai.Success),
				),
			},
			expectBehaviorID: "currentBehavior",
		},
	}

	for _, tc := range tests {
		tc := tc // capture range variable
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := newBrainFixture(t).withBehaviorPackages(ctrl, tc.packages...)

			if tc.current.Running {
				f.withCurrentAction(ctrl, t, tc.current)
			}

			if tc.setup != nil {
				tc.setup(t, f)
			}

			decision := f.brain.Decide(f.brain.Context)

			if decision.Behavior == nil || decision.Behavior.ID() != tc.expectBehaviorID {
				t.Errorf("Expected behavior ID %s, got %v", tc.expectBehaviorID, decision.Behavior)
			}
		})
	}
}

func TestBrain_Update(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name string

		current  currentData
		packages []behaviorPackageData

		expectHistory []actionHistoryEntry
	}

	tests := []testCase{
		{
			name: "Start new behavior when none running",
			packages: []behaviorPackageData{
				Package(true, Behavior("behavior1", 0.9, ai.Success)),
			},
			expectHistory: []actionHistoryEntry{
				{ID: "behavior1_action", Status: ai.Running, Function: "Start"},
				{ID: "behavior1_action", Status: ai.Success, Function: "Update"},
			},
		},
		{
			name:    "Continue current behavior when below interrupt threshold",
			current: Current("current", 0.6, ai.Running),
			packages: []behaviorPackageData{
				Package(true,
					Behavior("current", 0.6, ai.Running),
					Behavior("new", 0.65, ai.Success),
				),
			},
			expectHistory: []actionHistoryEntry{
				{ID: "current_action", Status: ai.Running, Function: "Update"},
			},
		},
		{
			name:    "Interrupt current behavior when new score much higher",
			current: Current("current", 0.6, ai.Running),
			packages: []behaviorPackageData{
				Package(true,
					Behavior("current", 0.6, ai.Running),
					Behavior("new", 0.95, ai.Success),
				),
			},
			expectHistory: []actionHistoryEntry{
				{ID: "current_action", Status: ai.Failure, Function: "Cancel"},
				{ID: "new_action", Status: ai.Running, Function: "Start"},
				{ID: "new_action", Status: ai.Success, Function: "Update"},
			},
		},
		{
			name:    "Finish action if update returns success",
			current: Current("current", 0.8, ai.Success),
			packages: []behaviorPackageData{
				Package(true, Behavior("current", 0.8, ai.Success)),
			},
			expectHistory: []actionHistoryEntry{
				{ID: "current_action", Status: ai.Success, Function: "Update"},
			},
		},
		{
			name: "No decision does nothing",
			packages: []behaviorPackageData{
				Package(false, Behavior("behavior1", 0.9, ai.Success)),
			},
			expectHistory: nil,
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			f := newBrainFixture(t).
				withBehaviorPackages(ctrl, tc.packages...)

			if tc.current.Running {
				f.withCurrentAction(ctrl, t, tc.current)
			}

			f.brain.Update(f.brain.Context)

			if len(f.actionHistories) != len(tc.expectHistory) {
				t.Fatalf("expected %d history entries got %d\n%v",
					len(tc.expectHistory),
					len(f.actionHistories),
					f.actionHistories,
				)
			}

			for i, expected := range tc.expectHistory {
				got := f.actionHistories[i]

				if got.ID != expected.ID ||
					got.Function != expected.Function ||
					got.Status != expected.Status {

					t.Fatalf(
						"history[%d] mismatch\nexpected %+v\ngot %+v",
						i,
						expected,
						got,
					)
				}
			}
		})
	}
}
