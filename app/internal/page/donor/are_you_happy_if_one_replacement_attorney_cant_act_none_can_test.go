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

func TestGetAreYouHappyIfOneReplacementAttorneyCantActNoneCan(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "/", nil)

	template := newMockTemplate(t)
	template.
		On("Execute", w, &areYouHappyIfOneReplacementAttorneyCantActNoneCanData{
			App:     testAppData,
			Options: actor.YesNoValues,
		}).
		Return(nil)

	err := AreYouHappyIfOneReplacementAttorneyCantActNoneCan(template.Execute, nil)(testAppData, w, r, &page.Lpa{})
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestGetAreYouHappyIfOneReplacementAttorneyCantActNoneCanWhenTemplateErrors(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "/", nil)

	template := newMockTemplate(t)
	template.
		On("Execute", w, mock.Anything).
		Return(expectedError)

	err := AreYouHappyIfOneReplacementAttorneyCantActNoneCan(template.Execute, nil)(testAppData, w, r, &page.Lpa{})
	resp := w.Result()

	assert.Equal(t, expectedError, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestPostAreYouHappyIfOneReplacementAttorneyCantActNoneCan(t *testing.T) {
	testcases := map[string]struct {
		happy    actor.YesNo
		lpa      *page.Lpa
		redirect page.LpaPath
	}{
		"yes": {
			happy: actor.Yes,
			lpa: &page.Lpa{
				ID:                           "lpa-id",
				ReplacementAttorneyDecisions: actor.AttorneyDecisions{HappyIfOneCannotActNoneCan: actor.Yes},
				Tasks:                        page.Tasks{YourDetails: actor.TaskCompleted, ChooseAttorneys: actor.TaskCompleted},
			},
			redirect: page.Paths.TaskList,
		},
		"no": {
			happy: actor.No,
			lpa: &page.Lpa{
				ID:                           "lpa-id",
				ReplacementAttorneyDecisions: actor.AttorneyDecisions{HappyIfOneCannotActNoneCan: actor.No},
				Tasks:                        page.Tasks{YourDetails: actor.TaskCompleted, ChooseAttorneys: actor.TaskCompleted},
			},
			redirect: page.Paths.AreYouHappyIfRemainingReplacementAttorneysCanContinueToAct,
		},
	}

	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			form := url.Values{
				"happy": {tc.happy.String()},
			}

			w := httptest.NewRecorder()
			r, _ := http.NewRequest(http.MethodPost, "/", strings.NewReader(form.Encode()))
			r.Header.Add("Content-Type", page.FormUrlEncoded)

			donorStore := newMockDonorStore(t)
			donorStore.
				On("Put", r.Context(), tc.lpa).
				Return(nil)

			err := AreYouHappyIfOneReplacementAttorneyCantActNoneCan(nil, donorStore)(testAppData, w, r, &page.Lpa{
				ID:    "lpa-id",
				Tasks: page.Tasks{YourDetails: actor.TaskCompleted, ChooseAttorneys: actor.TaskCompleted},
			})
			resp := w.Result()

			assert.Nil(t, err)
			assert.Equal(t, http.StatusFound, resp.StatusCode)
			assert.Equal(t, tc.redirect.Format("lpa-id"), resp.Header.Get("Location"))
		})
	}
}

func TestPostAreYouHappyIfOneReplacementAttorneyCantActNoneCanWhenStoreErrors(t *testing.T) {
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

	err := AreYouHappyIfOneReplacementAttorneyCantActNoneCan(nil, donorStore)(testAppData, w, r, &page.Lpa{})

	assert.Equal(t, expectedError, err)
}

func TestPostAreYouHappyIfOneReplacementAttorneyCantActNoneCanWhenValidationErrors(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodPost, "/", strings.NewReader(""))
	r.Header.Add("Content-Type", page.FormUrlEncoded)

	template := newMockTemplate(t)
	template.
		On("Execute", w, mock.MatchedBy(func(data *areYouHappyIfOneReplacementAttorneyCantActNoneCanData) bool {
			return assert.Equal(t, validation.With("happy", validation.SelectError{Label: "yesIfYouAreHappyIfOneReplacementAttorneyCantActNoneCan"}), data.Errors)
		})).
		Return(nil)

	err := AreYouHappyIfOneReplacementAttorneyCantActNoneCan(template.Execute, nil)(testAppData, w, r, &page.Lpa{})
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
