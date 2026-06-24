package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func (app *App) exportarChecklistControleTXT(w http.ResponseWriter, r *http.Request) {
	reuniao, err := app.pegarReuniaoDaURL(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = app.garantirItensChecklist(reuniao)
	if err != nil {
		http.Error(w, "Erro ao garantir itens do checklist: "+err.Error(), http.StatusInternalServerError)
		return
	}

	itens, err := app.listarItensChecklist(reuniao.ID)
	if err != nil {
		http.Error(w, "Erro ao listar itens do checklist: "+err.Error(), http.StatusInternalServerError)
		return
	}

	texto := gerarTextoChecklistControle(reuniao, itens)

	nomeArquivo := fmt.Sprintf(
		"checklist_controle_%d_%s.txt",
		reuniao.ID,
		limparNomeArquivo(reuniao.Produtor),
	)

	caminho := filepath.Join(app.PastaDados, "exports", nomeArquivo)

	err = os.WriteFile(caminho, []byte(texto), 0644)
	if err != nil {
		http.Error(w, "Erro ao salvar arquivo: "+err.Error(), http.StatusInternalServerError)
		return
	}

	app.telaArquivoGerado(w, "Checklist controlado TXT gerado", caminho)
}

func gerarTextoChecklistControle(reuniao Reuniao, itens []ItemChecklist) string {
	var b strings.Builder

	b.WriteString("CHECKLIST CONTROLADO DA REUNIÃO\n")
	b.WriteString("==================================================\n\n")

	b.WriteString("DADOS DA REUNIÃO\n")
	b.WriteString("--------------------------------------------------\n")
	b.WriteString("Produtor: " + reuniao.Produtor + "\n")
	b.WriteString("Telefone: " + reuniao.Telefone + "\n")
	b.WriteString("Município/UF: " + reuniao.Municipio + "/" + reuniao.UF + "\n")
	b.WriteString("Banco pretendido: " + reuniao.Banco + "\n")
	b.WriteString("Tipo de projeto: " + reuniao.TipoProjeto + "\n")
	b.WriteString("Atividade: " + reuniao.Atividade + "\n")
	b.WriteString("Data da reunião: " + reuniao.CriadoEm + "\n\n")

	total := len(itens)
	pendentes := 0
	recebidos := 0
	naoAplica := 0

	for _, item := range itens {
		switch item.Status {
		case "Recebido":
			recebidos++
		case "Não se aplica":
			naoAplica++
		default:
			pendentes++
		}
	}

	b.WriteString("RESUMO DO ANDAMENTO\n")
	b.WriteString("--------------------------------------------------\n")
	b.WriteString(fmt.Sprintf("Total de itens: %d\n", total))
	b.WriteString(fmt.Sprintf("Pendentes: %d\n", pendentes))
	b.WriteString(fmt.Sprintf("Recebidos: %d\n", recebidos))
	b.WriteString(fmt.Sprintf("Não se aplica: %d\n\n", naoAplica))

	grupoAtual := ""

	for _, item := range itens {
		if item.Grupo != grupoAtual {
			grupoAtual = item.Grupo
			b.WriteString(strings.ToUpper(grupoAtual) + "\n")
			b.WriteString("--------------------------------------------------\n")
		}

		b.WriteString("[")
		b.WriteString(item.Status)
		b.WriteString("] ")
		b.WriteString(item.Item)
		b.WriteString("\n")

		if strings.TrimSpace(item.Observacao) != "" {
			b.WriteString("Observação: ")
			b.WriteString(item.Observacao)
			b.WriteString("\n")
		}

		b.WriteString("\n")
	}

	return b.String()
}
