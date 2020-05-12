package nw_retry

import (
	"github.com/watermint/toolbox/essentials/go/es_goroutine"
	"github.com/watermint/toolbox/essentials/http/es_response"
	"github.com/watermint/toolbox/essentials/log/esl"
	"github.com/watermint/toolbox/essentials/network/nw_client"
	"github.com/watermint/toolbox/essentials/network/nw_ratelimit"
	"github.com/watermint/toolbox/infra/api/api_context"
	"github.com/watermint/toolbox/quality/infra/qt_errors"
	"time"
)

func NewRetry(client nw_client.Rest) nw_client.Rest {
	return &Retry{
		client: client,
	}
}

type Retry struct {
	client nw_client.Rest
}

func (z *Retry) Call(ctx api_context.Context, req nw_client.RequestBuilder) (res es_response.Response) {
	l := ctx.Log().With(
		esl.String("Url", req.Endpoint()),
		esl.String("Routine", es_goroutine.GetGoRoutineName()),
	)

	res = z.client.Call(ctx, req)

	// return when on success
	if res.IsSuccess() {
		return res
	}
	if res.TransportError() == nil {
		return res
	}

	switch er := res.TransportError().(type) {
	case *ErrorRateLimit:
		l.Debug("Rate limit, waiting for reset",
			esl.Int("limit", er.Limit),
			esl.Int("remaining", er.Remaining),
			esl.String("reset", er.Reset.Format(time.RFC3339)))
		nw_ratelimit.UpdateRetryAfter(ctx.ClientHash(), req.Endpoint(), er.Reset)
		return z.Call(ctx, req)

	default:
		if re, cont := qt_errors.ErrorsForTest(ctx.Log(), er); cont || re == nil {
			return res
		}
		abort := nw_ratelimit.AddError(ctx.ClientHash(), req.Endpoint(), er)
		if abort {
			l.Debug("Abort retry due to retries exceeds retry limit")
			return res
		}
		l.Debug("Retrying", esl.Error(er))
		return z.Call(ctx, req)
	}
}