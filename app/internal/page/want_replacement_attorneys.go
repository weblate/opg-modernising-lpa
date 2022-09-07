package page

import (
	"net/http"

	"github.com/ministryofjustice/opg-go-common/template"
)

type wantReplacementAttorneysData struct {
	App    AppData
	Errors map[string]string
	Want   string
}

func WantReplacementAttorneys(tmpl template.Template, dataStore DataStore) Handler {
	return func(appData AppData, w http.ResponseWriter, r *http.Request) error {
		var lpa Lpa
		if err := dataStore.Get(r.Context(), appData.SessionID, &lpa); err != nil {
			return err
		}

		data := &wantReplacementAttorneysData{
			App:  appData,
			Want: lpa.WantReplacementAttorneys,
		}

		if r.Method == http.MethodPost {
			form := readWantReplacementAttorneysForm(r)
			data.Errors = form.Validate()

			if len(data.Errors) == 0 {
				lpa.WantReplacementAttorneys = form.Want
				if err := dataStore.Put(r.Context(), appData.SessionID, lpa); err != nil {
					return err
				}
				appData.Lang.Redirect(w, r, taskListPath, http.StatusFound)
				return nil
			}
		}

		return tmpl(w, data)
	}
}

type wantReplacementAttorneysForm struct {
	Want string
}

func readWantReplacementAttorneysForm(r *http.Request) *wantReplacementAttorneysForm {
	return &wantReplacementAttorneysForm{
		Want: postFormString(r, "want"),
	}
}

func (f *wantReplacementAttorneysForm) Validate() map[string]string {
	errors := map[string]string{}

	if f.Want != "yes" && f.Want != "no" {
		errors["want"] = "selectWantReplacementAttorneys"
	}

	return errors
}