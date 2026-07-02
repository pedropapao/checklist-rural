package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type InfoSimplesCARImovelItem struct {
	CAR         string `json:"car"`
	Nome        string `json:"nome"`
	Situacao    string `json:"situacao"`
	Condicao    string `json:"condicao_cadastro"`
	Area        any    `json:"area"`
	Municipio   string `json:"municipio"`
	UF          string `json:"uf"`
	SiteReceipt string `json:"site_receipt"`
}

func consultarInfoSimplesCARImovel(car string) (InfoSimplesResposta, error) {
	config := lerConfigInfoSimples()

	car = normalizarCAR(car)
	if err := validarCARParaConsulta(car); err != nil {
		return InfoSimplesResposta{}, err
	}

	return consultarInfoSimplesPOSTForm(config.CARImovelURL, config.Token, map[string]string{
		"car": car,
	})
}

func primeiroCARImovelInfoSimples(resp InfoSimplesResposta) (map[string]any, error) {
	var itens []map[string]any

	if len(resp.Data) == 0 || string(resp.Data) == "null" {
		return nil, fmt.Errorf("resposta sem dados do imóvel")
	}

	if err := json.Unmarshal(resp.Data, &itens); err != nil {
		return nil, fmt.Errorf("não consegui interpretar o retorno do CAR / Imóvel: %w", err)
	}

	if len(itens) == 0 {
		return nil, fmt.Errorf("nenhum imóvel retornado pela InfoSimples")
	}

	return itens[0], nil
}

func valorStringMapaInfoSimples(m map[string]any, chaves ...string) string {
	for _, chave := range chaves {
		if v, ok := m[chave]; ok && v != nil {
			switch t := v.(type) {
			case string:
				return strings.TrimSpace(t)
			case float64:
				return fmt.Sprintf("%.4f", t)
			case int:
				return fmt.Sprintf("%d", t)
			default:
				return strings.TrimSpace(fmt.Sprint(t))
			}
		}
	}
	return ""
}

func numeroStringMapaInfoSimples(m map[string]any, chaves ...string) string {
	for _, chave := range chaves {
		if v, ok := m[chave]; ok && v != nil {
			switch t := v.(type) {
			case float64:
				return fmt.Sprintf("%.4f ha", t)
			case int:
				return fmt.Sprintf("%d ha", t)
			case string:
				txt := strings.TrimSpace(t)
				if txt != "" {
					return txt
				}
			default:
				txt := strings.TrimSpace(fmt.Sprint(t))
				if txt != "" {
					return txt
				}
			}
		}
	}
	return ""
}

func resumoCARImovelInfoSimples(item map[string]any, resp InfoSimplesResposta) string {
	var b strings.Builder

	b.WriteString("\n\n--- Consulta InfoSimples CAR / Imóvel ---\n")
	b.WriteString("Data da consulta: " + time.Now().Format("02/01/2006 15:04") + "\n")

	car := valorStringMapaInfoSimples(item, "car", "numero_car", "codigo_car")
	if car != "" {
		b.WriteString("CAR: " + car + "\n")
	}

	nome := valorStringMapaInfoSimples(item, "nome", "nome_imovel", "imovel", "denominacao")
	if nome != "" {
		b.WriteString("Nome do imóvel: " + nome + "\n")
	}

	situacao := valorStringMapaInfoSimples(item, "situacao", "situacao_car", "status")
	if situacao != "" {
		b.WriteString("Situação: " + situacao + "\n")
	}

	condicao := valorStringMapaInfoSimples(item, "condicao_cadastro", "condicao", "condicao_do_cadastro")
	if condicao != "" {
		b.WriteString("Condição do cadastro: " + condicao + "\n")
	}

	area := numeroStringMapaInfoSimples(item, "area", "area_total", "area_imovel")
	if area != "" {
		b.WriteString("Área: " + area + "\n")
	}

	municipio := valorStringMapaInfoSimples(item, "municipio", "endereco_municipio")
	uf := valorStringMapaInfoSimples(item, "uf", "endereco_uf")
	if municipio != "" || uf != "" {
		b.WriteString("Município/UF: " + municipio + "/" + uf + "\n")
	}

	if resp.Header.Price != "" {
		b.WriteString("Preço informado pela API: R$ " + resp.Header.Price + "\n")
	}

	links := extrairSiteReceipts(resp)
	if len(links) > 0 {
		b.WriteString("Comprovante/recibo da consulta: " + links[0] + "\n")
	}

	return b.String()
}

