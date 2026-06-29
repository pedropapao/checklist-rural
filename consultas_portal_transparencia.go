package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

type ResultadoPortalTransparencia struct {
	Encontrou bool
	Resumo    string
}

func consultarCEISPortalTransparencia(cnpj string) ResultadoPortalTransparencia {
	return consultarSancoesPortalTransparencia("ceis", cnpj)
}

func consultarCNEPPortalTransparencia(cnpj string) ResultadoPortalTransparencia {
	return consultarSancoesPortalTransparencia("cnep", cnpj)
}

func consultarSancoesPortalTransparencia(tipo string, cnpjInformado string) ResultadoPortalTransparencia {
	cnpj := limparSomenteNumeros(cnpjInformado)

	if len(cnpj) != 14 {
		return ResultadoPortalTransparencia{
			Encontrou: false,
			Resumo:    "CNPJ inválido para consulta.",
		}
	}

	chave := strings.TrimSpace(os.Getenv("PORTAL_TRANSPARENCIA_TOKEN"))
	if chave == "" {
		return ResultadoPortalTransparencia{
			Encontrou: false,
			Resumo:    "Token da API do Portal da Transparência não configurado. Defina PORTAL_TRANSPARENCIA_TOKEN.",
		}
	}

	url := fmt.Sprintf("https://api.portaldatransparencia.gov.br/api-de-dados/%s?cnpj=%s&pagina=1", tipo, cnpj)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return ResultadoPortalTransparencia{
			Encontrou: false,
			Resumo:    "Erro ao montar requisição: " + err.Error(),
		}
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("chave-api-dados", chave)

	cliente := http.Client{Timeout: 20 * time.Second}

	resp, err := cliente.Do(req)
	if err != nil {
		return ResultadoPortalTransparencia{
			Encontrou: false,
			Resumo:    "Erro ao consultar Portal da Transparência: " + err.Error(),
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return ResultadoPortalTransparencia{
			Encontrou: false,
			Resumo:    fmt.Sprintf("A API recusou a consulta. Verifique o token. Status %d.", resp.StatusCode),
		}
	}

	if resp.StatusCode != http.StatusOK {
		return ResultadoPortalTransparencia{
			Encontrou: false,
			Resumo:    fmt.Sprintf("Portal da Transparência retornou status %d.", resp.StatusCode),
		}
	}

	var resposta []map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&resposta); err != nil {
		return ResultadoPortalTransparencia{
			Encontrou: false,
			Resumo:    "Erro ao ler resposta da API: " + err.Error(),
		}
	}

	if len(resposta) == 0 {
		return ResultadoPortalTransparencia{
			Encontrou: false,
			Resumo:    "Nenhum registro encontrado.",
		}
	}

	resumos := []string{}
	limite := len(resposta)
	if limite > 5 {
		limite = 5
	}

	for i := 0; i < limite; i++ {
		item := resposta[i]
		partes := []string{}

		for _, campo := range []string{
			"nome",
			"nomeSancionado",
			"razaoSocial",
			"cpfCnpj",
			"cnpj",
			"tipoSancao",
			"descricaoSancao",
			"orgaoSancionador",
			"dataInicioSancao",
			"dataFimSancao",
		} {
			if valor, ok := item[campo]; ok && valor != nil {
				partes = append(partes, fmt.Sprintf("%s: %v", campo, valor))
			}
		}

		if len(partes) == 0 {
			b, _ := json.Marshal(item)
			partes = append(partes, string(b))
		}

		resumos = append(resumos, strings.Join(partes, " | "))
	}

	resumo := fmt.Sprintf("Foram encontrados %d registro(s).", len(resposta))
	if len(resumos) > 0 {
		resumo += "\n" + strings.Join(resumos, "\n")
	}

	return ResultadoPortalTransparencia{
		Encontrou: true,
		Resumo:    resumo,
	}
}
