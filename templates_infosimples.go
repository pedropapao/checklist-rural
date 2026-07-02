package main

const infosimplesConfigHTML = `
<section class="card">
	<h2>InfoSimples</h2>

	<p class="pequeno">
		Configure aqui o token da InfoSimples. Esse arquivo fica salvo somente no seu computador.
		Não envie esse token para GitHub, WhatsApp ou prints.
	</p>

	{{if .Mensagem}}
		<div class="alerta sucesso">{{.Mensagem}}</div>
	{{end}}

	{{if .Erro}}
		<div class="alerta perigo">{{.Erro}}</div>
	{{end}}

	<div class="grade">
		<div>
			<strong>Status do token:</strong><br>
			<span>{{.TokenMascarado}}</span>
		</div>
		<div>
			<strong>Arquivo local:</strong><br>
			<span class="pequeno">{{.CaminhoConfig}}</span>
		</div>
	</div>

	<form method="post" class="formulario">
		<label>Token InfoSimples</label>
		<input type="password" name="token" value="{{.Config.Token}}" placeholder="Cole aqui seu token da InfoSimples">

		<label>Endpoint CAR / Demonstrativo</label>
		<input type="text" name="car_demonstrativo_url" value="{{.Config.CARDemonstrativoURL}}">

		<label>Endpoint CAR / Demonstrativo PDF</label>
		<input type="text" name="car_demonstrativo_pdf_url" value="{{.Config.CARDemonstrativoPDFURL}}">

		<label>Endpoint CAR / Download Shapefile</label>
		<input type="text" name="car_download_shapefile_url" value="{{.Config.CARDownloadShapefileURL}}">

		<div class="barra-acoes">
			<button class="botao" type="submit">Salvar InfoSimples</button>
			<a class="botao secundario" href="/infosimples-car">Testar CAR / Demonstrativo</a>
			<a class="botao secundario" href="/configuracoes-api">Voltar às APIs</a>
		</div>
	</form>
</section>
`

const infosimplesCARHTML = `
<section class="card">
	<h2>InfoSimples - CAR / Demonstrativo</h2>

	<p class="pequeno">
		Esta tela faz uma consulta real na InfoSimples usando o número do CAR.
		A consulta pode consumir crédito da sua conta. Confira o CAR antes de enviar.
	</p>

	<div class="grade">
		<div>
			<strong>Token:</strong><br>
			<span>{{.TokenMascarado}}</span>
		</div>
		<div>
			<strong>Endpoint:</strong><br>
			<span class="pequeno">{{.Config.CARDemonstrativoURL}}</span>
		</div>
	</div>

	{{if .Reuniao.ID}}
		<div class="alerta">
			<strong>Reunião:</strong> {{.Reuniao.Produtor}} —
			{{.Reuniao.Municipio}}/{{.Reuniao.UF}}
		</div>
	{{end}}

	{{if .Erro}}
		<div class="alerta perigo">{{.Erro}}</div>
	{{end}}

	{{if .Sucesso}}
		<div class="alerta sucesso">{{.Sucesso}}</div>
	{{end}}

	<form method="post" class="formulario">
		<label>Número do CAR</label>
		<input type="text" name="car" value="{{.CAR}}" placeholder="Ex: MG-0000000-...">

		<div class="alerta perigo">
			Atenção: ao clicar em consultar, a InfoSimples pode cobrar a consulta do seu saldo.
			Não teste várias vezes com CAR errado.
		</div>

		<div class="barra-acoes">
			<button class="botao" type="submit" name="acao" value="consultar">Consultar sem salvar</button>

			{{if .Reuniao.ID}}
				<button class="botao" type="submit" name="acao" value="consultar_salvar">
					Consultar e salvar na reunião
				</button>
			{{end}}

			<a class="botao secundario" href="/infosimples-config">Configurar token</a>
			{{if .Reuniao.ID}}
				<a class="botao secundario" href="/infosimples-car-pdf?id={{.Reuniao.ID}}">Demonstrativo PDF</a>
				<a class="botao secundario" href="/infosimples-car-shapefile?id={{.Reuniao.ID}}">Download Shapefile</a>
				<a class="botao secundario" href="/investigacao?id={{.Reuniao.ID}}">Voltar à investigação</a>
			{{end}}
			<a class="botao secundario" href="/">Início</a>
		</div>
	</form>
</section>

{{if .Consultou}}
<section class="card">
	<h3>Resumo da resposta</h3>

	<div class="grade">
		<div>
			<strong>Código:</strong><br>
			{{.Resposta.Code}}
		</div>
		<div>
			<strong>Mensagem:</strong><br>
			{{.Resposta.CodeMessage}}
		</div>
		<div>
			<strong>Cobrável:</strong><br>
			{{.Resposta.Header.Billable}}
		</div>
		<div>
			<strong>Preço informado:</strong><br>
			R$ {{.Resposta.Header.Price}}
		</div>
		<div>
			<strong>Serviço:</strong><br>
			{{.Resposta.Header.Service}}
		</div>
		<div>
			<strong>Itens retornados:</strong><br>
			{{.Resposta.DataCount}}
		</div>
	</div>

	{{if .Resposta.Errors}}
		<h4>Erros</h4>
		<ul>
			{{range .Resposta.Errors}}
				<li>{{.}}</li>
			{{end}}
		</ul>
	{{end}}
</section>

<section class="card">
	<h3>JSON bruto retornado</h3>
	<p class="pequeno">
		Antes de integrarmos com os campos da reunião, vamos conferir quais campos a InfoSimples retorna.
	</p>
	<pre style="white-space: pre-wrap; overflow:auto; max-height:600px;">{{.JSONBonito}}</pre>
</section>
{{end}}
`

