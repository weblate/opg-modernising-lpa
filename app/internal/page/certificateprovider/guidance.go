package certificateprovider

import (
	"net/http"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-modernising-lpa/internal/actor"
	"github.com/ministryofjustice/opg-modernising-lpa/internal/page"
	"github.com/ministryofjustice/opg-modernising-lpa/internal/validation"
)

type guidanceData struct {
	App                 page.AppData
	Errors              validation.List
	Lpa                 *page.Lpa
	CertificateProvider *actor.CertificateProvider
}

func Guidance(tmpl template.Template, lpaStore LpaStore, certificateProviderStore CertificateProviderStore) page.Handler {
	return func(appData page.AppData, w http.ResponseWriter, r *http.Request) error {
		data := &guidanceData{
			App: appData,
		}

		if lpaStore != nil {
			lpa, err := lpaStore.Get(r.Context())
			if err != nil {
				return err
			}
			data.Lpa = lpa
		}

		if certificateProviderStore != nil {
			certificateProvider, err := certificateProviderStore.Get(r.Context())
			if err != nil {
				return err
			}
			data.CertificateProvider = certificateProvider
		}

		return tmpl(w, data)
	}
}