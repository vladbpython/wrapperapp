package adapters

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

type HttpClient struct {
	AppName string            //Наименование клиента
	Headers map[string]string // заголовки
	Client  *http.Client      // экземпляр клиента
}

func (a *HttpClient) MakeRequest(method, url string, query map[string]string,
	headers map[string]string, body []byte, JSONRequest, JSONResponse bool) ([]byte, int, error) {

	req, _ := http.NewRequest(method, url, bytes.NewReader(body))
	for k, v := range a.Headers {
		req.Header.Set(k, v)
	}
	q := req.URL.Query()
	for k, v := range query {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()

	for k, v := range headers {
		req.Header.Set(k, v)
	}
	if JSONRequest {
		req.Header.Set("content-type", "application/json")
	}

	if JSONResponse {
		req.Header.Set("accept", "application/json")
	}
	resp, err := a.Client.Do(req)
	if err != nil {

		return []byte(""), 503, err
	}
	defer resp.Body.Close()
	data, _ := ioutil.ReadAll(resp.Body)
	return data, resp.StatusCode, err

}

func (a *HttpClient) MakeRequestWithModel(Method, Url string, query map[string]string,
	headers map[string]string, body []byte, model interface{}, JSONRequest, JSONResponse bool) error {

	data, _, err := a.MakeRequest(Method, Url, query, headers, body, JSONRequest, JSONResponse)
	json.Unmarshal(data, model)
	return err
}

func NewHttpClient(appName string, timeOut time.Duration, headers map[string]string) *HttpClient {
	return &HttpClient{
		AppName: appName,
		Client: &http.Client{
			Timeout: time.Second * timeOut,
		},
		Headers: headers,
	}
}
