package main

const configuracoesAPIHTML = `
<h2>Configurações das APIs</h2>

<div class="barra-acoes">
	<a class="botao secundario" href="/">Início</a>
	<a class="botao secundario" href="/reunioes">Reuniões</a>
</div>

<div class="card destaque">
	<h3>Para que serve esta tela?</h3>
	<p>
		Aqui ficam as chaves e URLs usadas nas consultas externas do sistema.
		Essas informações ficam salvas apenas no seu computador.
	</p>

	<p class="pequeno">
		Arquivo local:
	</p>

	<pre>{{.CaminhoConfig}}</pre>
</div>

{{if .Mensagem}}
<div class="card destaque">
	<p>{{.Mensagem}}</p>
</div>
{{end}}

<div class="grid">
	<div class="card">
		<h3>Status das credenciais</h3>

		<p><strong>Conecta Gov / SICAR:</strong> {{.ConectaGovTokenMascarado}}</p>
		<p><strong>Portal da Transparência:</strong> {{.PortalTransparenciaMascarado}}</p>

		<p class="pequeno">
			Se o Conecta Gov não estiver liberado, o SICAR por CNPJ não vai puxar CAR automaticamente.
			Nesse caso, informe o CAR manualmente na investigação.
		</p>
	</div>

	<div class="card alerta">
		<h3>Segurança</h3>
		<p>
			Não envie suas chaves por WhatsApp, e-mail ou GitHub.
			Se alguma chave for exposta, gere outra no serviço oficial.
		</p>
	</div>
</div>

<form method="POST" action="/configuracoes-api">
	<div class="card">
		<h3>1. Conecta Gov / SICAR</h3>

		<p class="pequeno">
			Use estes campos somente se você tiver acesso liberado ao Conecta Gov/Serpro.
			Essas credenciais permitem gerar token para consultar APIs como SICAR CPF/CNPJ,
			SICAR Imóvel e SICAR Tema.
		</p>

		<label>Client ID / Chave Conecta Gov</label>
		<input 
			type="text" 
			name="conecta_gov_client_id" 
			value="{{.Config.ConectaGovClientID}}" 
			placeholder="Cole aqui o Client ID / chave"
		>

		<label>Client Secret / Senha Conecta Gov</label>
		<textarea 
			name="conecta_gov_client_secret" 
			rows="3" 
			placeholder="Cole aqui o Client Secret / senha"
		>{{.Config.ConectaGovClientSecret}}</textarea>

		<label>Token Conecta Gov manual</label>
		<textarea 
			name="conecta_gov_token" 
			rows="3" 
			placeholder="Opcional. Normalmente deixe vazio para o app gerar o token automaticamente."
		>{{.Config.ConectaGovToken}}</textarea>

		<p class="pequeno">
			Recomendado: preencher Client ID e Client Secret, deixando o token manual vazio.
		</p>
	</div>

	<div class="card">
		<h3>2. URLs SICAR</h3>

		<p class="pequeno">
			As URLs precisam conter <strong>%s</strong>, que é onde o app coloca o CNPJ ou o CAR.
			Se você não souber alterar, mantenha os valores padrão.
		</p>

		<label>URL SICAR por CNPJ</label>
		<input 
			type="text" 
			name="sicar_cnpj_url" 
			value="{{.Config.SICARCNPJURL}}" 
			placeholder="https://.../%s"
		>

		<label>URL SICAR por CAR / Imóvel</label>
		<input 
			type="text" 
			name="sicar_car_url" 
			value="{{.Config.SICARCARURL}}" 
			placeholder="https://.../%s"
		>

		<label>URL SICAR Tema / Área / Polígono</label>
		<input 
			type="text" 
			name="sicar_tema_url" 
			value="{{.Config.SICARTemaURL}}" 
			placeholder="https://.../%s"
		>
	</div>

	<div class="card">
		<h3>3. Portal da Transparência</h3>

		<p class="pequeno">
			Essa chave é usada para consultas como CEIS e CNEP.
		</p>

		<label>Token Portal da Transparência</label>
		<textarea 
			name="portal_transparencia_token" 
			rows="3" 
			placeholder="Cole aqui a chave da API do Portal da Transparência"
		>{{.Config.PortalTransparenciaToken}}</textarea>
	</div>

	<div class="card">
		<h3>4. Salvar e testar</h3>

		<button type="submit">Salvar configurações</button>

		<a class="botao secundario" href="/testar-auth-conecta-gov" target="_blank">
			Testar autenticação Conecta Gov
		</a>

		<p class="pequeno">
			Primeiro salve as configurações. Depois clique em testar autenticação.
			Se gerar token, a autenticação básica está funcionando.
		</p>
	</div>
</form>
`
