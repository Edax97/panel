package panelServer

import (
	"fmt"
	"math"
	"strconv"
)

type TransformerParams struct {
	IaPrim  string
	IbPrim  string
	IcPrim  string
	VabPrim string
	VbcPrim string
	VcaPrim string
	Pprim   string
	Qprim   string
	Sprim   string
}

func NewTransformerParams() *TransformerParams {
	return &TransformerParams{"NaN", "NaN", "NaN", "NaN", "NaN", "NaN", "NaN", "NaN", "NaN"}
}

func CalcTransformer(IaStr, IbStr, IcStr, VabStr, VbcStr, VcaStr, PftlStr string) (p *TransformerParams) {
	p = NewTransformerParams()
	Ia, err := strconv.ParseFloat(IaStr, 64)
	if err != nil {
		return
	}
	Ib, err := strconv.ParseFloat(IbStr, 64)
	if err != nil {
		return
	}
	Ic, err := strconv.ParseFloat(IcStr, 64)
	if err != nil {
		return
	}
	Vab, err := strconv.ParseFloat(VabStr, 64)
	if err != nil {
		return
	}
	Vbc, err := strconv.ParseFloat(VbcStr, 64)
	if err != nil {
		return
	}
	Vca, err := strconv.ParseFloat(VcaStr, 64)
	if err != nil {
		return
	}
	Pftl, err := strconv.ParseFloat(PftlStr, 64)
	if err != nil {
		return
	}

	IaPrim := Ia / 25
	IbPrim := Ib / 25
	IcPrim := Ic / 25
	VabPrim := Vab * 25
	VbcPrim := Vbc * 25
	VcaPrim := Vca * 25
	Iavg := (Ia + Ib + Ic) / 3
	Vavg := (Vab + Vbc + Vca) / 3
	Sprim := Iavg * Vavg * math.Sqrt(3)
	Pprim := Sprim * Pftl
	Qprim := math.Sqrt(math.Max(Sprim*Sprim-Pprim*Pprim, 0))

	return &TransformerParams{
		IaPrim:  fmt.Sprintf("%f", IaPrim),
		IbPrim:  fmt.Sprintf("%f", IbPrim),
		IcPrim:  fmt.Sprintf("%f", IcPrim),
		VabPrim: fmt.Sprintf("%f", VabPrim),
		VbcPrim: fmt.Sprintf("%f", VbcPrim),
		VcaPrim: fmt.Sprintf("%f", VcaPrim),
		Pprim:   fmt.Sprintf("%f", Pprim),
		Qprim:   fmt.Sprintf("%f", Qprim),
		Sprim:   fmt.Sprintf("%f", Sprim),
	}
}
