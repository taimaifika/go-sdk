package otel

import (
	"context"
	"errors"
	"flag"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"

	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutlog"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"

	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

// Default values for configuration.
const (
	otelProtocolHTTP = "http"
	otelProtocolGRPC = "grpc"

	defaultOtelEndpoint = ""
	defaultNameService  = ""
	defaultVersion      = ""
	defaultOtelProtocol = otelProtocolGRPC
	defaultIsEnabled    = true
)

// Config
type Config struct {
}

// OtelPlugin
type OtelPlugin struct {
	Config
	name      string
	prefix    string
	isEnabled bool
	ctx       context.Context

	// otel attributes
	serviceName    string
	serviceVersion string

	// otel exporter
	exporterOtlpEndpoint string
	exporterOtlpProtocol string

	// otel features
	isEnabledTrace  bool
	isEnabledMetric bool
	isEnabledLog    bool

	shutdown func(context.Context) error
}

// New creates a new OtelPlugin.
func NewOtelPlugin(name string) *OtelPlugin {
	return &OtelPlugin{
		name:           name,
		prefix:         name,
		ctx:            context.Background(),
		serviceName:    defaultNameService,
		serviceVersion: defaultNameService,
		isEnabled:      defaultIsEnabled,
	}
}

// IsEnabled returns the value of isEnabled.
func (op *OtelPlugin) IsEnabled() bool {
	return op.isEnabled
}

// setupOTelSDK bootstraps the OpenTelemetry pipeline.
// If it does not return an error, make sure to call shutdown for proper cleanup.
func (op *OtelPlugin) SetupOTelSDK() (shutdown func(context.Context) error, err error) {
	var shutdownFuncs []func(context.Context) error

	// shutdown calls cleanup functions registered via shutdownFuncs.
	// The errors from the calls are joined.
	// Each registered cleanup will be invoked once.
	shutdown = func(ctx context.Context) error {
		var err error
		for _, fn := range shutdownFuncs {
			err = errors.Join(err, fn(ctx))
		}
		shutdownFuncs = nil
		return err
	}

	// handleErr calls shutdown for cleanup and makes sure that all errors are returned.
	handleErr := func(inErr error) {
		err = errors.Join(inErr, shutdown(op.ctx))
	}

	// Set up propagator.
	prop := newPropagator()
	otel.SetTextMapPropagator(prop)

	// Set up trace provider.
	if op.isEnabledTrace {
		tracerProvider, providerErr := op.newTraceProvider()
		if providerErr != nil {
			handleErr(providerErr)
			return
		}
		shutdownFuncs = append(shutdownFuncs, tracerProvider.Shutdown)
		otel.SetTracerProvider(tracerProvider)
	}

	// Set up meter provider.
	if op.isEnabledMetric {
		meterProvider, providerErr := op.newMeterProvider()
		if providerErr != nil {
			handleErr(providerErr)
			return
		}
		shutdownFuncs = append(shutdownFuncs, meterProvider.Shutdown)
		otel.SetMeterProvider(meterProvider)
	}

	// Set up logger provider.
	if op.isEnabledLog {
		loggerProvider, providerErr := op.newLoggerProvider()
		if providerErr != nil {
			handleErr(providerErr)
			return
		}
		shutdownFuncs = append(shutdownFuncs, loggerProvider.Shutdown)
		global.SetLoggerProvider(loggerProvider)
	}
	return
}

// newPropagator creates a new propagator.
func newPropagator() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
}

// newTraceProvider creates a new trace provider.
func (op *OtelPlugin) newTraceProvider() (*trace.TracerProvider, error) {
	var traceExporter trace.SpanExporter

	if op.isOtlpProtocolEnabled() {
		// Exporter to otlp
		otlpTraceExporter, err := op.newOtlpTraceExporter()
		if err != nil {
			return nil, err
		}
		traceExporter = otlpTraceExporter
	} else {
		// Exporter to stdout
		stdoutTraceExporter, err := stdouttrace.New(
			stdouttrace.WithPrettyPrint(),
		)
		if err != nil {
			return nil, err
		}
		traceExporter = stdoutTraceExporter
	}

	// Resource attributes
	res := op.newResource()

	traceProvider := trace.NewTracerProvider(
		trace.WithBatcher(traceExporter,
			// Default is 5s. Set to 1s for demonstrative purposes.
			trace.WithBatchTimeout(time.Second)),
		trace.WithResource(res),
	)
	return traceProvider, nil
}

// newOtlpTraceExporter creates a new OTLP trace exporter. (gRPC or HTTP)
func (op *OtelPlugin) newOtlpTraceExporter() (trace.SpanExporter, error) {
	if op.exporterOtlpProtocol == otelProtocolHTTP {
		return otlptracehttp.New(op.ctx)
	}
	return otlptracegrpc.New(op.ctx)
}

// newResource creates a new resource with service.name and service.namespace.
func (op *OtelPlugin) newResource() *resource.Resource {
	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(op.serviceName),
		semconv.ServiceVersionKey.String(op.serviceVersion),
	)
	return res
}

