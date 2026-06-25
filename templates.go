package main

func htmlBase(conteudo string) string {
	return `
<!DOCTYPE html>
<html lang="pt-BR">
<head>
	<meta charset="UTF-8">
	<title>{{.Titulo}}</title>
	<style>
		body {
			font-family: Arial, sans-serif;
			background: #f4f6f5;
			margin: 0;
			padding: 0;
			color: #222;
		}

		header {
			background: #1f5f3b;
			color: white;
			padding: 20px;
		}

		nav {
			background: #17472d;
			padding: 10px 20px;
		}

		nav a {
			background: transparent;
			color: white;
			padding: 8px 12px;
			border-radius: 5px;
			text-decoration: none;
			margin-right: 8px;
		}

		nav a:hover {
			background: rgba(255,255,255,0.18);
		}

		main {
			max-width: 1150px;
			margin: 30px auto;
			background: white;
			padding: 25px;
			border-radius: 10px;
			box-shadow: 0 2px 8px rgba(0,0,0,0.12);
		}

		h1, h2 {
			margin-top: 0;
		}

		a.botao, button {
			background: #1f5f3b;
			color: white;
			padding: 10px 14px;
			border-radius: 6px;
			text-decoration: none;
			border: none;
			cursor: pointer;
			display: inline-block;
			margin: 5px 5px 5px 0;
			font-size: 14px;
		}

		a.secundario {
			background: #555;
		}

		a.alerta {
			background: #8a5a00;
		}

		a.perigo, button.perigo {
			background: #9b1c1c;
		}

		label {
			display: block;
			margin-top: 14px;
			font-weight: bold;
		}

		input, select, textarea {
			width: 100%;
			padding: 10px;
			margin-top: 5px;
			border: 1px solid #ccc;
			border-radius: 6px;
			box-sizing: border-box;
			font-size: 15px;
		}

		input[type="checkbox"] {
			width: auto;
			margin-right: 8px;
		}

		.check {
			display: block;
			font-weight: normal;
			background: #f7f7f7;
			padding: 8px;
			border-radius: 6px;
			margin-top: 8px;
		}

		textarea {
			min-height: 90px;
		}

		table {
			width: 100%;
			border-collapse: collapse;
			margin-top: 20px;
			font-size: 14px;
		}

		th, td {
			border: 1px solid #ddd;
			padding: 8px;
			text-align: left;
			vertical-align: top;
		}

		th {
			background: #e9f2ed;
		}

		.card {
			background: #eef6f1;
			padding: 18px;
			border-radius: 8px;
			margin-bottom: 15px;
		}

		.card-perigo {
			background: #fff1f1;
			border: 1px solid #d08a8a;
			padding: 18px;
			border-radius: 8px;
			margin-bottom: 15px;
		}

		.grid {
			display: grid;
			grid-template-columns: 1fr 1fr;
			gap: 12px;
		}

		pre {
			background: #f7f7f7;
			padding: 18px;
			border-radius: 8px;
			white-space: pre-wrap;
			font-family: Consolas, monospace;
			font-size: 14px;
			line-height: 1.45;
		}

		.pequeno {
			color: #666;
			font-size: 13px;
		}
	</style>
</head>
<body>
	<header>
		<h1>Checklist Rural</h1>
		<p>Atendimento, triagem e checklist para reunião com produtor</p>
	</header>

	<nav>
		<a href="/">Início</a>
		<a href="/nova-reuniao">Nova reunião</a>
		<a href="/reunioes">Reuniões salvas</a>
		<a href="/backups">Backups</a>
		<a href="/pastas">Pastas</a>
	</nav>

	<main>
		` + conteudo + `
	</main>
</body>
</html>
`
}

const inicioHTML = `
<h2>Painel inicial</h2>

<div class="card">
	<p>Este sistema salva as reuniões em banco SQLite local no seu computador.</p>
	<p>Use para registrar atendimento, triagem, checklist automático e resumo para WhatsApp.</p>
</div>

<a class="botao" href="/nova-reuniao">Nova reunião</a>
<a class="botao secundario" href="/reunioes">Ver reuniões salvas</a>
<a class="botao" href="/backups">Backups</a>
<a class="botao" href="/pastas">Pastas</a>
`

