package page

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/sessions"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCertificateProviderLogin(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "/?sessionId=session-id&lpaId=lpa-id", nil)

	client := &mockOneLoginClient{}
	client.
		On("AuthCodeURL", "i am random", "i am random", "cy", true).
		Return("http://auth")

	sessionsStore := &mockSessionsStore{}

	session := sessions.NewSession(sessionsStore, "params")

	session.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   600,
		SameSite: http.SameSiteLaxMode,
		HttpOnly: true,
		Secure:   true,
	}
	session.Values = map[any]any{
		"one-login": &OneLoginSession{
			State:               "i am random",
			Nonce:               "i am random",
			Locale:              "cy",
			CertificateProvider: true,
			Identity:            true,
			SessionID:           "session-id",
			LpaID:               "lpa-id",
		},
	}

	sessionsStore.
		On("Save", r, w, session).
		Return(nil)

	CertificateProviderLogin(nil, client, sessionsStore, func(int) string { return "i am random" })(AppData{Lang: Cy, Paths: Paths}, w, r)
	resp := w.Result()

	assert.Equal(t, http.StatusFound, resp.StatusCode)
	assert.Equal(t, "http://auth", resp.Header.Get("Location"))

	mock.AssertExpectationsForObjects(t, client, sessionsStore)
}

func TestCertificateProviderLoginDefaultLocale(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "/?sessionId=session-id&lpaId=lpa-id", nil)

	client := &mockOneLoginClient{}
	client.
		On("AuthCodeURL", "i am random", "i am random", "en", true).
		Return("http://auth")

	sessionsStore := &mockSessionsStore{}

	session := sessions.NewSession(sessionsStore, "params")

	session.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   600,
		SameSite: http.SameSiteLaxMode,
		HttpOnly: true,
		Secure:   true,
	}
	session.Values = map[any]any{
		"one-login": &OneLoginSession{
			State:               "i am random",
			Nonce:               "i am random",
			Locale:              "en",
			CertificateProvider: true,
			Identity:            true,
			SessionID:           "session-id",
			LpaID:               "lpa-id",
		},
	}

	sessionsStore.
		On("Save", r, w, session).
		Return(nil)

	CertificateProviderLogin(nil, client, sessionsStore, func(int) string { return "i am random" })(appData, w, r)
	resp := w.Result()

	assert.Equal(t, http.StatusFound, resp.StatusCode)
	assert.Equal(t, "http://auth", resp.Header.Get("Location"))

	mock.AssertExpectationsForObjects(t, client, sessionsStore)
}

func TestCertificateProviderLoginWhenStoreSaveError(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "/", nil)

	logger := &mockLogger{}
	logger.
		On("Print", expectedError)

	client := &mockOneLoginClient{}
	client.
		On("AuthCodeURL", "i am random", "i am random", "en", true).
		Return("http://auth?locale=en")

	sessionsStore := &mockSessionsStore{}
	sessionsStore.
		On("Save", r, w, mock.Anything).
		Return(expectedError)

	CertificateProviderLogin(logger, client, sessionsStore, func(int) string { return "i am random" })(appData, w, r)
	resp := w.Result()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	mock.AssertExpectationsForObjects(t, logger, client, sessionsStore)
}