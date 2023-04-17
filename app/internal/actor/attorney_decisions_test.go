package actor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakeAttorneyDecisions(t *testing.T) {
	testcases := map[string]struct {
		existing AttorneyDecisions
		how      string
		details  string
		expected AttorneyDecisions
	}{
		"without details": {
			existing: AttorneyDecisions{HappyIfOneCannotActNoneCan: "yes"},
			how:      Jointly,
			details:  "hey",
			expected: AttorneyDecisions{How: Jointly},
		},
		"with details": {
			existing: AttorneyDecisions{HappyIfOneCannotActNoneCan: "yes"},
			how:      JointlyForSomeSeverallyForOthers,
			details:  "hey",
			expected: AttorneyDecisions{How: JointlyForSomeSeverallyForOthers, Details: "hey"},
		},
		"same how without details": {
			existing: AttorneyDecisions{How: Jointly, HappyIfOneCannotActNoneCan: "yes"},
			how:      Jointly,
			details:  "hey",
			expected: AttorneyDecisions{How: Jointly, HappyIfOneCannotActNoneCan: "yes"},
		},
		"same how with details": {
			existing: AttorneyDecisions{How: JointlyForSomeSeverallyForOthers, Details: "what", HappyIfOneCannotActNoneCan: "yes"},
			how:      JointlyForSomeSeverallyForOthers,
			details:  "hey",
			expected: AttorneyDecisions{How: JointlyForSomeSeverallyForOthers, Details: "hey", HappyIfOneCannotActNoneCan: "yes"},
		},
	}

	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.expected, MakeAttorneyDecisions(tc.existing, tc.how, tc.details))
		})
	}
}

func TestAttorneyDecisionsRequiresHappiness(t *testing.T) {
	testcases := map[string]struct {
		attorneyCount int
		how           string
		expected      bool
	}{
		"jointly attorneys": {
			attorneyCount: 2,
			how:           Jointly,
			expected:      true,
		},
		"jointly for some severally for others attorney": {
			attorneyCount: 2,
			how:           JointlyForSomeSeverallyForOthers,
			expected:      true,
		},
		"not for jointly and severally attorney": {
			attorneyCount: 2,
			how:           JointlyAndSeverally,
			expected:      false,
		},
		"not for single attorney": {
			attorneyCount: 1,
			how:           Jointly,
			expected:      false,
		},
	}

	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			decisions := AttorneyDecisions{How: tc.how}

			assert.Equal(t, tc.expected, decisions.RequiresHappiness(tc.attorneyCount))
		})
	}
}

func TestAttorneyDecisionsIsComplete(t *testing.T) {
	testcases := map[string]struct {
		attorneyCount int
		decisions     AttorneyDecisions
		expected      bool
	}{
		"jointly attorneys, happy": {
			attorneyCount: 2,
			decisions:     AttorneyDecisions{How: Jointly, HappyIfOneCannotActNoneCan: "yes"},
			expected:      true,
		},
		"jointly for some severally for others attorney, happy": {
			attorneyCount: 2,
			decisions:     AttorneyDecisions{How: JointlyForSomeSeverallyForOthers, HappyIfOneCannotActNoneCan: "yes"},
			expected:      true,
		},
		"jointly attorneys, unhappy": {
			attorneyCount: 2,
			decisions:     AttorneyDecisions{How: Jointly, HappyIfOneCannotActNoneCan: "no", HappyIfRemainingCanContinueToAct: "no"},
			expected:      true,
		},
		"jointly attorneys, mixed happy": {
			attorneyCount: 2,
			decisions:     AttorneyDecisions{How: Jointly, HappyIfOneCannotActNoneCan: "no", HappyIfRemainingCanContinueToAct: "yes"},
			expected:      true,
		},
		"jointly attorneys, unhappy missing": {
			attorneyCount: 2,
			decisions:     AttorneyDecisions{How: Jointly, HappyIfOneCannotActNoneCan: "no"},
			expected:      false,
		},
		"jointly and severally attorney": {
			attorneyCount: 2,
			decisions:     AttorneyDecisions{How: JointlyAndSeverally},
			expected:      true,
		},
		"single attorney": {
			attorneyCount: 1,
			decisions:     AttorneyDecisions{How: Jointly},
			expected:      true,
		},
		"missing how": {
			attorneyCount: 1,
			decisions:     AttorneyDecisions{},
			expected:      false,
		},
	}

	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.expected, tc.decisions.IsComplete(tc.attorneyCount))
		})
	}
}