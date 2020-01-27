package instana

type jsonSpan struct {
	TraceID   int64     `json:"t"`
	ParentID  int64     `json:"p,omitempty"`
	SpanID    int64     `json:"s"`
	Timestamp uint64    `json:"ts"`
	Duration  uint64    `json:"d"`
	Name      string    `json:"n"`
	From      *fromS    `json:"f"`
	Kind      int       `json:"k"`
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
