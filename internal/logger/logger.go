package logger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"sync"
)

// 模块常量
const (
	ModuleMain       = "main"
	ModuleEncryption = "encryption"
	ModuleScheduler  = "scheduler"
	ModuleBitwarden  = "bitwarden"
	ModuleDatabase   = "database"
)

var defaultLogger *slog.Logger

// CustomTextHandler 自定义文本处理器
// 输出格式: time="..." level=INFO module=xxx msg="..." key=value
type CustomTextHandler struct {
	out    io.Writer
	level  slog.Level
	module string
	attrs  []slog.Attr
	mu     *sync.Mutex
}

func (h *CustomTextHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.level
}

func (h *CustomTextHandler) Handle(_ context.Context, r slog.Record) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	timeStr := r.Time.Format("2006/01/02 15:04:05.000")
	levelStr := r.Level.String()

	// 构建基础输出: 时间 级别 [模块] 消息
	var output string
	if h.module != "" {
		output = fmt.Sprintf("%s %s [%s] %s", timeStr, levelStr, h.module, r.Message)
	} else {
		output = fmt.Sprintf("%s %s %s", timeStr, levelStr, r.Message)
	}

	// 追加预设属性
	for _, attr := range h.attrs {
		output += " " + formatAttr(attr)
	}

	// 追加日志调用时的属性
	r.Attrs(func(a slog.Attr) bool {
		output += " " + formatAttr(a)
		return true
	})

	_, err := fmt.Fprintln(h.out, output)
	return err
}

func (h *CustomTextHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newAttrs := make([]slog.Attr, len(h.attrs), len(h.attrs)+len(attrs))
	copy(newAttrs, h.attrs)
	newAttrs = append(newAttrs, attrs...)
	return &CustomTextHandler{out: h.out, level: h.level, module: h.module, attrs: newAttrs, mu: h.mu}
}

func (h *CustomTextHandler) WithGroup(_ string) slog.Handler {
	return h
}

// formatAttr 格式化属性为 key=value 形式
func formatAttr(a slog.Attr) string {
	if a.Value.Kind() == slog.KindString {
		return fmt.Sprintf("%s=%q", a.Key, a.Value.String())
	}
	return fmt.Sprintf("%s=%v", a.Key, a.Value.Any())
}

// 全局互斥锁，用于所有 CustomTextHandler 实例
var handlerMu sync.Mutex

// Init 初始化全局 logger
func Init(level slog.Level) {
	handler := &CustomTextHandler{
		out:   os.Stdout,
		level: level,
		mu:    &handlerMu,
	}
	defaultLogger = slog.New(handler)
	slog.SetDefault(defaultLogger)
}

// Get 获取全局 logger
func Get() *slog.Logger {
	if defaultLogger == nil {
		Init(slog.LevelInfo)
	}
	return defaultLogger
}

// WithContext 从 context 获取 logger
func WithContext(ctx context.Context) *slog.Logger {
	if logger := ctx.Value("logger"); logger != nil {
		if l, ok := logger.(*slog.Logger); ok {
			return l
		}
	}
	return Get()
}

// 便捷方法（兼容旧 API）
func Info(msg string, args ...any) {
	Get().Info(msg, args...)
}

func Error(msg string, args ...any) {
	Get().Error(msg, args...)
}

func Debug(msg string, args ...any) {
	Get().Debug(msg, args...)
}

func Warn(msg string, args ...any) {
	Get().Warn(msg, args...)
}

// Fatal 记录错误并退出
func Fatal(msg string, args ...any) {
	Get().Error(msg, args...)
	os.Exit(1)
}

// InfoContext 带 context 的 Info
func InfoContext(ctx context.Context, msg string, args ...any) {
	WithContext(ctx).InfoContext(ctx, msg, args...)
}

// ErrorContext 带 context 的 Error
func ErrorContext(ctx context.Context, msg string, args ...any) {
	WithContext(ctx).ErrorContext(ctx, msg, args...)
}

// DebugContext 带 context 的 Debug
func DebugContext(ctx context.Context, msg string, args ...any) {
	WithContext(ctx).DebugContext(ctx, msg, args...)
}

// WarnContext 带 context 的 Warn
func WarnContext(ctx context.Context, msg string, args ...any) {
	WithContext(ctx).WarnContext(ctx, msg, args...)
}

// Module 返回带模块标识的 logger，module 字段位于 level 后面
func Module(name string) *slog.Logger {
	handler := &CustomTextHandler{
		out:    os.Stdout,
		level:  slog.LevelDebug,
		module: name,
		mu:     &handlerMu,
	}
	return slog.New(handler)
}
