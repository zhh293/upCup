package ai

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

func Run(sentence string, host string, port int) (string, error) {
	url := "http://" + host + ":" + strconv.Itoa(port) + "/"
	payload := strings.NewReader("-----011000010111000001101001\r\nContent-Disposition: form-data; name=\"sentence\"\r\n\r\n" + sentence + "\r\n-----011000010111000001101001--\r\n\r\n")
	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return "", err
	}
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("content-type", "multipart/form-data; boundary=---011000010111000001101001")
	
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	
	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("AI服务返回非200状态码: %d", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}
