package signin

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
)

var b64 = base64.URLEncoding.WithPadding(base64.NoPadding)

type TokenRequestBody struct {
	GrantType           string `json:"grant_type"`
	AuthorizationCode   string `json:"code"`
	RedirectUri         string `json:"redirect_uri"`
	ClientAssertionType string `json:"client_assertion_type"`
	ClientAssertion     string `json:"client_assertion"`
}

type TokenResponseBody struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	IdToken      string `json:"id_token"`
}

func (c *Client) GetToken(redirectUri, clientID, JTI, code string) (string, error) {
	log.Println("GetToken()")

	data, _ := json.Marshal(map[string]interface{}{
		"aud": []string{"https://oidc.integration.account.gov.uk/token"},
		"iss": clientID,
		"sub": clientID,
		"exp": time.Now().Add(5 * time.Minute).Unix(),
		"jti": JTI,
		"iat": time.Now().Unix(),
	})

	privateKey, err := c.secretsClient.PrivateKey()
	if err != nil {
		return "", err
	}

	signedAssertion := signJwt(string(data), privateKey)

	body := &TokenRequestBody{
		GrantType:           "authorization_code",
		AuthorizationCode:   code,
		RedirectUri:         redirectUri,
		ClientAssertionType: "urn:ietf:params:oauth:client-assertion-type:jwt-bearer",
		ClientAssertion:     signedAssertion,
	}

	encodedPostBody := new(bytes.Buffer)
	err = json.NewEncoder(encodedPostBody).Encode(body)

	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", c.DiscoverData.TokenEndpoint, encodedPostBody)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	pubKey, err := c.secretsClient.PublicKey()
	if err != nil {
		return "", err
	}

	var tokenResponse TokenResponseBody

	err = json.NewDecoder(res.Body).Decode(&tokenResponse)
	if err != nil {
		return "", err
	}

	_, err = jwt.Parse(tokenResponse.IdToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return pubKey, nil
	})

	return tokenResponse.IdToken, err
}

func signJwt(data string, privateKey *rsa.PrivateKey) string {
	header := `{"alg":"RS256"}`

	toSign := b64.EncodeToString([]byte(header)) + "." + b64.EncodeToString([]byte(data))

	digest := sha256.Sum256([]byte(toSign))
	sig, err := privateKey.Sign(rand.Reader, digest[:], crypto.SHA256)
	if err != nil {
		panic(err)
	}

	return toSign + "." + b64.EncodeToString(sig)
}
