package main

import (
	"archive/zip"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
	"strings"
)

type ResumoMapaGeoref struct {
	TotalArquivos     int
	TotalFeatures     int
	TotalPontos       int
	TotalLinhas       int
	TotalPoligonos    int
	AreaTotalHectares float64
	Mensagens         []string
	PorArquivo        []ResumoArquivoMapa
	Talhoes           []ResumoTalhaoMapa
}

type ResumoArquivoMapa struct {
	Nome           string
	Tipo           string
	TotalFeatures  int
	TotalPontos    int
	TotalLinhas    int
	TotalPoligonos int
	AreaHectares   float64
	Mensagens      []string
}

type ResumoTalhaoMapa struct {
	Nome         string
	Origem       string
	Tipo         string
	AreaHectares float64
}

type kmlMapaDocumento struct {
	Placemarks []kmlMapaPlacemark `xml:"Document>Placemark"`
	Folders    []kmlMapaFolder    `xml:"Document>Folder"`
}

type kmlMapaFolder struct {
	Placemarks []kmlMapaPlacemark `xml:"Placemark"`
}

type kmlMapaPlacemark struct {
	Name       string `xml:"name"`
	Point      string `xml:"Point>coordinates"`
	LineString string `xml:"LineString>coordinates"`
	Polygon    string `xml:"Polygon>outerBoundaryIs>LinearRing>coordinates"`
}

func (app *App) montarGeoJSONDaReuniao(reuniaoID int) (string, ResumoMapaGeoref, error) {
	arquivos, err := app.listarArquivosGeoref(reuniaoID)
	if err != nil {
		return "", ResumoMapaGeoref{}, err
	}

	features := []map[string]any{}
	resumo := ResumoMapaGeoref{
		TotalArquivos: len(arquivos),
	}

	for _, arquivo := range arquivos {
		novas, mensagens := lerArquivoGeorefParaFeatures(arquivo)

		resumoArquivo := resumirArquivoMapa(arquivo, novas, mensagens)
		resumo.PorArquivo = append(resumo.PorArquivo, resumoArquivo)

		if len(mensagens) > 0 {
			resumo.Mensagens = append(resumo.Mensagens, mensagens...)
		}

		features = append(features, novas...)
	}

	resumo.Talhoes = enriquecerFeaturesComAreaETalhoes(features)

	fc := map[string]any{
		"type":     "FeatureCollection",
		"features": features,
	}

	normalizado, err := json.Marshal(fc)
	if err != nil {
		return "", resumo, err
	}

	var fcNormal map[string]any
	_ = json.Unmarshal(normalizado, &fcNormal)

	resumo = resumirFeatureCollectionMapa(fcNormal, resumo)
	resumo.TotalFeatures = len(features)

	return string(normalizado), resumo, nil
}

func resumirArquivoMapa(arquivo ArquivoGeoref, features []map[string]any, mensagens []string) ResumoArquivoMapa {
	fc := map[string]any{
		"type":     "FeatureCollection",
		"features": mapasParaAny(features),
	}

	parcial := ResumoMapaGeoref{}
	parcial = resumirFeatureCollectionMapa(fc, parcial)
	parcial.TotalFeatures = len(features)

	return ResumoArquivoMapa{
		Nome:           arquivo.NomeOriginal,
		Tipo:           arquivo.Tipo,
		TotalFeatures:  parcial.TotalFeatures,
		TotalPontos:    parcial.TotalPontos,
		TotalLinhas:    parcial.TotalLinhas,
		TotalPoligonos: parcial.TotalPoligonos,
		AreaHectares:   parcial.AreaTotalHectares,
		Mensagens:      mensagens,
	}
}

func mapasParaAny(features []map[string]any) []any {
	lista := []any{}

	for _, f := range features {
		lista = append(lista, f)
	}

	return lista
}

