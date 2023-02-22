package donor

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ministryofjustice/opg-modernising-lpa/internal/identity"
	"github.com/ministryofjustice/opg-modernising-lpa/internal/page"
	"github.com/stretchr/testify/assert"
)

func TestGetIdentityWithYotiCallback(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "/?token=a-token", nil)
	now := time.Now()
	userData := identity.UserData{FullName: "a-full-name", RetrievedAt: now}

	lpaStore := newMockLpaStore(t)
	lpaStore.On("Get", r.Context()).Return(&page.Lpa{}, nil)
	lpaStore.On("Put", r.Context(), &page.Lpa{YotiUserData: userData}).Return(nil)

	yotiClient := newMockYotiClient(t)
	yotiClient.On("User", "a-token").Return(userData, nil)

	template := newMockTemplate(t)
	template.
		On("Execute", w, &identityWithYotiCallbackData{
			App:         testAppData,
			FullName:    "a-full-name",
			ConfirmedAt: now,
		}).
		Return(nil)

	err := IdentityWithYotiCallback(template.Execute, yotiClient, lpaStore)(testAppData, w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestGetIdentityWithYotiCallbackWhenError(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "/?token=a-token", nil)

	lpaStore := newMockLpaStore(t)
	lpaStore.On("Get", r.Context()).Return(&page.Lpa{}, nil)

	yotiClient := newMockYotiClient(t)
	yotiClient.On("User", "a-token").Return(identity.UserData{}, expectedError)

	err := IdentityWithYotiCallback(nil, yotiClient, lpaStore)(testAppData, w, r)

	assert.Equal(t, expectedError, err)
}

func TestGetIdentityWithYotiCallbackWhenGetDataStoreError(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "/?token=a-token", nil)

	lpaStore := newMockLpaStore(t)
	lpaStore.On("Get", r.Context()).Return(&page.Lpa{}, expectedError)

	err := IdentityWithYotiCallback(nil, nil, lpaStore)(testAppData, w, r)

	assert.Equal(t, expectedError, err)
}

func TestGetIdentityWithYotiCallbackWhenPutDataStoreError(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "/?token=a-token", nil)
	now := time.Now()
	userData := identity.UserData{FullName: "a-full-name", RetrievedAt: now}

	lpaStore := newMockLpaStore(t)
	lpaStore.On("Get", r.Context()).Return(&page.Lpa{}, nil)
	lpaStore.On("Put", r.Context(), &page.Lpa{YotiUserData: userData}).Return(expectedError)

	yotiClient := newMockYotiClient(t)
	yotiClient.On("User", "a-token").Return(userData, nil)

	err := IdentityWithYotiCallback(nil, yotiClient, lpaStore)(testAppData, w, r)

	assert.Equal(t, expectedError, err)
}

func TestGetIdentityWithYotiCallbackWhenReturning(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "/?token=a-token", nil)
	now := time.Date(2012, time.January, 1, 2, 3, 4, 5, time.UTC)
	userData := identity.UserData{OK: true, FullName: "a-full-name", RetrievedAt: now}

	lpaStore := newMockLpaStore(t)
	lpaStore.On("Get", r.Context()).Return(&page.Lpa{YotiUserData: userData}, nil)

	template := newMockTemplate(t)
	template.
		On("Execute", w, &identityWithYotiCallbackData{
			App:         testAppData,
			FullName:    "a-full-name",
			ConfirmedAt: now,
		}).
		Return(nil)

	err := IdentityWithYotiCallback(template.Execute, nil, lpaStore)(testAppData, w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestPostIdentityWithYotiCallback(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodPost, "/", nil)

	lpaStore := newMockLpaStore(t)
	lpaStore.On("Get", r.Context()).Return(&page.Lpa{
		IdentityOption: identity.EasyID,
	}, nil)

	err := IdentityWithYotiCallback(nil, nil, lpaStore)(testAppData, w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusFound, resp.StatusCode)
	assert.Equal(t, "/lpa/lpa-id"+page.Paths.ReadYourLpa, resp.Header.Get("Location"))
}
