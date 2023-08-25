package donor

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ministryofjustice/opg-modernising-lpa/internal/actor"
	"github.com/ministryofjustice/opg-modernising-lpa/internal/page"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetHowToSendEvidence(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "/about-payment", nil)

	template := newMockTemplate(t)
	template.
		On("Execute", w, &howToSendEvidenceData{App: testAppData}).
		Return(nil)

	err := HowToSendEvidence(template.Execute, nil)(testAppData, w, r, &page.Lpa{})
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestGetHowToSendEvidenceWhenTemplateErrors(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "/about-payment", nil)

	template := newMockTemplate(t)
	template.
		On("Execute", w, &howToSendEvidenceData{App: testAppData}).
		Return(expectedError)

	err := HowToSendEvidence(template.Execute, nil)(testAppData, w, r, &page.Lpa{})
	resp := w.Result()

	assert.Equal(t, expectedError, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestPostHowToSendEvidence(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodPost, "/about-payment", nil)

	lpa := &page.Lpa{ID: "lpa-id", Donor: actor.Donor{Email: "a@b.com"}}

	payer := newMockPayer(t)
	payer.
		On("Pay", testAppData, w, r, lpa).
		Return(nil)

	err := HowToSendEvidence(nil, payer)(testAppData, w, r, lpa)
	assert.Nil(t, err)
}

func TestPostHowToSendEvidenceWhenPayerErrors(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodPost, "/about-payment", nil)

	payer := newMockPayer(t)
	payer.
		On("Pay", testAppData, w, r, mock.Anything).
		Return(expectedError)

	err := HowToSendEvidence(nil, payer)(testAppData, w, r, &page.Lpa{})
	assert.Equal(t, expectedError, err)
}