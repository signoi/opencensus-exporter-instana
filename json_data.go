package instana

type jsonSpan struct {
	TraceID   string    `json:"traceId"`
	ParentID  string    `json:"parentId,omitempty"`
	SpanID    string    `json:"spanId"`
	Timestamp uint64    `json:"timestamp"`
	Duration  uint64    `json:"duration"`
	Name      string    `json:"name"`
	Kind      int       `json:"kind"`
	Error     bool      `json:"error"`
	Ec        int       `json:"ec,omitempty"`
	Lang      string    `json:"ta,omitempty"`
	Data      *jsonData `json:"data"`
}

type jsonData struct {
	Service string       `json:"service,omitempty"`
	SDK     *jsonSDKData `json:"sdk"`
}

type jsonCustomData struct {
	Tags    map[string]interface{}            `json:"tags,omitempty"`
	Logs    map[uint64]map[string]interface{} `json:"logs,omitempty"`
	Baggage map[string]string                 `json:"baggage,omitempty"`
}

type jsonSDKData struct {
	Name      string          `json:"name"`
	Type      string          `json:"type,omitempty"`
	Arguments string          `json:"arguments,omitempty"`
	Return    string          `json:"return,omitempty"`
	Custom    *jsonCustomData `json:"custom,omitempty"`
}

type fromS struct {
	PID    string `json:"e"`
	HostID string `json:"h"`
}
