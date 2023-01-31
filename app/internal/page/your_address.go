package page

import (
	"net/http"

	"github.com/ministryofjustice/opg-go-common/template"
	"github.com/ministryofjustice/opg-modernising-lpa/internal/place"
	"github.com/ministryofjustice/opg-modernising-lpa/internal/validation"
)

type yourAddressData struct {
	App       AppData
	Errors    validation.List
	Addresses []place.Address
	Form      *yourAddressForm
}

func YourAddress(logger Logger, tmpl template.Template, addressClient AddressClient, lpaStore LpaStore) Handler {
	return func(appData AppData, w http.ResponseWriter, r *http.Request) error {
		lpa, err := lpaStore.Get(r.Context())
		if err != nil {
			return err
		}

		data := &yourAddressData{
			App:  appData,
			Form: &yourAddressForm{},
		}

		if lpa.You.Address.Line1 != "" {
			data.Form.Action = "manual"
			data.Form.Address = &lpa.You.Address
		}

		if r.Method == http.MethodPost {
			data.Form = readYourAddressForm(r)
			data.Errors = data.Form.Validate()

			if data.Form.Action == "manual" && data.Errors.None() {
				lpa.You.Address = *data.Form.Address
				if err := lpaStore.Put(r.Context(), lpa); err != nil {
					return err
				}

				return appData.Redirect(w, r, lpa, Paths.WhoIsTheLpaFor)
			}

			if data.Form.Action == "select" && data.Errors.None() {
				data.Form.Action = "manual"
			}

			if data.Form.Action == "lookup" && data.Errors.None() ||
				data.Form.Action == "select" && data.Errors.Any() {
				addresses, err := addressClient.LookupPostcode(r.Context(), data.Form.LookupPostcode)
				if err != nil {
					logger.Print(err)
					data.Errors.Add("lookup-postcode", validation.CustomError{Label: "couldNotLookupPostcode"})
				}

				data.Addresses = addresses
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

type yourAddressForm struct {
	Action         string
	LookupPostcode string
	Address        *place.Address
}

func readYourAddressForm(r *http.Request) *yourAddressForm {
	f := &yourAddressForm{}
	f.Action = r.PostFormValue("action")

	switch f.Action {
	case "lookup":
		f.LookupPostcode = postFormString(r, "lookup-postcode")

	case "select":
		f.LookupPostcode = postFormString(r, "lookup-postcode")
		selectAddress := r.PostFormValue("select-address")
		if selectAddress != "" {
			f.Address = DecodeAddress(selectAddress)
		}

	case "manual":
		f.Address = &place.Address{
			Line1:      postFormString(r, "address-line-1"),
			Line2:      postFormString(r, "address-line-2"),
			Line3:      postFormString(r, "address-line-3"),
			TownOrCity: postFormString(r, "address-town"),
			Postcode:   postFormString(r, "address-postcode"),
		}
	}

	return f
}

func (f *yourAddressForm) Validate() validation.List {
	var errors validation.List

	switch f.Action {
	case "lookup":
		errors.String("lookup-postcode", "postcode", f.LookupPostcode,
			validation.Empty())

	case "select":
		errors.Address("select-address", "address", f.Address,
			validation.AddressSelected())

	case "manual":
		errors.String("address-line-1", "addressLine1", f.Address.Line1,
			validation.Empty(),
			validation.StringTooLong(50))
		errors.String("address-line-2", "addressLine2Label", f.Address.Line2,
			validation.StringTooLong(50))
		errors.String("address-line-3", "addressLine3Label", f.Address.Line3,
			validation.StringTooLong(50))
		errors.String("address-town", "townOrCity", f.Address.TownOrCity,
			validation.Empty())
	}

	return errors
}
