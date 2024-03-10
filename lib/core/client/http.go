package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"gopkg.in/resty.v1"
)

type Info struct {
	Method  string
	Url     string
	Auth    string
	Payload map[string]interface{}
	Headers map[string]interface{}
}

func (i *Info) prepare() *resty.Request {
	restyInstance := resty.New()
	restyInstance.SetTimeout(time.Duration(1 * time.Minute)).
		SetDebug(false)
	// logFile, _ := os.OpenFile("storages/logs/go-resty.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	// restyInstance.SetLogger(logFile)

	req := restyInstance.R()

	if len(i.Headers) > 0 {
		for k, v := range i.Headers {
			req.SetHeader(k, v.(string))
		}
	} else {
		req.SetHeader("Content-Type", "application/json").
			SetHeader("Accept", "application/json")

		if i.Auth != "" {
			//req.SetAuthToken(i.Auth)
			req.SetHeader("Authorization", i.Auth)
		}
	}

	return req
}

func (i *Info) Dispatch(target interface{}) (error, string) {
	var (
		err  error
		resp *resty.Response
	)

	req := i.prepare()
	if i.Method == "" {
		return errors.New("Missing HTTP method for dispatch request to core."), ""
	}

	if len(i.Payload) == 0 || i.Payload == nil {
		resp, err = req.Execute(i.Method, i.Url)
	} else {
		resp, err = req.SetBody(i.Payload).Execute(i.Method, i.Url)
	}

	if err != nil {
		return err, ""
	}

	fmt.Println("response >>> ", resp.String())

	if resp.String() != "" {
		err = json.Unmarshal([]byte(resp.String()), target)
		if err != nil {
			return err, resp.String()
		}
	}

	return nil, resp.String()
}

func (i *Info) DispatchWithOriginalReturn() (*resty.Response, error) {
	var (
		err  error
		resp *resty.Response
	)

	req := i.prepare()
	if i.Method == "" {
		return nil, errors.New("Missing HTTP method for dispatch request to core.")
	}

	if len(i.Payload) == 0 || i.Payload == nil {
		resp, err = req.Execute(i.Method, i.Url)
	} else {
		resp, err = req.SetBody(i.Payload).Execute(i.Method, i.Url)
	}

	if err != nil {
		return nil, err
	}

	return resp, err
}
