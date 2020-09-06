package logger

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

var (
	baseTimestamp = time.Now()
)

type fieldKey string

// FieldMap allows customization of the key names for default fields.
type FieldMap map[fieldKey]string

func (f FieldMap) resolve(key fieldKey) string {
	if k, ok := f[key]; ok {
		return k
	}

	return string(key)
}

const (
	fileSufixLevel                = 2
	funcSufixLevel                = 1
	defaultTimestampFormat        = time.RFC3339
	defaultHiddenField     string = "Identificador"
	warnLevel                     = "warn"
	FieldKeyMsg                   = "msg"
	FieldKeyLevel                 = "level"
	FieldKeyTime                  = "time"
	fieldKeyLogrusError           = "logrus_error"
	FieldKeyFunc                  = "func"
	FieldKeyFile                  = "file"
)

type TextFormatter struct {

	// Force formatted layout, even for non-TTY output.
	ForceFormatting bool

	// Disable timestamp logging. useful when output is redirected to logging
	// system that already adds timestamps.
	DisableTimestamp bool

	// Disable the conversion of the log levels to uppercase
	DisableUppercase bool

	// Enable logging the full timestamp when a TTY is attached instead of just
	// the time passed since beginning of execution.
	FullTimestamp bool

	// Enable logging of the function
	EnableCallerFunc bool

	// Timestamp format to use for display when a full timestamp is printed.
	TimestampFormat string

	// The fields are sorted by default for a consistent output. For applications
	// that log extremely frequently and don't use the JSON formatter this may not
	// be desired.
	DisableSorting bool

	// Wrap empty fields in quotes if true.
	QuoteEmptyFields bool

	// Can be set to the override the default quoting character "
	// with something else. For example: ', or `.
	QuoteCharacter string

	// Pad msg field with spaces on the right for display.
	// The value for this parameter will be the size of padding.
	// Its default value is zero, which means no padding will be applied for msg.
	SpacePadding int

	//Hidden field will not be printed, usually used to split log with MwHook
	HiddenField string

	// FieldMap allows users to customize the names of keys for default fields.
	// As an example:
	// formatter := &JSONFormatter{
	//   	FieldMap: FieldMap{
	// 		 FieldKeyTime:  "@timestamp",
	// 		 FieldKeyLevel: "@level",
	// 		 FieldKeyMsg:   "@message",
	// 		 FieldKeyFunc:  "@caller",
	//    },
	// }
	FieldMap FieldMap

	sync.Once
}

func miniTS() int {
	return int(time.Since(baseTimestamp) / time.Second)
}

func (f *TextFormatter) init() {
	if len(f.QuoteCharacter) == 0 {
		f.QuoteCharacter = "\""
	}
}

func (f *TextFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var (
		b       *bytes.Buffer
		fileVal string
	)
	data := make(logrus.Fields)
	for k, v := range entry.Data {
		data[k] = v
	}
	prefixFieldClashes(data, f.FieldMap, entry.HasCaller())
	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}
	fixedKeys := make([]string, 0, 4+len(data))
	fieldKeyFunc := f.FieldMap.resolve(FieldKeyFunc)
	fieldKeyFile := f.FieldMap.resolve(FieldKeyFile)
	fieldKeyTime := f.FieldMap.resolve(FieldKeyTime)
	fieldKeyLevel := f.FieldMap.resolve(FieldKeyLevel)
	fieldKeyMsg := f.FieldMap.resolve(FieldKeyMsg)

	if !f.DisableTimestamp {
		fixedKeys = append(fixedKeys, fieldKeyTime)
	}
	fixedKeys = append(fixedKeys, fieldKeyLevel)
	if entry.Message != "" {
		fixedKeys = append(fixedKeys, fieldKeyMsg)
	}
	if entry.HasCaller() {
		fixedKeys = append(fixedKeys, fieldKeyFunc, fieldKeyFile)
	}

	if !f.DisableSorting {
		sort.Strings(keys)
		fixedKeys = append(fixedKeys, keys...)
	}
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	f.Do(func() { f.init() })

	isFormatted := f.ForceFormatting

	timestampFormat := f.TimestampFormat
	if timestampFormat == "" {
		timestampFormat = defaultTimestampFormat
	}

	if entry.HasCaller() {
		fileVal = fmt.Sprintf("%s:%d", extractSufix(entry.Caller.File, fileSufixLevel), entry.Caller.Line)
		if f.EnableCallerFunc {
			fileVal = fmt.Sprintf("%s - %s", fileVal, extractSufix(entry.Caller.Function, funcSufixLevel))
		}
	}

	if isFormatted {
		f.printFormatted(b, entry, keys, timestampFormat, fileVal)
	} else {
		for _, key := range fixedKeys {
			var value interface{}
			switch {
			case key == fieldKeyTime:
				value = entry.Time.Format(timestampFormat)
			case key == fieldKeyLevel:
				value = entry.Level.String()
			case key == fieldKeyMsg:
				value = entry.Message
			case key == fieldKeyFunc && entry.HasCaller():
				if !f.EnableCallerFunc {
					continue
				}
				value = entry.Caller.Function
			case key == fieldKeyFile && entry.HasCaller():
				value = fileVal
			default:
				value = data[key]
			}
			f.appendKeyValue(b, key, value)
		}
	}

	b.WriteByte('\n')
	return b.Bytes(), nil
}

