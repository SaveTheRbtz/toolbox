package diag

import (
	"github.com/watermint/toolbox/essentials/log/esl"
	"github.com/watermint/toolbox/essentials/model/mo_filter"
	"github.com/watermint/toolbox/essentials/model/mo_string"
	"github.com/watermint/toolbox/essentials/network/nw_capture"
	"github.com/watermint/toolbox/infra/control/app_control"
	"github.com/watermint/toolbox/infra/control/app_job"
	"github.com/watermint/toolbox/infra/recipe/rc_exec"
	"github.com/watermint/toolbox/infra/recipe/rc_recipe"
	"github.com/watermint/toolbox/infra/report/rp_model"
	"math"
	"time"
)

type ThroughputRow struct {
	Time               string `json:"time"`
	Concurrency        int    `json:"concurrency"`
	SuccessConcurrency int    `json:"success_concurrency"`
	SuccessSent        int64  `json:"success_sent"`
	SuccessReceived    int64  `json:"success_received"`
	FailureConcurrency int    `json:"failure_concurrency"`
	FailureSent        int64  `json:"failure_sent"`
	FailureReceived    int64  `json:"failure_received"`
}

type Throughput struct {
	rc_recipe.RemarkSecret
	rc_recipe.RemarkConsole
	JobId      mo_string.OptionalString
	Path       mo_string.OptionalString
	Endpoint   mo_filter.Filter
	Bucket     int // bucket time in milli seconds
	Report     rp_model.RowReport
	TimeFormat string

	buckets map[time.Time]*ThroughputRow
}

func (z *Throughput) Preset() {
	z.Report.SetModel(&ThroughputRow{})
	z.Endpoint.SetOptions(
		mo_filter.NewNameFilter(),
		mo_filter.NewNamePrefixFilter(),
		mo_filter.NewNameSuffixFilter(),
	)
	z.Bucket = 1000
	z.TimeFormat = "2006-01-02 15:04:05.999"
}

func (z *Throughput) handleRecord(history app_job.History, rec nw_capture.Record) {
	l := esl.Default()

	if !z.Endpoint.Accept(rec.Req.RequestUrl) {
		l.Debug("Skip", esl.String("endpoint", rec.Req.RequestUrl))
		return
	}

	t, err := time.Parse("2006-01-02T15:04:05.999Z0700", rec.Time)
	// skip the record with invalid time format
	if err != nil {
		l.Debug("Unable to parse time", esl.Error(err), esl.String("time", rec.Time))
		return
	}

	var sent, recv int64

	bt := t.Truncate(time.Duration(z.Bucket) * time.Millisecond)
	latency := rec.Latency / 1_000_000 // ns -> ms
	bl := int64(math.Ceil(float64(latency) / float64(z.Bucket)))
	be := bt.Add(time.Duration(rec.Latency) * time.Nanosecond)
	sent = rec.Req.ContentLength
	if rec.Res != nil {
		recv = rec.Res.ContentLength
	}
	if bl < 1 {
		bl = 1
	}
	bs := sent / bl
	br := recv / bl

	for b := bt; b.Equal(be) || b.Before(be); b = b.Add(time.Duration(z.Bucket) * time.Millisecond) {
		if bucket, ok := z.buckets[b]; ok {
			if rec.IsSuccess() {
				bucket.Concurrency++
				bucket.SuccessConcurrency++
				bucket.SuccessSent += bs
				bucket.SuccessReceived += br
			} else {
				bucket.Concurrency++
				bucket.FailureConcurrency++
				bucket.FailureSent += bs
				bucket.FailureReceived += br
			}
		} else {
			if rec.IsSuccess() {
				z.buckets[b] = &ThroughputRow{
					Time:               b.Format(z.TimeFormat),
					Concurrency:        1,
					SuccessSent:        bs,
					SuccessReceived:    br,
					SuccessConcurrency: 1,
					FailureSent:        0,
					FailureReceived:    0,
					FailureConcurrency: 0,
				}
			} else {
				z.buckets[b] = &ThroughputRow{
					Time:               b.Format(z.TimeFormat),
					Concurrency:        1,
					SuccessSent:        0,
					SuccessReceived:    0,
					SuccessConcurrency: 0,
					FailureSent:        bs,
					FailureReceived:    br,
					FailureConcurrency: 1,
				}
			}
		}
	}
}

func (z *Throughput) Exec(c app_control.Control) error {
	l := c.Log()
	if err := z.Report.Open(); err != nil {
		return err
	}

	z.buckets = make(map[time.Time]*ThroughputRow)

	loader := CaptureLoader{
		Ctl:   c,
		JobId: z.JobId,
		Path:  z.Path,
	}
	err := loader.Load(z.handleRecord)
	if err != nil {
		l.Debug("Unable to load", esl.Error(err))
		return err
	}

	for _, bucket := range z.buckets {
		z.Report.Row(bucket)
	}

	return nil
}

func (z *Throughput) Test(c app_control.Control) error {
	return rc_exec.Exec(c, &Throughput{}, rc_recipe.NoCustomValues)
}
