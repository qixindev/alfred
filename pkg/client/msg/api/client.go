package api

import (
	"alfred/pkg/utils"
	"crypto/tls"
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

// 免证书校验
func certificateValidation() *http.Client {
	tr := &http.Transport{
		//InsecureSkipVerify用来控制客户端是否证书和服务器主机名。如果设置为true, 则不会校验证书以及证书中的主机名和服务器主机名是否一致。
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	return &http.Client{Transport: tr}
}

// PostClient 发送POST请求
func PostClient(url string, postBody any, result any) error {
	sendBody := ""
	if postBody != nil {
		marshal, err := json.Marshal(postBody)
		if err != nil {
			return err
		}
		sendBody = string(marshal)
	}

	resp, err := certificateValidation().Post(url, "application/json;charset=UTF-8", strings.NewReader(sendBody))
	if err != nil {
		return err
	}

	defer utils.DeferErr(resp.Body.Close)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if result != nil {
		if err = json.Unmarshal(body, result); err != nil {
			return err
		}
	}

	return nil
}

// GetClient 发送GET请求
func GetClient(url string, result any) error {
	resp, err := certificateValidation().Get(url)
	if err != nil {
		return err
	}

	defer utils.DeferErr(resp.Body.Close)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if result != nil {
		if err = json.Unmarshal(body, result); err != nil {
			return err
		}
	}

	return nil
}
