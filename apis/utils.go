package apis

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"golang.org/x/net/proxy"
)

var (
	httpTransport    = &http.Transport{}
	torHttpClient    = &http.Client{Transport: httpTransport}
	directHttpClient = &http.Client{}
	proxyAddr        = "127.0.0.1:9050"
)

func init() {
	dialer, err := proxy.SOCKS5("tcp", proxyAddr, nil, proxy.Direct)
	if err != nil {
		panic(err)
	}
	httpTransport.Dial = dialer.Dial
}

func TorGET(url string) (string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	resp, err := torHttpClient.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", errors.New("Wrong status code")
	}

	bm, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(bm), nil
}

func DirectGET(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", errors.New("Wrong status code")
	}

	bm, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(bm), nil
}

func TorPOST(url string, params url.Values) (string, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBufferString(params.Encode()))
	if err != nil {
		return "", err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(params.Encode())))
	resp, err := torHttpClient.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	bm, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(bm), nil
}

func DirectPOST(url string, params url.Values) (string, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBufferString(params.Encode()))
	if err != nil {
		return "", err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(params.Encode())))
	resp, err := directHttpClient.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	bm, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(bm), nil
}

func TorRawPOST(url string, params string) (string, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBufferString(params))
	if err != nil {
		return "", err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(params)))
	resp, err := torHttpClient.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	bm, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(bm), nil
}
