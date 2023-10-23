// Code generated by mockery v2.35.4. DO NOT EDIT.

package mocks

import (
	config "github.com/smartcontractkit/chainlink/v2/core/chains/evm/config"
	coreconfig "github.com/smartcontractkit/chainlink/v2/core/config"

	mock "github.com/stretchr/testify/mock"

	time "time"

	uuid "github.com/google/uuid"

	zapcore "go.uber.org/zap/zapcore"
)

// ChainScopedConfig is an autogenerated mock type for the ChainScopedConfig type
type ChainScopedConfig struct {
	mock.Mock
}

// AppID provides a mock function with given fields:
func (_m *ChainScopedConfig) AppID() uuid.UUID {
	ret := _m.Called()

	var r0 uuid.UUID
	if rf, ok := ret.Get(0).(func() uuid.UUID); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(uuid.UUID)
		}
	}

	return r0
}

// AuditLogger provides a mock function with given fields:
func (_m *ChainScopedConfig) AuditLogger() coreconfig.AuditLogger {
	ret := _m.Called()

	var r0 coreconfig.AuditLogger
	if rf, ok := ret.Get(0).(func() coreconfig.AuditLogger); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(coreconfig.AuditLogger)
		}
	}

	return r0
}

// AutoPprof provides a mock function with given fields:
func (_m *ChainScopedConfig) AutoPprof() coreconfig.AutoPprof {
	ret := _m.Called()

	var r0 coreconfig.AutoPprof
	if rf, ok := ret.Get(0).(func() coreconfig.AutoPprof); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(coreconfig.AutoPprof)
		}
	}

	return r0
}

// CosmosEnabled provides a mock function with given fields:
func (_m *ChainScopedConfig) CosmosEnabled() bool {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// Database provides a mock function with given fields:
func (_m *ChainScopedConfig) Database() coreconfig.Database {
	ret := _m.Called()

	var r0 coreconfig.Database
	if rf, ok := ret.Get(0).(func() coreconfig.Database); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(coreconfig.Database)
		}
	}

	return r0
}

// EVM provides a mock function with given fields:
func (_m *ChainScopedConfig) EVM() config.EVM {
	ret := _m.Called()

	var r0 config.EVM
	if rf, ok := ret.Get(0).(func() config.EVM); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(config.EVM)
		}
	}

	return r0
}

// EVMEnabled provides a mock function with given fields:
func (_m *ChainScopedConfig) EVMEnabled() bool {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// EVMRPCEnabled provides a mock function with given fields:
func (_m *ChainScopedConfig) EVMRPCEnabled() bool {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// Feature provides a mock function with given fields:
func (_m *ChainScopedConfig) Feature() coreconfig.Feature {
	ret := _m.Called()

	var r0 coreconfig.Feature
	if rf, ok := ret.Get(0).(func() coreconfig.Feature); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(coreconfig.Feature)
		}
	}

	return r0
}

// FluxMonitor provides a mock function with given fields:
func (_m *ChainScopedConfig) FluxMonitor() coreconfig.FluxMonitor {
	ret := _m.Called()

	var r0 coreconfig.FluxMonitor
	if rf, ok := ret.Get(0).(func() coreconfig.FluxMonitor); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(coreconfig.FluxMonitor)
		}
	}

	return r0
}

// Insecure provides a mock function with given fields:
func (_m *ChainScopedConfig) Insecure() coreconfig.Insecure {
	ret := _m.Called()

	var r0 coreconfig.Insecure
	if rf, ok := ret.Get(0).(func() coreconfig.Insecure); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(coreconfig.Insecure)
		}
	}

	return r0
}

// InsecureFastScrypt provides a mock function with given fields:
func (_m *ChainScopedConfig) InsecureFastScrypt() bool {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// JobPipeline provides a mock function with given fields:
func (_m *ChainScopedConfig) JobPipeline() coreconfig.JobPipeline {
	ret := _m.Called()

	var r0 coreconfig.JobPipeline
	if rf, ok := ret.Get(0).(func() coreconfig.JobPipeline); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(coreconfig.JobPipeline)
		}
	}

	return r0
}

// Keeper provides a mock function with given fields:
func (_m *ChainScopedConfig) Keeper() coreconfig.Keeper {
	ret := _m.Called()

	var r0 coreconfig.Keeper
	if rf, ok := ret.Get(0).(func() coreconfig.Keeper); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(coreconfig.Keeper)
		}
	}

	return r0
}

// Log provides a mock function with given fields:
func (_m *ChainScopedConfig) Log() coreconfig.Log {
	ret := _m.Called()

	var r0 coreconfig.Log
	if rf, ok := ret.Get(0).(func() coreconfig.Log); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(coreconfig.Log)
		}
	}

	return r0
}

// LogConfiguration provides a mock function with given fields: log
func (_m *ChainScopedConfig) LogConfiguration(log coreconfig.LogfFn) {
	_m.Called(log)
}

// Mercury provides a mock function with given fields:
func (_m *ChainScopedConfig) Mercury() coreconfig.Mercury {
	ret := _m.Called()

	var r0 coreconfig.Mercury
	if rf, ok := ret.Get(0).(func() coreconfig.Mercury); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(coreconfig.Mercury)
		}
	}

	return r0
}

