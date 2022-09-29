package cmdbutils

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"
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

	fmt.Println(strconv.FormatInt(requestTime, 10))
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