const novaReuniaoHTML = `
<h2>Nova reunião com produtor</h2>

<form method="POST" action="/salvar-reuniao">
	<h3>Dados principais</h3>

	<label>Nome do produtor</label>
	<input type="text" name="produtor" required>

	<label>Telefone / WhatsApp</label>
	<input type="text" name="telefone">

	<label>Município</label>
	<input type="text" name="municipio">

	<label>UF</label>
	<input type="text" name="uf" maxlength="2">

	<label>Banco pretendido</label>
	<select name="banco">
		<option value="Ainda não definido">Ainda não definido</option>
		<option value="Banco do Brasil">Banco do Brasil</option>
		<option value="Sicoob">Sicoob</option>
	</select>

	<label>Tipo de projeto</label>
	<select name="tipo_projeto">
		<option value="Ainda não definido">Ainda não definido</option>
		<option value="Custeio agrícola">Custeio agrícola</option>
		<option value="Custeio pecuário">Custeio pecuário</option>
		<option value="Investimento">Investimento</option>
		<option value="Comercialização">Comercialização</option>
		<option value="Industrialização">Industrialização</option>
	</select>

	<label>Atividade principal</label>
	<select name="atividade">
		<option value="Ainda não definida">Ainda não definida</option>
		<option value="Agrícola">Agrícola</option>
		<option value="Pecuária de corte">Pecuária de corte</option>
		<option value="Pecuária de leite">Pecuária de leite</option>
		<option value="Irrigação">Irrigação</option>
		<option value="Máquinas e equipamentos">Máquinas e equipamentos</option>
		<option value="Obras e benfeitorias">Obras e benfeitorias</option>
		<option value="Outro">Outro</option>
	</select>

	<h3>Classificação do produtor</h3>

	<label>Receita Bruta Agropecuária Anual - RBA</label>
	<input type="text" name="renda_anual" value="{{printf "%.2f" .Reuniao.RendaAnual}}" placeholder="Ex: 500000,00">

	<p class="pequeno">
		Informe a renda bruta agropecuária anual aproximada para classificar o produtor.
	</p>

	<h3>Triagem guiada</h3>

<div class="card">
	<h4>1. Situação cadastral do produtor</h4>

	<label class="check">
		<input type="checkbox" name="cadastro_banco" value="sim">
		O produtor já tem cadastro no banco pretendido?
	</label>

	<label class="check">
		<input type="checkbox" name="financiamento_ativo" value="sim">
		O produtor possui financiamento rural ativo?
	</label>

	<label class="check">
		<input type="checkbox" name="restricao_cadastral" value="sim">
		O produtor informou alguma restrição cadastral conhecida?
	</label>
</div>

<div class="card">
	<h4>2. Situação da propriedade</h4>

	<label class="check">
		<input type="checkbox" name="imovel_proprio" value="sim">
		A área do projeto é imóvel próprio?
	</label>

	<label class="check">
		<input type="checkbox" name="imovel_arrendado" value="sim">
		A área é arrendada, parceria ou comodato?
	</label>

	<label class="check">
		<input type="checkbox" name="tem_car" value="sim">
		A propriedade possui CAR?
	</label>
</div>

<div class="card">
	<h4>3. Situação ambiental e uso da área</h4>

	<label class="check">
		<input type="checkbox" name="usa_agua" value="sim">
		O projeto usa água, irrigação, poço, represa ou pecuária intensiva?
	</label>

	<label class="check">
		<input type="checkbox" name="tem_supressao" value="sim">
		Haverá supressão vegetal, abertura de área ou limpeza de vegetação nativa?
	</label>
</div>

<div class="card">
	<h4>4. Tipo do projeto pretendido</h4>

	<label class="check">
		<input type="checkbox" name="tem_pecuaria" value="sim">
		O projeto envolve pecuária?
	</label>

	<label class="check">
		<input type="checkbox" name="tem_investimento" value="sim">
		O projeto envolve investimento?
	</label>

	<label class="check">
		<input type="checkbox" name="tem_obra" value="sim">
		O projeto possui obra, construção, reforma ou benfeitoria?
	</label>

	<label class="check">
		<input type="checkbox" name="precisa_zarc" value="sim">
		O projeto precisa validar ZARC?
	</label>
</div>

<label>Observações iniciais</label>
	<textarea name="observacoes"></textarea>

	<br><br>
	<button type="submit">Salvar reunião</button>
	<a class="botao secundario" href="/">Voltar</a>
</form>
`

