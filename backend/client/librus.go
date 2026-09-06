package client

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	ApiURL   = "https://synergia.librus.pl/gateway/api/2.0"
	AuthInit = "https://synergia.librus.pl/loguj/portalRodzina?v="
	AuthURL  = "https://api.librus.pl/OAuth/Authorization?client_id=46"
	MFAURL   = "https://api.librus.pl/OAuth/Authorization/PerformLogin?client_id=46"
	GrantURL = "https://api.librus.pl/OAuth/Authorization/Grant?client_id=46"
)

type LibrusClient struct {
	HTTPClient *http.Client
}

func NewClient() *LibrusClient {
	jar, _ := cookiejar.New(nil)
	return &LibrusClient{
		HTTPClient: &http.Client{
			Jar:     jar,
			Timeout: 10 * time.Second,
		},
	}
}

func (c *LibrusClient) getBaner() string {
	timestamp := strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10)
	turnips := c.transformString(timestamp)

	randomStr := strconv.FormatFloat(rand.Float64(), 'f', -1, 64)
	preTurnips := c.transformString(randomStr)

	return preTurnips + "_" + turnips
}

func (c *LibrusClient) transformString(input string) string {
	var result strings.Builder
	for _, char := range input {
		result.WriteRune(char + 20)
	}
	return result.String()
}

func (c *LibrusClient) Authenticate(username, password string) error {
	// Ensure HTTPClient and Jar are initialized if struct was instantiated manually
	if c.HTTPClient == nil {
		jar, _ := cookiejar.New(nil)
		c.HTTPClient = &http.Client{
			Jar:     jar,
			Timeout: 10 * time.Second,
		}
	} else if c.HTTPClient.Jar == nil {
		jar, _ := cookiejar.New(nil)
		c.HTTPClient.Jar = jar
	}

	// 1. Initial warm-up request
	initReqURL := fmt.Sprintf("%s%d", AuthInit, time.Now().Unix())
	req, err := http.NewRequest("GET", initReqURL, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Referer", "https://portal.librus.pl/")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	resp.Body.Close()

	// Handle redirect if present
	if resp.StatusCode == http.StatusFound || resp.StatusCode == http.StatusSeeOther {
		loc := resp.Header.Get("Location")
		if loc != "" {
			reqRedir, _ := http.NewRequest("GET", loc, nil)
			reqRedir.Header.Set("Referer", "https://portal.librus.pl/")
			if respRedir, err := c.HTTPClient.Do(reqRedir); err == nil {
				respRedir.Body.Close()
			}
		}
	}

	// 2. Touch OAuth page
	reqAuth, _ := http.NewRequest("GET", AuthURL, nil)
	if respAuth, err := c.HTTPClient.Do(reqAuth); err == nil {
		respAuth.Body.Close()
	}

	// 3. Post credentials using JavaScript-compatible payload & X-Baner header
	formData := url.Values{}
	formData.Set("action", "login")
	formData.Set("login", username)
	formData.Set("pass", password)

	postReq, err := http.NewRequest("POST", AuthURL, strings.NewReader(formData.Encode()))
	if err != nil {
		return err
	}

	postReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	postReq.Header.Set("Accept", "application/json")
	postReq.Header.Set("X-Requested-With", "XMLHttpRequest")
	postReq.Header.Set("X-Baner", c.getBaner())
	postReq.Header.Set("Referer", AuthURL)

	postResp, err := c.HTTPClient.Do(postReq)
	if err != nil {
		return err
	}
	defer postResp.Body.Close()

	if postResp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(postResp.Body)
		return fmt.Errorf("authentication failed [status %d]: %s", postResp.StatusCode, string(body))
	}

	// 4. Complete MFA & Grant sequence
	for _, targetURL := range []string{MFAURL, MFAURL, GrantURL} {
		gReq, _ := http.NewRequest("GET", targetURL, nil)
		gResp, err := c.HTTPClient.Do(gReq)
		if err != nil {
			return err
		}
		gResp.Body.Close()
		if gResp.StatusCode != http.StatusOK {
			return fmt.Errorf("failed during handshake stage: %s", targetURL)
		}
	}

	// 5. Activate session
	return c.Activate()
}

func (c *LibrusClient) Activate() error {
	tokenReq, _ := http.NewRequest("GET", ApiURL+"/Auth/TokenInfo", nil)
	tokenResp, err := c.HTTPClient.Do(tokenReq)
	if err != nil {
		return err
	}
	defer tokenResp.Body.Close()

	if tokenResp.StatusCode != http.StatusOK {
		return fmt.Errorf("token info request failed with status: %d", tokenResp.StatusCode)
	}

	var tokenInfo struct {
		UserIdentifier string `json:"UserIdentifier"`
	}
	if err := json.NewDecoder(tokenResp.Body).Decode(&tokenInfo); err != nil {
		return err
	}

	userReq, _ := http.NewRequest("GET", fmt.Sprintf("%s/Auth/UserInfo/%s", ApiURL, tokenInfo.UserIdentifier), nil)
	userResp, err := c.HTTPClient.Do(userReq)
	if err != nil {
		return err
	}
	defer userResp.Body.Close()

	if userResp.StatusCode != http.StatusOK {
		return fmt.Errorf("user info request failed")
	}

	// Session ping
	msgReq, _ := http.NewRequest("GET", "https://synergia.librus.pl/wiadomosci2", nil)
	if msgResp, err := c.HTTPClient.Do(msgReq); err == nil {
		msgResp.Body.Close()
	}

	return nil
}
