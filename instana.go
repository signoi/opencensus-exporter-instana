package instana

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"go.opencensus.io/trace"
)

// Exporter satisfies the opencensus interface i.e trace.Exporter which has a single method
// i.e ExportSpan(span trace.SpanData)
type Exporter struct {
	// Service name
	serviceName string
	// Prefix
	prefix string
	// batch wraps around the slice buffer and
	batch *Batch
}

// Dispatcher defines the interface that holds single method
// i.e Dispatch whose whose job is to emit the jsonSpan to
// instana agent
type Dispatcher interface {
	Dispatch([]*jsonSpan) error
}

// TraceDispatcher implements
type TraceDispatcher struct {
	agentHost string
	agentPort int
}

// Dispatch transmits the slice of jsonSpan to instana agent
func (td *TraceDispatcher) Dispatch(jsonSpans []*jsonSpan) error {

	path := "com.instana.plugin.generic.trace"

	url := fmt.Sprintf("http://%s:%d/%s", td.agentHost, td.agentPort, path)
	fmt.Println("URL ", url)
	data, err := json.Marshal(jsonSpans)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
	if err != nil {
		return err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	b, err := ioutil.ReadAll(res.Body)
	fmt.Println("Response", string(b))

	if err != nil {
		return err
	}

	return nil
}

// Batch wraps around the buffer of slice of json spans and provides buffered approach
// for transmitting spans to instana agent
type Batch struct {
	instanaJSONSpans []*jsonSpan
	// TODO : add buffferbyte limit size for each individual instana json span
	limit       int
	dispatcher  Dispatcher
	servicename string
}

// Add appends span to buffer.
func (b *Batch) Add(span *jsonSpan) {
	b.instanaJSONSpans = append(b.instanaJSONSpans, span)
}

// This function runs in go routine
// that calls the handler whenever batchsize is met.
func (b *Batch) run() {
	for {
		select {
		case <-time.After(time.Second * 3):
			// check to see if we have the buffer size of
			if len(b.instanaJSONSpans) >= b.limit {
				// call the handler function
				// passing in the chunk
				err := b.dispatcher.Dispatch(b.instanaJSONSpans)
				if err != nil {
					log.Println("Failed to sent to instana agent")
				}
				if err == nil {
					// clear the buffer
					b.instanaJSONSpans = b.instanaJSONSpans[:0]
					log.Println("Log sent to instana agent")
				}
			}
		}
	}
}

// Compile time checks to satisfy trace.Exporter interface
var _ trace.Exporter = (*Exporter)(nil)

// ExportSpan implements the trace.Exporter interface defined
// by opencensus trace.
// This function converts opencensus span data to span specification
// of instana and appends it to instana JSON spans buffer
func (e *Exporter) ExportSpan(data *trace.SpanData) {
	instanaSpan := e.ToInstanaSpan(data)

	e.batch.Add(instanaSpan)
}

// ToInstanaSpan converts opencensus span data to span data around instana specification.
func (e *Exporter) ToInstanaSpan(data *trace.SpanData) *jsonSpan {
	// Drive ...
	fmt.Println("data", data.SpanID)
	jData := &jsonData{}
	jData.SDK = &jsonSDKData{
		Name: data.Name,
		Custom: &jsonCustomData{
			Tags: data.Attributes,
		},
	}
	jData.Service = e.serviceName

	jS := &jsonSpan{
		TraceID:   bytesToInt64(data.TraceID[0:8]),
		Data:      jData,
		SpanID:    bytesToInt64(data.SpanID[:]),
		Duration:  uint64(data.EndTime.Sub(data.StartTime).Nanoseconds()) / uint64(time.Millisecond),
		Kind:      data.SpanKind,
		Timestamp: uint64(data.StartTime.UnixNano()) / uint64(time.Millisecond),
		ParentID:  bytesToInt64(data.ParentSpanID[:]),
		Name:      "sdk",
		Error:     data.Status.Code != 0,
		Lang:      "go",
	}

	return jS
}

// NewExporter ...
func NewExporter(servicename, host string, port int) *Exporter {

	// Create batch struct
	batch := &Batch{
		limit: 1,
		dispatcher: &TraceDispatcher{
			agentHost: host,
			agentPort: port,
		},
	}

	go batch.run()

	return &Exporter{

		batch:       batch,
		serviceName: servicename,
	}
}

func bytesToInt64(buf []byte) int64 {
	u := binary.BigEndian.Uint64(buf)
	return int64(u)
}
