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

	adicionarItem(&itens, "Documentos pessoais e cadastrais", "RG ou documento oficial com foto")
	adicionarItem(&itens, "Documentos pessoais e cadastrais", "CPF")
	adicionarItem(&itens, "Documentos pessoais e cadastrais", "Certidão de estado civil atualizada")
	adicionarItem(&itens, "Documentos pessoais e cadastrais", "Comprovante de residência atualizado")
	adicionarItem(&itens, "Documentos pessoais e cadastrais", "Declaração de IRPF completa")
	adicionarItem(&itens, "Documentos pessoais e cadastrais", "Recibo de entrega do IRPF")
	adicionarItem(&itens, "Documentos pessoais e cadastrais", "Relação de bens e direitos assinada")

	if r.CadastroBanco == "sim" {
		adicionarItem(&itens, "Cadastro bancário", "Conferir cadastro já existente no banco")
	} else {
		adicionarItem(&itens, "Cadastro bancário", "Fazer ou atualizar cadastro do produtor no banco")
	}

	adicionarItem(&itens, "Documentos fundiários", "Matrícula do imóvel atualizada")
	adicionarItem(&itens, "Documentos fundiários", "Certidão de inteiro teor")
	adicionarItem(&itens, "Documentos fundiários", "Certidão de ônus reais")
	adicionarItem(&itens, "Documentos fundiários", "Certidão de ações reipersecutórias")
	adicionarItem(&itens, "Documentos fundiários", "Conferir validade da matrícula")
	adicionarItem(&itens, "Documentos fundiários", "CCIR / Incra do último exercício")
	adicionarItem(&itens, "Documentos fundiários", "Comprovante de quitação do CCIR")
	adicionarItem(&itens, "Documentos fundiários", "ITR do último exercício")
	adicionarItem(&itens, "Documentos fundiários", "DARF do ITR quitado")

	if r.ImovelArrendado == "sim" {
		adicionarItem(&itens, "Documentos fundiários", "Contrato de arrendamento/parceria/comodato vigente")
		adicionarItem(&itens, "Documentos fundiários", "Contrato cobrindo todo o prazo do financiamento")
		adicionarItem(&itens, "Documentos fundiários", "Carta de anuência do proprietário")
		adicionarItem(&itens, "Documentos fundiários", "Autorização para penhor e execução do projeto")
	}

	if r.TemCAR == "sim" {
		adicionarItem(&itens, "Documentos ambientais", "Recibo do CAR ativo")
		adicionarItem(&itens, "Documentos ambientais", "Conferir situação do CAR")
	} else {
		adicionarItem(&itens, "Documentos ambientais", "Providenciar ou confirmar CAR da propriedade")
	}

	adicionarItem(&itens, "Documentos ambientais", "Verificar embargos ambientais")
	adicionarItem(&itens, "Documentos ambientais", "Verificar APP e Reserva Legal")
	adicionarItem(&itens, "Documentos ambientais", "Verificar unidade de conservação")
	adicionarItem(&itens, "Documentos ambientais", "Verificar terra indígena ou área quilombola")

	if r.UsaAgua == "sim" {
		adicionarItem(&itens, "Água e irrigação", "Outorga de uso de água")
		adicionarItem(&itens, "Água e irrigação", "Certidão de dispensa de outorga, se aplicável")
		adicionarItem(&itens, "Água e irrigação", "Licenciamento ou inexigibilidade ambiental")
		adicionarItem(&itens, "Água e irrigação", "Projeto técnico do sistema de irrigação")
		adicionarItem(&itens, "Água e irrigação", "Orçamento do sistema de irrigação")
	}

	if r.TemSupressao == "sim" {
		adicionarItem(&itens, "Supressão vegetal", "ASV - Autorização de Supressão Vegetal")
		adicionarItem(&itens, "Supressão vegetal", "Conferir compatibilidade com CAR, APP e Reserva Legal")
	}

	adicionarItem(&itens, "Documentos técnicos gerais", "Projeto técnico")
	adicionarItem(&itens, "Documentos técnicos gerais", "Orçamento técnico")
	adicionarItem(&itens, "Documentos técnicos gerais", "Memorial descritivo, se aplicável")
	adicionarItem(&itens, "Documentos técnicos gerais", "Croqui de localização")
	adicionarItem(&itens, "Documentos técnicos gerais", "Coordenadas geográficas em graus decimais")
	adicionarItem(&itens, "Documentos técnicos gerais", "ART ou TRT")
	adicionarItem(&itens, "Documentos técnicos gerais", "Comprovante de pagamento da ART/TRT")
	adicionarItem(&itens, "Documentos técnicos gerais", "Assinatura do técnico")
	adicionarItem(&itens, "Documentos técnicos gerais", "Assinatura do produtor")
	adicionarItem(&itens, "Documentos técnicos gerais", "Contrato de prestação de serviços técnicos")

	tipo := strings.ToLower(r.TipoProjeto)
	atividade := strings.ToLower(r.Atividade)

	if strings.Contains(tipo, "custeio agrícola") || r.PrecisaZARC == "sim" {
		adicionarItem(&itens, "Custeio agrícola", "Cultura financiada")
		adicionarItem(&itens, "Custeio agrícola", "Área plantada")
		adicionarItem(&itens, "Custeio agrícola", "Talhão / gleba")
		adicionarItem(&itens, "Custeio agrícola", "Variedade / cultivar")
		adicionarItem(&itens, "Custeio agrícola", "Data prevista de plantio")
		adicionarItem(&itens, "Custeio agrícola", "Data prevista de colheita")
		adicionarItem(&itens, "Custeio agrícola", "Análise de solo")
		adicionarItem(&itens, "Custeio agrícola", "Recomendação agronômica")
		adicionarItem(&itens, "Custeio agrícola", "Insumos previstos")
		adicionarItem(&itens, "Custeio agrícola", "Validar ZARC")
		adicionarItem(&itens, "Custeio agrícola", "Conferir município dentro do ZARC")
		adicionarItem(&itens, "Custeio agrícola", "Conferir janela de plantio")
		adicionarItem(&itens, "Custeio agrícola", "Conferir Proagro ou seguro rural")
	}

	if strings.Contains(tipo, "custeio pecuário") || r.TemPecuaria == "sim" || strings.Contains(atividade, "pecuária") {
		adicionarItem(&itens, "Pecuária", "Atividade pecuária definida")
		adicionarItem(&itens, "Pecuária", "Ficha sanitária do rebanho")
		adicionarItem(&itens, "Pecuária", "Declaração de vacinação")
		adicionarItem(&itens, "Pecuária", "Quantidade atual de animais")
		adicionarItem(&itens, "Pecuária", "Categorias do rebanho")
		adicionarItem(&itens, "Pecuária", "Capacidade de suporte")
		adicionarItem(&itens, "Pecuária", "Pastagens disponíveis")
		adicionarItem(&itens, "Pecuária", "Manejo previsto")
		adicionarItem(&itens, "Pecuária", "Evolução de rebanho")
		adicionarItem(&itens, "Pecuária", "Projeção zootécnica")
	}

	if strings.Contains(tipo, "investimento") || r.TemInvestimento == "sim" {
		adicionarItem(&itens, "Investimento", "Objetivo do investimento")
		adicionarItem(&itens, "Investimento", "Justificativa técnica")
		adicionarItem(&itens, "Investimento", "Mínimo de 3 propostas comerciais com CNPJ")
		adicionarItem(&itens, "Investimento", "Validade das propostas comerciais")
		adicionarItem(&itens, "Investimento", "Cronograma físico-financeiro")
		adicionarItem(&itens, "Investimento", "Cronograma de desembolso")
		adicionarItem(&itens, "Investimento", "Comprovação de capacidade de pagamento")
		adicionarItem(&itens, "Investimento", "Compatibilidade do investimento com a atividade")
	}

	if strings.Contains(atividade, "máquinas") || strings.Contains(atividade, "equipamentos") {
		adicionarItem(&itens, "Máquinas e equipamentos", "Tipo de máquina/equipamento")
		adicionarItem(&itens, "Máquinas e equipamentos", "Marca")
		adicionarItem(&itens, "Máquinas e equipamentos", "Modelo")
		adicionarItem(&itens, "Máquinas e equipamentos", "Ano, se usado")
		adicionarItem(&itens, "Máquinas e equipamentos", "Novo ou usado")
		adicionarItem(&itens, "Máquinas e equipamentos", "Proposta comercial")
		adicionarItem(&itens, "Máquinas e equipamentos", "Fornecedor com CNPJ")
		adicionarItem(&itens, "Máquinas e equipamentos", "Condições de entrega")
	}

	if r.TemObra == "sim" || strings.Contains(atividade, "obras") || strings.Contains(atividade, "benfeitorias") {
		adicionarItem(&itens, "Obras e benfeitorias", "Descrição da obra")
		adicionarItem(&itens, "Obras e benfeitorias", "Local da obra")
		adicionarItem(&itens, "Obras e benfeitorias", "Projeto técnico da construção")
		adicionarItem(&itens, "Obras e benfeitorias", "Memorial descritivo")
		adicionarItem(&itens, "Obras e benfeitorias", "Orçamento de materiais")
		adicionarItem(&itens, "Obras e benfeitorias", "Orçamento de mão de obra")
		adicionarItem(&itens, "Obras e benfeitorias", "Cronograma físico-financeiro")
		adicionarItem(&itens, "Obras e benfeitorias", "Licença ou autorização, se exigida")
		adicionarItem(&itens, "Obras e benfeitorias", "Fotos do local")
		adicionarItem(&itens, "Obras e benfeitorias", "Coordenadas do local")
	}

	banco := strings.ToLower(r.Banco)

	if strings.Contains(banco, "brasil") {
		adicionarItem(&itens, "Banco do Brasil", "Ficha cadastral Banco do Brasil")
		adicionarItem(&itens, "Banco do Brasil", "Projeto/orçamento emitido no sistema oficial BB")
		adicionarItem(&itens, "Banco do Brasil", "Projeto/orçamento assinado")
		adicionarItem(&itens, "Banco do Brasil", "Croqui com coordenadas em graus decimais")
		adicionarItem(&itens, "Banco do Brasil", "Matrícula com validade de 30 dias")
		adicionarItem(&itens, "Banco do Brasil", "Certidões de ônus e ações")
		adicionarItem(&itens, "Banco do Brasil", "CCIR quitado")
		adicionarItem(&itens, "Banco do Brasil", "ITR e DARF quitados")
		adicionarItem(&itens, "Banco do Brasil", "Anuência do proprietário, se aplicável")
		adicionarItem(&itens, "Banco do Brasil", "3 propostas comerciais para investimento")
		adicionarItem(&itens, "Banco do Brasil", "Cronograma físico-financeiro para investimento fixo")
		adicionarItem(&itens, "Banco do Brasil", "Ficha sanitária ou vacinação para pecuária")
	}

	if strings.Contains(banco, "sicoob") {
		adicionarItem(&itens, "Sicoob", "Planilha oficial Sicoob preenchida")
		adicionarItem(&itens, "Sicoob", "Planilha oficial Sicoob assinada pelo produtor")
		adicionarItem(&itens, "Sicoob", "Planilha oficial Sicoob assinada pelo técnico")
		adicionarItem(&itens, "Sicoob", "PUC - Proposta de Utilização de Crédito")
		adicionarItem(&itens, "Sicoob", "PUC assinada")
		adicionarItem(&itens, "Sicoob", "Carta de internalização")
		adicionarItem(&itens, "Sicoob", "Contrato de prestação de serviços técnicos")
		adicionarItem(&itens, "Sicoob", "ART ou TRT emitida, quitada e assinada")
		adicionarItem(&itens, "Sicoob", "Cronograma de desembolso")
		adicionarItem(&itens, "Sicoob", "Cronograma de reembolso")
		adicionarItem(&itens, "Sicoob", "Orçamentos detalhados com CNPJ")
		adicionarItem(&itens, "Sicoob", "Projeção zootécnica para pecuária")
		adicionarItem(&itens, "Sicoob", "Plano de manejo do rebanho")
		adicionarItem(&itens, "Sicoob", "Evolução de rebanho por eras/categorias")
	}

	return itens
}

func adicionarItem(itens *[]ItemChecklist, grupo string, item string) {
	*itens = append(*itens, ItemChecklist{
		Grupo: grupo,
		Item:  item,
	})
}
