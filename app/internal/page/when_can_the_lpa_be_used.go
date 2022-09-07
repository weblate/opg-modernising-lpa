package page

import (
	"net/http"

	"github.com/ministryofjustice/opg-go-common/template"
)

type whenCanTheLpaBeUsedData struct {
	App       AppData
	Errors    map[string]string
	When      string
	Completed bool
}

func WhenCanTheLpaBeUsed(tmpl template.Template, dataStore DataStore) Handler {
	return func(appData AppData, w http.ResponseWriter, r *http.Request) error {
		var lpa Lpa
		if err := dataStore.Get(r.Context(), appData.SessionID, &lpa); err != nil {
			return err
		}

		data := &whenCanTheLpaBeUsedData{
			App:       appData,
			When:      lpa.WhenCanTheLpaBeUsed,
			Completed: lpa.Tasks.WhenCanTheLpaBeUsed == TaskCompleted,
		}

		if r.Method == http.MethodPost {
			form := readWhenCanTheLpaBeUsedForm(r)
			data.Errors = form.Validate()

			if len(data.Errors) == 0 || form.AnswerLater {
				if form.AnswerLater {
					lpa.Tasks.WhenCanTheLpaBeUsed = TaskInProgress
				} else {
					lpa.Tasks.WhenCanTheLpaBeUsed = TaskCompleted
					lpa.WhenCanTheLpaBeUsed = form.When
				}
				if err := dataStore.Put(r.Context(), appData.SessionID, lpa); err != nil {
					return err
				}
				appData.Lang.Redirect(w, r, restrictionsPath, http.StatusFound)
				return nil
			}
		}

		return tmpl(w, data)
	}
}

type whenCanTheLpaBeUsedForm struct {
	AnswerLater bool
	When        string
}

func readWhenCanTheLpaBeUsedForm(r *http.Request) *whenCanTheLpaBeUsedForm {
	return &whenCanTheLpaBeUsedForm{
		AnswerLater: postFormString(r, "answer-later") == "1",
		When:        postFormString(r, "when"),
	}
}

func (f *whenCanTheLpaBeUsedForm) Validate() map[string]string {
	errors := map[string]string{}

	if f.When != "when-registered" && f.When != "when-capacity-lost" {
		errors["when"] = "selectWhenCanTheLpaBeUsed"
	}

	return errors
}
