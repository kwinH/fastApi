package util

import (
	"context"
	"fmt"
	"github.com/spf13/viper"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	oteltrace "go.opentelemetry.io/otel/trace"
)

func TracerProvider() (*tracesdk.TracerProvider, error) {
	//endpoint := "http://127.0.0.1:14268/api/traces"
	// Create the Jaeger exporter
	endpoint := viper.GetString("telemetry.endpoint")
	serviceName := viper.GetString("telemetry.name")
	sampler := viper.GetFloat64("telemetry.sampler")

	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(endpoint)))
	if err != nil {
		return nil, err
	}
	tp := tracesdk.NewTracerProvider(
		// Always be sure to batch in production.
		tracesdk.WithSampler(tracesdk.AlwaysSample()),
		tracesdk.WithBatcher(exp),
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(serviceName),
		),
		),
		tracesdk.WithSampler(tracesdk.ParentBased(tracesdk.TraceIDRatioBased(sampler))),
	)

	// Register our TracerProvider as the global so any imported
	// instrumentation in the future will default to using it.
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	return tp, nil
}

func ContextWithSpanContext(c context.Context, traceId, spanId, tracerName, spanName string) (ctx context.Context, span oteltrace.Span, err error) {
	ctx = c
	if !viper.IsSet("telemetry") {
		return
	}

	if spanId == "" {
		tracer := otel.Tracer(tracerName)
		ctx, span = tracer.Start(
			c,
			spanName,
		)
		return
	}

	// 从消息中恢复 TraceID 和 SpanID
	traceID, err := oteltrace.TraceIDFromHex(traceId)
	if err != nil {
		err = fmt.Errorf("invalid TraceID: %v", err)
		return
	}

	//spanID, err := oteltrace.SpanIDFromHex(spanId)
	//if err != nil {
	//	err = fmt.Errorf("invalid SpanID: %v", err)
	//	return
	//}

	// 构造 SpanContext
	spanContext := oteltrace.NewSpanContext(oteltrace.SpanContextConfig{
		TraceID: traceID,
		//		SpanID:     spanID,
		TraceFlags: oteltrace.FlagsSampled, // 标记为被采样的 Trace
	})

	// 使用原始 TraceContext 恢复上下文
	ctx = oteltrace.ContextWithSpanContext(c, spanContext)

	// 创建新的 Span 处理消息
	tracer := otel.Tracer(tracerName)
	ctx, span = tracer.Start(ctx, spanName)
	return
}
