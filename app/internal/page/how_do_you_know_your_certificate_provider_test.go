package page

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetHowDoYouKnowYourCertificateProvider(t *testing.T) {
	w := httptest.NewRecorder()

	dataStore := &mockDataStore{}
	dataStore.
		On("Get", mock.Anything, "session-id").
		Return(nil)

	template := &mockTemplate{}
	template.
		On("Func", w, &howDoYouKnowYourCertificateProviderData{
			App:  appData,
			Form: &howDoYouKnowYourCertificateProviderForm{},
		}).
		Return(nil)

	r, _ := http.NewRequest(http.MethodGet, "/", nil)

	err := HowDoYouKnowYourCertificateProvider(template.Func, dataStore)(appData, w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, template, dataStore)
}

func TestGetHowDoYouKnowYourCertificateProviderWhenStoreErrors(t *testing.T) {
	w := httptest.NewRecorder()

	dataStore := &mockDataStore{}
	dataStore.
		On("Get", mock.Anything, "session-id").
		Return(expectedError)

	r, _ := http.NewRequest(http.MethodGet, "/", nil)

	err := HowDoYouKnowYourCertificateProvider(nil, dataStore)(appData, w, r)

	assert.Equal(t, expectedError, err)
	mock.AssertExpectationsForObjects(t, dataStore)
}

func TestGetHowDoYouKnowYourCertificateProviderFromStore(t *testing.T) {
	w := httptest.NewRecorder()

	certificateProvider := CertificateProvider{
		Relationship: []string{"friend"},
	}
	dataStore := &mockDataStore{
		data: Lpa{
			CertificateProvider: certificateProvider,
		},
	}
	dataStore.
		On("Get", mock.Anything, "session-id").
		Return(nil)

	template := &mockTemplate{}
	template.
		On("Func", w, &howDoYouKnowYourCertificateProviderData{
			App:                 appData,
			CertificateProvider: certificateProvider,
			Form:                &howDoYouKnowYourCertificateProviderForm{How: []string{"friend"}},
		}).
		Return(nil)

	r, _ := http.NewRequest(http.MethodGet, "/", nil)

	err := HowDoYouKnowYourCertificateProvider(template.Func, dataStore)(appData, w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, template)
}

func TestGetHowDoYouKnowYourCertificateProviderWhenTemplateErrors(t *testing.T) {
	w := httptest.NewRecorder()

	dataStore := &mockDataStore{}
	dataStore.
		On("Get", mock.Anything, "session-id").
		Return(nil)

	template := &mockTemplate{}
	template.
		On("Func", w, mock.Anything).
		Return(expectedError)

	r, _ := http.NewRequest(http.MethodGet, "/", nil)

	err := HowDoYouKnowYourCertificateProvider(template.Func, dataStore)(appData, w, r)
	resp := w.Result()

	assert.Equal(t, expectedError, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, template)
}

func TestPostHowDoYouKnowYourCertificateProvider(t *testing.T) {
	testCases := map[string]struct {
		form                url.Values
		certificateProvider CertificateProvider
		taskState           TaskState
		redirect            string
	}{
		"professional": {
			form: url.Values{"how": {"legal-professional", "health-professional"}},
			certificateProvider: CertificateProvider{
				FirstNames:   "John",
				Relationship: []string{"legal-professional", "health-professional"},
			},
			taskState: TaskCompleted,
			redirect:  taskListPath,
		},
		"other": {
			form: url.Values{"how": {"other"}, "description": {"This"}},
			certificateProvider: CertificateProvider{
				FirstNames:              "John",
				Relationship:            []string{"other"},
				RelationshipDescription: "This",
				RelationshipLength:      "gte-2-years",
			},
			taskState: TaskInProgress,
			redirect:  howLongHaveYouKnownCertificateProviderPath,
		},
		"lay": {
			form: url.Values{"how": {"friend", "neighbour", "colleague"}},
			certificateProvider: CertificateProvider{
				FirstNames:         "John",
				Relationship:       []string{"friend", "neighbour", "colleague"},
				RelationshipLength: "gte-2-years",
			},
			taskState: TaskInProgress,
			redirect:  howLongHaveYouKnownCertificateProviderPath,
		},
		"mixed": {
			form: url.Values{"how": {"legal-professional", "friend", "other"}, "description": {"This"}},
			certificateProvider: CertificateProvider{
				FirstNames:              "John",
				Relationship:            []string{"legal-professional", "friend", "other"},
				RelationshipDescription: "This",
				RelationshipLength:      "gte-2-years",
			},
			taskState: TaskInProgress,
			redirect:  howLongHaveYouKnownCertificateProviderPath,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			w := httptest.NewRecorder()

			dataStore := &mockDataStore{
				data: Lpa{
					CertificateProvider: CertificateProvider{FirstNames: "John", Relationship: []string{"what"}, RelationshipLength: "gte-2-years"},
				},
			}
			dataStore.
				On("Get", mock.Anything, "session-id").
				Return(nil)
			dataStore.
				On("Put", mock.Anything, "session-id", Lpa{
					CertificateProvider: tc.certificateProvider,
					Tasks: Tasks{
						CertificateProvider: tc.taskState,
					},
				}).
				Return(nil)

			r, _ := http.NewRequest(http.MethodPost, "/", strings.NewReader(tc.form.Encode()))
			r.Header.Add("Content-Type", formUrlEncoded)

			err := HowDoYouKnowYourCertificateProvider(nil, dataStore)(appData, w, r)
			resp := w.Result()

			assert.Nil(t, err)
			assert.Equal(t, http.StatusFound, resp.StatusCode)
			assert.Equal(t, tc.redirect, resp.Header.Get("Location"))
			mock.AssertExpectationsForObjects(t, dataStore)
		})
	}
}

