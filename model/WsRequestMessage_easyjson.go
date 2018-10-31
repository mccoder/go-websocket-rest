// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package model

import (
	json "encoding/json"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjsonA1eac286DecodeGithubComMccoderGoWebsocketRestModel(in *jlexer.Lexer, out *WsRequestMessage) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "id":
			out.Id = string(in.String())
		case "action":
			out.Action = string(in.String())
		case "params":
			out.Params = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonA1eac286EncodeGithubComMccoderGoWebsocketRestModel(out *jwriter.Writer, in WsRequestMessage) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"id\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Id))
	}
	{
		const prefix string = ",\"action\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Action))
	}
	if in.Params != "" {
		const prefix string = ",\"params\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Params))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v WsRequestMessage) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonA1eac286EncodeGithubComMccoderGoWebsocketRestModel(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v WsRequestMessage) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonA1eac286EncodeGithubComMccoderGoWebsocketRestModel(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *WsRequestMessage) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonA1eac286DecodeGithubComMccoderGoWebsocketRestModel(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *WsRequestMessage) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonA1eac286DecodeGithubComMccoderGoWebsocketRestModel(l, v)
}
