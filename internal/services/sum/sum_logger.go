package sum

import (
	"context"
	"time"

	"go.uber.org/zap"
)

type ISumCalculator interface {
	SumNumbers(v interface{}) int
}

type SumLogger struct {
	logger *zap.Logger
	next   ISumCalculator
}

func NewSumLogger(next ISumCalculator, logger *zap.Logger) *SumLogger {
	return &SumLogger{next: next, logger: logger.Named("sum-service")}
}

// its always good idea to have context param for every service: for deadlines, for tracing
func (l *SumLogger) SumNumbers(ctx context.Context, v interface{}) int {
	defer func(start time.Time) {
		// todo here we can log/trace request/response and errors
		l.logger.Info("SumNumbers", zap.Duration("took", time.Since(start)))
	}(time.Now())

	return l.next.SumNumbers(v)
}
