package main

import (
	"strconv"
	"strings"
)

func gerarChecklistDaReuniao(r Reuniao) string {
	var b strings.Builder

	b.WriteString("CHECKLIST DA REUNIÃO\n")
	b.WriteString("==================================================\n\n")

	b.WriteString("DADOS DO ATENDIMENTO\n")
	b.WriteString("--------------------------------------------------\n")
	b.WriteString("Produtor: " + r.Produtor + "\n")
	b.WriteString("Telefone: " + r.Telefone + "\n")
	b.WriteString("Município/UF: " + r.Municipio + "/" + r.UF + "\n")
	b.WriteString("Banco pretendido: " + r.Banco + "\n")
	b.WriteString("Tipo de projeto: " + r.TipoProjeto + "\n")
	b.WriteString("Atividade: " + r.Atividade + "\n")
	b.WriteString("Data da reunião: " + r.CriadoEm + "\n\n")

	adicionarDocumentosPessoais(&b, r)
	adicionarDocumentosFundiarios(&b, r)
	adicionarDocumentosAmbientais(&b, r)
	adicionarDocumentosTecnicos(&b)
	adicionarCondicionais(&b, r)
	adicionarBanco(&b, r)
	adicionarAlertas(&b, r)
	adicionarProximoPasso(&b)

	return b.String()
}

func adicionarDocumentosPessoais(b *strings.Builder, r Reuniao) {
	b.WriteString("DOCUMENTOS PESSOAIS E CADASTRAIS\n")
	b.WriteString("--------------------------------------------------\n")
	b.WriteString("[ ] RG ou documento oficial com foto\n")
	b.WriteString("[ ] CPF\n")
	b.WriteString("[ ] Certidão de estado civil atualizada\n")
	b.WriteString("[ ] Comprovante de residência atualizado\n")
	b.WriteString("[ ] Declaração de IRPF completa\n")
	b.WriteString("[ ] Recibo de entrega do IRPF\n")
	b.WriteString("[ ] Relação de bens e direitos assinada\n")

	if r.CadastroBanco == "sim" {
		b.WriteString("[ ] Conferir cadastro já existente no banco\n")
	} else {
		b.WriteString("[ ] Fazer ou atualizar cadastro do produtor no banco\n")
	}

	b.WriteString("\n")
}

func adicionarDocumentosFundiarios(b *strings.Builder, r Reuniao) {
	b.WriteString("DOCUMENTOS FUNDIÁRIOS\n")
	b.WriteString("--------------------------------------------------\n")
	b.WriteString("[ ] Matrícula do imóvel atualizada\n")
	b.WriteString("[ ] Certidão de inteiro teor\n")
	b.WriteString("[ ] Certidão de ônus reais\n")
	b.WriteString("[ ] Certidão de ações reipersecutórias\n")
	b.WriteString("[ ] Conferir validade da matrícula: 30 dias\n")
	b.WriteString("[ ] CCIR / Incra do último exercício\n")
	b.WriteString("[ ] Comprovante de quitação do CCIR\n")
	b.WriteString("[ ] ITR do último exercício\n")
	b.WriteString("[ ] DARF do ITR quitado\n")

	if r.ImovelArrendado == "sim" {
		b.WriteString("[ ] Contrato de arrendamento/parceria/comodato vigente\n")
		b.WriteString("[ ] Contrato cobrindo todo o prazo do financiamento\n")
		b.WriteString("[ ] Carta de anuência do proprietário\n")
		b.WriteString("[ ] Autorização para penhor e execução do projeto\n")
	}

	b.WriteString("\n")
}

func adicionarDocumentosAmbientais(b *strings.Builder, r Reuniao) {
	b.WriteString("DOCUMENTOS AMBIENTAIS\n")
	b.WriteString("--------------------------------------------------\n")

	if r.TemCAR == "sim" {
		b.WriteString("[ ] Recibo do CAR ativo\n")
		b.WriteString("[ ] Conferir situação do CAR\n")
	} else {
		b.WriteString("[ ] Providenciar/confirmar CAR da propriedade\n")
	}

	b.WriteString("[ ] Verificar embargos ambientais\n")
	b.WriteString("[ ] Verificar APP e Reserva Legal\n")
	b.WriteString("[ ] Verificar unidade de conservação\n")
	b.WriteString("[ ] Verificar terra indígena ou área quilombola\n")

	if r.UsaAgua == "sim" {
		b.WriteString("[ ] Outorga de uso de água\n")
		b.WriteString("[ ] Certidão de dispensa de outorga, se aplicável\n")
		b.WriteString("[ ] Licenciamento ou inexigibilidade ambiental\n")
	}

	if r.TemSupressao == "sim" {
		b.WriteString("[ ] ASV - Autorização de Supressão Vegetal\n")
		b.WriteString("[ ] Conferir compatibilidade com CAR, APP e Reserva Legal\n")
	}

	b.WriteString("\n")
}

