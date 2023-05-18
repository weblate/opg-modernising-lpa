package attorney

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gorilla/sessions"
	"github.com/ministryofjustice/opg-modernising-lpa/internal/actor"
	"github.com/ministryofjustice/opg-modernising-lpa/internal/dynamo"
	"github.com/ministryofjustice/opg-modernising-lpa/internal/page"
	"github.com/ministryofjustice/opg-modernising-lpa/internal/sesh"
	"github.com/ministryofjustice/opg-modernising-lpa/internal/validation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func (m *mockSessionStore) ExpectGet(r *http.Request, values map[any]any, err error) {
	session := sessions.NewSession(m, "session")
	session.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400,
		SameSite: http.SameSiteLaxMode,
		HttpOnly: true,
		Secure:   true,
	}
	session.Values = values
	m.
		On("Get", r, "session").
		Return(session, err)
}

func (m *mockSessionStore) ExpectSet(r *http.Request, w http.ResponseWriter, values map[any]any, err error) {
	savedSession := sessions.NewSession(m, "session")
	savedSession.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400,
		SameSite: http.SameSiteLaxMode,
		HttpOnly: true,
		Secure:   true,
	}
	savedSession.Values = values
	m.
		On("Save", r, w, savedSession).
		Return(err)
}

func TestGetEnterReferenceNumber(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "/", nil)

	data := enterReferenceNumberData{
		App:  testAppData,
		Form: &enterReferenceNumberForm{},
	}

	template := newMockTemplate(t)
	template.
		On("Execute", w, data).
		Return(nil)

	err := EnterReferenceNumber(template.Execute, nil, nil, nil)(testAppData, w, r)

	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestGetEnterReferenceNumberOnTemplateError(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "/", nil)

	data := enterReferenceNumberData{
		App:  testAppData,
		Form: &enterReferenceNumberForm{},
	}

	template := newMockTemplate(t)
	template.
		On("Execute", w, data).
		Return(expectedError)

	err := EnterReferenceNumber(template.Execute, nil, nil, nil)(testAppData, w, r)

	resp := w.Result()

	assert.Equal(t, expectedError, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestPostEnterReferenceNumber(t *testing.T) {
	testcases := map[string]struct {
		shareCode     actor.ShareCodeData
		session       *sesh.AttorneySession
		isReplacement bool
	}{
		"attorney": {
			shareCode: actor.ShareCodeData{LpaID: "lpa-id", SessionID: "aGV5", AttorneyID: "attorney-id"},
			session:   &sesh.AttorneySession{Sub: "hey", LpaID: "lpa-id", AttorneyID: "attorney-id"},
		},
		"replacement attorney": {
			shareCode:     actor.ShareCodeData{LpaID: "lpa-id", SessionID: "aGV5", AttorneyID: "attorney-id", IsReplacementAttorney: true},
			session:       &sesh.AttorneySession{Sub: "hey", LpaID: "lpa-id", AttorneyID: "attorney-id", IsReplacementAttorney: true},
			isReplacement: true,
		},
	}

	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			form := url.Values{
				"reference-number": {"a Ref-Number12"},
			}

			w := httptest.NewRecorder()
			r, _ := http.NewRequest(http.MethodPost, "/", strings.NewReader(form.Encode()))
			r.Header.Add("Content-Type", page.FormUrlEncoded)

			shareCodeStore := newMockShareCodeStore(t)
			shareCodeStore.
				On("Get", r.Context(), actor.TypeAttorney, "aRefNumber12").
				Return(tc.shareCode, nil)

			attorneyStore := newMockAttorneyStore(t)
			attorneyStore.
				On("Create", mock.MatchedBy(func(ctx context.Context) bool {
					session, _ := page.SessionDataFromContext(ctx)

					return assert.Equal(t, &page.SessionData{SessionID: "aGV5", LpaID: "lpa-id"}, session)
				}), tc.isReplacement).
				Return(&actor.AttorneyProvidedDetails{}, nil)

			sessionStore := newMockSessionStore(t)
			sessionStore.
				ExpectGet(r,
					map[any]any{"attorney": &sesh.AttorneySession{Sub: "hey"}}, nil)
			sessionStore.
				ExpectSet(r, w, map[any]any{"attorney": tc.session},
					nil)

			err := EnterReferenceNumber(nil, shareCodeStore, sessionStore, attorneyStore)(testAppData, w, r)

			resp := w.Result()

			assert.Nil(t, err)
			assert.Equal(t, http.StatusFound, resp.StatusCode)
			assert.Equal(t, page.Paths.Attorney.CodeOfConduct, resp.Header.Get("Location"))
		})
	}
}

