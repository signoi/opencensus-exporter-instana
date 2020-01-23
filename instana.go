package instana

import (
	"time"

	"go.opencensus.io/trace"
)

// Exporter satisfies the opencensus interface i.e trace.Exporter which has a single method
// i.e ExportSpan(span trace.SpanData)
type Exporter struct {
	// Buffer for json spans that will
	// eventually be transmitted to instana agent
	instanaJSONSpans []*jsonSpan
	// Agent host
	agentHost string
	// Agent port
	agentPort int
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

// DispatcherFunc is function type that implements Dispatcher interface.
type DispatcherFunc func([]*jsonSpan) error

// Dispatch is a method on DispatcherFunc that implemetsn Dispatcher interface.
func (df DispatcherFunc) Dispatch(jsonSpans []*jsonSpan) error {
	return df(jsonSpans)
}

var defaultDispatcher = DispatcherFunc(func(jsonSpans []*jsonSpan) error {
	return nil
})

type Batch struct {
	instanaJSONSpans []*jsonSpan
	// TODO : add buffferbyte limit size for each individual instana json span
	limit      int
	dispatcher Dispatcher
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
				if err == nil {
					// clear the buffer
					b.instanaJSONSpans = b.instanaJSONSpans[:0]
				}
			}
		}
	}
}

var _ trace.Exporter = (*Exporter)(nil)

// ExportSpan implements the trace.Exporter interface defined
// by opencensus trace.
// This function converts opencensus span data to span specification
// of instana and appends it to instana JSON spans buffer
func (e *Exporter) ExportSpan(data *trace.SpanData) {
	instanaSpan := ToInstanaSpan(data)

	e.batch.Add(instanaSpan)
}

// ToInstanaSpan converts opencensus span data to span data around instana specification.
func ToInstanaSpan(data *trace.SpanData) *jsonSpan {

	return nil
}

// NewExporter ...
func NewExporter(host string, port int) *Exporter {
	// Create batch struct
	batch := &Batch{
		limit:      1,
		dispatcher: defaultDispatcher,
	}

	go batch.run()

	return &Exporter{
		agentHost: host,
		agentPort: port,
		batch:     batch,
	}
}
