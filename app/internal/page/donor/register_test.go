package donor

import (
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/sessions"
	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-modernising-lpa/internal/identity"
	"github.com/ministryofjustice/opg-modernising-lpa/internal/notify"
	"github.com/ministryofjustice/opg-modernising-lpa/internal/onelogin"
	"github.com/ministryofjustice/opg-modernising-lpa/internal/page"
	"github.com/ministryofjustice/opg-modernising-lpa/internal/pay"
	"github.com/ministryofjustice/opg-modernising-lpa/internal/place"
	"github.com/ministryofjustice/opg-modernising-lpa/internal/sesh"
	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	mux := http.NewServeMux()
	Register(mux, &log.Logger{}, template.Templates{}, nil, nil, &onelogin.Client{}, &place.Client{}, "http://public.url", &pay.Client{}, &identity.YotiClient{}, "yoti-scenario-id", &notify.Client{}, nil, nil, nil)

	assert.Implements(t, (*http.Handler)(nil), mux)
}

func TestMakeHandle(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "/path?a=b", nil)

	sessionStore := newMockSessionStore(t)
	sessionStore.
		On("Get", r, "session").
		Return(&sessions.Session{Values: map[interface{}]interface{}{"donor": &sesh.DonorSession{Sub: "random"}}}, nil)

	mux := http.NewServeMux()
	handle := makeHandle(mux, sessionStore, None, nil)
	handle("/path", RequireSession|CanGoBack, func(appData page.AppData, hw http.ResponseWriter, hr *http.Request) error {
		assert.Equal(t, page.AppData{
			Page:      "/path",
			CanGoBack: true,
			SessionID: "cmFuZG9t",
			IsDonor:   true,
		}, appData)
		assert.Equal(t, w, hw)
		assert.Equal(t, &page.SessionData{SessionID: "cmFuZG9t"}, page.SessionDataFromContext(hr.Context()))
		hw.WriteHeader(http.StatusTeapot)
		return nil
	})

	mux.ServeHTTP(w, r)
	resp := w.Result()

	assert.Equal(t, http.StatusTeapot, resp.StatusCode)
}

func TestMakeHandleExistingSessionData(t *testing.T) {
	ctx := page.ContextWithSessionData(context.Background(), &page.SessionData{LpaID: "123"})
	w := httptest.NewRecorder()
	r, _ := http.NewRequestWithContext(ctx, http.MethodGet, "/path?a=b", nil)

	sessionStore := newMockSessionStore(t)
	sessionStore.
		On("Get", r, "session").
		Return(&sessions.Session{Values: map[interface{}]interface{}{"donor": &sesh.DonorSession{Sub: "random"}}}, nil)

	mux := http.NewServeMux()
	handle := makeHandle(mux, sessionStore, None, nil)
	handle("/path", RequireSession|CanGoBack, func(appData page.AppData, hw http.ResponseWriter, hr *http.Request) error {
		assert.Equal(t, page.AppData{
			Page:      "/path",
			SessionID: "cmFuZG9t",
			CanGoBack: true,
			LpaID:     "123",
			IsDonor:   true,
		}, appData)
		assert.Equal(t, w, hw)
		assert.Equal(t, &page.SessionData{LpaID: "123", SessionID: "cmFuZG9t"}, page.SessionDataFromContext(hr.Context()))
		hw.WriteHeader(http.StatusTeapot)
		return nil
	})

	mux.ServeHTTP(w, r)
	resp := w.Result()

	assert.Equal(t, http.StatusTeapot, resp.StatusCode)
}

func TestMakeHandleErrors(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "/path", nil)

	errorHandler := newMockErrorHandler(t)
	errorHandler.
		On("Execute", w, r, expectedError)

	sessionStore := newMockSessionStore(t)
	sessionStore.
		On("Get", r, "session").
		Return(&sessions.Session{Values: map[interface{}]interface{}{"donor": &sesh.DonorSession{Sub: "random"}}}, nil)

	mux := http.NewServeMux()
	handle := makeHandle(mux, sessionStore, None, errorHandler.Execute)
	handle("/path", RequireSession, func(appData page.AppData, hw http.ResponseWriter, hr *http.Request) error {
		return expectedError
	})

	mux.ServeHTTP(w, r)
}

func TestMakeHandleSessionError(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "/path", nil)

	sessionStore := newMockSessionStore(t)
	sessionStore.
		On("Get", r, "session").
		Return(&sessions.Session{}, expectedError)

	mux := http.NewServeMux()
	handle := makeHandle(mux, sessionStore, None, nil)
	handle("/path", RequireSession, func(appData page.AppData, hw http.ResponseWriter, hr *http.Request) error { return nil })

	mux.ServeHTTP(w, r)
	resp := w.Result()

	assert.Equal(t, http.StatusFound, resp.StatusCode)
	assert.Equal(t, page.Paths.Start, resp.Header.Get("Location"))
}

func TestMakeHandleSessionMissing(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "/path", nil)

	sessionStore := newMockSessionStore(t)
	sessionStore.
		On("Get", r, "session").
		Return(&sessions.Session{Values: map[interface{}]interface{}{}}, nil)

	mux := http.NewServeMux()
	handle := makeHandle(mux, sessionStore, None, nil)
	handle("/path", RequireSession, func(appData page.AppData, hw http.ResponseWriter, hr *http.Request) error { return nil })

	mux.ServeHTTP(w, r)
	resp := w.Result()

	assert.Equal(t, http.StatusFound, resp.StatusCode)
	assert.Equal(t, page.Paths.Start, resp.Header.Get("Location"))
}

func TestMakeHandleNoSessionRequired(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "/path", nil)

	mux := http.NewServeMux()
	handle := makeHandle(mux, nil, None, nil)
	handle("/path", None, func(appData page.AppData, hw http.ResponseWriter, hr *http.Request) error {
		assert.Equal(t, page.AppData{
			Page:    "/path",
			IsDonor: true,
		}, appData)
		assert.Equal(t, w, hw)
		assert.Equal(t, r.WithContext(page.ContextWithAppData(r.Context(), page.AppData{Page: "/path", IsDonor: true})), hr)
		hw.WriteHeader(http.StatusTeapot)
		return nil
	})

	mux.ServeHTTP(w, r)
	resp := w.Result()

	assert.Equal(t, http.StatusTeapot, resp.StatusCode)
}

func TestRouteToLpa(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "/lpa/123/somewhere%2Fwhat", nil)

	routeToLpa(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/somewhere/what", r.URL.Path)
		assert.Equal(t, "/somewhere%2Fwhat", r.URL.RawPath)

		w.WriteHeader(http.StatusTeapot)
	}), nil).ServeHTTP(w, r)

	res := w.Result()

	assert.Equal(t, http.StatusTeapot, res.StatusCode)
}

func TestRouteToLpaWithoutID(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "/lpa/", nil)

	notFoundHandler := newMockHandler(t)
	notFoundHandler.
		On("Execute", page.AppData{}, w, r).
		Return(nil)

	routeToLpa(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTeapot)
	}), notFoundHandler.Execute).ServeHTTP(w, r)

	res := w.Result()
	assert.Equal(t, http.StatusOK, res.StatusCode)
}