func TestPostEnterReferenceNumberOnDonorStoreError(t *testing.T) {
	form := url.Values{
		"reference-number": {"  aRefNumber12  "},
	}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodPost, "/", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", page.FormUrlEncoded)

	shareCodeStore := newMockShareCodeStore(t)
	shareCodeStore.
		On("Get", r.Context(), actor.TypeAttorney, "aRefNumber12").
		Return(actor.ShareCodeData{LpaID: "lpa-id", SessionID: "aGV5", Identity: true}, expectedError)

	err := EnterReferenceNumber(nil, shareCodeStore, nil, nil)(testAppData, w, r)

	resp := w.Result()

	assert.Equal(t, expectedError, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestPostEnterReferenceNumberOnDonorStoreNotFoundError(t *testing.T) {
	form := url.Values{
		"reference-number": {"a Ref-Number12"},
	}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodPost, "/", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", page.FormUrlEncoded)

	data := enterReferenceNumberData{
		App:    testAppData,
		Form:   &enterReferenceNumberForm{ReferenceNumber: "aRefNumber12", ReferenceNumberRaw: "a Ref-Number12"},
		Errors: validation.With("reference-number", validation.CustomError{Label: "incorrectAttorneyReferenceNumber"}),
	}

	template := newMockTemplate(t)
	template.
		On("Execute", w, data).
		Return(nil)

	shareCodeStore := newMockShareCodeStore(t)
	shareCodeStore.
		On("Get", r.Context(), actor.TypeAttorney, "aRefNumber12").
		Return(actor.ShareCodeData{LpaID: "lpa-id", SessionID: "aGV5", Identity: true}, dynamo.NotFoundError{})

	err := EnterReferenceNumber(template.Execute, shareCodeStore, nil, nil)(testAppData, w, r)

	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestPostEnterReferenceNumberOnSessionGetError(t *testing.T) {
	form := url.Values{
		"reference-number": {"aRefNumber12"},
	}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodPost, "/", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", page.FormUrlEncoded)

	shareCodeStore := newMockShareCodeStore(t)
	shareCodeStore.
		On("Get", r.Context(), actor.TypeAttorney, "aRefNumber12").
		Return(actor.ShareCodeData{LpaID: "lpa-id", SessionID: "aGV5", Identity: true}, nil)

	sessionStore := newMockSessionStore(t)
	sessionStore.
		ExpectGet(r,
			map[any]any{"attorney": &sesh.AttorneySession{Sub: "hey"}}, expectedError)

	err := EnterReferenceNumber(nil, shareCodeStore, sessionStore, nil)(testAppData, w, r)

	assert.Equal(t, expectedError, err)
}

func TestPostEnterReferenceNumberOnSessionSetError(t *testing.T) {
	form := url.Values{
		"reference-number": {"aRefNumber12"},
	}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodPost, "/", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", page.FormUrlEncoded)

	shareCodeStore := newMockShareCodeStore(t)
	shareCodeStore.
		On("Get", r.Context(), actor.TypeAttorney, "aRefNumber12").
		Return(actor.ShareCodeData{LpaID: "lpa-id", SessionID: "aGV5", Identity: true}, nil)

	sessionStore := newMockSessionStore(t)
	sessionStore.
		ExpectGet(r,
			map[any]any{"attorney": &sesh.AttorneySession{Sub: "hey"}}, nil)
	sessionStore.
		ExpectSet(r, w, map[any]any{"attorney": &sesh.AttorneySession{Sub: "hey", LpaID: "lpa-id"}},
			expectedError)

	err := EnterReferenceNumber(nil, shareCodeStore, sessionStore, nil)(testAppData, w, r)

	resp := w.Result()

	assert.Equal(t, expectedError, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestPostEnterReferenceNumberOnAttorneyStoreError(t *testing.T) {
	form := url.Values{
		"reference-number": {"a Ref-Number12"},
	}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodPost, "/", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", page.FormUrlEncoded)

	shareCodeStore := newMockShareCodeStore(t)
	shareCodeStore.
		On("Get", r.Context(), actor.TypeAttorney, "aRefNumber12").
		Return(actor.ShareCodeData{LpaID: "lpa-id", SessionID: "aGV5", Identity: true}, nil)

	attorneyStore := newMockAttorneyStore(t)
	attorneyStore.
		On("Create", mock.Anything, false).
		Return(&actor.AttorneyProvidedDetails{}, expectedError)

	sessionStore := newMockSessionStore(t)
	sessionStore.
		ExpectGet(r,
			map[any]any{"attorney": &sesh.AttorneySession{Sub: "hey"}}, nil)
	sessionStore.
		ExpectSet(r, w, map[any]any{"attorney": &sesh.AttorneySession{Sub: "hey", LpaID: "lpa-id"}},
			nil)

	err := EnterReferenceNumber(nil, shareCodeStore, sessionStore, attorneyStore)(testAppData, w, r)

	resp := w.Result()

	assert.Equal(t, expectedError, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestPostEnterReferenceNumberOnValidationError(t *testing.T) {
	form := url.Values{
		"reference-number": {""},
	}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodPost, "/", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", page.FormUrlEncoded)

	data := enterReferenceNumberData{
		App:    testAppData,
		Form:   &enterReferenceNumberForm{},
		Errors: validation.With("reference-number", validation.EnterError{Label: "twelveCharactersAttorneyReferenceNumber"}),
	}

	template := newMockTemplate(t)
	template.
		On("Execute", w, data).
		Return(nil)

	err := EnterReferenceNumber(template.Execute, nil, nil, nil)(testAppData, w, r)

	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestValidateEnterReferenceNumberForm(t *testing.T) {
	testCases := map[string]struct {
		form   *enterReferenceNumberForm
		errors validation.List
	}{
		"valid": {
			form:   &enterReferenceNumberForm{ReferenceNumber: "abcdef123456"},
			errors: nil,
		},
		"too short": {
			form: &enterReferenceNumberForm{ReferenceNumber: "1"},
			errors: validation.With("reference-number", validation.StringLengthError{
				Label:  "attorneyReferenceNumberMustBeTwelveCharacters",
				Length: 12,
			}),
		},
		"too long": {
			form: &enterReferenceNumberForm{ReferenceNumber: "abcdef1234567"},
			errors: validation.With("reference-number", validation.StringLengthError{
				Label:  "attorneyReferenceNumberMustBeTwelveCharacters",
				Length: 12,
			}),
		},
		"empty": {
			form: &enterReferenceNumberForm{},
			errors: validation.With("reference-number", validation.EnterError{
				Label: "twelveCharactersAttorneyReferenceNumber",
			}),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.errors, tc.form.Validate())
		})
	}

}
