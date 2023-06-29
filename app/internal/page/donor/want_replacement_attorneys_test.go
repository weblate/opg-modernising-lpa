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

func TestGetWantReplacementAttorneys(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "/", nil)

	template := newMockTemplate(t)
	template.
		On("Execute", w, &wantReplacementAttorneysData{
			App:     testAppData,
			Lpa:     &page.Lpa{},
			Form:    &wantReplacementAttorneysForm{},
			Options: actor.YesNoValues,
		}).
		Return(nil)

	err := WantReplacementAttorneys(template.Execute, nil)(testAppData, w, r, &page.Lpa{})
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestGetWantReplacementAttorneysWithExistingReplacementAttorneys(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "/", nil)

	template := newMockTemplate(t)

	err := WantReplacementAttorneys(template.Execute, nil)(testAppData, w, r, &page.Lpa{ID: "lpa-id", ReplacementAttorneys: actor.Attorneys{{FirstNames: "this"}}})
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusFound, resp.StatusCode)
	assert.Equal(t, page.Paths.ChooseReplacementAttorneysSummary.Format("lpa-id"), resp.Header.Get("Location"))
}

func TestGetWantReplacementAttorneysFromStore(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "/", nil)

	template := newMockTemplate(t)
	template.
		On("Execute", w, &wantReplacementAttorneysData{
			App: testAppData,
			Lpa: &page.Lpa{WantReplacementAttorneys: actor.Yes},
			Form: &wantReplacementAttorneysForm{
				Want: actor.Yes,
			},
			Options: actor.YesNoValues,
		}).
		Return(nil)

	err := WantReplacementAttorneys(template.Execute, nil)(testAppData, w, r, &page.Lpa{WantReplacementAttorneys: actor.Yes})
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestGetWantReplacementAttorneysWhenTemplateErrors(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "/", nil)

	template := newMockTemplate(t)
	template.
		On("Execute", w, mock.Anything).
		Return(expectedError)

	err := WantReplacementAttorneys(template.Execute, nil)(testAppData, w, r, &page.Lpa{})
	resp := w.Result()

	assert.Equal(t, expectedError, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestPostWantReplacementAttorneys(t *testing.T) {
	testCases := map[string]struct {
		want                         actor.YesNo
		existingReplacementAttorneys actor.Attorneys
		expectedReplacementAttorneys actor.Attorneys
		taskState                    actor.TaskState
		redirect                     page.LpaPath
	}{
		"yes": {
			want:                         actor.Yes,
			existingReplacementAttorneys: actor.Attorneys{{ID: "123"}},
			expectedReplacementAttorneys: actor.Attorneys{{ID: "123"}},
			taskState:                    actor.TaskInProgress,
			redirect:                     page.Paths.ChooseReplacementAttorneys,
		},
		"no": {
			want: actor.No,
			existingReplacementAttorneys: actor.Attorneys{
				{ID: "123"},
				{ID: "345"},
			},
			expectedReplacementAttorneys: actor.Attorneys{},
			taskState:                    actor.TaskCompleted,
			redirect:                     page.Paths.TaskList,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			form := url.Values{
				"want": {tc.want.String()},
			}

			w := httptest.NewRecorder()
			r, _ := http.NewRequest(http.MethodPost, "/", strings.NewReader(form.Encode()))
			r.Header.Add("Content-Type", page.FormUrlEncoded)

			donorStore := newMockDonorStore(t)
			donorStore.
				On("Put", r.Context(), &page.Lpa{
					ID:                       "lpa-id",
					WantReplacementAttorneys: tc.want,
					ReplacementAttorneys:     tc.expectedReplacementAttorneys,
					Tasks:                    page.Tasks{YourDetails: actor.TaskCompleted, ChooseAttorneys: actor.TaskCompleted, ChooseReplacementAttorneys: tc.taskState},
				}).
				Return(nil)

			err := WantReplacementAttorneys(nil, donorStore)(testAppData, w, r, &page.Lpa{
				ID:                   "lpa-id",
				ReplacementAttorneys: tc.existingReplacementAttorneys,
				Tasks:                page.Tasks{YourDetails: actor.TaskCompleted, ChooseAttorneys: actor.TaskCompleted},
			})
			resp := w.Result()

			assert.Nil(t, err)
			assert.Equal(t, http.StatusFound, resp.StatusCode)
			assert.Equal(t, tc.redirect.Format("lpa-id"), resp.Header.Get("Location"))
		})
	}
}

func TestPostWantReplacementAttorneysWhenStoreErrors(t *testing.T) {
	form := url.Values{
		"want": {actor.Yes.String()},
	}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodPost, "/", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", page.FormUrlEncoded)

	donorStore := newMockDonorStore(t)
	donorStore.
		On("Put", r.Context(), mock.Anything).
		Return(expectedError)

	err := WantReplacementAttorneys(nil, donorStore)(testAppData, w, r, &page.Lpa{})

	assert.Equal(t, expectedError, err)
}

func TestPostWantReplacementAttorneysWhenValidationErrors(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodPost, "/", strings.NewReader(""))
	r.Header.Add("Content-Type", page.FormUrlEncoded)

	template := newMockTemplate(t)
	template.
		On("Execute", w, mock.MatchedBy(func(data *wantReplacementAttorneysData) bool {
			return assert.Equal(t, validation.With("want", validation.SelectError{Label: "yesToAddReplacementAttorneys"}), data.Errors)
		})).
		Return(nil)

	err := WantReplacementAttorneys(template.Execute, nil)(testAppData, w, r, &page.Lpa{})
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestReadWantReplacementAttorneysForm(t *testing.T) {
	form := url.Values{
		"want": {actor.Yes.String()},
	}

	r, _ := http.NewRequest(http.MethodPost, "/", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", page.FormUrlEncoded)

	result := readWantReplacementAttorneysForm(r)

	assert.Equal(t, actor.Yes, result.Want)
}

func TestWantReplacementAttorneysFormValidate(t *testing.T) {
	testCases := map[string]struct {
		form   *wantReplacementAttorneysForm
		errors validation.List
	}{
		"valid": {
			form: &wantReplacementAttorneysForm{},
		},
		"invalid": {
			form: &wantReplacementAttorneysForm{
				Error: expectedError,
			},
			errors: validation.With("want", validation.SelectError{Label: "yesToAddReplacementAttorneys"}),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.errors, tc.form.Validate())
		})
	}
}
