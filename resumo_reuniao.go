package main

type ReuniaoComResumo struct {
	Reuniao
	Resumo ResumoChecklist
}

func (app *App) resumoDaReuniao(reuniao Reuniao) ResumoChecklist {
	err := app.garantirItensChecklist(reuniao)
	if err != nil {
		return ResumoChecklist{}
	}

	itens, err := app.listarItensChecklist(reuniao.ID)
	if err != nil {
		return ResumoChecklist{}
	}

	return calcularResumoChecklist(itens)
}
