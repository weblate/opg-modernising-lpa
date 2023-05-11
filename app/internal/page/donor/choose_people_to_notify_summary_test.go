package donor

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/ministryofjustice/opg-modernising-lpa/internal/actor"
	"github.com/ministryofjustice/opg-modernising-lpa/internal/page"
	"github.com/ministryofjustice/opg-modernising-lpa/internal/validation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetChoosePeopleToNotifySummary(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "/", nil)

	lpa := &page.Lpa{PeopleToNotify: actor.PeopleToNotify{{}}}

	lpaStore := newMockLpaStore(t)
	lpaStore.
		On("Get", r.Context()).
		Return(lpa, nil)

	template := newMockTemplate(t)
	template.
		On("Execute", w, &choosePeopleToNotifySummaryData{
			App:  testAppData,
			Lpa:  lpa,
			Form: &choosePeopleToNotifySummaryForm{},
		}).
		Return(nil)

	err := ChoosePeopleToNotifySummary(nil, template.Execute, lpaStore)(testAppData, w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestGetChoosePeopleToNotifySummaryWhenNoPeopleToNotify(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "/", nil)

	lpaStore := newMockLpaStore(t)
	lpaStore.
		On("Get", r.Context()).
		Return(&page.Lpa{
			Tasks: page.Tasks{
				YourDetails:                page.TaskCompleted,
				ChooseAttorneys:            page.TaskCompleted,
				ChooseReplacementAttorneys: page.TaskCompleted,
				WhenCanTheLpaBeUsed:        page.TaskCompleted,
				Restrictions:               page.TaskCompleted,
				CertificateProvider:        page.TaskCompleted,
			},
		}, nil)

	err := ChoosePeopleToNotifySummary(nil, nil, lpaStore)(testAppData, w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusFound, resp.StatusCode)
	assert.Equal(t, "/lpa/lpa-id"+page.Paths.DoYouWantToNotifyPeople, resp.Header.Get("Location"))
}

func TestGetChoosePeopleToNotifySummaryWhenStoreErrors(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "/", nil)

	lpaStore := newMockLpaStore(t)
	lpaStore.
		On("Get", r.Context()).
		Return(&page.Lpa{PeopleToNotify: actor.PeopleToNotify{{}}}, expectedError)

	logger := &mockLogger{}
	logger.
		On("Print", "error getting lpa from store: err").
		Return(nil)

	err := ChoosePeopleToNotifySummary(logger, nil, lpaStore)(testAppData, w, r)
	resp := w.Result()

	assert.Equal(t, expectedError, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestPostChoosePeopleToNotifySummaryAddPersonToNotify(t *testing.T) {
	form := url.Values{
		"add-person-to-notify": {"yes"},
	}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodPost, "/", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", page.FormUrlEncoded)

	lpaStore := newMockLpaStore(t)
	lpaStore.
		On("Get", r.Context()).
		Return(&page.Lpa{PeopleToNotify: actor.PeopleToNotify{{ID: "123"}}}, nil)

	err := ChoosePeopleToNotifySummary(nil, nil, lpaStore)(testAppData, w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusFound, resp.StatusCode)
	assert.Equal(t, fmt.Sprintf("/lpa/lpa-id%s?addAnother=1", page.Paths.ChoosePeopleToNotify), resp.Header.Get("Location"))
}

func TestPostChoosePeopleToNotifySummaryNoFurtherPeopleToNotify(t *testing.T) {
	form := url.Values{
		"add-person-to-notify": {"no"},
	}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodPost, "/", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", page.FormUrlEncoded)

	lpaStore := newMockLpaStore(t)
	lpaStore.
		On("Get", r.Context()).
		Return(&page.Lpa{
			PeopleToNotify: actor.PeopleToNotify{{ID: "123"}},
			Tasks: page.Tasks{
				YourDetails:                page.TaskCompleted,
				ChooseAttorneys:            page.TaskCompleted,
				ChooseReplacementAttorneys: page.TaskCompleted,
				WhenCanTheLpaBeUsed:        page.TaskCompleted,
				Restrictions:               page.TaskCompleted,
				CertificateProvider:        page.TaskCompleted,
				PeopleToNotify:             page.TaskCompleted,
			},
		}, nil)

	err := ChoosePeopleToNotifySummary(nil, nil, lpaStore)(testAppData, w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusFound, resp.StatusCode)
	assert.Equal(t, "/lpa/lpa-id"+page.Paths.CheckYourLpa, resp.Header.Get("Location"))
}

func TestPostChoosePeopleToNotifySummaryFormValidation(t *testing.T) {
	form := url.Values{
		"add-person-to-notify": {""},
	}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodPost, "/", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", page.FormUrlEncoded)

	lpaStore := newMockLpaStore(t)
	lpaStore.
		On("Get", r.Context()).
		Return(&page.Lpa{PeopleToNotify: actor.PeopleToNotify{{}}}, nil)

	validationError := validation.With("add-person-to-notify", validation.SelectError{Label: "yesToAddAnotherPersonToNotify"})

	template := newMockTemplate(t)
	template.
		On("Execute", w, mock.MatchedBy(func(data *choosePeopleToNotifySummaryData) bool {
			return assert.Equal(t, validationError, data.Errors)
		})).
		Return(nil)

	err := ChoosePeopleToNotifySummary(nil, template.Execute, lpaStore)(testAppData, w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestChoosePeopleToNotifySummaryFormValidate(t *testing.T) {
	testCases := map[string]struct {
		form   *choosePeopleToNotifySummaryForm
		errors validation.List
	}{
		"yes": {
			form: &choosePeopleToNotifySummaryForm{
				AddPersonToNotify: "yes",
			},
		},
		"no": {
			form: &choosePeopleToNotifySummaryForm{
				AddPersonToNotify: "no",
			},
		},
		"missing": {
			form:   &choosePeopleToNotifySummaryForm{},
			errors: validation.With("add-person-to-notify", validation.SelectError{Label: "yesToAddAnotherPersonToNotify"}),
		},
		"invalid": {
			form: &choosePeopleToNotifySummaryForm{
				AddPersonToNotify: "what",
			},
			errors: validation.With("add-person-to-notify", validation.SelectError{Label: "yesToAddAnotherPersonToNotify"}),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.errors, tc.form.Validate())
		})
	}
}
