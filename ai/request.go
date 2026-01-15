package ai

import (
	"io"
	"net/http"
	"strconv"
	"strings"
)

func Run(sentence string, port int) string {
	url := "http://127.0.0.1:" + strconv.Itoa(port) + "/"
	payload := strings.NewReader("-----011000010111000001101001\r\nContent-Disposition: form-data; name=\"sentence\"\r\n\r\n" + sentence + "\r\n-----011000010111000001101001--\r\n\r\n")
	req, _ := http.NewRequest("POST", url, payload)
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Accept-Encoding", "gzip, deflate, br")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("content-type", "multipart/form-data; boundary=---011000010111000001101001")
	res, e := http.DefaultClient.Do(req)
	if e != nil {
		return "err"
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	return string(body)
}
