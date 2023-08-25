package donor

import (
	"net/http"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-modernising-lpa/internal/actor"
	"github.com/ministryofjustice/opg-modernising-lpa/internal/page"
	"github.com/ministryofjustice/opg-modernising-lpa/internal/validation"
)

type applicationReasonData struct {
	App     page.AppData
	Errors  validation.List
	Form    *applicationReasonForm
	Options page.ApplicationReasonOptions
}

func ApplicationReason(tmpl template.Template, donorStore DonorStore) Handler {
	return func(appData page.AppData, w http.ResponseWriter, r *http.Request, lpa *page.Lpa) error {
		data := &applicationReasonData{
			App: appData,
			Form: &applicationReasonForm{
				ApplicationReason: lpa.ApplicationReason,
			},
			Options: page.ApplicationReasonValues,
		}

		if r.Method == http.MethodPost {
			data.Form = readApplicationReasonForm(r)
			data.Errors = data.Form.Validate()

			if data.Errors.None() {
				redirect := page.Paths.PreviousApplicationNumber
				lpa.ApplicationReason = data.Form.ApplicationReason

				if lpa.ApplicationReason.IsNewApplication() {
					redirect = page.Paths.TaskList
					lpa.Tasks.YourDetails = actor.TaskCompleted
				}

				if err := donorStore.Put(r.Context(), lpa); err != nil {
					return err
				}

				return appData.Redirect(w, r, lpa, redirect.Format(lpa.ID))
			}
		}

		return tmpl(w, data)
	}
}

type applicationReasonForm struct {
	ApplicationReason page.ApplicationReason
	Error             error
}

func readApplicationReasonForm(r *http.Request) *applicationReasonForm {
	applicationReason, err := page.ParseApplicationReason(page.PostFormString(r, "application-reason"))

	return &applicationReasonForm{
		ApplicationReason: applicationReason,
		Error:             err,
	}
}

func (f *applicationReasonForm) Validate() validation.List {
	var errors validation.List

	errors.Error("application-reason", "theReasonForMakingTheApplication", f.Error,
		validation.Selected())

	return errors
}