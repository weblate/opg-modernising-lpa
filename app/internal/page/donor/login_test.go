package donor

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/sessions"
	"github.com/ministryofjustice/opg-modernising-lpa/internal/sesh"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestLogin(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "/?locale=cy", nil)

	client := newMockOneLoginClient(t)
	client.
		On("AuthCodeURL", "i am random", "i am random", "cy", false).
		Return("http://auth")

	sessionStore := newMockSessionStore(t)

	session := sessions.NewSession(sessionStore, "params")

	session.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   600,
		SameSite: http.SameSiteLaxMode,
		HttpOnly: true,
		Secure:   true,
	}
	session.Values = map[any]any{
		"one-login": &sesh.OneLoginSession{State: "i am random", Nonce: "i am random", Locale: "cy"},
	}

	sessionStore.
		On("Save", r, w, session).
		Return(nil)

	Login(nil, client, sessionStore, func(int) string { return "i am random" })(w, r)
	resp := w.Result()

	assert.Equal(t, http.StatusFound, resp.StatusCode)
	assert.Equal(t, "http://auth", resp.Header.Get("Location"))
}

func TestLoginDefaultLocale(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "/", nil)

	client := newMockOneLoginClient(t)
	client.
		On("AuthCodeURL", "i am random", "i am random", "en", false).
		Return("http://auth")

	sessionStore := newMockSessionStore(t)

	session := sessions.NewSession(sessionStore, "params")

	session.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   600,
		SameSite: http.SameSiteLaxMode,
		HttpOnly: true,
		Secure:   true,
	}
	session.Values = map[any]any{
		"one-login": &sesh.OneLoginSession{State: "i am random", Nonce: "i am random", Locale: "en"},
	}

	sessionStore.
		On("Save", r, w, session).
		Return(nil)

	Login(nil, client, sessionStore, func(int) string { return "i am random" })(w, r)
	resp := w.Result()

	assert.Equal(t, http.StatusFound, resp.StatusCode)
	assert.Equal(t, "http://auth", resp.Header.Get("Location"))
}

func TestLoginWhenStoreSaveError(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "/", nil)

	logger := newMockLogger(t)
	logger.
		On("Print", expectedError)

	client := newMockOneLoginClient(t)
	client.
		On("AuthCodeURL", "i am random", "i am random", "en", false).
		Return("http://auth?locale=en")

	sessionStore := newMockSessionStore(t)
	sessionStore.
		On("Save", r, w, mock.Anything).
		Return(expectedError)

	Login(logger, client, sessionStore, func(int) string { return "i am random" })(w, r)
	resp := w.Result()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
