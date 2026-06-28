package main

import "strings"

type LeituraInicial struct {
	Enquadramento  string
	Caminho        string
	Resumo         string
	Alertas        []string
	ProximosPassos []string
}

func montarLeituraInicial(r Reuniao) LeituraInicial {
	leitura := LeituraInicial{
		Enquadramento:  r.ClassificacaoProdutor,
		Caminho:        "A confirmar conforme documentos, cadastro e política do banco.",
		Alertas:        []string{},
		ProximosPassos: []string{},
	}

	finalidade := strings.ToLower(r.FinalidadeCredito)
	tipoProjeto := strings.ToLower(r.TipoProjeto + " " + r.Atividade + " " + r.FinalidadeCredito)

	if leitura.Enquadramento == "" {
		leitura.Enquadramento = "Não classificado"
	}

	// Caminho inicial
	if r.PossuiCAF == "sim" && r.ClassificacaoProdutor == "Pequeno produtor" {
		leitura.Caminho = "Possível caminho para agricultura familiar / Pronaf, se o CAF estiver ativo e o produtor atender às demais exigências."
	} else if r.ClassificacaoProdutor == "Médio produtor" {
		leitura.Caminho = "Possível caminho para médio produtor / Pronamp, dependendo da finalidade, cadastro e limite aprovado."
	} else if r.ClassificacaoProdutor == "Grande produtor" {
		leitura.Caminho = "Possível caminho por linhas gerais do crédito rural, conforme cadastro, limite e política do banco."
	} else if r.PossuiCAF == "sim" {
		leitura.Caminho = "Existe indicação de CAF. Conferir se o produtor pode ser tratado como agricultura familiar."
	}

	// Resumo da demanda
	if r.FinalidadeCredito != "" {
		leitura.Resumo = "Finalidade principal informada: " + r.FinalidadeCredito + "."
	} else {
		leitura.Resumo = "A finalidade principal do crédito ainda não foi informada."
		leitura.Alertas = append(leitura.Alertas, "Definir se o crédito é para custeio, investimento, máquinas, obra, irrigação ou pecuária.")
	}

	// Alertas de banco/cadastro
	if r.CadastroBanco != "sim" {
		leitura.Alertas = append(leitura.Alertas, "Conferir cadastro do produtor no banco ou cooperativa.")
	}

	if r.FinanciamentoAtivo == "sim" {
		leitura.Alertas = append(leitura.Alertas, "Verificar financiamentos rurais ativos, saldo devedor, vencimentos e limite disponível.")
	}

	if r.RestricaoCadastral == "sim" {
		leitura.Alertas = append(leitura.Alertas, "Existe indicação de restrição ou pendência cadastral. Resolver isso antes de avançar com segurança.")
	}

	// Alertas de terra
	if r.ImovelArrendado == "sim" {
		leitura.Alertas = append(leitura.Alertas, "Área arrendada, parceria ou comodato: conferir contrato, prazo, autorização e anuência do proprietário.")
	}

	if r.ImovelProprio != "sim" && r.ImovelArrendado != "sim" {
		leitura.Alertas = append(leitura.Alertas, "Definir claramente a situação da área: própria, arrendada, parceria ou comodato.")
	}

	if r.TemCAR != "sim" {
		leitura.Alertas = append(leitura.Alertas, "Conferir CAR da propriedade antes de avançar.")
	}

	// Alertas ambientais
	if r.UsaAgua == "sim" || strings.Contains(finalidade, "irrig") || strings.Contains(tipoProjeto, "irrig") {
		leitura.Alertas = append(leitura.Alertas, "Projeto envolve água ou irrigação: verificar outorga, dispensa ou regularidade do uso da água.")
	}

	if r.TemSupressao == "sim" {
		leitura.Alertas = append(leitura.Alertas, "Projeto envolve abertura de área ou vegetação: conferir autorização ambiental antes de assumir o projeto.")
	}

	// Alertas técnicos
	if r.TemOrcamento != "sim" {
		leitura.Alertas = append(leitura.Alertas, "Solicitar orçamento, proposta comercial ou memória de cálculo do valor pretendido.")
	}

	if r.TemObra == "sim" || strings.Contains(finalidade, "obra") || strings.Contains(finalidade, "benfeitoria") {
		leitura.Alertas = append(leitura.Alertas, "Projeto envolve obra ou benfeitoria: conferir orçamento detalhado, croqui, memorial e responsabilidade técnica.")
	}

	if r.TemInvestimento == "sim" || strings.Contains(finalidade, "invest") || strings.Contains(finalidade, "máquina") || strings.Contains(finalidade, "maquina") {
		leitura.Alertas = append(leitura.Alertas, "Projeto envolve investimento: conferir orçamento, nota/proposta, capacidade de pagamento e garantias.")
	}

	if r.TemPecuaria == "sim" || strings.Contains(finalidade, "pecu") || strings.Contains(tipoProjeto, "pecu") {
		leitura.Alertas = append(leitura.Alertas, "Projeto envolve pecuária: conferir rebanho, vacinação, pastagem, alimentação e documentos sanitários quando aplicável.")
	}

	if r.PrecisaZARC == "sim" || strings.Contains(finalidade, "custeio agrícola") || strings.Contains(tipoProjeto, "lavoura") {
		leitura.Alertas = append(leitura.Alertas, "Projeto agrícola/custeio: conferir ZARC, cultura, município, safra, solo, cultivar e janela de plantio.")
	}

	// Próximos passos
	leitura.ProximosPassos = append(leitura.ProximosPassos, "Gerar ou revisar o checklist da reunião.")
	leitura.ProximosPassos = append(leitura.ProximosPassos, "Separar pendências por produtor, imóvel, banco, ambiental e técnico.")

	if len(leitura.Alertas) == 0 {
		leitura.Alertas = append(leitura.Alertas, "Nenhum alerta crítico foi identificado na pré-análise, mas os documentos ainda precisam ser conferidos.")
	}

	return leitura
}
