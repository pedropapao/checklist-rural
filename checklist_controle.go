package main

import (
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type ResumoChecklist struct {
	Total               int
	Pendentes           int
	Recebidos           int
	NaoSeAplica         int
	PercentualConcluido int
}

func calcularResumoChecklist(itens []ItemChecklist) ResumoChecklist {
	var resumo ResumoChecklist

	resumo.Total = len(itens)

	for _, item := range itens {
		switch item.Status {
		case "Recebido":
			resumo.Recebidos++
		case "Não se aplica":
			resumo.NaoSeAplica++
		default:
			resumo.Pendentes++
		}
	}

	if resumo.Total > 0 {
		resumo.PercentualConcluido = ((resumo.Recebidos + resumo.NaoSeAplica) * 100) / resumo.Total
	}

	return resumo
}

func (app *App) garantirItensChecklist(reuniao Reuniao) error {
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
			return err
		}
	}

	return nil
}

func (app *App) listarItensChecklist(reuniaoID int) ([]ItemChecklist, error) {
	linhas, err := app.DB.Query(`
		SELECT
			id,
			reuniao_id,
			grupo,
			item,
			status,
			observacao,
			criado_em,
			atualizado_em
		FROM checklist_itens
		WHERE reuniao_id = ?
		ORDER BY grupo, id
	`, reuniaoID)

	if err != nil {
		return nil, err
	}
	defer linhas.Close()

	var itens []ItemChecklist

	for linhas.Next() {
		var item ItemChecklist

		err := linhas.Scan(
			&item.ID,
			&item.ReuniaoID,
			&item.Grupo,
			&item.Item,
			&item.Status,
			&item.Observacao,
			&item.CriadoEm,
			&item.AtualizadoEm,
		)

		if err != nil {
			return nil, err
		}

		itens = append(itens, item)
	}

	return itens, nil
}

func (app *App) telaChecklistControle(w http.ResponseWriter, r *http.Request) {
	reuniao, err := app.pegarReuniaoDaURL(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = app.garantirItensChecklist(reuniao)
	if err != nil {
		http.Error(w, "Erro ao gerar itens do checklist: "+err.Error(), http.StatusInternalServerError)
		return
	}

	itens, err := app.listarItensChecklist(reuniao.ID)
	if err != nil {
		http.Error(w, "Erro ao listar itens do checklist: "+err.Error(), http.StatusInternalServerError)
		return
	}

	resumo := calcularResumoChecklist(itens)
	grupos := agruparItensChecklistParaTela(itens)

	tpl := template.Must(template.New("checklistControle").Parse(htmlBase(checklistControleHTML)))

	dados := map[string]any{
		"Titulo":  "Controle do checklist",
		"Reuniao": reuniao,
		"Itens":   itens,
		"Grupos":  grupos,
		"Resumo":  resumo,
	}

	tpl.Execute(w, dados)
}

func (app *App) salvarItensChecklist(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/reunioes", http.StatusSeeOther)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Erro ao ler formulário", http.StatusBadRequest)
		return
	}

	reuniaoIDTexto := r.FormValue("reuniao_id")

	reuniaoID, err := strconv.Atoi(reuniaoIDTexto)
	if err != nil {
		http.Error(w, "ID da reunião inválido", http.StatusBadRequest)
		return
	}

	agora := time.Now().Format("02/01/2006 15:04")

	for nomeCampo, valores := range r.PostForm {
		if !strings.HasPrefix(nomeCampo, "status_") {
			continue
		}

		itemIDTexto := strings.TrimPrefix(nomeCampo, "status_")

		itemID, err := strconv.Atoi(itemIDTexto)
		if err != nil {
			continue
		}

		status := "Pendente"
		if len(valores) > 0 {
			status = valores[0]
		}

		if status != "Pendente" && status != "Recebido" && status != "Não se aplica" {
			status = "Pendente"
		}

		observacao := r.FormValue("observacao_" + itemIDTexto)

		_, err = app.DB.Exec(`
			UPDATE checklist_itens
			SET status = ?, observacao = ?, atualizado_em = ?
			WHERE id = ? AND reuniao_id = ?
		`,
			status,
			observacao,
			agora,
			itemID,
			reuniaoID,
		)

		if err != nil {
			http.Error(w, "Erro ao salvar item do checklist: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	http.Redirect(w, r, "/checklist-controle?id="+strconv.Itoa(reuniaoID), http.StatusSeeOther)
}
