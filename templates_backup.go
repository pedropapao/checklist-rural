package main

const backupsHTML = `
<h2>Backups do banco de dados</h2>

<div class="card">
	<p>Esta tela cria e restaura cópias de segurança do banco SQLite local.</p>
	<p>Os backups ficam salvos em:</p>
	<pre>~/Documentos/ChecklistRural/backups</pre>
</div>

<form method="POST" action="/criar-backup">
	<button type="submit">Criar backup agora</button>
	<a class="botao secundario" href="/">Início</a>
</form>

<h3>Backups existentes</h3>

<table>
	<tr>
		<th>Arquivo</th>
		<th>Tamanho</th>
		<th>Data</th>
		<th>Caminho</th>
		<th>Ações</th>
	</tr>

	{{range .Backups}}
	<tr>
		<td>{{.Nome}}</td>
		<td>{{.Tamanho}} bytes</td>
		<td>{{.Modificado}}</td>
		<td><pre>{{.Caminho}}</pre></td>
		<td>
			<a class="botao perigo" href="/confirmar-restaurar-backup?arquivo={{.Nome}}">Restaurar</a>
		</td>
	</tr>
	{{else}}
	<tr>
		<td colspan="5">Nenhum backup criado ainda.</td>
	</tr>
	{{end}}
</table>
`

const confirmarRestaurarBackupHTML = `
<h2>Confirmar restauração de backup</h2>

<div class="card-perigo">
	<p><strong>Atenção:</strong> esta ação vai substituir o banco de dados atual pelo backup selecionado.</p>

	<p>Antes de restaurar, o sistema criará automaticamente uma cópia de segurança do banco atual.</p>

	<p><strong>Backup selecionado:</strong></p>
	<pre>{{.Arquivo}}</pre>

	<p><strong>Caminho:</strong></p>
	<pre>{{.Caminho}}</pre>
</div>

<form method="POST" action="/restaurar-backup">
	<input type="hidden" name="arquivo" value="{{.Arquivo}}">

	<button class="perigo" type="submit">Sim, restaurar este backup</button>
	<a class="botao secundario" href="/backups">Cancelar</a>
</form>
`
