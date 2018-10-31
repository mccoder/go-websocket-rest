package model

type WsRequestMessage struct {
	Id     string `json:"id"`
	Action string `json:"action"`
	Params string `json:"params,omitempty"`
}

func (m *WsRequestMessage) ParamsByteArray() []byte {
	return []byte(m.Params)
}

func (m *WsRequestMessage) Clean() {
	m.Action = ""
	m.Id = ""
	m.Params = ""
}
