package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/spf13/viper"
)

type ExpireTime struct {
	time.Time
}

const webApiExpireTimeFormatwithT = time.RFC3339Nano                    // "2025-10-29T17:49:25.709Z"
const webApiExpireTimeFormatWithSpace = "2006-01-02 15:04:05.000 -0700" // "2025-10-29 17:59:45.090 +0000"

func (at *ExpireTime) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	if s == "null" {
		return nil
	}

	t, err := time.Parse(webApiExpireTimeFormatwithT, s)
	if err == nil {
		at.Time = t
		return nil
	}

	t, err = time.Parse(webApiExpireTimeFormatWithSpace, s)
	if err == nil {
		at.Time = t
		return nil
	}

	return fmt.Errorf("could not parse time in either format: %s", s)
}

type TokenManager struct {
	AccessToken string        `json:"accessToken"`
	Expiration  ExpireTime    `json:"expiredDate"`
	mux         sync.RWMutex  `json:"-"`
	buffer      time.Duration `json:"-"`
}

func NewTokenManager() *TokenManager {
	m := &TokenManager{
		buffer: 5 * time.Minute,
	}

	if err := m.loadToken(); err != nil {
		fmt.Println("No cached token found, will fetch a new one.")
	}

	return m
}

func (m *TokenManager) saveToken() error {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return err
	}

	cliConfigDir := filepath.Join(configDir, "griddb-cloud-cli")
	tokenPath := filepath.Join(cliConfigDir, "token.json")

	data, err := json.Marshal(m)
	if err != nil {
		return err
	}

	return os.WriteFile(tokenPath, data, 0600)
}

func (m *TokenManager) loadToken() error {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return err
	}

	tokenPath := filepath.Join(configDir, "griddb-cloud-cli", "token.json")

	data, err := os.ReadFile(tokenPath)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, &m)
}

func (m *TokenManager) getBearerToken() error {
	fmt.Println("--- REFRESHING TOKEN ---")

	if !(viper.IsSet("cloud_url")) {
		log.Fatal("Please provide a `cloud_url` in your config file! You can copy this directly from your Cloud dashboard")
	}

	configCloudURL := viper.GetString("cloud_url")
	parsedURL, err := url.Parse(configCloudURL)
	if err != nil {
		fmt.Println("Error parsing URL:", err)
		return err
	}
	newPath := path.Dir(parsedURL.Path)
	parsedURL.Path = newPath
	authEndpoint, _ := parsedURL.Parse("./authenticate")

	authURL := authEndpoint.String()
	method := "POST"

	user := viper.GetString("cloud_username")
	pass := viper.GetString("cloud_pass")
	payloadStr := fmt.Sprintf(`{"username": "%s", "password": "%s" }`, user, pass)
	payload := strings.NewReader(payloadStr)

	client := &http.Client{}
	req, err := http.NewRequest(method, authURL, payload)
	if err != nil {
		log.Fatal(err)
	}
	defer req.Body.Close()

	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("error with client DO: ", err)
	}
	CheckForErrors(resp)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(body, &m); err != nil {
		log.Fatalf("Error unmarshaling access token %s", err)
	}

	if err := m.saveToken(); err != nil {
		fmt.Println("Warning: Could not save token to cache:", err)
	}

	return nil
}

func (m *TokenManager) getAndAddValidToken(req *http.Request) error {
	m.mux.RLock()
	needsRefresh := time.Now().UTC().After(m.Expiration.Add(-m.buffer))
	m.mux.RUnlock()

	if needsRefresh {
		m.mux.Lock()
		if time.Now().UTC().After(m.Expiration.Add(-m.buffer)) {
			if err := m.getBearerToken(); err != nil {
				m.mux.Unlock()
				return err
			}
		}
		m.mux.Unlock()
	}

	m.mux.RLock()
	defer m.mux.RUnlock()

	req.Header.Add("Authorization", "Bearer "+m.AccessToken)
	req.Header.Add("Content-Type", "application/json")
	return nil
}

func MakeNewRequest(method, endpoint string, body io.Reader) (req *http.Request, e error) {

	if !(viper.IsSet("cloud_url")) {
		log.Fatal("Please provide a `cloud_url` in your config file! You can copy this directly from your Cloud dashboard")
	}

	url := viper.GetString("cloud_url")
	req, err := http.NewRequest(method, url+endpoint, body)
	if err != nil {
		fmt.Println("error with request:", err)
		return req, err
	}
	tokenManager.getAndAddValidToken(req)
	return req, nil
}
