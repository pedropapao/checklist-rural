package main

const checklistControleHTML = `
<h2>Controle do checklist</h2>

<style>

	.legenda-status {
		display: flex;
		flex-wrap: wrap;
		gap: 10px;
		margin: 12px 0;
	}

	.legenda-item {
		display: inline-flex;
		align-items: center;
		gap: 6px;
		padding: 6px 10px;
		border-radius: 999px;
		border: 1px solid #ddd;
		background: #fff;
		font-size: 0.9em;
	}

	.bolinha-status {
		width: 12px;
		height: 12px;
		border-radius: 999px;
		display: inline-block;
	}

	.bolinha-pendente {
		background: #f59e0b;
	}

	.bolinha-recebido {
		background: #16a34a;
	}

	.bolinha-nao-aplica {
		background: #6b7280;
	}

	.item-checklist[data-status="Pendente"] {
		border-left: 7px solid #f59e0b;
		background: #fffaf0;
	}

	.item-checklist[data-status="Recebido"] {
		border-left: 7px solid #16a34a;
		background: #f0fdf4;
	}

	.item-checklist[data-status="Não se aplica"] {
		border-left: 7px solid #6b7280;
		background: #f9fafb;
		opacity: 0.85;
	}

	.item-checklist[data-status="Recebido"] .item-titulo::before {
		content: "✓ ";
		color: #16a34a;
		font-weight: bold;
	}

	.item-checklist[data-status="Pendente"] .item-titulo::before {
		content: "! ";
		color: #f59e0b;
		font-weight: bold;
	}

	.item-checklist[data-status="Não se aplica"] .item-titulo::before {
		content: "– ";
		color: #6b7280;
		font-weight: bold;
	}

	.item-checklist[data-status="Pendente"] .status-atual {
		background: #fef3c7;
		border-color: #f59e0b;
		color: #92400e;
	}

	.item-checklist[data-status="Recebido"] .status-atual {
		background: #dcfce7;
		border-color: #16a34a;
		color: #166534;
	}

	.item-checklist[data-status="Não se aplica"] .status-atual {
		background: #e5e7eb;
		border-color: #6b7280;
		color: #374151;
	}


	.painel-checklist {
		display: grid;
		grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
		gap: 12px;
		margin-bottom: 14px;
	}

	.caixa-resumo {
		padding: 12px;
		border: 1px solid #ddd;
		border-radius: 10px;
		background: #fafafa;
	}

	.barra-ferramentas {
		position: sticky;
		top: 0;
		z-index: 10;
		padding: 12px;
		margin: 14px 0;
		border: 1px solid #ddd;
		border-radius: 12px;
		background: #ffffff;
		box-shadow: 0 2px 8px rgba(0,0,0,0.05);
	}

	.barra-ferramentas input,
	.barra-ferramentas select {
		width: 100%;
		margin-bottom: 8px;
	}

	.botoes-pequenos {
		display: flex;
		flex-wrap: wrap;
		gap: 8px;
	}

	.botao-pequeno {
		padding: 7px 10px;
		border: 1px solid #bbb;
		border-radius: 8px;
		background: #f7f7f7;
		cursor: pointer;
	}

	.etapa-checklist {
		margin: 14px 0;
		border: 1px solid #ddd;
		border-radius: 14px;
		background: #fff;
		overflow: hidden;
	}

	.etapa-checklist summary {
		cursor: pointer;
		padding: 14px;
		background: #f5f5f5;
		font-weight: bold;
	}

	.etapa-corpo {
		padding: 12px;
	}

	.linha-resumo-etapa {
		display: flex;
		flex-wrap: wrap;
		gap: 8px;
		margin-top: 8px;
		font-size: 0.9em;
		font-weight: normal;
	}

	.selo {
		display: inline-block;
		padding: 3px 8px;
		border-radius: 999px;
		border: 1px solid #ccc;
		background: #fff;
		font-size: 0.85em;
	}

	.item-checklist {
		border: 1px solid #e0e0e0;
		border-radius: 12px;
		padding: 12px;
		margin-bottom: 10px;
		background: #ffffff;
	}

	.item-topo {
		display: flex;
		justify-content: space-between;
		gap: 10px;
		align-items: flex-start;
	}

	.item-titulo {
		font-weight: bold;
		line-height: 1.35;
	}

	.botao-ajuda {
		margin-left: 8px;
		padding: 3px 8px;
		border: 1px solid #999;
		border-radius: 999px;
		background: #f7f7f7;
		cursor: pointer;
		font-size: 0.85em;
	}

	.modal-explicacao {
		display: none;
		position: fixed;
		z-index: 9999;
		left: 0;
		top: 0;
		width: 100%;
		height: 100%;
		background: rgba(0,0,0,0.45);
	}

	.modal-conteudo {
		background: #fff;
		margin: 8% auto;
		padding: 20px;
		border-radius: 14px;
		max-width: 600px;
		box-shadow: 0 8px 30px rgba(0,0,0,0.25);
	}

	.modal-fechar {
		float: right;
		font-size: 24px;
		font-weight: bold;
		cursor: pointer;
	}

	.status-atual {
		white-space: nowrap;
		font-size: 0.85em;
		padding: 4px 8px;
		border-radius: 999px;
		border: 1px solid #ccc;
		background: #f8f8f8;
	}

	.item-campos {
		display: grid;
		grid-template-columns: minmax(160px, 220px) 1fr;
		gap: 10px;
		margin-top: 10px;
	}

	.item-campos textarea {
		min-height: 70px;
	}

	.item-meta {
		margin-top: 6px;
		font-size: 0.85em;
		color: #666;
	}

	.aviso-vazio {
		padding: 14px;
		border: 1px dashed #ccc;
		border-radius: 12px;
		background: #fafafa;
	}

	@media (max-width: 700px) {
		.item-topo {
			display: block;
		}

		.botao-ajuda {
		margin-left: 8px;
		padding: 3px 8px;
		border: 1px solid #999;
		border-radius: 999px;
		background: #f7f7f7;
		cursor: pointer;
		font-size: 0.85em;
	}

	.modal-explicacao {
		display: none;
		position: fixed;
		z-index: 9999;
		left: 0;
		top: 0;
		width: 100%;
		height: 100%;
		background: rgba(0,0,0,0.45);
	}

	.modal-conteudo {
		background: #fff;
		margin: 8% auto;
		padding: 20px;
		border-radius: 14px;
		max-width: 600px;
		box-shadow: 0 8px 30px rgba(0,0,0,0.25);
	}

	.modal-fechar {
		float: right;
		font-size: 24px;
		font-weight: bold;
		cursor: pointer;
	}

	.status-atual {
			display: inline-block;
			margin-top: 8px;
		}

		.item-campos {
			grid-template-columns: 1fr;
		}
	}
</style>

<div class="card">
	<h3>Resumo da reunião</h3>

	<div class="painel-checklist">
		<div class="caixa-resumo">
			<strong>Produtor</strong><br>
			{{.Reuniao.Produtor}}
		</div>

		<div class="caixa-resumo">
			<strong>Banco</strong><br>
			{{.Reuniao.Banco}}
		</div>

		<div class="caixa-resumo">
			<strong>Tipo</strong><br>
			{{.Reuniao.TipoProjeto}}
		</div>

		<div class="caixa-resumo">
			<strong>Classificação</strong><br>
			{{.Reuniao.ClassificacaoProdutor}}
		</div>

		<div class="caixa-resumo">
			<strong>RBA</strong><br>
			R$ {{printf "%.2f" .Reuniao.RendaAnual}}
		</div>
	</div>
</div>

<div class="card">
	<h3>Resumo do andamento</h3>

	<div class="painel-checklist">
		<div class="caixa-resumo">
			<strong>Total</strong><br>
			{{.Resumo.Total}}
		</div>

		<div class="caixa-resumo">
			<strong>Concluído</strong><br>
			{{.Resumo.PercentualConcluido}}%
		</div>

		<div class="caixa-resumo">
			<strong>Pendentes</strong><br>
			{{.Resumo.Pendentes}}
		</div>

		<div class="caixa-resumo">
			<strong>Recebidos</strong><br>
			{{.Resumo.Recebidos}}
		</div>

		<div class="caixa-resumo">
			<strong>Não se aplica</strong><br>
			{{.Resumo.NaoSeAplica}}
		</div>
	</div>
</div>


<div class="card">
	<h3>Legenda visual</h3>

	<div class="legenda-status">
		<span class="legenda-item">
			<span class="bolinha-status bolinha-pendente"></span>
			Pendente: ainda falta resolver
		</span>

		<span class="legenda-item">
			<span class="bolinha-status bolinha-recebido"></span>
			Recebido: documento já chegou
		</span>

		<span class="legenda-item">
			<span class="bolinha-status bolinha-nao-aplica"></span>
			Não se aplica: não precisa neste caso
		</span>
	</div>
</div>


<a class="botao secundario" href="/detalhes?id={{.Reuniao.ID}}">Voltar para detalhes</a>
<a class="botao" href="/checklist?id={{.Reuniao.ID}}">Ver checklist em texto</a>
<a class="botao alerta" href="/whatsapp?id={{.Reuniao.ID}}">Resumo WhatsApp</a>
<a class="botao alerta" href="/whatsapp-pendencias?id={{.Reuniao.ID}}">WhatsApp pendências</a>
<a class="botao" href="/exportar-checklist-controle-txt?id={{.Reuniao.ID}}">Exportar controle TXT</a>

<div class="barra-ferramentas">
	<label>Buscar documento ou pendência</label>
	<input type="text" id="buscaChecklist" placeholder="Digite: CAR, matrícula, outorga, orçamento, ZARC, nota fiscal...">

	<label>Filtrar por status</label>
	<select id="filtroStatus">
		<option value="">Todos os status</option>
		<option value="Pendente">Somente pendentes</option>
		<option value="Recebido">Somente recebidos</option>
		<option value="Não se aplica">Somente não se aplica</option>
	</select>

	<div class="botoes-pequenos">
		<button type="button" class="botao-pequeno" onclick="abrirTodasEtapas()">Abrir todas as etapas</button>
		<button type="button" class="botao-pequeno" onclick="fecharTodasEtapas()">Fechar todas as etapas</button>
		<button type="button" class="botao-pequeno" onclick="mostrarPendentes()">Ver pendentes</button>
		<button type="button" class="botao-pequeno" onclick="limparFiltrosChecklist()">Limpar filtros</button>
	</div>
</div>

<form method="POST" action="/salvar-itens-checklist">
	<input type="hidden" name="reuniao_id" value="{{.Reuniao.ID}}">

	{{range .Grupos}}
	<details class="etapa-checklist" open>
		<summary>
			{{.Etapa}}

			<div class="linha-resumo-etapa">
				<span class="selo">Total: {{.Total}}</span>
				<span class="selo">Pendentes: {{.Pendentes}}</span>
				<span class="selo">Recebidos: {{.Recebidos}}</span>
				<span class="selo">Não se aplica: {{.NaoSeAplica}}</span>
			</div>
		</summary>

		<div class="etapa-corpo">
			<p class="pequeno">{{.Descricao}}</p>

			{{range .Itens}}
			<div class="item-checklist" data-status="{{.Status}}" data-texto="{{.Grupo}} {{.Item}} {{.Observacao}}">
				<div class="item-topo">
					<div>
						<div class="item-titulo">
							{{.Item}}
							<button type="button" class="botao-ajuda" onclick="abrirExplicacao('{{js .Item}}', '{{js .Explicacao}}')">?</button>
						</div>
						<div class="item-meta">Grupo original: {{.Grupo}}</div>
					</div>

					<div class="status-atual">{{.Status}}</div>
				</div>

				<div class="item-campos">
					<div>
						<label>Status</label>
						<select name="status_{{.ID}}" onchange="atualizarStatusVisual(this)">
							<option value="Pendente" {{if eq .Status "Pendente"}}selected{{end}}>Pendente</option>
							<option value="Recebido" {{if eq .Status "Recebido"}}selected{{end}}>Recebido</option>
							<option value="Não se aplica" {{if eq .Status "Não se aplica"}}selected{{end}}>Não se aplica</option>
						</select>
					</div>

					<div>
						<label>Observação</label>
						<textarea name="observacao_{{.ID}}" placeholder="Ex: produtor ficou de enviar, documento vencido, falta assinatura...">{{.Observacao}}</textarea>
					</div>
				</div>
			</div>
			{{end}}
		</div>
	</details>
	{{else}}
	<div class="aviso-vazio">
		Nenhum item de checklist gerado para esta reunião.
	</div>
	{{end}}

	<br>
	<button type="submit">Salvar checklist</button>
</form>

<div id="modalExplicacao" class="modal-explicacao">
	<div class="modal-conteudo">
		<span class="modal-fechar" onclick="fecharExplicacao()">&times;</span>
		<h3 id="modalTitulo"></h3>
		<p id="modalTexto"></p>
	</div>
</div>

<script>
	function abrirExplicacao(titulo, texto) {
		document.getElementById("modalTitulo").textContent = titulo;
		document.getElementById("modalTexto").textContent = texto;
		document.getElementById("modalExplicacao").style.display = "block";
	}

	function fecharExplicacao() {
		document.getElementById("modalExplicacao").style.display = "none";
	}

	window.onclick = function(event) {
		const modal = document.getElementById("modalExplicacao");
		if (event.target === modal) {
			fecharExplicacao();
		}
	}

	function textoNormalizado(valor) {
		return (valor || "")
			.toLowerCase()
			.normalize("NFD")
			.replace(/[\u0300-\u036f]/g, "");
	}

	function aplicarFiltrosChecklist() {
		const busca = textoNormalizado(document.getElementById("buscaChecklist").value);
		const status = document.getElementById("filtroStatus").value;

		document.querySelectorAll(".item-checklist").forEach(function(item) {
			const texto = textoNormalizado(item.dataset.texto);
			const statusItem = item.dataset.status || "";

			const passaBusca = busca === "" || texto.includes(busca);
			const passaStatus = status === "" || statusItem === status;

			item.style.display = passaBusca && passaStatus ? "" : "none";
		});

		document.querySelectorAll(".etapa-checklist").forEach(function(etapa) {
			const visiveis = etapa.querySelectorAll('.item-checklist:not([style*="display: none"])').length;
			etapa.style.display = visiveis > 0 ? "" : "none";

			if (busca !== "" || status !== "") {
				etapa.open = true;
			}
		});
	}

	function abrirTodasEtapas() {
		document.querySelectorAll(".etapa-checklist").forEach(function(etapa) {
			etapa.open = true;
		});
	}

	function fecharTodasEtapas() {
		document.querySelectorAll(".etapa-checklist").forEach(function(etapa) {
			etapa.open = false;
		});
	}

	function mostrarPendentes() {
		document.getElementById("filtroStatus").value = "Pendente";
		aplicarFiltrosChecklist();
	}

	function limparFiltrosChecklist() {
		document.getElementById("buscaChecklist").value = "";
		document.getElementById("filtroStatus").value = "";
		document.querySelectorAll(".etapa-checklist").forEach(function(etapa) {
			etapa.style.display = "";
			etapa.open = true;
		});
		document.querySelectorAll(".item-checklist").forEach(function(item) {
			item.style.display = "";
		});
	}

	function atualizarStatusVisual(select) {
		const item = select.closest(".item-checklist");
		const status = select.value;

		item.dataset.status = status;

		const selo = item.querySelector(".status-atual");
		if (selo) {
			selo.textContent = status;
		}

		aplicarFiltrosChecklist();
	}

	document.getElementById("buscaChecklist").addEventListener("input", aplicarFiltrosChecklist);
	document.getElementById("filtroStatus").addEventListener("change", aplicarFiltrosChecklist);
</script>
`