func adicionarDocumentosTecnicos(b *strings.Builder) {
	b.WriteString("DOCUMENTOS TÉCNICOS GERAIS\n")
	b.WriteString("--------------------------------------------------\n")
	b.WriteString("[ ] Projeto técnico\n")
	b.WriteString("[ ] Orçamento técnico\n")
	b.WriteString("[ ] Memorial descritivo, se aplicável\n")
	b.WriteString("[ ] Croqui de localização\n")
	b.WriteString("[ ] Coordenadas geográficas em graus decimais\n")
	b.WriteString("[ ] ART ou TRT\n")
	b.WriteString("[ ] Comprovante de pagamento da ART/TRT\n")
	b.WriteString("[ ] Assinatura do técnico\n")
	b.WriteString("[ ] Assinatura do produtor\n")
	b.WriteString("[ ] Contrato de prestação de serviços técnicos\n\n")
}

func adicionarCondicionais(b *strings.Builder, r Reuniao) {
	tipo := strings.ToLower(r.TipoProjeto)
	atividade := strings.ToLower(r.Atividade)

	if strings.Contains(tipo, "custeio agrícola") || r.PrecisaZARC == "sim" {
		b.WriteString("CUSTEIO AGRÍCOLA\n")
		b.WriteString("--------------------------------------------------\n")
		b.WriteString("[ ] Cultura financiada\n")
		b.WriteString("[ ] Área plantada\n")
		b.WriteString("[ ] Talhão / gleba\n")
		b.WriteString("[ ] Variedade / cultivar\n")
		b.WriteString("[ ] Data prevista de plantio\n")
		b.WriteString("[ ] Data prevista de colheita\n")
		b.WriteString("[ ] Análise de solo\n")
		b.WriteString("[ ] Recomendação agronômica\n")
		b.WriteString("[ ] Insumos previstos\n")
		b.WriteString("[ ] Validar ZARC\n")
		b.WriteString("[ ] Conferir município dentro do ZARC\n")
		b.WriteString("[ ] Conferir janela de plantio\n")
		b.WriteString("[ ] Conferir Proagro ou seguro rural\n\n")
	}

	if strings.Contains(tipo, "custeio pecuário") || r.TemPecuaria == "sim" || strings.Contains(atividade, "pecuária") {
		b.WriteString("PECUÁRIA\n")
		b.WriteString("--------------------------------------------------\n")
		b.WriteString("[ ] Atividade pecuária definida\n")
		b.WriteString("[ ] Ficha sanitária do rebanho\n")
		b.WriteString("[ ] Declaração de vacinação\n")
		b.WriteString("[ ] Quantidade atual de animais\n")
		b.WriteString("[ ] Categorias do rebanho\n")
		b.WriteString("[ ] Capacidade de suporte\n")
		b.WriteString("[ ] Pastagens disponíveis\n")
		b.WriteString("[ ] Manejo previsto\n")
		b.WriteString("[ ] Evolução de rebanho\n")
		b.WriteString("[ ] Projeção zootécnica\n\n")
	}

	if strings.Contains(tipo, "investimento") || r.TemInvestimento == "sim" {
		b.WriteString("INVESTIMENTO\n")
		b.WriteString("--------------------------------------------------\n")
		b.WriteString("[ ] Objetivo do investimento\n")
		b.WriteString("[ ] Justificativa técnica\n")
		b.WriteString("[ ] Mínimo de 3 propostas comerciais com CNPJ\n")
		b.WriteString("[ ] Validade das propostas comerciais\n")
		b.WriteString("[ ] Cronograma físico-financeiro\n")
		b.WriteString("[ ] Cronograma de desembolso\n")
		b.WriteString("[ ] Comprovação de capacidade de pagamento\n")
		b.WriteString("[ ] Compatibilidade do investimento com a atividade\n\n")
	}

	if strings.Contains(atividade, "máquinas") || strings.Contains(atividade, "equipamentos") {
		b.WriteString("MÁQUINAS E EQUIPAMENTOS\n")
		b.WriteString("--------------------------------------------------\n")
		b.WriteString("[ ] Tipo de máquina/equipamento\n")
		b.WriteString("[ ] Marca\n")
		b.WriteString("[ ] Modelo\n")
		b.WriteString("[ ] Ano, se usado\n")
		b.WriteString("[ ] Novo ou usado\n")
		b.WriteString("[ ] Proposta comercial\n")
		b.WriteString("[ ] Fornecedor com CNPJ\n")
		b.WriteString("[ ] Condições de entrega\n\n")
	}

	if r.TemObra == "sim" || strings.Contains(atividade, "obras") || strings.Contains(atividade, "benfeitorias") {
		b.WriteString("OBRAS, CONSTRUÇÕES E BENFEITORIAS\n")
		b.WriteString("--------------------------------------------------\n")
		b.WriteString("[ ] Descrição da obra\n")
		b.WriteString("[ ] Local da obra\n")
		b.WriteString("[ ] Projeto técnico da construção\n")
		b.WriteString("[ ] Memorial descritivo\n")
		b.WriteString("[ ] Orçamento de materiais\n")
		b.WriteString("[ ] Orçamento de mão de obra\n")
		b.WriteString("[ ] Cronograma físico-financeiro\n")
		b.WriteString("[ ] Licença ou autorização, se exigida\n")
		b.WriteString("[ ] Fotos do local\n")
		b.WriteString("[ ] Coordenadas do local\n\n")
	}

	if r.UsaAgua == "sim" || strings.Contains(atividade, "irrigação") {
		b.WriteString("IRRIGAÇÃO / USO DE ÁGUA\n")
		b.WriteString("--------------------------------------------------\n")
		b.WriteString("[ ] Tipo de irrigação\n")
		b.WriteString("[ ] Fonte de água\n")
		b.WriteString("[ ] Vazão necessária\n")
		b.WriteString("[ ] Área irrigada\n")
		b.WriteString("[ ] Outorga de uso de água ou dispensa\n")
		b.WriteString("[ ] Licença ambiental\n")
		b.WriteString("[ ] Projeto técnico do sistema\n")
		b.WriteString("[ ] Orçamento do sistema\n")
		b.WriteString("[ ] Energia elétrica disponível\n")
		b.WriteString("[ ] Croqui do sistema\n\n")
	}
}