const reunioesHTML = `
<h2>Reuniões salvas</h2>

<a class="botao" href="/nova-reuniao">Nova reunião</a>
<a class="botao secundario" href="/">Início</a>

<div class="card">
	<h3>Filtros</h3>

	<form method="GET" action="/reunioes">
		<label>Buscar por produtor, telefone, município, banco, projeto ou observação</label>
		<input type="text" name="busca" value="{{.Filtros.Busca}}" placeholder="Digite aqui para buscar">

		<div class="grid">
			<div>
				<label>Banco</label>
				<select name="banco">
					<option value="" {{if eq .Filtros.Banco ""}}selected{{end}}>Todos</option>
					<option value="Banco do Brasil" {{if eq .Filtros.Banco "Banco do Brasil"}}selected{{end}}>Banco do Brasil</option>
					<option value="Sicoob" {{if eq .Filtros.Banco "Sicoob"}}selected{{end}}>Sicoob</option>
					<option value="Ainda não definido" {{if eq .Filtros.Banco "Ainda não definido"}}selected{{end}}>Ainda não definido</option>
				</select>
			</div>

			<div>
				<label>Situação do checklist</label>
				<select name="situacao">
					<option value="" {{if eq .Filtros.Situacao ""}}selected{{end}}>Todas</option>
					<option value="pendentes" {{if eq .Filtros.Situacao "pendentes"}}selected{{end}}>Com pendências</option>
					<option value="andamento" {{if eq .Filtros.Situacao "andamento"}}selected{{end}}>Em andamento</option>
					<option value="concluidas" {{if eq .Filtros.Situacao "concluidas"}}selected{{end}}>Concluídas</option>
				</select>
			</div>
		</div>

		<br>
		<button type="submit">Filtrar</button>
		<a class="botao secundario" href="/reunioes">Limpar filtros</a>
	</form>
</div>

<table>
	<tr>
		<th>ID</th>
		<th>Produtor</th>
		<th>Contato</th>
		<th>Projeto</th>
		<th>Triagem</th>
		<th>Andamento</th>
		<th>Data</th>
		<th>Observações</th>
		<th>Ações</th>
	</tr>

	{{range .Reunioes}}
	<tr>
		<td>{{.ID}}</td>
		<td>
			<strong>{{.Produtor}}</strong><br>
			{{.Municipio}}/{{.UF}}
		</td>
		<td>{{.Telefone}}</td>
		<td>
			Banco: {{.Banco}}<br>
			Tipo: {{.TipoProjeto}}<br>
			Atividade: {{.Atividade}}<br>
			Classificação: {{.ClassificacaoProdutor}}
		</td>
		<td>
			Cadastro banco: {{.CadastroBanco}}<br>
			Financ. ativo: {{.FinanciamentoAtivo}}<br>
			Restrição: {{.RestricaoCadastral}}<br>
			Imóvel próprio: {{.ImovelProprio}}<br>
			Arrendado/parceria: {{.ImovelArrendado}}<br>
			CAR: {{.TemCAR}}<br>
			Água/irrigação: {{.UsaAgua}}<br>
			Pecuária: {{.TemPecuaria}}<br>
			Investimento: {{.TemInvestimento}}<br>
			Obra: {{.TemObra}}<br>
			Supressão: {{.TemSupressao}}<br>
			ZARC: {{.PrecisaZARC}}
		</td>
		<td>
			<strong>{{.Resumo.PercentualConcluido}}%</strong> concluído<br>
			Total: {{.Resumo.Total}}<br>
			Pendentes: {{.Resumo.Pendentes}}<br>
			Recebidos: {{.Resumo.Recebidos}}<br>
			Não se aplica: {{.Resumo.NaoSeAplica}}
		</td>
		<td>{{.CriadoEm}}</td>
		<td>{{.Observacoes}}</td>
		<td>
			<a class="botao" href="/detalhes?id={{.ID}}">Detalhes</a>
			<a class="botao alerta" href="/editar-reuniao?id={{.ID}}">Editar</a>
			<a class="botao perigo" href="/confirmar-excluir?id={{.ID}}">Excluir</a>
			<a class="botao" href="/checklist?id={{.ID}}">Checklist</a>
			<a class="botao" href="/checklist-controle?id={{.ID}}">Controle</a>
			<a class="botao alerta" href="/whatsapp?id={{.ID}}">WhatsApp</a>
			<a class="botao secundario" href="/exportar-checklist-txt?id={{.ID}}">TXT</a>
		</td>
	</tr>
	{{else}}
	<tr>
		<td colspan="9">Nenhuma reunião encontrada.</td>
	</tr>
	{{end}}
</table>
`

