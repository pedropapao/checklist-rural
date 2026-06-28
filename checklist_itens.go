package main

import (
	"net/http"
	"strings"
	"time"
)

type ItemChecklist struct {
	ID           int
	ReuniaoID    int
	Grupo        string
	Item         string
	Status       string
	Observacao   string
	Explicacao   string
	CriadoEm     string
	AtualizadoEm string
}

func (app *App) criarTabelaItensChecklist() error {
	_, err := app.DB.Exec(`
		CREATE TABLE IF NOT EXISTS checklist_itens (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			reuniao_id INTEGER NOT NULL,
			grupo TEXT NOT NULL,
			item TEXT NOT NULL,
			status TEXT NOT NULL DEFAULT 'Pendente',
			observacao TEXT NOT NULL DEFAULT '',
			criado_em TEXT NOT NULL,
			atualizado_em TEXT NOT NULL,
			UNIQUE(reuniao_id, grupo, item)
		)
	`)
	return err
}

func (app *App) gerarItensChecklist(w http.ResponseWriter, r *http.Request) {
	reuniao, err := app.pegarReuniaoDaURL(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	itens := montarItensChecklistDaReuniao(reuniao)
	agora := time.Now().Format("02/01/2006 15:04")

	for _, item := range itens {
		_, err := app.DB.Exec(`
			INSERT OR IGNORE INTO checklist_itens
			(reuniao_id, grupo, item, status, observacao, criado_em, atualizado_em)
			VALUES (?, ?, ?, ?, ?, ?, ?)
		`,
			reuniao.ID,
			item.Grupo,
			item.Item,
			"Pendente",
			"",
			agora,
			agora,
		)

		if err != nil {
			http.Error(w, "Erro ao gerar item do checklist: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	http.Redirect(w, r, "/checklist?id="+r.URL.Query().Get("id"), http.StatusSeeOther)
}

func montarItensChecklistDaReuniao(r Reuniao) []ItemChecklist {
	var itens []ItemChecklist

	// 1. Documentos pessoais e cadastrais
	adicionarItem(&itens, "Documentos pessoais e cadastrais", "Documentos pessoais básicos do produtor: RG/CNH, CPF, estado civil e comprovante de residência")
	adicionarItem(&itens, "Documentos pessoais e cadastrais", "IRPF completo com recibo de entrega e relação de bens/direitos")

	// 2. Cadastro bancário
	if r.CadastroBanco == "sim" {
		adicionarItem(&itens, "Cadastro bancário", "Conferir cadastro existente no banco e pendências cadastrais")
	} else {
		adicionarItem(&itens, "Cadastro bancário", "Fazer ou atualizar cadastro do produtor no banco")
	}

	if r.FinanciamentoAtivo == "sim" {
		adicionarItem(&itens, "Cadastro bancário", "Conferir financiamentos rurais ativos e comprometimento de limite")
	}

	if r.RestricaoCadastral == "sim" {
		adicionarItem(&itens, "Cadastro bancário", "Verificar restrição cadastral antes de avançar com o projeto")
	}

	// 3. Documentos fundiários
	adicionarItem(&itens, "Documentos fundiários", "Matrícula/certidões do imóvel atualizadas: inteiro teor, ônus, ações e validade")
	adicionarItem(&itens, "Documentos fundiários", "CCIR/Incra atualizado com comprovante de quitação")
	adicionarItem(&itens, "Documentos fundiários", "ITR do último exercício com DARF quitado")

	if r.ImovelArrendado == "sim" {
		adicionarItem(&itens, "Documentos fundiários", "Contrato de arrendamento/parceria/comodato vigente, com prazo suficiente para o financiamento")
		adicionarItem(&itens, "Documentos fundiários", "Anuência do proprietário e autorização para execução do projeto/garantia, se exigido")
	}

	// 4. Documentos ambientais
	if r.TemCAR == "sim" {
		adicionarItem(&itens, "Documentos ambientais", "CAR da propriedade: recibo ativo e conferência da situação")
	} else {
		adicionarItem(&itens, "Documentos ambientais", "Providenciar ou confirmar CAR da propriedade")
	}

	adicionarItem(&itens, "Documentos ambientais", "Consulta ambiental da área: embargos, APP, Reserva Legal, unidade de conservação, terra indígena ou quilombola")

	// 5. Água e irrigação
	if r.UsaAgua == "sim" {
		adicionarItem(&itens, "Água e irrigação", "Regularidade do uso da água: outorga, dispensa, licenciamento ou inexigibilidade ambiental")
		adicionarItem(&itens, "Água e irrigação", "Projeto e orçamento do sistema de irrigação ou uso de água")
	}

	// 6. Supressão vegetal
	if r.TemSupressao == "sim" {
		adicionarItem(&itens, "Supressão vegetal", "Autorização de supressão vegetal e compatibilidade com CAR, APP e Reserva Legal")
	}

	// 7. Documentos técnicos gerais
	adicionarItem(&itens, "Documentos técnicos gerais", "Projeto técnico completo: objetivo, orçamento, memorial, croqui e coordenadas")
	adicionarItem(&itens, "Documentos técnicos gerais", "Responsabilidade técnica: ART/TRT emitida, quitada, assinada e contrato técnico")
	adicionarItem(&itens, "Documentos técnicos gerais", "Assinaturas obrigatórias do produtor e do técnico")

	tipo := strings.ToLower(r.TipoProjeto)
	atividade := strings.ToLower(r.Atividade)

	// 8. Custeio agrícola
	if strings.Contains(tipo, "custeio agrícola") || r.PrecisaZARC == "sim" {
		adicionarItem(&itens, "Custeio agrícola", "Dados da lavoura: cultura, área, talhão/gleba, variedade, plantio e colheita")
		adicionarItem(&itens, "Custeio agrícola", "Base agronômica: análise de solo, recomendação técnica e insumos previstos")
		adicionarItem(&itens, "Custeio agrícola", "Validação ZARC: município, cultura e janela de plantio")
		adicionarItem(&itens, "Custeio agrícola", "Conferir Proagro ou seguro rural, quando aplicável")
	}

	// 9. Pecuária
	if strings.Contains(tipo, "custeio pecuário") || r.TemPecuaria == "sim" || strings.Contains(atividade, "pecuária") {
		adicionarItem(&itens, "Pecuária", "Atividade pecuária definida e finalidade do crédito")
		adicionarItem(&itens, "Pecuária", "Sanidade do rebanho: ficha sanitária e declaração de vacinação")
		adicionarItem(&itens, "Pecuária", "Rebanho atual: quantidade de animais e categorias")
		adicionarItem(&itens, "Pecuária", "Suporte produtivo: pastagens, capacidade de suporte e manejo previsto")
		adicionarItem(&itens, "Pecuária", "Evolução do rebanho e projeção zootécnica")
	}

	// 10. Investimento
	if strings.Contains(tipo, "investimento") || r.TemInvestimento == "sim" {
		adicionarItem(&itens, "Investimento", "Objetivo, justificativa técnica e compatibilidade do investimento com a atividade")
		adicionarItem(&itens, "Investimento", "Propostas comerciais com CNPJ e validade")
		adicionarItem(&itens, "Investimento", "Cronograma físico-financeiro e cronograma de desembolso")
		adicionarItem(&itens, "Investimento", "Comprovação da capacidade de pagamento")
	}

	// 11. Máquinas e equipamentos
	if strings.Contains(atividade, "máquinas") || strings.Contains(atividade, "equipamentos") {
		adicionarItem(&itens, "Máquinas e equipamentos", "Identificação da máquina/equipamento: tipo, marca, modelo, ano e condição novo/usado")
		adicionarItem(&itens, "Máquinas e equipamentos", "Proposta comercial do fornecedor com CNPJ e condições de entrega")
	}

	// 12. Obras e benfeitorias
	if r.TemObra == "sim" || strings.Contains(atividade, "obras") || strings.Contains(atividade, "benfeitorias") {
		adicionarItem(&itens, "Obras e benfeitorias", "Caracterização da obra: descrição, local, fotos e coordenadas")
		adicionarItem(&itens, "Obras e benfeitorias", "Projeto da construção, memorial descritivo e orçamento de materiais/mão de obra")
		adicionarItem(&itens, "Obras e benfeitorias", "Cronograma físico-financeiro e licença/autorização, se exigida")
	}

	banco := strings.ToLower(r.Banco)

	// 13. Banco do Brasil
	if strings.Contains(banco, "brasil") {
		adicionarItem(&itens, "Banco do Brasil", "Conferir ficha cadastral e exigências específicas do Banco do Brasil")
		adicionarItem(&itens, "Banco do Brasil", "Projeto/orçamento no padrão ou sistema oficial do BB, quando exigido")
		adicionarItem(&itens, "Banco do Brasil", "Conferir documentos do imóvel conforme validade exigida pelo BB")
		adicionarItem(&itens, "Banco do Brasil", "Conferir documentos específicos por finalidade: investimento, pecuária, custeio ou garantia")
	}

	// 14. Sicoob
	if strings.Contains(banco, "sicoob") {
		adicionarItem(&itens, "Sicoob", "Planilha oficial Sicoob preenchida e assinada pelo produtor e técnico")
		adicionarItem(&itens, "Sicoob", "PUC - Proposta de Utilização de Crédito preenchida e assinada")
		adicionarItem(&itens, "Sicoob", "Carta de internalização, contrato técnico e ART/TRT emitida, quitada e assinada")
		adicionarItem(&itens, "Sicoob", "Cronogramas de desembolso e reembolso")
		adicionarItem(&itens, "Sicoob", "Orçamentos detalhados com CNPJ")
		adicionarItem(&itens, "Sicoob", "Documentos zootécnicos para pecuária: projeção, manejo e evolução do rebanho")
	}

	return itens
}

func adicionarItem(itens *[]ItemChecklist, grupo string, item string) {
	*itens = append(*itens, ItemChecklist{
		Grupo: grupo,
		Item:  item,
	})
}
