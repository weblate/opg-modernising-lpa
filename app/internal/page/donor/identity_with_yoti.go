package donor

import (
	"net/http"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-modernising-lpa/app/internal/page"
	"github.com/ministryofjustice/opg-modernising-lpa/app/internal/sesh"
	"github.com/ministryofjustice/opg-modernising-lpa/app/internal/validation"
)

type identityWithYotiData struct {
	App         page.AppData
	Errors      validation.List
	ClientSdkID string
	ScenarioID  string
}

func IdentityWithYoti(tmpl template.Template, sessionStore SessionStore, yotiClient YotiClient) Handler {
	return func(appData page.AppData, w http.ResponseWriter, r *http.Request, lpa *page.Lpa) error {
		if lpa.DonorIdentityConfirmed() || yotiClient.IsTest() {
			return appData.Redirect(w, r, lpa, page.Paths.IdentityWithYotiCallback.Format(lpa.ID))
		}

		if err := sesh.SetYoti(sessionStore, r, w, &sesh.YotiSession{
			Locale: appData.Lang.String(),
			LpaID:  appData.LpaID,
		}); err != nil {
			return err
		}

		data := &identityWithYotiData{
			App:         appData,
			ClientSdkID: yotiClient.SdkID(),
			ScenarioID:  yotiClient.ScenarioID(),
		}

		return tmpl(w, data)
	}
}
