package onelogin

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/ministryofjustice/opg-modernising-lpa/app/internal/date"
	"github.com/ministryofjustice/opg-modernising-lpa/app/internal/identity"
	"github.com/ministryofjustice/opg-modernising-lpa/app/internal/secrets"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUserInfo(t *testing.T) {
	expectedUserInfo := UserInfo{Email: "email@example.com"}

	data, _ := json.Marshal(expectedUserInfo)

	httpClient := newMockHttpClient(t)
	httpClient.
		On("Do", mock.MatchedBy(func(r *http.Request) bool {
			return assert.Equal(t, http.MethodGet, r.Method) &&
				assert.Equal(t, "http://user-info", r.URL.String()) &&
				assert.Equal(t, "Bearer hey", r.Header.Get("Authorization"))
		})).
		Return(&http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewReader(data)),
		}, nil)

	c := &Client{
		httpClient: httpClient,
		openidConfiguration: openidConfiguration{
			UserinfoEndpoint: "http://user-info",
		},
	}

	userinfo, err := c.UserInfo(context.Background(), "hey")
	assert.Nil(t, err)
	assert.Equal(t, expectedUserInfo, userinfo)
}

func TestUserInfoWhenRequestError(t *testing.T) {
	httpClient := newMockHttpClient(t)
	httpClient.
		On("Do", mock.Anything).
		Return(&http.Response{}, expectedError)

	c := &Client{
		httpClient: httpClient,
		openidConfiguration: openidConfiguration{
			UserinfoEndpoint: "http://user-info",
		},
	}

	_, err := c.UserInfo(context.Background(), "hey")
	assert.Equal(t, expectedError, err)
}

func TestParseIdentityClaim(t *testing.T) {
	privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	issuedAt := time.Now().Add(-time.Minute).Round(time.Second)

	publicKeyBytes, _ := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)

	secretsClient := newMockSecretsClient(t)
	secretsClient.
		On("SecretBytes", ctx, secrets.GovUkOneLoginIdentityPublicKey).
		Return(pem.EncodeToMemory(
			&pem.Block{
				Type:  "PUBLIC KEY",
				Bytes: publicKeyBytes,
			},
		), nil)

	c := &Client{secretsClient: secretsClient}

	namePart := []map[string]any{
		{
			"validFrom": "2020-03-01",
			"nameParts": []map[string]string{
				{
					"value": "Alice",
					"type":  "GivenName",
				},
				{
					"value": "Jane",
					"type":  "GivenName",
				},
				{
					"value": "Laura",
					"type":  "GivenName",
				},
				{
					"value": "Doe",
					"type":  "FamilyName",
				},
			},
		},
		{
			"validUntil": "2020-03-01",
			"nameParts": []map[string]string{
				{
					"value": "Alice",
					"type":  "GivenName",
				},
				{
					"value": "Eod",
					"type":  "FamilyName",
				},
			},
		},
	}

	birthDatePart := []map[string]any{
		{
			"value": "1970-01-02",
		},
	}

	vc := map[string]any{
		"credentialSubject": map[string]any{
			"name":      namePart,
			"birthDate": birthDatePart,
		},
	}

	mustSign := func(token *jwt.Token, key any) string {
		s, err := token.SignedString(key)

		assert.Nil(t, err)
		return s
	}

	testcases := map[string]struct {
		token    string
		userData identity.UserData
		error    error
		jwtError uint32
	}{
		"with required claims": {
			token: mustSign(jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
				"iat": issuedAt.Unix(),
				"vc":  vc,
			}), privateKey),
			userData: identity.UserData{
				OK:          true,
				Provider:    identity.OneLogin,
				FirstNames:  "Alice Jane Laura",
				LastName:    "Doe",
				DateOfBirth: date.New("1970", "01", "02"),
				RetrievedAt: issuedAt,
			},
		},
		"missing": {
			error: errors.New("UserInfo missing CoreIdentityJWT property"),
		},
		"without name": {
			token: mustSign(jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
				"iat": issuedAt.Unix(),
			}), privateKey),
			userData: identity.UserData{OK: false, Provider: identity.OneLogin},
		},
		"without dob": {
			token: mustSign(jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
				"iat": issuedAt.Unix(),
				"vc": map[string]any{
					"credentialSubject": map[string]any{
						"name": namePart,
					},
				},
			}), privateKey),
			userData: identity.UserData{OK: false, Provider: identity.OneLogin},
		},
		"with invalid dob": {
			token: mustSign(jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
				"iat": issuedAt.Unix(),
				"vc": map[string]any{
					"credentialSubject": map[string]any{
						"name": namePart,
						"birthDate": []map[string]any{
							{
								"value": "1970-100-02",
							},
						},
					},
				},
			}), privateKey),
			userData: identity.UserData{OK: false, Provider: identity.OneLogin},
		},
		"without iat": {
			token: mustSign(jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
				"vc": vc,
			}), privateKey),
			userData: identity.UserData{OK: false, Provider: identity.OneLogin},
		},
		"with unexpected signing method": {
			token: mustSign(jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"iat": issuedAt.Unix(),
				"vc":  vc,
			}), []byte("a key")),
			jwtError: jwt.ValidationErrorSignatureInvalid,
		},
		"with malformed token": {
			token:    "what token",
			jwtError: jwt.ValidationErrorMalformed,
		},
		"with invalid token": {
			token: mustSign(jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
				"iat": time.Now().Add(time.Minute).Unix(),
				"vc":  vc,
			}), privateKey),
			jwtError: jwt.ValidationErrorIssuedAt,
		},
	}

	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			userInfo := UserInfo{
				CoreIdentityJWT: tc.token,
			}

			userData, err := c.ParseIdentityClaim(context.Background(), userInfo)
			if verr, ok := err.(*jwt.ValidationError); tc.jwtError > 0 {
				assert.True(t, ok)
				assert.Equal(t, tc.jwtError, verr.Errors)
			} else {
				assert.Equal(t, tc.error, err)
			}
			assert.Equal(t, tc.userData, userData)
		})
	}
}

func TestParseIdentityClaimWhenSecretBytesError(t *testing.T) {
	secretsClient := newMockSecretsClient(t)
	secretsClient.
		On("SecretBytes", ctx, secrets.GovUkOneLoginIdentityPublicKey).
		Return([]byte{}, expectedError)

	c := &Client{secretsClient: secretsClient}

	_, err := c.ParseIdentityClaim(context.Background(), UserInfo{})
	assert.Equal(t, expectedError, err)
}

func TestParseIdentityClaimWhenPublicKeyInvalid(t *testing.T) {
	secretsClient := newMockSecretsClient(t)
	secretsClient.
		On("SecretBytes", ctx, secrets.GovUkOneLoginIdentityPublicKey).
		Return([]byte("not a key"), nil)

	c := &Client{secretsClient: secretsClient}

	_, err := c.ParseIdentityClaim(context.Background(), UserInfo{})
	assert.NotNil(t, err)
}
