package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type InfoSimplesCARDemonstrativoItem struct {
	CAR                       string  `json:"car"`
	CondicaoCadastro          string  `json:"condicao_cadastro"`
	Situacao                  string  `json:"situacao"`
	SiteReceipt               string  `json:"site_receipt"`
	AreaPreservacaoPermanente float64 `json:"area_preservacao_permanente"`
	AreaUsoRestrito           float64 `json:"area_uso_restrito"`

	Imovel struct {
		AnaliseData       string  `json:"analise_data"`
		Area              float64 `json:"area"`
		EnderecoLatitude  string  `json:"endereco_latitude"`
		EnderecoLongitude string  `json:"endereco_longitude"`
		EnderecoMunicipio string  `json:"endereco_municipio"`
		EnderecoUF        string  `json:"endereco_uf"`
		ModulosFiscais    float64 `json:"modulos_fiscais"`
		RegistroData      string  `json:"registro_data"`
		RetificacaoData   string  `json:"retificacao_data"`
	} `json:"imovel"`

	RegularidadeAmbiental struct {
		AreaPreservacaoPermanenteRecompor float64 `json:"area_preservacao_permanente_recompor"`
		AreaReservaLegalRecompor          float64 `json:"area_reserva_legal_recompor"`
		PassivoExcedenteReservaLegal      float64 `json:"passivo_excedente_reserva_legal"`
	} `json:"regularidade_ambiental"`

	Reserva struct {
		AreaAverbada           float64 `json:"area_averbada"`
		AreaAverbadaDocumental float64 `json:"area_averbada_documental"`
		AreaLegalDeclarada     float64 `json:"area_legal_declarada"`
		AreaLegalProposta      float64 `json:"area_legal_proposta"`
		AreaNaoAverbada        float64 `json:"area_nao_averbada"`
		Situacao               string  `json:"situacao"`
	} `json:"reserva"`

	Solo struct {
		AreaNativa                 float64 `json:"area_nativa"`
		AreaServidaoAdministrativa float64 `json:"area_servidao_administrativa"`
		AreaUso                    float64 `json:"area_uso"`
	} `json:"solo"`
}

func primeiroCARDemonstrativoInfoSimples(resp InfoSimplesResposta) (InfoSimplesCARDemonstrativoItem, error) {
	var itens []InfoSimplesCARDemonstrativoItem

	if len(resp.Data) == 0 || string(resp.Data) == "null" {
		return InfoSimplesCARDemonstrativoItem{}, fmt.Errorf("resposta sem dados do CAR")
	}

	if err := json.Unmarshal(resp.Data, &itens); err != nil {
		return InfoSimplesCARDemonstrativoItem{}, fmt.Errorf("não consegui interpretar o retorno do CAR: %w", err)
	}

	if len(itens) == 0 {
		return InfoSimplesCARDemonstrativoItem{}, fmt.Errorf("nenhum CAR retornado pela InfoSimples")
	}

	return itens[0], nil
}

func textoResumoCARDemonstrativoInfoSimples(item InfoSimplesCARDemonstrativoItem, resp InfoSimplesResposta) string {
	var b strings.Builder

	b.WriteString("\n\n--- Consulta InfoSimples CAR / Demonstrativo ---\n")
	b.WriteString("Data da consulta: " + time.Now().Format("02/01/2006 15:04") + "\n")
	b.WriteString("CAR: " + item.CAR + "\n")
	b.WriteString("Situação: " + item.Situacao + "\n")
	b.WriteString("Condição do cadastro: " + item.CondicaoCadastro + "\n")
	b.WriteString(fmt.Sprintf("Área do imóvel: %.4f ha\n", item.Imovel.Area))
	b.WriteString(fmt.Sprintf("Módulos fiscais: %.4f\n", item.Imovel.ModulosFiscais))
	b.WriteString("Município/UF: " + item.Imovel.EnderecoMunicipio + "/" + item.Imovel.EnderecoUF + "\n")
	b.WriteString("Latitude: " + item.Imovel.EnderecoLatitude + "\n")
	b.WriteString("Longitude: " + item.Imovel.EnderecoLongitude + "\n")
	b.WriteString("Registro: " + item.Imovel.RegistroData + "\n")
	b.WriteString("Retificação: " + item.Imovel.RetificacaoData + "\n")
	b.WriteString("Reserva legal: " + item.Reserva.Situacao + "\n")
	b.WriteString(fmt.Sprintf("Área de uso do solo: %.4f ha\n", item.Solo.AreaUso))
	b.WriteString(fmt.Sprintf("Área nativa: %.4f ha\n", item.Solo.AreaNativa))
	b.WriteString(fmt.Sprintf("APP a recompor: %.4f ha\n", item.RegularidadeAmbiental.AreaPreservacaoPermanenteRecompor))
	b.WriteString(fmt.Sprintf("Reserva legal a recompor: %.4f ha\n", item.RegularidadeAmbiental.AreaReservaLegalRecompor))

	if resp.Header.Price != "" {
		b.WriteString("Preço informado pela API: R$ " + resp.Header.Price + "\n")
	}

	if item.SiteReceipt != "" {
		b.WriteString("Comprovante/recibo da consulta: " + item.SiteReceipt + "\n")
	}

	return b.String()
}

func aplicarCARDemonstrativoInfoSimplesNosDados(d DadosExternosReuniao, item InfoSimplesCARDemonstrativoItem, resp InfoSimplesResposta) DadosExternosReuniao {
	if strings.TrimSpace(item.CAR) != "" {
		d.CAR = item.CAR
	}

	if strings.TrimSpace(item.Situacao) != "" {
		d.SituacaoCAR = item.Situacao
	}

	if item.Imovel.Area > 0 {
		d.AreaTotal = fmt.Sprintf("%.4f ha", item.Imovel.Area)
	}

	if strings.TrimSpace(item.Imovel.EnderecoMunicipio) != "" {
		d.Municipio = item.Imovel.EnderecoMunicipio
	}

	if strings.TrimSpace(item.Imovel.EnderecoUF) != "" {
		d.UF = item.Imovel.EnderecoUF
	}

	d.FonteConsulta = "InfoSimples CAR / Demonstrativo"

	if strings.TrimSpace(item.SiteReceipt) != "" {
		d.LinkFonte = item.SiteReceipt
	}

	d.UltimaInvestigacaoOnline = time.Now().Format("02/01/2006 15:04")

	resumo := textoResumoCARDemonstrativoInfoSimples(item, resp)
	if !strings.Contains(d.Observacoes, "--- Consulta InfoSimples CAR / Demonstrativo ---") {
		d.Observacoes = strings.TrimSpace(d.Observacoes + resumo)
	} else {
		d.Observacoes = strings.TrimSpace(d.Observacoes + "\n\n" + resumo)
	}

	return d
}
