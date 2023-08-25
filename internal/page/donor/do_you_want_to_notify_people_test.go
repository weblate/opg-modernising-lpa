package donor

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/ministryofjustice/opg-modernising-lpa/internal/actor"
	"github.com/ministryofjustice/opg-modernising-lpa/internal/form"
	"github.com/ministryofjustice/opg-modernising-lpa/internal/page"
	"github.com/ministryofjustice/opg-modernising-lpa/internal/validation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetDoYouWantToNotifyPeople(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "/", nil)

	template := newMockTemplate(t)
	template.
		On("Execute", w, &doYouWantToNotifyPeopleData{
			App:     testAppData,
			Lpa:     &page.Lpa{},
			Form:    &form.YesNoForm{},
			Options: form.YesNoValues,
		}).
		Return(nil)

	err := DoYouWantToNotifyPeople(template.Execute, nil)(testAppData, w, r, &page.Lpa{})
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestGetDoYouWantToNotifyPeopleFromStore(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "/", nil)

	template := newMockTemplate(t)
	template.
		On("Execute", w, &doYouWantToNotifyPeopleData{
			App: testAppData,
			Lpa: &page.Lpa{
				DoYouWantToNotifyPeople: form.Yes,
			},
			Form: &form.YesNoForm{
				YesNo: form.Yes,
			},
			Options: form.YesNoValues,
		}).
		Return(nil)

	err := DoYouWantToNotifyPeople(template.Execute, nil)(testAppData, w, r, &page.Lpa{
		DoYouWantToNotifyPeople: form.Yes,
	})
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestGetDoYouWantToNotifyPeopleHowAttorneysWorkTogether(t *testing.T) {
	testCases := map[string]struct {
		howWorkTogether  actor.AttorneysAct
		expectedTransKey string
	}{
		"jointly": {
			howWorkTogether:  actor.Jointly,
			expectedTransKey: "jointlyDescription",
		},
		"jointly and severally": {
			howWorkTogether:  actor.JointlyAndSeverally,
			expectedTransKey: "jointlyAndSeverallyDescription",
		},
		"jointly for some severally for others": {
			howWorkTogether:  actor.JointlyForSomeSeverallyForOthers,
			expectedTransKey: "jointlyForSomeSeverallyForOthersDescription",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r, _ := http.NewRequest(http.MethodGet, "/", nil)

			template := newMockTemplate(t)
			template.
				On("Execute", w, &doYouWantToNotifyPeopleData{
					App: testAppData,
					Lpa: &page.Lpa{
						DoYouWantToNotifyPeople: form.Yes,
						AttorneyDecisions:       actor.AttorneyDecisions{How: tc.howWorkTogether},
					},
					Form: &form.YesNoForm{
						YesNo: form.Yes,
					},
					Options:         form.YesNoValues,
					HowWorkTogether: tc.expectedTransKey,
				}).
				Return(nil)

			err := DoYouWantToNotifyPeople(template.Execute, nil)(testAppData, w, r, &page.Lpa{
				DoYouWantToNotifyPeople: form.Yes,
				AttorneyDecisions:       actor.AttorneyDecisions{How: tc.howWorkTogether},
			})
			resp := w.Result()

			assert.Nil(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
		})
	}
}

func TestGetDoYouWantToNotifyPeopleFromStoreWithPeople(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "/", nil)

	template := newMockTemplate(t)

	err := DoYouWantToNotifyPeople(template.Execute, nil)(testAppData, w, r, &page.Lpa{
		ID: "lpa-id",
		PeopleToNotify: actor.PeopleToNotify{
			{ID: "123"},
		},
	})
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusFound, resp.StatusCode)
	assert.Equal(t, page.Paths.ChoosePeopleToNotifySummary.Format("lpa-id"), resp.Header.Get("Location"))
}