func adicionarBanco(b *strings.Builder, r Reuniao) {
	banco := strings.ToLower(r.Banco)

	if strings.Contains(banco, "brasil") {
		b.WriteString("BANCO DO BRASIL\n")
		b.WriteString("--------------------------------------------------\n")
		b.WriteString("[ ] Ficha cadastral Banco do Brasil\n")
		b.WriteString("[ ] Projeto/orçamento emitido no sistema oficial BB\n")
		b.WriteString("[ ] Projeto/orçamento assinado\n")
		b.WriteString("[ ] Croqui com coordenadas em graus decimais\n")
		b.WriteString("[ ] Matrícula com validade de 30 dias\n")
		b.WriteString("[ ] Certidões de ônus e ações\n")
		b.WriteString("[ ] CCIR quitado\n")
		b.WriteString("[ ] ITR e DARF quitados\n")
		b.WriteString("[ ] Anuência do proprietário, se aplicável\n")
		b.WriteString("[ ] 3 propostas comerciais para investimento\n")
		b.WriteString("[ ] Cronograma físico-financeiro para investimento fixo\n")
		b.WriteString("[ ] Ficha sanitária ou vacinação para pecuária\n\n")
	}

	if strings.Contains(banco, "sicoob") {
		b.WriteString("SICOOB\n")
		b.WriteString("--------------------------------------------------\n")
		b.WriteString("[ ] Planilha oficial Sicoob preenchida\n")
		b.WriteString("[ ] Planilha oficial Sicoob assinada pelo produtor\n")
		b.WriteString("[ ] Planilha oficial Sicoob assinada pelo técnico\n")
		b.WriteString("[ ] PUC - Proposta de Utilização de Crédito\n")
		b.WriteString("[ ] PUC assinada\n")
		b.WriteString("[ ] Carta de internalização\n")
		b.WriteString("[ ] Contrato de prestação de serviços técnicos\n")
		b.WriteString("[ ] ART ou TRT emitida, quitada e assinada\n")
		b.WriteString("[ ] Cronograma de desembolso\n")
		b.WriteString("[ ] Cronograma de reembolso\n")
		b.WriteString("[ ] Orçamentos detalhados com CNPJ\n")
		b.WriteString("[ ] Projeção zootécnica para pecuária\n")
		b.WriteString("[ ] Plano de manejo do rebanho\n")
		b.WriteString("[ ] Evolução de rebanho por eras/categorias\n\n")
	}

	if strings.Contains(banco, "ainda") || banco == "" {
		b.WriteString("BANCO AINDA NÃO DEFINIDO\n")
		b.WriteString("--------------------------------------------------\n")
		b.WriteString("[ ] Conferir se a operação seguirá pelo Banco do Brasil\n")
		b.WriteString("[ ] Conferir se a operação seguirá pelo Sicoob\n")
		b.WriteString("[ ] Validar checklist específico após definição do banco\n\n")
	}
}