func enriquecerFeaturesComAreaETalhoes(features []map[string]any) []ResumoTalhaoMapa {
	talhoes := []ResumoTalhaoMapa{}

	for i, feature := range features {
		geom, ok := feature["geometry"].(map[string]any)
		if !ok {
			continue
		}

		props, ok := feature["properties"].(map[string]any)
		if !ok {
			props = map[string]any{}
		}

		tipo, _ := geom["type"].(string)
		coords := geom["coordinates"]

		nome, _ := props["name"].(string)
		if strings.TrimSpace(nome) == "" {
			nome = fmt.Sprintf("Talhão/área %d", i+1)
		}

		origem, _ := props["origem"].(string)

		area := 0.0

		switch strings.ToLower(tipo) {
		case "polygon":
			area = areaPoligonoHectares(coords)

		case "multipolygon":
			area = areaMultiPoligonoHectares(coords)
		}

		props["tipo_geometria"] = tipo

		if area > 0 {
			props["area_ha"] = area
		}

		feature["properties"] = props

		talhoes = append(talhoes, ResumoTalhaoMapa{
			Nome:         nome,
			Origem:       origem,
			Tipo:         tipo,
			AreaHectares: area,
		})
	}

	return talhoes
}

func lerArquivoGeorefParaFeatures(arquivo ArquivoGeoref) ([]map[string]any, []string) {
	ext := strings.ToLower(filepath.Ext(arquivo.Caminho))

	switch ext {
	case ".geojson", ".json":
		return lerGeoJSONFeatures(arquivo.Caminho, arquivo.NomeOriginal)

	case ".kml":
		conteudo, err := os.ReadFile(arquivo.Caminho)
		if err != nil {
			return nil, []string{"Erro ao ler KML " + arquivo.NomeOriginal + ": " + err.Error()}
		}
		return lerKMLFeatures(conteudo, arquivo.NomeOriginal)

	case ".kmz":
		return lerKMZFeatures(arquivo.Caminho, arquivo.NomeOriginal)

	case ".shp":
		conteudo, err := os.ReadFile(arquivo.Caminho)
		if err != nil {
			return nil, []string{"Erro ao ler SHP " + arquivo.NomeOriginal + ": " + err.Error()}
		}
		return lerSHPFeatures(conteudo, arquivo.NomeOriginal)

	case ".zip":
		return lerZIPGeograficoFeatures(arquivo.Caminho, arquivo.NomeOriginal)

	default:
		return nil, []string{"Arquivo " + arquivo.NomeOriginal + " guardado, mas ainda não é desenhado no mapa."}
	}
}

func lerGeoJSONFeatures(caminho string, nomeOrigem string) ([]map[string]any, []string) {
	conteudo, err := os.ReadFile(caminho)
	if err != nil {
		return nil, []string{"Erro ao ler GeoJSON " + nomeOrigem + ": " + err.Error()}
	}

	var raiz map[string]any
	if err := json.Unmarshal(conteudo, &raiz); err != nil {
		return nil, []string{"Erro ao interpretar GeoJSON " + nomeOrigem + ": " + err.Error()}
	}

	tipo, _ := raiz["type"].(string)
	tipo = strings.ToLower(tipo)

	features := []map[string]any{}

	switch tipo {
	case "featurecollection":
		lista, ok := raiz["features"].([]any)
		if !ok {
			return nil, []string{"GeoJSON " + nomeOrigem + " não possui lista de features."}
		}

		for _, item := range lista {
			if feature, ok := item.(map[string]any); ok {
				garantirPropriedadesOrigem(feature, nomeOrigem)
				features = append(features, feature)
			}
		}

	case "feature":
		garantirPropriedadesOrigem(raiz, nomeOrigem)
		features = append(features, raiz)

	default:
		if _, ok := raiz["coordinates"]; ok {
			features = append(features, criarFeatureMapa(nomeOrigem, nomeOrigem, raiz))
		}
	}

	return features, nil
}

func garantirPropriedadesOrigem(feature map[string]any, origem string) {
	props, ok := feature["properties"].(map[string]any)
	if !ok {
		props = map[string]any{}
	}

	if _, ok := props["origem"]; !ok {
		props["origem"] = origem
	}

	feature["properties"] = props
}

