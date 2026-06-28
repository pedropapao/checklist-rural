package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func (app *App) telaInicial(w http.ResponseWriter, r *http.Request) {
	tpl := template.Must(template.New("inicio").Parse(htmlBase(inicioHTML)))

	dados := map[string]any{
		"Titulo": "Checklist Rural",
	}

	tpl.Execute(w, dados)
}

func (app *App) telaNovaReuniao(w http.ResponseWriter, r *http.Request) {
	tpl := template.Must(template.New("nova").Parse(htmlBase(novaReuniaoHTML)))

	dados := map[string]any{
		"Titulo": "Nova reunião",
	}

	tpl.Execute(w, dados)
}

func (app *App) salvarReuniao(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/nova-reuniao", http.StatusSeeOther)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Erro ao ler formulário", http.StatusBadRequest)
		return
	}

	produtor := r.FormValue("produtor")
	telefone := r.FormValue("telefone")
	municipio := r.FormValue("municipio")
	uf := r.FormValue("uf")
	banco := r.FormValue("banco")
	tipoProjeto := r.FormValue("tipo_projeto")
	atividade := r.FormValue("atividade")
	observacoes := r.FormValue("observacoes")
	rendaAnualTexto := r.FormValue("renda_anual")
	finalidadeCredito := r.FormValue("finalidade_credito")
	valorPretendidoTexto := r.FormValue("valor_pretendido")

	rendaAnual, err := normalizarValorMonetario(rendaAnualTexto)
	if err != nil {
		http.Error(w, "Renda anual inválida", http.StatusBadRequest)
		return
	}

	valorPretendido, err := normalizarValorMonetario(valorPretendidoTexto)
	if err != nil {
		http.Error(w, "Valor pretendido inválido", http.StatusBadRequest)
		return
	}

	classificacaoProdutor := classificarProdutorPorRBA(rendaAnual)

	criadoEm := time.Now().Format("02/01/2006 15:04")

	if produtor == "" {
		http.Error(w, "Informe o nome do produtor", http.StatusBadRequest)
		return
	}

	_, err = app.DB.Exec(`
		INSERT INTO reunioes 
		(
			produtor, telefone, municipio, uf, banco, tipo_projeto, atividade,
			renda_anual, classificacao_produtor, possui_caf, finalidade_credito, valor_pretendido, tem_orcamento,
			cadastro_banco, financiamento_ativo, restricao_cadastral,
			imovel_proprio, imovel_arrendado, tem_car, usa_agua,
			tem_pecuaria, tem_investimento, tem_obra, tem_supressao, precisa_zarc,
			observacoes, criado_em
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`,
		produtor,
		telefone,
		municipio,
		uf,
		banco,
		tipoProjeto,
		atividade,
		rendaAnual,
		classificacaoProdutor,
		checkValor(r, "possui_caf"),
		finalidadeCredito,
		valorPretendido,
		checkValor(r, "tem_orcamento"),
		checkValor(r, "cadastro_banco"),
		checkValor(r, "financiamento_ativo"),
		checkValor(r, "restricao_cadastral"),
		checkValor(r, "imovel_proprio"),
		checkValor(r, "imovel_arrendado"),
		checkValor(r, "tem_car"),
		checkValor(r, "usa_agua"),
		checkValor(r, "tem_pecuaria"),
		checkValor(r, "tem_investimento"),
		checkValor(r, "tem_obra"),
		checkValor(r, "tem_supressao"),
		checkValor(r, "precisa_zarc"),
		observacoes,
		criadoEm,
	)

	if err != nil {
		http.Error(w, "Erro ao salvar reunião: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/reunioes", http.StatusSeeOther)
}

