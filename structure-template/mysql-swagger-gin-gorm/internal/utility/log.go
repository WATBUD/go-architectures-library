package utility

import (
	"os"

	log "github.com/sirupsen/logrus"
)

func SetupLog() {
	log.Infoln("SetupLog()...")
	log.SetFormatter(&log.TextFormatter{
		ForceColors:            true,
		FullTimestamp:          true,
		TimestampFormat:        "2006-01-02 15:04:05",
		DisableLevelTruncation: true,
	})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
	log.Infoln("This is a test log entry")
}
