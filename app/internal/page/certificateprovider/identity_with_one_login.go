package certificateprovider

import (
	"net/http"

	"github.com/ministryofjustice/opg-modernising-lpa/app/internal/localize"
	"github.com/ministryofjustice/opg-modernising-lpa/app/internal/page"
	"github.com/ministryofjustice/opg-modernising-lpa/app/internal/sesh"
)

func IdentityWithOneLogin(logger Logger, oneLoginClient OneLoginClient, store sesh.Store, randomString func(int) string) page.Handler {
	return func(appData page.AppData, w http.ResponseWriter, r *http.Request) error {
		locale := ""
		if appData.Lang == localize.Cy {
			locale = "cy"
		}

		state := randomString(12)
		nonce := randomString(12)

		authCodeURL := oneLoginClient.AuthCodeURL(state, nonce, locale, true)

		if err := sesh.SetOneLogin(store, r, w, &sesh.OneLoginSession{
			State:    state,
			Nonce:    nonce,
			Locale:   locale,
			Redirect: page.Paths.CertificateProvider.IdentityWithOneLoginCallback,
		}); err != nil {
			logger.Print(err)
			return nil
		}

		http.Redirect(w, r, authCodeURL, http.StatusFound)
		return nil
	}
}
