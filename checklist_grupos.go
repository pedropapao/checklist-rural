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
			Nome:      "Documentos pessoais e cadastrais",
			Etapa:     "1. Documentos pessoais e cadastrais",
			Descricao: "Identificação do produtor, estado civil, endereço, imposto de renda e patrimônio.",
		},
		{
			Nome:      "Cadastro bancário",
			Etapa:     "2. Cadastro bancário",
			Descricao: "Cadastro do produtor no banco ou cooperativa, atualização cadastral e pendências bancárias.",
		},
		{
			Nome:      "Documentos fundiários",
			Etapa:     "3. Documentos fundiários",
			Descricao: "Documentos da terra, matrícula, CCIR, ITR, arrendamento, parceria, comodato e anuências.",
		},
		{
			Nome:      "Documentos ambientais",
			Etapa:     "4. Documentos ambientais",
			Descricao: "CAR, embargos, APP, Reserva Legal, unidade de conservação, terra indígena ou quilombola.",
		},
		{
			Nome:      "Água e irrigação",
			Etapa:     "5. Água e irrigação",
			Descricao: "Outorga, dispensa, licenciamento e projeto técnico quando houver uso de água ou irrigação.",
		},
		{
			Nome:      "Supressão vegetal",
			Etapa:     "6. Supressão vegetal",
			Descricao: "Autorização e conferência quando houver retirada de vegetação ou abertura de área.",
		},
		{
			Nome:      "Documentos técnicos gerais",
			Etapa:     "7. Documentos técnicos gerais",
			Descricao: "Projeto, orçamento, croqui, coordenadas, ART/TRT, assinaturas e contrato técnico.",
		},
		{
			Nome:      "Custeio agrícola",
			Etapa:     "8. Custeio agrícola",
			Descricao: "Itens ligados à lavoura, área, talhão, plantio, colheita, insumos, ZARC, Proagro ou seguro.",
		},
		{
			Nome:      "Pecuária",
			Etapa:     "9. Pecuária",
			Descricao: "Rebanho, vacinação, capacidade de suporte, pastagens, manejo e projeção zootécnica.",
		},
		{
			Nome:      "Investimento",
			Etapa:     "10. Investimento",
			Descricao: "Objetivo, justificativa, propostas comerciais, cronogramas e capacidade de pagamento.",
		},
		{
			Nome:      "Máquinas e equipamentos",
			Etapa:     "11. Máquinas e equipamentos",
			Descricao: "Tipo, marca, modelo, ano, proposta comercial, fornecedor e condições de entrega.",
		},
		{
			Nome:      "Obras e benfeitorias",
			Etapa:     "12. Obras e benfeitorias",
			Descricao: "Descrição da obra, local, projeto, memorial, orçamento, fotos e coordenadas.",
		},
		{
			Nome:      "Banco do Brasil",
			Etapa:     "13. Banco do Brasil",
			Descricao: "Documentos e conferências específicas para operações no Banco do Brasil.",
		},
		{
			Nome:      "Sicoob",
			Etapa:     "14. Sicoob",
			Descricao: "Planilhas, PUC, carta de internalização, cronogramas e documentos próprios do Sicoob.",
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
		item.Explicacao = explicacaoItemChecklist(item.Item)

		grupo, existe := mapa[item.Grupo]
		if !existe {
			continue
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
		if grupo != nil && grupo.Total > 0 {
			resultado = append(resultado, *grupo)
		}
	}

	return resultado
}

