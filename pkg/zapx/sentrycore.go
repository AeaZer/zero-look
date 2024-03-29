package zapx

import (
	"time"

	"github.com/getsentry/sentry-go"
	"go.uber.org/zap/zapcore"
)

// 将zap的Level转换为sentry的Level
func sentryLevel(lvl zapcore.Level) sentry.Level {
	switch lvl {
	case zapcore.DebugLevel:
		return sentry.LevelDebug
	case zapcore.InfoLevel:
		return sentry.LevelInfo
	case zapcore.WarnLevel:
		return sentry.LevelWarning
	case zapcore.ErrorLevel:
		return sentry.LevelError
	case zapcore.DPanicLevel:
		return sentry.LevelFatal
	case zapcore.PanicLevel:
		return sentry.LevelFatal
	case zapcore.FatalLevel:
		return sentry.LevelFatal
	default:
		return sentry.LevelFatal
	}
}

// SentryCoreConfig 定义 Sentry Core 的配置参数.
type SentryCoreConfig struct {
	Tags              map[string]string
	DisableStacktrace bool
	Level             zapcore.Level
	FlushTimeout      time.Duration
	Hub               *sentry.Hub
}

// sentryCore sentry core 的 Core 结构体，用于实现 Core 接口
type sentryCore struct {
	client               *sentry.Client    // sentry客户端
	cfg                  *SentryCoreConfig // core配置
	zapcore.LevelEnabler                   // LevelEnabler接口
	flushTimeout         time.Duration     // sentry上报的flush时间

	fields map[string]interface{} // 保存Fields
}

// NewSentryCore 生成Core对象
func NewSentryCore(cfg SentryCoreConfig, sentryClient *sentry.Client) zapcore.Core {
	core := sentryCore{
		client:       sentryClient,
		cfg:          &cfg,
		LevelEnabler: cfg.Level,
		flushTimeout: 3 * time.Second,
		fields:       make(map[string]interface{}),
	}
	if cfg.FlushTimeout > 0 {
		core.flushTimeout = cfg.FlushTimeout
	}
	return &core
}

// With接口方法的实际实现，对传入fields进行设置日志打印时的打印解析方式并添加到已有的fields中
func (c *sentryCore) with(fs []zapcore.Field) *sentryCore {
	// Copy our map.
	m := make(map[string]interface{}, len(c.fields))
	for k, v := range c.fields {
		m[k] = v
	}

	// Add fields to an in-memory encoder.
	enc := zapcore.NewMapObjectEncoder()
	for _, f := range fs {
		f.AddTo(enc)
	}

	// Merge the two maps.
	for k, v := range enc.Fields {
		m[k] = v
	}

	return &sentryCore{
		client:       c.client,
		cfg:          c.cfg,
		fields:       m,
		LevelEnabler: c.LevelEnabler,
	}
}

// With 实现Core接口的With方法
func (c *sentryCore) With(fs []zapcore.Field) zapcore.Core {
	return c.with(fs)
}

// Check 实现Core接口的Check方法，只有大于在core配置中的的Level才会被打印
func (c *sentryCore) Check(ent zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if c.cfg.Level.Enabled(ent.Level) {
		return ce.AddCore(ent, c)
	}
	return ce
}

// Write 实现Core接口的Write方法，对sentry进行上报，Fields作为Extra信息上报
func (c *sentryCore) Write(ent zapcore.Entry, fs []zapcore.Field) error {
	clone := c.with(fs)

	event := sentry.NewEvent()
	event.Message = ent.Message
	event.Timestamp = ent.Time
	event.Level = sentryLevel(ent.Level)
	event.Extra = clone.fields
	event.Tags = c.cfg.Tags

	if !c.cfg.DisableStacktrace {
		trace := sentry.NewStacktrace()
		if trace != nil {
			event.Exception = []sentry.Exception{{
				Type:       ent.Message,
				Value:      ent.Caller.TrimmedPath(),
				Stacktrace: trace,
			}}
		}
	}

	hub := c.cfg.Hub
	if hub == nil {
		hub = sentry.CurrentHub()
	}
	_ = c.client.CaptureEvent(event, nil, hub.Scope())

	if ent.Level > zapcore.ErrorLevel {
		c.client.Flush(c.flushTimeout)
	}
	return nil
}

// Sync 实现Core接口的Sync方法
func (c *sentryCore) Sync() error {
	c.client.Flush(c.flushTimeout)
	return nil
}
