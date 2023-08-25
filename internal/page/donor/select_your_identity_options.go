package donor

import (
	"net/http"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-modernising-lpa/internal/actor"
	"github.com/ministryofjustice/opg-modernising-lpa/internal/identity"
	"github.com/ministryofjustice/opg-modernising-lpa/internal/page"
	"github.com/ministryofjustice/opg-modernising-lpa/internal/validation"
)

type selectYourIdentityOptionsData struct {
	App    page.AppData
	Errors validation.List
	Form   *selectYourIdentityOptionsForm
	Page   int
}

func SelectYourIdentityOptions(tmpl template.Template, donorStore DonorStore, pageIndex int) Handler {
	return func(appData page.AppData, w http.ResponseWriter, r *http.Request, lpa *page.Lpa) error {
		data := &selectYourIdentityOptionsData{
			App:  appData,
			Page: pageIndex,
			Form: &selectYourIdentityOptionsForm{
				Selected: lpa.DonorIdentityOption,
			},
		}

		if r.Method == http.MethodPost {
			data.Form = readSelectYourIdentityOptionsForm(r)
			data.Errors = data.Form.Validate(pageIndex)

			if data.Form.None {
				switch pageIndex {
				case 0:
					return appData.Redirect(w, r, lpa, page.Paths.SelectYourIdentityOptions1.Format(lpa.ID))
				case 1:
					return appData.Redirect(w, r, lpa, page.Paths.SelectYourIdentityOptions2.Format(lpa.ID))
				default:
					// will go to vouching flow when that is built
					return appData.Redirect(w, r, lpa, page.Paths.TaskList.Format(lpa.ID))
				}
			}

			if data.Errors.None() {
				lpa.DonorIdentityOption = data.Form.Selected
				lpa.Tasks.ConfirmYourIdentityAndSign = actor.TaskInProgress

				if err := donorStore.Put(r.Context(), lpa); err != nil {
					return err
				}

				return appData.Redirect(w, r, lpa, page.Paths.YourChosenIdentityOptions.Format(lpa.ID))
			}
		}

		return tmpl(w, data)
	}
}

type selectYourIdentityOptionsForm struct {
	Selected identity.Option
	None     bool
}

func readSelectYourIdentityOptionsForm(r *http.Request) *selectYourIdentityOptionsForm {
	option := page.PostFormString(r, "option")

	return &selectYourIdentityOptionsForm{
		Selected: identity.ReadOption(option),
		None:     option == "none",
	}
}

func (f *selectYourIdentityOptionsForm) Validate(pageIndex int) validation.List {
	var errors validation.List

	if f.Selected == identity.UnknownOption && !f.None {
		if pageIndex == 0 {
			errors.Add("option", validation.SelectError{Label: "fromTheListedOptions"})
		} else {
			errors.Add("option", validation.SelectError{Label: "whichDocumentYouWillUse"})
		}
	}

	return errors
}