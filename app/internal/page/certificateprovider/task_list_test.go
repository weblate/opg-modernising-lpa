package certificateprovider

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ministryofjustice/opg-modernising-lpa/app/internal/actor"
	"github.com/ministryofjustice/opg-modernising-lpa/app/internal/page"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetTaskList(t *testing.T) {
	testCases := map[string]struct {
		lpa                 *page.Lpa
		certificateProvider *actor.CertificateProviderProvidedDetails
		appData             page.AppData
		expected            func([]taskListItem) []taskListItem
	}{
		"empty": {
			lpa:                 &page.Lpa{ID: "lpa-id"},
			certificateProvider: &actor.CertificateProviderProvidedDetails{},
			appData:             testAppData,
			expected: func(items []taskListItem) []taskListItem {
				items[1].Disabled = true
				items[2].Disabled = true
				items[3].Disabled = true

				return items
			},
		},
		"paid": {
			lpa: &page.Lpa{
				ID: "lpa-id",
				Tasks: page.Tasks{
					PayForLpa: actor.PaymentTaskCompleted,
				},
			},
			certificateProvider: &actor.CertificateProviderProvidedDetails{
				Tasks: actor.CertificateProviderTasks{
					ConfirmYourDetails: actor.TaskCompleted,
				},
			},
			appData: testAppData,
			expected: func(items []taskListItem) []taskListItem {
				items[0].State = actor.TaskCompleted
				items[1].Disabled = true
				items[2].Disabled = true
				items[3].Disabled = true

				return items
			},
		},
		"submitted": {
			lpa: &page.Lpa{
				ID:        "lpa-id",
				Submitted: time.Now(),
			},
			certificateProvider: &actor.CertificateProviderProvidedDetails{
				Tasks: actor.CertificateProviderTasks{
					ConfirmYourDetails: actor.TaskCompleted,
				},
			},
			appData: testAppData,
			expected: func(items []taskListItem) []taskListItem {
				items[0].State = actor.TaskCompleted
				items[1].Disabled = true
				items[3].Disabled = true

				return items
			},
		},
		"all": {
			lpa: &page.Lpa{
				ID:        "lpa-id",
				Submitted: time.Now(),
				Tasks: page.Tasks{
					PayForLpa: actor.PaymentTaskCompleted,
				},
			},
			certificateProvider: &actor.CertificateProviderProvidedDetails{
				Tasks: actor.CertificateProviderTasks{
					ConfirmYourDetails:    actor.TaskCompleted,
					ConfirmYourIdentity:   actor.TaskCompleted,
					ReadTheLpa:            actor.TaskCompleted,
					ProvideTheCertificate: actor.TaskCompleted,
				},
			},
			appData: testAppData,
			expected: func(items []taskListItem) []taskListItem {
				items[0].State = actor.TaskCompleted
				items[1].State = actor.TaskCompleted
				items[2].State = actor.TaskCompleted
				items[3].State = actor.TaskCompleted

				return items
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r, _ := http.NewRequest(http.MethodGet, "/", nil)

			donorStore := newMockDonorStore(t)
			donorStore.
				On("GetAny", r.Context()).
				Return(tc.lpa, nil)

			certificateProviderStore := newMockCertificateProviderStore(t)
			certificateProviderStore.
				On("Get", r.Context()).
				Return(tc.certificateProvider, nil)

			template := newMockTemplate(t)
			template.
				On("Execute", w, &taskListData{
					App: tc.appData,
					Lpa: tc.lpa,
					Items: tc.expected([]taskListItem{
						{Name: "confirmYourDetails", Path: page.Paths.CertificateProvider.EnterDateOfBirth.Format("lpa-id")},
						{Name: "confirmYourIdentity", Path: page.Paths.CertificateProvider.WhatYoullNeedToConfirmYourIdentity.Format("lpa-id")},
						{Name: "readTheLpa", Path: page.Paths.CertificateProvider.ReadTheLpa.Format("lpa-id")},
						{Name: "provideTheCertificateForThisLpa", Path: page.Paths.CertificateProvider.ProvideCertificate.Format("lpa-id")},
					}),
				}).
				Return(nil)

			err := TaskList(template.Execute, donorStore, certificateProviderStore)(tc.appData, w, r)
			resp := w.Result()

			assert.Nil(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
		})
	}
}

func TestGetTaskListWhenDonorStoreErrors(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "/", nil)

	donorStore := newMockDonorStore(t)
	donorStore.
		On("GetAny", r.Context()).
		Return(&page.Lpa{}, expectedError)

	err := TaskList(nil, donorStore, nil)(testAppData, w, r)

	assert.Equal(t, expectedError, err)
}

func TestGetTaskListWhenCertificateProviderStoreErrors(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "/", nil)

	donorStore := newMockDonorStore(t)
	donorStore.
		On("GetAny", r.Context()).
		Return(&page.Lpa{ID: "lpa-id"}, nil)

	certificateProviderStore := newMockCertificateProviderStore(t)
	certificateProviderStore.
		On("Get", mock.Anything).
		Return(nil, expectedError)

	err := TaskList(nil, donorStore, certificateProviderStore)(testAppData, w, r)

	assert.Equal(t, expectedError, err)
}

func TestGetTaskListWhenTemplateErrors(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "/", nil)

	donorStore := newMockDonorStore(t)
	donorStore.
		On("GetAny", r.Context()).
		Return(&page.Lpa{ID: "lpa-id"}, nil)

	certificateProviderStore := newMockCertificateProviderStore(t)
	certificateProviderStore.
		On("Get", r.Context()).
		Return(&actor.CertificateProviderProvidedDetails{}, nil)

	template := newMockTemplate(t)
	template.
		On("Execute", w, mock.Anything).
		Return(expectedError)

	err := TaskList(template.Execute, donorStore, certificateProviderStore)(testAppData, w, r)
	resp := w.Result()

	assert.Equal(t, expectedError, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
