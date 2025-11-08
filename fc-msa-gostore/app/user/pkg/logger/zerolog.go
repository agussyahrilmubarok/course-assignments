package logger

import (
	"os"
	"path/filepath"
	"sync"

	"example.com/user/pkg/config"
	"github.com/rs/zerolog"
)

var (
	lock           = &sync.Mutex{}
	loggerInstance *zerolog.Logger
)

func GetLogger(cfg *config.Config) (*zerolog.Logger, error) {
	if loggerInstance == nil {
		lock.Lock()
		defer lock.Unlock()

		if loggerInstance == nil {
			logDir := filepath.Dir(cfg.Logger.Filepath)
			if err := os.MkdirAll(logDir, os.ModePerm); err != nil {
				return nil, err
			}

			logFile, err := os.OpenFile(cfg.Logger.Filepath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
			if err != nil {
				return nil, err
			}

			multi := zerolog.MultiLevelWriter(os.Stdout, logFile)

			level, err := zerolog.ParseLevel(cfg.Logger.Level)
			if err != nil {
				level = zerolog.InfoLevel
			}
			zerolog.SetGlobalLevel(level)

			logger := zerolog.New(multi).With().Timestamp().Logger()

			loggerInstance = &logger
			//fmt.Println("Creating zerolog singleton instance now.")
		} else {
			//fmt.Println("Zerolog singleton instance already created.")
		}
	} else {
		// fmt.Println("Zerolog singleton instance already created.")
	}

	return loggerInstance, nil
}