func lerKMLFeatures(conteudo []byte, nomeOrigem string) ([]map[string]any, []string) {
	var doc kmlMapaDocumento

	if err := xml.Unmarshal(conteudo, &doc); err != nil {
		return nil, []string{"Erro ao interpretar KML " + nomeOrigem + ": " + err.Error()}
	}

	placemarks := []kmlMapaPlacemark{}
	placemarks = append(placemarks, doc.Placemarks...)

	for _, pasta := range doc.Folders {
		placemarks = append(placemarks, pasta.Placemarks...)
	}

	features := []map[string]any{}

	for i, p := range placemarks {
		nome := strings.TrimSpace(p.Name)
		if nome == "" {
			nome = fmt.Sprintf("%s - item %d", nomeOrigem, i+1)
		}

		if strings.TrimSpace(p.Point) != "" {
			coords := parseCoordenadasKML(p.Point)
			if len(coords) > 0 {
				features = append(features, criarFeatureMapa(nome, nomeOrigem, map[string]any{
					"type":        "Point",
					"coordinates": coords[0],
				}))
			}
		}

		if strings.TrimSpace(p.LineString) != "" {
			coords := parseCoordenadasKML(p.LineString)
			if len(coords) > 1 {
				features = append(features, criarFeatureMapa(nome, nomeOrigem, map[string]any{
					"type":        "LineString",
					"coordinates": coords,
				}))
			}
		}

		if strings.TrimSpace(p.Polygon) != "" {
			coords := parseCoordenadasKML(p.Polygon)
			if len(coords) > 2 {
				coords = fecharAnelSeNecessario(coords)

				features = append(features, criarFeatureMapa(nome, nomeOrigem, map[string]any{
					"type":        "Polygon",
					"coordinates": []any{coords},
				}))
			}
		}
	}

	return features, nil
}

func parseCoordenadasKML(texto string) []any {
	texto = strings.TrimSpace(texto)
	partes := strings.Fields(texto)

	coords := []any{}

	for _, parte := range partes {
		valores := strings.Split(parte, ",")
		if len(valores) < 2 {
			continue
		}

		lon, ok1 := parseFloatSeguro(valores[0])
		lat, ok2 := parseFloatSeguro(valores[1])

		if ok1 && ok2 {
			coords = append(coords, []float64{lon, lat})
		}
	}

	return coords
}

func lerKMZFeatures(caminho string, nomeOrigem string) ([]map[string]any, []string) {
	zr, err := zip.OpenReader(caminho)
	if err != nil {
		return nil, []string{"Erro ao abrir KMZ " + nomeOrigem + ": " + err.Error()}
	}
	defer zr.Close()

	for _, f := range zr.File {
		if strings.HasSuffix(strings.ToLower(f.Name), ".kml") {
			rc, err := f.Open()
			if err != nil {
				return nil, []string{"Erro ao abrir KML dentro do KMZ " + nomeOrigem + ": " + err.Error()}
			}

			conteudo, err := io.ReadAll(rc)
			rc.Close()

			if err != nil {
				return nil, []string{"Erro ao ler KML dentro do KMZ " + nomeOrigem + ": " + err.Error()}
			}

			return lerKMLFeatures(conteudo, nomeOrigem)
		}
	}

	return nil, []string{"KMZ " + nomeOrigem + " não contém arquivo KML."}
}

