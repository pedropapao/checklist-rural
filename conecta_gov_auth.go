package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type respostaTokenConectaGov struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
}

func obterTokenConectaGov() (string, error) {
	// 1. Se o usuário colou token manual, usa ele.
	tokenManual := strings.TrimSpace(os.Getenv("CONECTA_GOV_TOKEN"))
	if tokenManual != "" {
		return tokenManual, nil
	}

	clientID := strings.TrimSpace(os.Getenv("CONECTA_GOV_CLIENT_ID"))
	clientSecret := strings.TrimSpace(os.Getenv("CONECTA_GOV_CLIENT_SECRET"))

	if clientID == "" || clientSecret == "" {
		return "", fmt.Errorf("Conecta Gov não configurado. Informe Client ID e Client Secret na tela Configurações das APIs")
	}

	tokenURL := strings.TrimSpace(os.Getenv("CONECTA_GOV_TOKEN_URL"))
	if tokenURL == "" {
		tokenURL = "https://apigateway.conectagov.estaleiro.serpro.gov.br/oauth2/jwt-token"
	}

	form := url.Values{}
	form.Set("grant_type", "client_credentials")

	req, err := http.NewRequest(http.MethodPost, tokenURL, strings.NewReader(form.Encode()))
	if err != nil {
		return "", err
	}

	credencial := base64.StdEncoding.EncodeToString([]byte(clientID + ":" + clientSecret))

	req.Header.Set("Authorization", "Basic "+credencial)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	cliente := http.Client{Timeout: 60 * time.Second}

	resp, err := cliente.Do(req)
	if err != nil {
		return "", fmt.Errorf("erro ao gerar token Conecta Gov: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		resposta := strings.TrimSpace(string(body))
		if len(resposta) > 1000 {
			resposta = resposta[:1000] + "..."
		}

		return "", fmt.Errorf("Conecta Gov recusou a geração do token. Status %d. Resposta: %s", resp.StatusCode, resposta)
	}

	var token respostaTokenConectaGov
	if err := json.Unmarshal(body, &token); err != nil {
		return "", fmt.Errorf("erro ao interpretar resposta do token Conecta Gov: %w", err)
	}

	if strings.TrimSpace(token.AccessToken) == "" {
		return "", fmt.Errorf("Conecta Gov respondeu, mas não retornou access_token")
	}

	return token.AccessToken, nil
}
