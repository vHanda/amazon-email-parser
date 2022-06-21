package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/davecgh/go-spew/spew"
	"github.com/jhillyerd/enmime"
	"github.com/shopspring/decimal"
)

type orderInfo struct {
	Name     string
	Merchant string
	ID       string
	currency string
	price    decimal.Decimal
}

type dumper struct {
	errOut, stdOut io.Writer
	exit           exitFunc
}
type exitFunc func(int)

func newDefaultDumper() *dumper {
	return &dumper{
		errOut: os.Stderr,
		stdOut: os.Stdout,
		exit:   os.Exit,
	}
}

func main() {
	d := newDefaultDumper()
	d.exit(d.dump(os.Args))
}

func extractNames(doc *goquery.Document) []string {
	names := []string{}

	i := 0
	selector := "table#itemDetails a"
	doc.Find(selector).Each(func(_ int, s *goquery.Selection) {
		title := s.Text()
		title = strings.TrimSpace(title)
		if len(title) == 0 {
			return
		}
		if title == "Gestionado por Amazon" {
			return
		}

		if (i+1)%2 == 1 {
			names = append(names, title)
		}
		i++
	})

	return names
}

func extractMerchants(doc *goquery.Document) []string {
	names := []string{}

	i := 0
	selector := "table#itemDetails a"
	doc.Find(selector).Each(func(_ int, s *goquery.Selection) {
		title := s.Text()
		title = strings.TrimSpace(title)
		if len(title) == 0 {
			return
		}
		if title == "Gestionado por Amazon" {
			return
		}

		if (i+1)%2 == 0 {
			names = append(names, title)
		}
		i++
	})

	return names
}

func filter(inputs []string, notMatch []string) []string {
	outputs := []string{}

	for _, i := range inputs {
		ignore := false
		for _, not := range notMatch {
			if i == not {
				ignore = true
				break
			}
		}

		if !ignore {
			outputs = append(outputs, i)
		}
	}
	return outputs
}

func extractCssSelector(doc *goquery.Document, selector string) []string {
	arr := []string{}

	doc.Find(selector).Each(func(i int, s *goquery.Selection) {
		title := s.Text()
		title = strings.TrimSpace(title)
		if len(title) == 0 {
			return
		}
		arr = append(arr, title)
	})

	return arr
}

func extractCurrencyAndPrice(input string) (string, decimal.Decimal, error) {
	parts := strings.Split(input, " ")
	num := parts[1]
	num = strings.ReplaceAll(num, ".", "")
	num = strings.Replace(num, ",", ".", 1)

	d, err := decimal.NewFromString(num)
	if err != nil {
		return "", decimal.Decimal{}, err
	}

	return parts[0], d, nil
}

func (d *dumper) dump(args []string) int {
	if len(args) < 2 {
		fmt.Fprintln(d.errOut, "Missing filename argument")
		return 1
	}

	reader, err := os.Open(args[1])
	if err != nil {
		fmt.Fprintln(d.errOut, "Failed to open file:", err)
		return 1
	}

	e, err := enmime.ReadEnvelope(reader)
	if err != nil {
		fmt.Fprintf(d.errOut, "Failed to read envelope:\n%+v\n", err)
		return 1
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader((e.HTML)))
	if err != nil {
		log.Fatal(err)
	}

	names := extractNames(doc)
	merchants := extractMerchants(doc)
	ids := extractCssSelector(doc, "#orderDetails a")
	prices := filter(extractCssSelector(doc, "td.price"), []string{"Mastercard", "Visa"})

	orders := []orderInfo{}
	for i := range names {
		pStr := prices[i]
		currency, price, err := extractCurrencyAndPrice(pStr)
		if err != nil {
			log.Fatalln(err)
		}

		o := orderInfo{
			Name:     names[i],
			Merchant: merchants[i],
			ID:       ids[i],
			currency: currency,
			price:    price,
		}

		orders = append(orders, o)
	}

	spew.Dump(orders)

	return 0
}
