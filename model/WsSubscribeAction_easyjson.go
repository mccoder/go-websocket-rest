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

func easyjsonE495e2aeDecodeGithubComMccoderGoWebsocketRestModel(in *jlexer.Lexer, out *WsSubscribeAction) {
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
		case "topicId":
			out.TopicId = string(in.String())
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
func easyjsonE495e2aeEncodeGithubComMccoderGoWebsocketRestModel(out *jwriter.Writer, in WsSubscribeAction) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"topicId\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.TopicId))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v WsSubscribeAction) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonE495e2aeEncodeGithubComMccoderGoWebsocketRestModel(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v WsSubscribeAction) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonE495e2aeEncodeGithubComMccoderGoWebsocketRestModel(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *WsSubscribeAction) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonE495e2aeDecodeGithubComMccoderGoWebsocketRestModel(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *WsSubscribeAction) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonE495e2aeDecodeGithubComMccoderGoWebsocketRestModel(l, v)
}
