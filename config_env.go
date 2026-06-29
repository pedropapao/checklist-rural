package main

import (
	"bufio"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type ConfigAPI struct {
	ConectaGovClientID       string
	ConectaGovClientSecret   string
	ConectaGovToken          string
	SICARCNPJURL             string
	SICARCARURL              string
	SICARTemaURL             string
	PortalTransparenciaToken string
}

func caminhoConfigEnv() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	pasta := filepath.Join(home, "Documentos", "ChecklistRural")
	if err := os.MkdirAll(pasta, 0755); err != nil {
		return "", err
	}

	return filepath.Join(pasta, "config.env"), nil
}

func carregarConfigEnv() {
	caminho, err := caminhoConfigEnv()
	if err != nil {
		return
	}

	arquivo, err := os.Open(caminho)
	if err != nil {
		return
	}
	defer arquivo.Close()

	scanner := bufio.NewScanner(arquivo)

	for scanner.Scan() {
		linha := strings.TrimSpace(scanner.Text())

		if linha == "" || strings.HasPrefix(linha, "#") {
			continue
		}

		partes := strings.SplitN(linha, "=", 2)
		if len(partes) != 2 {
			continue
		}

		chave := strings.TrimSpace(partes[0])
		valor := strings.TrimSpace(partes[1])
		valor = strings.Trim(valor, `"`)

		if chave != "" && os.Getenv(chave) == "" {
			os.Setenv(chave, valor)
		}
	}
}

func lerConfigAPI() ConfigAPI {
	return ConfigAPI{
		ConectaGovClientID:       os.Getenv("CONECTA_GOV_CLIENT_ID"),
		ConectaGovClientSecret:   os.Getenv("CONECTA_GOV_CLIENT_SECRET"),
		ConectaGovToken:          os.Getenv("CONECTA_GOV_TOKEN"),
		SICARCNPJURL:             os.Getenv("SICAR_CNPJ_URL"),
		SICARCARURL:              os.Getenv("SICAR_CAR_URL"),
		SICARTemaURL:             os.Getenv("SICAR_TEMA_URL"),
		PortalTransparenciaToken: os.Getenv("PORTAL_TRANSPARENCIA_TOKEN"),
	}
}

func salvarConfigAPI(c ConfigAPI) error {
	caminho, err := caminhoConfigEnv()
	if err != nil {
		return err
	}

	conteudo := []string{
		"# Configurações privadas do Checklist Rural",
		"# Não envie este arquivo para o GitHub.",
		"",
		"CONECTA_GOV_CLIENT_ID=" + strings.TrimSpace(c.ConectaGovClientID),
		"CONECTA_GOV_CLIENT_SECRET=" + strings.TrimSpace(c.ConectaGovClientSecret),
		"CONECTA_GOV_TOKEN=" + strings.TrimSpace(c.ConectaGovToken),
		"SICAR_CNPJ_URL=" + strings.TrimSpace(c.SICARCNPJURL),
		"SICAR_CAR_URL=" + strings.TrimSpace(c.SICARCARURL),
		"SICAR_TEMA_URL=" + strings.TrimSpace(c.SICARTemaURL),
		"PORTAL_TRANSPARENCIA_TOKEN=" + strings.TrimSpace(c.PortalTransparenciaToken),
		"",
	}

	err = os.WriteFile(caminho, []byte(strings.Join(conteudo, "\n")), 0600)
	if err != nil {
		return err
	}

	os.Setenv("CONECTA_GOV_CLIENT_ID", strings.TrimSpace(c.ConectaGovClientID))
	os.Setenv("CONECTA_GOV_CLIENT_SECRET", strings.TrimSpace(c.ConectaGovClientSecret))
	os.Setenv("CONECTA_GOV_TOKEN", strings.TrimSpace(c.ConectaGovToken))
	os.Setenv("SICAR_CNPJ_URL", strings.TrimSpace(c.SICARCNPJURL))
	os.Setenv("SICAR_CAR_URL", strings.TrimSpace(c.SICARCARURL))
	os.Setenv("SICAR_TEMA_URL", strings.TrimSpace(c.SICARTemaURL))
	os.Setenv("PORTAL_TRANSPARENCIA_TOKEN", strings.TrimSpace(c.PortalTransparenciaToken))

	return nil
}

func tokenMascarado(valor string) string {
	valor = strings.TrimSpace(valor)

	if valor == "" {
		return "não configurado"
	}

	if len(valor) <= 8 {
		return "configurado"
	}

	return valor[:4] + "..." + valor[len(valor)-4:]
}

func (app *App) telaConfiguracoesAPI(w http.ResponseWriter, r *http.Request) {
	mensagem := ""

	if r.Method == http.MethodPost {
		config := ConfigAPI{
			ConectaGovClientID:       r.FormValue("conecta_gov_client_id"),
			ConectaGovClientSecret:   r.FormValue("conecta_gov_client_secret"),
			ConectaGovToken:          r.FormValue("conecta_gov_token"),
			SICARCNPJURL:             r.FormValue("sicar_cnpj_url"),
			SICARCARURL:              r.FormValue("sicar_car_url"),
			SICARTemaURL:             r.FormValue("sicar_tema_url"),
			PortalTransparenciaToken: r.FormValue("portal_transparencia_token"),
		}

		err := salvarConfigAPI(config)
		if err != nil {
			http.Error(w, "Erro ao salvar configurações: "+err.Error(), http.StatusInternalServerError)
			return
		}

		mensagem = "Configurações salvas com sucesso."
	}

	config := aplicarPadroesConfigAPI(lerConfigAPI())

	dados := map[string]any{
		"Titulo":                       "Configurações das APIs",
		"Config":                       config,
		"Mensagem":                     mensagem,
		"ConectaGovTokenMascarado":     tokenMascarado(config.ConectaGovToken),
		"PortalTransparenciaMascarado": tokenMascarado(config.PortalTransparenciaToken),
		"CaminhoConfig":                func() string { c, _ := caminhoConfigEnv(); return c }(),
	}

	tpl := templateMust("config_api", configuracoesAPIHTML)
	tpl.Execute(w, dados)
}

func templateMust(nome string, conteudo string) *template.Template {
	return template.Must(template.New(nome).Parse(htmlBase(conteudo)))
}

const urlPadraoSICARCNPJ = "https://apigateway.conectagov.estaleiro.serpro.gov.br/api-sicar-cpfcnpj/v1/%s"
const urlPadraoSICARImovel = "https://apigateway.conectagov.estaleiro.serpro.gov.br/api-sicar-imovel/v1/%s"
const urlPadraoSICARTema = "https://apigateway.conectagov.estaleiro.serpro.gov.br/api-sicar-tema/v1/%s"

func aplicarPadroesConfigAPI(c ConfigAPI) ConfigAPI {
	if strings.TrimSpace(c.SICARCNPJURL) == "" {
		c.SICARCNPJURL = urlPadraoSICARCNPJ
	}

	if strings.TrimSpace(c.SICARCARURL) == "" {
		c.SICARCARURL = urlPadraoSICARImovel
	}

	if strings.TrimSpace(c.SICARTemaURL) == "" {
		c.SICARTemaURL = urlPadraoSICARTema
	}

	return c
}
