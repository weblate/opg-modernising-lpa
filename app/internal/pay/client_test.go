package pay

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	amount      = 5
	reference   = "abc123"
	description = "A payment"
	returnUrl   = "/example/url"
	email       = "a@example.org"
	language    = "en"
	created     = time.Date(2022, time.January, 5, 0, 0, 0, 0, time.UTC)
	apiToken    = "fake-token"
)

func TestCreatePayment(t *testing.T) {
	t.Run("POSTs required body content to expected GOVUK Pay create payment endpoint", func(t *testing.T) {
		body := CreatePaymentBody{
			Amount:      amount,
			Reference:   reference,
			Description: description,
			ReturnUrl:   returnUrl,
			Email:       email,
			Language:    language,
		}

		expectedCPResponse := CreatePaymentResponse{
			CreatedDate: created,
			State: State{
				Status:   "created",
				Finished: false,
			},
			Links: map[string]Link{
				"self": {
					Href:   "https://publicapi.payments.service.gov.uk/v1/payments/hu20sqlact5260q2nanm0q8u93",
					Method: "GET",
				},
				"next_url": {
					Href:   "https://www.payments.service.gov.uk/secure/bb0a272c-8eaf-468d-b3xf-ae5e000d2231",
					Method: "GET",
				},
			},
			Amount:          amount,
			Reference:       reference,
			Description:     description,
			ReturnUrl:       returnUrl,
			PaymentId:       "hu20sqlact5260q2nanm0q8u93",
			PaymentProvider: "worldpay",
			ProviderId:      "10987654321",
		}

		server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			defer req.Body.Close()

			reqBody, _ := io.ReadAll(req.Body)
			expectedReqBody := `{"amount": 5,"reference" : "abc123","description": "A payment","return_url": "/example/url","email": "a@example.org","language": "en"}`

			assert.Equal(t, req.URL.String(), "/v1/payments", "URL did not match")
			assert.Equal(t, req.Header.Get("Authorization"), "Bearer fake-token", "Authorization token did not match")
			assert.JSONEq(t, expectedReqBody, string(reqBody), "Request body did not match")

			rw.WriteHeader(http.StatusCreated)
			rw.Write(generateCreatePaymentResponseBodyJsonString())
		}))

		defer server.Close()

		payClient, _ := New(server.URL, apiToken, server.Client())

		actualCPResponse, err := payClient.CreatePayment(body)
		if err != nil {
			t.Fatal(err, "An error unexpectedly occurred during CreatePayment")
		}

		assert.Equal(t, expectedCPResponse, actualCPResponse, "Return value did not match")
	})

	t.Run("Returns an error if unable to create a request", func(t *testing.T) {
		body := CreatePaymentBody{
			Amount:      amount,
			Reference:   reference,
			Description: description,
			ReturnUrl:   returnUrl,
			Email:       email,
			Language:    language,
		}

		expectedCPResponse := CreatePaymentResponse{}

		server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {}))

		defer server.Close()

		payClient, _ := New(server.URL+"`invalid-url-format", apiToken, server.Client())

		actualCPResponse, err := payClient.CreatePayment(body)
		if err == nil {
			t.Fatal("Expected an error but received nil")
		}

		assert.Equal(t, expectedCPResponse, actualCPResponse, "Return value did not match")
	})

	t.Run("Returns an error if unable to decode response", func(t *testing.T) {
		body := CreatePaymentBody{
			Amount:      amount,
			Reference:   reference,
			Description: description,
			ReturnUrl:   returnUrl,
			Email:       email,
			Language:    language,
		}

		expectedCPResponse := CreatePaymentResponse{}

		server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			defer req.Body.Close()

			assert.Equal(t, req.URL.String(), "/v1/payments", "URL did not match")

			rw.WriteHeader(http.StatusCreated)
			rw.Write([]byte("not JSON"))
		}))

		defer server.Close()

		payClient, _ := New(server.URL, apiToken, server.Client())

		actualCPResponse, err := payClient.CreatePayment(body)
		if err == nil {
			t.Fatal("Expected an error but received nil")
		}

		assert.Equal(t, expectedCPResponse, actualCPResponse, "Return value did not match")
	})
}

func generateCreatePaymentResponseBodyJsonString() []byte {
	return []byte(fmt.Sprintf(`
{
  "created_date": "%s",
  "state": {
    "status": "created",
    "finished": false
  },
  "_links": {
    "self": {
      "href": "https://publicapi.payments.service.gov.uk/v1/payments/hu20sqlact5260q2nanm0q8u93",
      "method": "GET"
   },
    "next_url": {
      "href": "https://www.payments.service.gov.uk/secure/bb0a272c-8eaf-468d-b3xf-ae5e000d2231",
      "method": "GET"
    }
  },
  "amount": 5,
  "reference" : "abc123",
  "description": "A payment",
  "return_url": "/example/url",
  "payment_id": "hu20sqlact5260q2nanm0q8u93",
  "payment_provider": "worldpay",
  "provider_id": "10987654321"
}`, created.Format(time.RFC3339)))
}

