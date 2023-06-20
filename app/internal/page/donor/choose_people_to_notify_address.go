package donor

import (
	"net/http"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-modernising-lpa/app/internal/actor"
	"github.com/ministryofjustice/opg-modernising-lpa/app/internal/form"
	"github.com/ministryofjustice/opg-modernising-lpa/app/internal/page"
	"github.com/ministryofjustice/opg-modernising-lpa/app/internal/place"
)

func ChoosePeopleToNotifyAddress(logger Logger, tmpl template.Template, addressClient AddressClient, donorStore DonorStore) Handler {
	return func(appData page.AppData, w http.ResponseWriter, r *http.Request, lpa *page.Lpa) error {
		personId := r.FormValue("id")
		personToNotify, found := lpa.PeopleToNotify.Get(personId)

		if found == false {
			return appData.Redirect(w, r, lpa, page.Paths.ChoosePeopleToNotify.Format(lpa.ID))
		}

		data := &chooseAddressData{
			App:        appData,
			ActorLabel: "personToNotify",
			FullName:   personToNotify.FullName(),
			ID:         personToNotify.ID,
			Form:       &form.AddressForm{},
		}

		if personToNotify.Address.Line1 != "" {
			data.Form.Action = "manual"
			data.Form.Address = &personToNotify.Address
		}

		if r.Method == http.MethodPost {
			data.Form = form.ReadAddressForm(r)
			data.Errors = data.Form.Validate(false)

			setAddress := func(address place.Address) error {
				personToNotify.Address = *data.Form.Address
				lpa.PeopleToNotify.Put(personToNotify)
				lpa.Tasks.PeopleToNotify = actor.TaskCompleted

				return donorStore.Put(r.Context(), lpa)
			}

			switch data.Form.Action {
			case "manual":
				if data.Errors.None() {
					if err := setAddress(*data.Form.Address); err != nil {
						return err
					}

					return appData.Redirect(w, r, lpa, page.Paths.ChoosePeopleToNotifySummary.Format(lpa.ID))
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

					return appData.Redirect(w, r, lpa, page.Paths.ChoosePeopleToNotifySummary.Format(lpa.ID))
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
