package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type InfoSimplesHeader struct {
	APIVersion                string         `json:"api_version"`
	APIVersionFull            string         `json:"api_version_full"`
	Product                   string         `json:"product"`
	Service                   string         `json:"service"`
	Parameters                map[string]any `json:"parameters"`
	ClientName                string         `json:"client_name"`
	TokenName                 string         `json:"token_name"`
	Billable                  bool           `json:"billable"`
	Price                     string         `json:"price"`
	RequestedAt               string         `json:"requested_at"`
	ElapsedTimeInMilliseconds int            `json:"elapsed_time_in_milliseconds"`
	RemoteIP                  string         `json:"remote_ip"`
	Signature                 string         `json:"signature"`
}

type InfoSimplesResposta struct {
	Code         int               `json:"code"`
	CodeMessage  string            `json:"code_message"`
	Errors       []string          `json:"errors"`
	Header       InfoSimplesHeader `json:"header"`
	DataCount    int               `json:"data_count"`
	Data         json.RawMessage   `json:"data"`
	SiteReceipts json.RawMessage   `json:"site_receipts"`
	Raw          string            `json:"-"`
}

func normalizarCAR(valor string) string {
	valor = strings.TrimSpace(valor)
	valor = strings.ReplaceAll(valor, " ", "")
	valor = strings.ReplaceAll(valor, "-", "")
	valor = strings.ToUpper(valor)
	return valor
}

func validarCARParaConsulta(car string) error {
	car = normalizarCAR(car)

	if car == "" {
		return errors.New("informe o número do CAR")
	}

	if len(car) < 20 {
		return errors.New("o número do CAR parece curto demais; confira antes de consultar para não gastar crédito")
	}

	return nil
}

func consultarInfoSimplesPOSTForm(endpoint string, token string, parametros map[string]string) (InfoSimplesResposta, error) {
	var resposta InfoSimplesResposta

	endpoint = strings.TrimSpace(endpoint)
	token = strings.TrimSpace(token)

	if endpoint == "" {
		return resposta, errors.New("endpoint da InfoSimples não configurado")
	}

	if token == "" {
		return resposta, errors.New("token da InfoSimples não configurado")
	}

	form := url.Values{}
	form.Set("token", token)
	form.Set("timeout", "300")

	for chave, valor := range parametros {
		form.Set(chave, strings.TrimSpace(valor))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 330*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		endpoint,
		strings.NewReader(form.Encode()),
	)
	if err != nil {
		return resposta, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return resposta, err
	}
	defer resp.Body.Close()

	corpo, err := io.ReadAll(resp.Body)
	if err != nil {
		return resposta, err
	}

	resposta.Raw = string(corpo)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return resposta, fmt.Errorf("InfoSimples respondeu HTTP %d: %s", resp.StatusCode, string(corpo))
	}

	if err := json.Unmarshal(corpo, &resposta); err != nil {
		return resposta, fmt.Errorf("erro ao ler JSON da InfoSimples: %w", err)
	}

	return resposta, nil
}

func jsonBonitoInfoSimples(raw string) string {
	var v any

	if err := json.Unmarshal([]byte(raw), &v); err != nil {
		return raw
	}

	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return raw
	}

	return string(b)
}
