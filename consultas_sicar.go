package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type ResultadoSICAR struct {
	Encontrou bool
	Resumo    string
	CARs      []string
}

func consultarSICARPorCNPJ(cnpjInformado string) ResultadoSICAR {
	cnpj := limparSomenteNumeros(cnpjInformado)

	if len(cnpj) != 14 {
		return ResultadoSICAR{
			Encontrou: false,
			Resumo:    "CNPJ inválido para consulta SICAR.",
		}
	}

	modeloURL := os.Getenv("SICAR_CNPJ_URL")
	if strings.TrimSpace(modeloURL) == "" {
		modeloURL = urlPadraoSICARCNPJ
	}

	url, err := montarURLSICAR(modeloURL, cnpj)
	if err != nil {
		return ResultadoSICAR{
			Encontrou: false,
			Resumo:    err.Error(),
		}
	}

	return consultarSICARURL(url, "SICAR por CNPJ")
}

func consultarSICARPorCAR(carInformado string) ResultadoSICAR {
	car := strings.TrimSpace(carInformado)

	if car == "" {
		return ResultadoSICAR{
			Encontrou: false,
			Resumo:    "Informe o número/código do CAR para consultar.",
		}
	}

	modeloURL := os.Getenv("SICAR_CAR_URL")
	if strings.TrimSpace(modeloURL) == "" {
		modeloURL = urlPadraoSICARImovel
	}

	url, err := montarURLSICAR(modeloURL, car)
	if err != nil {
		return ResultadoSICAR{
			Encontrou: false,
			Resumo:    err.Error(),
		}
	}

	return consultarSICARURL(url, "SICAR por CAR")
}

func consultarSICARURL(url string, tipo string) ResultadoSICAR {
	token, err := obterTokenConectaGov()
	if err != nil {
		return ResultadoSICAR{
			Encontrou: false,
			Resumo:    err.Error(),
		}
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return ResultadoSICAR{
			Encontrou: false,
			Resumo:    "Erro ao montar requisição SICAR: " + err.Error(),
		}
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	cliente := http.Client{Timeout: 30 * time.Second}

	resp, err := cliente.Do(req)
	if err != nil {
		return ResultadoSICAR{
			Encontrou: false,
			Resumo:    "Erro ao consultar SICAR: " + err.Error(),
		}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ResultadoSICAR{
			Encontrou: false,
			Resumo:    "Erro ao ler resposta SICAR: " + err.Error(),
		}
	}

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return ResultadoSICAR{
			Encontrou: false,
			Resumo:    fmt.Sprintf("A API recusou a consulta. Verifique token/credencial Conecta Gov. Status %d.", resp.StatusCode),
		}
	}

	if resp.StatusCode == http.StatusNotFound {
		return ResultadoSICAR{
			Encontrou: false,
			Resumo:    "Nenhum registro encontrado no SICAR para essa informação.",
		}
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		resumoCurto := string(body)
		if len(resumoCurto) > 1000 {
			resumoCurto = resumoCurto[:1000] + "..."
		}

		return ResultadoSICAR{
			Encontrou: false,
			Resumo:    fmt.Sprintf("SICAR retornou status %d. Resposta: %s", resp.StatusCode, resumoCurto),
		}
	}

	resumo := resumirJSONGenerico(body)

	return ResultadoSICAR{
		Encontrou: true,
		Resumo:    tipo + " consultado com sucesso.\n\n" + resumo,
	}
}

func resumirJSONGenerico(body []byte) string {
	var valor any

	if err := json.Unmarshal(body, &valor); err != nil {
		texto := string(body)
		if len(texto) > 2000 {
			texto = texto[:2000] + "..."
		}
		return texto
	}

	formatado, err := json.MarshalIndent(valor, "", "  ")
	if err != nil {
		texto := string(body)
		if len(texto) > 2000 {
			texto = texto[:2000] + "..."
		}
		return texto
	}

	resumo := string(formatado)
	if len(resumo) > 4000 {
		resumo = resumo[:4000] + "\n..."
	}

	return resumo
}

func acrescentarBlocoObservacaoSICAR(observacoes string, titulo string, resumo string) string {
	observacoes = strings.TrimSpace(observacoes)

	bloco := strings.TrimSpace(titulo) + "\n" + strings.TrimSpace(resumo)

	if observacoes == "" {
		return bloco
	}

	return observacoes + "\n\n" + bloco
}

func extrairCARsDeJSON(body []byte) []string {
	var valor any

	if err := json.Unmarshal(body, &valor); err != nil {
		return []string{}
	}

	encontrados := []string{}
	vistos := map[string]bool{}

	percorrerJSONParaCAR(valor, &encontrados, vistos)

	return encontrados
}