func adicionarAlertas(b *strings.Builder, r Reuniao) {
	b.WriteString("PENDÊNCIAS E ALERTAS AUTOMÁTICOS\n")
	b.WriteString("--------------------------------------------------\n")

	temAlerta := false

	if r.RestricaoCadastral == "sim" {
		b.WriteString("⚠️ Produtor informou restrição cadastral conhecida. Conferir antes de avançar.\n")
		temAlerta = true
	}

	if r.TemCAR == "nao" {
		b.WriteString("⚠️ CAR não confirmado na reunião.\n")
		temAlerta = true
	}

	if r.ImovelArrendado == "sim" {
		b.WriteString("⚠️ Conferir contrato e anuência do proprietário.\n")
		temAlerta = true
	}

	if r.UsaAgua == "sim" {
		b.WriteString("⚠️ Conferir outorga ou dispensa de outorga.\n")
		temAlerta = true
	}

	if r.TemPecuaria == "sim" {
		b.WriteString("⚠️ Conferir ficha sanitária, vacinação e evolução do rebanho.\n")
		temAlerta = true
	}

	if r.TemInvestimento == "sim" {
		b.WriteString("⚠️ Conferir 3 propostas comerciais e cronograma físico-financeiro.\n")
		temAlerta = true
	}

	if r.TemObra == "sim" {
		b.WriteString("⚠️ Conferir projeto técnico, memorial e orçamento da obra.\n")
		temAlerta = true
	}

	if r.TemSupressao == "sim" {
		b.WriteString("⚠️ Conferir ASV antes de avançar.\n")
		temAlerta = true
	}

	if r.PrecisaZARC == "sim" {
		b.WriteString("⚠️ Validar ZARC antes de fechar o custeio.\n")
		temAlerta = true
	}

	if !temAlerta {
		b.WriteString("✅ Nenhum alerta automático crítico pela triagem inicial.\n")
	}

	b.WriteString("\n")
}

func adicionarProximoPasso(b *strings.Builder) {
	b.WriteString("PRÓXIMO PASSO DA REUNIÃO\n")
	b.WriteString("--------------------------------------------------\n")
	b.WriteString("[ ] Definir pendências principais\n")
	b.WriteString("[ ] Combinar prazo para envio dos documentos\n")
	b.WriteString("[ ] Enviar resumo da reunião ao produtor\n")
}

func gerarMensagemWhatsApp(r Reuniao) string {
	var pendencias []string

	pendencias = append(pendencias, "RG ou documento oficial com foto")
	pendencias = append(pendencias, "CPF")
	pendencias = append(pendencias, "Comprovante de residência atualizado")
	pendencias = append(pendencias, "IRPF completo com recibo")
	pendencias = append(pendencias, "Matrícula atualizada do imóvel")
	pendencias = append(pendencias, "CCIR/Incra")
	pendencias = append(pendencias, "ITR e DARF quitado")

	if r.TemCAR == "nao" {
		pendencias = append(pendencias, "CAR da propriedade")
	}

	if r.ImovelArrendado == "sim" {
		pendencias = append(pendencias, "Contrato de arrendamento/parceria/comodato")
		pendencias = append(pendencias, "Carta de anuência do proprietário")
	}

	if r.UsaAgua == "sim" {
		pendencias = append(pendencias, "Outorga de uso de água ou certidão de dispensa")
	}

	if r.TemInvestimento == "sim" {
		pendencias = append(pendencias, "3 propostas comerciais com CNPJ")
		pendencias = append(pendencias, "Cronograma físico-financeiro")
	}

	if r.TemPecuaria == "sim" {
		pendencias = append(pendencias, "Ficha sanitária do rebanho")
		pendencias = append(pendencias, "Declaração de vacinação")
		pendencias = append(pendencias, "Evolução/projeção do rebanho")
	}

	if r.PrecisaZARC == "sim" {
		pendencias = append(pendencias, "Informações para validação do ZARC")
	}

	var b strings.Builder

	nome := r.Produtor
	if nome == "" {
		nome = "produtor"
	}

	b.WriteString("Olá, " + nome + ". Segue a lista inicial de pendências para andamento do projeto rural:\n\n")

	for i, item := range pendencias {
		b.WriteString(strconv.Itoa(i+1) + ". " + item + "\n")
	}

	b.WriteString("\nAssim que enviar esses documentos, damos continuidade à análise e montagem do projeto.\n")

	return b.String()
}
