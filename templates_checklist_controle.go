package main

const checklistControleHTML = `
<h2>Controle do checklist</h2>

<div class="card">
	<strong>Produtor:</strong> {{.Reuniao.Produtor}}<br>
	<strong>Banco:</strong> {{.Reuniao.Banco}}<br>
	<strong>Tipo:</strong> {{.Reuniao.TipoProjeto}}<br>
	<strong>Atividade:</strong> {{.Reuniao.Atividade}}<br>
</div>

<div class="card">
	<h3>Resumo do andamento</h3>

	<div class="grid">
		<div>
			<strong>Total de itens:</strong><br>
			{{.Resumo.Total}}
		</div>

		<div>
			<strong>Concluído:</strong><br>
			{{.Resumo.PercentualConcluido}}%
		</div>

		<div>
			<strong>Pendentes:</strong><br>
			{{.Resumo.Pendentes}}
		</div>

		<div>
			<strong>Recebidos:</strong><br>
			{{.Resumo.Recebidos}}
		</div>

		<div>
			<strong>Não se aplica:</strong><br>
			{{.Resumo.NaoSeAplica}}
		</div>
	</div>
</div>

<a class="botao secundario" href="/detalhes?id={{.Reuniao.ID}}">Voltar para detalhes</a>
<a class="botao" href="/checklist?id={{.Reuniao.ID}}">Ver checklist em texto</a>
<a class="botao alerta" href="/whatsapp?id={{.Reuniao.ID}}">Resumo WhatsApp</a>
<a class="botao alerta" href="/whatsapp-pendencias?id={{.Reuniao.ID}}">WhatsApp pendências</a>
<a class="botao" href="/exportar-checklist-controle-txt?id={{.Reuniao.ID}}">Exportar controle TXT</a>

<form method="POST" action="/salvar-itens-checklist">
	<input type="hidden" name="reuniao_id" value="{{.Reuniao.ID}}">

	<table>
		<tr>
			<th>Grupo</th>
			<th>Item</th>
			<th>Status</th>
			<th>Observação</th>
		</tr>

		{{range .Itens}}
		<tr>
			<td>{{.Grupo}}</td>
			<td>{{.Item}}</td>
			<td>
				<select name="status_{{.ID}}">
					<option value="Pendente" {{if eq .Status "Pendente"}}selected{{end}}>Pendente</option>
					<option value="Recebido" {{if eq .Status "Recebido"}}selected{{end}}>Recebido</option>
					<option value="Não se aplica" {{if eq .Status "Não se aplica"}}selected{{end}}>Não se aplica</option>
				</select>
			</td>
			<td>
				<textarea name="observacao_{{.ID}}">{{.Observacao}}</textarea>
			</td>
		</tr>
		{{else}}
		<tr>
			<td colspan="4">Nenhum item de checklist gerado para esta reunião.</td>
		</tr>
		{{end}}
	</table>

	<br>
	<button type="submit">Salvar checklist</button>
</form>
`