func lerZIPGeograficoFeatures(caminho string, nomeOrigem string) ([]map[string]any, []string) {
	zr, err := zip.OpenReader(caminho)
	if err != nil {
		return nil, []string{"Erro ao abrir ZIP " + nomeOrigem + ": " + err.Error()}
	}
	defer zr.Close()

	features := []map[string]any{}
	mensagens := []string{}

	for _, f := range zr.File {
		nomeBaixo := strings.ToLower(f.Name)

		if !(strings.HasSuffix(nomeBaixo, ".kml") ||
			strings.HasSuffix(nomeBaixo, ".geojson") ||
			strings.HasSuffix(nomeBaixo, ".json") ||
			strings.HasSuffix(nomeBaixo, ".shp")) {
			continue
		}

		rc, err := f.Open()
		if err != nil {
			mensagens = append(mensagens, "Erro ao abrir "+f.Name+" dentro do ZIP: "+err.Error())
			continue
		}

		conteudo, err := io.ReadAll(rc)
		rc.Close()

		if err != nil {
			mensagens = append(mensagens, "Erro ao ler "+f.Name+" dentro do ZIP: "+err.Error())
			continue
		}

		switch {
		case strings.HasSuffix(nomeBaixo, ".kml"):
			fs, ms := lerKMLFeatures(conteudo, nomeOrigem+" / "+filepath.Base(f.Name))
			features = append(features, fs...)
			mensagens = append(mensagens, ms...)

		case strings.HasSuffix(nomeBaixo, ".geojson") || strings.HasSuffix(nomeBaixo, ".json"):
			tmp, err := os.CreateTemp("", "georef-*.geojson")
			if err != nil {
				mensagens = append(mensagens, "Erro temporário ao ler GeoJSON dentro do ZIP: "+err.Error())
				continue
			}

			_, _ = tmp.Write(conteudo)
			tmp.Close()

			fs, ms := lerGeoJSONFeatures(tmp.Name(), nomeOrigem+" / "+filepath.Base(f.Name))
			_ = os.Remove(tmp.Name())

			features = append(features, fs...)
			mensagens = append(mensagens, ms...)

		case strings.HasSuffix(nomeBaixo, ".shp"):
			fs, ms := lerSHPFeatures(conteudo, nomeOrigem+" / "+filepath.Base(f.Name))
			features = append(features, fs...)
			mensagens = append(mensagens, ms...)
		}
	}

	if len(features) == 0 && len(mensagens) == 0 {
		mensagens = append(mensagens, "ZIP "+nomeOrigem+" não contém KML, GeoJSON ou SHP reconhecido.")
	}

	return features, mensagens
}

func lerSHPFeatures(conteudo []byte, nomeOrigem string) ([]map[string]any, []string) {
	if len(conteudo) < 100 {
		return nil, []string{"SHP " + nomeOrigem + " inválido ou muito pequeno."}
	}

	r := bytes.NewReader(conteudo)
	_, _ = r.Seek(100, io.SeekStart)

	features := []map[string]any{}
	mensagens := []string{}

	for {
		header := make([]byte, 8)
		if _, err := io.ReadFull(r, header); err != nil {
			break
		}

		contentLengthWords := binary.BigEndian.Uint32(header[4:8])
		contentLengthBytes := int(contentLengthWords) * 2

		if contentLengthBytes <= 0 {
			continue
		}

		registro := make([]byte, contentLengthBytes)
		if _, err := io.ReadFull(r, registro); err != nil {
			break
		}

		fs, err := converterRegistroSHP(registro, nomeOrigem)
		if err != nil {
			mensagens = append(mensagens, err.Error())
			continue
		}

		features = append(features, fs...)
	}

	if len(features) == 0 && len(mensagens) == 0 {
		mensagens = append(mensagens, "SHP "+nomeOrigem+" não gerou geometria reconhecida.")
	}

	return features, mensagens
}

