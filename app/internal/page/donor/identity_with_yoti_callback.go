package donor

import (
	"net/http"
	"time"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-modernising-lpa/internal/page"
	"github.com/ministryofjustice/opg-modernising-lpa/internal/validation"
)

type identityWithYotiCallbackData struct {
	App             page.AppData
	Errors          validation.List
	FullName        string
	ConfirmedAt     time.Time
	CouldNotConfirm bool
}

func IdentityWithYotiCallback(tmpl template.Template, yotiClient YotiClient, lpaStore LpaStore) page.Handler {
	return func(appData page.AppData, w http.ResponseWriter, r *http.Request) error {
		lpa, err := lpaStore.Get(r.Context())
		if err != nil {
			return err
		}

		if r.Method == http.MethodPost {
			if lpa.DonorIdentityConfirmed() {
				return appData.Redirect(w, r, lpa, page.Paths.ReadYourLpa)
			} else {
				return appData.Redirect(w, r, lpa, page.Paths.SelectYourIdentityOptions1)
			}
		}

		data := &identityWithYotiCallbackData{App: appData}

		if lpa.DonorIdentityConfirmed() {
			data.FullName = lpa.DonorIdentityUserData.FirstNames + " " + lpa.DonorIdentityUserData.LastName
			data.ConfirmedAt = lpa.DonorIdentityUserData.RetrievedAt

			return tmpl(w, data)
		}

		user, err := yotiClient.User(r.FormValue("token"))
		if err != nil {
			return err
		}

		lpa.DonorIdentityUserData = user

		if lpa.DonorIdentityConfirmed() {
			if err := lpaStore.Put(r.Context(), lpa); err != nil {
				return err
			}

			data.FullName = user.FirstNames + " " + user.LastName
			data.ConfirmedAt = user.RetrievedAt
		} else {
			data.CouldNotConfirm = true
		}

		return tmpl(w, data)
	}
}
