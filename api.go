package mpesa

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
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

	// SANDBOX is the sandbox env tag
	SANDBOX = iota
	// PRODUCTION is the production env tag
	PRODUCTION
)

// Service is an Mpesa Service
type Service struct {
	AppKey    string
	AppSecret string
	Env       int
}

// New return a new Mpesa Service
func New(appKey, appSecret string, env int) (Service, error) {
	return Service{appKey, appSecret, env}, nil
}

//Generate Mpesa Daraja Access Token
func (s Service) authenticate() (string, error) {
	b := []byte(s.AppKey + ":" + s.AppSecret)
	encoded := base64.StdEncoding.EncodeToString(b)

	url := s.baseURL() + "oauth/v1/generate?grant_type=client_credentials"
	req, err := http.NewRequest(http.MethodGet, url, strings.NewReader(encoded))
	if err != nil {
		return "", err
	}
	req.Header.Add("authorization", "Basic "+encoded)
	req.Header.Add("cache-control", "no-cache")

	client := &http.Client{Timeout: 10 * time.Second}
	res, err := client.Do(req)
	if res != nil {
		defer res.Body.Close()
	}
	if err != nil {
		return "", fmt.Errorf("could not send auth request: %v", err)
	}

	var authResponse authResponse
	err = json.NewDecoder(res.Body).Decode(&authResponse)
	if err != nil {
		return "", fmt.Errorf("could not decode auth response: %v", err)
	}

	accessToken := authResponse.AccessToken
	log.Println("Received access_token: ", accessToken)
	return accessToken, nil
}

// STKPushSimulation sends an STK push?
func (s Service) MPESAExpressSimulation(mpesaExpress MPESAExpress) (string, error) {
	body, err := json.Marshal(mpesaExpress)
	if err != nil {
		return "", nil
	}
	auth, err := s.authenticate()
	if err != nil {
		return "", nil
	}

	headers := make(map[string]string)
	headers["content-type"] = "application/json"
	headers["authorization"] = "Bearer " + auth
	headers["cache-control"] = "no-cache"

	url := s.baseURL() + "mpesa/stkpush/v1/processrequest"
	return s.newStringRequest(url, body, headers)
}

// STKPushTransactionStatus gets a status
func (s Service) MPESAExpressTransactionStatus(mpesaExpress MPESAExpress) (string, error) {
	body, err := json.Marshal(mpesaExpress)
	if err != nil {
		return "", nil
	}

	auth, err := s.authenticate()
	if err != nil {
		return "", nil
	}

	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"
	headers["Authorization"] = "Bearer " + auth

	url := s.baseURL() + "mpesa/stkpushquery/v1/query"
	return s.newStringRequest(url, body, headers)
}

// C2BRegisterURL requests
func (s Service) C2BRegisterURL(c2BRegisterURL C2BRegisterURL) (string, error) {
	body, err := json.Marshal(c2BRegisterURL)
	if err != nil {
		return "", err
	}

	auth, err := s.authenticate()
	if err != nil {
		return "", nil
	}

	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"
	headers["Authorization"] = "Bearer " + auth
	headers["Cache-Control"] = "no-cache"

	url := s.baseURL() + "mpesa/c2b/v1/registerurl"
	return s.newStringRequest(url, body, headers)
}

// C2BSimulation sends a new request
func (s Service) C2BSimulation(c2b C2B) (string, error) {
	body, err := json.Marshal(c2b)
	if err != nil {
		return "", err
	}

	auth, err := s.authenticate()
	if err != nil {
		return "", nil
	}

	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"
	headers["Authorization"] = "Bearer " + auth
	headers["cache-control"] = "no-cache"

	url := s.baseURL() + "mpesa/c2b/v1/simulate"
	return s.newStringRequest(url, body, headers)
}

// B2CRequest sends a new request
func (s Service) B2CRequest(b2c B2C) (string, error) {
	body, err := json.Marshal(b2c)
	if err != nil {
		return "", err
	}

	auth, err := s.authenticate()
	if err != nil {
		return "", nil
	}

	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"
	headers["Authorization"] = "Bearer " + auth
	headers["cache-control"] = "no-cache"

	url := s.baseURL() + "mpesa/b2c/v1/paymentrequest"
	return s.newStringRequest(url, body, headers)
}

// B2BRequest sends a new request
func (s Service) B2BRequest(b2b B2B) (string, error) {
	body, err := json.Marshal(b2b)
	if err != nil {
		return "", nil
	}
	auth, err := s.authenticate()
	if err != nil {
		return "", nil
	}

	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"
	headers["Authorization"] = "Bearer " + auth
	headers["cache-control"] = "no-cache"

	url := s.baseURL() + "mpesa/b2b/v1/paymentrequest"
	return s.newStringRequest(url, body, headers)
}

// Reversal requests a reversal?
func (s Service) Reversal(reversal Reversal) (string, error) {
	body, err := json.Marshal(reversal)
	if err != nil {
		return "", err
	}

	auth, err := s.authenticate()
	if err != nil {
		return "", nil
	}

	headers := make(map[string]string)
	headers["Content-Type"] = "application/json"
	headers["Authorization"] = "Bearer " + auth
	headers["cache-control"] = "no-cache"

	url := s.baseURL() + "safaricom/reversal/v1/request"
	return s.newStringRequest(url, body, headers)
}

// BalanceInquiry sends a balance inquiry
func (s Service) BalanceInquiry(balanceInquiry BalanceInquiry) (string, error) {
	auth, err := s.authenticate()
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

	url := s.baseURL() + "mpesa/accountbalance/v1/query"
	return s.newStringRequest(url, body, headers)
}

func (s Service) newStringRequest(url string, body []byte, headers map[string]string) (string, error) {
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return "", nil
	}

	for key, value := range headers {
		request.Header.Set(key, value)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	res, err := client.Do(request)
	if res != nil {
		defer res.Body.Close()
	}
	if err != nil {
		return "", err
	}

	stringBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	log.Println("Response received")
	return string(stringBody), nil
}

func (s Service) baseURL() string {
	if s.Env == PRODUCTION {
		return "https://api.safaricom.co.ke/"
	}
	return "https://sandbox.safaricom.co.ke/"
}
