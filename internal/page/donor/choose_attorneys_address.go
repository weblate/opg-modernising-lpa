package donor

import (
	"net/http"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-modernising-lpa/internal/form"
	"github.com/ministryofjustice/opg-modernising-lpa/internal/page"
	"github.com/ministryofjustice/opg-modernising-lpa/internal/place"
)

func ChooseAttorneysAddress(logger Logger, tmpl template.Template, addressClient AddressClient, donorStore DonorStore) Handler {
	return func(appData page.AppData, w http.ResponseWriter, r *http.Request, lpa *page.Lpa) error {
		attorneyId := r.FormValue("id")
		attorney, found := lpa.Attorneys.Get(attorneyId)

		if found == false {
			return appData.Redirect(w, r, lpa, page.Paths.ChooseAttorneys.Format(lpa.ID))
		}

		data := &chooseAddressData{
			App:        appData,
			ActorLabel: "attorney",
			FullName:   attorney.FullName(),
			ID:         attorney.ID,
			CanSkip:    true,
			Form:       &form.AddressForm{},
		}

		if attorney.Address.Line1 != "" {
			data.Form.Action = "manual"
			data.Form.Address = &attorney.Address
		}

		if r.Method == http.MethodPost {
			data.Form = form.ReadAddressForm(r)
			data.Errors = data.Form.Validate(false)

			setAddress := func(address place.Address) error {
				attorney.Address = address
				lpa.Attorneys.Put(attorney)
				lpa.Tasks.ChooseAttorneys = page.ChooseAttorneysState(lpa.Attorneys, lpa.AttorneyDecisions)
				lpa.Tasks.ChooseReplacementAttorneys = page.ChooseReplacementAttorneysState(lpa)

				return donorStore.Put(r.Context(), lpa)
			}

			switch data.Form.Action {
			case "skip":
				if err := setAddress(place.Address{}); err != nil {
					return err
				}

				return appData.Redirect(w, r, lpa, page.Paths.ChooseAttorneysSummary.Format(lpa.ID))

			case "manual":
				if data.Errors.None() {
					if err := setAddress(*data.Form.Address); err != nil {
						return err
					}

					return appData.Redirect(w, r, lpa, page.Paths.ChooseAttorneysSummary.Format(lpa.ID))
				}

			case "postcode-select":
				if data.Errors.None() {
					data.Form.Action = "manual"
				} else {
					lookupAddress(r.Context(), logger, addressClient, data, false)
				}

			case "postcode-lookup":
				if data.Errors.None() {
					lookupAddress(r.Context(), logger, addressClient, data, false)
				} else {
					data.Form.Action = "postcode"
				}

			case "reuse":
				data.Addresses = lpa.ActorAddresses()

			case "reuse-select":
				if data.Errors.None() {
					if err := setAddress(*data.Form.Address); err != nil {
						return err
					}

					return appData.Redirect(w, r, lpa, page.Paths.ChooseAttorneysSummary.Format(lpa.ID))
				} else {
					data.Addresses = lpa.ActorAddresses()
				}
			}
		}

		if r.Method == http.MethodGet {
			action := r.FormValue("action")
			if action == "manual" {
				data.Form.Action = "manual"
				data.Form.Address = &place.Address{}
			}
		}

		return tmpl(w, data)
	}
}