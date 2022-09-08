// @Author : detaohe
// @File   : zap
// @Description:
// @Date   : 2022/9/4 18:13

package util

import (
	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"os"
	"path/filepath"
	"time"
)

type sysLogger struct {
	logger   *zap.Logger
	filename string
	service  string
	maxAge   time.Duration
	rotate   time.Duration
	debug    bool
	global   bool
}

type LogOption func(logger *sysLogger)

func WithMaxAge(maxAge time.Duration) LogOption {
	return func(logger *sysLogger) {
		logger.maxAge = maxAge
	}
}
func WithRotate(rotate time.Duration) LogOption {
	return func(logger *sysLogger) {
		logger.rotate = rotate
	}
}
func WithDebug(debug bool) LogOption {
	return func(logger *sysLogger) {
		logger.debug = debug
	}
}

func InitZapLogger(logPath string, service string, options ...LogOption) (*zap.Logger, error) {
	//if !strings.HasSuffix(logPath, "/") {
	//	logPath += "/"
	//}
	l := &sysLogger{
		logger:   nil,
		filename: filepath.Join(logPath, service+".log"),
		service:  service,
		maxAge:   time.Hour * 24 * 7,
		rotate:   time.Hour * 24,
		debug:    true,
		global:   true,
	}
	for _, o := range options {
		o(l)
	}
	level := zap.LevelEnablerFunc(func(lv zapcore.Level) bool {
		return lv >= zapcore.DebugLevel
	})
	zapWriter := l.rotateWriter()
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(NewDevelopmentEncoderConfig()),
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(zapWriter)),
		level,
	)
	l.logger = zap.New(core, zap.AddCaller(), zap.Fields(zap.String("svc", l.service)))
	if l.global {
		zap.ReplaceGlobals(l.logger)
		return zap.L(), nil
	}
	return l.logger, nil
}

func (s *sysLogger) rotateWriter() io.Writer {
	writer, err := rotatelogs.New(
		s.filename+".%Y%m%d",
		rotatelogs.WithLinkName(s.filename),
		rotatelogs.WithMaxAge(s.maxAge),
		rotatelogs.WithRotationTime(s.rotate),
	)
	if err != nil {
		panic(err)
	}
	return writer
}

func NewDevelopmentEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		// Keys can be anything except the empty string.
		TimeKey:        "T",
		LevelKey:       "L",
		NameKey:        "N",
		CallerKey:      "C",
		MessageKey:     "M",
		StacktraceKey:  "S",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}
