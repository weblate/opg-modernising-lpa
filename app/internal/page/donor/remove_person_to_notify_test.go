package donor

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/ministryofjustice/opg-modernising-lpa/internal/actor"
	"github.com/ministryofjustice/opg-modernising-lpa/internal/page"
	"github.com/ministryofjustice/opg-modernising-lpa/internal/place"
	"github.com/ministryofjustice/opg-modernising-lpa/internal/validation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetRemovePersonToNotify(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "/?id=123", nil)

	logger := newMockLogger(t)

	personToNotify := actor.PersonToNotify{
		ID: "123",
		Address: place.Address{
			Line1: "1 Road way",
		},
	}

	template := newMockTemplate(t)
	template.
		On("Execute", w, &removePersonToNotifyData{
			App:            testAppData,
			PersonToNotify: personToNotify,
			Errors:         nil,
			Form:           &removePersonToNotifyForm{},
		}).
		Return(nil)

	lpaStore := newMockLpaStore(t)
	lpaStore.
		On("Get", r.Context()).
		Return(&page.Lpa{PeopleToNotify: actor.PeopleToNotify{personToNotify}}, nil)

	err := RemovePersonToNotify(logger, template.Execute, lpaStore)(testAppData, w, r)

	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestGetRemovePersonToNotifyErrorOnStore(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "/?id=123", nil)

	logger := newMockLogger(t)
	logger.
		On("Print", "error getting lpa from store: err").
		Return(nil)

	template := newMockTemplate(t)

	lpaStore := newMockLpaStore(t)
	lpaStore.
		On("Get", r.Context()).
		Return(&page.Lpa{}, expectedError)

	err := RemovePersonToNotify(logger, template.Execute, lpaStore)(testAppData, w, r)

	resp := w.Result()

	assert.Equal(t, expectedError, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestGetRemovePersonToNotifyAttorneyDoesNotExist(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "/?id=invalid-id", nil)

	logger := newMockLogger(t)

	template := newMockTemplate(t)

	personToNotify := actor.PersonToNotify{
		ID: "123",
		Address: place.Address{
			Line1: "1 Road way",
		},
	}

	lpaStore := newMockLpaStore(t)
	lpaStore.
		On("Get", r.Context()).
		Return(&page.Lpa{PeopleToNotify: actor.PeopleToNotify{personToNotify}}, nil)

	err := RemovePersonToNotify(logger, template.Execute, lpaStore)(testAppData, w, r)

	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusFound, resp.StatusCode)
	assert.Equal(t, "/lpa/lpa-id"+page.Paths.ChoosePeopleToNotifySummary, resp.Header.Get("Location"))
}

func TestPostRemovePersonToNotify(t *testing.T) {
	form := url.Values{
		"remove-person-to-notify": {"yes"},
	}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodPost, "/?id=without-address", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", page.FormUrlEncoded)

	logger := newMockLogger(t)
	template := newMockTemplate(t)

	personToNotifyWithAddress := actor.PersonToNotify{
		ID: "with-address",
		Address: place.Address{
			Line1: "1 Road way",
		},
	}

	personToNotifyWithoutAddress := actor.PersonToNotify{
		ID:      "without-address",
		Address: place.Address{},
	}

	lpaStore := newMockLpaStore(t)
	lpaStore.
		On("Get", r.Context()).
		Return(&page.Lpa{PeopleToNotify: actor.PeopleToNotify{personToNotifyWithoutAddress, personToNotifyWithAddress}}, nil)
	lpaStore.
		On("Put", r.Context(), &page.Lpa{PeopleToNotify: actor.PeopleToNotify{personToNotifyWithAddress}}).
		Return(nil)

	err := RemovePersonToNotify(logger, template.Execute, lpaStore)(testAppData, w, r)

	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusFound, resp.StatusCode)
	assert.Equal(t, "/lpa/lpa-id"+page.Paths.ChoosePeopleToNotifySummary, resp.Header.Get("Location"))
}

func TestPostRemovePersonToNotifyWithFormValueNo(t *testing.T) {
	form := url.Values{
		"remove-person-to-notify": {"no"},
	}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodPost, "/?id=without-address", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", page.FormUrlEncoded)

	logger := newMockLogger(t)
	template := newMockTemplate(t)

	personToNotifyWithAddress := actor.PersonToNotify{
		ID: "with-address",
		Address: place.Address{
			Line1: "1 Road way",
		},
	}

	personToNotifyWithoutAddress := actor.PersonToNotify{
		ID:      "without-address",
		Address: place.Address{},
	}

	lpaStore := newMockLpaStore(t)
	lpaStore.
		On("Get", r.Context()).
		Return(&page.Lpa{PeopleToNotify: actor.PeopleToNotify{personToNotifyWithoutAddress, personToNotifyWithAddress}}, nil)

	err := RemovePersonToNotify(logger, template.Execute, lpaStore)(testAppData, w, r)

	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusFound, resp.StatusCode)
	assert.Equal(t, "/lpa/lpa-id"+page.Paths.ChoosePeopleToNotifySummary, resp.Header.Get("Location"))
}

