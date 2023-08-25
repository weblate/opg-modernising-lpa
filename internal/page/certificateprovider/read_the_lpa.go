package certificateprovider

import (
	"net/http"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-modernising-lpa/internal/actor"
	"github.com/ministryofjustice/opg-modernising-lpa/internal/page"
	"github.com/ministryofjustice/opg-modernising-lpa/internal/validation"
)

type readTheLpaData struct {
	App                 page.AppData
	Errors              validation.List
	Lpa                 *page.Lpa
	CertificateProvider *actor.CertificateProviderProvidedDetails
}

func ReadTheLpa(tmpl template.Template, donorStore DonorStore, certificateProviderStore CertificateProviderStore) page.Handler {
	return func(appData page.AppData, w http.ResponseWriter, r *http.Request) error {
		lpa, err := donorStore.GetAny(r.Context())
		if err != nil {
			return err
		}

		certificateProvider, err := certificateProviderStore.Get(r.Context())
		if err != nil {
			return err
		}

		if r.Method == http.MethodPost {
			if lpa.Submitted.IsZero() || !lpa.Tasks.PayForLpa.IsCompleted() {
				return appData.Redirect(w, r, nil, page.Paths.CertificateProvider.TaskList.Format(lpa.ID))
			}

			certificateProvider.Tasks.ReadTheLpa = actor.TaskCompleted
			if err := certificateProviderStore.Put(r.Context(), certificateProvider); err != nil {
				return err
			}

			return appData.Redirect(w, r, nil, page.Paths.CertificateProvider.WhatHappensNext.Format(lpa.ID))
		}

		data := &readTheLpaData{
			App:                 appData,
			Lpa:                 lpa,
			CertificateProvider: certificateProvider,
		}

		return tmpl(w, data)
	}
}