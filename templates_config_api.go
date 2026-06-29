package main

const configuracoesAPIHTML = `
<h2>Configurações das APIs</h2>

<p class="pequeno">
	Essas configurações ficam salvas apenas no seu computador, no arquivo:
</p>

<pre>{{.CaminhoConfig}}</pre>

{{if .Mensagem}}
<div class="card destaque">
	<p>{{.Mensagem}}</p>
</div>
{{end}}

<div class="card">
	<h3>Status atual</h3>

	<p><strong>Conecta Gov / SICAR:</strong> {{.ConectaGovTokenMascarado}}</p>
	<p><strong>Portal da Transparência:</strong> {{.PortalTransparenciaMascarado}}</p>
</div>

<form method="POST" action="/configuracoes-api">
	<div class="card">
		<h3>Conecta Gov / SICAR</h3>

		<label>Token Conecta Gov</label>
		<textarea name="conecta_gov_token" rows="4" placeholder="Cole aqui o token do Conecta Gov/SICAR">{{.Config.ConectaGovToken}}</textarea>

		<label>URL SICAR por CNPJ</label>
		<input type="text" name="sicar_cnpj_url" value="{{.Config.SICARCNPJURL}}" placeholder="https://.../%s">

		<p class="pequeno">
			Use <strong>%s</strong> no lugar onde o CNPJ deve entrar.
		</p>

		<label>URL SICAR por CAR</label>
		<input type="text" name="sicar_car_url" value="{{.Config.SICARCARURL}}" placeholder="https://.../%s">

		<p class="pequeno">
			Use <strong>%s</strong> no lugar onde o código CAR deve entrar.
		</p>
	</div>

	<div class="card">
		<h3>Portal da Transparência</h3>

		<label>Token Portal da Transparência</label>
		<textarea name="portal_transparencia_token" rows="3" placeholder="Cole aqui a chave da API do Portal da Transparência">{{.Config.PortalTransparenciaToken}}</textarea>
	</div>

	<button type="submit">Salvar configurações</button>

	<a class="botao secundario" href="/">Voltar</a>
</form>

<div class="card alerta">
	<h3>Atenção</h3>
	<p>
		Não envie esse arquivo de configuração para GitHub, WhatsApp ou e-mail.
		Ele pode conter tokens de acesso.
	</p>
</div>
`
