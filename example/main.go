package main

import instana "github.com/robusgauli/opencensus-exporter-instana"

func main() {
	instana.NewExporter("localhost", 4000)
	for {

	}
}
