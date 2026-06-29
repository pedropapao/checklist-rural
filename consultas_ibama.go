package main

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type ResultadoIbama struct {
	Encontrou bool
	Resumo    string
}

const urlIbamaTermoEmbargoJSON = "https://dadosabertos.ibama.gov.br/dados/SIFISC/termo_embargo/termo_embargo/termo_embargo_json.zip"

func consultarEmbargoIbamaPorCNPJ(cnpjInformado string) ResultadoIbama {
	cnpj := limparSomenteNumeros(cnpjInformado)

	if len(cnpj) != 14 {
		return ResultadoIbama{
			Encontrou: false,
			Resumo:    "CNPJ inválido para consulta Ibama.",
		}
	}

	caminhoZip, err := baixarOuUsarCacheIbama()
	if err != nil {
		return ResultadoIbama{
			Encontrou: false,
			Resumo:    "Erro ao obter base do Ibama: " + err.Error(),
		}
	}

	resultado, err := procurarCNPJNoZipIbama(caminhoZip, cnpj)
	if err != nil {
		return ResultadoIbama{
			Encontrou: false,
			Resumo:    "Erro ao pesquisar base do Ibama: " + err.Error(),
		}
	}

	return resultado
}

func baixarOuUsarCacheIbama() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	pastaCache := filepath.Join(home, "Documentos", "ChecklistRural", "cache")
	err = os.MkdirAll(pastaCache, 0755)
	if err != nil {
		return "", err
	}

	caminhoZip := filepath.Join(pastaCache, "ibama_termo_embargo_json.zip")

	info, err := os.Stat(caminhoZip)
	if err == nil {
		idade := time.Since(info.ModTime())
		if idade < 7*24*time.Hour && info.Size() > 0 {
			return caminhoZip, nil
		}
	}

	resp, err := http.Get(urlIbamaTermoEmbargoJSON)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("download retornou status %d", resp.StatusCode)
	}

	arquivoTemp := caminhoZip + ".tmp"

	out, err := os.Create(arquivoTemp)
	if err != nil {
		return "", err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", err
	}

	err = os.Rename(arquivoTemp, caminhoZip)
	if err != nil {
		return "", err
	}

	return caminhoZip, nil
}

func procurarCNPJNoZipIbama(caminhoZip string, cnpj string) (ResultadoIbama, error) {
	zipReader, err := zip.OpenReader(caminhoZip)
	if err != nil {
		return ResultadoIbama{}, err
	}
	defer zipReader.Close()

	for _, arquivo := range zipReader.File {
		nome := strings.ToLower(arquivo.Name)

		if !strings.HasSuffix(nome, ".json") {
			continue
		}

		f, err := arquivo.Open()
		if err != nil {
			return ResultadoIbama{}, err
		}
		defer f.Close()

		return procurarCNPJNoJSONIbama(f, cnpj)
	}

	return ResultadoIbama{
		Encontrou: false,
		Resumo:    "Nenhum arquivo JSON encontrado dentro do ZIP do Ibama.",
	}, nil
}

func procurarCNPJNoJSONIbama(r io.Reader, cnpj string) (ResultadoIbama, error) {
	var raiz any

	decoder := json.NewDecoder(r)
	if err := decoder.Decode(&raiz); err != nil {
		return ResultadoIbama{}, err
	}

	registros := extrairRegistrosIbama(raiz)

	if len(registros) == 0 {
		return ResultadoIbama{
			Encontrou: false,
			Resumo:    "O JSON do Ibama foi lido, mas nenhum registro foi identificado no formato esperado.",
		}, nil
	}

	encontrados := []string{}
	total := 0

	for _, item := range registros {
		campoCNPJ := valorStringMapa(item, "CPF_CNPJ_EMBARGADO")
		if campoCNPJ == "" {
			campoCNPJ = valorStringMapa(item, "cpf_cnpj_embargado")
		}
		if campoCNPJ == "" {
			campoCNPJ = valorStringMapa(item, "cpfCnpjEmbargado")
		}
		if campoCNPJ == "" {
			campoCNPJ = valorStringMapa(item, "cpfCnpj")
		}
		if campoCNPJ == "" {
			campoCNPJ = valorStringMapa(item, "cnpj")
		}

		campoCNPJLimpo := limparSomenteNumeros(campoCNPJ)

		if campoCNPJLimpo == cnpj {
			total++

			if len(encontrados) < 5 {
				encontrados = append(encontrados, resumirRegistroIbama(item))
			}
		}
	}

	if total == 0 {
		return ResultadoIbama{
			Encontrou: false,
			Resumo:    "Nenhum termo de embargo encontrado para este CNPJ na base aberta do Ibama.",
		}, nil
	}

	resumo := fmt.Sprintf("Foram encontrados %d registro(s) de termo de embargo para este CNPJ.", total)

	if len(encontrados) > 0 {
		resumo += "\n\nPrimeiros registros:\n" + strings.Join(encontrados, "\n\n")
	}

	return ResultadoIbama{
		Encontrou: true,
		Resumo:    resumo,
	}, nil
}

