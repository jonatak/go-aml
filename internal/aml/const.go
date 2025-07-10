package aml

import "time"

const (
	MaxAmount            float64       = 10000.0        // maximum amount for a normal activity.
	SuspiciousWindowSize time.Duration = 24 * time.Hour // the window we want to check on.
	MaxBurstGrpcQuery    int           = 10000
)
