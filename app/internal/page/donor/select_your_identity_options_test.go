package donor

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/ministryofjustice/opg-modernising-lpa/internal/actor"
	"github.com/ministryofjustice/opg-modernising-lpa/internal/identity"
	"github.com/ministryofjustice/opg-modernising-lpa/internal/page"
	"github.com/ministryofjustice/opg-modernising-lpa/internal/validation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetSelectYourIdentityOptions(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "/", nil)

	donorStore := newMockDonorStore(t)
	donorStore.
		On("Get", r.Context()).
		Return(&page.Lpa{}, nil)

	template := newMockTemplate(t)
	template.
		On("Execute", w, &selectYourIdentityOptionsData{
			App:  testAppData,
			Page: 2,
			Form: &selectYourIdentityOptionsForm{},
		}).
		Return(nil)

	err := SelectYourIdentityOptions(template.Execute, donorStore, 2)(testAppData, w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestGetSelectYourIdentityOptionsWhenStoreErrors(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "/", nil)

	donorStore := newMockDonorStore(t)
	donorStore.
		On("Get", r.Context()).
		Return(&page.Lpa{}, expectedError)

	err := SelectYourIdentityOptions(nil, donorStore, 0)(testAppData, w, r)

	assert.Equal(t, expectedError, err)
}

func TestGetSelectYourIdentityOptionsFromStore(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "/", nil)

	donorStore := newMockDonorStore(t)
	donorStore.
		On("Get", r.Context()).
		Return(&page.Lpa{
			DonorIdentityOption: identity.Passport,
		}, nil)

	template := newMockTemplate(t)
	template.
		On("Execute", w, &selectYourIdentityOptionsData{
			App:  testAppData,
			Form: &selectYourIdentityOptionsForm{Selected: identity.Passport},
		}).
		Return(nil)

	err := SelectYourIdentityOptions(template.Execute, donorStore, 0)(testAppData, w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestGetSelectYourIdentityOptionsWhenTemplateErrors(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "/", nil)

	donorStore := newMockDonorStore(t)
	donorStore.
		On("Get", r.Context()).
		Return(&page.Lpa{}, nil)

	template := newMockTemplate(t)
	template.
		On("Execute", w, mock.Anything).
		Return(expectedError)

	err := SelectYourIdentityOptions(template.Execute, donorStore, 0)(testAppData, w, r)
	resp := w.Result()

	assert.Equal(t, expectedError, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestPostSelectYourIdentityOptions(t *testing.T) {
	form := url.Values{
		"option": {"passport"},
	}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodPost, "/", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", page.FormUrlEncoded)

	donorStore := newMockDonorStore(t)
	donorStore.
		On("Get", r.Context()).
		Return(&page.Lpa{}, nil)
	donorStore.
		On("Put", r.Context(), &page.Lpa{
			DonorIdentityOption: identity.Passport,
			Tasks: page.Tasks{
				ConfirmYourIdentityAndSign: actor.TaskInProgress,
			},
		}).
		Return(nil)

	err := SelectYourIdentityOptions(nil, donorStore, 0)(testAppData, w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusFound, resp.StatusCode)
	assert.Equal(t, "/lpa/lpa-id"+page.Paths.YourChosenIdentityOptions, resp.Header.Get("Location"))
}

func TestPostSelectYourIdentityOptionsNone(t *testing.T) {
	for pageIndex, nextPath := range map[int]string{
		0: "/lpa/lpa-id" + page.Paths.SelectYourIdentityOptions1,
		1: "/lpa/lpa-id" + page.Paths.SelectYourIdentityOptions2,
		2: "/lpa/lpa-id" + page.Paths.TaskList,
	} {
		t.Run(fmt.Sprintf("Page%d", pageIndex), func(t *testing.T) {
			form := url.Values{
				"option": {"none"},
			}

			w := httptest.NewRecorder()
			r, _ := http.NewRequest(http.MethodPost, "/", strings.NewReader(form.Encode()))
			r.Header.Add("Content-Type", page.FormUrlEncoded)

			donorStore := newMockDonorStore(t)
			donorStore.
				On("Get", r.Context()).
				Return(&page.Lpa{}, nil)

			err := SelectYourIdentityOptions(nil, donorStore, pageIndex)(testAppData, w, r)
			resp := w.Result()

			assert.Nil(t, err)
			assert.Equal(t, http.StatusFound, resp.StatusCode)
			assert.Equal(t, nextPath, resp.Header.Get("Location"))
		})
	}
}

func TestPostSelectYourIdentityOptionsWhenStoreErrors(t *testing.T) {
	form := url.Values{
		"option": {"passport"},
	}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodPost, "/", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", page.FormUrlEncoded)

	donorStore := newMockDonorStore(t)
	donorStore.
		On("Get", r.Context()).
		Return(&page.Lpa{}, nil)
	donorStore.
		On("Put", r.Context(), mock.Anything).
		Return(expectedError)

	err := SelectYourIdentityOptions(nil, donorStore, 0)(testAppData, w, r)

	assert.Equal(t, expectedError, err)
}

func TestPostSelectYourIdentityOptionsWhenValidationErrors(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodPost, "/", strings.NewReader(""))
	r.Header.Add("Content-Type", page.FormUrlEncoded)

	donorStore := newMockDonorStore(t)
	donorStore.
		On("Get", r.Context()).
		Return(&page.Lpa{}, nil)

	template := newMockTemplate(t)
	template.
		On("Execute", w, &selectYourIdentityOptionsData{
			App:    testAppData,
			Form:   &selectYourIdentityOptionsForm{},
			Errors: validation.With("option", validation.SelectError{Label: "fromTheListedOptions"}),
		}).
		Return(nil)

	err := SelectYourIdentityOptions(template.Execute, donorStore, 0)(testAppData, w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestReadSelectYourIdentityOptionsForm(t *testing.T) {
	form := url.Values{
		"option": {"passport"},
	}

	r, _ := http.NewRequest(http.MethodPost, "/", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", page.FormUrlEncoded)

	result := readSelectYourIdentityOptionsForm(r)

	assert.Equal(t, identity.Passport, result.Selected)
	assert.False(t, result.None)
}

func TestSelectYourIdentityOptionsFormValidate(t *testing.T) {
	testCases := map[string]struct {
		form      *selectYourIdentityOptionsForm
		errors    validation.List
		pageIndex int
	}{
		"valid": {
			form: &selectYourIdentityOptionsForm{
				Selected: identity.EasyID,
			},
		},
		"none": {
			form: &selectYourIdentityOptionsForm{
				Selected: identity.UnknownOption,
				None:     true,
			},
		},
		"missing": {
			form:   &selectYourIdentityOptionsForm{},
			errors: validation.With("option", validation.SelectError{Label: "fromTheListedOptions"}),
		},
		"invalid": {
			form: &selectYourIdentityOptionsForm{
				Selected: identity.UnknownOption,
			},
			errors: validation.With("option", validation.SelectError{Label: "fromTheListedOptions"}),
		},
		"missing after first page": {
			form:      &selectYourIdentityOptionsForm{},
			errors:    validation.With("option", validation.SelectError{Label: "whichDocumentYouWillUse"}),
			pageIndex: 1,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.errors, tc.form.Validate(tc.pageIndex))
		})
	}
}