const editarReuniaoHTML = `
<h2>Editar reunião</h2>

<form method="POST" action="/atualizar-reuniao">
	<input type="hidden" name="id" value="{{.Reuniao.ID}}">

	<h3>Dados principais</h3>

	<label>Nome do produtor</label>
	<input type="text" name="produtor" value="{{.Reuniao.Produtor}}" required>

	<label>Telefone / WhatsApp</label>
	<input type="text" name="telefone" value="{{.Reuniao.Telefone}}">

	<label>Município</label>
	<input type="text" name="municipio" value="{{.Reuniao.Municipio}}">

	<label>UF</label>
	<input type="text" name="uf" maxlength="2" value="{{.Reuniao.UF}}">

	<label>Banco pretendido</label>
	<select name="banco">
		<option value="Ainda não definido" {{if eq .Reuniao.Banco "Ainda não definido"}}selected{{end}}>Ainda não definido</option>
		<option value="Banco do Brasil" {{if eq .Reuniao.Banco "Banco do Brasil"}}selected{{end}}>Banco do Brasil</option>
		<option value="Sicoob" {{if eq .Reuniao.Banco "Sicoob"}}selected{{end}}>Sicoob</option>
	</select>

	<label>Tipo de projeto</label>
	<select name="tipo_projeto">
		<option value="Ainda não definido" {{if eq .Reuniao.TipoProjeto "Ainda não definido"}}selected{{end}}>Ainda não definido</option>
		<option value="Custeio agrícola" {{if eq .Reuniao.TipoProjeto "Custeio agrícola"}}selected{{end}}>Custeio agrícola</option>
		<option value="Custeio pecuário" {{if eq .Reuniao.TipoProjeto "Custeio pecuário"}}selected{{end}}>Custeio pecuário</option>
		<option value="Investimento" {{if eq .Reuniao.TipoProjeto "Investimento"}}selected{{end}}>Investimento</option>
		<option value="Comercialização" {{if eq .Reuniao.TipoProjeto "Comercialização"}}selected{{end}}>Comercialização</option>
		<option value="Industrialização" {{if eq .Reuniao.TipoProjeto "Industrialização"}}selected{{end}}>Industrialização</option>
	</select>

	<label>Atividade principal</label>
	<select name="atividade">
		<option value="Ainda não definida" {{if eq .Reuniao.Atividade "Ainda não definida"}}selected{{end}}>Ainda não definida</option>
		<option value="Agrícola" {{if eq .Reuniao.Atividade "Agrícola"}}selected{{end}}>Agrícola</option>
		<option value="Pecuária de corte" {{if eq .Reuniao.Atividade "Pecuária de corte"}}selected{{end}}>Pecuária de corte</option>
		<option value="Pecuária de leite" {{if eq .Reuniao.Atividade "Pecuária de leite"}}selected{{end}}>Pecuária de leite</option>
		<option value="Irrigação" {{if eq .Reuniao.Atividade "Irrigação"}}selected{{end}}>Irrigação</option>
		<option value="Máquinas e equipamentos" {{if eq .Reuniao.Atividade "Máquinas e equipamentos"}}selected{{end}}>Máquinas e equipamentos</option>
		<option value="Obras e benfeitorias" {{if eq .Reuniao.Atividade "Obras e benfeitorias"}}selected{{end}}>Obras e benfeitorias</option>
		<option value="Outro" {{if eq .Reuniao.Atividade "Outro"}}selected{{end}}>Outro</option>
	</select>

	<h3>Classificação do produtor</h3>

	<label>Receita Bruta Agropecuária Anual - RBA</label>
	<input type="text" name="renda_anual" value="{{printf "%.2f" .Reuniao.RendaAnual}}" placeholder="Ex: 500000,00">

	<p class="pequeno">
		Informe a renda bruta agropecuária anual aproximada para classificar o produtor.
	</p>

	<h3>Triagem guiada</h3>

<div class="card">
	<h4>1. Situação cadastral do produtor</h4>

	<label class="check">
		<input type="checkbox" name="cadastro_banco" value="sim" {{if eq .Reuniao.CadastroBanco "sim"}}checked{{end}}>
		O produtor já tem cadastro no banco pretendido?
	</label>

	<label class="check">
		<input type="checkbox" name="financiamento_ativo" value="sim" {{if eq .Reuniao.FinanciamentoAtivo "sim"}}checked{{end}}>
		O produtor possui financiamento rural ativo?
	</label>

	<label class="check">
		<input type="checkbox" name="restricao_cadastral" value="sim" {{if eq .Reuniao.RestricaoCadastral "sim"}}checked{{end}}>
		O produtor informou alguma restrição cadastral conhecida?
	</label>
</div>

<div class="card">
	<h4>2. Situação da propriedade</h4>

	<label class="check">
		<input type="checkbox" name="imovel_proprio" value="sim" {{if eq .Reuniao.ImovelProprio "sim"}}checked{{end}}>
		A área do projeto é imóvel próprio?
	</label>

	<label class="check">
		<input type="checkbox" name="imovel_arrendado" value="sim" {{if eq .Reuniao.ImovelArrendado "sim"}}checked{{end}}>
		A área é arrendada, parceria ou comodato?
	</label>

	<label class="check">
		<input type="checkbox" name="tem_car" value="sim" {{if eq .Reuniao.TemCAR "sim"}}checked{{end}}>
		A propriedade possui CAR?
	</label>
</div>

<div class="card">
	<h4>3. Situação ambiental e uso da área</h4>

	<label class="check">
		<input type="checkbox" name="usa_agua" value="sim" {{if eq .Reuniao.UsaAgua "sim"}}checked{{end}}>
		O projeto usa água, irrigação, poço, represa ou pecuária intensiva?
	</label>

	<label class="check">
		<input type="checkbox" name="tem_supressao" value="sim" {{if eq .Reuniao.TemSupressao "sim"}}checked{{end}}>
		Haverá supressão vegetal, abertura de área ou limpeza de vegetação nativa?
	</label>
</div>

<div class="card">
	<h4>4. Tipo do projeto pretendido</h4>

	<label class="check">
		<input type="checkbox" name="tem_pecuaria" value="sim" {{if eq .Reuniao.TemPecuaria "sim"}}checked{{end}}>
		O projeto envolve pecuária?
	</label>

	<label class="check">
		<input type="checkbox" name="tem_investimento" value="sim" {{if eq .Reuniao.TemInvestimento "sim"}}checked{{end}}>
		O projeto envolve investimento?
	</label>

	<label class="check">
		<input type="checkbox" name="tem_obra" value="sim" {{if eq .Reuniao.TemObra "sim"}}checked{{end}}>
		O projeto possui obra, construção, reforma ou benfeitoria?
	</label>

	<label class="check">
		<input type="checkbox" name="precisa_zarc" value="sim" {{if eq .Reuniao.PrecisaZARC "sim"}}checked{{end}}>
		O projeto precisa validar ZARC?
	</label>
</div>

<label>Observações iniciais</label>
	<textarea name="observacoes">{{.Reuniao.Observacoes}}</textarea>

	<br><br>
	<button type="submit">Salvar alterações</button>
	<a class="botao secundario" href="/detalhes?id={{.Reuniao.ID}}">Cancelar</a>
</form>
`

