package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/jhillyerd/enmime"
)

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

	selector := "table#itemDetails a"
	doc.Find(selector).Each(func(i int, s *goquery.Selection) {
		title := s.Text()
		fmt.Printf("Review %d: %s\n", i, title)
	})

	orderNumberSelector := "#orderDetails a"
	doc.Find(orderNumberSelector).Each(func(i int, s *goquery.Selection) {
		fmt.Printf("Order: %s\n", s.Text())
		// fmt.Printf("Date: %s\n", s.Find("span").Text())
	})
	return 0
}

//