const infosimplesCARPDFHTML = `
<section class="card">
	<h2>InfoSimples - CAR / Demonstrativo PDF</h2>

	<p class="pequeno">
		Esta consulta pode consumir crédito. Ela deve retornar link/arquivo de visualização do demonstrativo em PDF ou HTML.
	</p>

	<div class="grade">
		<div>
			<strong>Token:</strong><br>
			<span>{{.TokenMascarado}}</span>
		</div>
		<div>
			<strong>Endpoint PDF:</strong><br>
			<span class="pequeno">{{.Config.CARDemonstrativoPDFURL}}</span>
		</div>
	</div>

	{{if .Reuniao.ID}}
		<div class="alerta">
			<strong>Reunião:</strong> {{.Reuniao.Produtor}} —
			{{.Reuniao.Municipio}}/{{.Reuniao.UF}}
		</div>
	{{end}}

	{{if .Erro}}
		<div class="alerta perigo">{{.Erro}}</div>
	{{end}}

	{{if .Sucesso}}
		<div class="alerta sucesso">{{.Sucesso}}</div>
	{{end}}

	<form method="post" class="formulario">
		<label>Número do CAR</label>
		<input type="text" name="car" value="{{.CAR}}" placeholder="Ex: MG-0000000-...">

		<div class="alerta perigo">
			Atenção: ao clicar em consultar, a InfoSimples pode cobrar a consulta do seu saldo.
		</div>

		<div class="barra-acoes">
			<button class="botao" type="submit" name="acao" value="consultar">Consultar PDF sem baixar</button>

			{{if .Reuniao.ID}}
				<button class="botao" type="submit" name="acao" value="consultar_baixar">
					Consultar e salvar PDF na reunião
				</button>
			{{end}}

			<a class="botao secundario" href="/infosimples-car?id={{.Reuniao.ID}}">Voltar ao demonstrativo</a>
			<a class="botao secundario" href="/infosimples-config">Configurar token</a>
			{{if .Reuniao.ID}}
				<a class="botao secundario" href="/investigacao?id={{.Reuniao.ID}}">Voltar à investigação</a>
			{{end}}
		</div>
	</form>
</section>

{{if .Consultou}}
<section class="card">
	<h3>Resumo da resposta</h3>

	<div class="grade">
		<div><strong>Código:</strong><br>{{.Resposta.Code}}</div>
		<div><strong>Mensagem:</strong><br>{{.Resposta.CodeMessage}}</div>
		<div><strong>Cobrável:</strong><br>{{.Resposta.Header.Billable}}</div>
		<div><strong>Preço informado:</strong><br>R$ {{.Resposta.Header.Price}}</div>
		<div><strong>Serviço:</strong><br>{{.Resposta.Header.Service}}</div>
		<div><strong>Itens retornados:</strong><br>{{.Resposta.DataCount}}</div>
	</div>

	{{if .Resposta.Errors}}
		<h4>Erros</h4>
		<ul>
			{{range .Resposta.Errors}}
				<li>{{.}}</li>
			{{end}}
		</ul>
	{{end}}
</section>

<section class="card">
	<h3>JSON bruto retornado</h3>
	<pre style="white-space: pre-wrap; overflow:auto; max-height:600px;">{{.JSONBonito}}</pre>
</section>
{{end}}
`

