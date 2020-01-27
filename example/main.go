package main

import (
	"fmt"

	instana "github.com/robusgauli/opencensus-exporter-instana"
)

func main() {
	exporter := instana.NewExporter("name", "localhost", 3000)
	fmt.Println(exporter)
	for {

	}
}