func (f *TextFormatter) printFormatted(b *bytes.Buffer, entry *logrus.Entry, keys []string, timestampFormat, fileVal string) {
	var (
		levelText string
		fields    string
	)
	if entry.Level != logrus.WarnLevel {
		levelText = entry.Level.String()
	} else {
		levelText = warnLevel
	}

	if !f.DisableUppercase {
		levelText = strings.ToUpper(levelText)
	}

	if f.HiddenField == "" {
		f.HiddenField = defaultHiddenField
	}

	level := fmt.Sprintf("[%-5s]", levelText)
	prefix := ""
	message := strings.TrimSuffix(entry.Message, "\n")

	if entry.HasCaller() {
		prefix = fmt.Sprintf("[%s]", fileVal)
	}
	messageFormat := "%s"
	if f.SpacePadding != 0 {
		messageFormat = fmt.Sprintf("%%-%ds", f.SpacePadding)
	}

	if f.DisableTimestamp {
		_, _ = fmt.Fprintf(b, "%s%s", level, prefix)
	} else {
		var timestamp string
		if !f.FullTimestamp {
			timestamp = fmt.Sprintf("[%04d]", miniTS())
		} else {
			timestamp = fmt.Sprintf("[%s]", entry.Time.Format(timestampFormat))
		}
		if prefix != "" {
			_, _ = fmt.Fprintf(b, "%s %s %s", timestamp, level, prefix)
		} else {
			_, _ = fmt.Fprintf(b, "%s %s", timestamp, level)
		}
	}
	for _, k := range keys {
		if k != f.HiddenField {
			v := entry.Data[k]
			fields += fmt.Sprintf(" [%s=%+v]", k, v)
		}
	}
	if fields != "" {
		_, _ = fmt.Fprintf(b, " –%s – "+messageFormat, fields, message)
	} else {
		_, _ = fmt.Fprintf(b, " – "+messageFormat, message)
	}
}

func (f *TextFormatter) needsQuoting(text string) bool {
	if f.QuoteEmptyFields && len(text) == 0 {
		return true
	}
	for _, ch := range text {
		if !((ch >= 'a' && ch <= 'z') ||
			(ch >= 'A' && ch <= 'Z') ||
			(ch >= '0' && ch <= '9') ||
			ch == '-' || ch == '.' || ch == '_' || ch == '/' || ch == '@' || ch == '^' || ch == '+') {
			return true
		}
	}
	return false
}

func extractSufix(data string, level int) string {
	n := 0
	for i := len(data) - 1; i > 0; i-- {
		if data[i] == '/' {
			n++
			if n >= level {
				data = data[i+1:]
				break
			}
		}
	}
	return data
}

func (f *TextFormatter) appendKeyValue(b *bytes.Buffer, key string, value interface{}) {
	if b.Len() > 0 {
		b.WriteByte(' ')
	}
	b.WriteString(key)
	b.WriteByte('=')
	f.appendValue(b, value)
}

func (f *TextFormatter) appendValue(b *bytes.Buffer, value interface{}) {
	switch value := value.(type) {
	case string:
		if !f.needsQuoting(value) {
			b.WriteString(value)
		} else {
			_, _ = fmt.Fprintf(b, "%s%v%s", f.QuoteCharacter, value, f.QuoteCharacter)
		}
	case error:
		errmsg := value.Error()
		if !f.needsQuoting(errmsg) {
			b.WriteString(errmsg)
		} else {
			_, _ = fmt.Fprintf(b, "%s%v%s", f.QuoteCharacter, errmsg, f.QuoteCharacter)
		}
	default:
		_, _ = fmt.Fprint(b, value)
	}
}

// This is to not silently overwrite `time`, `msg` and `level` fields when
// dumping it. If this code wasn't there doing:
//
//  logrus.WithField("level", 1).Info("hello")
//
// would just silently drop the user provided level. Instead with this code
// it'll be logged as:
//
//  {"level": "info", "fields.level": 1, "msg": "hello", "time": "..."}
func prefixFieldClashes(data logrus.Fields, fieldMap FieldMap, reportCaller bool) {
	timeKey := fieldMap.resolve(FieldKeyTime)
	if t, ok := data[timeKey]; ok {
		data["fields."+timeKey] = t
		delete(data, timeKey)
	}

	msgKey := fieldMap.resolve(FieldKeyMsg)
	if m, ok := data[msgKey]; ok {
		data["fields."+msgKey] = m
		delete(data, msgKey)
	}

	levelKey := fieldMap.resolve(FieldKeyLevel)
	if l, ok := data[levelKey]; ok {
		data["fields."+levelKey] = l
		delete(data, levelKey)
	}

	logrusErrKey := fieldMap.resolve(fieldKeyLogrusError)
	if l, ok := data[logrusErrKey]; ok {
		data["fields."+logrusErrKey] = l
		delete(data, logrusErrKey)
	}

	// If reportCaller is not set, 'func' will not conflict.
	if reportCaller {
		funcKey := fieldMap.resolve(FieldKeyFunc)
		if l, ok := data[funcKey]; ok {
			data["fields."+funcKey] = l
		}
		fileKey := fieldMap.resolve(FieldKeyFile)
		if l, ok := data[fileKey]; ok {
			data["fields."+fileKey] = l
			delete(data, fileKey)
		}
	}
}