func converterRegistroSHP(registro []byte, nomeOrigem string) ([]map[string]any, error) {
	if len(registro) < 4 {
		return nil, fmt.Errorf("Registro SHP muito pequeno em %s", nomeOrigem)
	}

	shapeType := int32(binary.LittleEndian.Uint32(registro[0:4]))

	switch shapeType {
	case 0:
		return nil, nil

	case 1:
		if len(registro) < 20 {
			return nil, fmt.Errorf("Ponto SHP inválido em %s", nomeOrigem)
		}

		x := math.Float64frombits(binary.LittleEndian.Uint64(registro[4:12]))
		y := math.Float64frombits(binary.LittleEndian.Uint64(registro[12:20]))

		return []map[string]any{
			criarFeatureMapa(nomeOrigem, nomeOrigem, map[string]any{
				"type":        "Point",
				"coordinates": []float64{x, y},
			}),
		}, nil

	case 3, 5:
		if len(registro) < 44 {
			return nil, fmt.Errorf("Polyline/Polygon SHP inválido em %s", nomeOrigem)
		}

		numParts := int(int32(binary.LittleEndian.Uint32(registro[36:40])))
		numPoints := int(int32(binary.LittleEndian.Uint32(registro[40:44])))

		if numParts <= 0 || numPoints <= 0 {
			return nil, nil
		}

		partesInicio := 44
		pontosInicio := partesInicio + 4*numParts

		if len(registro) < pontosInicio+(16*numPoints) {
			return nil, fmt.Errorf("SHP incompleto em %s", nomeOrigem)
		}

		partes := make([]int, numParts)
		for i := 0; i < numParts; i++ {
			partes[i] = int(int32(binary.LittleEndian.Uint32(registro[partesInicio+i*4 : partesInicio+i*4+4])))
		}

		pontos := make([][]float64, numPoints)
		for i := 0; i < numPoints; i++ {
			pos := pontosInicio + i*16
			x := math.Float64frombits(binary.LittleEndian.Uint64(registro[pos : pos+8]))
			y := math.Float64frombits(binary.LittleEndian.Uint64(registro[pos+8 : pos+16]))
			pontos[i] = []float64{x, y}
		}

		if shapeType == 3 {
			features := []map[string]any{}

			for i := 0; i < numParts; i++ {
				inicio := partes[i]
				fim := numPoints

				if i+1 < numParts {
					fim = partes[i+1]
				}

				linha := []any{}
				for _, p := range pontos[inicio:fim] {
					linha = append(linha, p)
				}

				features = append(features, criarFeatureMapa(nomeOrigem, nomeOrigem, map[string]any{
					"type":        "LineString",
					"coordinates": linha,
				}))
			}

			return features, nil
		}

		aneis := []any{}

		for i := 0; i < numParts; i++ {
			inicio := partes[i]
			fim := numPoints

			if i+1 < numParts {
				fim = partes[i+1]
			}

			anel := []any{}
			for _, p := range pontos[inicio:fim] {
				anel = append(anel, p)
			}

			anel = fecharAnelSeNecessario(anel)
			aneis = append(aneis, anel)
		}

		return []map[string]any{
			criarFeatureMapa(nomeOrigem, nomeOrigem, map[string]any{
				"type":        "Polygon",
				"coordinates": aneis,
			}),
		}, nil

	default:
		return nil, fmt.Errorf("Tipo SHP %d ainda não suportado em %s", shapeType, nomeOrigem)
	}
}

func criarFeatureMapa(nome string, origem string, geometria map[string]any) map[string]any {
	return map[string]any{
		"type":     "Feature",
		"geometry": geometria,
		"properties": map[string]any{
			"name":   nome,
			"origem": origem,
		},
	}
}

func parseFloatSeguro(valor string) (float64, bool) {
	valor = strings.TrimSpace(valor)
	valor = strings.ReplaceAll(valor, ",", ".")

	var n float64
	_, err := fmt.Sscanf(valor, "%f", &n)

	return n, err == nil
}

func fecharAnelSeNecessario(coords []any) []any {
	if len(coords) < 2 {
		return coords
	}

	lon1, lat1, ok1 := coordenadaLonLat(coords[0])
	lon2, lat2, ok2 := coordenadaLonLat(coords[len(coords)-1])

	if ok1 && ok2 && (lon1 != lon2 || lat1 != lat2) {
		coords = append(coords, coords[0])
	}

	return coords
}

