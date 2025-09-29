package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

const LoginPath = "/em-edm/sessions"

type LoginResponse struct {
	AuthToken struct {
		TokenType string `json:"token_type"`
		Token     string `json:"access_token"`
		Role      string `json:"role"`
	} `json:"AccessToken"`
	Error struct {
		Code int `json:"ErrorCode"`
	} `json:"Error"`
}
type PanelAPI struct {
	comma rune
}

func GetCSVRequest(url, token string) *http.Response {
	get, err := http.NewRequest("GET", fmt.Sprintf("%s/csv", url), nil)
	if err != nil {
		return nil
	}
	get.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	get.Header.Set("Accept", "text/csv, /")
	get.Header.Add("Connection", "keep-alive")
	get.Header.Add("Referer", fmt.Sprintf("%s/public/settings/equipment-management/local-export", url))
	get.Header.Add("Sec-Fetch-Dest", "empty")
	get.Header.Add("Sec-Fetch-Mode", "cors")
	get.Header.Add("Sec-Fetch-Site", "same-origin")
	get.Header.Add("Sec-GPC", "1")
	client := &http.Client{}
	res, err := client.Do(get)
	if err != nil {
		return nil
	}
	return res
}

func (p PanelAPI) fetchCSV(url string, user string, pass string) (io.Reader, error) {
	//Login
	var loginRes LoginResponse
	loginPayload := strings.NewReader(fmt.Sprintf(`{"scheme": "BASIC", "user":"%s","password":"%s"}`, user, pass))
	resp, err := http.Post(fmt.Sprintf("%s%s", url, LoginPath), "application/json", loginPayload)
	if err != nil {
		return nil, err
	}
	respJson, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(respJson, &loginRes)
	if err != nil || loginRes.Error.Code != 0 {
		return nil, fmt.Errorf("Login failed")
	}

	//CSV download
	token := loginRes.AuthToken.Token
	csvResponse := GetCSVRequest(url, token)
	if csvResponse == nil {
		return nil, fmt.Errorf("Error getting CSV")
	}
	defer func() {
		_ = csvResponse.Body.Close()
	}()
	if csvResponse.StatusCode < 200 || csvResponse.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return csvResponse.Body, nil
}

func (p PanelAPI) saveCSV(data io.Reader, fileName string) error {
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = io.Copy(file, data)
	if err != nil {
		return fmt.Errorf("error saving CSV data: %w", err)
	}
	return nil
}
