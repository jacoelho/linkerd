package linkerd

import (
	"context"
	"encoding/base64"
	"encoding/binary"

	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

const (
	ld5ContextHeaderKey = "l5d-ctx-trace"
)

const (
	flagDebug = 1 << iota
	flagSamplingKnown
	flagSampled
)

type Propagator struct{}

// Asserts that the propagator implements the otel.TextMapPropagator interface at compile time.
var _ propagation.TextMapPropagator = Propagator{}

func New() propagation.TextMapPropagator {
	return Propagator{}
}

// Inject injects a context to the carrier following linkerd trace format.
func (l5d Propagator) Inject(ctx context.Context, carrier propagation.TextMapCarrier) {
	sc := trace.SpanFromContext(ctx).SpanContext()

	if !sc.TraceID().IsValid() || !sc.SpanID().IsValid() {
		return
	}

	spanID := sc.SpanID()
	traceID := sc.TraceID()

	var buf [40]byte
	copy64be(buf[:8], spanID[:])
	// skip parent
	copy64be(buf[16:24], traceID[8:])
	copy64be(buf[32:], traceID[:8])
	if sc.IsSampled() {
		buf[31] = flagSamplingKnown | flagSampled
	}

	carrier.Set(ld5ContextHeaderKey, base64.StdEncoding.EncodeToString(buf[:]))
}

// Extract gets a context from the carrier if it contains linkerd trace header.
func (l5d Propagator) Extract(ctx context.Context, carrier propagation.TextMapCarrier) context.Context {
	h := carrier.Get(ld5ContextHeaderKey)
	if h == "" {
		return ctx
	}

	traceBytes, err := base64.StdEncoding.DecodeString(h)
	if err != nil {
		return ctx
	}

	if len(traceBytes) != 32 && len(traceBytes) != 40 {
		return ctx
	}

	var scc trace.SpanContextConfig
	copy64be(scc.SpanID[:], traceBytes[:8])
	// skip parent not supported
	copy64be(scc.TraceID[8:], traceBytes[16:24])
	if len(traceBytes) == 40 {
		copy64be(scc.TraceID[:8], traceBytes[32:])
	}

	flags := traceBytes[31]
	if flags&flagSamplingKnown != 0 && flags&flagSampled != 0 {
		scc.TraceFlags |= trace.FlagsSampled
	}

	spanCtx := trace.NewSpanContext(scc)
	if !spanCtx.IsValid() {
		return ctx
	}

	return trace.ContextWithRemoteSpanContext(ctx, spanCtx)
}

// Fields returns list of fields used by HTTPTextFormat.
func (l5d Propagator) Fields() []string {
	return []string{ld5ContextHeaderKey}
}

func copy64be(dst, src []byte) {
	binary.BigEndian.PutUint64(dst, binary.BigEndian.Uint64(src))
}
