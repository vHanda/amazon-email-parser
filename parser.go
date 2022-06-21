package main

import (
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/shopspring/decimal"
)

type OrderInfo struct {
	Name     string
	Merchant string
	ID       string
	Currency string
	Price    decimal.Decimal
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

func parseHTML(html string) ([]OrderInfo, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader((html)))
	if err != nil {
		log.Fatal(err)
	}

	names := extractNames(doc)
	merchants := extractMerchants(doc)
	ids := extractCssSelector(doc, "#orderDetails a")
	prices := filter(extractCssSelector(doc, "td.price"), []string{"Mastercard", "Visa"})

	orders := []OrderInfo{}
	for i := range names {
		pStr := prices[i]
		currency, price, err := extractCurrencyAndPrice(pStr)
		if err != nil {
			log.Fatalln(err)
		}

		o := OrderInfo{
			Name:     names[i],
			Merchant: merchants[i],
			ID:       ids[i],
			Currency: currency,
			Price:    price,
		}

		orders = append(orders, o)
	}

	return orders, nil
}
