package model

const TypeTopic = "topic"

type WsResponceMessage struct {
	Id    string      `json:"id"`
	Type  string      `json:"type,omitempty"`
	Error string      `json:"error,omitempty"`
	Code  int16       `json:"code"`
	Data  interface{} `json:"data,omitempty"`
}

func (m *WsResponceMessage) Clear() {
	m.Code = 500
	m.Error = ""
	m.Data = nil
	m.Type = ""
}