func (app *App) listarReunioes(w http.ResponseWriter, r *http.Request) {
	linhas, err := app.DB.Query(`
		SELECT 
			id,
			COALESCE(produtor, ''),
			COALESCE(telefone, ''),
			COALESCE(municipio, ''),
			COALESCE(uf, ''),
			COALESCE(banco, ''),
			COALESCE(tipo_projeto, ''),
			COALESCE(atividade, ''),
			COALESCE(renda_anual, 0),
			COALESCE(classificacao_produtor, ''),
			COALESCE(possui_caf, 'nao'),
			COALESCE(finalidade_credito, ''),
			COALESCE(valor_pretendido, 0),
			COALESCE(tem_orcamento, 'nao'),
			COALESCE(cadastro_banco, 'nao'),
			COALESCE(financiamento_ativo, 'nao'),
			COALESCE(restricao_cadastral, 'nao'),
			COALESCE(imovel_proprio, 'nao'),
			COALESCE(imovel_arrendado, 'nao'),
			COALESCE(tem_car, 'nao'),
			COALESCE(usa_agua, 'nao'),
			COALESCE(tem_pecuaria, 'nao'),
			COALESCE(tem_investimento, 'nao'),
			COALESCE(tem_obra, 'nao'),
			COALESCE(tem_supressao, 'nao'),
			COALESCE(precisa_zarc, 'nao'),
			COALESCE(observacoes, ''),
			COALESCE(criado_em, '')
		FROM reunioes
		ORDER BY id DESC
	`)
	if err != nil {
		http.Error(w, "Erro ao listar reuniões: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer linhas.Close()

	var reunioes []ReuniaoComResumo

	for linhas.Next() {
		var reuniao Reuniao

		err := linhas.Scan(
			&reuniao.ID,
			&reuniao.Produtor,
			&reuniao.Telefone,
			&reuniao.Municipio,
			&reuniao.UF,
			&reuniao.Banco,
			&reuniao.TipoProjeto,
			&reuniao.Atividade,
			&reuniao.RendaAnual,
			&reuniao.ClassificacaoProdutor,
			&reuniao.PossuiCAF,
			&reuniao.FinalidadeCredito,
			&reuniao.ValorPretendido,
			&reuniao.TemOrcamento,
			&reuniao.CadastroBanco,
			&reuniao.FinanciamentoAtivo,
			&reuniao.RestricaoCadastral,
			&reuniao.ImovelProprio,
			&reuniao.ImovelArrendado,
			&reuniao.TemCAR,
			&reuniao.UsaAgua,
			&reuniao.TemPecuaria,
			&reuniao.TemInvestimento,
			&reuniao.TemObra,
			&reuniao.TemSupressao,
			&reuniao.PrecisaZARC,
			&reuniao.Observacoes,
			&reuniao.CriadoEm,
		)

		if err != nil {
			http.Error(w, "Erro ao ler reunião: "+err.Error(), http.StatusInternalServerError)
			return
		}

		resumo := app.resumoDaReuniao(reuniao)
		reunioes = append(reunioes, ReuniaoComResumo{Reuniao: reuniao, Resumo: resumo})
	}

	filtros := FiltrosReuniao{
		Busca:    r.URL.Query().Get("busca"),
		Banco:    r.URL.Query().Get("banco"),
		Situacao: r.URL.Query().Get("situacao"),
	}

	reunioes = filtrarReunioes(reunioes, filtros)

	tpl := template.Must(template.New("reunioes").Parse(htmlBase(reunioesHTML)))

	dados := map[string]any{
		"Titulo":   "Reuniões salvas",
		"Reunioes": reunioes,
		"Filtros":  filtros,
	}

	tpl.Execute(w, dados)
}

func (app *App) buscarReuniaoPorID(id int) (Reuniao, error) {
	var reuniao Reuniao

	err := app.DB.QueryRow(`
		SELECT 
			id,
			COALESCE(produtor, ''),
			COALESCE(telefone, ''),
			COALESCE(municipio, ''),
			COALESCE(uf, ''),
			COALESCE(banco, ''),
			COALESCE(tipo_projeto, ''),
			COALESCE(atividade, ''),
			COALESCE(renda_anual, 0),
			COALESCE(classificacao_produtor, ''),
			COALESCE(possui_caf, 'nao'),
			COALESCE(finalidade_credito, ''),
			COALESCE(valor_pretendido, 0),
			COALESCE(tem_orcamento, 'nao'),
			COALESCE(cadastro_banco, 'nao'),
			COALESCE(financiamento_ativo, 'nao'),
			COALESCE(restricao_cadastral, 'nao'),
			COALESCE(imovel_proprio, 'nao'),
			COALESCE(imovel_arrendado, 'nao'),
			COALESCE(tem_car, 'nao'),
			COALESCE(usa_agua, 'nao'),
			COALESCE(tem_pecuaria, 'nao'),
			COALESCE(tem_investimento, 'nao'),
			COALESCE(tem_obra, 'nao'),
			COALESCE(tem_supressao, 'nao'),
			COALESCE(precisa_zarc, 'nao'),
			COALESCE(observacoes, ''),
			COALESCE(criado_em, '')
		FROM reunioes
		WHERE id = ?
	`, id).Scan(
		&reuniao.ID,
		&reuniao.Produtor,
		&reuniao.Telefone,
		&reuniao.Municipio,
		&reuniao.UF,
		&reuniao.Banco,
		&reuniao.TipoProjeto,
		&reuniao.Atividade,
		&reuniao.RendaAnual,
		&reuniao.ClassificacaoProdutor,
		&reuniao.PossuiCAF,
		&reuniao.FinalidadeCredito,
		&reuniao.ValorPretendido,
		&reuniao.TemOrcamento,
		&reuniao.CadastroBanco,
		&reuniao.FinanciamentoAtivo,
		&reuniao.RestricaoCadastral,
		&reuniao.ImovelProprio,
		&reuniao.ImovelArrendado,
		&reuniao.TemCAR,
		&reuniao.UsaAgua,
		&reuniao.TemPecuaria,
		&reuniao.TemInvestimento,
		&reuniao.TemObra,
		&reuniao.TemSupressao,
		&reuniao.PrecisaZARC,
		&reuniao.Observacoes,
		&reuniao.CriadoEm,
	)

	if err != nil {
		return reuniao, err
	}

	return reuniao, nil
}

func (app *App) telaEditarReuniao(w http.ResponseWriter, r *http.Request) {
	reuniao, err := app.pegarReuniaoDaURL(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tpl := template.Must(template.New("editar").Parse(htmlBase(editarReuniaoHTML)))

	dados := map[string]any{
		"Titulo":  "Editar reunião",
		"Reuniao": reuniao,
		"Leitura": montarLeituraInicial(reuniao),
	}

	tpl.Execute(w, dados)
}

func (app *App) atualizarReuniao(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/reunioes", http.StatusSeeOther)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Erro ao ler formulário", http.StatusBadRequest)
		return
	}

	idTexto := r.FormValue("id")

	id, err := strconv.Atoi(idTexto)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	produtor := r.FormValue("produtor")
	telefone := r.FormValue("telefone")
	municipio := r.FormValue("municipio")
	uf := r.FormValue("uf")
	banco := r.FormValue("banco")
	tipoProjeto := r.FormValue("tipo_projeto")
	atividade := r.FormValue("atividade")
	observacoes := r.FormValue("observacoes")
	rendaAnualTexto := r.FormValue("renda_anual")
	finalidadeCredito := r.FormValue("finalidade_credito")
	valorPretendidoTexto := r.FormValue("valor_pretendido")

	rendaAnual, err := normalizarValorMonetario(rendaAnualTexto)
	if err != nil {
		http.Error(w, "Renda anual inválida", http.StatusBadRequest)
		return
	}

	valorPretendido, err := normalizarValorMonetario(valorPretendidoTexto)
	if err != nil {
		http.Error(w, "Valor pretendido inválido", http.StatusBadRequest)
		return
	}

	classificacaoProdutor := classificarProdutorPorRBA(rendaAnual)

	if produtor == "" {
		http.Error(w, "Informe o nome do produtor", http.StatusBadRequest)
		return
	}

	_, err = app.DB.Exec(`
		UPDATE reunioes SET
			produtor = ?,
			telefone = ?,
			municipio = ?,
			uf = ?,
			banco = ?,
			tipo_projeto = ?,
			atividade = ?,
			renda_anual = ?,
			classificacao_produtor = ?,
			possui_caf = ?,
			finalidade_credito = ?,
			valor_pretendido = ?,
			tem_orcamento = ?,
			cadastro_banco = ?,
			financiamento_ativo = ?,
			restricao_cadastral = ?,
			imovel_proprio = ?,
			imovel_arrendado = ?,
			tem_car = ?,
			usa_agua = ?,
			tem_pecuaria = ?,
			tem_investimento = ?,
			tem_obra = ?,
			tem_supressao = ?,
			precisa_zarc = ?,
			observacoes = ?
		WHERE id = ?
	`,
		produtor,
		telefone,
		municipio,
		uf,
		banco,
		tipoProjeto,
		atividade,
		rendaAnual,
		classificacaoProdutor,
		checkValor(r, "possui_caf"),
		finalidadeCredito,
		valorPretendido,
		checkValor(r, "tem_orcamento"),
		checkValor(r, "cadastro_banco"),
		checkValor(r, "financiamento_ativo"),
		checkValor(r, "restricao_cadastral"),
		checkValor(r, "imovel_proprio"),
		checkValor(r, "imovel_arrendado"),
		checkValor(r, "tem_car"),
		checkValor(r, "usa_agua"),
		checkValor(r, "tem_pecuaria"),
		checkValor(r, "tem_investimento"),
		checkValor(r, "tem_obra"),
		checkValor(r, "tem_supressao"),
		checkValor(r, "precisa_zarc"),
		observacoes,
		id,
	)

	if err != nil {
		http.Error(w, "Erro ao atualizar reunião: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/detalhes?id="+strconv.Itoa(id), http.StatusSeeOther)
}

func (app *App) telaConfirmarExcluir(w http.ResponseWriter, r *http.Request) {
	reuniao, err := app.pegarReuniaoDaURL(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tpl := template.Must(template.New("confirmarExcluir").Parse(htmlBase(confirmarExcluirHTML)))

	dados := map[string]any{
		"Titulo":  "Confirmar exclusão",
		"Reuniao": reuniao,
		"Leitura": montarLeituraInicial(reuniao),
	}

	tpl.Execute(w, dados)
}

func (app *App) excluirReuniao(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/reunioes", http.StatusSeeOther)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Erro ao ler formulário", http.StatusBadRequest)
		return
	}

	idTexto := r.FormValue("id")

	id, err := strconv.Atoi(idTexto)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	_, err = app.DB.Exec("DELETE FROM reunioes WHERE id = ?", id)
	if err != nil {
		http.Error(w, "Erro ao excluir reunião: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/reunioes", http.StatusSeeOther)
}

func (app *App) telaDetalhes(w http.ResponseWriter, r *http.Request) {
	reuniao, err := app.pegarReuniaoDaURL(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tpl := template.Must(template.New("detalhes").Parse(htmlBase(detalhesHTML)))

	dados := map[string]any{
		"Titulo":  "Detalhes da reunião",
		"Reuniao": reuniao,
		"Leitura": montarLeituraInicial(reuniao),
		"Resumo":  app.resumoDaReuniao(reuniao),
	}

	tpl.Execute(w, dados)
}

func (app *App) telaChecklist(w http.ResponseWriter, r *http.Request) {
	reuniao, err := app.pegarReuniaoDaURL(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	checklist := gerarChecklistDaReuniao(reuniao)

	tpl := template.Must(template.New("checklist").Parse(htmlBase(checklistHTML)))

	dados := map[string]any{
		"Titulo":    "Checklist automático",
		"Reuniao":   reuniao,
		"Checklist": checklist,
	}

	tpl.Execute(w, dados)
}

func (app *App) telaWhatsApp(w http.ResponseWriter, r *http.Request) {
	reuniao, err := app.pegarReuniaoDaURL(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	mensagem := gerarMensagemWhatsApp(reuniao)

	tpl := template.Must(template.New("whatsapp").Parse(htmlBase(whatsappHTML)))

	dados := map[string]any{
		"Titulo":   "Resumo WhatsApp",
		"Reuniao":  reuniao,
		"Mensagem": mensagem,
	}

	tpl.Execute(w, dados)
}

func (app *App) exportarChecklistTXT(w http.ResponseWriter, r *http.Request) {
	reuniao, err := app.pegarReuniaoDaURL(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	texto := gerarChecklistDaReuniao(reuniao)
	nomeArquivo := fmt.Sprintf("checklist_%d_%s.txt", reuniao.ID, limparNomeArquivo(reuniao.Produtor))
	caminho := filepath.Join(app.PastaDados, "exports", nomeArquivo)

	err = os.WriteFile(caminho, []byte(texto), 0644)
	if err != nil {
		http.Error(w, "Erro ao salvar arquivo: "+err.Error(), http.StatusInternalServerError)
		return
	}

	app.telaArquivoGerado(w, "Checklist TXT gerado", caminho)
}

func (app *App) exportarWhatsAppTXT(w http.ResponseWriter, r *http.Request) {
	reuniao, err := app.pegarReuniaoDaURL(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	texto := gerarMensagemWhatsApp(reuniao)
	nomeArquivo := fmt.Sprintf("whatsapp_%d_%s.txt", reuniao.ID, limparNomeArquivo(reuniao.Produtor))
	caminho := filepath.Join(app.PastaDados, "exports", nomeArquivo)

	err = os.WriteFile(caminho, []byte(texto), 0644)
	if err != nil {
		http.Error(w, "Erro ao salvar arquivo: "+err.Error(), http.StatusInternalServerError)
		return
	}

	app.telaArquivoGerado(w, "Mensagem WhatsApp gerada", caminho)
}

func (app *App) telaArquivoGerado(w http.ResponseWriter, titulo string, caminho string) {
	tpl := template.Must(template.New("arquivo").Parse(htmlBase(arquivoGeradoHTML)))

	dados := map[string]any{
		"Titulo":  titulo,
		"Caminho": caminho,
	}

	tpl.Execute(w, dados)
}

func (app *App) pegarReuniaoDaURL(r *http.Request) (Reuniao, error) {
	idTexto := r.URL.Query().Get("id")

	id, err := strconv.Atoi(idTexto)
	if err != nil {
		return Reuniao{}, fmt.Errorf("ID da reunião inválido")
	}

	reuniao, err := app.buscarReuniaoPorID(id)
	if err != nil {
		return Reuniao{}, fmt.Errorf("reunião não encontrada")
	}

	return reuniao, nil
}
