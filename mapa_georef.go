package main

import (
	"fmt"
	"html/template"
	"net/http"
)

func (app *App) telaMapaGeorreferenciamento(w http.ResponseWriter, r *http.Request) {
	reuniao, err := app.pegarReuniaoDaURL(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	geojson, resumo, err := app.montarGeoJSONDaReuniao(reuniao.ID)
	if err != nil {
		http.Error(w, "Erro ao montar mapa: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if r.URL.Query().Get("exportar") == "1" {
		nomeArquivo := fmt.Sprintf("mapa_georreferenciado_reuniao_%d.html", reuniao.ID)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Content-Disposition", "attachment; filename="+nomeArquivo)
	}

	dados := map[string]any{
		"Titulo":  "Mapa georreferenciado",
		"Reuniao": reuniao,
		"GeoJSON": template.JS(geojson),
		"Resumo":  resumo,
	}

	tpl := template.Must(template.New("mapa_georef").Parse(mapaGeorefHTML))
	tpl.Execute(w, dados)
}

const mapaGeorefHTML = `
<!DOCTYPE html>
<html lang="pt-br">
<head>
	<meta charset="UTF-8">
	<title>Mapa georreferenciado</title>

	<link rel="stylesheet" href="https://unpkg.com/leaflet@1.9.4/dist/leaflet.css">

	<style>
		:root {
			--verde: #1f7a4d;
			--verde-escuro: #145c38;
			--fundo: #f4f7f3;
			--card: #ffffff;
			--texto: #1f2933;
			--muted: #6b7280;
			--borda: #d9e2dc;
			--alerta: #fff7e6;
		}

		* {
			box-sizing: border-box;
		}

		body {
			margin: 0;
			font-family: Arial, Helvetica, sans-serif;
			background: var(--fundo);
			color: var(--texto);
		}

		header {
			background: linear-gradient(135deg, var(--verde), var(--verde-escuro));
			color: white;
			padding: 18px 24px;
		}

		header h1 {
			margin: 0;
			font-size: 24px;
		}

		header p {
			margin: 6px 0 0;
			color: #e7f5ec;
		}

		main {
			max-width: 1280px;
			margin: 0 auto;
			padding: 18px;
		}

		.barra-acoes {
			display: flex;
			flex-wrap: wrap;
			gap: 8px;
			margin-bottom: 16px;
		}

		.botao {
			display: inline-block;
			background: var(--verde);
			color: white;
			text-decoration: none;
			border-radius: 10px;
			padding: 10px 14px;
			font-weight: bold;
			font-size: 14px;
		}

		.botao.secundario {
			background: #eef6f0;
			color: var(--verde-escuro);
			border: 1px solid #cfe5d6;
		}

		.grid {
			display: grid;
			grid-template-columns: 320px 1fr;
			gap: 16px;
		}

		.card {
			background: var(--card);
			border: 1px solid var(--borda);
			border-radius: 14px;
			padding: 16px;
			box-shadow: 0 6px 18px rgba(15, 23, 42, 0.05);
		}

		.card h3 {
			margin: 0 0 10px;
			color: var(--verde-escuro);
		}

		.pequeno {
			font-size: 13px;
			color: var(--muted);
			line-height: 1.45;
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
			margin: 2px;
		}

		#mapa {
			height: calc(100vh - 150px);
			min-height: 560px;
			width: 100%;
			border-radius: 14px;
			border: 1px solid var(--borda);
			overflow: hidden;
		}

		.alerta {
			background: var(--alerta);
			border-color: #f2d28a;
		}

		ul {
			padding-left: 20px;
		}

		table {
			width: 100%;
			border-collapse: collapse;
			background: white;
			margin-top: 10px;
			border-radius: 10px;
			overflow: hidden;
			font-size: 12px;
		}

		th {
			background: #eef6f0;
			color: var(--verde-escuro);
			text-align: left;
			padding: 8px;
		}

		td {
			border-top: 1px solid var(--borda);
			padding: 8px;
			vertical-align: top;
		}


		@media (max-width: 900px) {
			.grid {
				grid-template-columns: 1fr;
			}

			#mapa {
				height: 560px;
			}
		}
	</style>
</head>
<body>
	<header>
		<h1>Mapa georreferenciado</h1>
		<p>{{.Reuniao.Produtor}} — {{.Reuniao.Municipio}}/{{.Reuniao.UF}}</p>
	</header>

	<main>
		<div class="barra-acoes">
			<a class="botao secundario" href="/georreferenciamento?id={{.Reuniao.ID}}">Voltar para georreferenciamento</a>
			<a class="botao secundario" href="/investigacao?id={{.Reuniao.ID}}">Voltar para investigação</a>
			<a class="botao secundario" href="/detalhes?id={{.Reuniao.ID}}">Detalhes</a>
			<a class="botao" href="/mapa-georef?id={{.Reuniao.ID}}&exportar=1">Exportar mapa HTML</a>
		</div>

		<div class="grid">
			<aside class="card">
				<h3>Resumo do mapa</h3>

				<p><strong>Arquivos analisados:</strong> {{.Resumo.TotalArquivos}}</p>
				<p><strong>Elementos desenhados:</strong> {{.Resumo.TotalFeatures}}</p>
				<p><strong>Pontos/coordenadas:</strong> {{.Resumo.TotalPontos}}</p>
				<p><strong>Linhas:</strong> {{.Resumo.TotalLinhas}}</p>
				<p><strong>Polígonos:</strong> {{.Resumo.TotalPoligonos}}</p>
				<p><strong>Área aproximada:</strong> {{printf "%.4f" .Resumo.AreaTotalHectares}} ha</p>

				<hr>

				<p class="pequeno">
					A área é uma estimativa calculada pelas coordenadas do arquivo.
					Para laudo técnico, confira projeção, datum e fonte do arquivo.
				</p>

				{{if .Resumo.Mensagens}}
				<div class="card alerta">
					<h3>Avisos</h3>
					<ul>
						{{range .Resumo.Mensagens}}
						<li class="pequeno">{{.}}</li>
						{{end}}
					</ul>
				</div>
				{{end}}

				<div>
					<span class="badge">GeoJSON</span>
					<span class="badge">KML</span>
					<span class="badge">KMZ</span>
					<span class="badge">SHP</span>
					<span class="badge">ZIP</span>
				</div>

				<hr>

				<h3>Área por arquivo</h3>

				{{if .Resumo.PorArquivo}}
				<table>
					<thead>
						<tr>
							<th>Arquivo</th>
							<th>Tipo</th>
							<th>Área ha</th>
						</tr>
					</thead>
					<tbody>
						{{range .Resumo.PorArquivo}}
						<tr>
							<td>{{.Nome}}</td>
							<td>{{.Tipo}}</td>
							<td>{{printf "%.4f" .AreaHectares}}</td>
						</tr>
						{{end}}
					</tbody>
				</table>
				{{else}}
				<p class="pequeno">Nenhum arquivo com geometria reconhecida.</p>
				{{end}}

				<hr>

				<h3>Talhões / feições</h3>

				{{if .Resumo.Talhoes}}
				<table>
					<thead>
						<tr>
							<th>Nome</th>
							<th>Tipo</th>
							<th>Área ha</th>
						</tr>
					</thead>
					<tbody>
						{{range .Resumo.Talhoes}}
						<tr>
							<td>{{.Nome}}</td>
							<td>{{.Tipo}}</td>
							<td>{{printf "%.4f" .AreaHectares}}</td>
						</tr>
						{{end}}
					</tbody>
				</table>
				{{else}}
				<p class="pequeno">Nenhum talhão identificado.</p>
				{{end}}
			</aside>

			<section>
				<div id="mapa"></div>
			</section>
		</div>
	</main>

	<script src="https://unpkg.com/leaflet@1.9.4/dist/leaflet.js"></script>

	<script>
		const geojson = {{.GeoJSON}};

		const mapa = L.map("mapa");

		L.tileLayer("https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png", {
			maxZoom: 20,
			attribution: "&copy; OpenStreetMap"
		}).addTo(mapa);

		function estiloFeature(feature) {
			return {
				weight: 3,
				fillOpacity: 0.25
			};
		}

		function popupFeature(feature, layer) {
			const props = feature.properties || {};
			const nome = props.name || props.nome || "Área";
			const origem = props.origem || "";
			const tipo = props.tipo_geometria || "";
			const area = props.area_ha || 0;

			let html = "<strong>" + nome + "</strong>";

			if (origem) {
				html += "<br><small>Origem: " + origem + "</small>";
			}

			if (tipo) {
				html += "<br><small>Tipo: " + tipo + "</small>";
			}

			if (area > 0) {
				html += "<br><strong>Área aprox.:</strong> " + Number(area).toFixed(4) + " ha";
			}

			layer.bindPopup(html);
		}

		const camada = L.geoJSON(geojson, {
			style: estiloFeature,
			onEachFeature: popupFeature,
			pointToLayer: function(feature, latlng) {
				return L.circleMarker(latlng, {
					radius: 7,
					weight: 2,
					fillOpacity: 0.8
				});
			}
		}).addTo(mapa);

		const bounds = camada.getBounds();

		if (bounds.isValid()) {
			mapa.fitBounds(bounds, {
				padding: [30, 30]
			});
		} else {
			mapa.setView([-15.78, -47.93], 4);
		}
	</script>
</body>
</html>
`