func TestPostHowDoYouKnowYourCertificateProviderWhenStoreErrors(t *testing.T) {
	w := httptest.NewRecorder()

	dataStore := &mockDataStore{}
	dataStore.
		On("Get", mock.Anything, "session-id").
		Return(nil)
	dataStore.
		On("Put", mock.Anything, "session-id", mock.Anything).
		Return(expectedError)

	form := url.Values{
		"how": {"friend"},
	}

	r, _ := http.NewRequest(http.MethodPost, "/", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)

	err := HowDoYouKnowYourCertificateProvider(nil, dataStore)(appData, w, r)

	assert.Equal(t, expectedError, err)
	mock.AssertExpectationsForObjects(t, dataStore)
}

func TestPostHowDoYouKnowYourCertificateProviderWhenValidationErrors(t *testing.T) {
	w := httptest.NewRecorder()

	dataStore := &mockDataStore{}
	dataStore.
		On("Get", mock.Anything, "session-id").
		Return(nil)

	template := &mockTemplate{}
	template.
		On("Func", w, &howDoYouKnowYourCertificateProviderData{
			App:  appData,
			Form: &howDoYouKnowYourCertificateProviderForm{},
			Errors: map[string]string{
				"how": "selectHowYouKnowCertificateProvider",
			},
		}).
		Return(nil)

	r, _ := http.NewRequest(http.MethodPost, "/", strings.NewReader(""))
	r.Header.Add("Content-Type", formUrlEncoded)

	err := HowDoYouKnowYourCertificateProvider(template.Func, dataStore)(appData, w, r)
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mock.AssertExpectationsForObjects(t, template)
}

func TestReadHowDoYouKnowYourCertificateProviderForm(t *testing.T) {
	form := url.Values{
		"how":         {"friend", "other"},
		"description": {"What"},
	}

	r, _ := http.NewRequest(http.MethodPost, "/", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", formUrlEncoded)

	result := readHowDoYouKnowYourCertificateProviderForm(r)

	assert.Equal(t, []string{"friend", "other"}, result.How)
	assert.Equal(t, "What", result.Description)
}

func TestHowDoYouKnowYourCertificateProviderFormValidate(t *testing.T) {
	testCases := map[string]struct {
		form   *howDoYouKnowYourCertificateProviderForm
		errors map[string]string
	}{
		"all": {
			form: &howDoYouKnowYourCertificateProviderForm{
				How:         []string{"friend", "neighbour", "colleague", "health-professional", "legal-professional", "other"},
				Description: "This",
			},
			errors: map[string]string{},
		},
		"missing": {
			form: &howDoYouKnowYourCertificateProviderForm{},
			errors: map[string]string{
				"how": "selectHowYouKnowCertificateProvider",
			},
		},
		"invalid-option": {
			form: &howDoYouKnowYourCertificateProviderForm{
				How: []string{"what"},
			},
			errors: map[string]string{
				"how": "selectHowYouKnowCertificateProvider",
			},
		},
		"other-missing-description": {
			form: &howDoYouKnowYourCertificateProviderForm{
				How: []string{"other"},
			},
			errors: map[string]string{
				"description": "enterDescription",
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.errors, tc.form.Validate())
		})
	}
}