func TestPostRemovePersonToNotifyErrorOnPutStore(t *testing.T) {
	form := url.Values{
		"remove-person-to-notify": {"yes"},
	}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodPost, "/?id=without-address", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", page.FormUrlEncoded)

	template := newMockTemplate(t)

	logger := newMockLogger(t)
	logger.
		On("Print", "error removing PersonToNotify from LPA: err").
		Return(nil)

	personToNotifyWithAddress := actor.PersonToNotify{
		ID: "with-address",
		Address: place.Address{
			Line1: "1 Road way",
		},
	}

	personToNotifyWithoutAddress := actor.PersonToNotify{
		ID:      "without-address",
		Address: place.Address{},
	}

	lpaStore := newMockLpaStore(t)
	lpaStore.
		On("Get", r.Context()).
		Return(&page.Lpa{PeopleToNotify: actor.PeopleToNotify{personToNotifyWithoutAddress, personToNotifyWithAddress}}, nil)
	lpaStore.
		On("Put", r.Context(), &page.Lpa{PeopleToNotify: actor.PeopleToNotify{personToNotifyWithAddress}}).
		Return(expectedError)

	err := RemovePersonToNotify(logger, template.Execute, lpaStore)(testAppData, w, r)

	resp := w.Result()

	assert.Equal(t, expectedError, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestRemovePersonToNotifyFormValidation(t *testing.T) {
	form := url.Values{
		"remove-person-to-notify": {""},
	}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodPost, "/?id=without-address", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", page.FormUrlEncoded)

	personToNotifyWithoutAddress := actor.PersonToNotify{
		ID:      "without-address",
		Address: place.Address{},
	}

	lpaStore := newMockLpaStore(t)
	lpaStore.
		On("Get", r.Context()).
		Return(&page.Lpa{PeopleToNotify: actor.PeopleToNotify{personToNotifyWithoutAddress}}, nil)

	validationError := validation.With("remove-person-to-notify", validation.SelectError{Label: "removePersonToNotify"})

	template := newMockTemplate(t)
	template.
		On("Execute", w, mock.MatchedBy(func(data *removePersonToNotifyData) bool {
			return assert.Equal(t, validationError, data.Errors)
		})).
		Return(nil)

	err := RemovePersonToNotify(nil, template.Execute, lpaStore)(testAppData, w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestRemovePersonToNotifyRemoveLastPerson(t *testing.T) {
	form := url.Values{
		"remove-person-to-notify": {"yes"},
	}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodPost, "/?id=without-address", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", page.FormUrlEncoded)

	logger := newMockLogger(t)
	template := newMockTemplate(t)

	personToNotifyWithoutAddress := actor.PersonToNotify{
		ID:      "without-address",
		Address: place.Address{},
	}

	lpaStore := newMockLpaStore(t)
	lpaStore.
		On("Get", r.Context()).
		Return(&page.Lpa{
			PeopleToNotify: actor.PeopleToNotify{personToNotifyWithoutAddress},
			Tasks:          page.Tasks{YourDetails: actor.TaskCompleted, ChooseAttorneys: actor.TaskCompleted, PeopleToNotify: actor.TaskCompleted},
		}, nil)
	lpaStore.
		On("Put", r.Context(), &page.Lpa{
			PeopleToNotify: actor.PeopleToNotify{},
			Tasks:          page.Tasks{YourDetails: actor.TaskCompleted, ChooseAttorneys: actor.TaskCompleted, PeopleToNotify: actor.TaskNotStarted},
		}).
		Return(nil)

	err := RemovePersonToNotify(logger, template.Execute, lpaStore)(testAppData, w, r)

	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusFound, resp.StatusCode)
	assert.Equal(t, "/lpa/lpa-id"+page.Paths.ChoosePeopleToNotifySummary, resp.Header.Get("Location"))
}

func TestRemovePersonToNotifyFormValidate(t *testing.T) {
	testCases := map[string]struct {
		form   *removePersonToNotifyForm
		errors validation.List
	}{
		"valid - yes": {
			form: &removePersonToNotifyForm{
				RemovePersonToNotify: "yes",
			},
		},
		"valid - no": {
			form: &removePersonToNotifyForm{
				RemovePersonToNotify: "no",
			},
		},
		"missing-value": {
			form: &removePersonToNotifyForm{
				RemovePersonToNotify: "",
			},
			errors: validation.With("remove-person-to-notify", validation.SelectError{Label: "removePersonToNotify"}),
		},
		"unexpected-value": {
			form: &removePersonToNotifyForm{
				RemovePersonToNotify: "not expected",
			},
			errors: validation.With("remove-person-to-notify", validation.SelectError{Label: "removePersonToNotify"}),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.errors, tc.form.Validate())
		})
	}
}
