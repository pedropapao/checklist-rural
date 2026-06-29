package main

import (
	"html/template"
	"net/http"
	"strconv"
)

func (app *App) telaRelatorioPreAnalise(w http.ResponseWriter, r *http.Request) {
	idTexto := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idTexto)
	if err != nil || id <= 0 {
		http.Error(w, "ID da reunião inválido", http.StatusBadRequest)
		return
	}

	reuniao, err := app.buscarReuniaoPorID(id)
	if err != nil {
		http.Error(w, "Reunião não encontrada: "+err.Error(), http.StatusNotFound)
		return
	}

	dados := map[string]any{
		"Reuniao":      reuniao,
		"Leitura":      montarLeituraInicial(reuniao),
		"LinhasBB":     sugerirLinhasBB(reuniao),
		"LinhasSicoob": sugerirLinhasSicoob(reuniao),
		"ResumoGeoref": app.resumoArquivosGeorreferenciamento(reuniao.ID),
	}

	tpl := template.Must(template.New("relatorio").Parse(relatorioPreAnaliseHTML))
	tpl.Execute(w, dados)
}

const relatorioPreAnaliseHTML = `
<!DOCTYPE html>
<html lang="pt-br">
<head>
	<meta charset="UTF-8">
	<title>Relatório da Pré-análise</title>

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
			--perigo: #fff1f1;
		}

		* {
			box-sizing: border-box;
		}

		body {
			margin: 0;
			font-family: Arial, Helvetica, sans-serif;
			background: var(--fundo);
			color: var(--texto);
			line-height: 1.45;
		}

		header {
			background: linear-gradient(135deg, var(--verde), var(--verde-escuro));
			color: white;
			padding: 28px 34px;
		}

		header h1 {
			margin: 0;
			font-size: 28px;
		}

		header p {
			margin: 8px 0 0;
			color: #e7f5ec;
		}

		main {
			max-width: 1120px;
			margin: 0 auto;
			padding: 26px;
		}

		.grid {
			display: grid;
			grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
			gap: 16px;
		}

		.card {
			background: var(--card);
			border: 1px solid var(--borda);
			border-radius: 14px;
			padding: 18px;
			margin-bottom: 16px;
			box-shadow: 0 6px 18px rgba(15, 23, 42, 0.05);
		}

		.card.destaque {
			background: #eef8f1;
			border-color: #cfe8d8;
		}

		.card.alerta {
			background: var(--alerta);
			border-color: #f2d28a;
		}

		.card.perigo {
			background: var(--perigo);
			border-color: #f2b8b8;
		}

		h2 {
			color: var(--verde-escuro);
			margin: 0 0 14px;
			font-size: 21px;
		}

		h3 {
			color: var(--verde-escuro);
			margin: 0 0 10px;
			font-size: 17px;
		}

		p {
			margin: 6px 0;
		}

		.pequeno {
			font-size: 13px;
			color: var(--muted);
		}

		.badge {
			display: inline-block;
			border-radius: 999px;
			padding: 4px 9px;
			font-size: 12px;
			font-weight: bold;
			background: #eef6f0;
			color: var(--verde-escuro);
			border: 1px solid #cfe5d6;
		}

		table {
			width: 100%;
			border-collapse: collapse;
			background: white;
			margin-top: 10px;
			border-radius: 10px;
			overflow: hidden;
		}

		th {
			background: #eef6f0;
			color: var(--verde-escuro);
			text-align: left;
			padding: 10px;
			font-size: 13px;
		}

		td {
			border-top: 1px solid var(--borda);
			padding: 10px;
			vertical-align: top;
			font-size: 13px;
		}

		.barra-acoes {
			margin-bottom: 18px;
		}

		.botao {
			display: inline-block;
			background: var(--verde);
			color: white;
			text-decoration: none;
			border-radius: 9px;
			padding: 10px 14px;
			font-weight: bold;
			font-size: 14px;
			margin-right: 6px;
		}

		.botao.secundario {
			background: #eef6f0;
			color: var(--verde-escuro);
			border: 1px solid #cfe5d6;
		}

		ul {
			margin-top: 8px;
		}

		pre {
			white-space: pre-wrap;
			background: #f8fafc;
			border: 1px solid var(--borda);
			border-radius: 10px;
			padding: 12px;
			font-size: 13px;
		}

		@media print {
			body {
				background: white;
			}

			header {
				padding: 20px 24px;
			}

			main {
				padding: 18px;
				max-width: 100%;
			}

			.barra-acoes {
				display: none;
			}

			.card {
				box-shadow: none;
				break-inside: avoid;
			}
		}
	</style>
</head>
<body>
	<header>
		<h1>Relatório da Pré-análise</h1>
		<p>Checklist Rural — documento inicial para organização do projeto de financiamento rural.</p>
	</header>

	<main>
		<div class="barra-acoes">
			<a class="botao secundario" href="/detalhes?id={{.Reuniao.ID}}">Voltar para detalhes</a>
			<a class="botao" href="javascript:window.print()">Imprimir / Salvar PDF</a>
		</div>

		<section class="card destaque">
			<h2>{{.Reuniao.Produtor}}</h2>
			<p>
				<strong>Município/UF:</strong> {{.Reuniao.Municipio}}/{{.Reuniao.UF}}
				&nbsp; | &nbsp;
				<strong>Banco:</strong> {{.Reuniao.Banco}}
				&nbsp; | &nbsp;
				<strong>Data:</strong> {{.Reuniao.CriadoEm}}
			</p>
		</section>

		<section>
			<h2>1. Dados principais</h2>

			<div class="grid">
				<div class="card">
					<h3>Produtor</h3>
					<p><strong>Nome:</strong> {{.Reuniao.Produtor}}</p>
					<p><strong>Telefone:</strong> {{.Reuniao.Telefone}}</p>
					<p><strong>Município/UF:</strong> {{.Reuniao.Municipio}}/{{.Reuniao.UF}}</p>
				</div>

				<div class="card">
					<h3>Projeto</h3>
					<p><strong>Banco:</strong> {{.Reuniao.Banco}}</p>
					<p><strong>Tipo:</strong> {{.Reuniao.TipoProjeto}}</p>
					<p><strong>Atividade:</strong> {{.Reuniao.Atividade}}</p>
					<p><strong>Finalidade:</strong> {{.Reuniao.FinalidadeCredito}}</p>
				</div>

				<div class="card">
					<h3>Enquadramento</h3>
					<p><strong>Renda anual:</strong> R$ {{printf "%.2f" .Reuniao.RendaAnual}}</p>
					<p><strong>Classificação:</strong> <span class="badge">{{.Reuniao.ClassificacaoProdutor}}</span></p>
					<p><strong>CAF:</strong> {{.Reuniao.PossuiCAF}}</p>
					<p><strong>Valor pretendido:</strong> R$ {{printf "%.2f" .Reuniao.ValorPretendido}}</p>
				</div>
			</div>
		</section>

		<section>
			<h2>2. Leitura inicial automática</h2>

			<div class="card">
				<p><strong>Enquadramento:</strong> {{.Leitura.Enquadramento}}</p>
				<p><strong>Caminho sugerido:</strong> {{.Leitura.Caminho}}</p>

				{{if .Leitura.Resumo}}
					<p><strong>Resumo:</strong> {{.Leitura.Resumo}}</p>
				{{end}}
			</div>

			{{if .Leitura.Alertas}}
			<div class="card alerta">
				<h3>Alertas</h3>
				<ul>
					{{range .Leitura.Alertas}}
						<li>{{.}}</li>
					{{end}}
				</ul>
			</div>
			{{end}}

			{{if .Leitura.ProximosPassos}}
			<div class="card">
				<h3>Próximos passos</h3>
				<ul>
					{{range .Leitura.ProximosPassos}}
						<li>{{.}}</li>
					{{end}}
				</ul>
			</div>
			{{end}}
		</section>

		<section>
			<h2>3. Situação inicial informada</h2>

			<div class="grid">
				<div class="card">
					<h3>Cadastro e crédito</h3>
					<p><strong>Cadastro no banco:</strong> {{.Reuniao.CadastroBanco}}</p>
					<p><strong>Financiamento ativo:</strong> {{.Reuniao.FinanciamentoAtivo}}</p>
					<p><strong>Restrição cadastral:</strong> {{.Reuniao.RestricaoCadastral}}</p>
				</div>

				<div class="card">
					<h3>Imóvel</h3>
					<p><strong>Imóvel próprio:</strong> {{.Reuniao.ImovelProprio}}</p>
					<p><strong>Imóvel arrendado:</strong> {{.Reuniao.ImovelArrendado}}</p>
					<p><strong>Tem CAR:</strong> {{.Reuniao.TemCAR}}</p>
				</div>

				<div class="card">
					<h3>Aspectos técnicos</h3>
					<p><strong>Usa água:</strong> {{.Reuniao.UsaAgua}}</p>
					<p><strong>Pecuária:</strong> {{.Reuniao.TemPecuaria}}</p>
					<p><strong>Investimento:</strong> {{.Reuniao.TemInvestimento}}</p>
					<p><strong>Obra:</strong> {{.Reuniao.TemObra}}</p>
					<p><strong>Supressão vegetal:</strong> {{.Reuniao.TemSupressao}}</p>
					<p><strong>ZARC:</strong> {{.Reuniao.PrecisaZARC}}</p>
				</div>
			</div>
		</section>

		<section>
			<h2>4. Investigação do produtor/imóvel</h2>

			<div class="grid">
				<div class="card">
					<h3>CNPJ e endereço</h3>
					<p><strong>CNPJ:</strong> {{if .DadosExternos.CNPJ}}{{.DadosExternos.CNPJ}}{{else}}-{{end}}</p>
					<p><strong>Razão social:</strong> {{if .DadosExternos.RazaoSocial}}{{.DadosExternos.RazaoSocial}}{{else}}-{{end}}</p>
					<p><strong>Situação:</strong> {{if .DadosExternos.SituacaoCadastral}}{{.DadosExternos.SituacaoCadastral}}{{else}}-{{end}}</p>
					<p><strong>Endereço:</strong> {{if .DadosExternos.Endereco}}{{.DadosExternos.Endereco}}{{else}}-{{end}}</p>
				</div>

				<div class="card">
					<h3>Imóvel rural</h3>
					<p><strong>Nome do imóvel:</strong> {{if .DadosExternos.NomeImovel}}{{.DadosExternos.NomeImovel}}{{else}}-{{end}}</p>
					<p><strong>CAR:</strong> {{if .DadosExternos.CAR}}{{.DadosExternos.CAR}}{{else}}-{{end}}</p>
					<p><strong>Situação do CAR:</strong> {{if .DadosExternos.SituacaoCAR}}{{.DadosExternos.SituacaoCAR}}{{else}}-{{end}}</p>
					<p><strong>Código INCRA:</strong> {{if .DadosExternos.CodigoINCRA}}{{.DadosExternos.CodigoINCRA}}{{else}}-{{end}}</p>
					<p><strong>CCIR:</strong> {{if .DadosExternos.CCIR}}{{.DadosExternos.CCIR}}{{else}}-{{end}}</p>
					<p><strong>Matrícula:</strong> {{if .DadosExternos.Matricula}}{{.DadosExternos.Matricula}}{{else}}-{{end}}</p>
				</div>

				<div class="card">
					<h3>Alertas públicos</h3>
					<p><strong>Ibama:</strong> {{if .DadosExternos.EmbargoIBAMA}}{{.DadosExternos.EmbargoIBAMA}}{{else}}não consultado{{end}}</p>
					<p><strong>CEIS:</strong> {{if .DadosExternos.CEISStatus}}{{.DadosExternos.CEISStatus}}{{else}}não consultado{{end}}</p>
					<p><strong>CNEP:</strong> {{if .DadosExternos.CNEPStatus}}{{.DadosExternos.CNEPStatus}}{{else}}não consultado{{end}}</p>
					<p><strong>SNCR:</strong> {{if .DadosExternos.SNCRStatus}}{{.DadosExternos.SNCRStatus}}{{else}}-{{end}}</p>
				</div>
			</div>
		</section>

		<section>
			<h2>5. Georreferenciamento</h2>

			{{if .ResumoGeoref.Total}}
			<div class="card">
				<p><strong>Total de arquivos importados:</strong> {{.ResumoGeoref.Total}}</p>

				<table>
					<thead>
						<tr>
							<th>Tipo</th>
							<th>Arquivo</th>
							<th>Data</th>
							<th>Observação</th>
						</tr>
					</thead>
					<tbody>
						{{range .ResumoGeoref.Arquivos}}
						<tr>
							<td>{{.Tipo}}</td>
							<td>{{.NomeOriginal}}</td>
							<td>{{.CriadoEm}}</td>
							<td>{{.Observacao}}</td>
						</tr>
						{{end}}
					</tbody>
				</table>
			</div>
			{{else}}
			<div class="card alerta">
				<p>Nenhum arquivo de georreferenciamento foi importado até o momento.</p>
			</div>
			{{end}}
		</section>

		<section>
			<h2>6. Linhas sugeridas</h2>

			<div class="grid">
				<div class="card">
					<h3>Banco do Brasil</h3>

					{{if .LinhasBB}}
						<ul>
							{{range .LinhasBB}}
								<li><strong>{{.Nome}}</strong> — {{.Motivo}}</li>
							{{end}}
						</ul>
					{{else}}
						<p class="pequeno">Nenhuma sugestão específica para Banco do Brasil.</p>
					{{end}}
				</div>

				<div class="card">
					<h3>Sicoob</h3>

					{{if .LinhasSicoob}}
						<ul>
							{{range .LinhasSicoob}}
								<li><strong>{{.Nome}}</strong> — {{.Motivo}}</li>
							{{end}}
						</ul>
					{{else}}
						<p class="pequeno">Nenhuma sugestão específica para Sicoob.</p>
					{{end}}
				</div>
			</div>
		</section>

		<section>
			<h2>7. Observações</h2>

			<div class="card">
				{{if .Reuniao.Observacoes}}
					<p>{{.Reuniao.Observacoes}}</p>
				{{else}}
					<p class="pequeno">Nenhuma observação registrada na reunião.</p>
				{{end}}
			</div>

			{{if .DadosExternos.Observacoes}}
			<div class="card">
				<h3>Observações da investigação</h3>
				<pre>{{.DadosExternos.Observacoes}}</pre>
			</div>
			{{end}}
		</section>
	</main>
</body>
</html>
`
