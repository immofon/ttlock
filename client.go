package ttlock

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

const (
	CNBaseURL = "https://cnapi.ttlock.com" // 中国大陆服务器地址
	EUBaseURL = "https://euapi.ttlock.com" // 欧洲服务器地址
)

// Client represents a TTLock API client
type Client struct {
	ClientID     string
	ClientSecret string
	Username     string
	Password     string // Plain text password; will be MD5 hashed automatically
	BaseURL      string
	HTTPClient   *http.Client

	access_token_resp      AccessTokenResponse
	access_token_resp_lock sync.RWMutex
}

// NewClient creates a new TTLock API client, defaulting to the China base URL
func NewClient(clientID, clientSecret string, username, password string) *Client {
	c := &Client{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Username:     username,
		Password:     password,
		BaseURL:      CNBaseURL,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}

	// Initialize access token
	tokenResp, err := c.GetAccessToken()
	if err != nil {
		panic(fmt.Sprintf("failed to get access token: %v", err))
	}
	c.update_access_token_resp(*tokenResp)

	// Start background goroutine to refresh access token
	go c.refresh_access_token()

	return c
}

// update_access_token_resp safely updates the access token response
func (c *Client) update_access_token_resp(resp AccessTokenResponse) {
	c.access_token_resp_lock.Lock()
	defer c.access_token_resp_lock.Unlock()
	c.access_token_resp = resp
}

// AccessToken retrieves the current access token
func (c *Client) AccessToken() string {
	c.access_token_resp_lock.RLock()
	defer c.access_token_resp_lock.RUnlock()
	return c.access_token_resp.AccessToken
}

// this function should be runed by goroutine in background by caller.
// It should contain a dead loop to periodically check and refresh the access token if needed.
// It refreshes the token depends on .ExpiresIn field by min(.ExpiresIn/2, 1 day)
func (c *Client) refresh_access_token() {
	for {
		// Calculate the sleep duration based on ExpiresIn
		c.access_token_resp_lock.RLock()
		expiresIn := c.access_token_resp.ExpiresIn
		refreshToken := c.access_token_resp.RefreshToken
		c.access_token_resp_lock.RUnlock()

		// Default to 1 day if ExpiresIn is too large
		refreshInterval := time.Duration(expiresIn/2) * time.Second
		if refreshInterval > 24*time.Hour {
			refreshInterval = 24 * time.Hour
		}

		// Sleep for the calculated duration
		time.Sleep(refreshInterval)

		// Refresh the access token with retry logic
		var err error
		for i := 0; i < 3; i++ {
			newTokenResp, retryErr := c.RefreshAccessToken(refreshToken)
			if retryErr != nil {
				err = retryErr
				fmt.Printf("failed to refresh access token (attempt %d): %v\n", i+1, retryErr)
				time.Sleep(2 * time.Second) // Wait before retrying
				continue
			}

			// Update the access token response on success
			c.update_access_token_resp(*newTokenResp)
			err = nil
			break
		}

		// Panic if all retries fail
		if err != nil {
			panic(fmt.Sprintf("failed to refresh access token after 3 attempts: %v", err))
		}
	}
}

// SetBaseURL sets the API base URL
func (c *Client) SetBaseURL(url string) {
	c.BaseURL = url
}

// AccessTokenResponse represents the response from the oauth2/token endpoint
type AccessTokenResponse struct {
	AccessToken  string `json:"access_token"`
	UID          int    `json:"uid"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"` // in seconds
	Errcode      int    `json:"errcode"`
	Errmsg       string `json:"errmsg"`
}

// GetAccessToken obtains an access token using the user's credentials.
// The password should be the plain text password; it will be MD5 hashed automatically.
func (c *Client) GetAccessToken() (*AccessTokenResponse, error) {
	endpoint := c.BaseURL + "/oauth2/token"

	// MD5 hash the password
	hasher := md5.New()
	hasher.Write([]byte(c.Password))
	md5Password := hex.EncodeToString(hasher.Sum(nil))

	data := url.Values{}
	data.Set("clientId", c.ClientID)
	data.Set("clientSecret", c.ClientSecret)
	data.Set("username", c.Username)
	data.Set("password", md5Password)

	req, err := http.NewRequest("POST", endpoint, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var tokenResp AccessTokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if tokenResp.Errcode != 0 {
		return nil, NewError(ErrorCode(tokenResp.Errcode))
	}

	return &tokenResp, nil
}

// RefreshAccessToken refreshes the access token using a refresh token.
func (c *Client) RefreshAccessToken(refreshToken string) (*AccessTokenResponse, error) {
	endpoint := c.BaseURL + "/oauth2/token"

	data := url.Values{}
	data.Set("clientId", c.ClientID)
	data.Set("clientSecret", c.ClientSecret)
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", refreshToken)

	req, err := http.NewRequest("POST", endpoint, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var tokenResp AccessTokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if tokenResp.Errcode != 0 {
		return nil, NewError(ErrorCode(tokenResp.Errcode))
	}

	return &tokenResp, nil
}
