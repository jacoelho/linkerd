# linkerd

[OpenTelemetry](https://opentelemetry.io/) linkerd trace propagator.

Currently supports `l5d-tctx-trace` set by `telemetry: io.l5d.tracelog`.

## Usage

```go
import (
	"net/http"

	"github.com/jacoelho/linkerd"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

var httpClient = http.Client{Transport: otelhttp.NewTransport(
    http.DefaultTransport,
    // ...
    otelhttp.WithPropagators(linkerd.New()),
),
}
```

## License

[license](./LICENSE)