func extrairRegistrosIbama(valor any) []map[string]any {
	registros := []map[string]any{}

	switch v := valor.(type) {
	case []any:
		for _, item := range v {
			if mapa, ok := item.(map[string]any); ok {
				registros = append(registros, normalizarRegistroIbama(mapa))
			}
		}

	case map[string]any:
		// Formato GeoJSON: {"type":"FeatureCollection","features":[{"properties":{...}}]}
		if features, ok := v["features"].([]any); ok {
			for _, feature := range features {
				if mapaFeature, ok := feature.(map[string]any); ok {
					if props, ok := mapaFeature["properties"].(map[string]any); ok {
						registros = append(registros, normalizarRegistroIbama(props))
					} else {
						registros = append(registros, normalizarRegistroIbama(mapaFeature))
					}
				}
			}
			return registros
		}

		// Outros formatos comuns: {"data":[...]}, {"records":[...]}, {"resultado":[...]}
		for _, chave := range []string{"data", "dados", "records", "registros", "resultado", "resultados"} {
			if lista, ok := v[chave].([]any); ok {
				for _, item := range lista {
					if mapa, ok := item.(map[string]any); ok {
						registros = append(registros, normalizarRegistroIbama(mapa))
					}
				}
				return registros
			}
		}

		// Se for um único registro em objeto
		registros = append(registros, normalizarRegistroIbama(v))
	}

	return registros
}

func normalizarRegistroIbama(m map[string]any) map[string]any {
	normalizado := map[string]any{}

	for chave, valor := range m {
		normalizado[chave] = valor
		normalizado[strings.ToUpper(chave)] = valor
		normalizado[strings.ToLower(chave)] = valor
	}

	return normalizado
}

func resumirRegistroIbama(item map[string]any) string {
	gruposCampos := [][]string{
		{"NOME_EMBARGADO", "nome_embargado", "nomeEmbargado", "nome"},
		{"CPF_CNPJ_EMBARGADO", "cpf_cnpj_embargado", "cpfCnpjEmbargado", "cpfCnpj", "cnpj"},
		{"NUM_TAD", "num_tad", "numTad", "numeroTad"},
		{"DAT_EMBARGO", "dat_embargo", "dataEmbargo"},
		{"DES_STATUS_FORMULARIO", "des_status_formulario", "status"},
		{"SIT_CANCELADO", "sit_cancelado", "cancelado"},
		{"SIT_DESEMBARGO", "sit_desembargo", "desembargo"},
		{"MUNICIPIO", "municipio", "município"},
		{"UF", "uf"},
		{"NOME_IMOVEL", "nome_imovel", "nomeImovel"},
		{"QTD_AREA_EMBARGADA", "qtd_area_embargada", "areaEmbargada"},
		{"DES_LOCALIZACAO", "des_localizacao", "localizacao", "localização"},
	}

	partes := []string{}

	for _, opcoes := range gruposCampos {
		valor := ""

		for _, campo := range opcoes {
			valor = strings.TrimSpace(valorStringMapa(item, campo))
			if valor != "" {
				partes = append(partes, opcoes[0]+": "+valor)
				break
			}
		}
	}

	if len(partes) == 0 {
		return "Registro encontrado, mas sem campos resumidos."
	}

	return strings.Join(partes, " | ")
}

func valorStringMapa(m map[string]any, chave string) string {
	valor, ok := m[chave]
	if !ok || valor == nil {
		return ""
	}

	return fmt.Sprintf("%v", valor)
}
