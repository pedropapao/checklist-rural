package main

import "strings"

func adicionarItensPorLinhaSicoob(itens []ItemChecklist, r Reuniao) []ItemChecklist {
	linhas := sugerirLinhasSicoob(r)

	for _, linha := range linhas {
		nomeLinha := strings.ToLower(linha.Nome)

		if strings.Contains(nomeLinha, "pronaf") {
			itens = adicionarItemSeNaoExiste(itens, ItemChecklist{
				Grupo:      "Enquadramento da linha de crédito",
				Item:       "CAF/DAP válido e enquadramento do produtor familiar",
				Status:     "Pendente",
				Explicacao: "Para operação Pronaf no Sicoob, conferir CAF/DAP válido, enquadramento familiar, atividade rural e documentos exigidos pela cooperativa.",
			})
		}

		if strings.Contains(nomeLinha, "pronamp") {
			itens = adicionarItemSeNaoExiste(itens, ItemChecklist{
				Grupo:      "Enquadramento da linha de crédito",
				Item:       "RBA e enquadramento do produtor como médio produtor",
				Status:     "Pendente",
				Explicacao: "Para Pronamp no Sicoob, conferir Receita Bruta Agropecuária Anual, documentos de renda e enquadramento como médio produtor.",
			})
		}

		if strings.Contains(nomeLinha, "custeio") {
			itens = adicionarItemSeNaoExiste(itens, ItemChecklist{
				Grupo:      "Custeio rural",
				Item:       "Plano/orçamento de custeio agropecuário para apresentar à cooperativa",
				Status:     "Pendente",
				Explicacao: "No custeio, organizar cultura ou atividade, safra/ciclo, área, insumos, quantidades, valores, produtividade esperada e capacidade de pagamento.",
			})
		}

		if strings.Contains(nomeLinha, "máquina") ||
			strings.Contains(nomeLinha, "maquina") ||
			strings.Contains(nomeLinha, "moderfrota") {
			itens = adicionarItemSeNaoExiste(itens, ItemChecklist{
				Grupo:      "Máquinas e equipamentos",
				Item:       "Proposta comercial da máquina/equipamento com especificações completas",
				Status:     "Pendente",
				Explicacao: "Para máquinas e equipamentos, reunir proposta comercial, marca, modelo, ano, fornecedor, valor, condição de novo/usado e finalidade na atividade rural.",
			})
		}

		if strings.Contains(nomeLinha, "irrig") || strings.Contains(nomeLinha, "proirriga") {
			itens = adicionarItemSeNaoExiste(itens, ItemChecklist{
				Grupo:      "Ambiental e uso da água",
				Item:       "Projeto de irrigação e regularidade do uso da água",
				Status:     "Pendente",
				Explicacao: "Para irrigação, conferir projeto técnico, orçamento, fonte de água, outorga ou dispensa, licenças e compatibilidade ambiental.",
			})
		}

		if strings.Contains(nomeLinha, "armazen") || strings.Contains(nomeLinha, "pca") {
			itens = adicionarItemSeNaoExiste(itens, ItemChecklist{
				Grupo:      "Investimento rural",
				Item:       "Projeto e orçamento da estrutura de armazenagem",
				Status:     "Pendente",
				Explicacao: "Para armazenagem, organizar projeto, memorial, capacidade, orçamento, licenças, matrícula/CAR e garantias.",
			})
		}

		if strings.Contains(nomeLinha, "industrial") {
			itens = adicionarItemSeNaoExiste(itens, ItemChecklist{
				Grupo:      "Industrialização rural",
				Item:       "Dados do beneficiamento, processamento ou agroindústria",
				Status:     "Pendente",
				Explicacao: "Para industrialização, levantar origem da produção, produto processado, orçamento, fluxo de produção, ciclo de comercialização e documentos sanitários/ambientais aplicáveis.",
			})
		}
	}

	return itens
}
