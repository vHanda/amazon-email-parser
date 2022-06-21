package main

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/shopspring/decimal"
	"gotest.tools/v3/assert"
)

func TestSpain(t *testing.T) {
	dataPath := filepath.Join("testdata", "2022-06-21-es.html")
	htmlData, err := ioutil.ReadFile(dataPath)
	assert.NilError(t, err)

	orders, err := parseHTML(string(htmlData))
	assert.NilError(t, err)

	assert.DeepEqual(t, orders, []OrderInfo{
		{
			ID:       "407-5248971-1092319",
			Name:     "PONY DANCE Cortinas Dormitorio Matrimonio - Cortinas Opacas Terimicas Aislantes para Salón con Trabillas, 2 Piezas, 140 x 160 cm, Moca",
			Merchant: "RYB HOME EU",
			Currency: "EUR",
			Price:    decimal.RequireFromString("26.95"),
		},
		{
			ID:       "407-0402347-9633967",
			Name:     "Aceite de Sésamo bio 500 ml Naturitas | Primera presión en frío | Sabor intenso y aromático | Sin conservantes",
			Merchant: "CURAE SOLUTIONS SL",
			Currency: "EUR",
			Price:    decimal.RequireFromString("8.49"),
		},
	})
}
