package cmdbutils

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"io"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

func GenSignature(accessKey string, secretKey string, requestTime int64,
	method string, uri string, params map[string]string, data map[string]interface{}, signature *string) error {
	//计算签名 核心的加密的函数
	//:param access_key: access_key
	//:param secret_key: secret_key
	//:param request_time: 请求发起时间戳
	//:param method: 请求方法 POST / GET / PUT / DELETE
	//:param uri: 请求URI
	//:param data: 请求数据，GET和POST的数据都放在这里，组成字典，通过不同的请求方式进行组合
	//:param content_type: content_type协议 application/json
	//:return:
	urlParams := ""
	if method == "GET" || method == "DELETE" {
		var keys []string
		for k := range params {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		var kvs []string
		for _, k := range keys {
			kvs = append(kvs, fmt.Sprintf("%s%s", k, data[k]))
		}
		urlParams = strings.Join(kvs, "")
	}

	bodyContent := ""
	if method == "POST" || method == "PUT" {
		bodyByte, err := json.Marshal(data)
		if err != nil {
			return err
		}
		if err = md5Str(string(bodyByte), &bodyContent); err != nil {
			return err
		}
	}

	signInfo := []string{method, uri, urlParams, "application/json", bodyContent, strconv.FormatInt(requestTime, 10), accessKey}
	signDecode := strings.Join(signInfo, "\n")

	signEncode := hmacSha1(secretKey, signDecode)
	*signature = signEncode
	return nil
}

func md5Str(s string, dst *string) error {
	w := md5.New()
	_, err := io.WriteString(w, s)
	if err != nil {
		return err
	}
	*dst = fmt.Sprintf("%x", w.Sum(nil))
	return nil
}

func hmacSha1(key string, originStr string) string {
	mac := hmac.New(sha1.New, []byte(key))
	mac.Write([]byte(originStr))
	res := hex.EncodeToString(mac.Sum(nil))
	return res
}

func CmdbPost(uri string, data map[string]interface{}, ak string, sk string, domain string, client *resty.Client, statusCode *int, body *[]byte) error {
	method := "POST"

	now := time.Now().Unix()
	var signature string
	uriParams := make(map[string]string)
	_ = GenSignature(ak, sk, now, method, uri, uriParams, data, &signature)

	fullUri := fmt.Sprintf("%s/%s", domain, uri)
	baseUrl, _ := url.Parse(fullUri)
	baseUrl.Path = uri
	params := url.Values{}
	params.Add("accesskey", ak)
	params.Add("signature", signature)
	params.Add("expires", strconv.FormatInt(now, 10))
	baseUrl.RawQuery = params.Encode()

	client.SetHeader("Content-Type", "application/json")
	// global.CmdbHttpClient.SetHeader("Host", "openapi.easyops-only.com")
	result, err := client.R().SetBody(data).Post(baseUrl.String())
	if err != nil {
		return err
	}
	*statusCode = result.StatusCode()
	*body = result.Body()
	return nil
}

func CmdbDelete(uri string, data map[string]string, ak string, sk string, domain string, client *resty.Client, statusCode *int, body *[]byte) error {
	method := "DELETE"

	now := time.Now().Unix()
	var signature string
	//uriParams := make(map[string]string)
	_ = GenSignature(ak, sk, now, method, uri, data, nil, &signature)

	fullUri := fmt.Sprintf("%s/%s", domain, uri)
	baseUrl, _ := url.Parse(fullUri)
	baseUrl.Path = uri
	params := url.Values{}
	params.Add("accesskey", ak)
	params.Add("signature", signature)
	params.Add("expires", strconv.FormatInt(now, 10))
	for k, v := range data {
		params.Add(k, v)
	}
	baseUrl.RawQuery = params.Encode()

	client.SetHeader("Content-Type", "application/json")
	// global.CmdbHttpClient.SetHeader("Host", "openapi.easyops-only.com")
	result, err := client.R().Delete(baseUrl.String())
	if err != nil {
		return err
	}
	*statusCode = result.StatusCode()
	*body = result.Body()
	return nil
}
