package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type ArquivoGeoref struct {
	ID           int
	ReuniaoID    int
	NomeOriginal string
	NomeSalvo    string
	Tipo         string
	Caminho      string
	Observacao   string
	CriadoEm     string
}

func pastaGeorreferenciamentoReuniao(reuniaoID int) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	pasta := filepath.Join(
		home,
		"Documentos",
		"ChecklistRural",
		"georreferenciamento",
		fmt.Sprintf("reuniao_%d", reuniaoID),
	)

	err = os.MkdirAll(pasta, 0755)
	if err != nil {
		return "", err
	}

	return pasta, nil
}

func extensaoPermitidaGeoref(nome string) bool {
	ext := strings.ToLower(filepath.Ext(nome))

	permitidas := map[string]bool{
		".kml":     true,
		".kmz":     true,
		".geojson": true,
		".json":    true,
		".gpx":     true,
		".csv":     true,
		".shp":     true,
		".dbf":     true,
		".shx":     true,
		".prj":     true,
		".pdf":     true,
		".zip":     true,
	}

	return permitidas[ext]
}

func tipoArquivoGeoref(nome string) string {
	ext := strings.ToLower(filepath.Ext(nome))

	switch ext {
	case ".kml":
		return "KML"
	case ".kmz":
		return "KMZ"
	case ".geojson", ".json":
		return "GeoJSON/JSON"
	case ".gpx":
		return "GPX"
	case ".csv":
		return "CSV"
	case ".shp", ".dbf", ".shx", ".prj":
		return "Shapefile"
	case ".pdf":
		return "PDF/planta"
	case ".zip":
		return "ZIP"
	default:
		return "Outro"
	}
}

func limparNomeArquivoGeoref(nome string) string {
	nome = filepath.Base(nome)
	nome = strings.TrimSpace(nome)

	trocador := strings.NewReplacer(
		" ", "_",
		"/", "_",
		"\\", "_",
		":", "_",
		"*", "_",
		"?", "_",
		"\"", "_",
		"<", "_",
		">", "_",
		"|", "_",
	)

	nome = trocador.Replace(nome)

	if nome == "" || nome == "." {
		nome = "arquivo"
	}

	return nome
}

func (app *App) listarArquivosGeoref(reuniaoID int) ([]ArquivoGeoref, error) {
	linhas, err := app.DB.Query(`
		SELECT
			id,
			reuniao_id,
			COALESCE(nome_original, ''),
			COALESCE(nome_salvo, ''),
			COALESCE(tipo, ''),
			COALESCE(caminho, ''),
			COALESCE(observacao, ''),
			COALESCE(criado_em, '')
		FROM arquivos_georef
		WHERE reuniao_id = ?
		ORDER BY id DESC
	`, reuniaoID)

	if err != nil {
		return nil, err
	}

	defer linhas.Close()

	arquivos := []ArquivoGeoref{}

	for linhas.Next() {
		var a ArquivoGeoref

		err := linhas.Scan(
			&a.ID,
			&a.ReuniaoID,
			&a.NomeOriginal,
			&a.NomeSalvo,
			&a.Tipo,
			&a.Caminho,
			&a.Observacao,
			&a.CriadoEm,
		)

		if err != nil {
			return nil, err
		}

		arquivos = append(arquivos, a)
	}

	return arquivos, nil
}

func (app *App) salvarArquivoGeorefRegistro(a ArquivoGeoref) error {
	_, err := app.DB.Exec(`
		INSERT INTO arquivos_georef (
			reuniao_id,
			nome_original,
			nome_salvo,
			tipo,
			caminho,
			observacao,
			criado_em
		)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`,
		a.ReuniaoID,
		a.NomeOriginal,
		a.NomeSalvo,
		a.Tipo,
		a.Caminho,
		a.Observacao,
		a.CriadoEm,
	)

	return err
}

