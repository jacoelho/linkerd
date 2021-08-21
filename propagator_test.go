package linkerd

import (
	"context"
	"encoding/base64"
	"net/http"
	"testing"

	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

func headerCarrierWithValue(value string) propagation.TextMapCarrier {
	h := make(http.Header, 1)
	h.Set(ld5ContextHeaderKey, value)

	return propagation.HeaderCarrier(h)
}

func printBinary(t *testing.T, s string) {
	t.Helper()

	traceBytes, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		t.Error(err)
	}

	t.Logf("%#q\n", traceBytes)
}

func TestPropagator_Extract(t *testing.T) {
	tests := []struct {
		input   string
		traceID string
		spanID  string
		sampled bool
	}{
		// {
		// 	input:   "w7oaZWDKDEgrUYn/oBOtc0EdGALJFR3tAAAAAAAAAAY=",
		// 	traceID: "0000000000000000411d1802c9151ded",
		// 	spanID:  "c3ba1a6560ca0c48",
		// 	sampled: true,
		// },
		{
			input:   "9BQdXcDJNdAAAAAAAAAAADKk2yD11ZLnAAAAAAAAAAYAAAAAAAAAAQ==",
			traceID: "411d1802c9151ded2b5189ffa013ad73",
			spanID:  "c3ba1a6560ca0c48",
			sampled: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.input, func(t *testing.T) {
			t.Parallel()

			ctx := New().Extract(context.Background(), headerCarrierWithValue(tt.input))
			gotSc := trace.SpanContextFromContext(ctx)

			if !gotSc.IsValid() {
				t.Error("expected span context to be valid")
				return
			}

			if got := gotSc.TraceID(); got.String() != tt.traceID {
				t.Errorf("expected TraceID %v, got %v", tt.traceID, got)
			}

			if got := gotSc.SpanID(); got.String() != tt.spanID {
				t.Errorf("expected SpanID %v, got %v", tt.spanID, got)
			}

			if got := gotSc.IsSampled(); got != tt.sampled {
				t.Errorf("expected IsSampled %v, got %v", tt.sampled, got)
			}
		})
	}
}

func TestPropagator_Inject(t *testing.T) {
	tests := []struct {
		want    string
		traceID string
		spanID  string
		sampled trace.TraceFlags
	}{
		{
			want:    "w7oaZWDKDEgrUYn/oBOtc0EdGALJFR3tAAAAAAAAAAY=",
			traceID: "0000000000000000411d1802c9151ded",
			spanID:  "c3ba1a6560ca0c48",
			sampled: trace.FlagsSampled,
		},
		// {
		// 	want:    "9BQdXcDJNdAAAAAAAAAAADKk2yD11ZLnAAAAAAAAAAYAAAAAAAAAAQ==",
		// 	traceID: "411d1802c9151ded2b5189ffa013ad73",
		// 	spanID:  "c3ba1a6560ca0c48",
		// 	sampled: trace.FlagsSampled,
		// },
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.want, func(t *testing.T) {
			t.Parallel()

			traceID, err := trace.TraceIDFromHex(tt.traceID)
			if err != nil {
				t.Errorf("unexpected error %v", err)
				return
			}

			spanID, err := trace.SpanIDFromHex(tt.spanID)
			if err != nil {
				t.Errorf("unexpected error %v", err)
				return
			}

			scc := trace.SpanContextConfig{
				TraceID:    traceID,
				SpanID:     spanID,
				TraceFlags: tt.sampled,
			}

			ctx := trace.ContextWithSpanContext(context.Background(), trace.NewSpanContext(scc))

			hc := propagation.HeaderCarrier(make(http.Header))
			New().Inject(ctx, hc)

			a, _ := base64.StdEncoding.DecodeString(tt.want)
			b, _ := base64.StdEncoding.DecodeString(hc.Get(ld5ContextHeaderKey))

			_ = a
			_ = b

			if got := hc.Get(ld5ContextHeaderKey); got != tt.want {
				t.Errorf("expected header value %v, got %v", tt.want, got)
			}
		})
	}
}
