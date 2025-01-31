package page

import (
	"net/http"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-modernising-lpa/internal/actor"
	"github.com/ministryofjustice/opg-modernising-lpa/internal/validation"
)

type howLongHaveYouKnownCertificateProviderData struct {
	App                 AppData
	Errors              validation.List
	CertificateProvider actor.CertificateProvider
	HowLong             string
}

func HowLongHaveYouKnownCertificateProvider(tmpl template.Template, lpaStore LpaStore) Handler {
	return func(appData AppData, w http.ResponseWriter, r *http.Request) error {
		lpa, err := lpaStore.Get(r.Context())
		if err != nil {
			return err
		}

		data := &howLongHaveYouKnownCertificateProviderData{
			App:                 appData,
			CertificateProvider: lpa.CertificateProvider,
			HowLong:             lpa.CertificateProvider.RelationshipLength,
		}

		if r.Method == http.MethodPost {
			form := readHowLongHaveYouKnownCertificateProviderForm(r)
			data.Errors = form.Validate()

			if data.Errors.None() {
				lpa.Tasks.CertificateProvider = TaskCompleted
				lpa.CertificateProvider.RelationshipLength = form.HowLong
				if err := lpaStore.Put(r.Context(), lpa); err != nil {
					return err
				}

				return appData.Redirect(w, r, lpa, Paths.DoYouWantToNotifyPeople)
			}
		}

		return tmpl(w, data)
	}
}

type howLongHaveYouKnownCertificateProviderForm struct {
	HowLong string
}

func readHowLongHaveYouKnownCertificateProviderForm(r *http.Request) *howLongHaveYouKnownCertificateProviderForm {
	return &howLongHaveYouKnownCertificateProviderForm{
		HowLong: postFormString(r, "how-long"),
	}
}

func (f *howLongHaveYouKnownCertificateProviderForm) Validate() validation.List {
	var errors validation.List

	errors.String("how-long", "howLongYouHaveKnownCertificateProvider", f.HowLong,
		validation.Select("gte-2-years", "lt-2-years"))

	if f.HowLong == "lt-2-years" {
		errors.Add("how-long", validation.CustomError{Label: "mustHaveKnownCertificateProviderTwoYears"})
	}

	return errors
}
