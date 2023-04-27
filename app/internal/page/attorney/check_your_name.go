package attorney

import (
	"net/http"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-modernising-lpa/internal/actor"
	"github.com/ministryofjustice/opg-modernising-lpa/internal/notify"
	"github.com/ministryofjustice/opg-modernising-lpa/internal/page"
	"github.com/ministryofjustice/opg-modernising-lpa/internal/validation"
)

type checkYourNameData struct {
	App      page.AppData
	Form     *checkYourNameForm
	Errors   validation.List
	Lpa      *page.Lpa
	Attorney actor.Attorney
}

func CheckYourName(tmpl template.Template, lpaStore LpaStore, notifyClient NotifyClient) page.Handler {
	return func(appData page.AppData, w http.ResponseWriter, r *http.Request) error {
		lpa, err := lpaStore.Get(r.Context())
		if err != nil {
			return err
		}

		attorneys := lpa.Attorneys
		if appData.IsReplacementAttorney {
			attorneys = lpa.ReplacementAttorneys
		}

		attorney, ok := attorneys.Get(appData.AttorneyID)
		if !ok {
			return appData.Redirect(w, r, lpa, page.Paths.Attorney.Start)
		}

		attorneyProvidedDetails := getProvidedDetails(appData, lpa)

		data := &checkYourNameData{
			App: appData,
			Form: &checkYourNameForm{
				IsNameCorrect: attorneyProvidedDetails.IsNameCorrect,
				CorrectedName: attorneyProvidedDetails.CorrectedName,
			},
			Lpa:      lpa,
			Attorney: attorney,
		}

		if r.Method == http.MethodPost {
			data.Form = readCheckYourNameForm(r)
			data.Errors = data.Form.Validate()

			if len(data.Errors) == 0 {
				previousCorrectedName := attorneyProvidedDetails.CorrectedName
				attorneyProvidedDetails.IsNameCorrect = data.Form.IsNameCorrect
				attorneyProvidedDetails.CorrectedName = data.Form.CorrectedName
				setProvidedDetails(appData, lpa, attorneyProvidedDetails)

				if data.Form.CorrectedName != "" && data.Form.CorrectedName != previousCorrectedName {
					attorneyProvidedDetails.CorrectedName = data.Form.CorrectedName
					_, err := notifyClient.Email(r.Context(), notify.Email{
						EmailAddress:    lpa.Donor.Email,
						TemplateID:      notifyClient.TemplateID(notify.AttorneyNameChangeEmail),
						Personalisation: map[string]string{"declaredName": attorneyProvidedDetails.CorrectedName},
					})
					if err != nil {
						return err
					}
				}

				tasks := getTasks(appData, lpa)
				if tasks.ConfirmYourDetails == page.TaskNotStarted {
					tasks.ConfirmYourDetails = page.TaskInProgress
					setTasks(appData, lpa, tasks)
				}

				if err := lpaStore.Put(r.Context(), lpa); err != nil {
					return err
				}

				appData.Redirect(w, r, lpa, page.Paths.Attorney.DateOfBirth)
				return nil
			}
		}

		return tmpl(w, data)
	}
}

type checkYourNameForm struct {
	IsNameCorrect string
	CorrectedName string
}

func readCheckYourNameForm(r *http.Request) *checkYourNameForm {

	return &checkYourNameForm{
		IsNameCorrect: page.PostFormString(r, "is-name-correct"),
		CorrectedName: page.PostFormString(r, "corrected-name"),
	}
}

func (f *checkYourNameForm) Validate() validation.List {
	errors := validation.List{}

	errors.String("is-name-correct", "yesIfTheNameIsCorrect", f.IsNameCorrect,
		validation.Select("yes", "no"))

	if f.IsNameCorrect == "no" && f.CorrectedName == "" {
		errors.String("corrected-name", "yourFullName", f.CorrectedName,
			validation.Empty(),
		)
	}

	return errors
}
