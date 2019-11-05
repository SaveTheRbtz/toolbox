package rp_model_impl

import (
	"bytes"
	"encoding/json"
	"github.com/watermint/toolbox/infra/control/app_control"
	"github.com/watermint/toolbox/infra/report/rp_model"
	"github.com/watermint/toolbox/infra/ui/app_msg"
	"go.uber.org/zap"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"sync"
)

func NewJsonForQuiet(name string, ctl app_control.Control) (r rp_model.Report, err error) {
	r = &Json{
		w:   os.Stdout,
		ctl: ctl,
	}
	return r, nil
}

func NewJson(name string, ctl app_control.Control, opts ...rp_model.ReportOpt) (r rp_model.Report, err error) {
	l := ctl.Log()
	p := filepath.Join(ctl.Workspace().Report(), name+".json")
	l.Debug("Create new csv report", zap.String("path", p))
	f, err := os.Create(p)
	if err != nil {
		l.Error("Unable to create file", zap.String("path", p), zap.Error(err))
		return nil, err
	}
	r = &Json{
		file: f,
		w:    f,
		ctl:  ctl,
	}
	return r, nil
}

type Json struct {
	file  *os.File
	w     io.Writer
	ctl   app_control.Control
	mutex sync.Mutex
}

func (z *Json) findRaw(row interface{}, orig interface{}) json.RawMessage {
	var rv reflect.Value
	switch r := row.(type) {
	case reflect.Value:
		rv = r
	default:
		rv = reflect.ValueOf(row)
		if rv.Kind() == reflect.Ptr {
			rv = rv.Elem()
		}
	}
	rt := rv.Type()
	_, e := rt.FieldByName("Raw")
	if !e {
		return nil
	}
	rvf := rv.FieldByName("Raw")
	if rvf.Type().Kind() != reflect.TypeOf(json.RawMessage{}).Kind() {
		return nil
	}
	return rvf.Bytes()
}

func (z *Json) Success(input interface{}, result interface{}) {
	ui := z.ctl.UI()
	z.Row(rp_model.TransactionRow{
		Status: ui.Text(rp_model.MsgSuccess.Key(), rp_model.MsgSuccess.Params()...),
		Input:  input,
		Result: result,
	})
}

func (z *Json) Failure(reason app_msg.Message, input interface{}, result interface{}) {
	ui := z.ctl.UI()
	z.Row(rp_model.TransactionRow{
		Status: ui.Text(rp_model.MsgFailure.Key(), rp_model.MsgFailure.Params()...),
		Reason: ui.Text(reason.Key(), reason.Params()...),
		Input:  input,
		Result: result,
	})
}

func (z *Json) Skip(reason app_msg.Message, input interface{}, result interface{}) {
	ui := z.ctl.UI()
	z.Row(rp_model.TransactionRow{
		Status: ui.Text(rp_model.MsgSkip.Key(), rp_model.MsgFailure.Params()...),
		Reason: ui.Text(reason.Key(), reason.Params()...),
		Input:  input,
		Result: result,
	})
}

func (z *Json) Row(row interface{}) {
	z.mutex.Lock()
	defer z.mutex.Unlock()

	raw := z.findRaw(row, row)
	if raw != nil {
		z.w.Write(raw)
		z.w.Write([]byte("\n"))
		return
	}

	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "")
	err := enc.Encode(row)
	if err != nil {
		z.ctl.Log().Debug("Unable to unmarshal", zap.Error(err))
		return
	}
	z.w.Write(buf.Bytes())
}

func (z *Json) Close() {
	if z.file != nil {
		z.file.Close()
	}
}