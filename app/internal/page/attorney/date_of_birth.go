package attorney

import (
	"net/http"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-modernising-lpa/internal/actor"
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

		attorney, ok := lpa.AttorneyProvidedDetails.Get(appData.AttorneyID)
		if !ok {
			attorney = actor.Attorney{ID: appData.AttorneyID}
			lpa.AttorneyProvidedDetails = append(lpa.AttorneyProvidedDetails, attorney)
		}

		data := &dateOfBirthData{
			App: appData,
			Lpa: lpa,
			Form: &dateOfBirthForm{
				Dob: attorney.DateOfBirth,
			},
		}

		if r.Method == http.MethodPost {
			data.Form = readDateOfBirthForm(r)
			data.Errors = data.Form.Validate()
			dobWarning := data.Form.DobWarning()

			if data.Errors.Any() || data.Form.IgnoreDobWarning != dobWarning {
				data.DobWarning = dobWarning
			}

			if data.Errors.None() && data.DobWarning == "" {
				attorney.DateOfBirth = data.Form.Dob
				lpa.AttorneyProvidedDetails.Put(attorney)

				if err := lpaStore.Put(r.Context(), lpa); err != nil {
					return err
				}

				return appData.Redirect(w, r, lpa, page.Paths.Attorney.NextPage)
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

func (f *dateOfBirthForm) Validate() validation.List {
	var errors validation.List

	errors.Date("date-of-birth", "dateOfBirth", f.Dob,
		validation.DateMissing(),
		validation.DateMustBeReal(),
		validation.DateMustBePast())

	if f.Dob.After(date.Today().AddDate(-18, 0, 0)) {
		errors.Add("date-of-birth", validation.CustomError{Label: "youAttorneyAreUnder18Error"})
	}

	return errors
}
