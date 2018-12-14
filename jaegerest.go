package jaegerest


import (
  "fmt"
  "io"
  "log"
  "net/http"

  "github.com/bpmason1/gaia"
  config "github.com/uber/jaeger-client-go/config"
  jaeger "github.com/uber/jaeger-client-go"
  opentracing "github.com/opentracing/opentracing-go"
  "github.com/opentracing/opentracing-go/ext"
)

//------------------------------------------------------------------------------
//------------------------------------------------------------------------------
func InitJaeger(serviceName string) (opentracing.Tracer, io.Closer) {
    agentHost := gaia.GetEnvWithDefault("JAEGER_AGENT_HOST", "127.0.0.1")
    agentPort := gaia.GetPortWithDefault("JAEGER_AGENT_PORT", 6831)
    agentAddr := fmt.Sprintf("%s:%d", agentHost, agentPort)

    cfg := &config.Configuration{
        Sampler: &config.SamplerConfig{
            Type:  "const",
            Param: 1,
        },
        Reporter: &config.ReporterConfig{
            LogSpans: true,
            LocalAgentHostPort: agentAddr,
        },
    }
    tracer, closer, err := cfg.New(serviceName, config.Logger(jaeger.StdLogger))
    if err != nil {
        log.Printf("ERROR: cannot init Jaeger: %v\n", err)
    }
    return tracer, closer
}

//------------------------------------------------------------------------------
//------------------------------------------------------------------------------
func TraceRequest( fn func(http.ResponseWriter, *http.Request) ) func(http.ResponseWriter, *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        spanCtx, _ := opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
        span := opentracing.GlobalTracer().StartSpan(r.URL.String(), ext.RPCServerOption(spanCtx))
        defer span.Finish()
        fn(w, r)
    }
}

