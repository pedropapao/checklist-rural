package main

const whatsappPendenciasHTML = `
<h2>Pendências para WhatsApp</h2>

<div class="card">
	<strong>Produtor:</strong> {{.Reuniao.Produtor}}<br>
	<strong>Banco:</strong> {{.Reuniao.Banco}}<br>
	<strong>Tipo:</strong> {{.Reuniao.TipoProjeto}}<br>
	<strong>Atividade:</strong> {{.Reuniao.Atividade}}<br>
</div>

<a class="botao secundario" href="/checklist-controle?id={{.Reuniao.ID}}">Voltar ao controle</a>
<a class="botao" href="/detalhes?id={{.Reuniao.ID}}">Detalhes</a>
<a class="botao" href="/exportar-whatsapp-pendencias-txt?id={{.Reuniao.ID}}">Exportar pendências TXT</a>

<p class="pequeno">Copie a mensagem abaixo e envie ao produtor pelo WhatsApp.</p>

<pre>{{.Mensagem}}</pre>
`
