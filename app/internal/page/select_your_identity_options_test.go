package page

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/ministryofjustice/opg-modernising-lpa/internal/validation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetSelectYourIdentityOptions(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "/", nil)

	lpaStore := &mockLpaStore{}
	lpaStore.
		On("Get", r.Context()).
		Return(&Lpa{}, nil)

	template := &mockTemplate{}
	template.
		On("Func", w, &selectYourIdentityOptionsData{
			App:  appData,
			Page: 2,
			Form: &selectYourIdentityOptionsForm{},
		}).
		Return(nil)

	err := SelectYourIdentityOptions(template.Func, lpaStore, 2)(appData, w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, template, lpaStore)
}

func TestGetSelectYourIdentityOptionsWhenStoreErrors(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "/", nil)

	lpaStore := &mockLpaStore{}
	lpaStore.
		On("Get", r.Context()).
		Return(&Lpa{}, expectedError)

	err := SelectYourIdentityOptions(nil, lpaStore, 0)(appData, w, r)

	assert.Equal(t, expectedError, err)
	mock.AssertExpectationsForObjects(t, lpaStore)
}

func TestGetSelectYourIdentityOptionsFromStore(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "/", nil)

	lpaStore := &mockLpaStore{}
	lpaStore.
		On("Get", r.Context()).
		Return(&Lpa{
			IdentityOption: Passport,
		}, nil)

	template := &mockTemplate{}
	template.
		On("Func", w, &selectYourIdentityOptionsData{
			App:  appData,
			Form: &selectYourIdentityOptionsForm{Selected: Passport},
		}).
		Return(nil)

	err := SelectYourIdentityOptions(template.Func, lpaStore, 0)(appData, w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, template)
}

func TestGetSelectYourIdentityOptionsWhenTemplateErrors(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "/", nil)

	lpaStore := &mockLpaStore{}
	lpaStore.
		On("Get", r.Context()).
		Return(&Lpa{}, nil)

	template := &mockTemplate{}
	template.
		On("Func", w, mock.Anything).
		Return(expectedError)

	err := SelectYourIdentityOptions(template.Func, lpaStore, 0)(appData, w, r)
	resp := w.Result()

	assert.Equal(t, expectedError, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, template)
}

func TestPostSelectYourIdentityOptions(t *testing.T) {
	form := url.Values{
		"option": {"passport"},
	}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodPost, "/", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)

	lpaStore := &mockLpaStore{}
	lpaStore.
		On("Get", r.Context()).
		Return(&Lpa{}, nil)
	lpaStore.
		On("Put", r.Context(), &Lpa{
			IdentityOption: Passport,
			Tasks: Tasks{
				ConfirmYourIdentityAndSign: TaskInProgress,
			},
		}).
		Return(nil)

	err := SelectYourIdentityOptions(nil, lpaStore, 0)(appData, w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusFound, resp.StatusCode)
	assert.Equal(t, "/lpa/lpa-id"+Paths.YourChosenIdentityOptions, resp.Header.Get("Location"))
	mock.AssertExpectationsForObjects(t, lpaStore)
}

func TestPostSelectYourIdentityOptionsNone(t *testing.T) {
	for page, nextPath := range map[int]string{
		0: "/lpa/lpa-id" + Paths.SelectYourIdentityOptions1,
		1: "/lpa/lpa-id" + Paths.SelectYourIdentityOptions2,
		2: "/lpa/lpa-id" + Paths.TaskList,
	} {
		t.Run(fmt.Sprintf("Page%d", page), func(t *testing.T) {
			form := url.Values{
				"option": {"none"},
			}

			w := httptest.NewRecorder()
			r, _ := http.NewRequest(http.MethodPost, "/", strings.NewReader(form.Encode()))
			r.Header.Add("Content-Type", formUrlEncoded)

			lpaStore := &mockLpaStore{}
			lpaStore.
				On("Get", r.Context()).
				Return(&Lpa{}, nil)

			err := SelectYourIdentityOptions(nil, lpaStore, page)(appData, w, r)
			resp := w.Result()

			assert.Nil(t, err)
			assert.Equal(t, http.StatusFound, resp.StatusCode)
			assert.Equal(t, nextPath, resp.Header.Get("Location"))
			mock.AssertExpectationsForObjects(t, lpaStore)
		})
	}
}

func TestPostSelectYourIdentityOptionsWhenStoreErrors(t *testing.T) {
	form := url.Values{
		"option": {"passport"},
	}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodPost, "/", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)

	lpaStore := &mockLpaStore{}
	lpaStore.
		On("Get", r.Context()).
		Return(&Lpa{}, nil)
	lpaStore.
		On("Put", r.Context(), mock.Anything).
		Return(expectedError)

	err := SelectYourIdentityOptions(nil, lpaStore, 0)(appData, w, r)

	assert.Equal(t, expectedError, err)
	mock.AssertExpectationsForObjects(t, lpaStore)
}

func TestPostSelectYourIdentityOptionsWhenValidationErrors(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodPost, "/", strings.NewReader(""))
	r.Header.Add("Content-Type", formUrlEncoded)

	lpaStore := &mockLpaStore{}
	lpaStore.
		On("Get", r.Context()).
		Return(&Lpa{}, nil)

	template := &mockTemplate{}
	template.
		On("Func", w, &selectYourIdentityOptionsData{
			App:    appData,
			Form:   &selectYourIdentityOptionsForm{},
			Errors: validation.With("option", validation.SelectError{Label: "fromTheListedOptions"}),
		}).
		Return(nil)

	err := SelectYourIdentityOptions(template.Func, lpaStore, 0)(appData, w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, template)
}

func TestReadSelectYourIdentityOptionsForm(t *testing.T) {
	form := url.Values{
		"option": {"passport"},
	}

	r, _ := http.NewRequest(http.MethodPost, "/", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)

	result := readSelectYourIdentityOptionsForm(r)

	assert.Equal(t, Passport, result.Selected)
	assert.False(t, result.None)
}

func TestSelectYourIdentityOptionsFormValidate(t *testing.T) {
	testCases := map[string]struct {
		form   *selectYourIdentityOptionsForm
		errors validation.List
	}{
		"valid": {
			form: &selectYourIdentityOptionsForm{
				Selected: EasyID,
			},
		},
		"none": {
			form: &selectYourIdentityOptionsForm{
				Selected: IdentityOptionUnknown,
				None:     true,
			},
		},
		"missing": {
			form:   &selectYourIdentityOptionsForm{},
			errors: validation.With("option", validation.SelectError{Label: "fromTheListedOptions"}),
		},
		"invalid": {
			form: &selectYourIdentityOptionsForm{
				Selected: IdentityOptionUnknown,
			},
			errors: validation.With("option", validation.SelectError{Label: "fromTheListedOptions"}),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.errors, tc.form.Validate())
		})
	}
}
