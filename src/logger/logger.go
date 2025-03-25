package logger

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/rs/zerolog"
)

var once sync.Once
var logger zerolog.Logger

func GetLogger() zerolog.Logger {
	once.Do(func() {
		logger = setUpLogger()
	})

	return logger
}

func setUpLogger() zerolog.Logger {
	_, err := os.Stat("logs/logs.log")

	if os.IsNotExist(err) {
		// Förklaring: rwx | 7 | Read, write and execute för user. Mer info finns här: https://stackoverflow.com/a/31151508
		os.MkdirAll("logs/", 0700)
	}

	runLogFile, logFileError := os.OpenFile("logs/logs.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0664)
	if logFileError != nil {
		fmt.Fprintf(os.Stderr, "Kunde inte öppna, hitta eller skapa filen logs.log\n%v", logFileError)
		// TODO: Kolla om det är rätt. Man kanske alltid ska stänga filen, även om det gick bra
		defer runLogFile.Close()
	}

	consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	multi := zerolog.MultiLevelWriter(consoleWriter, runLogFile)
	log := zerolog.New(multi).With().Timestamp().Caller().Stack().Logger()

	return log
}
