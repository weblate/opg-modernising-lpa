package donor

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/ministryofjustice/opg-modernising-lpa/app/internal/actor"
	"github.com/ministryofjustice/opg-modernising-lpa/app/internal/page"
	"github.com/ministryofjustice/opg-modernising-lpa/app/internal/validation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetAreYouHappyIfRemainingReplacementAttorneysCanContinueToAct(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "/", nil)

	template := newMockTemplate(t)
	template.
		On("Execute", w, &areYouHappyIfRemainingReplacementAttorneysCanContinueToActData{
			App:     testAppData,
			Options: actor.YesNoValues,
		}).
		Return(nil)

	err := AreYouHappyIfRemainingReplacementAttorneysCanContinueToAct(template.Execute, nil)(testAppData, w, r, &page.Lpa{})
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestGetAreYouHappyIfRemainingReplacementAttorneysCanContinueToActWhenTemplateErrors(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "/", nil)

	template := newMockTemplate(t)
	template.
		On("Execute", w, mock.Anything).
		Return(expectedError)

	err := AreYouHappyIfRemainingReplacementAttorneysCanContinueToAct(template.Execute, nil)(testAppData, w, r, &page.Lpa{})
	resp := w.Result()

	assert.Equal(t, expectedError, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestPostAreYouHappyIfRemainingReplacementAttorneysCanContinueToAct(t *testing.T) {
	for _, happy := range []actor.YesNo{actor.Yes, actor.No} {
		t.Run(happy.String(), func(t *testing.T) {
			form := url.Values{
				"happy": {happy.String()},
			}

			w := httptest.NewRecorder()
			r, _ := http.NewRequest(http.MethodPost, "/", strings.NewReader(form.Encode()))
			r.Header.Add("Content-Type", page.FormUrlEncoded)

			donorStore := newMockDonorStore(t)
			donorStore.
				On("Put", r.Context(), &page.Lpa{
					ID:                           "lpa-id",
					ReplacementAttorneyDecisions: actor.AttorneyDecisions{HappyIfRemainingCanContinueToAct: happy},
					Tasks:                        page.Tasks{YourDetails: actor.TaskCompleted, ChooseAttorneys: actor.TaskCompleted},
				}).
				Return(nil)

			err := AreYouHappyIfRemainingReplacementAttorneysCanContinueToAct(nil, donorStore)(testAppData, w, r, &page.Lpa{
				ID:    "lpa-id",
				Tasks: page.Tasks{YourDetails: actor.TaskCompleted, ChooseAttorneys: actor.TaskCompleted},
			})
			resp := w.Result()

			assert.Nil(t, err)
			assert.Equal(t, http.StatusFound, resp.StatusCode)
			assert.Equal(t, page.Paths.TaskList.Format("lpa-id"), resp.Header.Get("Location"))
		})
	}
}

func TestPostAreYouHappyIfRemainingReplacementAttorneysCanContinueToActWhenStoreErrors(t *testing.T) {
	form := url.Values{
		"happy": {actor.Yes.String()},
	}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodPost, "/", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", page.FormUrlEncoded)

	donorStore := newMockDonorStore(t)
	donorStore.
		On("Put", r.Context(), mock.Anything).
		Return(expectedError)

	err := AreYouHappyIfRemainingReplacementAttorneysCanContinueToAct(nil, donorStore)(testAppData, w, r, &page.Lpa{})

	assert.Equal(t, expectedError, err)
}

func TestPostAreYouHappyIfRemainingReplacementAttorneysCanContinueToActWhenValidationErrors(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodPost, "/", strings.NewReader(""))
	r.Header.Add("Content-Type", page.FormUrlEncoded)

	template := newMockTemplate(t)
	template.
		On("Execute", w, mock.MatchedBy(func(data *areYouHappyIfRemainingReplacementAttorneysCanContinueToActData) bool {
			return assert.Equal(t, validation.With("happy", validation.SelectError{Label: "yesIfYouAreHappyIfRemainingReplacementAttorneysCanContinueToAct"}), data.Errors)
		})).
		Return(nil)

	err := AreYouHappyIfRemainingReplacementAttorneysCanContinueToAct(template.Execute, nil)(testAppData, w, r, &page.Lpa{})
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
