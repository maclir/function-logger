package logger

import (
	"context"
	"fmt"
	"os"

	"cloud.google.com/go/logging"
	"contrib.go.opencensus.io/exporter/stackdriver"
	"go.opencensus.io/trace"
	"google.golang.org/genproto/googleapis/api/monitoredres"
)

// New creates wraps a new instance of logging.Logger with
// function attributes.
func New() (*logging.Logger, error) {
	projectID := os.Getenv("GCP_PROJECT")
	if projectID == "" {
		return nil, fmt.Errorf("GCP_PROJECT environment variable unset or missing")
	}

	functionName := os.Getenv("FUNCTION_NAME")
	if functionName == "" {
		return nil, fmt.Errorf("FUNCTION_NAME environment variable unset or missing")
	}

	region := os.Getenv("FUNCTION_REGION")
	if region == "" {
		return nil, fmt.Errorf("FUNCTION_REGION environment variable unset or missing")
	}

	stackdriverExporter, err := stackdriver.NewExporter(stackdriver.Options{ProjectID: projectID})
	if err != nil {
		return nil, err
	}

	trace.RegisterExporter(stackdriverExporter)
	trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})

	client, err := logging.NewClient(context.Background(), projectID)
	if err != nil {
		return nil, err
	}

	monitoredResource := monitoredres.MonitoredResource{
		Type: "cloud_function",
		Labels: map[string]string{
			"function_name": functionName,
			"region":        region,
		},
	}

	commonResource := logging.CommonResource(&monitoredResource)
	logger := client.Logger(functionName, commonResource)

	return logger, nil
}