func TestGetPayment(t *testing.T) {
	paymentId := "fake-id-value"

	t.Run("GETs payment information using a payment ID", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			defer req.Body.Close()

			assert.Equal(t, req.URL.String(), fmt.Sprintf("/v1/payments/%s", paymentId), "URL did not match")
			assert.Equal(t, req.Header.Get("Authorization"), "Bearer fake-token", "Authorization token did not match")

			rw.WriteHeader(http.StatusCreated)
			rw.Write(generateGetPaymentResponseBodyJsonBytes())
		}))

		defer server.Close()

		payClient, _ := New(server.URL, apiToken, server.Client())

		actualGPResponse, err := payClient.GetPayment(paymentId)
		if err != nil {
			t.Fatal(err)
		}

		expectedGPResponse := GetPaymentResponse{
			CreatedDate: created,
			Amount:      amount,
			State: State{
				Status:   "success",
				Finished: true,
			},
			Description: description,
			Reference:   reference,
			Language:    language,
			Email:       email,
			CardDetails: CardDetails{
				CardBrand:             "Visa",
				CardType:              "debit",
				LastDigitsCardNumber:  "1234",
				FirstDigitsCardNumber: "123456",
				ExpiryDate:            "04/24",
				CardholderName:        "Sherlock Holmes",
				BillingAddress: BillingAddress{
					Line1:    "221 Baker Street",
					Line2:    "Flat b",
					Postcode: "NW1 6XE",
					City:     "London",
					Country:  "GB",
				},
			},
			PaymentId: "hu20sqlact5260q2nanm0q8u93",
			AuthorisationSummary: AuthorisationSummary{
				ThreeDSecure: ThreeDSecure{
					Required: true,
				},
			},
			RefundSummary: RefundSummary{
				Status:          "available",
				AmountAvailable: 4000,
				AmountSubmitted: 80,
			},
			SettlementSummary: SettlementSummary{
				CaptureSubmitTime: created.Format(time.RFC3339),
				CapturedDate:      "2022-01-05",
				SettledDate:       "2022-01-05",
			},
			DelayedCapture:         false,
			Moto:                   false,
			CorporateCardSurcharge: 250,
			TotalAmount:            4000,
			Fee:                    200,
			NetAmount:              3800,
			PaymentProvider:        "worldpay",
			ProviderId:             "10987654321",
			ReturnUrl:              "https://your.service.gov.uk/completed",
		}
		assert.Equal(t, expectedGPResponse, actualGPResponse, "Return value did not match")
	})

	t.Run("Returns an error if unable to create a request", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {}))

		defer server.Close()

		payClient, _ := New(server.URL+"`invalid-url-format", apiToken, server.Client())

		actualGPResponse, err := payClient.GetPayment(paymentId)

		if err == nil {
			t.Fatal("Expected an error but received nil")
		}

		assert.Equal(t, GetPaymentResponse{}, actualGPResponse, "Return value did not match")
	})

	t.Run("Returns an error if unable to decode response", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			defer req.Body.Close()

			assert.Equal(t, req.URL.String(), fmt.Sprintf("/v1/payments/%s", paymentId), "URL did not match")

			rw.WriteHeader(http.StatusCreated)
			rw.Write([]byte("still not JSON"))
		}))

		defer server.Close()

		payClient, _ := New(server.URL, apiToken, server.Client())

		actualGPResponse, err := payClient.GetPayment(paymentId)

		if err == nil {
			t.Fatal("Expected an error but received nil")
		}

		assert.Equal(t, GetPaymentResponse{}, actualGPResponse, "Return value did not match")
	})
}

func generateGetPaymentResponseBodyJsonBytes() []byte {
	return []byte(fmt.Sprintf(`
{
  "created_date": "%s",
  "amount": %v,
  "state": {
    "status": "success",
    "finished": true
  },
  "description": "%s",
  "reference": "%s",
  "language": "%s",
  "metadata": {
    "ledger_code": "AB100",
    "an_internal_reference_number": 200
  },
  "email": "%s",
  "card_details": {
    "card_brand": "Visa",
    "card_type": "debit",
    "last_digits_card_number": "1234",
    "first_digits_card_number": "123456",
    "expiry_date": "04/24",
    "cardholder_name": "Sherlock Holmes",
    "billing_address": {
        "line1": "221 Baker Street",
        "line2": "Flat b",
        "postcode": "NW1 6XE",
        "city": "London",
        "country": "GB"
    }
  },
  "payment_id": "hu20sqlact5260q2nanm0q8u93",
  "authorisation_summary": {
    "three_d_secure": {
      "required": true
    }
  },
  "refund_summary": {
    "status": "available",
    "amount_available": 4000,
    "amount_submitted": 80
  },
  "settlement_summary": {
    "capture_submit_time": "%s",
    "captured_date": "2022-01-05",
    "settled_date": "2022-01-05"
  },
  "delayed_capture": false,
  "moto": false,
  "corporate_card_surcharge": 250,
  "total_amount": 4000,
  "fee": 200,
  "net_amount": 3800,
  "payment_provider": "worldpay",
  "provider_id": "10987654321",
  "return_url": "https://your.service.gov.uk/completed"
}`,
		created.Format(time.RFC3339),
		amount,
		description,
		reference,
		language,
		email,
		created.Format(time.RFC3339)))
}