func resumirFeatureCollectionMapa(fc map[string]any, resumo ResumoMapaGeoref) ResumoMapaGeoref {
	features, ok := fc["features"].([]any)
	if !ok {
		return resumo
	}

	for _, item := range features {
		feature, ok := item.(map[string]any)
		if !ok {
			continue
		}

		geom, ok := feature["geometry"].(map[string]any)
		if !ok {
			continue
		}

		tipo, _ := geom["type"].(string)
		coords := geom["coordinates"]

		switch strings.ToLower(tipo) {
		case "point":
			resumo.TotalPontos++

		case "multipoint":
			resumo.TotalPontos += contarCoordenadasRecursivo(coords)

		case "linestring", "multilinestring":
			resumo.TotalLinhas++
			resumo.TotalPontos += contarCoordenadasRecursivo(coords)

		case "polygon":
			resumo.TotalPoligonos++
			resumo.TotalPontos += contarCoordenadasRecursivo(coords)
			resumo.AreaTotalHectares += areaPoligonoHectares(coords)

		case "multipolygon":
			resumo.TotalPoligonos++
			resumo.TotalPontos += contarCoordenadasRecursivo(coords)
			resumo.AreaTotalHectares += areaMultiPoligonoHectares(coords)
		}
	}

	return resumo
}

func contarCoordenadasRecursivo(valor any) int {
	if _, _, ok := coordenadaLonLat(valor); ok {
		return 1
	}

	lista, ok := valor.([]any)
	if !ok {
		return 0
	}

	total := 0
	for _, item := range lista {
		total += contarCoordenadasRecursivo(item)
	}

	return total
}

func coordenadaLonLat(valor any) (float64, float64, bool) {
	switch v := valor.(type) {
	case []any:
		if len(v) < 2 {
			return 0, 0, false
		}

		lon, ok1 := numeroParaFloat(v[0])
		lat, ok2 := numeroParaFloat(v[1])

		return lon, lat, ok1 && ok2

	case []float64:
		if len(v) < 2 {
			return 0, 0, false
		}

		return v[0], v[1], true
	}

	return 0, 0, false
}

func numeroParaFloat(valor any) (float64, bool) {
	switch v := valor.(type) {
	case float64:
		return v, true
	case float32:
		return float64(v), true
	case int:
		return float64(v), true
	case int64:
		return float64(v), true
	case json.Number:
		n, err := v.Float64()
		return n, err == nil
	default:
		return 0, false
	}
}

func areaPoligonoHectares(coords any) float64 {
	aneis, ok := coords.([]any)
	if !ok || len(aneis) == 0 {
		return 0
	}

	total := 0.0

	for i, anelValor := range aneis {
		anel, ok := anelValor.([]any)
		if !ok {
			continue
		}

		area := areaAnelHectares(anel)

		if i == 0 {
			total += area
		} else {
			total -= area
		}
	}

	if total < 0 {
		total = -total
	}

	return total
}

func areaMultiPoligonoHectares(coords any) float64 {
	poligonos, ok := coords.([]any)
	if !ok {
		return 0
	}

	total := 0.0

	for _, poligono := range poligonos {
		total += areaPoligonoHectares(poligono)
	}

	return total
}

func areaAnelHectares(anel []any) float64 {
	if len(anel) < 3 {
		return 0
	}

	latMedia := 0.0
	totalValidos := 0.0

	for _, c := range anel {
		_, lat, ok := coordenadaLonLat(c)
		if ok {
			latMedia += lat
			totalValidos++
		}
	}

	if totalValidos == 0 {
		return 0
	}

	latMedia = latMedia / totalValidos

	const raioTerra = 6378137.0
	lat0 := latMedia * math.Pi / 180.0

	pontos := [][2]float64{}

	for _, c := range anel {
		lon, lat, ok := coordenadaLonLat(c)
		if !ok {
			continue
		}

		x := raioTerra * (lon * math.Pi / 180.0) * math.Cos(lat0)
		y := raioTerra * (lat * math.Pi / 180.0)

		pontos = append(pontos, [2]float64{x, y})
	}

	if len(pontos) < 3 {
		return 0
	}

	soma := 0.0

	for i := 0; i < len(pontos); i++ {
		j := (i + 1) % len(pontos)
		soma += pontos[i][0]*pontos[j][1] - pontos[j][0]*pontos[i][1]
	}

	areaM2 := math.Abs(soma) / 2.0

	return areaM2 / 10000.0
}
