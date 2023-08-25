package certificateprovider

import (
	"net/http"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-modernising-lpa/internal/actor"
	"github.com/ministryofjustice/opg-modernising-lpa/internal/page"
	"github.com/ministryofjustice/opg-modernising-lpa/internal/validation"
)

type confirmYourDetailsData struct {
	App                 page.AppData
	Errors              validation.List
	Lpa                 *page.Lpa
	CertificateProvider *actor.CertificateProviderProvidedDetails
}

func ConfirmYourDetails(tmpl template.Template, donorStore DonorStore, certificateProviderStore CertificateProviderStore) page.Handler {
	return func(appData page.AppData, w http.ResponseWriter, r *http.Request) error {
		certificateProvider, err := certificateProviderStore.Get(r.Context())
		if err != nil {
			return err
		}

		lpa, err := donorStore.GetAny(r.Context())
		if err != nil {
			return err
		}

		if r.Method == http.MethodPost {
			redirect := page.Paths.CertificateProvider.YourRole.Format(certificateProvider.LpaID)
			if certificateProvider.Tasks.ConfirmYourDetails.Completed() || !lpa.Submitted.IsZero() {
				redirect = page.Paths.CertificateProvider.TaskList.Format(certificateProvider.LpaID)
			}

			certificateProvider.Tasks.ConfirmYourDetails = actor.TaskCompleted

			if err := certificateProviderStore.Put(r.Context(), certificateProvider); err != nil {
				return err
			}

			return appData.Redirect(w, r, nil, redirect)
		}

		data := &confirmYourDetailsData{
			App:                 appData,
			CertificateProvider: certificateProvider,
			Lpa:                 lpa,
		}

		return tmpl(w, data)
	}
}