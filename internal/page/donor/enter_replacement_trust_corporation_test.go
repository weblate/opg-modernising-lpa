package donor

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/ministryofjustice/opg-modernising-lpa/internal/actor"
	"github.com/ministryofjustice/opg-modernising-lpa/internal/page"
	"github.com/ministryofjustice/opg-modernising-lpa/internal/validation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetEnterReplacementTrustCorporation(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "/", nil)

	template := newMockTemplate(t)
	template.
		On("Execute", w, &enterReplacementTrustCorporationData{
			App:  testAppData,
			Form: &enterTrustCorporationForm{},
		}).
		Return(nil)

	err := EnterReplacementTrustCorporation(template.Execute, nil)(testAppData, w, r, &page.Lpa{})
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestGetEnterReplacementTrustCorporationWhenTemplateErrors(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "/", nil)

	template := newMockTemplate(t)
	template.
		On("Execute", w, &enterReplacementTrustCorporationData{
			App:  testAppData,
			Form: &enterTrustCorporationForm{},
		}).
		Return(expectedError)

	err := EnterReplacementTrustCorporation(template.Execute, nil)(testAppData, w, r, &page.Lpa{})
	resp := w.Result()

	assert.Equal(t, expectedError, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestPostEnterReplacementTrustCorporation(t *testing.T) {
	form := url.Values{
		"name":           {"Co co."},
		"company-number": {"453345"},
		"email":          {"name@example.com"},
	}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodPost, "/", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", page.FormUrlEncoded)

	donorStore := newMockDonorStore(t)
	donorStore.
		On("Put", r.Context(), &page.Lpa{
			ID: "lpa-id",
			ReplacementAttorneys: actor.Attorneys{
				TrustCorporation: actor.TrustCorporation{
					Name:          "Co co.",
					CompanyNumber: "453345",
					Email:         "name@example.com",
				},
			},
			Tasks: page.Tasks{
				ChooseReplacementAttorneys: actor.TaskInProgress,
			},
		}).
		Return(nil)

	err := EnterReplacementTrustCorporation(nil, donorStore)(testAppData, w, r, &page.Lpa{
		ID: "lpa-id",
	})
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusFound, resp.StatusCode)
	assert.Equal(t, page.Paths.EnterReplacementTrustCorporationAddress.Format("lpa-id"), resp.Header.Get("Location"))
}

func TestPostEnterReplacementTrustCorporationWhenValidationError(t *testing.T) {
	form := url.Values{
		"company-number": {"453345"},
		"email":          {"name@example.com"},
	}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodPost, "/", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", page.FormUrlEncoded)

	template := newMockTemplate(t)
	template.
		On("Execute", w, mock.MatchedBy(func(data *enterReplacementTrustCorporationData) bool {
			return assert.Equal(t, validation.With("name", validation.EnterError{Label: "companyName"}), data.Errors)
		})).
		Return(nil)

	err := EnterReplacementTrustCorporation(template.Execute, nil)(testAppData, w, r, &page.Lpa{
		Donor: actor.Donor{FirstNames: "Jane", LastName: "Doe"},
	})
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestPostEnterReplacementTrustCorporationWhenStoreErrors(t *testing.T) {
	form := url.Values{
		"name":           {"Inc co."},
		"company-number": {"64365634"},
		"email":          {"john@example.com"},
	}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodPost, "/", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", page.FormUrlEncoded)

	donorStore := newMockDonorStore(t)
	donorStore.
		On("Put", r.Context(), mock.Anything).
		Return(expectedError)

	err := EnterReplacementTrustCorporation(nil, donorStore)(testAppData, w, r, &page.Lpa{})

	assert.Equal(t, expectedError, err)
}
