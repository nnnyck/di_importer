package models

type DIData struct {
	MAWB   string
	HAWB   string
	Campo1 string
	Campo2 string
}

type DadosLinha struct {
	HAWB  string
	Index int
}

type DadosExtraidos struct {
	LinhaIndex  int
	HAWB        string
	MAWB        string
	DataChegada string
	VMLE        string
	Volumes     string
	PesoBruto   string
	PesoLiquido string
	Aeroporto   string
}
