package page

import (
	"net/http"

	"github.com/ministryofjustice/opg-go-common/template"
)

type yourChosenIdentityOptionsData struct {
	App          AppData
	Errors       map[string]string
	Selected     []IdentityOption
	FirstChoice  IdentityOption
	SecondChoice IdentityOption
	You          Person
}

func YourChosenIdentityOptions(tmpl template.Template, dataStore DataStore) Handler {
	return func(appData AppData, w http.ResponseWriter, r *http.Request) error {
		var lpa Lpa
		if err := dataStore.Get(r.Context(), appData.SessionID, &lpa); err != nil {
			return err
		}

		if r.Method == http.MethodPost {
			appData.Lang.Redirect(w, r, lpa.IdentityOptions.NextPath(IdentityOptionUnknown), http.StatusFound)
			return nil
		}

		data := &yourChosenIdentityOptionsData{
			App:          appData,
			Selected:     lpa.IdentityOptions.Selected,
			FirstChoice:  lpa.IdentityOptions.First,
			SecondChoice: lpa.IdentityOptions.Second,
			You:          lpa.You,
		}

		return tmpl(w, data)
	}
}