func (app *App) excluirArquivoGeoref(id int, reuniaoID int) error {
	var caminho string

	err := app.DB.QueryRow(`
		SELECT COALESCE(caminho, '')
		FROM arquivos_georef
		WHERE id = ? AND reuniao_id = ?
	`, id, reuniaoID).Scan(&caminho)

	if err == sql.ErrNoRows {
		return nil
	}

	if err != nil {
		return err
	}

	if caminho != "" {
		_ = os.Remove(caminho)
	}

	_, err = app.DB.Exec(`
		DELETE FROM arquivos_georef
		WHERE id = ? AND reuniao_id = ?
	`, id, reuniaoID)

	return err
}

func (app *App) telaGeorreferenciamento(w http.ResponseWriter, r *http.Request) {
	idTexto := r.URL.Query().Get("id")

	reuniaoID, err := strconv.Atoi(idTexto)
	if err != nil || reuniaoID <= 0 {
		http.Error(w, "ID da reunião inválido", http.StatusBadRequest)
		return
	}

	reuniao, err := app.buscarReuniaoPorID(reuniaoID)
	if err != nil {
		http.Error(w, "Reunião não encontrada: "+err.Error(), http.StatusNotFound)
		return
	}

	mensagem := ""

	if r.Method == http.MethodPost {
		acao := r.FormValue("acao")

		switch acao {
		case "upload":
			mensagem = app.processarUploadGeoref(w, r, reuniaoID)

		case "excluir":
			arquivoID, _ := strconv.Atoi(r.FormValue("arquivo_id"))
			err := app.excluirArquivoGeoref(arquivoID, reuniaoID)
			if err != nil {
				http.Error(w, "Erro ao excluir arquivo: "+err.Error(), http.StatusInternalServerError)
				return
			}
			http.Redirect(w, r, "/georreferenciamento?id="+strconv.Itoa(reuniaoID), http.StatusSeeOther)
			return
		}
	}

	arquivos, err := app.listarArquivosGeoref(reuniaoID)
	if err != nil {
		http.Error(w, "Erro ao listar arquivos: "+err.Error(), http.StatusInternalServerError)
		return
	}

	dados := map[string]any{
		"Titulo":   "Georreferenciamento",
		"Reuniao":  reuniao,
		"Arquivos": arquivos,
		"Mensagem": mensagem,
	}

	tpl := template.Must(template.New("georef").Parse(htmlBase(georreferenciamentoHTML)))
	tpl.Execute(w, dados)
}

func (app *App) processarUploadGeoref(w http.ResponseWriter, r *http.Request, reuniaoID int) string {
	err := r.ParseMultipartForm(60 << 20) // 60 MB
	if err != nil {
		return "Erro ao ler upload: " + err.Error()
	}

	arquivo, cabecalho, err := r.FormFile("arquivo")
	if err != nil {
		return "Selecione um arquivo para importar."
	}
	defer arquivo.Close()

	if !extensaoPermitidaGeoref(cabecalho.Filename) {
		return "Tipo de arquivo não permitido. Use KML, KMZ, GeoJSON, GPX, CSV, SHP/DBF/SHX/PRJ, PDF ou ZIP."
	}

	pasta, err := pastaGeorreferenciamentoReuniao(reuniaoID)
	if err != nil {
		return "Erro ao criar pasta da reunião: " + err.Error()
	}

	nomeOriginal := limparNomeArquivoGeoref(cabecalho.Filename)
	ext := strings.ToLower(filepath.Ext(nomeOriginal))
	base := strings.TrimSuffix(nomeOriginal, ext)

	nomeSalvo := fmt.Sprintf(
		"%s_%s%s",
		base,
		time.Now().Format("20060102_150405"),
		ext,
	)

	caminhoFinal := filepath.Join(pasta, nomeSalvo)

	destino, err := os.Create(caminhoFinal)
	if err != nil {
		return "Erro ao criar arquivo no disco: " + err.Error()
	}
	defer destino.Close()

	_, err = io.Copy(destino, arquivo)
	if err != nil {
		return "Erro ao salvar arquivo: " + err.Error()
	}

	registro := ArquivoGeoref{
		ReuniaoID:    reuniaoID,
		NomeOriginal: nomeOriginal,
		NomeSalvo:    nomeSalvo,
		Tipo:         tipoArquivoGeoref(nomeOriginal),
		Caminho:      caminhoFinal,
		Observacao:   r.FormValue("observacao"),
		CriadoEm:     time.Now().Format("02/01/2006 15:04"),
	}

	err = app.salvarArquivoGeorefRegistro(registro)
	if err != nil {
		return "Arquivo salvo na pasta, mas erro ao registrar no banco: " + err.Error()
	}

	return "Arquivo importado com sucesso."
}

