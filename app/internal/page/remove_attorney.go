package page

import (
	"fmt"
	"net/http"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-modernising-lpa/internal/actor"
	"github.com/ministryofjustice/opg-modernising-lpa/internal/validation"
)

type removeAttorneyData struct {
	App      AppData
	Attorney actor.Attorney
	Errors   validation.List
	Form     *removeAttorneyForm
}

func RemoveAttorney(logger Logger, tmpl template.Template, lpaStore LpaStore) Handler {
	return func(appData AppData, w http.ResponseWriter, r *http.Request) error {
		lpa, err := lpaStore.Get(r.Context())
		if err != nil {
			logger.Print(fmt.Sprintf("error getting lpa from store: %s", err.Error()))
			return err
		}

		id := r.FormValue("id")
		attorney, found := lpa.Attorneys.Get(id)

		if found == false {
			return appData.Redirect(w, r, lpa, Paths.ChooseAttorneysSummary)
		}

		data := &removeAttorneyData{
			App:      appData,
			Attorney: attorney,
			Form:     &removeAttorneyForm{},
		}

		if r.Method == http.MethodPost {
			data.Form = readRemoveAttorneyForm(r, "yesToRemoveAttorney")
			data.Errors = data.Form.Validate()

			if data.Form.RemoveAttorney == "yes" && data.Errors.None() {
				lpa.Attorneys.Delete(attorney)
				if len(lpa.Attorneys) == 0 {
					lpa.Tasks.ChooseAttorneys = TaskInProgress
				}

				err = lpaStore.Put(r.Context(), lpa)

				if err != nil {
					logger.Print(fmt.Sprintf("error removing Attorney from LPA: %s", err.Error()))
					return err
				}

				if len(lpa.Attorneys) == 0 {
					return appData.Redirect(w, r, lpa, Paths.ChooseAttorneys)
				}

				return appData.Redirect(w, r, lpa, Paths.ChooseAttorneysSummary)
			}

			if data.Form.RemoveAttorney == "no" {
				return appData.Redirect(w, r, lpa, Paths.ChooseAttorneysSummary)
			}

		}

		return tmpl(w, data)
	}
}

type removeAttorneyForm struct {
	RemoveAttorney string
	errorLabel     string
}

func readRemoveAttorneyForm(r *http.Request, errorLabel string) *removeAttorneyForm {
	return &removeAttorneyForm{
		RemoveAttorney: postFormString(r, "remove-attorney"),
		errorLabel:     errorLabel,
	}
}

func (f *removeAttorneyForm) Validate() validation.List {
	var errors validation.List

	errors.String("remove-attorney", f.errorLabel, f.RemoveAttorney,
		validation.Select("yes", "no"))

	return errors
}
