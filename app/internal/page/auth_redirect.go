package page

import (
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/ministryofjustice/opg-modernising-lpa/internal/signin"
)

type authRedirectClient interface {
	Exchange(code, nonce string) (string, error)
	UserInfo(string) (signin.UserInfo, error)
}

func AuthRedirect(logger Logger, c authRedirectClient, store sessions.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := store.Get(r, "params")
		if err != nil {
			logger.Print(err)
			return
		}

		if s, ok := session.Values["state"].(string); !ok || s != r.FormValue("state") {
			logger.Print("state missing from session or incorrect")
			return
		}

		nonce, ok := session.Values["nonce"].(string)
		if !ok {
			logger.Print("nonce missing from session")
			return
		}

		jwt, err := c.Exchange(r.FormValue("code"), nonce)
		if err != nil {
			logger.Print(err)
			return
		}

		userInfo, err := c.UserInfo(jwt)
		if err != nil {
			logger.Print(err)
			return
		}

		session.Values = map[interface{}]interface{}{"email": userInfo.Email}
		if err := store.Save(r, w, session); err != nil {
			logger.Print(err)
			return
		}

		http.Redirect(w, r, lpaTypePath, http.StatusFound)
	}
}
