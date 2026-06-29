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
	
/* UI MODERNA CHECKLIST RURAL */
:root {
	--verde: #1f7a4d;
	--verde-escuro: #145c38;
	--verde-claro: #eef8f1;
	--fundo: #f4f7f3;
	--card: #ffffff;
	--texto: #1f2933;
	--muted: #6b7280;
	--borda: #d9e2dc;
	--alerta: #fff7e6;
	--perigo: #fff1f1;
}

body {
	background: var(--fundo) !important;
	color: var(--texto);
	font-family: Arial, Helvetica, sans-serif;
	margin: 0;
}

.container, main {
	max-width: 1180px;
	margin: 0 auto;
	padding: 24px;
}

h1, h2, h3 {
	color: var(--verde-escuro);
}

h2 {
	font-size: 24px;
	margin-top: 0;
	margin-bottom: 18px;
}

h3 {
	font-size: 18px;
	margin-top: 0;
}

.card {
	background: var(--card);
	border: 1px solid var(--borda);
	border-radius: 16px;
	padding: 20px;
	margin-bottom: 18px;
	box-shadow: 0 8px 22px rgba(15, 23, 42, 0.05);
}

.card.destaque {
	background: var(--verde-claro);
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

.grid {
	display: grid;
	grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
	gap: 16px;
}

label {
	display: block;
	font-weight: bold;
	margin-top: 12px;
	margin-bottom: 6px;
	color: var(--texto);
}

input, select, textarea {
	width: 100%;
	padding: 11px 12px;
	border: 1px solid var(--borda);
	border-radius: 10px;
	font-size: 14px;
	background: #fff;
	color: var(--texto);
}

textarea {
	resize: vertical;
}

button, .botao, a.botao {
	display: inline-block;
	border: 0;
	background: var(--verde);
	color: white !important;
	text-decoration: none;
	border-radius: 10px;
	padding: 11px 15px;
	font-weight: bold;
	font-size: 14px;
	cursor: pointer;
	margin: 4px 4px 4px 0;
}

button:hover, .botao:hover, a.botao:hover {
	background: var(--verde-escuro);
}

.botao.secundario, a.botao.secundario {
	background: #eef6f0;
	color: var(--verde-escuro) !important;
	border: 1px solid #cfe5d6;
}

.botao.secundario:hover, a.botao.secundario:hover {
	background: #dcefe4;
}

table {
	width: 100%;
	border-collapse: collapse;
	background: white;
	border-radius: 12px;
	overflow: hidden;
	margin-top: 12px;
}

th {
	background: #eef6f0;
	color: var(--verde-escuro);
	text-align: left;
	padding: 12px;
	font-size: 14px;
}

td {
	border-top: 1px solid var(--borda);
	padding: 12px;
	vertical-align: top;
	font-size: 14px;
}

tr:hover td {
	background: #fafdfb;
}

.pequeno {
	font-size: 13px;
	color: var(--muted);
	line-height: 1.45;
}

pre {
	background: #0f172a;
	color: #e5e7eb;
	border-radius: 12px;
	padding: 14px;
	overflow-x: auto;
}

.topo-sistema {
	background: linear-gradient(135deg, var(--verde), var(--verde-escuro));
	color: white;
	padding: 18px 24px;
	margin-bottom: 24px;
}

.topo-sistema .topo-conteudo {
	max-width: 1180px;
	margin: 0 auto;
	display: flex;
	justify-content: space-between;
	align-items: center;
	gap: 16px;
	flex-wrap: wrap;
}

.topo-sistema strong {
	font-size: 20px;
}

.topo-sistema nav {
	display: flex;
	gap: 8px;
	flex-wrap: wrap;
}

.topo-sistema a {
	color: white !important;
	text-decoration: none;
	background: rgba(255,255,255,0.14);
	padding: 8px 11px;
	border-radius: 9px;
	font-size: 13px;
	font-weight: bold;
}

.topo-sistema a:hover {
	background: rgba(255,255,255,0.24);
}

.barra-acoes {
	display: flex;
	flex-wrap: wrap;
	gap: 8px;
	margin: 12px 0 20px;
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

@media (max-width: 700px) {
	.container, main {
		padding: 14px;
	}

	.card {
		padding: 16px;
	}

	table {
		display: block;
		overflow-x: auto;
	}

	.topo-sistema .topo-conteudo {
		align-items: flex-start;
	}
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

	<h3>Pré-análise do produtor e do projeto</h3>

<p class="pequeno">
	Preencha esta parte conversando com o produtor. As respostas ajudam a classificar o caso, montar o checklist e futuramente sugerir possíveis linhas de crédito.
</p>

<div class="card">
	<h4>1. Quem é o produtor?</h4>

	<label>Qual é a Receita Bruta Agropecuária Anual aproximada?</label>
	<input type="text" name="renda_anual" placeholder="Ex: 500000,00">

	<label class="check">
		<input type="checkbox" name="possui_caf" value="sim">
		O produtor possui CAF ativo para agricultura familiar?
	</label>

	<p class="pequeno">
		Essa etapa ajuda a separar se o produtor pode caminhar para agricultura familiar, médio produtor ou outro enquadramento.
	</p>
</div>

<div class="card">
	<h4>2. O que o produtor quer financiar?</h4>

	<label>Qual é a finalidade principal do crédito?</label>
	<select name="finalidade_credito">
		<option value="">Selecione</option>
		<option value="Custeio agrícola">Custeio agrícola / safra</option>
		<option value="Custeio pecuário">Custeio pecuário / animais</option>
		<option value="Investimento">Investimento na propriedade</option>
		<option value="Máquinas e equipamentos">Máquinas e equipamentos</option>
		<option value="Obras e benfeitorias">Obras, construção ou benfeitorias</option>
		<option value="Irrigação">Irrigação ou uso de água</option>
	</select>

	<label>Qual valor aproximado o produtor pretende financiar?</label>
	<input type="text" name="valor_pretendido" placeholder="Ex: 250000,00">

	<label class="check">
		<input type="checkbox" name="tem_orcamento" value="sim">
		O produtor já tem orçamento, proposta comercial ou valor definido?
	</label>

	<p class="pequeno">
		Aqui você identifica se a demanda é custeio, investimento, máquina, obra, irrigação ou pecuária.
	</p>
</div>

<div class="card">
	<h4>3. Relação com o banco</h4>

	<label class="check">
		<input type="checkbox" name="cadastro_banco" value="sim">
		O produtor já tem cadastro nesse banco ou cooperativa?
	</label>

	<label class="check">
		<input type="checkbox" name="financiamento_ativo" value="sim">
		O produtor possui financiamento rural ativo hoje?
	</label>

	<label class="check">
		<input type="checkbox" name="restricao_cadastral" value="sim">
		Existe alguma restrição ou pendência no CPF/CNPJ que ele saiba?
	</label>

	<p class="pequeno">
		Essa parte ajuda a saber se o projeto pode andar ou se primeiro será necessário resolver cadastro, limite ou pendência.
	</p>
</div>

<div class="card">
	<h4>4. Situação da terra</h4>

	<label class="check">
		<input type="checkbox" name="imovel_proprio" value="sim">
		A área onde será feito o projeto é própria?
	</label>

	<label class="check">
		<input type="checkbox" name="imovel_arrendado" value="sim">
		A área é arrendada, parceria ou comodato?
	</label>

	<label class="check">
		<input type="checkbox" name="tem_car" value="sim">
		A propriedade já possui CAR?
	</label>

	<p class="pequeno">
		Se a área não for própria, o checklist vai precisar considerar contrato, anuência e autorização do proprietário.
	</p>
</div>

<div class="card">
	<h4>5. Ambiental e uso da água</h4>

	<label class="check">
		<input type="checkbox" name="usa_agua" value="sim">
		O projeto vai usar água, poço, represa, rio, córrego ou irrigação?
	</label>

	<label class="check">
		<input type="checkbox" name="tem_supressao" value="sim">
		Vai precisar abrir área, retirar vegetação ou limpar vegetação nativa?
	</label>

	<p class="pequeno">
		Essa etapa aponta riscos ambientais que podem exigir outorga, dispensa, licença ou autorização.
	</p>
</div>

<div class="card">
	<h4>6. Pontos técnicos do projeto</h4>

	<label class="check">
		<input type="checkbox" name="tem_pecuaria" value="sim">
		O projeto envolve pecuária, rebanho, pastagem, leite ou corte?
	</label>

	<label class="check">
		<input type="checkbox" name="tem_investimento" value="sim">
		O projeto envolve compra, melhoria, estrutura ou investimento na propriedade?
	</label>

	<label class="check">
		<input type="checkbox" name="tem_obra" value="sim">
		O projeto envolve obra, construção, reforma, curral, barracão ou benfeitoria?
	</label>

	<label class="check">
		<input type="checkbox" name="precisa_zarc" value="sim">
		É lavoura/custeio agrícola e precisa conferir ZARC?
	</label>

	<p class="pequeno">
		Essas respostas definem quais blocos entram no checklist técnico.
	</p>
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

<div class="barra-acoes">
	<a class="botao" href="/nova-reuniao">Nova reunião</a>
	<a class="botao secundario" href="/">Início</a>
</div>

<div class="card destaque">
	<h3>Buscar projeto</h3>
	<p class="pequeno">
		Pesquise por produtor, município, banco, atividade, tipo de projeto ou classificação.
	</p>

	<input 
		type="text" 
		id="buscaReunioes" 
		placeholder="Digite para filtrar as reuniões..."
		onkeyup="filtrarReunioes()"
	>
</div>

{{if .Reunioes}}
<div class="grid" id="listaReunioes">
	{{range .Reunioes}}
	<div class="card card-reuniao" data-busca="{{.Produtor}} {{.Municipio}} {{.UF}} {{.Banco}} {{.TipoProjeto}} {{.Atividade}} {{.ClassificacaoProdutor}} {{.FinalidadeCredito}}">
		<h3>{{.Produtor}}</h3>

		<p class="pequeno">
			{{.Municipio}}/{{.UF}} — {{.Banco}}
		</p>

		<p>
			<span class="badge">{{.TipoProjeto}}</span>
			<span class="badge">{{.Atividade}}</span>
			{{if .ClassificacaoProdutor}}
				<span class="badge">{{.ClassificacaoProdutor}}</span>
			{{end}}
		</p>

		<p><strong>Finalidade:</strong> {{if .FinalidadeCredito}}{{.FinalidadeCredito}}{{else}}-{{end}}</p>
		<p><strong>Data:</strong> {{.CriadoEm}}</p>

		<div class="barra-acoes">
			<a class="botao" href="/detalhes?id={{.ID}}">Abrir</a>
			<a class="botao secundario" href="/editar-reuniao?id={{.ID}}">Editar</a>
			<a class="botao secundario" href="/investigacao?id={{.ID}}">Investigação</a>
			<a class="botao secundario" href="/relatorio?id={{.ID}}" target="_blank">Relatório</a>
		</div>
	</div>
	{{end}}
</div>
{{else}}
<div class="card alerta">
	<h3>Nenhuma reunião cadastrada</h3>
	<p>
		Comece criando sua primeira reunião de projeto rural.
	</p>

	<a class="botao" href="/nova-reuniao">Criar primeira reunião</a>
</div>
{{end}}

<script>
function filtrarReunioes() {
	const campo = document.getElementById("buscaReunioes");
	const termo = campo.value.toLowerCase().trim();
	const cards = document.querySelectorAll(".card-reuniao");

	cards.forEach(function(card) {
		const texto = card.getAttribute("data-busca").toLowerCase();

		if (texto.includes(termo)) {
			card.style.display = "";
		} else {
			card.style.display = "none";
		}
	});
}
</script>

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

	<h3>Pré-análise do produtor e do projeto</h3>

<p class="pequeno">
	Preencha esta parte conversando com o produtor. As respostas ajudam a classificar o caso, montar o checklist e futuramente sugerir possíveis linhas de crédito.
</p>

<div class="card">
	<h4>1. Quem é o produtor?</h4>

	<label>Qual é a Receita Bruta Agropecuária Anual aproximada?</label>
	<input type="text" name="renda_anual" value="{{printf "%.2f" .Reuniao.RendaAnual}}" placeholder="Ex: 500000,00">

	<label class="check">
		<input type="checkbox" name="possui_caf" value="sim" {{if eq .Reuniao.PossuiCAF "sim"}}checked{{end}}>
		O produtor possui CAF ativo para agricultura familiar?
	</label>

	<p class="pequeno">
		Essa etapa ajuda a separar se o produtor pode caminhar para agricultura familiar, médio produtor ou outro enquadramento.
	</p>
</div>

<div class="card">
	<h4>2. O que o produtor quer financiar?</h4>

	<label>Qual é a finalidade principal do crédito?</label>
	<select name="finalidade_credito">
		<option value="">Selecione</option>
		<option value="Custeio agrícola" {{if eq .Reuniao.FinalidadeCredito "Custeio agrícola"}}selected{{end}}>Custeio agrícola / safra</option>
		<option value="Custeio pecuário" {{if eq .Reuniao.FinalidadeCredito "Custeio pecuário"}}selected{{end}}>Custeio pecuário / animais</option>
		<option value="Investimento" {{if eq .Reuniao.FinalidadeCredito "Investimento"}}selected{{end}}>Investimento na propriedade</option>
		<option value="Máquinas e equipamentos" {{if eq .Reuniao.FinalidadeCredito "Máquinas e equipamentos"}}selected{{end}}>Máquinas e equipamentos</option>
		<option value="Obras e benfeitorias" {{if eq .Reuniao.FinalidadeCredito "Obras e benfeitorias"}}selected{{end}}>Obras, construção ou benfeitorias</option>
		<option value="Irrigação" {{if eq .Reuniao.FinalidadeCredito "Irrigação"}}selected{{end}}>Irrigação ou uso de água</option>
	</select>

	<label>Qual valor aproximado o produtor pretende financiar?</label>
	<input type="text" name="valor_pretendido" value="{{printf "%.2f" .Reuniao.ValorPretendido}}" placeholder="Ex: 250000,00">

	<label class="check">
		<input type="checkbox" name="tem_orcamento" value="sim" {{if eq .Reuniao.TemOrcamento "sim"}}checked{{end}}>
		O produtor já tem orçamento, proposta comercial ou valor definido?
	</label>

	<p class="pequeno">
		Aqui você identifica se a demanda é custeio, investimento, máquina, obra, irrigação ou pecuária.
	</p>
</div>

<div class="card">
	<h4>3. Relação com o banco</h4>

	<label class="check">
		<input type="checkbox" name="cadastro_banco" value="sim" {{if eq .Reuniao.CadastroBanco "sim"}}checked{{end}}>
		O produtor já tem cadastro nesse banco ou cooperativa?
	</label>

	<label class="check">
		<input type="checkbox" name="financiamento_ativo" value="sim" {{if eq .Reuniao.FinanciamentoAtivo "sim"}}checked{{end}}>
		O produtor possui financiamento rural ativo hoje?
	</label>

	<label class="check">
		<input type="checkbox" name="restricao_cadastral" value="sim" {{if eq .Reuniao.RestricaoCadastral "sim"}}checked{{end}}>
		Existe alguma restrição ou pendência no CPF/CNPJ que ele saiba?
	</label>

	<p class="pequeno">
		Essa parte ajuda a saber se o projeto pode andar ou se primeiro será necessário resolver cadastro, limite ou pendência.
	</p>
</div>

<div class="card">
	<h4>4. Situação da terra</h4>

	<label class="check">
		<input type="checkbox" name="imovel_proprio" value="sim" {{if eq .Reuniao.ImovelProprio "sim"}}checked{{end}}>
		A área onde será feito o projeto é própria?
	</label>

	<label class="check">
		<input type="checkbox" name="imovel_arrendado" value="sim" {{if eq .Reuniao.ImovelArrendado "sim"}}checked{{end}}>
		A área é arrendada, parceria ou comodato?
	</label>

	<label class="check">
		<input type="checkbox" name="tem_car" value="sim" {{if eq .Reuniao.TemCAR "sim"}}checked{{end}}>
		A propriedade já possui CAR?
	</label>

	<p class="pequeno">
		Se a área não for própria, o checklist vai precisar considerar contrato, anuência e autorização do proprietário.
	</p>
</div>

<div class="card">
	<h4>5. Ambiental e uso da água</h4>

	<label class="check">
		<input type="checkbox" name="usa_agua" value="sim" {{if eq .Reuniao.UsaAgua "sim"}}checked{{end}}>
		O projeto vai usar água, poço, represa, rio, córrego ou irrigação?
	</label>

	<label class="check">
		<input type="checkbox" name="tem_supressao" value="sim" {{if eq .Reuniao.TemSupressao "sim"}}checked{{end}}>
		Vai precisar abrir área, retirar vegetação ou limpar vegetação nativa?
	</label>

	<p class="pequeno">
		Essa etapa aponta riscos ambientais que podem exigir outorga, dispensa, licença ou autorização.
	</p>
</div>

<div class="card">
	<h4>6. Pontos técnicos do projeto</h4>

	<label class="check">
		<input type="checkbox" name="tem_pecuaria" value="sim" {{if eq .Reuniao.TemPecuaria "sim"}}checked{{end}}>
		O projeto envolve pecuária, rebanho, pastagem, leite ou corte?
	</label>

	<label class="check">
		<input type="checkbox" name="tem_investimento" value="sim" {{if eq .Reuniao.TemInvestimento "sim"}}checked{{end}}>
		O projeto envolve compra, melhoria, estrutura ou investimento na propriedade?
	</label>

	<label class="check">
		<input type="checkbox" name="tem_obra" value="sim" {{if eq .Reuniao.TemObra "sim"}}checked{{end}}>
		O projeto envolve obra, construção, reforma, curral, barracão ou benfeitoria?
	</label>

	<label class="check">
		<input type="checkbox" name="precisa_zarc" value="sim" {{if eq .Reuniao.PrecisaZARC "sim"}}checked{{end}}>
		É lavoura/custeio agrícola e precisa conferir ZARC?
	</label>

	<p class="pequeno">
		Essas respostas definem quais blocos entram no checklist técnico.
	</p>
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

const detalhesReuniaoModernoHTML = `
<h2>Detalhes da reunião</h2>

<div class="barra-acoes">
	<a class="botao secundario" href="/reunioes">Voltar</a>
	<a class="botao secundario" href="/editar-reuniao?id={{.Reuniao.ID}}">Editar</a>
	<a class="botao perigo" href="/confirmar-excluir?id={{.Reuniao.ID}}">Excluir</a>
</div>

<div class="grid">
	<div class="card destaque">
		<h3>Produtor</h3>
		<p class="pequeno">Dados principais da reunião.</p>

		<p><strong>Nome:</strong> {{.Reuniao.Produtor}}</p>
		<p><strong>Telefone:</strong> {{.Reuniao.Telefone}}</p>
		<p><strong>Município/UF:</strong> {{.Reuniao.Municipio}}/{{.Reuniao.UF}}</p>
		<p><strong>Data:</strong> {{.Reuniao.CriadoEm}}</p>
	</div>

	<div class="card">
		<h3>Projeto</h3>
		<p class="pequeno">Resumo da intenção de financiamento.</p>

		<p><strong>Banco:</strong> {{.Reuniao.Banco}}</p>
		<p><strong>Tipo:</strong> {{.Reuniao.TipoProjeto}}</p>
		<p><strong>Atividade:</strong> {{.Reuniao.Atividade}}</p>
		<p><strong>Finalidade:</strong> {{.Reuniao.FinalidadeCredito}}</p>
	</div>

	<div class="card">
		<h3>Enquadramento</h3>
		<p class="pequeno">Classificação inicial do produtor.</p>

		<p><strong>Renda anual:</strong> R$ {{printf "%.2f" .Reuniao.RendaAnual}}</p>
		<p><strong>Classificação:</strong> <span class="badge">{{.Reuniao.ClassificacaoProdutor}}</span></p>
		<p><strong>Possui CAF:</strong> {{.Reuniao.PossuiCAF}}</p>
		<p><strong>Valor pretendido:</strong> R$ {{printf "%.2f" .Reuniao.ValorPretendido}}</p>
	</div>
</div>

<div class="card">
	<h3>Fluxo da análise</h3>
	<p class="pequeno">
		Siga a ordem abaixo para manter o projeto organizado.
	</p>

	<div class="grid">
		<div class="card">
			<h3>1. Investigação</h3>
			<p class="pequeno">CNPJ, CEP, CAR, imóvel, Ibama, CEIS/CNEP, SIGEF/SNCR e fontes.</p>
			<a class="botao" href="/investigacao?id={{.Reuniao.ID}}">Abrir investigação</a>
		</div>

		<div class="card">
			<h3>2. Georreferenciamento</h3>
			<p class="pequeno">KML, KMZ, GeoJSON, shapefile, PDF, croqui ou planta.</p>

			{{if .ResumoGeoref.Total}}
				<p><strong>Arquivos:</strong> {{.ResumoGeoref.Total}}</p>
			{{else}}
				<p class="pequeno">Nenhum arquivo importado ainda.</p>
			{{end}}

			<a class="botao" href="/georreferenciamento?id={{.Reuniao.ID}}">Abrir georreferenciamento</a>
		</div>

		<div class="card">
			<h3>3. Relatório</h3>
			<p class="pequeno">Gere a pré-análise para imprimir ou salvar em PDF.</p>
			<a class="botao" href="/relatorio?id={{.Reuniao.ID}}" target="_blank">Abrir relatório</a>
		</div>
	</div>
</div>

<div class="card">
	<h3>Ações rápidas</h3>

	<div class="barra-acoes">
		<a class="botao secundario" href="/whatsapp?id={{.Reuniao.ID}}">Resumo WhatsApp</a>
		<a class="botao secundario" href="/whatsapp-pendencias?id={{.Reuniao.ID}}">WhatsApp pendências</a>
		<a class="botao secundario" href="/exportar-checklist-txt?id={{.Reuniao.ID}}">Exportar TXT</a>
		<a class="botao secundario" href="/exportar-checklist-controle-txt?id={{.Reuniao.ID}}">Exportar controle TXT</a>
	</div>
</div>

{{if .ResumoGeoref.Total}}
<div class="card">
	<h3>Arquivos georreferenciados</h3>

	<table>
		<thead>
			<tr>
				<th>Tipo</th>
				<th>Arquivo</th>
				<th>Data</th>
			</tr>
		</thead>
		<tbody>
			{{range .ResumoGeoref.Arquivos}}
			<tr>
				<td>{{.Tipo}}</td>
				<td>{{.NomeOriginal}}</td>
				<td>{{.CriadoEm}}</td>
			</tr>
			{{end}}
		</tbody>
	</table>
</div>
{{end}}

<div class="card">
	<h3>Leitura inicial</h3>

	<p><strong>Enquadramento:</strong> {{.Leitura.Enquadramento}}</p>
	<p><strong>Caminho sugerido:</strong> {{.Leitura.Caminho}}</p>

	{{if .Leitura.Alertas}}
		<h3>Alertas</h3>
		<ul>
			{{range .Leitura.Alertas}}
				<li>{{.}}</li>
			{{end}}
		</ul>
	{{end}}

	{{if .Leitura.ProximosPassos}}
		<h3>Próximos passos</h3>
		<ul>
			{{range .Leitura.ProximosPassos}}
				<li>{{.}}</li>
			{{end}}
		</ul>
	{{end}}
</div>

<div class="card">
	<h3>Observações</h3>

	{{if .Reuniao.Observacoes}}
		<p>{{.Reuniao.Observacoes}}</p>
	{{else}}
		<p class="pequeno">Nenhuma observação registrada.</p>
	{{end}}
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

<a class="botao secundario" href="/configuracoes-api">Configurações das APIs</a>

<a class="botao" href="/reunioes">Voltar para reuniões</a>
<a class="botao secundario" href="/">Início</a>
`
