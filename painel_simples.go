package main

import (
	"html/template"
	"net/http"
)

func (app *App) painelSimples(w http.ResponseWriter, r *http.Request) {
	dados := map[string]any{
		"Titulo": "Checklist Rural",
	}

	tpl := template.Must(template.New("painel").Parse(painelSimplesHTML))
	tpl.Execute(w, dados)
}

const painelSimplesHTML = `
<!DOCTYPE html>
<html lang="pt-br">
<head>
	<meta charset="UTF-8">
	<title>{{.Titulo}}</title>
	<style>
		:root {
			--verde: #1f7a4d;
			--verde-escuro: #145c38;
			--fundo: #f4f7f3;
			--card: #ffffff;
			--texto: #1f2933;
			--muted: #6b7280;
			--borda: #d9e2dc;
			--alerta: #fff7e6;
		}

		* {
			box-sizing: border-box;
		}

		body {
			margin: 0;
			font-family: Arial, Helvetica, sans-serif;
			background: var(--fundo);
			color: var(--texto);
		}

		header {
			background: linear-gradient(135deg, var(--verde), var(--verde-escuro));
			color: white;
			padding: 28px 34px;
		}

		header h1 {
			margin: 0;
			font-size: 30px;
		}

		header p {
			margin: 8px 0 0;
			color: #e7f5ec;
			font-size: 15px;
		}

		main {
			max-width: 1180px;
			margin: 0 auto;
			padding: 28px;
		}

		.secao {
			margin-bottom: 28px;
		}

		.secao h2 {
			font-size: 20px;
			margin: 0 0 14px;
		}

		.grid {
			display: grid;
			grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
			gap: 18px;
		}

		.card {
			background: var(--card);
			border: 1px solid var(--borda);
			border-radius: 16px;
			padding: 22px;
			box-shadow: 0 8px 22px rgba(15, 23, 42, 0.06);
		}

		.card h3 {
			margin: 0 0 8px;
			font-size: 18px;
			color: var(--verde-escuro);
		}

		.card p {
			margin: 0 0 18px;
			color: var(--muted);
			line-height: 1.45;
			font-size: 14px;
		}

		.botao {
			display: inline-block;
			background: var(--verde);
			color: white;
			text-decoration: none;
			border-radius: 10px;
			padding: 11px 15px;
			font-weight: bold;
			font-size: 14px;
		}

		.botao:hover {
			background: var(--verde-escuro);
		}

		.botao.secundario {
			background: #eef6f0;
			color: var(--verde-escuro);
			border: 1px solid #cfe5d6;
		}

		.botao.secundario:hover {
			background: #dcefe4;
		}

		.aviso {
			background: var(--alerta);
			border: 1px solid #f2d28a;
			border-radius: 14px;
			padding: 16px 18px;
			color: #6b4e00;
			line-height: 1.45;
		}

		.atalhos {
			display: flex;
			flex-wrap: wrap;
			gap: 10px;
			margin-top: 16px;
		}

		footer {
			text-align: center;
			color: var(--muted);
			padding: 20px;
			font-size: 13px;
		}
	</style>
</head>
<body>
	<header>
		<h1>Checklist Rural</h1>
		<p>Sistema pessoal para organizar reunião, investigação, pré-análise, checklist e relatório de projeto rural.</p>
	</header>

	<main>
		<section class="secao">
			<h2>Começar</h2>

			<div class="grid">
				<div class="card">
					<h3>Nova reunião</h3>
					<p>Cadastre produtor, banco, atividade, finalidade e dados iniciais do projeto.</p>
					<a class="botao" href="/nova-reuniao">Criar reunião</a>
				</div>

				<div class="card">
					<h3>Reuniões salvas</h3>
					<p>Abra reuniões antigas, continue análise, gere checklist ou relatório.</p>
					<a class="botao" href="/reunioes">Ver reuniões</a>
				</div>

				<div class="card">
					<h3>Investigação rápida</h3>
					<p>Use consultas externas, CNPJ, CEP, CAR, Ibama e georreferenciamento.</p>
					<a class="botao secundario" href="/reunioes">Escolher reunião</a>
				</div>
			</div>
		</section>

		<section class="secao">
			<h2>Ferramentas</h2>

			<div class="grid">
				<div class="card">
					<h3>Consultas externas</h3>
					<p>Área de apoio para consultas de CNPJ, CEP e outras fontes públicas.</p>
					<a class="botao secundario" href="/consultas-externas">Abrir consultas</a>
				</div>

				<div class="card">
					<h3>Configurações das APIs</h3>
					<p>Configure tokens, URLs do SICAR, Portal da Transparência e integrações.</p>
					<a class="botao secundario" href="/configuracoes-api">Abrir configurações</a>
				</div>

				<div class="card">
					<h3>Organização</h3>
					<p>Acesse pastas, arquivos, backups e documentos gerados pelo sistema.</p>
					<div class="atalhos">
						<a class="botao secundario" href="/pastas">Pastas</a>
						<a class="botao secundario" href="/backup">Backup</a>
					</div>
				</div>
			</div>
		</section>

		<section class="secao">
			<div class="aviso">
				<strong>Fluxo recomendado:</strong>
				crie a reunião, faça a investigação do produtor/imóvel, importe documentos/georreferenciamento,
				revise o checklist e só depois gere o relatório da pré-análise.
			</div>
		</section>
	</main>

	<footer>
		Checklist Rural — uso local e pessoal
	</footer>
</body>
</html>
`
