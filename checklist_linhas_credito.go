package main

import "strings"

func adicionarItensPorLinhaCredito(itens []ItemChecklist, reuniao Reuniao) []ItemChecklist {
	linhas := sugerirLinhasBB(reuniao)

	for _, linha := range linhas {
		nomeLinha := strings.ToLower(linha.Nome)

		if strings.Contains(nomeLinha, "pronaf") {
			itens = adicionarItemSeNaoExiste(itens, ItemChecklist{
				Grupo:      "Enquadramento da linha de crédito",
				Item:       "CAF ativo e enquadramento do produtor na agricultura familiar",
				Status:     "Pendente",
				Explicacao: "Para operações Pronaf, conferir se o produtor possui CAF ativo e se o enquadramento familiar está compatível com a finalidade do crédito.",
			})
		}

		if strings.Contains(nomeLinha, "pronamp") {
			itens = adicionarItemSeNaoExiste(itens, ItemChecklist{
				Grupo:      "Enquadramento da linha de crédito",
				Item:       "Comprovação da RBA e enquadramento como médio produtor",
				Status:     "Pendente",
				Explicacao: "Para operações Pronamp, conferir a Receita Bruta Agropecuária Anual e os documentos que demonstram o enquadramento do produtor como médio produtor.",
			})
		}

		if strings.Contains(nomeLinha, "custeio") {
			itens = adicionarItemSeNaoExiste(itens, ItemChecklist{
				Grupo:      "Custeio rural",
				Item:       "Orçamento de custeio com cultura, safra, área e insumos",
				Status:     "Pendente",
				Explicacao: "No custeio, o banco precisa entender o que será financiado: cultura ou atividade, safra/ciclo, área, insumos, quantidades e valores.",
			})

			itens = adicionarItemSeNaoExiste(itens, ItemChecklist{
				Grupo:      "Custeio rural",
				Item:       "Conferência de ZARC quando for lavoura sujeita ao zoneamento",
				Status:     "Pendente",
				Explicacao: "Para custeio agrícola, conferir município, cultura, safra, tipo de solo, cultivar e janela de plantio quando houver exigência de ZARC.",
			})
		}

		if strings.Contains(nomeLinha, "investimento") ||
			strings.Contains(nomeLinha, "mais alimentos") ||
			strings.Contains(nomeLinha, "moderfrota") ||
			strings.Contains(nomeLinha, "máquina") ||
			strings.Contains(nomeLinha, "maquina") {

			itens = adicionarItemSeNaoExiste(itens, ItemChecklist{
				Grupo:      "Investimento rural",
				Item:       "Orçamento ou proposta comercial do bem ou serviço financiado",
				Status:     "Pendente",
				Explicacao: "Em investimento, normalmente é necessário apresentar orçamento, proposta comercial, descrição do item, valor, fornecedor e condições de pagamento.",
			})

			itens = adicionarItemSeNaoExiste(itens, ItemChecklist{
				Grupo:      "Investimento rural",
				Item:       "Análise da capacidade de pagamento do investimento",
				Status:     "Pendente",
				Explicacao: "Antes de avançar, conferir se a atividade rural gera receita suficiente para pagar o financiamento dentro do prazo pretendido.",
			})
		}
	}

	return itens
}

func adicionarItemSeNaoExiste(itens []ItemChecklist, novo ItemChecklist) []ItemChecklist {
	for _, item := range itens {
		if strings.EqualFold(strings.TrimSpace(item.Item), strings.TrimSpace(novo.Item)) {
			return itens
		}
	}

	return append(itens, novo)
}
