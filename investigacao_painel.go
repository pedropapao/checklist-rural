package main

import (
	"html/template"
	"net/http"
)

func (app *App) painelInvestigacao(w http.ResponseWriter, r *http.Request) {
	reuniao, err := app.pegarReuniaoDaURL(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dadosExternos := app.buscarDadosExternosPorReuniao(reuniao.ID)
	resumoGeoref := app.resumoArquivosGeorreferenciamento(reuniao.ID)

	dados := map[string]any{
		"Titulo":       "Investigação",
		"Reuniao":      reuniao,
		"Dados":        dadosExternos,
		"ResumoGeoref": resumoGeoref,
	}

	tpl := template.Must(template.New("investigacao").Parse(htmlBase(investigacaoPainelHTML)))
	tpl.Execute(w, dados)
}

const investigacaoPainelHTML = `
<h2>Investigação do produtor e imóvel</h2>

<div class="barra-acoes">
	<a class="botao secundario" href="/detalhes?id={{.Reuniao.ID}}">Voltar para detalhes</a>
	<a class="botao" href="/investigacao-editar?id={{.Reuniao.ID}}">Editar investigação rápida</a>
	<a class="botao secundario" href="/dados-externos-reuniao?id={{.Reuniao.ID}}">Editar investigação completa</a>
	<a class="botao secundario" href="/georreferenciamento?id={{.Reuniao.ID}}">Georreferenciamento</a>
		<a class="botao" href="/infosimples-car?id={{.Reuniao.ID}}">Consultar CAR InfoSimples</a>
</div>

<div class="grid">
	<div class="card destaque">
		<h3>Produtor</h3>
		<p class="pequeno">Dados principais da reunião.</p>

		<p><strong>Nome:</strong> {{.Reuniao.Produtor}}</p>
		<p><strong>Município/UF:</strong> {{.Reuniao.Municipio}}/{{.Reuniao.UF}}</p>
		<p><strong>Banco:</strong> {{.Reuniao.Banco}}</p>
		<p><strong>Atividade:</strong> {{.Reuniao.Atividade}}</p>
	</div>

	<div class="card">
		<h3>CNPJ e endereço</h3>
		<p class="pequeno">Informações cadastrais e localização básica.</p>

		<p><strong>CNPJ:</strong> {{if .Dados.CNPJ}}{{.Dados.CNPJ}}{{else}}<span class="badge">não informado</span>{{end}}</p>
		<p><strong>Razão social:</strong> {{if .Dados.RazaoSocial}}{{.Dados.RazaoSocial}}{{else}}-{{end}}</p>
		<p><strong>Situação cadastral:</strong> {{if .Dados.SituacaoCadastral}}{{.Dados.SituacaoCadastral}}{{else}}-{{end}}</p>
		<p><strong>CEP:</strong> {{if .Dados.CEP}}{{.Dados.CEP}}{{else}}-{{end}}</p>
		<p><strong>Endereço:</strong> {{if .Dados.Endereco}}{{.Dados.Endereco}}{{else}}-{{end}}</p>
	</div>

	<div class="card">
		<h3>Imóvel rural</h3>
		<p class="pequeno">CAR, matrícula, Incra e informações do imóvel.</p>

		<p><strong>Nome do imóvel:</strong> {{if .Dados.NomeImovel}}{{.Dados.NomeImovel}}{{else}}-{{end}}</p>
		<p><strong>CAR:</strong> {{if .Dados.CAR}}{{.Dados.CAR}}{{else}}<span class="badge">não informado</span>{{end}}</p>
		<p><strong>Situação do CAR:</strong> {{if .Dados.SituacaoCAR}}{{.Dados.SituacaoCAR}}{{else}}-{{end}}</p>
		<p><strong>Código INCRA:</strong> {{if .Dados.CodigoINCRA}}{{.Dados.CodigoINCRA}}{{else}}-{{end}}</p>
		<p><strong>CCIR:</strong> {{if .Dados.CCIR}}{{.Dados.CCIR}}{{else}}-{{end}}</p>
		<p><strong>Matrícula:</strong> {{if .Dados.Matricula}}{{.Dados.Matricula}}{{else}}-{{end}}</p>
	</div>
</div>

<div class="grid">
	<div class="card">
		<h3>Alertas públicos</h3>
		<p class="pequeno">Situação em bases externas consultadas.</p>

		<p><strong>Ibama:</strong> {{if .Dados.EmbargoIBAMA}}{{.Dados.EmbargoIBAMA}}{{else}}não consultado{{end}}</p>
		<p><strong>CEIS:</strong> {{if .Dados.CEISStatus}}{{.Dados.CEISStatus}}{{else}}não consultado{{end}}</p>
		<p><strong>CNEP:</strong> {{if .Dados.CNEPStatus}}{{.Dados.CNEPStatus}}{{else}}não consultado{{end}}</p>
		<p><strong>SNCR:</strong> {{if .Dados.SNCRStatus}}{{.Dados.SNCRStatus}}{{else}}não informado{{end}}</p>
	</div>

	<div class="card">
		<h3>Georreferenciamento</h3>
		<p class="pequeno">Arquivos e localização espacial da área do projeto.</p>

		{{if .ResumoGeoref.Total}}
			<p><strong>Arquivos importados:</strong> {{.ResumoGeoref.Total}}</p>

			<table>
				<thead>
					<tr>
						<th>Tipo</th>
						<th>Arquivo</th>
					</tr>
				</thead>
				<tbody>
					{{range .ResumoGeoref.Arquivos}}
					<tr>
						<td>{{.Tipo}}</td>
						<td>{{.NomeOriginal}}</td>
					</tr>
					{{end}}
				</tbody>
			</table>
		{{else}}
			<p class="pequeno">Nenhum arquivo georreferenciado importado ainda.</p>
		{{end}}

		<div class="barra-acoes">
			<a class="botao" href="/georreferenciamento?id={{.Reuniao.ID}}">Abrir georreferenciamento</a>
		</div>
	</div>

	<div class="card">
		<h3>Fontes e atualização</h3>
		<p class="pequeno">Controle da origem das informações.</p>

		<p><strong>Fonte:</strong> {{if .Dados.FonteConsulta}}{{.Dados.FonteConsulta}}{{else}}-{{end}}</p>
		<p><strong>Última investigação:</strong> {{if .Dados.UltimaInvestigacaoOnline}}{{.Dados.UltimaInvestigacaoOnline}}{{else}}-{{end}}</p>

		{{if .Dados.LinkFonte}}
			<p><a class="botao secundario" href="{{.Dados.LinkFonte}}" target="_blank">Abrir fonte</a></p>
		{{end}}
	</div>
</div>

<div class="card">
	<h3>O que falta conferir?</h3>

	<ul>
		{{if not .Dados.CNPJ}}<li>Informar ou consultar CNPJ.</li>{{end}}
		{{if not .Dados.CAR}}<li>Informar CAR ou anexar documento que contenha o CAR.</li>{{end}}
		{{if not .Dados.CodigoINCRA}}<li>Informar código INCRA, se aplicável.</li>{{end}}
		{{if not .Dados.CCIR}}<li>Informar CCIR, se disponível.</li>{{end}}
		{{if not .Dados.Matricula}}<li>Informar matrícula do imóvel.</li>{{end}}
		{{if not .ResumoGeoref.Total}}<li>Importar KML, GeoJSON, croqui, planta, PDF ou arquivo de localização.</li>{{end}}
	</ul>
</div>

<div class="card">
	<h3>Observações da investigação</h3>

	{{if .Dados.Observacoes}}
		<pre>{{.Dados.Observacoes}}</pre>
	{{else}}
		<p class="pequeno">Nenhuma observação registrada.</p>
	{{end}}
</div>
`
