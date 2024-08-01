package logwrt

import (
    "errors"
    "fmt"
    "log/slog"
    "os"
    "path/filepath"
    "regexp"
    "strings"
    "sync"
    "time"
)

const defaultGetLogsLimit = 500
const logHistoryLimit = 1000
const logFileRolloverCheckAfterLines = 5000
const logFileRolloverSizeMiB = 2

type LogEntry struct {
    Timestamp time.Time `json:"timestamp"`
    LogType   string    `json:"log-type"`
    Message   string    `json:"message"`
}

func (e *LogEntry) Format() string {
    return fmt.Sprintf("%s | %s | %s\n", e.LogType, e.Timestamp.Format(time.RFC3339Nano), e.Message)
}

func NewLogEntryFromLogLine(logLine string) (LogEntry, error) {
    split := strings.Split(logLine, "|")
    if len(split) != 3 {
        return LogEntry{}, errors.New("string malformed")
    }
    timestampStr := strings.TrimSpace(split[0])
    logType := strings.TrimSpace(split[1])
    message := strings.TrimSpace(split[2])

    timestamp, err := time.Parse(time.RFC3339Nano, timestampStr)
    if err != nil {
        return LogEntry{}, err
    }

    return LogEntry{
        Timestamp: timestamp,
        LogType:   logType,
        Message:   message,
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

    history                   []LogEntry
    linesWrittenSinceRollover int
    getLogsLimit              int

    fileWriter *os.File
    lock       sync.RWMutex
}

func NewLogWriter(logDirectory string, logName string) (*LogWriter, error) {
    newFileWriter, fileName, err := generateNewLogFile(logDirectory, logName)
    if err != nil {
        return nil, err
    }

    lw := LogWriter{
        name:            logName,
        logDir:          logDirectory,
        currentFileName: fileName,

        history:      make([]LogEntry, 0, 100),
        getLogsLimit: defaultGetLogsLimit,

        fileWriter: newFileWriter,
    }

    return &lw, nil
}

func generateNewLogFile(logDir string, name string) (fileWriter *os.File, fileName string, err error) {
    logFileName := fmt.Sprintf("%s_%s.log", time.Now().UTC().Format(time.RFC3339), name)
    logFilePath := filepath.Join(logDir, logFileName)

    if err := os.MkdirAll(logDir, 0777); err != nil {
        return nil, "", err
    }

    _, err = os.Stat(logFilePath)
    if err == nil {
        return nil, "", errors.New("file already exists")
    } else if os.IsNotExist(err) {
        f, err := os.Create(logFilePath)
        if err != nil {
            return nil, "", err
        }
        return f, logFileName, nil
    } else {
        return nil, "", err
    }
}

func (w *LogWriter) GetCurrentLogFilePath() string {
    return filepath.Join(w.logDir, w.currentFileName)
}

func (w *LogWriter) GetLogsLimit() int {
    return w.getLogsLimit
}

func (w *LogWriter) WriteLogEntry(entry LogEntry) error {
    w.lock.Lock()
    defer w.lock.Unlock()

    _, err := w.fileWriter.WriteString(entry.Format())
    if err != nil {
        return err
    }

    w.history = append(w.history, entry)
    w.linesWrittenSinceRollover += 1

    w.trimHistoryIfTooLarge()
    w.rolloverLogFileIfTooLarge()
    return nil
}

func (w *LogWriter) WriteLog(timestamp time.Time, logType string, msg string) error {
    logEntry := NewLogEntry(timestamp, logType, msg)
    return w.WriteLogEntry(logEntry)
}

func (w *LogWriter) trimHistoryIfTooLarge() {
    currentHistoryCount := len(w.history)
    if currentHistoryCount < logHistoryLimit {
        return
    }

    trimmedCount := currentHistoryCount / 2
    trimmedHistory := make([]LogEntry, trimmedCount)
    copy(trimmedHistory, w.history)
    w.history = trimmedHistory
    slog.Debug("log history has been trimmed", "log-writer-name", w.name, "before-trimmed-count", currentHistoryCount, "after-trimmed-count", trimmedCount)
}

func (w *LogWriter) rolloverLogFileIfTooLarge() {
    if w.linesWrittenSinceRollover < logFileRolloverCheckAfterLines {
        return
    }

    fileInfo, err := os.Stat(w.GetCurrentLogFilePath())
    if err != nil {
        slog.Debug("failed to check for log file rollover", "error", err)
        return
    }

    if fileInfo.Size() > (1024*1024)*logFileRolloverSizeMiB {

        newFileWriter, newFileName, err := generateNewLogFile(w.logDir, w.name)
        if err != nil {
            slog.Debug("log file rollover: failed to generate new log file. trying again", "error", err)
            time.Sleep(1000)

            newFileWriter, newFileName, err = generateNewLogFile(w.logDir, w.name)
            if err != nil {
                slog.Debug("failed to rollover log file. failed to generate new log file again", "error", err)
                return
            }
        }

        oldLogFilePath := w.GetCurrentLogFilePath()

        _ = w.fileWriter.Close()
        w.fileWriter = newFileWriter
        w.linesWrittenSinceRollover = 0
        w.currentFileName = newFileName

        slog.Debug("log file rolled over",
            "rollover-size-mib", logFileRolloverSizeMiB,
            "before-rollover-file-path", oldLogFilePath,
            "after-rollover-file-path", w.GetCurrentLogFilePath(),
        )
    }
}

func (w *LogWriter) Close() {
    if w.fileWriter != nil {
        _ = w.fileWriter.Close()
    }
}

func (w *LogWriter) GetLogs(last int) ([]LogEntry, error) {
    w.lock.RLock()
    defer w.lock.RUnlock()

    if last > w.getLogsLimit {
        return nil, fmt.Errorf("requested log count is bigger then limit %v", w.getLogsLimit)
    }

    if last < len(w.history) {
        result := make([]LogEntry, last)
        copy(result, w.history)
        return result, nil
    }

    if len(w.history) < last {
        last = len(w.history)
    }

    result := make([]LogEntry, 0, last)
    for i := 0; i < last; i++ {
        result = append(result, w.history[i])
    }

    return result, nil
}

func (w *LogWriter) GetLogsSince(since time.Time) ([]LogEntry, error) {
    w.lock.RLock()
    defer w.lock.RUnlock()

    result := make([]LogEntry, 0)
    for _, e := range w.history {
        if e.Timestamp.After(since) {
            result = append(result, e)
            if len(result) > w.getLogsLimit {
                return nil, fmt.Errorf("logs since %v are more then the get logs limit %v", since, w.getLogsLimit)
            }
        }
    }

    return result, nil
}

func (w *LogWriter) GetPastLogFiles() ([]string, error) {
    entries, err := os.ReadDir(w.logDir)
    if err != nil {
        return nil, err
    }

    pastLogFiles := make([]string, 0)
    for _, e := range entries {
        if e.Name() == w.currentFileName {
            continue
        }

        if !strings.Contains(e.Name(), w.name) {
            continue
        }
        pastLogFiles = append(pastLogFiles, e.Name())
    }

    return pastLogFiles, nil
}

// GetContentOfPastLogFile This operation can be expensive
func (w *LogWriter) GetContentOfPastLogFile(pastLogFileName string) ([]LogEntry, error) {
    if !regexp.MustCompile(`^[a-zA-Z\d_\-:]+.log$`).MatchString(pastLogFileName) {
        return nil, fmt.Errorf("past log file name %q is not valid", pastLogFileName)
    }

    if !strings.Contains(pastLogFileName, w.name) {
        return nil, errors.New("requested log file was not written with the current log writer")
    }
    pastLogFilePath := filepath.Join(w.logDir, pastLogFileName)

    contentBytes, err := os.ReadFile(pastLogFilePath)
    if err != nil {
        return nil, err
    }

    content := string(contentBytes)
    lines := strings.Split(content, "\n")
    result := make([]LogEntry, 0, len(lines))
    for _, line := range lines {
        if strings.TrimSpace(line) == "" {
            continue
        }

        entry, err := NewLogEntryFromLogLine(line)
        if err != nil {
            return nil, err
        }

        result = append(result, entry)
    }

    return result, nil
}