const infosimplesCARShapefileHTML = `
<section class="card">
	<h2>InfoSimples - CAR / Download Shapefile</h2>

	<p class="pequeno">
		Esta consulta pode consumir crédito. Ela deve retornar um arquivo ZIP com os shapefiles do imóvel rural.
		Depois de salvo, o ZIP será registrado no georreferenciamento da reunião.
	</p>

	<div class="grade">
		<div>
			<strong>Token:</strong><br>
			<span>{{.TokenMascarado}}</span>
		</div>
		<div>
			<strong>Endpoint Shapefile:</strong><br>
			<span class="pequeno">{{.Config.CARDownloadShapefileURL}}</span>
		</div>
	</div>

	{{if .Reuniao.ID}}
		<div class="alerta">
			<strong>Reunião:</strong> {{.Reuniao.Produtor}} —
			{{.Reuniao.Municipio}}/{{.Reuniao.UF}}
		</div>
	{{end}}

	{{if .Erro}}
		<div class="alerta perigo">{{.Erro}}</div>
	{{end}}

	{{if .Sucesso}}
		<div class="alerta sucesso">{{.Sucesso}}</div>
	{{end}}

	<form method="post" class="formulario">
		<label>Número do CAR</label>
		<input type="text" name="car" value="{{.CAR}}" placeholder="Ex: MG-0000000-...">

		<div class="alerta perigo">
			Atenção: ao clicar em consultar, a InfoSimples pode cobrar a consulta do seu saldo.
			Use somente com CAR conferido.
		</div>

		<div class="barra-acoes">
			<button class="botao" type="submit" name="acao" value="consultar">Consultar sem baixar</button>

			{{if .Reuniao.ID}}
				<button class="botao" type="submit" name="acao" value="consultar_baixar">
					Consultar, baixar ZIP e enviar para o mapa
				</button>
			{{end}}

			<a class="botao secundario" href="/infosimples-car?id={{.Reuniao.ID}}">Voltar ao demonstrativo</a>
			<a class="botao secundario" href="/infosimples-car-pdf?id={{.Reuniao.ID}}">Demonstrativo PDF</a>
			{{if .Reuniao.ID}}
				<a class="botao secundario" href="/georreferenciamento?id={{.Reuniao.ID}}">Georreferenciamento</a>
				<a class="botao secundario" href="/mapa-georef?id={{.Reuniao.ID}}">Ver mapa</a>
				<a class="botao secundario" href="/investigacao?id={{.Reuniao.ID}}">Voltar à investigação</a>
			{{end}}
		</div>
	</form>
</section>

{{if .Consultou}}
<section class="card">
	<h3>Resumo da resposta</h3>

	<div class="grade">
		<div><strong>Código:</strong><br>{{.Resposta.Code}}</div>
		<div><strong>Mensagem:</strong><br>{{.Resposta.CodeMessage}}</div>
		<div><strong>Cobrável:</strong><br>{{.Resposta.Header.Billable}}</div>
		<div><strong>Preço informado:</strong><br>R$ {{.Resposta.Header.Price}}</div>
		<div><strong>Serviço:</strong><br>{{.Resposta.Header.Service}}</div>
		<div><strong>Itens retornados:</strong><br>{{.Resposta.DataCount}}</div>
	</div>

	{{if .Resposta.Errors}}
		<h4>Erros</h4>
		<ul>
			{{range .Resposta.Errors}}
				<li>{{.}}</li>
			{{end}}
		</ul>
	{{end}}
</section>

<section class="card">
	<h3>JSON bruto retornado</h3>
	<pre style="white-space: pre-wrap; overflow:auto; max-height:600px;">{{.JSONBonito}}</pre>
</section>
{{end}}
`