func explicacaoItemChecklist(item string) string {
	explicacoes := map[string]string{
		"RG ou documento oficial com foto":     "Serve para comprovar quem é o produtor. Pode ser RG, CNH ou outro documento oficial aceito pelo banco.",
		"CPF":                                  "Identifica o produtor perante banco, Receita Federal e sistemas de consulta. É usado para cadastro e análise de crédito.",
		"Certidão de estado civil atualizada":  "Mostra se o produtor é solteiro, casado, divorciado ou viúvo. Pode indicar necessidade de assinatura do cônjuge.",
		"Comprovante de residência atualizado": "Serve para atualizar o cadastro do produtor no banco e confirmar o endereço de contato.",
		"Declaração de IRPF completa":          "Mostra renda, bens, atividade e situação patrimonial do produtor. Ajuda na análise da capacidade de pagamento.",
		"Recibo de entrega do IRPF":            "Comprova que a declaração de imposto de renda foi enviada à Receita Federal.",
		"Relação de bens e direitos assinada":  "Lista o patrimônio do produtor, como imóveis, veículos, máquinas, rebanho e outros bens.",

		"Conferir cadastro já existente no banco":          "Usado quando o produtor já é cliente. Serve para verificar se o cadastro está ativo, atualizado e sem pendências.",
		"Fazer ou atualizar cadastro do produtor no banco": "Usado quando o produtor ainda não tem cadastro ou está desatualizado. Sem cadastro correto, o banco não avança.",

		"Matrícula do imóvel atualizada":                     "Documento principal do cartório. Mostra quem é o dono da terra, área, localização e situação jurídica.",
		"Certidão de inteiro teor":                           "É uma cópia completa da matrícula, com todos os registros e averbações. Ajuda a entender o histórico do imóvel.",
		"Certidão de ônus reais":                             "Mostra se o imóvel possui dívida, hipoteca, penhora, alienação ou outra restrição registrada.",
		"Certidão de ações reipersecutórias":                 "Mostra se existe ação judicial que possa afetar a propriedade do imóvel.",
		"Conferir validade da matrícula":                     "Alguns bancos exigem matrícula recente. Documento antigo pode ser recusado na análise.",
		"CCIR / Incra do último exercício":                   "Documento ligado ao cadastro rural no Incra. Ajuda a comprovar a identificação cadastral do imóvel.",
		"Comprovante de quitação do CCIR":                    "Mostra que o CCIR foi pago. Sem quitação, pode virar pendência documental.",
		"ITR do último exercício":                            "É o Imposto Territorial Rural. O banco usa para conferir regularidade fiscal do imóvel.",
		"DARF do ITR quitado":                                "Comprova que o imposto territorial rural foi pago.",
		"Contrato de arrendamento/parceria/comodato vigente": "Usado quando o produtor não é dono da terra. Comprova que ele tem direito de usar a área.",
		"Contrato cobrindo todo o prazo do financiamento":    "O prazo do contrato da terra precisa cobrir o período do financiamento.",
		"Carta de anuência do proprietário":                  "É a autorização do dono da terra permitindo que o produtor faça o projeto naquela área.",
		"Autorização para penhor e execução do projeto":      "Serve para o proprietário autorizar o uso da produção, bens ou atividade financiada como garantia, quando necessário.",

		"Recibo do CAR ativo":                          "Comprova que o imóvel está inscrito no Cadastro Ambiental Rural.",
		"Conferir situação do CAR":                     "Não basta ter CAR. É preciso verificar pendências, sobreposição, cancelamento ou inconsistência.",
		"Providenciar ou confirmar CAR da propriedade": "Usado quando o produtor não sabe se tem CAR ou ainda não apresentou o documento.",
		"Verificar embargos ambientais":                "Serve para saber se a área ou produtor possui impedimento ambiental. Embargo pode travar o financiamento.",
		"Verificar APP e Reserva Legal":                "Confere se o projeto não está usando área protegida de forma irregular.",
		"Verificar unidade de conservação":             "Serve para saber se a área está dentro ou próxima de unidade de conservação.",
		"Verificar terra indígena ou área quilombola":  "Confere se existe sobreposição ou conflito com áreas protegidas por legislação específica.",

		"Outorga de uso de água":                        "Documento que autoriza o uso de água de rio, córrego, poço, represa ou outro recurso hídrico.",
		"Certidão de dispensa de outorga, se aplicável": "Quando a outorga não é obrigatória, o órgão ambiental pode emitir uma dispensa.",
		"Licenciamento ou inexigibilidade ambiental":    "Mostra se a atividade precisa de licença ambiental ou se está dispensada.",
		"Projeto técnico do sistema de irrigação":       "Explica como será o sistema de irrigação: bomba, tubos, vazão, área irrigada, fonte de água e funcionamento.",
		"Orçamento do sistema de irrigação":             "Mostra o custo dos equipamentos e serviços de irrigação.",

		"ASV - Autorização de Supressão Vegetal":                "É a autorização para retirar vegetação nativa. Sem ela, não se deve financiar abertura de área.",
		"Conferir compatibilidade com CAR, APP e Reserva Legal": "Verifica se a supressão pretendida não atinge área protegida ou irregular.",

		"Projeto técnico":                            "Documento que explica o financiamento: objetivo, valores, área, atividade, produção e forma de pagamento.",
		"Orçamento técnico":                          "Mostra quanto o projeto vai custar. Pode envolver insumos, máquinas, obras, serviços ou equipamentos.",
		"Memorial descritivo, se aplicável":          "Explica detalhes técnicos de uma obra, benfeitoria, sistema ou investimento.",
		"Croqui de localização":                      "Desenho ou mapa simples mostrando onde fica a área do projeto dentro do imóvel.",
		"Coordenadas geográficas em graus decimais":  "Localizam a área financiada com precisão. Usadas em custeio, investimento, fiscalização e banco.",
		"ART ou TRT":                                 "Documento de responsabilidade técnica. Mostra que um profissional habilitado assumiu o projeto.",
		"Comprovante de pagamento da ART/TRT":        "Mostra que a responsabilidade técnica foi emitida e paga corretamente.",
		"Assinatura do técnico":                      "Comprova que o profissional concorda e responde pelas informações técnicas.",
		"Assinatura do produtor":                     "Comprova que o produtor conferiu e concorda com o projeto apresentado.",
		"Contrato de prestação de serviços técnicos": "Formaliza o serviço entre produtor e técnico ou projetista.",

		"Cultura financiada":                "Identifica o que será plantado: soja, milho, café, mandioca, feijão ou outra cultura.",
		"Área plantada":                     "Mostra quantos hectares serão financiados.",
		"Talhão / gleba":                    "Identifica exatamente qual parte da propriedade será usada.",
		"Variedade / cultivar":              "Mostra qual semente ou material genético será usado.",
		"Data prevista de plantio":          "Indica quando a lavoura será implantada.",
		"Data prevista de colheita":         "Ajuda a calcular prazo de pagamento e receita esperada.",
		"Análise de solo":                   "Mostra a condição do solo e ajuda a justificar adubação e correção.",
		"Recomendação agronômica":           "Orientação técnica de calagem, adubação, sementes, defensivos e manejo.",
		"Insumos previstos":                 "Lista sementes, fertilizantes, defensivos, calcário, combustível e demais custos previstos.",
		"Validar ZARC":                      "Confere se cultura, município, solo e época de plantio estão dentro do zoneamento agrícola.",
		"Conferir município dentro do ZARC": "Verifica se aquela cultura tem zoneamento permitido no município.",
		"Conferir janela de plantio":        "Verifica se a data de plantio está dentro do período recomendado.",
		"Conferir Proagro ou seguro rural":  "Confere se a operação terá cobertura de risco por Proagro ou seguro rural.",

		"Atividade pecuária definida": "Especifica se é corte, leite, cria, recria, engorda, ovinos, suínos ou outra atividade.",
		"Ficha sanitária do rebanho":  "Documento que mostra informações sanitárias dos animais.",
		"Declaração de vacinação":     "Comprova que o rebanho está com vacinação exigida em dia.",
		"Quantidade atual de animais": "Mostra o tamanho atual do rebanho.",
		"Categorias do rebanho":       "Divide os animais por tipo: matrizes, bezerros, novilhas, touros, garrotes etc.",
		"Capacidade de suporte":       "Mostra se a área tem pastagem suficiente para manter o rebanho.",
		"Pastagens disponíveis":       "Identifica a área e condição das pastagens.",
		"Manejo previsto":             "Explica como o produtor vai alimentar, cuidar e conduzir o rebanho.",
		"Evolução de rebanho":         "Mostra como o rebanho deve crescer ou mudar ao longo do tempo.",
		"Projeção zootécnica":         "Estimativa técnica de produção, nascimento, venda, mortalidade, ganho de peso ou produção de leite.",

		"Objetivo do investimento":                        "Explica o que o produtor quer comprar, construir ou melhorar.",
		"Justificativa técnica":                           "Explica por que aquele investimento é necessário e como melhora a atividade.",
		"Mínimo de 3 propostas comerciais com CNPJ":       "Serve para comparar preços e comprovar que o valor financiado está dentro do mercado.",
		"Validade das propostas comerciais":               "Confere se os orçamentos ainda estão válidos.",
		"Cronograma físico-financeiro":                    "Mostra as etapas da execução e quanto será gasto em cada fase.",
		"Cronograma de desembolso":                        "Mostra quando o banco deve liberar cada parte do dinheiro.",
		"Comprovação de capacidade de pagamento":          "Mostra que a atividade do produtor consegue pagar o financiamento.",
		"Compatibilidade do investimento com a atividade": "Verifica se o bem financiado faz sentido para a produção do produtor.",

		"Tipo de máquina/equipamento": "Define o que será comprado: trator, carreta, pulverizador, ordenhadeira ou outro equipamento.",
		"Marca":                       "Identifica o fabricante.",
		"Modelo":                      "Identifica a versão exata do equipamento.",
		"Ano, se usado":               "Importante para avaliar valor, vida útil e aceitação pelo banco.",
		"Novo ou usado":               "O banco pode ter regras diferentes para bem novo e usado.",
		"Proposta comercial":          "Documento do fornecedor com preço, descrição e condições de venda.",
		"Fornecedor com CNPJ":         "Mostra que quem está vendendo é uma empresa formalizada.",
		"Condições de entrega":        "Explica prazo, frete, montagem, garantia ou entrega técnica.",

		"Descrição da obra":                  "Explica o que será construído ou reformado.",
		"Local da obra":                      "Mostra onde a obra será feita dentro do imóvel.",
		"Projeto técnico da construção":      "Documento técnico com detalhes da construção.",
		"Memorial descritivo":                "Explica materiais, medidas, método de execução e padrão da obra.",
		"Orçamento de materiais":             "Lista os materiais e seus custos.",
		"Orçamento de mão de obra":           "Mostra o custo dos serviços para executar a obra.",
		"Licença ou autorização, se exigida": "Algumas obras precisam de autorização ambiental, municipal ou outra.",
		"Fotos do local":                     "Comprovam a situação antes da obra.",
		"Coordenadas do local":               "Localizam exatamente onde a benfeitoria será feita.",

		"Ficha cadastral Banco do Brasil":                     "Documento usado para cadastro ou atualização do cliente no Banco do Brasil.",
		"Projeto/orçamento emitido no sistema oficial BB":     "Algumas operações precisam seguir modelo ou sistema próprio do Banco do Brasil.",
		"Projeto/orçamento assinado":                          "Comprova concordância do produtor e responsabilidade do técnico.",
		"Croqui com coordenadas em graus decimais":            "Ajuda o banco a localizar a área financiada.",
		"Matrícula com validade de 30 dias":                   "Matrícula recente para análise da situação atual do imóvel.",
		"Certidões de ônus e ações":                           "Verificam restrições, dívidas, disputas ou riscos sobre o imóvel.",
		"CCIR quitado":                                        "Comprova regularidade cadastral rural.",
		"ITR e DARF quitados":                                 "Comprova regularidade do imposto rural.",
		"Anuência do proprietário, se aplicável":              "Usado quando o produtor não é dono da terra.",
		"3 propostas comerciais para investimento":            "Compara preços e evita valor fora de mercado.",
		"Cronograma físico-financeiro para investimento fixo": "Mostra etapas e desembolso do investimento.",
		"Ficha sanitária ou vacinação para pecuária":          "Comprova regularidade sanitária do rebanho.",

		"Planilha oficial Sicoob preenchida":             "Modelo da cooperativa com dados do produtor, projeto e proposta.",
		"Planilha oficial Sicoob assinada pelo produtor": "Confirma que o produtor conferiu e aceitou os dados.",
		"Planilha oficial Sicoob assinada pelo técnico":  "Confirma responsabilidade técnica pelas informações.",
		"PUC - Proposta de Utilização de Crédito":        "Documento que mostra como o crédito será usado.",
		"PUC assinada":                            "Formaliza que o produtor concorda com a utilização do crédito.",
		"Carta de internalização":                 "Documento interno usado para encaminhar ou formalizar a operação dentro do Sicoob.",
		"ART ou TRT emitida, quitada e assinada":  "Comprova responsabilidade técnica válida.",
		"Cronograma de reembolso":                 "Mostra como e quando o produtor vai pagar.",
		"Orçamentos detalhados com CNPJ":          "Comprovam custos, fornecedor e formalidade da compra.",
		"Projeção zootécnica para pecuária":       "Mostra a evolução produtiva esperada do rebanho.",
		"Plano de manejo do rebanho":              "Explica alimentação, reprodução, sanidade, lotação e condução dos animais.",
		"Evolução de rebanho por eras/categorias": "Organiza o rebanho por idade e categoria para projetar produção e receita.",
	}

	if explicacao, existe := explicacoes[item]; existe {
		return explicacao
	}

	return "Item usado para conferência documental do projeto rural. Verifique se é obrigatório para o banco, para o tipo de projeto ou para a situação do produtor."
}
