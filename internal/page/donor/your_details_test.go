package donor

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/sessions"
	"github.com/ministryofjustice/opg-modernising-lpa/internal/actor"
	"github.com/ministryofjustice/opg-modernising-lpa/internal/date"
	"github.com/ministryofjustice/opg-modernising-lpa/internal/page"
	"github.com/ministryofjustice/opg-modernising-lpa/internal/place"
	"github.com/ministryofjustice/opg-modernising-lpa/internal/sesh"
	"github.com/ministryofjustice/opg-modernising-lpa/internal/validation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetYourDetails(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "/", nil)

	template := newMockTemplate(t)
	template.
		On("Execute", w, &yourDetailsData{
			App:  testAppData,
			Form: &yourDetailsForm{},
		}).
		Return(nil)

	err := YourDetails(template.Execute, nil, nil)(testAppData, w, r, &page.Lpa{})
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestGetYourDetailsFromStore(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "/", nil)

	template := newMockTemplate(t)
	template.
		On("Execute", w, &yourDetailsData{
			App: testAppData,
			Form: &yourDetailsForm{
				FirstNames: "John",
			},
		}).
		Return(nil)

	err := YourDetails(template.Execute, nil, nil)(testAppData, w, r, &page.Lpa{
		Donor: actor.Donor{
			FirstNames: "John",
		},
	})
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestGetYourDetailsWhenTemplateErrors(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "/", nil)

	template := newMockTemplate(t)
	template.
		On("Execute", w, &yourDetailsData{
			App:  testAppData,
			Form: &yourDetailsForm{},
		}).
		Return(expectedError)

	err := YourDetails(template.Execute, nil, nil)(testAppData, w, r, &page.Lpa{})
	resp := w.Result()

	assert.Equal(t, expectedError, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestPostYourDetails(t *testing.T) {
	validBirthYear := strconv.Itoa(time.Now().Year() - 40)

	testCases := map[string]struct {
		form   url.Values
		person actor.Donor
	}{
		"valid": {
			form: url.Values{
				"first-names":         {"John"},
				"last-name":           {"Doe"},
				"date-of-birth-day":   {"2"},
				"date-of-birth-month": {"1"},
				"date-of-birth-year":  {validBirthYear},
			},
			person: actor.Donor{
				FirstNames:  "John",
				LastName:    "Doe",
				DateOfBirth: date.New(validBirthYear, "1", "2"),
				Address:     place.Address{Line1: "abc"},
				Email:       "name@example.com",
			},
		},
		"warning ignored": {
			form: url.Values{
				"first-names":         {"John"},
				"last-name":           {"Doe"},
				"date-of-birth-day":   {"2"},
				"date-of-birth-month": {"1"},
				"date-of-birth-year":  {"1900"},
				"ignore-dob-warning":  {"dateOfBirthIsOver100"},
			},
			person: actor.Donor{
				FirstNames:  "John",
				LastName:    "Doe",
				DateOfBirth: date.New("1900", "1", "2"),
				Address:     place.Address{Line1: "abc"},
				Email:       "name@example.com",
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			w := httptest.NewRecorder()

			r, _ := http.NewRequest(http.MethodPost, "/", strings.NewReader(tc.form.Encode()))
			r.Header.Add("Content-Type", page.FormUrlEncoded)

			donorStore := newMockDonorStore(t)
			donorStore.
				On("Put", r.Context(), &page.Lpa{
					ID:    "lpa-id",
					Donor: tc.person,
					Tasks: page.Tasks{YourDetails: actor.TaskInProgress},
				}).
				Return(nil)

			sessionStore := newMockSessionStore(t)
			sessionStore.
				On("Get", r, "session").
				Return(&sessions.Session{Values: map[any]any{"session": &sesh.LoginSession{Sub: "xyz", Email: "name@example.com"}}}, nil)

			err := YourDetails(nil, donorStore, sessionStore)(testAppData, w, r, &page.Lpa{
				ID: "lpa-id",
				Donor: actor.Donor{
					FirstNames: "John",
					Address:    place.Address{Line1: "abc"},
				},
			})
			resp := w.Result()

			assert.Nil(t, err)
			assert.Equal(t, http.StatusFound, resp.StatusCode)
			assert.Equal(t, page.Paths.YourAddress.Format("lpa-id"), resp.Header.Get("Location"))
		})
	}
}

func TestPostYourDetailsWhenTaskCompleted(t *testing.T) {
	validBirthYear := strconv.Itoa(time.Now().Year() - 40)

	form := url.Values{
		"first-names":         {"John"},
		"last-name":           {"Doe"},
		"date-of-birth-day":   {"2"},
		"date-of-birth-month": {"1"},
		"date-of-birth-year":  {validBirthYear},
	}

	w := httptest.NewRecorder()

	r, _ := http.NewRequest(http.MethodPost, "/", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", page.FormUrlEncoded)

	donorStore := newMockDonorStore(t)
	donorStore.
		On("Put", r.Context(), &page.Lpa{
			ID: "lpa-id",
			Donor: actor.Donor{
				FirstNames:  "John",
				LastName:    "Doe",
				DateOfBirth: date.New(validBirthYear, "1", "2"),
				Address:     place.Address{Line1: "abc"},
				Email:       "name@example.com",
			},
			Tasks: page.Tasks{YourDetails: actor.TaskCompleted},
		}).
		Return(nil)

	sessionStore := newMockSessionStore(t)
	sessionStore.
		On("Get", r, "session").
		Return(&sessions.Session{Values: map[any]any{"session": &sesh.LoginSession{Sub: "xyz", Email: "name@example.com"}}}, nil)

	err := YourDetails(nil, donorStore, sessionStore)(testAppData, w, r, &page.Lpa{
		ID: "lpa-id",
		Donor: actor.Donor{
			FirstNames: "John",
			Address:    place.Address{Line1: "abc"},
		},
		Tasks: page.Tasks{YourDetails: actor.TaskCompleted},
	})
	resp := w.Result()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusFound, resp.StatusCode)
	assert.Equal(t, page.Paths.YourAddress.Format("lpa-id"), resp.Header.Get("Location"))
}

func TestPostYourDetailsWhenInputRequired(t *testing.T) {
	testCases := map[string]struct {
		form        url.Values
		dataMatcher func(t *testing.T, data *yourDetailsData) bool
	}{
		"validation error": {
			form: url.Values{
				"last-name":           {"Doe"},
				"date-of-birth-day":   {"2"},
				"date-of-birth-month": {"1"},
				"date-of-birth-year":  {"1990"},
			},
			dataMatcher: func(t *testing.T, data *yourDetailsData) bool {
				return assert.Equal(t, validation.With("first-names", validation.EnterError{Label: "firstNames"}), data.Errors)
			},
		},
		"dob warning": {
			form: url.Values{
				"first-names":         {"John"},
				"last-name":           {"Doe"},
				"date-of-birth-day":   {"2"},
				"date-of-birth-month": {"1"},
				"date-of-birth-year":  {"1900"},
			},
			dataMatcher: func(t *testing.T, data *yourDetailsData) bool {
				return assert.Equal(t, "dateOfBirthIsOver100", data.DobWarning)
			},
		},
		"dob warning ignored but other errors": {
			form: url.Values{
				"first-names":         {"John"},
				"date-of-birth-day":   {"2"},
				"date-of-birth-month": {"1"},
				"date-of-birth-year":  {"1900"},
				"ignore-dob-warning":  {"dateOfBirthIsOver100"},
			},
			dataMatcher: func(t *testing.T, data *yourDetailsData) bool {
				return assert.Equal(t, "dateOfBirthIsOver100", data.DobWarning)
			},
		},
		"other dob warning ignored": {
			form: url.Values{
				"first-names":         {"John"},
				"last-name":           {"Doe"},
				"date-of-birth-day":   {"2"},
				"date-of-birth-month": {"1"},
				"date-of-birth-year":  {"1900"},
				"ignore-dob-warning":  {"dateOfBirthIsUnder18"},
			},
			dataMatcher: func(t *testing.T, data *yourDetailsData) bool {
				return assert.Equal(t, "dateOfBirthIsOver100", data.DobWarning)
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r, _ := http.NewRequest(http.MethodPost, "/", strings.NewReader(tc.form.Encode()))
			r.Header.Add("Content-Type", page.FormUrlEncoded)

			template := newMockTemplate(t)
			template.
				On("Execute", w, mock.MatchedBy(func(data *yourDetailsData) bool {
					return tc.dataMatcher(t, data)
				})).
				Return(nil)

			sessionStore := newMockSessionStore(t)
			sessionStore.
				On("Get", mock.Anything, "session").
				Return(&sessions.Session{Values: map[any]any{"session": &sesh.LoginSession{Sub: "xyz", Email: "name@example.com"}}}, nil)

			err := YourDetails(template.Execute, nil, sessionStore)(testAppData, w, r, &page.Lpa{})
			resp := w.Result()

			assert.Nil(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
		})
	}
}

func TestPostYourDetailsWhenStoreErrors(t *testing.T) {
	form := url.Values{
		"first-names":         {"John"},
		"last-name":           {"Doe"},
		"date-of-birth-day":   {"2"},
		"date-of-birth-month": {"1"},
		"date-of-birth-year":  {"1990"},
	}

	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodPost, "/", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", page.FormUrlEncoded)

	donorStore := newMockDonorStore(t)
	donorStore.
		On("Put", r.Context(), mock.Anything).
		Return(expectedError)

	sessionStore := newMockSessionStore(t)
	sessionStore.
		On("Get", mock.Anything, "session").
		Return(&sessions.Session{Values: map[any]any{"session": &sesh.LoginSession{Sub: "xyz", Email: "name@example.com"}}}, nil)

	err := YourDetails(nil, donorStore, sessionStore)(testAppData, w, r, &page.Lpa{
		Donor: actor.Donor{
			FirstNames: "John",
			Address:    place.Address{Line1: "abc"},
		},
	})

	assert.Equal(t, expectedError, err)
}

func TestPostYourDetailsWhenSessionProblem(t *testing.T) {
	testCases := map[string]struct {
		session *sessions.Session
		error   error
	}{
		"store error": {
			session: &sessions.Session{Values: map[any]any{"session": &sesh.LoginSession{Sub: "xyz", Email: "name@example.com"}}},
			error:   expectedError,
		},
		"missing donor session": {
			session: &sessions.Session{},
			error:   nil,
		},
		"missing email": {
			session: &sessions.Session{Values: map[any]any{"session": &sesh.LoginSession{Sub: "xyz"}}},
			error:   nil,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			form := url.Values{
				"first-names":         {"John"},
				"last-name":           {"Doe"},
				"date-of-birth-day":   {"2"},
				"date-of-birth-month": {"1"},
				"date-of-birth-year":  {"1990"},
			}

			w := httptest.NewRecorder()
			r, _ := http.NewRequest(http.MethodPost, "/", strings.NewReader(form.Encode()))
			r.Header.Add("Content-Type", page.FormUrlEncoded)

			sessionStore := newMockSessionStore(t)
			sessionStore.
				On("Get", mock.Anything, "session").
				Return(tc.session, tc.error)

			err := YourDetails(nil, nil, sessionStore)(testAppData, w, r, &page.Lpa{})

			assert.NotNil(t, err)
		})
	}
}

func TestReadYourDetailsForm(t *testing.T) {
	assert := assert.New(t)

	form := url.Values{
		"first-names":         {"  John "},
		"last-name":           {"Doe"},
		"other-names":         {"Somebody"},
		"date-of-birth-day":   {"2"},
		"date-of-birth-month": {"1"},
		"date-of-birth-year":  {"1990"},
		"ignore-dob-warning":  {"xyz"},
	}

	r, _ := http.NewRequest(http.MethodPost, "/", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", page.FormUrlEncoded)

	result := readYourDetailsForm(r)

	assert.Equal("John", result.FirstNames)
	assert.Equal("Doe", result.LastName)
	assert.Equal("Somebody", result.OtherNames)
	assert.Equal(date.New("1990", "1", "2"), result.Dob)
	assert.Equal("xyz", result.IgnoreDobWarning)
}

func TestYourDetailsFormValidate(t *testing.T) {
	now := date.Today()
	validDob := now.AddDate(-18, 0, -1)

	testCases := map[string]struct {
		form   *yourDetailsForm
		errors validation.List
	}{
		"valid": {
			form: &yourDetailsForm{
				FirstNames: "A",
				LastName:   "B",
				Dob:        validDob,
			},
		},
		"max-length": {
			form: &yourDetailsForm{
				FirstNames: strings.Repeat("x", 53),
				LastName:   strings.Repeat("x", 61),
				OtherNames: strings.Repeat("x", 50),
				Dob:        validDob,
			},
		},
		"missing-all": {
			form: &yourDetailsForm{},
			errors: validation.
				With("first-names", validation.EnterError{Label: "firstNames"}).
				With("last-name", validation.EnterError{Label: "lastName"}).
				With("date-of-birth", validation.EnterError{Label: "dateOfBirth"}),
		},
		"too-long": {
			form: &yourDetailsForm{
				FirstNames: strings.Repeat("x", 54),
				LastName:   strings.Repeat("x", 62),
				OtherNames: strings.Repeat("x", 51),
				Dob:        validDob,
			},
			errors: validation.
				With("first-names", validation.StringTooLongError{Label: "firstNames", Length: 53}).
				With("last-name", validation.StringTooLongError{Label: "lastName", Length: 61}).
				With("other-names", validation.StringTooLongError{Label: "otherNamesLabel", Length: 50}),
		},
		"future-dob": {
			form: &yourDetailsForm{
				FirstNames: "A",
				LastName:   "B",
				Dob:        now.AddDate(0, 0, 1),
			},
			errors: validation.With("date-of-birth", validation.DateMustBePastError{Label: "dateOfBirth"}),
		},
		"invalid-dob": {
			form: &yourDetailsForm{
				FirstNames: "A",
				LastName:   "B",
				Dob:        date.New("2000", "22", "2"),
			},
			errors: validation.With("date-of-birth", validation.DateMustBeRealError{Label: "dateOfBirth"}),
		},
		"invalid-missing-dob": {
			form: &yourDetailsForm{
				FirstNames: "A",
				LastName:   "B",
				Dob:        date.New("1", "", "1"),
			},
			errors: validation.With("date-of-birth", validation.DateMissingError{Label: "dateOfBirth", MissingMonth: true}),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.errors, tc.form.Validate())
		})
	}
}