// newMeterProvider creates a new meter provider.
func (op *OtelPlugin) newMeterProvider() (*metric.MeterProvider, error) {
	var metricExporter metric.Exporter
	if op.isOtlpProtocolEnabled() {
		// Exporter to otlp
		otlpMetricExporter, err := op.newOtlpMetricExporter()
		if err != nil {
			return nil, err
		}
		metricExporter = otlpMetricExporter
	} else {
		// Exporter to stdout
		stdoutMetricExporter, err := stdoutmetric.New()
		if err != nil {
			return nil, err
		}
		metricExporter = stdoutMetricExporter
	}

	meterProvider := metric.NewMeterProvider(
		metric.WithReader(metric.NewPeriodicReader(metricExporter,
			// Default is 1m. Set to 3s for demonstrative purposes.
			metric.WithInterval(3*time.Second))),
	)
	return meterProvider, nil
}

// newOtlpMetricExporter creates a new OTLP metric exporter. (gRPC or HTTP)
func (op *OtelPlugin) newOtlpMetricExporter() (metric.Exporter, error) {
	if op.exporterOtlpProtocol == otelProtocolHTTP {
		return otlpmetrichttp.New(op.ctx)
	}
	return otlpmetricgrpc.New(op.ctx)
}

// newLoggerProvider creates a new logger provider.
func (op *OtelPlugin) newLoggerProvider() (*log.LoggerProvider, error) {
	var logExporter log.Exporter
	if op.isOtlpProtocolEnabled() {
		// Exporter to otlp
		otlpLogExporter, err := op.newOtlpLogExporter()
		if err != nil {
			return nil, err
		}
		logExporter = otlpLogExporter
	} else {
		// Exporter to stdout
		stdoutLogExporter, err := stdoutlog.New()
		if err != nil {
			return nil, err
		}
		logExporter = stdoutLogExporter
	}

	loggerProvider := log.NewLoggerProvider(
		log.WithProcessor(log.NewBatchProcessor(logExporter)),
	)
	return loggerProvider, nil
}

// newOtlpLogExporter creates a new OTLP log exporter. (gRPC or HTTP)
func (op *OtelPlugin) newOtlpLogExporter() (log.Exporter, error) {
	if op.exporterOtlpProtocol == otelProtocolHTTP {
		return otlploghttp.New(op.ctx)
	}
	return otlploggrpc.New(op.ctx)
}

// IsOtlpProtocolEnabled returns true if the otlp protocol is enabled.
func (op *OtelPlugin) isOtlpProtocolEnabled() bool {
	return op.exporterOtlpEndpoint != defaultOtelEndpoint
}

// Implement PrefixRunnable interface
// Get returns the service.
func (op *OtelPlugin) Get() interface{} {
	return op
}

// Prefix returns the prefix of the service.
func (op *OtelPlugin) Prefix() string {
	return op.prefix
}

// Name returns the name of the service.
func (op *OtelPlugin) Name() string {
	return op.name
}

// Configure configures the service.
func (op *OtelPlugin) Configure() error {
	// Check if the servicename is empty
	if op.serviceName == "" {
		return errors.New("otel service name is empty")
	}

	// Check if the serviceversion is empty
	if op.serviceVersion == "" {
		return errors.New("otel service version is empty")
	}

	// Check if the exporterOtlpEndpoint is empty
	if op.exporterOtlpEndpoint == "" {
		return errors.New("if OTEL_IS_ENABLED=true, then otel exporter otlp endpoint is not empty, e.g. http://localhost:4317")
	}

	return nil
}

// Run runs the service.
func (op *OtelPlugin) Run() (err error) {
	if !op.isEnabled {
		return nil
	}

	// Configure the service
	if err := op.Configure(); err != nil {
		return err
	}

	// Setup OpenTelemetry SDK
	op.shutdown, err = op.SetupOTelSDK()

	if err != nil {
		return err
	}

	return nil
}

// Stop stops the service.
func (op *OtelPlugin) Stop() <-chan bool {
	c := make(chan bool)
	go func() {
		c <- true
		op.shutdown(op.ctx)
	}()
	return c
}

// GetPrefix returns the prefix of the service.
func (op *OtelPlugin) GetPrefix() string {
	return op.prefix
}

// InitFlags initializes the flags.
func (op *OtelPlugin) InitFlags() {
	flag.BoolVar(&op.isEnabled, op.prefix+"-is-enabled", defaultIsEnabled, "Enable otel service")

	// otel attributes
	// OTEL_SERVICE_NAME
	flag.StringVar(&op.serviceName, op.prefix+"-service-name", defaultNameService, "Service name")
	// OTEL_SERVICE_VERSION
	flag.StringVar(&op.serviceVersion, op.prefix+"-service-version", defaultVersion, "Service version, e.g. 1.0.0")

	// otel exporter
	// OTEL_EXPORTER_OTLP_ENDPOINT
	flag.StringVar(&op.exporterOtlpEndpoint, op.prefix+"-exporter-otlp-endpoint", defaultOtelEndpoint, "Otel otlp endpoint, e.g. http://localhost:4317")
	// OTEL_EXPORTER_OTLP_PROTOCOL
	flag.StringVar(&op.exporterOtlpProtocol, op.prefix+"-exporter-otlp-protocol", defaultOtelProtocol, "Otel protocol, e.g. http or grpc")

	// otel features
	flag.BoolVar(&op.isEnabledTrace, op.prefix+"-is-enabled-trace", true, "Enable otel trace")
	flag.BoolVar(&op.isEnabledMetric, op.prefix+"-is-enabled-metric", true, "Enable otel metric")
	flag.BoolVar(&op.isEnabledLog, op.prefix+"-is-enabled-log", true, "Enable otel log")
}