const confirmarExcluirHTML = `
<h2>Confirmar exclusão</h2>

<div class="card-perigo">
	<p><strong>Atenção:</strong> esta ação vai apagar definitivamente esta reunião do banco de dados.</p>

	<p><strong>Produtor:</strong> {{.Reuniao.Produtor}}</p>
	<p><strong>Telefone:</strong> {{.Reuniao.Telefone}}</p>
	<p><strong>Município/UF:</strong> {{.Reuniao.Municipio}}/{{.Reuniao.UF}}</p>
	<p><strong>Banco:</strong> {{.Reuniao.Banco}}</p>
	<p><strong>Tipo:</strong> {{.Reuniao.TipoProjeto}}</p>
	<p><strong>Atividade:</strong> {{.Reuniao.Atividade}}</p>
</div>

<form method="POST" action="/excluir-reuniao">
	<input type="hidden" name="id" value="{{.Reuniao.ID}}">

	<button class="perigo" type="submit">Sim, excluir reunião</button>
	<a class="botao secundario" href="/detalhes?id={{.Reuniao.ID}}">Cancelar</a>
</form>
`

const detalhesHTML = `
<h2>Detalhes da reunião</h2>

<a class="botao secundario" href="/reunioes">Voltar para reuniões</a>
<a class="botao alerta" href="/editar-reuniao?id={{.Reuniao.ID}}">Editar reunião</a>
<a class="botao perigo" href="/confirmar-excluir?id={{.Reuniao.ID}}">Excluir reunião</a>
<a class="botao" href="/checklist?id={{.Reuniao.ID}}">Gerar checklist</a>
<a class="botao" href="/checklist-controle?id={{.Reuniao.ID}}">Controlar checklist</a>
<a class="botao alerta" href="/whatsapp?id={{.Reuniao.ID}}">Resumo WhatsApp</a>
<a class="botao alerta" href="/whatsapp-pendencias?id={{.Reuniao.ID}}">WhatsApp pendências</a>
<a class="botao secundario" href="/exportar-checklist-txt?id={{.Reuniao.ID}}">Exportar TXT</a>
<a class="botao" href="/exportar-checklist-controle-txt?id={{.Reuniao.ID}}">Exportar controle TXT</a>

<div class="card">
	<h3>Produtor</h3>
	<p><strong>Nome:</strong> {{.Reuniao.Produtor}}</p>
	<p><strong>Telefone/WhatsApp:</strong> {{.Reuniao.Telefone}}</p>
	<p><strong>Município/UF:</strong> {{.Reuniao.Municipio}}/{{.Reuniao.UF}}</p>
	<p><strong>Data da reunião:</strong> {{.Reuniao.CriadoEm}}</p>
</div>

<div class="card">
	<h3>Projeto pretendido</h3>
	<p><strong>Banco:</strong> {{.Reuniao.Banco}}</p>
	<p><strong>Tipo de projeto:</strong> {{.Reuniao.TipoProjeto}}</p>
	<p><strong>Atividade:</strong> {{.Reuniao.Atividade}}</p>
	<p><strong>Renda anual:</strong> R$ {{printf "%.2f" .Reuniao.RendaAnual}}</p>
	<p><strong>Classificação:</strong> {{.Reuniao.ClassificacaoProdutor}}</p>
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

<div class="card">
	<h3>Triagem guiada</h3>

<div class="card">
	<h4>1. Situação cadastral do produtor</h4>

	<label class="check">
		<input type="checkbox" name="cadastro_banco" value="sim">
		O produtor já tem cadastro no banco pretendido?
	</label>

	<label class="check">
		<input type="checkbox" name="financiamento_ativo" value="sim">
		O produtor possui financiamento rural ativo?
	</label>

	<label class="check">
		<input type="checkbox" name="restricao_cadastral" value="sim">
		O produtor informou alguma restrição cadastral conhecida?
	</label>
</div>

<div class="card">
	<h4>2. Classificação pela renda rural</h4>

	<label>Receita Bruta Agropecuária Anual - RBA</label>
	<input type="text" name="renda_anual" placeholder="Ex: 500000,00">

	<p class="pequeno">
		Informe a renda bruta agropecuária anual aproximada do produtor.
		O sistema usará esse valor para classificar como pequeno, médio ou grande produtor.
	</p>
</div>

<div class="card">
	<h4>3. Situação da propriedade</h4>

	<label class="check">
		<input type="checkbox" name="imovel_proprio" value="sim">
		A área do projeto é imóvel próprio?
	</label>

	<label class="check">
		<input type="checkbox" name="imovel_arrendado" value="sim">
		A área é arrendada, parceria ou comodato?
	</label>

	<label class="check">
		<input type="checkbox" name="tem_car" value="sim">
		A propriedade possui CAR?
	</label>
</div>

<div class="card">
	<h4>4. Situação ambiental e uso da área</h4>

	<label class="check">
		<input type="checkbox" name="usa_agua" value="sim">
		O projeto usa água, irrigação, poço, represa ou pecuária intensiva?
	</label>

	<label class="check">
		<input type="checkbox" name="tem_supressao" value="sim">
		Haverá supressão vegetal, abertura de área ou limpeza de vegetação nativa?
	</label>
</div>

<div class="card">
	<h4>5. Tipo do projeto pretendido</h4>

	<label class="check">
		<input type="checkbox" name="tem_pecuaria" value="sim">
		O projeto envolve pecuária?
	</label>

	<label class="check">
		<input type="checkbox" name="tem_investimento" value="sim">
		O projeto envolve investimento?
	</label>

	<label class="check">
		<input type="checkbox" name="tem_obra" value="sim">
		O projeto possui obra, construção, reforma ou benfeitoria?
	</label>

	<label class="check">
		<input type="checkbox" name="precisa_zarc" value="sim">
		O projeto precisa validar ZARC?
	</label>
</div>

<div class="card">
	<h3>Observações da reunião</h3>
	<p>{{.Reuniao.Observacoes}}</p>
</div>
`

