package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/NeowayLabs/logger"
)

var logKC = logger.Namespace("kripto.bot")
var insecureTls = os.Getenv("KRIPTO_INSECURE_TLS")
var skipVerify bool = false

var (
	Username string
	Password string
)

type (
	Credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	Result struct {
		Token string `json:"token"`
	}

	Secret struct {
		App  string            `json:"app"`
		Vars map[string]string `json:"vars"`
	}
)

func init() {

	if "" != insecureTls {
		skipVerify = true
	}
}

func main() {

	var app = os.Getenv("KRIPTO_APP")
	var endpoint = os.Getenv("KRIPTO_SERVER_ENDPOINT")

	result, err := authenticate(Username, Password, buildAuthUrl(endpoint))
	if err != nil {
		logKC.Fatal("Unauthorized! %v", err)
	}

	sec, err := getVars(result.Token, buildSecretsUrl(endpoint, app))
	if err != nil {
		logKC.Fatal("Bad token! %v", err)
	}

	err = setVars(sec)
	if err != nil {
		logKC.Fatal("Output error! %v", err)
	}

}

func authenticate(u, p, e string) (*Result, error) {

	c := &Credentials{
		Username: Username,
		Password: Password,
	}

	q, err := json.Marshal(c)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, e, bytes.NewBuffer(q))
	if err != nil {
		return nil, err
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: skipVerify},
	}

	client := &http.Client{Transport: tr}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var result Result
	err = json.NewDecoder(resp.Body).Decode(&result)

	if err != nil {
		return nil, err
	}

	return &result, nil

}

func getVars(t, url string) (*Secret, error) {

	req, err := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Add("Authorization", t)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: skipVerify},
	}

	client := &http.Client{Transport: tr}
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var sec Secret
	err = json.NewDecoder(resp.Body).Decode(&sec)

	if err != nil {
		return nil, err
	}

	return &sec, nil
}

func setVars(s *Secret) error {

	filename := fmt.Sprintf("%s.env", s.App)
	_ = os.Remove(filename)

	var text string = ""
	for k, v := range s.Vars {
		text = fmt.Sprintf("export %s=%s\n", k, v)
		err := out(filename, text)
		if err != nil {
			return err
		}
	}

	return nil
}

func out(filename, text string) error {

	f, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}

	defer f.Close()

	if _, err = f.WriteString(text); err != nil {
		return err
	}

	return nil
}

func buildAuthUrl(e string) string {
	return fmt.Sprintf("%s/authenticate", e)
}

func buildSecretsUrl(e, a string) string {
	return fmt.Sprintf("%s/secrets?app=%s", e, a)
}