func TestGetDoYouWantToNotifyPeopleWhenTemplateErrors(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "/", nil)

	template := newMockTemplate(t)
	template.
		On("Execute", w, mock.Anything).
		Return(expectedError)

	err := DoYouWantToNotifyPeople(template.Execute, nil)(testAppData, w, r, &page.Lpa{})
	resp := w.Result()

	assert.Equal(t, expectedError, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestPostDoYouWantToNotifyPeople(t *testing.T) {
	testCases := []struct {
		YesNo            form.YesNo
		ExistingAnswer   form.YesNo
		ExpectedRedirect page.LpaPath
		ExpectedStatus   actor.TaskState
	}{
		{
			YesNo:            form.Yes,
			ExistingAnswer:   form.No,
			ExpectedRedirect: page.Paths.ChoosePeopleToNotify,
			ExpectedStatus:   actor.TaskInProgress,
		},
		{
			YesNo:            form.No,
			ExistingAnswer:   form.Yes,
			ExpectedRedirect: page.Paths.TaskList,
			ExpectedStatus:   actor.TaskCompleted,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.YesNo.String(), func(t *testing.T) {
			form := url.Values{
				"yes-no": {tc.YesNo.String()},
			}

			w := httptest.NewRecorder()
			r, _ := http.NewRequest(http.MethodPost, "/", strings.NewReader(form.Encode()))
			r.Header.Add("Content-Type", page.FormUrlEncoded)

			donorStore := newMockDonorStore(t)
			donorStore.
				On("Put", r.Context(), &page.Lpa{
					ID:                      "lpa-id",
					DoYouWantToNotifyPeople: tc.YesNo,
					Tasks: page.Tasks{
						YourDetails:                actor.TaskCompleted,
						ChooseAttorneys:            actor.TaskCompleted,
						ChooseReplacementAttorneys: actor.TaskCompleted,
						WhenCanTheLpaBeUsed:        actor.TaskCompleted,
						Restrictions:               actor.TaskCompleted,
						CertificateProvider:        actor.TaskCompleted,
						PeopleToNotify:             tc.ExpectedStatus,
					},
				}).
				Return(nil)

			err := DoYouWantToNotifyPeople(nil, donorStore)(testAppData, w, r, &page.Lpa{
				ID:                      "lpa-id",
				DoYouWantToNotifyPeople: tc.ExistingAnswer,
				Tasks: page.Tasks{
					YourDetails:                actor.TaskCompleted,
					ChooseAttorneys:            actor.TaskCompleted,
					ChooseReplacementAttorneys: actor.TaskCompleted,
					WhenCanTheLpaBeUsed:        actor.TaskCompleted,
					Restrictions:               actor.TaskCompleted,
					CertificateProvider:        actor.TaskCompleted,
				},
			})
			resp := w.Result()

			assert.Nil(t, err)
			assert.Equal(t, http.StatusFound, resp.StatusCode)
			assert.Equal(t, tc.ExpectedRedirect.Format("lpa-id"), resp.Header.Get("Location"))
		})
	}
}

func TestPostDoYouWantToNotifyPeopleWhenStoreErrors(t *testing.T) {
	f := url.Values{
		"yes-no": {form.Yes.String()},
	}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodPost, "/", strings.NewReader(f.Encode()))
	r.Header.Add("Content-Type", page.FormUrlEncoded)

	donorStore := newMockDonorStore(t)
	donorStore.
		On("Put", r.Context(), &page.Lpa{
			DoYouWantToNotifyPeople: form.Yes,
			Tasks:                   page.Tasks{PeopleToNotify: actor.TaskInProgress},
		}).
		Return(expectedError)

	err := DoYouWantToNotifyPeople(nil, donorStore)(testAppData, w, r, &page.Lpa{})

	assert.Equal(t, expectedError, err)
}

func TestPostDoYouWantToNotifyPeopleWhenValidationErrors(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodPost, "/", strings.NewReader("nope"))
	r.Header.Add("Content-Type", page.FormUrlEncoded)

	template := newMockTemplate(t)
	template.
		On("Execute", w, mock.MatchedBy(func(data *doYouWantToNotifyPeopleData) bool {
			return assert.Equal(t, validation.With("yes-no", validation.SelectError{Label: "yesToNotifySomeoneAboutYourLpa"}), data.Errors)
		})).
		Return(nil)

	err := DoYouWantToNotifyPeople(template.Execute, nil)(testAppData, w, r, &page.Lpa{})
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}