package main

import (
	"net/http"
)

func (app *App) painelSimples(w http.ResponseWriter, r *http.Request) {
	dados := map[string]any{
		"Titulo": "GO MCR",
	}

	tpl := templateMust("painel_simples", painelSimplesHTML)
	tpl.Execute(w, dados)
}

const painelSimplesHTML = `
<section class="home-limpa">

	<div class="home-cabecalho">
		<div>
			<p class="home-tag">GO MCR</p>
			<h1>Checklist Rural</h1>
			<p class="home-subtitulo">
				Organizador particular para triagem, investigação, CAR, georreferenciamento,
				checklist e relatório de projetos rurais.
			</p>
		</div>

		<div class="home-status">
			<strong>Fluxo recomendado</strong>
			<span>1. Consultar produtor</span>
			<span>2. Abrir reunião</span>
			<span>3. Conferir imóvel/CAR</span>
			<span>4. Gerar checklist</span>
		</div>
	</div>

	<div class="home-grade-principal">

		<div class="home-card destaque">
			<h2>Começar atendimento</h2>
			<p>Use esta parte quando chegar um produtor novo.</p>

			<div class="home-acoes">
				<a class="botao grande" href="/infosimples-cpf">Consultar CPF antes da triagem</a>
				<a class="botao secundario grande" href="/nova-reuniao">Nova reunião manual</a>
				<a class="botao secundario grande" href="/consultas-externas">Consultar CNPJ / CEP</a>
			</div>
		</div>

		<div class="home-card">
			<h2>Continuar trabalho</h2>
			<p>Acesse reuniões já criadas e continue a análise.</p>

			<div class="home-atalhos">
				<a href="/reunioes">Reuniões salvas</a>
				<a href="/investigacao">Investigação</a>
				<a href="/georreferenciamento">Georreferenciamento</a>
				<a href="/checklist">Checklist</a>
				<a href="/relatorio">Relatório</a>
			</div>
		</div>

	</div>

	<div class="home-grade-secundaria">

		<div class="home-mini-card">
			<h3>InfoSimples / CAR</h3>
			<p>Consultas automáticas e documentos do imóvel.</p>
			<div class="home-mini-acoes">
				<a href="/infosimples-car">Consultar CAR</a>
				<a href="/infosimples-automatico">Atualizar tudo</a>
				<a href="/infosimples-config">Configurar token</a>
			</div>
		</div>

		<div class="home-mini-card">
			<h3>Mapa e arquivos</h3>
			<p>KML, KMZ, GeoJSON, shapefile e mapa do imóvel.</p>
			<div class="home-mini-acoes">
				<a href="/georreferenciamento">Enviar arquivos</a>
				<a href="/mapa-georef">Abrir mapa</a>
				<a href="/infosimples-car-shapefile">Baixar shapefile</a>
			</div>
		</div>

		<div class="home-mini-card">
			<h3>Sistema</h3>
			<p>Manutenção local do programa.</p>
			<div class="home-mini-acoes">
				<a href="/backups">Backups</a>
				<a href="/pastas">Pastas</a>
				<a href="/configuracoes-api">APIs</a>
			</div>
		</div>

	</div>

</section>

<style>
.home-limpa {
	display: grid;
	gap: 22px;
}

.home-cabecalho {
	display: grid;
	grid-template-columns: 1fr 280px;
	gap: 20px;
	align-items: stretch;
}

.home-cabecalho,
.home-card,
.home-mini-card {
	background: #ffffff;
	border: 1px solid #e5e7eb;
	border-radius: 22px;
	box-shadow: 0 10px 26px rgba(15, 23, 42, 0.08);
}

.home-cabecalho {
	padding: 28px;
}

.home-tag {
	display: inline-block;
	margin: 0 0 10px 0;
	padding: 6px 12px;
	border-radius: 999px;
	background: #dcfce7;
	color: #166534;
	font-size: 13px;
	font-weight: 800;
	letter-spacing: .05em;
}

.home-cabecalho h1 {
	margin: 0;
	font-size: 38px;
	color: #0f172a;
}

.home-subtitulo {
	margin: 10px 0 0 0;
	max-width: 760px;
	color: #475569;
	font-size: 16px;
	line-height: 1.6;
}

.home-status {
	background: #f8fafc;
	border: 1px solid #e2e8f0;
	border-radius: 18px;
	padding: 18px;
	display: grid;
	gap: 8px;
	align-content: center;
}

.home-status strong {
	color: #0f172a;
	margin-bottom: 4px;
}

.home-status span {
	color: #475569;
	font-size: 14px;
}

.home-grade-principal {
	display: grid;
	grid-template-columns: 1.2fr .8fr;
	gap: 20px;
}

.home-card {
	padding: 24px;
}

.home-card.destaque {
	border-color: #bbf7d0;
	background: linear-gradient(180deg, #ffffff 0%, #f0fdf4 100%);
}

.home-card h2,
.home-mini-card h3 {
	margin: 0 0 8px 0;
	color: #0f172a;
}

.home-card p,
.home-mini-card p {
	margin: 0 0 18px 0;
	color: #64748b;
	line-height: 1.5;
}

.home-acoes {
	display: grid;
	gap: 10px;
}

.botao.grande {
	display: block;
	text-align: center;
	padding: 14px 18px;
	font-size: 15px;
}

.home-atalhos {
	display: grid;
	grid-template-columns: 1fr;
	gap: 10px;
}

.home-atalhos a,
.home-mini-acoes a {
	display: block;
	text-decoration: none;
	color: #166534;
	background: #f8fafc;
	border: 1px solid #e2e8f0;
	border-radius: 14px;
	padding: 12px 14px;
	font-weight: 700;
}

.home-atalhos a:hover,
.home-mini-acoes a:hover {
	background: #ecfdf5;
	border-color: #86efac;
}

.home-grade-secundaria {
	display: grid;
	grid-template-columns: repeat(3, 1fr);
	gap: 20px;
}

.home-mini-card {
	padding: 20px;
}

.home-mini-acoes {
	display: grid;
	gap: 8px;
}

.home-mini-acoes a {
	font-size: 14px;
	padding: 10px 12px;
}

@media (max-width: 900px) {
	.home-cabecalho,
	.home-grade-principal,
	.home-grade-secundaria {
		grid-template-columns: 1fr;
	}
}
</style>
`