const checklistHTML = `
<h2>Checklist automático</h2>

<div class="card">
	<strong>Produtor:</strong> {{.Reuniao.Produtor}}<br>
	<strong>Banco:</strong> {{.Reuniao.Banco}}<br>
	<strong>Tipo:</strong> {{.Reuniao.TipoProjeto}}<br>
	<strong>Atividade:</strong> {{.Reuniao.Atividade}}<br>
</div>

<a class="botao secundario" href="/reunioes">Voltar para reuniões</a>
<a class="botao" href="/exportar-checklist-txt?id={{.Reuniao.ID}}">Exportar checklist TXT</a>
<a class="botao" href="/checklist-controle?id={{.Reuniao.ID}}">Controlar checklist</a>
<a class="botao alerta" href="/whatsapp?id={{.Reuniao.ID}}">Gerar resumo WhatsApp</a>

<pre>{{.Checklist}}</pre>
`

const whatsappHTML = `
<h2>Resumo para WhatsApp</h2>

<div class="card">
	<strong>Produtor:</strong> {{.Reuniao.Produtor}}<br>
	<strong>Banco:</strong> {{.Reuniao.Banco}}<br>
	<strong>Tipo:</strong> {{.Reuniao.TipoProjeto}}<br>
	<strong>Atividade:</strong> {{.Reuniao.Atividade}}<br>
</div>

<a class="botao secundario" href="/reunioes">Voltar para reuniões</a>
<a class="botao" href="/exportar-whatsapp-txt?id={{.Reuniao.ID}}">Exportar WhatsApp TXT</a>

<p class="pequeno">Copie a mensagem abaixo e envie ao produtor.</p>

<pre>{{.Mensagem}}</pre>
`

const arquivoGeradoHTML = `
<h2>Arquivo gerado</h2>

<div class="card">
	<p>Arquivo salvo em:</p>
	<pre>{{.Caminho}}</pre>
</div>

<a class="botao" href="/reunioes">Voltar para reuniões</a>
<a class="botao secundario" href="/">Início</a>
`
