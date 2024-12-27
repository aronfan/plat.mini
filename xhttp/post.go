package xhttp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/aronfan/plat.mini/xlog"
	"go.uber.org/zap"
)

func Post(uri string, headers map[string][]string, body []byte) (int, []byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest("POST", uri, bytes.NewBuffer(body))
	if err != nil {
		return 0, nil, err
	}

	req.Header = headers
	resp, err := client.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer resp.Body.Close()

	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, nil, err
	}

	return resp.StatusCode, respData, nil
}

func PostFormData(uri string, headers map[string][]string, values *url.Values) ([]byte, error) {
	if headers == nil {
		headers = map[string][]string{
			"Content-Type": {"application/x-www-form-urlencoded"},
		}
	}

	statusCode, respData, err := Post(uri, headers, []byte(values.Encode()))
	if err != nil {
		return nil, err
	}

	if statusCode != http.StatusOK {
		xlog.Debug("PostFormData", zap.String("resp-data", string(respData)))
		return nil, fmt.Errorf("statusCode=%d", statusCode)
	}

	return respData, nil
}

func PostJSON[T any, U any](uri string, headers map[string][]string, body *T) (*U, error) {
	if headers == nil {
		headers = map[string][]string{
			"Accept":       {"application/json"},
			"Content-Type": {"application/json"},
		}
	}

	xlog.Debug("PostJSON", zap.String("req-uri:", uri))
	xlog.Debug("PostJSON", zap.Any("req-headers:", headers))
	xlog.Debug("PostJSON", zap.Any("req-body", body))

	reqData, _ := json.Marshal(body)
	statusCode, respData, err := Post(uri, headers, reqData)
	if err != nil {
		return nil, err
	}

	if statusCode != http.StatusOK {
		xlog.Debug("PostJSON", zap.String("resp-data", string(respData)))
		return nil, fmt.Errorf("statusCode=%d", statusCode)
	}

	respBody := new(U)
	if err := json.Unmarshal([]byte(respData), respBody); err != nil {
		return nil, err
	}

	xlog.Debug("PostJSON", zap.Any("resp-body", respBody))
	return respBody, nil
}
