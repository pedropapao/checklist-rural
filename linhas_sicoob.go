package main

import "strings"

type LinhaSicoobSugerida struct {
	Nome    string
	Motivo  string
	Atencao string
}

func sugerirLinhasSicoob(r Reuniao) []LinhaSicoobSugerida {
	banco := strings.ToLower(r.Banco)
	if !strings.Contains(banco, "sicoob") &&
		!strings.Contains(banco, "cooperativa") &&
		!strings.Contains(banco, "coop") {
		return nil
	}

	finalidade := strings.ToLower(r.FinalidadeCredito + " " + r.TipoProjeto + " " + r.Atividade)

	ehCusteio := strings.Contains(finalidade, "custeio") ||
		strings.Contains(finalidade, "safra") ||
		strings.Contains(finalidade, "lavoura")

	ehPecuario := strings.Contains(finalidade, "pecu") ||
		strings.Contains(finalidade, "gado") ||
		strings.Contains(finalidade, "leite") ||
		strings.Contains(finalidade, "rebanho")

	ehInvestimento := strings.Contains(finalidade, "investimento") ||
		strings.Contains(finalidade, "máquina") ||
		strings.Contains(finalidade, "maquina") ||
		strings.Contains(finalidade, "equipamento") ||
		strings.Contains(finalidade, "obra") ||
		strings.Contains(finalidade, "benfeitoria") ||
		strings.Contains(finalidade, "irrigação") ||
		strings.Contains(finalidade, "irrigacao")

	ehMaquina := strings.Contains(finalidade, "máquina") ||
		strings.Contains(finalidade, "maquina") ||
		strings.Contains(finalidade, "equipamento") ||
		strings.Contains(finalidade, "implemento") ||
		strings.Contains(finalidade, "trator") ||
		strings.Contains(finalidade, "colheitadeira") ||
		strings.Contains(finalidade, "plantadeira") ||
		strings.Contains(finalidade, "pulverizador")

	ehIrrigacao := strings.Contains(finalidade, "irrigação") ||
		strings.Contains(finalidade, "irrigacao") ||
		r.UsaAgua == "sim"

	ehArmazenagem := strings.Contains(finalidade, "armaz") ||
		strings.Contains(finalidade, "silo") ||
		strings.Contains(finalidade, "armazenagem")

	ehIndustrializacao := strings.Contains(finalidade, "industrial") ||
		strings.Contains(finalidade, "beneficiamento") ||
		strings.Contains(finalidade, "processamento") ||
		strings.Contains(finalidade, "agroindústria") ||
		strings.Contains(finalidade, "agroindustria")

	linhas := []LinhaSicoobSugerida{}

	if r.PossuiCAF == "sim" && ehCusteio {
		linhas = append(linhas, LinhaSicoobSugerida{
			Nome:    "Pronaf Custeio Sicoob",
			Motivo:  "O produtor informou possuir CAF e a finalidade parece ser custeio agrícola ou pecuário.",
			Atencao: "Confirmar CAF/DAP válido, enquadramento familiar, atividade, orçamento e condições com a cooperativa.",
		})
	}

	if r.PossuiCAF == "sim" && ehInvestimento {
		linhas = append(linhas, LinhaSicoobSugerida{
			Nome:    "Pronaf Investimento / Mais Alimentos",
			Motivo:  "O produtor informou possuir CAF e a demanda parece ser investimento na estrutura produtiva.",
			Atencao: "Confirmar linha aplicável, orçamento, proposta, finalidade financiável, capacidade de pagamento e regras da cooperativa.",
		})
	}

	if r.ClassificacaoProdutor == "Médio produtor" && ehCusteio {
		linhas = append(linhas, LinhaSicoobSugerida{
			Nome:    "Pronamp Custeio Sicoob",
			Motivo:  "A RBA classificou o produtor como médio produtor e a finalidade indicada é custeio.",
			Atencao: "Confirmar enquadramento, RBA, orçamento de custeio, limite e condições de contratação com a cooperativa.",
		})
	}

	if r.ClassificacaoProdutor == "Médio produtor" && ehInvestimento {
		linhas = append(linhas, LinhaSicoobSugerida{
			Nome:    "Pronamp Investimento Sicoob",
			Motivo:  "A RBA classificou o produtor como médio produtor e a finalidade parece ser investimento.",
			Atencao: "Confirmar orçamento, proposta, garantias, capacidade de pagamento, prazo e disponibilidade da linha na cooperativa.",
		})
	}

	if ehMaquina {
		linhas = append(linhas, LinhaSicoobSugerida{
			Nome:    "Moderfrota / investimento em máquinas e equipamentos",
			Motivo:  "A demanda envolve máquina, equipamento, implemento agrícola, trator, colheitadeira ou item semelhante.",
			Atencao: "Confirmar se o bem é novo ou usado, proposta comercial, fabricante, ano/modelo, valor, garantia e linha aceita pela cooperativa.",
		})
	}

	if ehIrrigacao {
		linhas = append(linhas, LinhaSicoobSugerida{
			Nome:    "Proirriga / investimento em irrigação",
			Motivo:  "A demanda envolve irrigação, uso de água ou estrutura hídrica.",
			Atencao: "Conferir projeto técnico, orçamento, outorga ou dispensa de uso da água e regularidade ambiental.",
		})
	}

	if ehArmazenagem {
		linhas = append(linhas, LinhaSicoobSugerida{
			Nome:    "PCA / armazenagem",
			Motivo:  "A demanda envolve construção, ampliação, modernização ou reforma de estrutura de armazenagem.",
			Atencao: "Conferir projeto, orçamento, memorial, licenças, capacidade de armazenagem e garantias.",
		})
	}

	if ehIndustrializacao {
		linhas = append(linhas, LinhaSicoobSugerida{
			Nome:    "Industrialização agropecuária",
			Motivo:  "A demanda envolve beneficiamento, processamento, armazenamento ou comercialização de produto agropecuário.",
			Atencao: "Confirmar origem da produção, itens financiáveis, orçamento, ciclo de comercialização e regras da cooperativa.",
		})
	}

	if len(linhas) == 0 && (ehCusteio || ehPecuario) {
		linhas = append(linhas, LinhaSicoobSugerida{
			Nome:    "Custeio agropecuário Sicoob",
			Motivo:  "A finalidade indicada é custear a produção rural, agrícola ou pecuária.",
			Atencao: "Confirmar atividade, safra/ciclo, orçamento, garantias, cadastro, limite e disponibilidade na cooperativa.",
		})
	}

	if len(linhas) == 0 && ehInvestimento {
		linhas = append(linhas, LinhaSicoobSugerida{
			Nome:    "Investimento rural Sicoob",
			Motivo:  "A finalidade indicada envolve investimento, estrutura, benfeitoria, máquina, equipamento ou melhoria produtiva.",
			Atencao: "Confirmar qual programa se aplica melhor: Pronaf, Pronamp, BNDES, RPL, FCO ou outra linha disponível.",
		})
	}

	if len(linhas) == 0 {
		linhas = append(linhas, LinhaSicoobSugerida{
			Nome:    "Caminho Sicoob a confirmar com a cooperativa",
			Motivo:  "As informações atuais ainda não permitem indicar um caminho principal com segurança.",
			Atencao: "Completar finalidade, RBA, CAF, orçamento, situação da terra, cadastro e documentos técnicos antes de definir a linha.",
		})
	}

	return linhas
}
