package pay

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
)

type Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

func New(baseURL, apiKey string, httpClient *http.Client) (Client, error) {
	return Client{
		baseURL:    baseURL,
		httpClient: httpClient,
		apiKey:     apiKey,
	}, nil
}

func (c *Client) CreatePayment(body CreatePaymentBody) (CreatePaymentResponse, error) {
	data, _ := json.Marshal(body)
	reader := bytes.NewReader(data)

	req, err := http.NewRequest("POST", c.baseURL+"/v1/payments", reader)
	if err != nil {
		return CreatePaymentResponse{}, err
	}
	req.Header.Add("Authorization", "Bearer "+c.apiKey)
	req.Header.Add("Content-Type", "application/json")

	requestDump, err := httputil.DumpRequest(req, true)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(requestDump))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return CreatePaymentResponse{}, err
	}

	defer resp.Body.Close()

	var createPaymentResp CreatePaymentResponse

	if err := json.NewDecoder(resp.Body).Decode(&createPaymentResp); err != nil {
		return CreatePaymentResponse{}, err
	}

	return createPaymentResp, nil
}

func (c *Client) GetPayment(paymentId string) (GetPaymentResponse, error) {
	req, err := http.NewRequest("GET", c.baseURL+"/v1/payments/"+paymentId, nil)
	if err != nil {
		return GetPaymentResponse{}, err
	}

	req.Header.Add("Authorization", "Bearer "+c.apiKey)
	req.Header.Add("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)

	defer resp.Body.Close()

	var getPaymentResponse GetPaymentResponse

	if err := json.NewDecoder(resp.Body).Decode(&getPaymentResponse); err != nil {
		return GetPaymentResponse{}, err
	}

	return getPaymentResponse, nil
}
