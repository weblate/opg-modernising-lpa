package page

import (
	"net/http"

	"github.com/ministryofjustice/opg-go-common/template"
)

type whoDoYouWantToBeCertificateProviderGuidanceData struct {
	App        AppData
	Errors     map[string]string
	NotStarted bool
}

func WhoDoYouWantToBeCertificateProviderGuidance(tmpl template.Template, dataStore DataStore) Handler {
	return func(appData AppData, w http.ResponseWriter, r *http.Request) error {
		var lpa Lpa
		if err := dataStore.Get(r.Context(), appData.SessionID, &lpa); err != nil {
			return err
		}

		data := &whoDoYouWantToBeCertificateProviderGuidanceData{
			App:        appData,
			NotStarted: lpa.Tasks.CertificateProvider == TaskNotStarted,
		}

		if r.Method == http.MethodPost {
			if postFormString(r, "will-do-this-later") == "1" {
				appData.Lang.Redirect(w, r, taskListPath, http.StatusFound)
				return nil
			}

			if data.NotStarted {
				lpa.Tasks.CertificateProvider = TaskInProgress
			}
			if err := dataStore.Put(r.Context(), appData.SessionID, lpa); err != nil {
				return err
			}
			appData.Lang.Redirect(w, r, certificateProviderDetailsPath, http.StatusFound)
			return nil
		}

		return tmpl(w, data)
	}
}