package donor

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/ministryofjustice/opg-modernising-lpa/app/internal/actor"
	"github.com/ministryofjustice/opg-modernising-lpa/app/internal/page"
	"github.com/ministryofjustice/opg-modernising-lpa/app/internal/validation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetHowWouldCertificateProviderPreferToCarryOutTheirRole(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "/", nil)

	template := newMockTemplate(t)
	template.
		On("Execute", w, &howWouldCertificateProviderPreferToCarryOutTheirRoleData{
			App:  testAppData,
			Form: &howWouldCertificateProviderPreferToCarryOutTheirRoleForm{},
		}).
		Return(nil)

	err := HowWouldCertificateProviderPreferToCarryOutTheirRole(template.Execute, nil)(testAppData, w, r, &page.Lpa{})
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestGetHowWouldCertificateProviderPreferToCarryOutTheirRoleFromStore(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "/", nil)

	template := newMockTemplate(t)
	template.
		On("Execute", w, &howWouldCertificateProviderPreferToCarryOutTheirRoleData{
			App:                 testAppData,
			CertificateProvider: actor.CertificateProvider{CarryOutBy: "paper"},
			Form:                &howWouldCertificateProviderPreferToCarryOutTheirRoleForm{CarryOutBy: "paper"},
		}).
		Return(nil)

	err := HowWouldCertificateProviderPreferToCarryOutTheirRole(template.Execute, nil)(testAppData, w, r, &page.Lpa{
		CertificateProvider: actor.CertificateProvider{CarryOutBy: "paper"},
	})
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestGetHowWouldCertificateProviderPreferToCarryOutTheirRoleWhenTemplateErrors(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "/", nil)

	template := newMockTemplate(t)
	template.
		On("Execute", w, &howWouldCertificateProviderPreferToCarryOutTheirRoleData{
			App:  testAppData,
			Form: &howWouldCertificateProviderPreferToCarryOutTheirRoleForm{},
		}).
		Return(expectedError)

	err := HowWouldCertificateProviderPreferToCarryOutTheirRole(template.Execute, nil)(testAppData, w, r, &page.Lpa{})
	resp := w.Result()

	assert.Equal(t, expectedError, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestPostHowWouldCertificateProviderPreferToCarryOutTheirRole(t *testing.T) {
	testCases := []struct {
		carryOutBy string
		email      string
	}{
		{
			carryOutBy: "paper",
		},
		{
			carryOutBy: "email",
			email:      "someone@example.com",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.carryOutBy, func(t *testing.T) {
			form := url.Values{
				"carry-out-by": {tc.carryOutBy},
				"email":        {tc.email},
			}

			w := httptest.NewRecorder()
			r, _ := http.NewRequest(http.MethodPost, "/", strings.NewReader(form.Encode()))
			r.Header.Add("Content-Type", page.FormUrlEncoded)

			donorStore := newMockDonorStore(t)
			donorStore.
				On("Put", r.Context(), &page.Lpa{
					ID:                  "lpa-id",
					CertificateProvider: actor.CertificateProvider{CarryOutBy: tc.carryOutBy, Email: tc.email},
				}).
				Return(nil)

			err := HowWouldCertificateProviderPreferToCarryOutTheirRole(nil, donorStore)(testAppData, w, r, &page.Lpa{ID: "lpa-id"})
			resp := w.Result()

			assert.Nil(t, err)
			assert.Equal(t, http.StatusFound, resp.StatusCode)
			assert.Equal(t, page.Paths.CertificateProviderAddress.Format("lpa-id"), resp.Header.Get("Location"))
		})
	}
}

func TestPostHowWouldCertificateProviderPreferToCarryOutTheirRoleWhenStoreErrors(t *testing.T) {
	form := url.Values{
		"carry-out-by": {"paper"},
	}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodPost, "/", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", page.FormUrlEncoded)

	donorStore := newMockDonorStore(t)
	donorStore.
		On("Put", r.Context(), mock.Anything).
		Return(expectedError)

	err := HowWouldCertificateProviderPreferToCarryOutTheirRole(nil, donorStore)(testAppData, w, r, &page.Lpa{})

	assert.Equal(t, expectedError, err)
}

func TestPostHowWouldCertificateProviderPreferToCarryOutTheirRoleWhenValidationErrors(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodPost, "/", strings.NewReader("nope"))
	r.Header.Add("Content-Type", page.FormUrlEncoded)

	template := newMockTemplate(t)
	template.
		On("Execute", w, &howWouldCertificateProviderPreferToCarryOutTheirRoleData{
			App:    testAppData,
			Errors: validation.With("carry-out-by", validation.SelectError{Label: "howYourCertificateProviderWouldPreferToCarryOutTheirRole"}),
			Form:   &howWouldCertificateProviderPreferToCarryOutTheirRoleForm{},
		}).
		Return(nil)

	err := HowWouldCertificateProviderPreferToCarryOutTheirRole(template.Execute, nil)(testAppData, w, r, &page.Lpa{})
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestReadHowWouldCertificateProviderPreferToCarryOutTheirRoleForm(t *testing.T) {
	form := url.Values{
		"carry-out-by": {"email"},
		"email":        {"someone@example.com"},
	}

	r, _ := http.NewRequest(http.MethodPost, "/", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", page.FormUrlEncoded)

	result := readHowWouldCertificateProviderPreferToCarryOutTheirRole(r)

	assert.Equal(t, "email", result.CarryOutBy)
	assert.Equal(t, "someone@example.com", result.Email)
}

func TestHowWouldCertificateProviderPreferToCarryOutTheirRoleFormValidate(t *testing.T) {
	testCases := map[string]struct {
		form   *howWouldCertificateProviderPreferToCarryOutTheirRoleForm
		errors validation.List
	}{
		"paper": {
			form: &howWouldCertificateProviderPreferToCarryOutTheirRoleForm{
				CarryOutBy: "paper",
			},
		},
		"email": {
			form: &howWouldCertificateProviderPreferToCarryOutTheirRoleForm{
				CarryOutBy: "email",
				Email:      "someone@example.com",
			},
		},
		"email invalid": {
			form: &howWouldCertificateProviderPreferToCarryOutTheirRoleForm{
				CarryOutBy: "email",
				Email:      "what",
			},
			errors: validation.With("email", validation.EmailError{Label: "certificateProvidersEmail"}),
		},
		"email missing": {
			form: &howWouldCertificateProviderPreferToCarryOutTheirRoleForm{
				CarryOutBy: "email",
			},
			errors: validation.With("email", validation.EnterError{Label: "certificateProvidersEmail"}),
		},
		"missing": {
			form:   &howWouldCertificateProviderPreferToCarryOutTheirRoleForm{},
			errors: validation.With("carry-out-by", validation.SelectError{Label: "howYourCertificateProviderWouldPreferToCarryOutTheirRole"}),
		},
		"invalid": {
			form: &howWouldCertificateProviderPreferToCarryOutTheirRoleForm{
				CarryOutBy: "what",
			},
			errors: validation.With("carry-out-by", validation.SelectError{Label: "howYourCertificateProviderWouldPreferToCarryOutTheirRole"}),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.errors, tc.form.Validate())
		})
	}
}
