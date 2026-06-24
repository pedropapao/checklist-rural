package main

import "strings"

type FiltrosReuniao struct {
	Busca    string
	Banco    string
	Situacao string
}

func filtrarReunioes(reunioes []ReuniaoComResumo, filtros FiltrosReuniao) []ReuniaoComResumo {
	var filtradas []ReuniaoComResumo

	busca := strings.ToLower(strings.TrimSpace(filtros.Busca))
	bancoFiltro := strings.ToLower(strings.TrimSpace(filtros.Banco))
	situacaoFiltro := strings.ToLower(strings.TrimSpace(filtros.Situacao))

	for _, reuniao := range reunioes {
		if busca != "" {
			texto := strings.ToLower(
				reuniao.Produtor + " " +
					reuniao.Telefone + " " +
					reuniao.Municipio + " " +
					reuniao.UF + " " +
					reuniao.Banco + " " +
					reuniao.TipoProjeto + " " +
					reuniao.Atividade + " " +
					reuniao.Observacoes,
			)

			if !strings.Contains(texto, busca) {
				continue
			}
		}

		if bancoFiltro != "" {
			bancoReuniao := strings.ToLower(reuniao.Banco)

			switch bancoFiltro {
			case "banco do brasil":
				if !strings.Contains(bancoReuniao, "brasil") {
					continue
				}
			case "sicoob":
				if !strings.Contains(bancoReuniao, "sicoob") {
					continue
				}
			case "ainda não definido":
				if !strings.Contains(bancoReuniao, "ainda") {
					continue
				}
			}
		}

		if situacaoFiltro != "" {
			switch situacaoFiltro {
			case "pendentes":
				if reuniao.Resumo.Pendentes <= 0 {
					continue
				}

			case "concluidas":
				if reuniao.Resumo.Total == 0 || reuniao.Resumo.Pendentes > 0 {
					continue
				}

			case "andamento":
				if reuniao.Resumo.Total == 0 {
					continue
				}

				if reuniao.Resumo.Pendentes == 0 {
					continue
				}

				if reuniao.Resumo.PercentualConcluido == 0 {
					continue
				}
			}
		}

		filtradas = append(filtradas, reuniao)
	}

	return filtradas
}
