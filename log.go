package coord_log

import (
	"context"
	cfg "github.com/jdy879526487/coord-cfg"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

var Logger *logrus.Logger

func init() {
	Logger = logrus.New()
	// 指定时间输出格式,显示字段等
	Logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "@timestamp",
			logrus.FieldKeyLevel: "@level",
			logrus.FieldKeyMsg:   "@message",
			logrus.FieldKeyFunc:  "@caller",
		},
		PrettyPrint: true,
	})

	level := logrus.DebugLevel
	if os.Getenv("RUN_LOCAL") != "" {
		level = logrus.WarnLevel
	}
	ctx := context.Background()
	Logger.SetLevel(level)

	// 如果设置了上报用的kafka,将日志上报
	topic := cfg.Get(ctx, KAFKA_TOPIC_KEY)
	broker := cfg.Get(ctx, KAFKA_BROKER_KEY)
	if topic != "" {
		kaCli, err := NewKafkaCli(topic, broker)
		if err != nil {
			logrus.Error("failed create kafka ", err)
			return
		}
		Logger.AddHook(kaCli)
	}
}

type argLogger func(args ...interface{})
type formatLogger func(format string, args ...interface{})

// 方便外部可以直接调用
var (
	Debug, Info, Warn, Warning, Error, Fatal argLogger = Logger.Debug, Logger.Info,
		Logger.Warn, Logger.Warning, Logger.Error, Logger.Fatal
	Debugf, Infof, Warnf, Warningf, Errorf, Fatalf formatLogger = Logger.Debugf, Logger.Infof,
		Logger.Warnf, Logger.Warningf, Logger.Errorf, Logger.Fatalf
)
