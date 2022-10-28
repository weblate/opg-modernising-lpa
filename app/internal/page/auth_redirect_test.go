package page

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/sessions"
	"github.com/ministryofjustice/opg-modernising-lpa/internal/signin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockAuthRedirectClient struct {
	mock.Mock
}

func (m *mockAuthRedirectClient) Exchange(ctx context.Context, code, nonce string) (string, error) {
	args := m.Called(ctx, code, nonce)
	return args.Get(0).(string), args.Error(1)
}

func (m *mockAuthRedirectClient) UserInfo(jwt string) (signin.UserInfo, error) {
	args := m.Called(jwt)
	return args.Get(0).(signin.UserInfo), args.Error(1)
}

func TestAuthRedirect(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "/?code=auth-code&state=my-state", nil)

	client := &mockAuthRedirectClient{}
	client.
		On("Exchange", r.Context(), "auth-code", "my-nonce").
		Return("a JWT", nil)
	client.
		On("UserInfo", "a JWT").
		Return(signin.UserInfo{Sub: "random", Email: "name@example.com"}, nil)

	sessionsStore := &mockSessionsStore{}

	session := sessions.NewSession(sessionsStore, "session")
	session.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400,
		SameSite: http.SameSiteLaxMode,
		HttpOnly: true,
		Secure:   true,
	}
	session.Values = map[interface{}]interface{}{"sub": "random", "email": "name@example.com"}

	sessionsStore.
		On("Get", r, "params").
		Return(&sessions.Session{Values: map[interface{}]interface{}{"state": "my-state", "nonce": "my-nonce"}}, nil)
	sessionsStore.
		On("Save", r, w, session).
		Return(nil)

	AuthRedirect(nil, client, sessionsStore, true)(w, r)
	resp := w.Result()

	assert.Equal(t, http.StatusFound, resp.StatusCode)
	assert.Equal(t, yourDetailsPath, resp.Header.Get("Location"))
	mock.AssertExpectationsForObjects(t, client, sessionsStore)
}

func TestAuthRedirectSessionMissing(t *testing.T) {
	testCases := map[string]struct {
		url         string
		session     *sessions.Session
		getErr      error
		expectedErr interface{}
	}{
		"missing session": {
			url:         "/?code=auth-code&state=my-state",
			session:     nil,
			getErr:      expectedError,
			expectedErr: expectedError,
		},
		"missing state": {
			url:         "/?code=auth-code&state=my-state",
			session:     &sessions.Session{Values: map[interface{}]interface{}{}},
			expectedErr: "state missing from session or incorrect",
		},
		"missing state from url": {
			url:         "/?code=auth-code",
			session:     &sessions.Session{Values: map[interface{}]interface{}{"state": "my-state"}},
			expectedErr: "state missing from session or incorrect",
		},
		"missing nonce": {
			url:         "/?code=auth-code&state=my-state",
			session:     &sessions.Session{Values: map[interface{}]interface{}{"state": "my-state"}},
			expectedErr: "nonce missing from session",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r, _ := http.NewRequest(http.MethodGet, tc.url, nil)

			logger := &mockLogger{}
			logger.
				On("Print", tc.expectedErr)

			sessionsStore := &mockSessionsStore{}
			sessionsStore.
				On("Get", r, "params").
				Return(tc.session, tc.getErr)

			AuthRedirect(logger, nil, sessionsStore, true)(w, r)
			resp := w.Result()

			assert.Equal(t, http.StatusOK, resp.StatusCode)
			mock.AssertExpectationsForObjects(t, logger, sessionsStore)
		})
	}
}

func TestAuthRedirectStateIncorrect(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "/?code=auth-code&state=hello", nil)

	logger := &mockLogger{}
	logger.
		On("Print", "state missing from session or incorrect")

	sessionsStore := &mockSessionsStore{}
	sessionsStore.
		On("Get", r, "params").
		Return(&sessions.Session{Values: map[interface{}]interface{}{"state": "my-state"}}, nil)

	AuthRedirect(logger, nil, sessionsStore, true)(w, r)
	resp := w.Result()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, logger, sessionsStore)
}

func TestAuthRedirectWhenExchangeErrors(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "/?code=auth-code&state=my-state", nil)

	logger := &mockLogger{}
	logger.
		On("Print", expectedError)

	client := &mockAuthRedirectClient{}
	client.
		On("Exchange", r.Context(), "auth-code", "my-nonce").
		Return("", expectedError)

	sessionsStore := &mockSessionsStore{}
	sessionsStore.
		On("Get", r, "params").
		Return(&sessions.Session{Values: map[interface{}]interface{}{"state": "my-state", "nonce": "my-nonce"}}, nil)

	AuthRedirect(logger, client, sessionsStore, true)(w, r)

	mock.AssertExpectationsForObjects(t, client, logger)
}

func TestAuthRedirectWhenUserInfoError(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "/?code=auth-code&state=my-state", nil)

	logger := &mockLogger{}
	logger.
		On("Print", expectedError)

	client := &mockAuthRedirectClient{}
	client.
		On("Exchange", r.Context(), "auth-code", "my-nonce").
		Return("a JWT", nil)
	client.
		On("UserInfo", "a JWT").
		Return(signin.UserInfo{}, expectedError)

	sessionsStore := &mockSessionsStore{}
	sessionsStore.
		On("Get", r, "params").
		Return(&sessions.Session{Values: map[interface{}]interface{}{"state": "my-state", "nonce": "my-nonce"}}, nil)

	AuthRedirect(logger, client, sessionsStore, true)(w, r)

	mock.AssertExpectationsForObjects(t, client, logger)
}
