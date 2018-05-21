package mpesa

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

// Env is the environment type
type Env string

const (
	// DEV is the development env tag
	SANDBOX = iota
	// PRODUCTION is the production env tag
	PRODUCTION
)

// Mpesa is mpesa
type Mpesa struct {
	ConsumerKey    string
	ConsumerSecret string
	Env            int
}

// New return a new Mpesa
func New(appKey, appSecret string, env int) (Mpesa, error) {
	return Mpesa{appKey, appSecret, env}, nil
}

//Generate Daraja Access Token
func (m Mpesa) authenticate() (string, error) {
	b := []byte(m.ConsumerKey + ":" + m.ConsumerSecret)
	encoded := base64.StdEncoding.EncodeToString(b)

	url := m.baseURL() + "oauth/v1/generate?grant_type=client_credentials"
	request, err := http.NewRequest(http.MethodGet, url, strings.NewReader(encoded))
	if err != nil {
		return "", err
	}
	request.Header.Add("authorization", "Basic "+encoded)
	request.Header.Add("cache-control", "no-cache")

	client := &http.Client{Timeout: 10 * time.Second}
	response, err := client.Do(request)
	if response != nil {
		defer response.Body.Close()
	}
	if err != nil {
		return "", err
	}

	var authResponse authResponse
	json.NewDecoder(response.Body).Decode(&authResponse)

	accessToken := authResponse.AccessToken
	log.Println("Received access_token: ", accessToken)
	return accessToken, nil
}

// STKPushSimulation sends an STK push?
func (m Mpesa) STKPushSimulation(stkPush STKPush) (string, error) {
	body, err := json.Marshal(stkPush)
	if err != nil {
		return "", nil
	}
	auth, err := m.authenticate()
	if err != nil {
		return "", nil
	}

	headers := make(map[string]string)
	headers["content-type"] = "application/json"
	headers["authorization"] = "Bearer " + auth
	headers["cache-control"] = "no-cache"

	url := m.baseURL() + "mpesa/stkpush/v1/processrequest"
	return m.newStringRequest(url, body, headers)
}

// STKPushTransactionStatus gets a status
func (m Mpesa) STKPushTransactionStatus(stkPush STKPush) (string, error) {
	body, err := json.Marshal(stkPush)
	if err != nil {
		return "", nil
	}

	auth, err := m.authenticate()
	if err != nil {
		return "", nil
	}

	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"
	headers["Authorization"] = "Bearer " + auth

	url := m.baseURL() + "mpesa/stkpushquery/v1/query"
	return m.newStringRequest(url, body, headers)
}

// RegisterURL requests
func (m Mpesa) C2BRegisterURL(c2BRegisterURL C2BRegisterURL) (string, error) {
	body, err := json.Marshal(c2BRegisterURL)
	if err != nil {
		return "", err
	}

	auth, err := m.authenticate()
	if err != nil {
		return "", nil
	}

	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"
	headers["Authorization"] = "Bearer " + auth
	headers["Cache-Control"] = "no-cache"

	url := m.baseURL() + "mpesa/c2b/v1/registerurl"
	return m.newStringRequest(url, body, headers)
}

// C2BSimulation sends a new request
func (m Mpesa) C2BSimulation(c2b C2B) (string, error) {
	body, err := json.Marshal(c2b)
	if err != nil {
		return "", err
	}

	auth, err := m.authenticate()
	if err != nil {
		return "", nil
	}

	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"
	headers["Authorization"] = "Bearer " + auth
	headers["cache-control"] = "no-cache"

	url := m.baseURL() + "mpesa/c2b/v1/simulate"
	return m.newStringRequest(url, body, headers)
}

// B2CRequest sends a new request
func (m Mpesa) B2CRequest(b2c B2C) (string, error) {
	body, err := json.Marshal(b2c)
	if err != nil {
		return "", err
	}

	auth, err := m.authenticate()
	if err != nil {
		return "", nil
	}

	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"
	headers["Authorization"] = "Bearer " + auth
	headers["cache-control"] = "no-cache"

	url := m.baseURL() + "mpesa/b2c/v1/paymentrequest"
	return m.newStringRequest(url, body, headers)
}

// B2BRequest sends a new request
func (m Mpesa) B2BRequest(b2b B2B) (string, error) {
	body, err := json.Marshal(b2b)
	if err != nil {
		return "", nil
	}
	auth, err := m.authenticate()
	if err != nil {
		return "", nil
	}

	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"
	headers["Authorization"] = "Bearer " + auth
	headers["cache-control"] = "no-cache"

	url := m.baseURL() + "mpesa/b2b/v1/paymentrequest"
	return m.newStringRequest(url, body, headers)
}

// Reversal requests a reversal?
func (m Mpesa) Reversal(reversal Reversal) (string, error) {
	body, err := json.Marshal(reversal)
	if err != nil {
		return "", err
	}

	auth, err := m.authenticate()
	if err != nil {
		return "", nil
	}

	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"
	headers["Authorization"] = "Bearer " + auth
	headers["cache-control"] = "no-cache"

	url := m.baseURL() + "safaricom/reversal/v1/request"
	return m.newStringRequest(url, body, headers)
}

// BalanceInquiry sends a balance inquiry
func (m Mpesa) BalanceInquiry(balanceInquiry BalanceInquiry) (string, error) {
	auth, err := m.authenticate()
	if err != nil {
		return "", nil
	}

	body, err := json.Marshal(balanceInquiry)
	if err != nil {
		return "", err
	}

	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"
	headers["Authorization"] = "Bearer " + auth
	headers["cache-control"] = "no-cache"
	headers["postman-token"] = "2aa448be-7d56-a796-065f-b378ede8b136"

	url := m.baseURL() + "safaricom/accountbalance/v1/query"
	return m.newStringRequest(url, body, headers)
}

func (m Mpesa) newStringRequest(url string, body []byte, headers map[string]string) (string, error) {
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return "", nil
	}

	for key, value := range headers {
		request.Header.Set(key, value)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	response, err := client.Do(request)
	if response != nil {
		defer response.Body.Close()
	}
	if err != nil {
		return "", err
	}

	respBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	log.Println("Response received")
	return string(respBody), nil
}

func (m Mpesa) baseURL() string {
	if m.Env == PRODUCTION {
		return "https://api.safaricom.co.ke/"
	}
	return "https://sandbox.safaricom.co.ke/"
}
