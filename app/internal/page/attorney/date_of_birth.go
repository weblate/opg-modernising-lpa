package attorney

import (
	"net/http"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-modernising-lpa/internal/date"
	"github.com/ministryofjustice/opg-modernising-lpa/internal/page"
	"github.com/ministryofjustice/opg-modernising-lpa/internal/validation"
)

type dateOfBirthData struct {
	App        page.AppData
	Lpa        *page.Lpa
	Form       *dateOfBirthForm
	Errors     validation.List
	DobWarning string
}

type dateOfBirthForm struct {
	Dob              date.Date
	IgnoreDobWarning string
}

func DateOfBirth(tmpl template.Template, lpaStore LpaStore) page.Handler {
	return func(appData page.AppData, w http.ResponseWriter, r *http.Request) error {
		lpa, err := lpaStore.Get(r.Context())
		if err != nil {
			return err
		}

		attorneyProvidedDetails := getProvidedDetails(appData, lpa)

		data := &dateOfBirthData{
			App: appData,
			Lpa: lpa,
			Form: &dateOfBirthForm{
				Dob: attorneyProvidedDetails.DateOfBirth,
			},
		}

		if r.Method == http.MethodPost {
			data.Form = readDateOfBirthForm(r)
			data.Errors = data.Form.Validate(appData.IsReplacementAttorney())
			dobWarning := data.Form.DobWarning()

			if data.Errors.Any() || data.Form.IgnoreDobWarning != dobWarning {
				data.DobWarning = dobWarning
			}

			if data.Errors.None() && data.DobWarning == "" {
				attorneyProvidedDetails.DateOfBirth = data.Form.Dob
				setProvidedDetails(appData, lpa, attorneyProvidedDetails)

				if err := lpaStore.Put(r.Context(), lpa); err != nil {
					return err
				}

				return appData.Redirect(w, r, lpa, page.Paths.Attorney.MobileNumber)
			}
		}

		return tmpl(w, data)
	}
}

func readDateOfBirthForm(r *http.Request) *dateOfBirthForm {
	return &dateOfBirthForm{
		Dob:              date.New(page.PostFormString(r, "date-of-birth-year"), page.PostFormString(r, "date-of-birth-month"), page.PostFormString(r, "date-of-birth-day")),
		IgnoreDobWarning: page.PostFormString(r, "ignore-dob-warning"),
	}
}

func (f *dateOfBirthForm) DobWarning() string {
	var (
		hundredYearsEarlier = date.Today().AddDate(-100, 0, 0)
	)

	if !f.Dob.IsZero() {
		if f.Dob.Before(hundredYearsEarlier) {
			return "dateOfBirthIsOver100"
		}
	}

	return ""
}

func (f *dateOfBirthForm) Validate(isReplacement bool) validation.List {
	var errors validation.List

	errors.Date("date-of-birth", "yourDateOfBirth", f.Dob,
		validation.DateMissing(),
		validation.DateMustBeReal(),
		validation.DateMustBePast())

	if f.Dob.After(date.Today().AddDate(-18, 0, 0)) {
		if isReplacement {
			errors.Add("date-of-birth", validation.CustomError{Label: "youReplacementAttorneyAreUnder18Error"})
		} else {
			errors.Add("date-of-birth", validation.CustomError{Label: "youAttorneyAreUnder18Error"})
		}
	}

	return errors
}
