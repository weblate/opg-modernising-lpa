package page

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetRemoveAttorney(t *testing.T) {
	w := httptest.NewRecorder()

	logger := &mockLogger{}

	template := &mockTemplate{}
	template.
		On("Func", w, &removeAttorneyData{
			App:      appData,
			Attorney: attorneyWithAddress,
			Errors:   nil,
			Form:     removeAttorneyForm{},
		}).
		Return(nil)

	lpaStore := &mockLpaStore{}
	lpaStore.
		On("Get", mock.Anything, "session-id").
		Return(&Lpa{Attorneys: []Attorney{attorneyWithAddress}}, nil)

	r, _ := http.NewRequest(http.MethodGet, "/?id=123", nil)

	err := RemoveAttorney(logger, template.Func, lpaStore)(appData, w, r)

	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, lpaStore, template)
}

func TestGetRemoveAttorneyErrorOnStore(t *testing.T) {
	w := httptest.NewRecorder()

	logger := &mockLogger{}
	logger.
		On("Print", "error getting lpa from store: err").
		Return(nil)

	template := &mockTemplate{}

	lpaStore := &mockLpaStore{}
	lpaStore.
		On("Get", mock.Anything, "session-id").
		Return(&Lpa{}, expectedError)

	r, _ := http.NewRequest(http.MethodGet, "/?id=123", nil)

	err := RemoveAttorney(logger, template.Func, lpaStore)(appData, w, r)

	resp := w.Result()

	assert.Equal(t, expectedError, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, lpaStore, logger)
}

func TestGetRemoveAttorneyAttorneyDoesNotExist(t *testing.T) {
	w := httptest.NewRecorder()

	logger := &mockLogger{}

	template := &mockTemplate{}

	lpaStore := &mockLpaStore{}
	lpaStore.
		On("Get", mock.Anything, "session-id").
		Return(&Lpa{Attorneys: []Attorney{attorneyWithAddress}}, nil)

	r, _ := http.NewRequest(http.MethodGet, "/?id=invalid-id", nil)

	err := RemoveAttorney(logger, template.Func, lpaStore)(appData, w, r)

	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusFound, resp.StatusCode)
	assert.Equal(t, "/choose-attorneys-summary", resp.Header.Get("Location"))
	mock.AssertExpectationsForObjects(t, lpaStore)
}

func TestPostRemoveAttorney(t *testing.T) {
	w := httptest.NewRecorder()

	logger := &mockLogger{}
	template := &mockTemplate{}

	lpaStore := &mockLpaStore{}
	lpaStore.
		On("Get", mock.Anything, "session-id").
		Return(&Lpa{Attorneys: []Attorney{attorneyWithoutAddress, attorneyWithAddress}}, nil)
	lpaStore.
		On("Put", mock.Anything, "session-id", &Lpa{Attorneys: []Attorney{attorneyWithAddress}}).
		Return(nil)

	form := url.Values{
		"remove-attorney": {"yes"},
	}

	r, _ := http.NewRequest(http.MethodPost, "/?id=without-address", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)

	err := RemoveAttorney(logger, template.Func, lpaStore)(appData, w, r)

	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusFound, resp.StatusCode)
	assert.Equal(t, "/choose-attorneys-summary", resp.Header.Get("Location"))
	mock.AssertExpectationsForObjects(t, lpaStore, template)
}

func TestPostRemoveAttorneyWithFormValueNo(t *testing.T) {
	w := httptest.NewRecorder()

	logger := &mockLogger{}
	template := &mockTemplate{}

	lpaStore := &mockLpaStore{}
	lpaStore.
		On("Get", mock.Anything, "session-id").
		Return(&Lpa{Attorneys: []Attorney{attorneyWithoutAddress, attorneyWithAddress}}, nil)

	form := url.Values{
		"remove-attorney": {"no"},
	}

	r, _ := http.NewRequest(http.MethodPost, "/?id=without-address", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)

	err := RemoveAttorney(logger, template.Func, lpaStore)(appData, w, r)

	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusFound, resp.StatusCode)
	assert.Equal(t, "/choose-attorneys-summary", resp.Header.Get("Location"))
	mock.AssertExpectationsForObjects(t, lpaStore, template)
}

func TestPostRemoveAttorneyErrorOnPutStore(t *testing.T) {
	w := httptest.NewRecorder()

	template := &mockTemplate{}

	logger := &mockLogger{}
	logger.
		On("Print", "error removing Attorney from LPA: err").
		Return(nil)

	lpaStore := &mockLpaStore{}
	lpaStore.
		On("Get", mock.Anything, "session-id").
		Return(&Lpa{Attorneys: []Attorney{attorneyWithoutAddress, attorneyWithAddress}}, nil)
	lpaStore.
		On("Put", mock.Anything, "session-id", &Lpa{Attorneys: []Attorney{attorneyWithAddress}}).
		Return(expectedError)

	form := url.Values{
		"remove-attorney": {"yes"},
	}

	r, _ := http.NewRequest(http.MethodPost, "/?id=without-address", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)

	err := RemoveAttorney(logger, template.Func, lpaStore)(appData, w, r)

	resp := w.Result()

	assert.Equal(t, expectedError, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, lpaStore, template, logger)
}
