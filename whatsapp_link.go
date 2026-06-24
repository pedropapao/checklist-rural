package main

import (
	"net/http"
	"net/url"
	"strings"
	"unicode"
)

func limparTelefoneWhatsApp(telefone string) string {
	var numeros strings.Builder

	for _, caractere := range telefone {
		if unicode.IsDigit(caractere) {
			numeros.WriteRune(caractere)
		}
	}

	limpo := numeros.String()

	if limpo == "" {
		return ""
	}

	// Se já começar com 55, mantém.
	if strings.HasPrefix(limpo, "55") {
		return limpo
	}

	// Se tiver DDD + número, adiciona código do Brasil.
	if len(limpo) == 10 || len(limpo) == 11 {
		return "55" + limpo
	}

	return limpo
}

func montarLinkWhatsApp(telefone string, mensagem string) string {
	numero := limparTelefoneWhatsApp(telefone)

	if numero == "" {
		return ""
	}

	mensagemEscapada := url.QueryEscape(mensagem)

	return "https://wa.me/" + numero + "?text=" + mensagemEscapada
}

func (app *App) abrirWhatsAppPendencias(w http.ResponseWriter, r *http.Request) {
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
	link := montarLinkWhatsApp(reuniao.Telefone, mensagem)

	if link == "" {
		http.Error(w, "Telefone do produtor não informado ou inválido", http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, link, http.StatusSeeOther)
}
