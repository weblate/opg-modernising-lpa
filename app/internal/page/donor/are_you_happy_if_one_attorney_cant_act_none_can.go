package donor

import (
	"net/http"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-modernising-lpa/internal/form"
	"github.com/ministryofjustice/opg-modernising-lpa/internal/page"
	"github.com/ministryofjustice/opg-modernising-lpa/internal/validation"
)

type areYouHappyIfOneAttorneyCantActNoneCanData struct {
	App    page.AppData
	Errors validation.List
	Happy  string
}

func AreYouHappyIfOneAttorneyCantActNoneCan(tmpl template.Template, donorStore DonorStore) page.Handler {
	return func(appData page.AppData, w http.ResponseWriter, r *http.Request) error {
		lpa, err := donorStore.Get(r.Context())
		if err != nil {
			return err
		}

		data := &areYouHappyIfOneAttorneyCantActNoneCanData{
			App:   appData,
			Happy: lpa.AttorneyDecisions.HappyIfOneCannotActNoneCan,
		}

		if r.Method == http.MethodPost {
			form := form.ReadHappyForm(r)
			data.Errors = form.Validate("yesIfYouAreHappyIfOneAttorneyCantActNoneCan")

			if data.Errors.None() {
				lpa.AttorneyDecisions.HappyIfOneCannotActNoneCan = form.Happy
				lpa.Tasks.ChooseAttorneys = page.ChooseAttorneysState(lpa.Attorneys, lpa.AttorneyDecisions)
				lpa.Tasks.ChooseReplacementAttorneys = page.ChooseReplacementAttorneysState(lpa)

				if err := donorStore.Put(r.Context(), lpa); err != nil {
					return err
				}

				if form.Happy == "yes" {
					return appData.Redirect(w, r, lpa, page.Paths.DoYouWantReplacementAttorneys)
				} else {
					return appData.Redirect(w, r, lpa, page.Paths.AreYouHappyIfRemainingAttorneysCanContinueToAct)
				}
			}
		}

		return tmpl(w, data)
	}
}
