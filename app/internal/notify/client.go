package notify

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type Doer interface {
	Do(*http.Request) (*http.Response, error)
}

type Client struct {
	isProduction bool
	baseURL      string
	doer         Doer
	issuer       string
	secretKey    []byte
	now          func() time.Time
}

func New(isProduction bool, baseURL, apiKey string, httpClient Doer) (*Client, error) {
	keyParts := strings.Split(apiKey, "-")
	if len(keyParts) != 11 {
		return nil, errors.New("invalid apiKey format")
	}

	if baseURL == "" {
		baseURL = "https://api.notifications.service.gov.uk"
	}

	return &Client{
		isProduction: isProduction,
		baseURL:      baseURL,
		doer:         httpClient,
		issuer:       strings.Join(keyParts[1:6], "-"),
		secretKey:    []byte(strings.Join(keyParts[6:11], "-")),
		now:          time.Now,
	}, nil
}

type Email struct {
	EmailAddress    string            `json:"email_address"`
	TemplateID      string            `json:"template_id"`
	Personalisation map[string]string `json:"personalisation,omitempty"`
	Reference       string            `json:"reference,omitempty"`
	EmailReplyToID  string            `json:"email_reply_to_id,omitempty"`
}

type Sms struct {
	PhoneNumber     string            `json:"phone_number"`
	TemplateID      string            `json:"template_id"`
	Personalisation map[string]string `json:"personalisation,omitempty"`
	Reference       string            `json:"reference,omitempty"`
}

type response struct {
	ID         string     `json:"id"`
	StatusCode int        `json:"status_code,omitempty"`
	Errors     errorsList `json:"errors,omitempty"`
}

type errorsList []errorItem

func (es errorsList) Error() string {
	s := "error sending message"
	for _, e := range es {
		s += ": " + e.Message
	}
	return s
}

type errorItem struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

type TemplateId int

const (
	SignatureCodeEmail TemplateId = iota
	SignatureCodeSms
)

func (c *Client) TemplateID(id TemplateId) string {
	if c.isProduction {
		switch id {
		case SignatureCodeEmail:
			return "95f7b0a2-1c3a-4ad9-818b-b358c549c88b"
		case SignatureCodeSms:
			return "a0997cbf-cfd9-4f01-acb2-f33b07074662"
		}
	} else {
		switch id {
		case SignatureCodeEmail:
			return "7e8564a0-2635-4f61-9155-0166ddbe5607"
		case SignatureCodeSms:
			return "0aa5b61c-ef30-410a-8473-915df9d343a5"
		}
	}

	return ""
}

func (c *Client) Email(ctx context.Context, email Email) (string, error) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(email); err != nil {
		return "", err
	}

	req, err := c.request(ctx, "/v2/notifications/email", buf)

	resp, err := c.doer.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var r response
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return "", err
	}

	if len(r.Errors) > 0 {
		return "", r.Errors
	}

	return r.ID, nil
}

func (c *Client) Sms(ctx context.Context, sms Sms) (string, error) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(sms); err != nil {
		return "", err
	}

	req, err := c.request(ctx, "/v2/notifications/sms", buf)

	resp, err := c.doer.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var r response
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return "", err
	}

	if len(r.Errors) > 0 {
		return "", r.Errors
	}

	return r.ID, nil
}

func (c *Client) request(ctx context.Context, url string, jsonBody bytes.Buffer) (*http.Request, error) {
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.RegisteredClaims{
		Issuer:   c.issuer,
		IssuedAt: jwt.NewNumericDate(c.now()),
	}).SignedString(c.secretKey)
	if err != nil {
		return &http.Request{}, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+url, &jsonBody)
	if err != nil {
		return &http.Request{}, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+token)

	return req, nil
}
