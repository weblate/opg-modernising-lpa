package attorney

import (
	"net/http"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-modernising-lpa/app/internal/actor"
	"github.com/ministryofjustice/opg-modernising-lpa/app/internal/page"
	"github.com/ministryofjustice/opg-modernising-lpa/app/internal/validation"
)

type confirmYourDetailsData struct {
	App                     page.AppData
	Errors                  validation.List
	Lpa                     *page.Lpa
	Attorney                actor.Attorney
	TrustCorporation        actor.TrustCorporation
	AttorneyProvidedDetails *actor.AttorneyProvidedDetails
}

func ConfirmYourDetails(tmpl template.Template, attorneyStore AttorneyStore, donorStore DonorStore) Handler {
	return func(appData page.AppData, w http.ResponseWriter, r *http.Request, attorneyProvidedDetails *actor.AttorneyProvidedDetails) error {
		if r.Method == http.MethodPost {
			attorneyProvidedDetails.Tasks.ConfirmYourDetails = actor.TaskCompleted

			if err := attorneyStore.Put(r.Context(), attorneyProvidedDetails); err != nil {
				return err
			}

			return appData.Redirect(w, r, nil, page.Paths.Attorney.ReadTheLpa.Format(attorneyProvidedDetails.LpaID))
		}

		lpa, err := donorStore.GetAny(r.Context())
		if err != nil {
			return err
		}

		data := &confirmYourDetailsData{
			App:                     appData,
			Lpa:                     lpa,
			AttorneyProvidedDetails: attorneyProvidedDetails,
		}

		if appData.IsTrustCorporation() {
			data.TrustCorporation = lpa.Attorneys.TrustCorporation
		} else {
			data.Attorney, _ = lpa.Attorneys.Get(attorneyProvidedDetails.ID)
		}

		return tmpl(w, data)
	}
}