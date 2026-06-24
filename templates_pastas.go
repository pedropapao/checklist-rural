package main

const pastasHTML = `
<h2>Pastas do sistema</h2>

<div class="card">
	<p>Estas são as pastas usadas pelo sistema neste computador.</p>
</div>

<div class="card">
	<h3>Pasta principal</h3>
	<p>Guarda o banco de dados SQLite.</p>
	<pre>{{.PastaDados}}</pre>
	<a class="botao" href="/abrir-pasta-dados">Abrir pasta principal</a>
</div>

<div class="card">
	<h3>Exports</h3>
	<p>Guarda os arquivos TXT exportados.</p>
	<pre>{{.PastaExports}}</pre>
	<a class="botao" href="/abrir-pasta-exports">Abrir pasta de exports</a>
</div>

<div class="card">
	<h3>Backups</h3>
	<p>Guarda os backups do banco de dados.</p>
	<pre>{{.PastaBackups}}</pre>
	<a class="botao" href="/abrir-pasta-backups">Abrir pasta de backups</a>
</div>

<a class="botao secundario" href="/">Início</a>
`
