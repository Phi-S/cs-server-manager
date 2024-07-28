package user_logs

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const SystemInfoLogType = "system-info"
const SystemErrorLogType = "system-error"
const ServerLogType = "server"
const SteamcmdLogType = "steamcmd"

type LogEntry struct {
	Timestamp time.Time `json:"timestamp"`
	LogType   string    `json:"log-type"`
	Message   string    `json:"message"`
}

func (e *LogEntry) Format() string {
	return fmt.Sprintf("%s | %s | %s\n", e.Timestamp.Format(time.RFC3339), e.LogType, e.Message)
}

func NewLogEntryFromLogLine(logLine string) (LogEntry, error) {
	split := strings.Split(logLine, "|")
	if len(split) != 3 {
		return LogEntry{}, errors.New("logLine string malformed")
	}

	t, err := time.Parse(time.RFC3339, split[0])
	if err != nil {
		return LogEntry{}, err
	}

	return LogEntry{
		Timestamp: t,
		LogType:   split[1],
		Message:   split[2],
	}, nil
}

func NewLogEntry(timestamp time.Time, logType string, msg string) LogEntry {
	logType = strings.ReplaceAll(logType, "\n", " ")
	logType = strings.ReplaceAll(logType, "|", " ")
	logType = strings.ReplaceAll(logType, "\r", "")

	msg = strings.ReplaceAll(msg, "\n", " ")
	msg = strings.ReplaceAll(msg, "|", " ")
	msg = strings.ReplaceAll(msg, "\r", "")

	return LogEntry{
		timestamp,
		logType,
		msg,
	}
}

type LogWriter struct {
	name            string
	logDir          string
	currentFileName string

	lastEntries []LogEntry

	writeInterval time.Duration
	fileWriter    *os.File
	lock          sync.RWMutex
}

func NewLogWriter(logDirIn string, nameIn string) (*LogWriter, error) {
	fileName := fmt.Sprintf("%s_%s.log", time.Now().UTC().Format(time.RFC3339), nameIn)

	filePath := filepath.Join(logDirIn, fileName)
	f, err := os.Create(filePath)
	if err != nil {
		return nil, err
	}

	lw := LogWriter{
		name:            nameIn,
		logDir:          logDirIn,
		currentFileName: fileName,

		lastEntries: make([]LogEntry, 100),

		writeInterval: time.Second,
		fileWriter:    f,
	}

	return &lw, nil
}

func (w *LogWriter) WriteLog(timestamp time.Time, logType string, msg string) error {
	w.lock.Lock()
	defer w.lock.Unlock()

	logEntry := NewLogEntry(timestamp, logType, msg)
	_, err := w.fileWriter.WriteString(logEntry.Format())
	if err != nil {
		return err
	}

	w.lastEntries = append(w.lastEntries, logEntry)
	return nil
}

func (w *LogWriter) WriteSystemInfoLog(timestamp time.Time, msg string) error {
	return w.WriteLog(timestamp, SystemInfoLogType, msg)
}

func (w *LogWriter) WriteSystemErrorLog(timestamp time.Time, msg string) error {
	return w.WriteLog(timestamp, SystemErrorLogType, msg)
}

func (w *LogWriter) WriteServerLog(timestamp time.Time, msg string) error {
	return w.WriteLog(timestamp, ServerLogType, msg)
}

func (w *LogWriter) WriteSteamcmdLog(timestamp time.Time, msg string) error {
	return w.WriteLog(timestamp, SteamcmdLogType, msg)
}

func (w *LogWriter) Close() {
	w.lock.Lock()
	defer w.lock.Unlock()

	if w.fileWriter != nil {
		w.fileWriter.Close()
	}
}

func (w *LogWriter) GetLogs(last int) []LogEntry {
	w.lock.RLock()
	defer w.lock.Unlock()

	count := last
	if len(w.lastEntries) < last {
		count = len(w.lastEntries)
	}

	result := make([]LogEntry, 0, count)
	for i := 0; i < count; i++ {
		result = append(result, w.lastEntries[i])
	}

	return result
}

func (w *LogWriter) GetLogsSince(since time.Time, count int) []LogEntry {
	w.lock.RLock()
	defer w.lock.Unlock()

	if len(w.lastEntries) < count {
		count = len(w.lastEntries)
	}

	result := make([]LogEntry, 0, count)
	for _, e := range w.lastEntries {
		if e.Timestamp.After(since) {
			result = append(result, e)
		}
	}

	return result
}
