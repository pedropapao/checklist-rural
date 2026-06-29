package main

type ResumoArquivosGeoref struct {
	Total    int
	Arquivos []ArquivoGeoref
}

func (app *App) resumoArquivosGeorreferenciamento(reuniaoID int) ResumoArquivosGeoref {
	arquivos, err := app.listarArquivosGeoref(reuniaoID)
	if err != nil {
		return ResumoArquivosGeoref{
			Total:    0,
			Arquivos: []ArquivoGeoref{},
		}
	}

	return ResumoArquivosGeoref{
		Total:    len(arquivos),
		Arquivos: arquivos,
	}
}
