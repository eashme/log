package coord_log

import (
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func init() {
	log.SetNoLock()
	log.SetReportCaller(true)
	log.SetFormatter()
	log.SetOutput()
	log.SetLevel()
	log.ReplaceHooks()
	log.SetFormatter()
	log.Writer()
}
