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
	}

	tpl := template.Must(template.New("relatorio").Parse(relatorioPreAnaliseHTML))
	tpl.Execute(w, dados)
}

const relatorioPreAnaliseHTML = `
<!DOCTYPE html>
<html lang="pt-br">
<head>
	<meta charset="UTF-8">
	<title>Relatório de Pré-análise Rural</title>
	<style>
		body {
			font-family: Arial, sans-serif;
			background: #f3f4f6;
			color: #111827;
			margin: 0;
			padding: 24px;
		}

		.pagina {
			max-width: 900px;
			margin: 0 auto;
			background: white;
			padding: 32px;
			border-radius: 12px;
			box-shadow: 0 8px 24px rgba(0,0,0,0.08);
		}

		.topo {
			border-bottom: 3px solid #166534;
			padding-bottom: 16px;
			margin-bottom: 24px;
		}

		h1 {
			margin: 0;
			color: #166534;
			font-size: 26px;
		}

		h2 {
			color: #166534;
			border-bottom: 1px solid #d1d5db;
			padding-bottom: 6px;
			margin-top: 28px;
			font-size: 19px;
		}

		h3 {
			margin-bottom: 6px;
			color: #1f2937;
			font-size: 16px;
		}

		.subtitulo {
			color: #4b5563;
			margin-top: 6px;
		}

		.grid {
			display: grid;
			grid-template-columns: 1fr 1fr;
			gap: 12px;
		}

		.campo {
			border: 1px solid #e5e7eb;
			border-radius: 8px;
			padding: 10px 12px;
			background: #fafafa;
		}

		.label {
			font-size: 12px;
			color: #6b7280;
			display: block;
			margin-bottom: 4px;
		}

		.valor {
			font-weight: bold;
		}

		.card {
			border: 1px solid #e5e7eb;
			border-radius: 10px;
			padding: 14px;
			margin: 10px 0;
			background: #ffffff;
		}

		.alerta {
			background: #fffbeb;
			border-color: #facc15;
		}

		.linha {
			background: #f0fdf4;
			border-color: #86efac;
		}

		ul {
			margin-top: 8px;
		}

		li {
			margin-bottom: 6px;
		}

		.botoes {
			max-width: 900px;
			margin: 0 auto 16px auto;
			display: flex;
			gap: 10px;
		}

		.botao {
			display: inline-block;
			background: #166534;
			color: white;
			padding: 10px 14px;
			border-radius: 8px;
			text-decoration: none;
			border: none;
			cursor: pointer;
			font-weight: bold;
		}

		.botao.secundario {
			background: #374151;
		}

		.rodape {
			margin-top: 30px;
			font-size: 12px;
			color: #6b7280;
			border-top: 1px solid #e5e7eb;
			padding-top: 12px;
		}

		@media print {
			body {
				background: white;
				padding: 0;
			}

			.botoes {
				display: none;
			}

			.pagina {
				box-shadow: none;
				border-radius: 0;
				padding: 12mm;
				max-width: none;
			}

			.card, .campo {
				break-inside: avoid;
			}
		}
	</style>
</head>
<body>

<div class="botoes">
	<a class="botao secundario" href="/reunioes">Voltar</a>
	<button class="botao" onclick="window.print()">Imprimir / salvar PDF</button>
</div>

<div class="pagina">
	<div class="topo">
		<h1>Relatório de Pré-análise Rural</h1>
		<p class="subtitulo">
			Resumo interno gerado a partir da entrevista inicial, classificação do produtor e leitura automática do caso.
		</p>
	</div>

	<h2>1. Dados do produtor</h2>

	<div class="grid">
		<div class="campo">
			<span class="label">Produtor</span>
			<span class="valor">{{.Reuniao.Produtor}}</span>
		</div>

		<div class="campo">
			<span class="label">Telefone</span>
			<span class="valor">{{.Reuniao.Telefone}}</span>
		</div>

		<div class="campo">
			<span class="label">Município/UF</span>
			<span class="valor">{{.Reuniao.Municipio}}/{{.Reuniao.UF}}</span>
		</div>

		<div class="campo">
			<span class="label">Banco/Cooperativa</span>
			<span class="valor">{{.Reuniao.Banco}}</span>
		</div>
	</div>

	<h2>2. Dados do projeto</h2>

	<div class="grid">
		<div class="campo">
			<span class="label">Tipo de projeto</span>
			<span class="valor">{{.Reuniao.TipoProjeto}}</span>
		</div>

		<div class="campo">
			<span class="label">Atividade</span>
			<span class="valor">{{.Reuniao.Atividade}}</span>
		</div>

		<div class="campo">
			<span class="label">Finalidade principal</span>
			<span class="valor">{{.Reuniao.FinalidadeCredito}}</span>
		</div>

		<div class="campo">
			<span class="label">Valor pretendido</span>
			<span class="valor">R$ {{printf "%.2f" .Reuniao.ValorPretendido}}</span>
		</div>
	</div>

	<h2>3. Enquadramento inicial</h2>

	<div class="grid">
		<div class="campo">
			<span class="label">RBA informada</span>
			<span class="valor">R$ {{printf "%.2f" .Reuniao.RendaAnual}}</span>
		</div>

		<div class="campo">
			<span class="label">Classificação pela RBA</span>
			<span class="valor">{{.Reuniao.ClassificacaoProdutor}}</span>
		</div>

		<div class="campo">
			<span class="label">Possui CAF?</span>
			<span class="valor">{{.Reuniao.PossuiCAF}}</span>
		</div>

		<div class="campo">
			<span class="label">Possui orçamento/proposta?</span>
			<span class="valor">{{.Reuniao.TemOrcamento}}</span>
		</div>
	</div>

	<h2>4. Pré-análise operacional</h2>

	<div class="grid">
		<div class="campo">
			<span class="label">Cadastro no banco/cooperativa</span>
			<span class="valor">{{.Reuniao.CadastroBanco}}</span>
		</div>

		<div class="campo">
			<span class="label">Financiamento ativo</span>
			<span class="valor">{{.Reuniao.FinanciamentoAtivo}}</span>
		</div>

		<div class="campo">
			<span class="label">Restrição cadastral</span>
			<span class="valor">{{.Reuniao.RestricaoCadastral}}</span>
		</div>

		<div class="campo">
			<span class="label">Imóvel próprio</span>
			<span class="valor">{{.Reuniao.ImovelProprio}}</span>
		</div>

		<div class="campo">
			<span class="label">Arrendamento/parceria/comodato</span>
			<span class="valor">{{.Reuniao.ImovelArrendado}}</span>
		</div>

		<div class="campo">
			<span class="label">Possui CAR</span>
			<span class="valor">{{.Reuniao.TemCAR}}</span>
		</div>

		<div class="campo">
			<span class="label">Uso de água/irrigação</span>
			<span class="valor">{{.Reuniao.UsaAgua}}</span>
		</div>

		<div class="campo">
			<span class="label">Supressão/abertura de área</span>
			<span class="valor">{{.Reuniao.TemSupressao}}</span>
		</div>
	</div>

	<h2>5. Leitura inicial automática</h2>

	<div class="card">
		<h3>Enquadramento provável</h3>
		<p>{{.Leitura.Enquadramento}}</p>
	</div>

	<div class="card">
		<h3>Caminho inicial</h3>
		<p>{{.Leitura.Caminho}}</p>
	</div>

	{{if .Leitura.Resumo}}
	<div class="card">
		<h3>Resumo da demanda</h3>
		<p>{{.Leitura.Resumo}}</p>
	</div>
	{{end}}

	<div class="card alerta">
		<h3>Atenções identificadas</h3>
		<ul>
			{{range .Leitura.Alertas}}
			<li>{{.}}</li>
			{{end}}
		</ul>
	</div>

	<div class="card">
		<h3>Próximos passos sugeridos</h3>
		<ul>
			{{range .Leitura.ProximosPassos}}
			<li>{{.}}</li>
			{{end}}
		</ul>
	</div>

	{{if .LinhasBB}}
	<h2>6. Possíveis linhas BB</h2>

	{{range .LinhasBB}}
	<div class="card linha">
		<h3>{{.Nome}}</h3>
		<p><strong>Motivo:</strong> {{.Motivo}}</p>
		<p><strong>Atenção:</strong> {{.Atencao}}</p>
	</div>
	{{end}}
	{{end}}

	{{if .LinhasSicoob}}
	<h2>6. Possíveis caminhos Sicoob</h2>

	{{range .LinhasSicoob}}
	<div class="card linha">
		<h3>{{.Nome}}</h3>
		<p><strong>Motivo:</strong> {{.Motivo}}</p>
		<p><strong>Atenção:</strong> {{.Atencao}}</p>
	</div>
	{{end}}
	{{end}}

	<h2>7. Observações registradas</h2>

	<div class="card">
		<p>{{.Reuniao.Observacoes}}</p>
	</div>

	<div class="rodape">
		Este relatório é uma ferramenta interna de organização da pré-análise. Não substitui análise cadastral, enquadramento oficial, normas do MCR, política de crédito do banco/cooperativa ou conferência documental.
	</div>
</div>

</body>
</html>
`
