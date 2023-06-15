package attorney

import (
	"net/http"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-modernising-lpa/app/internal/actor"
	"github.com/ministryofjustice/opg-modernising-lpa/app/internal/page"
	"github.com/ministryofjustice/opg-modernising-lpa/app/internal/validation"
)

type guidanceData struct {
	App    page.AppData
	Errors validation.List
	Lpa    *page.Lpa
}

func Guidance(tmpl template.Template, donorStore DonorStore) Handler {
	return func(appData page.AppData, w http.ResponseWriter, r *http.Request, _ *actor.AttorneyProvidedDetails) error {
		data := &guidanceData{
			App: appData,
		}

		if donorStore != nil {
			lpa, err := donorStore.GetAny(r.Context())
			if err != nil {
				return err
			}
			data.Lpa = lpa
		}

		return tmpl(w, data)
	}
}
