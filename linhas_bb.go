package main

import "strings"

type LinhaCreditoSugerida struct {
	Nome    string
	Motivo  string
	Atencao string
}

func sugerirLinhasBB(r Reuniao) []LinhaCreditoSugerida {
	banco := strings.ToLower(r.Banco)
	if !strings.Contains(banco, "brasil") && !strings.Contains(banco, "bb") {
		return nil
	}

	finalidade := strings.ToLower(r.FinalidadeCredito + " " + r.TipoProjeto + " " + r.Atividade)

	ehCusteio := strings.Contains(finalidade, "custeio")
	ehCusteioAgricola := strings.Contains(finalidade, "custeio agrícola") ||
		strings.Contains(finalidade, "custeio agricola") ||
		strings.Contains(finalidade, "lavoura") ||
		strings.Contains(finalidade, "safra")

	ehCusteioPecuario := strings.Contains(finalidade, "custeio pecuário") ||
		strings.Contains(finalidade, "custeio pecuario") ||
		strings.Contains(finalidade, "pecuária") ||
		strings.Contains(finalidade, "pecuaria") ||
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
		strings.Contains(finalidade, "trator")

	linhas := []LinhaCreditoSugerida{}

	if r.PossuiCAF == "sim" && ehCusteio {
		linhas = append(linhas, LinhaCreditoSugerida{
			Nome:    "Pronaf Custeio",
			Motivo:  "O produtor informou possuir CAF e a finalidade indicada é custeio da atividade rural.",
			Atencao: "Confirmar CAF ativo, enquadramento no Pronaf, limite disponível, cultura/atividade, orçamento e exigências da agência.",
		})
	}

	if r.PossuiCAF == "sim" && ehInvestimento {
		linhas = append(linhas, LinhaCreditoSugerida{
			Nome:    "Pronaf Mais Alimentos",
			Motivo:  "O produtor informou possuir CAF e a demanda parece ser investimento, estrutura, máquina, equipamento ou melhoria produtiva.",
			Atencao: "Confirmar enquadramento no Pronaf, finalidade financiável, orçamento/proposta, capacidade de pagamento e limite aplicável.",
		})
	}

	if r.ClassificacaoProdutor == "Médio produtor" && ehCusteio {
		linhas = append(linhas, LinhaCreditoSugerida{
			Nome:    "Pronamp Custeio",
			Motivo:  "A RBA classificou o produtor como médio produtor e a finalidade indicada é custeio.",
			Atencao: "Confirmar renda, enquadramento, limite, cadastro, atividade financiada e orçamento de custeio.",
		})
	}

	if r.ClassificacaoProdutor == "Médio produtor" && ehInvestimento {
		linhas = append(linhas, LinhaCreditoSugerida{
			Nome:    "Pronamp Investimento",
			Motivo:  "A RBA classificou o produtor como médio produtor e a finalidade indicada é investimento ou melhoria da atividade rural.",
			Atencao: "Confirmar orçamento, proposta, garantias, capacidade de pagamento, prazo e política vigente do banco.",
		})
	}

	if ehMaquina && r.ClassificacaoProdutor != "Pequeno produtor" {
		linhas = append(linhas, LinhaCreditoSugerida{
			Nome:    "Moderfrota ou linha de investimento para máquinas",
			Motivo:  "A demanda envolve máquina, equipamento, implemento ou trator.",
			Atencao: "Confirmar se a máquina é nova ou usada, nota/proposta, fabricante, ano, capacidade de pagamento e linha aceita pela agência.",
		})
	}

	if len(linhas) == 0 && (ehCusteioAgricola || ehCusteioPecuario || ehCusteio) {
		linhas = append(linhas, LinhaCreditoSugerida{
			Nome:    "Custeio Agropecuário BB",
			Motivo:  "A finalidade indicada é custeio da produção rural.",
			Atencao: "Confirmar atividade, orçamento, safra/ciclo, garantias, limite cadastral e política de crédito do banco.",
		})
	}

	if len(linhas) == 0 && ehInvestimento {
		linhas = append(linhas, LinhaCreditoSugerida{
			Nome:    "Linha de investimento agro a confirmar",
			Motivo:  "A demanda parece ser investimento, estrutura, máquina, equipamento, obra ou irrigação.",
			Atencao: "Confirmar no banco qual linha se ajusta melhor à finalidade, ao enquadramento do produtor e ao item financiado.",
		})
	}

	if len(linhas) == 0 {
		linhas = append(linhas, LinhaCreditoSugerida{
			Nome:    "Linha BB a confirmar com a agência",
			Motivo:  "As informações da pré-análise ainda não permitem indicar um caminho principal com segurança.",
			Atencao: "Completar finalidade, RBA, CAF, orçamento, situação da terra e cadastro antes de definir a linha.",
		})
	}

	return linhas
}
