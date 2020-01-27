package main

import (
	"fmt"

	instana "github.com/signoi/opencensus-exporter-instana"
)

func main() {
	fmt.Println(instana.NewExporter("service name", "localhost", 4000))
}
