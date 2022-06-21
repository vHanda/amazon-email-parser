package main

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"

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

	orders, err := parseHTML(e.HTML)
	if err != nil {
		log.Fatal(err)
	}

	var b bytes.Buffer
	csvWriter := csv.NewWriter(bufio.NewWriter(&b))

	var data [][]string

	header := []string{"ID", "name", "merchant", "currency", "price"}
	data = append(data, header)

	for _, o := range orders {
		row := []string{o.ID, o.Name, o.Merchant, o.Currency, o.Price.String()}
		data = append(data, row)
	}

	err = csvWriter.WriteAll(data)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(b.String())

	return 0
}