func aplicarCARImovelInfoSimplesNosDados(d DadosExternosReuniao, item map[string]any, resp InfoSimplesResposta) DadosExternosReuniao {
	car := valorStringMapaInfoSimples(item, "car", "numero_car", "codigo_car")
	if car != "" {
		d.CAR = car
	}

	nome := valorStringMapaInfoSimples(item, "nome", "nome_imovel", "imovel", "denominacao")
	if nome != "" {
		d.NomeImovel = nome
	}

	situacao := valorStringMapaInfoSimples(item, "situacao", "situacao_car", "status")
	if situacao != "" {
		d.SituacaoCAR = situacao
	}

	area := numeroStringMapaInfoSimples(item, "area", "area_total", "area_imovel")
	if area != "" {
		d.AreaTotal = area
	}

	municipio := valorStringMapaInfoSimples(item, "municipio", "endereco_municipio")
	if municipio != "" {
		d.Municipio = municipio
	}

	uf := valorStringMapaInfoSimples(item, "uf", "endereco_uf")
	if uf != "" {
		d.UF = uf
	}

	d.FonteConsulta = "InfoSimples CAR / Imóvel"

	links := extrairSiteReceipts(resp)
	if len(links) > 0 {
		d.LinkFonte = links[0]
	}

	d.UltimaInvestigacaoOnline = time.Now().Format("02/01/2006 15:04")
	d.Observacoes = strings.TrimSpace(d.Observacoes + resumoCARImovelInfoSimples(item, resp))

	return d
}

func (app *App) telaInfoSimplesCARImovel(w http.ResponseWriter, r *http.Request) {
	tela := TelaInfoSimplesCAR{
		Titulo:         "InfoSimples - CAR Imóvel",
		Config:         lerConfigInfoSimples(),
		TokenMascarado: tokenMascarado(lerConfigInfoSimples().Token),
	}

	tela.ConfigCaminho, _ = caminhoConfigInfoSimples()

	reuniaoID, _ := strconv.Atoi(r.URL.Query().Get("id"))
	if reuniaoID > 0 {
		reuniao, err := app.buscarReuniaoPorID(reuniaoID)
		if err == nil {
			tela.Reuniao = reuniao
			tela.Dados = app.buscarDadosExternosPorReuniao(reuniaoID)
			tela.CAR = tela.Dados.CAR
		}
	}

	if r.Method == http.MethodPost {
		acao := r.FormValue("acao")
		tela.CAR = strings.TrimSpace(r.FormValue("car"))

		resposta, err := consultarInfoSimplesCARImovel(tela.CAR)
		tela.Consultou = true
		tela.Resposta = resposta
		tela.JSONBonito = jsonBonitoInfoSimples(resposta.Raw)

		if err != nil {
			tela.Erro = err.Error()
		} else {
			tela.Sucesso = "Consulta CAR / Imóvel realizada. Confira o retorno antes de salvar."

			if acao == "consultar_salvar" {
				if tela.Reuniao.ID <= 0 {
					tela.Erro = "Para salvar, abra esta tela a partir de uma reunião."
				} else {
					item, err := primeiroCARImovelInfoSimples(resposta)
					if err != nil {
						tela.Erro = err.Error()
					} else {
						dadosAtualizados := aplicarCARImovelInfoSimplesNosDados(tela.Dados, item, resposta)
						err = app.salvarDadosExternos(dadosAtualizados)
						if err != nil {
							tela.Erro = "Erro ao salvar dados na reunião: " + err.Error()
						} else {
							tela.Dados = dadosAtualizados
							tela.Sucesso = "Consulta CAR / Imóvel realizada e salva na investigação da reunião."
						}
					}
				}
			}
		}
	}

	tpl := templateMust("infosimples_car_imovel", infosimplesCARImovelHTML)
	tpl.Execute(w, tela)
}
