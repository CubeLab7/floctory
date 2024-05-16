package softlinePayment

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type Service struct {
	config *Config
}

const (
	ping       = "/v2/ping"
	leads      = "/v2/exchange/leads"
	phoneLeads = "/v2/exchange/phone-leads"
)

func New(config *Config) *Service {
	return &Service{
		config: config,
	}
}

func sendRequest(config *Config, inputs *SendParams) (respBody []byte, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("Floctory! SendRequest: %v", err)
		}
	}()

	baseURL, err := url.Parse(config.URI)
	if err != nil {
		return respBody, fmt.Errorf("can't parse URI from config: %w", err)
	}

	// Добавляем путь из inputs.Path к базовому URL
	baseURL.Path += inputs.Path

	// Устанавливаем параметры запроса из queryParams
	query := baseURL.Query()
	for key, value := range inputs.QueryParams {
		query.Set(key, value)
	}

	baseURL.RawQuery = query.Encode()

	finalUrl := baseURL.String()

	req, err := http.NewRequest(inputs.HttpMethod, finalUrl, inputs.Body)
	if err != nil {
		return respBody, fmt.Errorf("can't create request! Err: %s", err)
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Accept", "application/json")

	httpClient := http.Client{
		Transport: &http.Transport{
			IdleConnTimeout: time.Second * time.Duration(config.IdleConnTimeoutSec),
		},
		Timeout: time.Second * time.Duration(config.RequestTimeoutSec),
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return respBody, fmt.Errorf("can't do request! Err: %s", err)
	}
	defer resp.Body.Close()

	inputs.HttpCode = resp.StatusCode

	respBody, err = io.ReadAll(resp.Body)
	if err != nil {
		return respBody, fmt.Errorf("can't read response body! Err: %w", err)
	}

	if resp.StatusCode == http.StatusInternalServerError {
		return respBody, fmt.Errorf("error: %v", string(respBody))
	}

	if err = json.Unmarshal(respBody, &inputs.Response); err != nil {
		return respBody, fmt.Errorf("can't unmarshall response: '%v'. Err: %w", string(respBody), err)
	}

	return
}

func (s *Service) Ping() (bool, error) {
	params := map[string]string{
		"pong":  "",
		"token": s.config.Token,
	}

	inputs := SendParams{
		Path:        ping,
		HttpMethod:  http.MethodGet,
		QueryParams: params,
	}

	if _, err := sendRequest(s.config, &inputs); err != nil {
		return false, err
	}

	if inputs.HttpCode != http.StatusOK {
		return false, nil
	}

	return true, nil
}

func (s *Service) ExchangeLeads(data Request) (response *ExchangeLeadsResponse, err error) {
	response = new(ExchangeLeadsResponse)

	params := map[string]string{
		"site_id":  fmt.Sprint(s.config.SiteID),
		"page":     fmt.Sprint(data.Page),
		"per_page": fmt.Sprint(data.PerPage),
		"from":     fmt.Sprint(data.From),
		"to":       fmt.Sprint(data.To),
		"token":    s.config.Token,
	}

	inputs := SendParams{
		Path:        leads,
		HttpMethod:  http.MethodGet,
		Response:    response,
		QueryParams: params,
	}

	if _, err = sendRequest(s.config, &inputs); err != nil {
		return
	}

	var hasNextData bool

	if len(response.Data) > 0 {
		hasNextData, err = s.CheckNext(data)
		if err != nil {
			return nil, errors.New("CheckNext! Desc: " + err.Error())
		}
	}

	response.HasNextData = hasNextData

	return
}

func (s *Service) ExchangePhoneLeads(data Request) (response *PhoneLeadsResponse, err error) {
	response = new(PhoneLeadsResponse)

	params := map[string]string{
		"site_id":  fmt.Sprint(s.config.SiteID),
		"page":     fmt.Sprint(data.Page),
		"per_page": fmt.Sprint(data.PerPage),
		"from":     fmt.Sprint(data.From),
		"to":       fmt.Sprint(data.To),
		"token":    s.config.Token,
	}

	inputs := SendParams{
		Path:        phoneLeads,
		HttpMethod:  http.MethodGet,
		Response:    response,
		QueryParams: params,
	}

	if _, err = sendRequest(s.config, &inputs); err != nil {
		return
	}

	var hasNextData bool

	if len(response.Data) > 0 {
		hasNextData, err = s.CheckNext(data)
		if err != nil {
			return nil, errors.New("CheckNext! Desc: " + err.Error())
		}
	}

	response.HasNextData = hasNextData

	return
}

func (s *Service) CheckNext(data Request) (bool, error) {
	response := new(ExchangeLeadsResponse)

	params := map[string]string{
		"site_id":  fmt.Sprint(s.config.SiteID),
		"page":     fmt.Sprint(data.Page + 1),
		"per_page": fmt.Sprint(data.PerPage),
		"from":     fmt.Sprint(data.From),
		"to":       fmt.Sprint(data.To),
		"token":    s.config.Token,
	}

	inputs := SendParams{
		Path:        leads,
		HttpMethod:  http.MethodGet,
		Response:    response,
		QueryParams: params,
	}

	if _, err := sendRequest(s.config, &inputs); err != nil {
		return false, nil
	}

	if len(response.Data) > 0 {
		return true, nil
	}

	return false, nil
}