func percorrerJSONParaCAR(valor any, encontrados *[]string, vistos map[string]bool) {
	switch v := valor.(type) {
	case map[string]any:
		for chave, conteudo := range v {
			chaveLimpa := strings.ToLower(chave)

			if texto, ok := conteudo.(string); ok {
				texto = strings.TrimSpace(texto)

				if pareceCodigoCAR(chaveLimpa, texto) {
					adicionarCAREncontrado(texto, encontrados, vistos)
				}
			}

			percorrerJSONParaCAR(conteudo, encontrados, vistos)
		}

	case []any:
		for _, item := range v {
			percorrerJSONParaCAR(item, encontrados, vistos)
		}

	case string:
		texto := strings.TrimSpace(v)
		if pareceFormatoCAR(texto) {
			adicionarCAREncontrado(texto, encontrados, vistos)
		}
	}
}

func pareceCodigoCAR(chave string, valor string) bool {
	if valor == "" {
		return false
	}

	if pareceFormatoCAR(valor) {
		return true
	}

	chave = strings.ToLower(chave)

	if strings.Contains(chave, "car") ||
		strings.Contains(chave, "codigocar") ||
		strings.Contains(chave, "codigo_car") ||
		strings.Contains(chave, "codigoimovel") ||
		strings.Contains(chave, "codigo_imovel") ||
		strings.Contains(chave, "codimovel") ||
		strings.Contains(chave, "cod_imovel") {
		return len(valor) >= 10
	}

	return false
}

func pareceFormatoCAR(valor string) bool {
	valor = strings.TrimSpace(strings.ToUpper(valor))

	if len(valor) < 20 {
		return false
	}

	if len(valor) > 80 {
		return false
	}

	if !strings.Contains(valor, "-") {
		return false
	}

	partes := strings.Split(valor, "-")
	if len(partes) < 3 {
		return false
	}

	if len(partes[0]) != 2 {
		return false
	}

	return true
}

func adicionarCAREncontrado(valor string, encontrados *[]string, vistos map[string]bool) {
	valor = strings.TrimSpace(strings.ToUpper(valor))

	if valor == "" {
		return
	}

	if vistos[valor] {
		return
	}

	vistos[valor] = true
	*encontrados = append(*encontrados, valor)
}

func substituirBlocoConsultaSICAR(observacoes string, titulo string, resumo string) string {
	observacoes = strings.TrimSpace(observacoes)
	titulo = strings.TrimSpace(titulo)
	resumo = strings.TrimSpace(resumo)

	novoBloco := titulo + "\n" + resumo

	if observacoes == "" {
		return novoBloco
	}

	partes := strings.Split(observacoes, "\n\n")
	filtradas := []string{}

	for _, parte := range partes {
		parteLimpa := strings.TrimSpace(parte)
		if parteLimpa == "" {
			continue
		}

		if strings.HasPrefix(parteLimpa, titulo) {
			continue
		}

		filtradas = append(filtradas, parteLimpa)
	}

	filtradas = append(filtradas, novoBloco)

	return strings.Join(filtradas, "\n\n")
}

func consultarSICARTemaPorCAR(carInformado string) ResultadoSICAR {
	car := strings.TrimSpace(carInformado)

	if car == "" {
		return ResultadoSICAR{
			Encontrou: false,
			Resumo:    "Informe o número/código do CAR para consultar temas, área ou polígono.",
		}
	}

	modeloURL := os.Getenv("SICAR_TEMA_URL")
	if strings.TrimSpace(modeloURL) == "" {
		modeloURL = urlPadraoSICARTema
	}

	url, err := montarURLSICAR(modeloURL, car)
	if err != nil {
		return ResultadoSICAR{
			Encontrou: false,
			Resumo:    err.Error(),
		}
	}

	return consultarSICARURL(url, "SICAR Tema / Área / Polígono")
}

func montarURLSICAR(modeloURL string, valor string) (string, error) {
	modeloURL = strings.TrimSpace(modeloURL)
	valor = strings.TrimSpace(valor)

	if modeloURL == "" {
		return "", fmt.Errorf("URL SICAR não configurada.")
	}

	if !strings.Contains(modeloURL, "%s") {
		return "", fmt.Errorf("URL SICAR inválida. A URL precisa conter %%s no lugar onde o CNPJ/CAR será inserido.")
	}

	return fmt.Sprintf(modeloURL, valor), nil
}

func limparBlocosSICARObservacoes(observacoes string) string {
	observacoes = strings.TrimSpace(observacoes)

	if observacoes == "" {
		return ""
	}

	partes := strings.Split(observacoes, "\n\n")
	filtradas := []string{}

	for _, parte := range partes {
		parteLimpa := strings.TrimSpace(parte)
		parteMinuscula := strings.ToLower(parteLimpa)

		if parteLimpa == "" {
			continue
		}

		if strings.HasPrefix(parteMinuscula, "consulta sicar por cnpj:") ||
			strings.HasPrefix(parteMinuscula, "consulta sicar por car:") ||
			strings.HasPrefix(parteMinuscula, "consulta sicar tema") {
			continue
		}

		filtradas = append(filtradas, parteLimpa)
	}

	return strings.Join(filtradas, "\n\n")
}
