package communication

type Caller string

const (
	CallerUDPClient Caller = "UDPClient"
	CallerUDPServer Caller = "UDPServer"
	CallerWebServer Caller = "WebServer"
)

type InternalMessage struct {
	Caller  Caller `json:"caller"`
	Target  string `json:"target"`
	Content string `json:"content"`
}

func (m *InternalMessage) ToBytes() []byte {
	return []byte(m.Content)
}

func (m *InternalMessage) FromBytes(data []byte) {
	m.Content = string(data)
}

func (m *InternalMessage) ToJSON() string {
	return `{"caller":"` + string(m.Caller) + `","target":"` + m.Target + `","content":"` + m.Content + `"}`
}
