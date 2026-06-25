package main

type GrupoChecklistTela struct {
	Nome        string
	Etapa       string
	Descricao   string
	Itens       []ItemChecklist
	Total       int
	Pendentes   int
	Recebidos   int
	NaoSeAplica int
}

func agruparItensChecklistParaTela(itens []ItemChecklist) []GrupoChecklistTela {
	ordem := []GrupoChecklistTela{
		{
			Nome:      "Produtor e cadastro",
			Etapa:     "1. Produtor e cadastro",
			Descricao: "Documentos pessoais, identificação e dados básicos do produtor.",
		},
		{
			Nome:      "Enquadramento e banco",
			Etapa:     "2. Enquadramento e banco",
			Descricao: "Cadastro bancário, restrições, financiamento ativo, renda e classificação.",
		},
		{
			Nome:      "Imóvel rural",
			Etapa:     "3. Imóvel rural",
			Descricao: "Documentos da terra, posse, arrendamento, CAR, matrícula e vínculo com a área.",
		},
		{
			Nome:      "Ambiental",
			Etapa:     "4. Ambiental",
			Descricao: "Água, outorga, licenças, supressão vegetal e situações ambientais sensíveis.",
		},
		{
			Nome:      "Projeto técnico",
			Etapa:     "5. Projeto técnico",
			Descricao: "Projeto, orçamento, croqui, coordenadas, ZARC, ART e documentos técnicos.",
		},
		{
			Nome:      "Documentos específicos",
			Etapa:     "6. Documentos específicos",
			Descricao: "Itens que dependem do tipo do projeto: custeio, pecuária, investimento, máquinas ou obras.",
		},
		{
			Nome:      "Pendências finais",
			Etapa:     "7. Pendências finais",
			Descricao: "Itens que não entraram nas etapas anteriores e precisam de conferência final.",
		},
	}

	mapa := make(map[string]*GrupoChecklistTela)

	for i := range ordem {
		grupo := ordem[i]
		mapa[grupo.Nome] = &GrupoChecklistTela{
			Nome:      grupo.Nome,
			Etapa:     grupo.Etapa,
			Descricao: grupo.Descricao,
			Itens:     []ItemChecklist{},
		}
	}

	for _, item := range itens {
		nomeGrupo := classificarGrupoChecklistTela(item.Grupo, item.Item)

		grupo, existe := mapa[nomeGrupo]
		if !existe {
			grupo = mapa["Pendências finais"]
		}

		grupo.Itens = append(grupo.Itens, item)
		grupo.Total++

		switch item.Status {
		case "Recebido":
			grupo.Recebidos++
		case "Não se aplica":
			grupo.NaoSeAplica++
		default:
			grupo.Pendentes++
		}
	}

	var resultado []GrupoChecklistTela

	for _, modelo := range ordem {
		grupo := mapa[modelo.Nome]
		if grupo == nil || grupo.Total == 0 {
			continue
		}

		resultado = append(resultado, *grupo)
	}

	return resultado
}

func classificarGrupoChecklistTela(grupo string, item string) string {
	texto := grupo + " " + item

	if contemAlgum(texto, []string{
		"CPF", "RG", "CNH", "documento pessoal", "estado civil",
		"comprovante de endereço", "telefone", "produtor",
	}) {
		return "Produtor e cadastro"
	}

	if contemAlgum(texto, []string{
		"banco", "cadastro", "restrição", "restricao", "financiamento ativo",
		"renda", "RBA", "Pronaf", "Pronamp", "classificação", "classificacao",
	}) {
		return "Enquadramento e banco"
	}

	if contemAlgum(texto, []string{
		"matrícula", "matricula", "CCIR", "ITR", "CAR", "arrendamento",
		"comodato", "parceria", "imóvel", "imovel", "propriedade", "terra",
	}) {
		return "Imóvel rural"
	}

	if contemAlgum(texto, []string{
		"ambiental", "licença", "licenca", "outorga", "água", "agua",
		"irrigação", "irrigacao", "supressão", "supressao", "vegetal",
		"APP", "reserva legal",
	}) {
		return "Ambiental"
	}

	if contemAlgum(texto, []string{
		"projeto técnico", "projeto tecnico", "orçamento", "orcamento",
		"croqui", "mapa", "coordenadas", "ZARC", "assistência técnica",
		"assistencia tecnica", "ART", "plano", "laudo",
	}) {
		return "Projeto técnico"
	}

	if contemAlgum(texto, []string{
		"custeio", "pecuária", "pecuaria", "investimento", "máquina",
		"maquina", "equipamento", "obra", "benfeitoria", "nota fiscal",
		"proposta", "fornecedor", "rebanho", "pastagem",
	}) {
		return "Documentos específicos"
	}

	return "Pendências finais"
}

func contemAlgum(texto string, palavras []string) bool {
	for _, palavra := range palavras {
		if contemTexto(texto, palavra) {
			return true
		}
	}

	return false
}

func contemTexto(texto string, palavra string) bool {
	texto = normalizarTextoSimples(texto)
	palavra = normalizarTextoSimples(palavra)

	if palavra == "" || len(texto) < len(palavra) {
		return false
	}

	return procurarTexto(texto, palavra)
}

func procurarTexto(texto string, palavra string) bool {
	for i := 0; i <= len(texto)-len(palavra); i++ {
		if texto[i:i+len(palavra)] == palavra {
			return true
		}
	}

	return false
}

func normalizarTextoSimples(texto string) string {
	novo := ""

	for _, letra := range texto {
		switch letra {
		case 'Á', 'À', 'Â', 'Ã', 'Ä', 'á', 'à', 'â', 'ã', 'ä':
			novo += "a"
		case 'É', 'È', 'Ê', 'Ë', 'é', 'è', 'ê', 'ë':
			novo += "e"
		case 'Í', 'Ì', 'Î', 'Ï', 'í', 'ì', 'î', 'ï':
			novo += "i"
		case 'Ó', 'Ò', 'Ô', 'Õ', 'Ö', 'ó', 'ò', 'ô', 'õ', 'ö':
			novo += "o"
		case 'Ú', 'Ù', 'Û', 'Ü', 'ú', 'ù', 'û', 'ü':
			novo += "u"
		case 'Ç', 'ç':
			novo += "c"
		default:
			if letra >= 'A' && letra <= 'Z' {
				novo += string(letra + 32)
			} else {
				novo += string(letra)
			}
		}
	}

	return novo
}
