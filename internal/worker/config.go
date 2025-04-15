package worker

import "time"

const (
	defaultTickerPeriod = time.Duration(10) * time.Second
	defaultTimeout      = time.Duration(5) * time.Second
	defaultFanOutPool   = 10
)

type SyncWorkerConfig struct {
	tickerPeriod time.Duration
	timeout      time.Duration
	fanOutPool   int
}

type SyncWorkerConfigBuilder struct {
	cfg *SyncWorkerConfig
}

func NewSyncWorkerConfigBuilder() *SyncWorkerConfigBuilder {
	return &SyncWorkerConfigBuilder{
		cfg: &SyncWorkerConfig{
			tickerPeriod: defaultTickerPeriod,
			timeout:      defaultTimeout,
			fanOutPool:   defaultFanOutPool,
		},
	}
}

func (b *SyncWorkerConfigBuilder) WithTimeout(timeout time.Duration) *SyncWorkerConfigBuilder {
	b.cfg.timeout = timeout
	return b
}

func (b *SyncWorkerConfigBuilder) Build() *SyncWorkerConfig {
	return b.cfg
}