const georreferenciamentoHTML = `
<h2>Georreferenciamento</h2>

<div class="barra-acoes">
	<a class="botao secundario" href="/investigacao?id={{.Reuniao.ID}}">Voltar para investigação</a>
	<a class="botao secundario" href="/detalhes?id={{.Reuniao.ID}}">Voltar para detalhes</a>
	<a class="botao" href="/mapa-georef?id={{.Reuniao.ID}}">Ver mapa</a>
</div>

<div class="card destaque">
	<h3>{{.Reuniao.Produtor}}</h3>
	<p class="pequeno">
		{{.Reuniao.Municipio}}/{{.Reuniao.UF}} — {{.Reuniao.Banco}} — {{.Reuniao.Atividade}}
	</p>
</div>

{{if .Mensagem}}
<div class="card destaque">
	<p>{{.Mensagem}}</p>
</div>
{{end}}

<div class="grid">
	<div class="card">
		<h3>Importar localização da área</h3>
		<p class="pequeno">
			Anexe KML, KMZ, GeoJSON, shapefile, croqui, planta, PDF ou ZIP.
			O sistema guarda tudo dentro da pasta da reunião.
		</p>

		<form method="POST" action="/georreferenciamento?id={{.Reuniao.ID}}" enctype="multipart/form-data">
			<input type="hidden" name="acao" value="upload">

			<label>Arquivo</label>
			<input type="file" name="arquivo" required>

			<label>Observação</label>
			<textarea name="observacao" rows="4" placeholder="Ex: área do projeto, talhão soja, croqui do produtor, planta do imóvel, arquivo do CAR"></textarea>

			<button type="submit">Importar arquivo</button>
		</form>
	</div>

	<div class="card alerta">
		<h3>O que devo anexar?</h3>
		<p class="pequeno">
			Use qualquer documento que ajude a localizar a área financiada.
		</p>

		<ul>
			<li>KML/KMZ do CAR ou talhão</li>
			<li>GeoJSON ou shapefile</li>
			<li>PDF de planta/croqui</li>
			<li>Mapa enviado pelo produtor</li>
			<li>ZIP com arquivos geográficos</li>
		</ul>
	</div>
</div>

<div class="card">
	<h3>Arquivos importados</h3>

	{{if .Arquivos}}
	<table>
		<thead>
			<tr>
				<th>Tipo</th>
				<th>Arquivo</th>
				<th>Data</th>
				<th>Observação / leitura automática</th>
				<th>Ação</th>
			</tr>
		</thead>
		<tbody>
			{{range .Arquivos}}
			<tr>
				<td><span class="badge">{{.Tipo}}</span></td>
				<td>{{.NomeOriginal}}</td>
				<td>{{.CriadoEm}}</td>
				<td><pre>{{.Observacao}}</pre></td>
				<td>
					<form method="POST" action="/georreferenciamento?id={{$.Reuniao.ID}}" onsubmit="return confirm('Excluir este arquivo?')">
						<input type="hidden" name="acao" value="excluir">
						<input type="hidden" name="arquivo_id" value="{{.ID}}">
						<button type="submit" class="perigo">Excluir</button>
					</form>
				</td>
			</tr>
			{{end}}
		</tbody>
	</table>
	{{else}}
	<div class="card alerta">
		<h3>Nenhum arquivo importado ainda</h3>
		<p>
			Importe pelo menos um arquivo de localização antes de fechar a investigação do projeto.
		</p>
	</div>
	{{end}}
</div>

<div class="card">
	<h3>Próxima conferência</h3>
	<p class="pequeno">
		Depois de importar, confira se o arquivo corresponde ao imóvel, CAR, talhão, gleba ou área financiada.
		Se o produtor não tiver arquivo geográfico, anexe ao menos croqui, planta ou PDF de localização.
	</p>
</div>
`
