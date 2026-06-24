package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func (app *App) telaWhatsAppPendencias(w http.ResponseWriter, r *http.Request) {
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

	mensagem := gerarMensagemPendenciasWhatsApp(reuniao, itens)

	tpl := template.Must(template.New("whatsappPendencias").Parse(htmlBase(whatsappPendenciasHTML)))

	dados := map[string]any{
		"Titulo":   "Pendências para WhatsApp",
		"Reuniao":  reuniao,
		"Mensagem": mensagem,
	}

	tpl.Execute(w, dados)
}

func (app *App) exportarWhatsAppPendenciasTXT(w http.ResponseWriter, r *http.Request) {
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

	texto := gerarMensagemPendenciasWhatsApp(reuniao, itens)

	nomeArquivo := fmt.Sprintf(
		"whatsapp_pendencias_%d_%s.txt",
		reuniao.ID,
		limparNomeArquivo(reuniao.Produtor),
	)

	caminho := filepath.Join(app.PastaDados, "exports", nomeArquivo)

	err = os.WriteFile(caminho, []byte(texto), 0644)
	if err != nil {
		http.Error(w, "Erro ao salvar arquivo: "+err.Error(), http.StatusInternalServerError)
		return
	}

	app.telaArquivoGerado(w, "WhatsApp de pendências TXT gerado", caminho)
}

func gerarMensagemPendenciasWhatsApp(reuniao Reuniao, itens []ItemChecklist) string {
	var pendentes []ItemChecklist

	for _, item := range itens {
		if item.Status == "Pendente" {
			pendentes = append(pendentes, item)
		}
	}

	nome := strings.TrimSpace(reuniao.Produtor)
	if nome == "" {
		nome = "produtor"
	}

	var b strings.Builder

	b.WriteString("Olá, " + nome + ". Segue a lista atualizada de pendências para continuidade do projeto rural:\n\n")

	if len(pendentes) == 0 {
		b.WriteString("No momento, não consta nenhuma pendência aberta no checklist.\n\n")
		b.WriteString("Assim que necessário, entrarei em contato para os próximos passos.")
		return b.String()
	}

	grupoAtual := ""
	contador := 1

	for _, item := range pendentes {
		if item.Grupo != grupoAtual {
			grupoAtual = item.Grupo
			b.WriteString("*" + grupoAtual + "*\n")
		}

		b.WriteString(strconv.Itoa(contador) + ". " + item.Item)

		observacao := strings.TrimSpace(item.Observacao)
		if observacao != "" {
			b.WriteString(" - " + observacao)
		}

		b.WriteString("\n")
		contador++
	}

	b.WriteString("\nAssim que conseguir enviar esses documentos/informações, damos continuidade à análise e montagem do projeto.")

	return b.String()
}
