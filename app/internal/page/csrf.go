package page

import (
	"bytes"
	"errors"
	"io"
	"mime"
	"mime/multipart"
	"net/http"

	"github.com/gorilla/sessions"
)

type contextKey string

var ErrCsrfInvalid = errors.New("CSRF token not valid")

const csrfTokenLength = 12

func ValidateCsrf(next http.Handler, store sessions.Store, randomString func(int) string, errorHandler ErrorHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		csrfSession, err := store.Get(r, "csrf")

		if r.Method == http.MethodPost {
			if err != nil {
				errorHandler(w, r, err)
				return
			}

			if !csrfValid(r, csrfSession) {
				errorHandler(w, r, ErrCsrfInvalid)
				return
			}
		}

		if csrfSession.IsNew {
			csrfSession.Values = map[any]any{"token": randomString(csrfTokenLength)}
			csrfSession.Options = &sessions.Options{
				MaxAge:   24 * 60 * 60,
				Secure:   true,
				HttpOnly: true,
				SameSite: http.SameSiteLaxMode,
			}
			_ = store.Save(r, w, csrfSession)
		}

		appData := AppDataFromContext(ctx)
		appData.CsrfToken, _ = csrfSession.Values["token"].(string)

		next.ServeHTTP(w, r.WithContext(ContextWithAppData(ctx, appData)))
	}
}

func csrfValid(r *http.Request, csrfSession *sessions.Session) bool {
	cookieValue, ok := csrfSession.Values["token"].(string)
	if !ok {
		return false
	}

	if mediaType, params, err := mime.ParseMediaType(r.Header.Get("Content-Type")); err == nil && mediaType == "multipart/form-data" {
		var buf bytes.Buffer
		reader := multipart.NewReader(io.TeeReader(r.Body, &buf), params["boundary"])

		part, err := reader.NextPart()
		if err != nil {
			return false
		}

		if part.FormName() != "csrf" {
			return false
		}

		lmt := io.LimitReader(part, csrfTokenLength+1)
		value, _ := io.ReadAll(lmt)

		r.Body = MultiReadCloser(io.NopCloser(&buf), r.Body)
		return string(value) == cookieValue
	}

	return r.PostFormValue("csrf") == cookieValue
}
