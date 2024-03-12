package logger

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type ServerLogger struct {
	file   *os.File
	logger *log.Logger
}

const (
	logsFileName  = "logs.txt"
	logTimeFormat = "2006-01-02 15:04:05.000 MST"
)

func NewServerLogger() ServerLogger {
	sl := ServerLogger{}
	sl.file = sl.openLogsFile()
	wrt := io.MultiWriter(os.Stdout, sl.file)
	sl.logger = log.New(wrt, "", 0)
	return sl
}

func (sl ServerLogger) HTTPLogger(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		handler.ServeHTTP(w, r)
		duration := time.Since(start)
		sl.Logln(
			fmt.Sprintf(
				"Method: %s, URL: %s, Duration: %s",
				r.Method,
				r.URL.Path,
				duration,
			),
		)
	})
}

func (sl ServerLogger) Logln(s string) {
	sl.logger.Println(prefix(s))
}

func (sl ServerLogger) openLogsFile() *os.File {
	file, err := os.OpenFile(logsFileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0o666)
	if err != nil {
		sl.Logln(fmt.Sprintf("Cannot open file \"%s\": %+v", logsFileName, err))
		return nil
	}
	fmt.Fprintf(file, "\n\n\n%s\n", prefix(">>>>> SERVER STARTED"))
	return file
}

func (sl *ServerLogger) CloseLogger() {
	if sl.file != nil {
		sl.file.Close()
	}
}

func prefix(s string) string {
	return fmt.Sprintf("%s --- %s", time.Now().Format(logTimeFormat), s)
}
