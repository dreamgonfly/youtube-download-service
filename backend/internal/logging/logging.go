package logging

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"time"
	"youtube-download-backend/internal/config"

	"gopkg.in/Graylog2/go-gelf.v2/gelf"
)

var Logger *logConfig

const Timezone = "Asia/Seoul"

type logConfig struct {
	timezone *time.Location
	local    *log.Logger
	gelf     *gelf.TCPWriter
	hostname string
}

func InitLocalLogger() {
	tz, _ := time.LoadLocation(Timezone)
	Logger = &logConfig{
		timezone: tz,
		local:    log.New(os.Stderr, "", 0),
	}
}

func InitMultiLogger() {
	tz, _ := time.LoadLocation(Timezone)

	gelfWriter, err := gelf.NewTCPWriter(config.Config.LogServer)
	if err != nil {
		log.Fatalf("gelf.NewTCPWriter: %s", err)
	}

	hostname, err := os.Hostname()
	if err != nil {
		log.Fatalf("os.Hostname: %s", err)
	}

	Logger = &logConfig{
		timezone: tz,
		local:    log.New(os.Stderr, "", 0),
		gelf:     gelfWriter,
		hostname: hostname,
	}
}

func (l *logConfig) Info(s string) {
	now := time.Now()
	t := now.In(l.timezone).Format("2006-01-02 15:04:05 MST")

	msg := fmt.Sprintf("%s INFO: %s", t, s)

	l.local.Println(msg)

	if l.gelf != nil {
		m := &gelf.Message{
			Version:  "1.1",
			Host:     l.hostname,
			Short:    msg,
			Full:     msg,
			TimeUnix: float64(now.UnixNano()) / float64(time.Second),
			Level:    6, // info
			Extra: map[string]interface{}{
				"_application": "youtube-download-backend",
				"_env":         config.Env,
				"_localTime":   t,
			},
		}

		if err := l.gelf.WriteMessage(m); err != nil {
			l.local.Println(err)
		}
	}
}

func (l *logConfig) Error(err error) {
	now := time.Now()
	t := now.In(l.timezone).Format("2006-01-02 15:04:05 MST")

	_, file, line, ok := runtime.Caller(1)
	if !ok {
		l.local.Println("runtime.Caller error")
	}

	var msg string = fmt.Sprintf("%s ERROR: %s %s:%d", t, err.Error(), file, line)

	l.local.Println(msg)

	if l.gelf != nil {
		m := &gelf.Message{
			Version:  "1.1",
			Host:     l.hostname,
			Short:    msg,
			Full:     msg,
			TimeUnix: float64(now.UnixNano()) / float64(time.Second),
			Level:    3, // error
			Extra: map[string]interface{}{
				"_application": "youtube-download-backend",
				"_env":         config.Env,
				"_localTime":   t,
				"_file":        file,
				"_line":        line,
			},
		}
		if err = l.gelf.WriteMessage(m); err != nil {
			l.local.Println(err)
		}
	}
}

type HTTPRequestData struct {
	Method    string
	URI       string
	Referer   string
	IPAddress string
	Status    int
	// number of bytes of the response sent
	Size int64
	// how long did it take to
	Duration  time.Duration
	UserAgent string
	Body      string
}

func (l *logConfig) LogHTTP(d *HTTPRequestData) {
	now := time.Now()
	t := now.In(l.timezone).Format("2006-01-02 15:04:05 MST")

	msg := fmt.Sprintf("%s INFO: %s %s %d %s %s", t, d.Method, d.URI, d.Status, d.Duration.String(), d.IPAddress)

	l.local.Println(msg)

	if l.gelf != nil {

		m := &gelf.Message{
			Version:  "1.1",
			Host:     l.hostname,
			Short:    msg,
			Full:     msg,
			TimeUnix: float64(now.UnixNano()) / float64(time.Second),
			Level:    6, // info
			Extra: map[string]interface{}{
				"_application": "youtube-download-backend",
				"_env":         config.Env,
				"_method":      d.Method,
				"_uri":         d.URI,
				"_referer":     d.Referer,
				"_ipaddr":      d.IPAddress,
				"_status":      d.Status,
				// number of bytes of the response sent
				"_size": d.Size,
				// how long did it take to
				"_duration":  d.Duration,
				"_userAgent": d.UserAgent,
				"_body":      d.Body,
				"_localTime": t,
			},
		}

		if err := l.gelf.WriteMessage(m); err != nil {
			l.local.Println(err)
		}
	}
}
