package auth

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

const ERROR_MSG_AUTH_FAIL = `Authorization failed!`
const ERROR_MSG_INVALID_ACCESS_TOKEN_TYPE = `Invalid access token type fetched!`
const ERROR_MSG_NO_ACCESS_TOKEN = `No access token fetched!`

type Authorizer interface {
	Authorize() (string, error)
}

type decoderResponse map[string]interface{}

type AuthorizerImpl struct {
	ApiUrl      string
	ApiKey      string
	ApiSecret   string
	Logger      *log.Logger
	UrlEncoding *base64.Encoding
	Ecoding     *base64.Encoding
}

func NewAuthorizer(apiUrl string, apiKey string, apiSecret string, logger *log.Logger) *AuthorizerImpl {
	return &AuthorizerImpl{
		ApiUrl:      apiUrl,
		ApiKey:      apiKey,
		ApiSecret:   apiSecret,
		Logger:      logger,
		UrlEncoding: base64.URLEncoding,
		Ecoding:     base64.StdEncoding,
	}
}

func (a *AuthorizerImpl) Authorize() (string, error) {
	b := bytes.NewBuffer([]byte(`grant_type=client_credentials`))
	req, err := http.NewRequest("POST", a.ApiUrl, b)
	if err != nil {
		return ``, err
	}
	req.Header.Set("Authorization", `Basic `+a.encodeApiKey())
	req.Header.Set("Content-Type", `application/x-www-form-urlencoded;charset=UTF-8`)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return ``, err
	}
	if resp.StatusCode != http.StatusOK {
		b := bytes.NewBuffer(make([]byte, 0))
		b.ReadFrom(resp.Body)
		a.Logger.Println(`Auth error`, resp.StatusCode, ERROR_MSG_AUTH_FAIL, b.String())
		return ``, errors.New(ERROR_MSG_AUTH_FAIL)
	}
	d := json.NewDecoder(resp.Body)
	dr := decoderResponse{}
	err = d.Decode(&dr)
	if err != nil {
		a.Logger.Printf(`Auth Decoding error %#v, decoded response %#v`, err, dr)
		return ``, err
	}
	tokenType := dr[`token_type`]
	if tokenType != `bearer` {
		return ``, errors.New(ERROR_MSG_INVALID_ACCESS_TOKEN_TYPE)
	}
	token, ok := dr[`access_token`].(string)
	if !ok {
		a.Logger.Printf(`Token type assertion error, decoded response %#v`, dr)
		return ``, err
	}
	if len(token) == 0 {
		a.Logger.Printf(`Token length error, decoded token %#v`, token)
		return ``, errors.New(ERROR_MSG_NO_ACCESS_TOKEN)
	}
	return token, nil
}

func (a *AuthorizerImpl) encodeApiKey() string {
	return a.Ecoding.EncodeToString([]byte(a.ApiKey + `:` + a.ApiSecret))
}