// OCR provides a mock function with given fields:
func (_m *ChainScopedConfig) OCR() coreconfig.OCR {
	ret := _m.Called()

	var r0 coreconfig.OCR
	if rf, ok := ret.Get(0).(func() coreconfig.OCR); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(coreconfig.OCR)
		}
	}

	return r0
}

// OCR2 provides a mock function with given fields:
func (_m *ChainScopedConfig) OCR2() coreconfig.OCR2 {
	ret := _m.Called()

	var r0 coreconfig.OCR2
	if rf, ok := ret.Get(0).(func() coreconfig.OCR2); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(coreconfig.OCR2)
		}
	}

	return r0
}

// P2P provides a mock function with given fields:
func (_m *ChainScopedConfig) P2P() coreconfig.P2P {
	ret := _m.Called()

	var r0 coreconfig.P2P
	if rf, ok := ret.Get(0).(func() coreconfig.P2P); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(coreconfig.P2P)
		}
	}

	return r0
}

// Password provides a mock function with given fields:
func (_m *ChainScopedConfig) Password() coreconfig.Password {
	ret := _m.Called()

	var r0 coreconfig.Password
	if rf, ok := ret.Get(0).(func() coreconfig.Password); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(coreconfig.Password)
		}
	}

	return r0
}

// Prometheus provides a mock function with given fields:
func (_m *ChainScopedConfig) Prometheus() coreconfig.Prometheus {
	ret := _m.Called()

	var r0 coreconfig.Prometheus
	if rf, ok := ret.Get(0).(func() coreconfig.Prometheus); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(coreconfig.Prometheus)
		}
	}

	return r0
}

// Pyroscope provides a mock function with given fields:
func (_m *ChainScopedConfig) Pyroscope() coreconfig.Pyroscope {
	ret := _m.Called()

	var r0 coreconfig.Pyroscope
	if rf, ok := ret.Get(0).(func() coreconfig.Pyroscope); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(coreconfig.Pyroscope)
		}
	}

	return r0
}

// RootDir provides a mock function with given fields:
func (_m *ChainScopedConfig) RootDir() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// Sentry provides a mock function with given fields:
func (_m *ChainScopedConfig) Sentry() coreconfig.Sentry {
	ret := _m.Called()

	var r0 coreconfig.Sentry
	if rf, ok := ret.Get(0).(func() coreconfig.Sentry); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(coreconfig.Sentry)
		}
	}

	return r0
}

// SetLogLevel provides a mock function with given fields: lvl
func (_m *ChainScopedConfig) SetLogLevel(lvl zapcore.Level) error {
	ret := _m.Called(lvl)

	var r0 error
	if rf, ok := ret.Get(0).(func(zapcore.Level) error); ok {
		r0 = rf(lvl)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SetLogSQL provides a mock function with given fields: logSQL
func (_m *ChainScopedConfig) SetLogSQL(logSQL bool) {
	_m.Called(logSQL)
}

// SetPasswords provides a mock function with given fields: keystore, vrf
func (_m *ChainScopedConfig) SetPasswords(keystore *string, vrf *string) {
	_m.Called(keystore, vrf)
}

// ShutdownGracePeriod provides a mock function with given fields:
func (_m *ChainScopedConfig) ShutdownGracePeriod() time.Duration {
	ret := _m.Called()

	var r0 time.Duration
	if rf, ok := ret.Get(0).(func() time.Duration); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(time.Duration)
	}

	return r0
}

// SolanaEnabled provides a mock function with given fields:
func (_m *ChainScopedConfig) SolanaEnabled() bool {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// StarkNetEnabled provides a mock function with given fields:
func (_m *ChainScopedConfig) StarkNetEnabled() bool {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// TelemetryIngress provides a mock function with given fields:
func (_m *ChainScopedConfig) TelemetryIngress() coreconfig.TelemetryIngress {
	ret := _m.Called()

	var r0 coreconfig.TelemetryIngress
	if rf, ok := ret.Get(0).(func() coreconfig.TelemetryIngress); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(coreconfig.TelemetryIngress)
		}
	}

	return r0
}

// Threshold provides a mock function with given fields:
func (_m *ChainScopedConfig) Threshold() coreconfig.Threshold {
	ret := _m.Called()

	var r0 coreconfig.Threshold
	if rf, ok := ret.Get(0).(func() coreconfig.Threshold); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(coreconfig.Threshold)
		}
	}

	return r0
}

// Tracing provides a mock function with given fields:
func (_m *ChainScopedConfig) Tracing() coreconfig.Tracing {
	ret := _m.Called()

	var r0 coreconfig.Tracing
	if rf, ok := ret.Get(0).(func() coreconfig.Tracing); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(coreconfig.Tracing)
		}
	}

	return r0
}

// Validate provides a mock function with given fields:
func (_m *ChainScopedConfig) Validate() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ValidateDB provides a mock function with given fields:
func (_m *ChainScopedConfig) ValidateDB() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// WebServer provides a mock function with given fields:
func (_m *ChainScopedConfig) WebServer() coreconfig.WebServer {
	ret := _m.Called()

	var r0 coreconfig.WebServer
	if rf, ok := ret.Get(0).(func() coreconfig.WebServer); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(coreconfig.WebServer)
		}
	}

	return r0
}

// NewChainScopedConfig creates a new instance of ChainScopedConfig. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewChainScopedConfig(t interface {
	mock.TestingT
	Cleanup(func())
}) *ChainScopedConfig {
	mock := &ChainScopedConfig{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
