package xhttp

import (
	"io"
	"net/http"
	"net/url"

	"github.com/aronfan/plat.mini/xlog"
	"go.uber.org/zap"
)

func Get(uri string, headers map[string][]string, params map[string]string) (int, []byte, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return 0, nil, err
	}

	q := u.Query()
	for k, v := range params {
		q.Add(k, v)
	}
	u.RawQuery = q.Encode()

	xlog.Debug("Get", zap.String("uri", u.String()))

	client := &http.Client{}
	req, err := http.NewRequest("GET", u.String(), nil)
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
