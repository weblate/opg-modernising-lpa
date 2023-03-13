package donor

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ministryofjustice/opg-modernising-lpa/internal/page"
	"github.com/stretchr/testify/assert"
)

func TestGetDashboard(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "/", nil)

	lpas := []*page.Lpa{{ID: "123"}, {ID: "456"}}

	lpaStore := newMockLpaStore(t)
	lpaStore.
		On("GetAll", r.Context()).
		Return(lpas, nil)

	template := newMockTemplate(t)
	template.
		On("Execute", w, &dashboardData{App: testAppData, Lpas: lpas}).
		Return(nil)

	err := Dashboard(template.Execute, lpaStore)(testAppData, w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestGetDashboardWhenDataStoreErrors(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "/", nil)

	lpas := []*page.Lpa{{}}

	lpaStore := newMockLpaStore(t)
	lpaStore.
		On("GetAll", r.Context()).
		Return(lpas, expectedError)

	err := Dashboard(nil, lpaStore)(testAppData, w, r)

	assert.Equal(t, expectedError, err)
}

func TestGetDashboardWhenTemplateErrors(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "/", nil)

	lpas := []*page.Lpa{{}}

	lpaStore := newMockLpaStore(t)
	lpaStore.
		On("GetAll", r.Context()).
		Return(lpas, nil)

	template := newMockTemplate(t)
	template.
		On("Execute", w, &dashboardData{App: testAppData, Lpas: lpas}).
		Return(expectedError)

	err := Dashboard(template.Execute, lpaStore)(testAppData, w, r)

	assert.Equal(t, expectedError, err)
}

func TestPostDashboard(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodPost, "/", nil)

	lpaStore := newMockLpaStore(t)
	lpaStore.
		On("Create", r.Context()).
		Return(&page.Lpa{ID: "123"}, nil)

	err := Dashboard(nil, lpaStore)(testAppData, w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusFound, resp.StatusCode)
	assert.Equal(t, "/lpa/lpa-id"+page.Paths.YourDetails, resp.Header.Get("Location"))
}