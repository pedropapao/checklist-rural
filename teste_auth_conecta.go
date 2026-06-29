package main

import (
	"html/template"
	"net/http"
	"time"
)

func (app *App) testarAuthConectaGov(w http.ResponseWriter, r *http.Request) {
	carregarConfigEnv()

	token, err := obterTokenConectaGov()

	dados := map[string]any{
		"Titulo":   "Teste Conecta Gov",
		"Sucesso":  err == nil,
		"Erro":     "",
		"Token":    "",
		"DataHora": time.Now().Format("02/01/2006 15:04"),
	}

	if err != nil {
		dados["Erro"] = err.Error()
	} else {
		dados["Token"] = tokenMascarado(token)
	}

	tpl := template.Must(template.New("teste_auth_conecta").Parse(htmlBase(testeAuthConectaHTML)))
	tpl.Execute(w, dados)
}

const testeAuthConectaHTML = `
<h2>Teste de autenticação Conecta Gov</h2>

<div class="barra-acoes">
	<a class="botao secundario" href="/configuracoes-api">Voltar para configurações</a>
	<a class="botao secundario" href="/">Início</a>
</div>

{{if .Sucesso}}
<div class="card destaque">
	<h3>Autenticação funcionando</h3>

	<p>
		O app conseguiu gerar um token no Conecta Gov.
	</p>

	<p><strong>Token:</strong> {{.Token}}</p>
	<p><strong>Horário:</strong> {{.DataHora}}</p>

	<p class="pequeno">
		Se a consulta SICAR ainda der erro ou timeout, o problema não é mais login básico.
		Pode ser falta de autorização específica para a API SICAR, IP não liberado, URL diferente ou instabilidade do gateway.
	</p>
</div>
{{else}}
<div class="card perigo">
	<h3>Não foi possível autenticar</h3>

	<p>
		O app não conseguiu gerar token no Conecta Gov.
	</p>

	<pre>{{.Erro}}</pre>

	<p class="pequeno">
		Confira Client ID, Client Secret e se sua credencial realmente está liberada no Conecta Gov/Serpro.
	</p>
</div>
{{end}}

<div class="card">
	<h3>Como interpretar</h3>

	<ul>
		<li><strong>Status 401:</strong> Client ID ou Client Secret incorretos.</li>
		<li><strong>Status 403:</strong> credencial existe, mas não tem permissão.</li>
		<li><strong>Timeout:</strong> gateway, rede, IP ou serviço indisponível.</li>
		<li><strong>Token gerado:</strong> autenticação básica funcionou.</li>
	</ul>
</div>
`