func TestYourDetailsFormDobWarning(t *testing.T) {
	now := date.Today()
	validDob := now.AddDate(-18, 0, -1)

	testCases := map[string]struct {
		form    *yourDetailsForm
		warning string
	}{
		"valid": {
			form: &yourDetailsForm{
				Dob: validDob,
			},
		},
		"future-dob": {
			form: &yourDetailsForm{
				Dob: now.AddDate(0, 0, 1),
			},
		},
		"dob-is-18": {
			form: &yourDetailsForm{
				Dob: now.AddDate(-18, 0, 0),
			},
		},
		"dob-under-18": {
			form: &yourDetailsForm{
				Dob: now.AddDate(-18, 0, 1),
			},
			warning: "dateOfBirthIsUnder18",
		},
		"dob-is-100": {
			form: &yourDetailsForm{
				Dob: now.AddDate(-100, 0, 0),
			},
		},
		"dob-over-100": {
			form: &yourDetailsForm{
				Dob: now.AddDate(-100, 0, -1),
			},
			warning: "dateOfBirthIsOver100",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.warning, tc.form.DobWarning())
		})
	}
}

func TestDonorMatches(t *testing.T) {
	lpa := &page.Lpa{
		Donor: actor.Donor{FirstNames: "a", LastName: "b"},
		Attorneys: actor.Attorneys{Attorneys: []actor.Attorney{
			{FirstNames: "c", LastName: "d"},
			{FirstNames: "e", LastName: "f"},
		}},
		ReplacementAttorneys: actor.Attorneys{Attorneys: []actor.Attorney{
			{FirstNames: "g", LastName: "h"},
			{FirstNames: "i", LastName: "j"},
		}},
		CertificateProvider: actor.CertificateProvider{FirstNames: "k", LastName: "l"},
		PeopleToNotify: actor.PeopleToNotify{
			{FirstNames: "m", LastName: "n"},
			{FirstNames: "o", LastName: "p"},
		},
	}

	assert.Equal(t, actor.TypeNone, donorMatches(lpa, "x", "y"))
	assert.Equal(t, actor.TypeNone, donorMatches(lpa, "a", "b"))
	assert.Equal(t, actor.TypeAttorney, donorMatches(lpa, "C", "D"))
	assert.Equal(t, actor.TypeAttorney, donorMatches(lpa, "e", "f"))
	assert.Equal(t, actor.TypeReplacementAttorney, donorMatches(lpa, "G", "H"))
	assert.Equal(t, actor.TypeReplacementAttorney, donorMatches(lpa, "i", "j"))
	assert.Equal(t, actor.TypeCertificateProvider, donorMatches(lpa, "k", "l"))
	assert.Equal(t, actor.TypePersonToNotify, donorMatches(lpa, "m", "n"))
	assert.Equal(t, actor.TypePersonToNotify, donorMatches(lpa, "O", "P"))
}

func TestDonorMatchesEmptyNamesIgnored(t *testing.T) {
	lpa := &page.Lpa{
		Donor: actor.Donor{FirstNames: "", LastName: ""},
		Attorneys: actor.Attorneys{Attorneys: []actor.Attorney{
			{FirstNames: "", LastName: ""},
		}},
		ReplacementAttorneys: actor.Attorneys{Attorneys: []actor.Attorney{
			{FirstNames: "", LastName: ""},
		}},
		CertificateProvider: actor.CertificateProvider{FirstNames: "", LastName: ""},
		PeopleToNotify: actor.PeopleToNotify{
			{FirstNames: "", LastName: ""},
		},
	}

	assert.Equal(t, actor.TypeNone, donorMatches(lpa, "", ""))
}