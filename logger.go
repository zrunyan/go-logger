package logger

import (
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"time"
)

//Set all of the constants for the different colors
const (
	FormatPrefix        = "\033[1;30m"              //Black or gray with a light color
	FormatFatal         = "\033[1;34m"              //Blue
	FormatError         = "\033[2;31m"              //Red
	FormatWarning       = "\033[1;33m"              //Yellow
	FormatNotice        = "\033[2;32m"              //Green
	FormatInfo          = ""                        //Standard white
	FormatDebug         = "\033[1;36m"              //Cyan
	FormatOff           = "\033[0m"                 //Stop formatting, turn it off
	FormatTimeDisplay   = "2006-01-02 15:04:05"     //The format we want for logs
	FormatStringDisplay = "%s[%s %s:%d] %s%s%s%s\n" //The format replacement for the start
)

//Set the log levels with an int value starting at 0
const (
	LoglevelOff     = iota
	LoglevelFatal   //Worst possible case[red]
	LoglevelError   //Something very bad happened and needs to be taken care of
	LoglevelWarning //Something kinda bad happened but not really an issue
	LoglevelInfo    //Informative of what is going on
	LoglevelNotice  //Notice of something changing / setting
	LoglevelDebug   //Debug information to help debug the code
)

//LogMessanger is the interface that allows us to write messages
type logMessanger interface {
	GetFormattedMessage() string
}

//Logger is the type struct
type Logger struct {
	localWriter io.Writer
	logLevel    int
	messageChan chan logMessanger
}

//SetLogLevel will change the current log level of the struct
func (l *Logger) SetLogLevel(newLevel int) {

	//Set it at the struct level here
	l.logLevel = newLevel
}

//GetLogLevel will return the current log level of the struct
func (l *Logger) GetLogLevel() int {

	return l.logLevel
}

//Fatal will format messages and write them as needed in Fatal format
func (l *Logger) Fatal(message ...interface{}) {

	//Check the log level to make sure we should do anything at all
	if l.logLevel < LoglevelFatal {

		return
	}

	//Grab the time of the log
	logTime := time.Now().Format(FormatTimeDisplay)

	//Get the file and line number that triggered the log
	fileName, line := l.getFileAndLine()

	//Create a new log message
	writeMessage := &LogMessage{
		fileName,
		LoglevelFatal,
		message,
		line,
		logTime,
		FormatFatal,
	}

	//Send the log message to the writer
	l.messageChan <- writeMessage
}

//Error will format messages and write them as needed in Error format
func (l *Logger) Error(message ...interface{}) {

	//Check the log level to make sure we should do anything at all
	if l.logLevel < LoglevelError {

		return
	}

	//Grab the time of the log
	logTime := time.Now().Format(FormatTimeDisplay)

	//Get the file and line number that triggered the log
	fileName, line := l.getFileAndLine()

	//Create a new log message
	writeMessage := &LogMessage{
		fileName,
		LoglevelFatal,
		message,
		line,
		logTime,
		FormatError,
	}

	//Send the log message to the writer
	l.messageChan <- writeMessage
}

//Warning will format messages and write them as needed in Warning format
func (l *Logger) Warning(message ...interface{}) {

	//Check the log level to make sure we should do anything at all
	if l.logLevel < LoglevelWarning {

		return
	}

	//Grab the time of the log
	logTime := time.Now().Format(FormatTimeDisplay)

	//Get the file and line number that triggered the log
	fileName, line := l.getFileAndLine()

	//Create a new log message
	writeMessage := &LogMessage{
		fileName,
		LoglevelWarning,
		message,
		line,
		logTime,
		FormatWarning,
	}

	//Send the log message to the writer
	l.messageChan <- writeMessage
}

//Info will format messages and write them as needed in Info format
func (l *Logger) Info(message ...interface{}) {

	//Check the log level to make sure we should do anything at all
	if l.logLevel < LoglevelInfo {

		return
	}

	//Grab the time of the log
	logTime := time.Now().Format(FormatTimeDisplay)

	//Get the file and line number that triggered the log
	fileName, line := l.getFileAndLine()

	//Create a new log message
	writeMessage := &LogMessage{
		fileName,
		LoglevelInfo,
		message,
		line,
		logTime,
		FormatInfo,
	}

	//Send the log message to the writer
	l.messageChan <- writeMessage
}

//Notice will format messages and write them as needed in Notice format
func (l *Logger) Notice(message ...interface{}) {

	//Check the log level to make sure we should do anything at all
	if l.logLevel < LoglevelNotice {

		return
	}

	//Grab the time of the log
	logTime := time.Now().Format(FormatTimeDisplay)

	//Get the file and line number that triggered the log
	fileName, line := l.getFileAndLine()

	//Create a new log message
	writeMessage := &LogMessage{
		fileName,
		LoglevelNotice,
		message,
		line,
		logTime,
		FormatNotice,
	}

	//Send the log message to the writer
	l.messageChan <- writeMessage
}

//Debug will format messages and write them as needed in Debug format
func (l *Logger) Debug(message ...interface{}) {

	//Check the log level to make sure we should do anything at all
	if l.logLevel < LoglevelDebug {

		return
	}

	//Grab the time of the log
	logTime := time.Now().Format(FormatTimeDisplay)

	//Get the file and line number that triggered the log
	fileName, line := l.getFileAndLine()

	//Create a new log message
	writeMessage := &LogMessage{
		fileName,
		LoglevelDebug,
		message,
		line,
		logTime,
		FormatDebug,
	}

	//Send the log message to the writer
	l.messageChan <- writeMessage
}

//getFileAndLine will grab the file name and line number of the calling function to log it
func (l *Logger) getFileAndLine() (fileName string, line int) {

	//Lets get the filename and line number from l function
	_, filePath, line, success := runtime.Caller(2)
	if !success {
		//TODO: Figure out what to do here since we are already in the Logger....
	}

	// Since the file comes in as a full path lets split it into the dir and filename
	_, fileName = filepath.Split(filePath)

	return
}

//startListening will actually start recieving the channel
func (l *Logger) startListening() {

	//Loop over the channel to get the data
	for n := range l.messageChan {
		output := n.GetFormattedMessage()
		io.WriteString(l.localWriter, output)
	}
}

//NewLogger will return a new version of the Logger
func NewLogger(logWriteFile interface{}) (*Logger, error) {

	var logWriter io.Writer

	if logWriteFile != nil && logWriteFile.(string) != "" {
		// Log file specified, create file if necessary and open for writing.

		LoggerFile := logWriteFile.(string)

		//Lets create an instance of a log file ready for writing here
		//TODO: Need the log file configurable for both a config file and constructor

		err := os.MkdirAll(path.Dir(LoggerFile), 0755)
		if err != nil {
			return nil, err
		}

		logFile, err := os.OpenFile(LoggerFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

		//TODO: Probably do something better here
		if err != nil {
			return nil, err
		}

		//Create a multi writer so we can output to both the log file and standard output
		logWriter = io.MultiWriter(logFile, os.Stdout)

	} else {
		// No log file given, log to stdout
		logWriter = os.Stdout
	}

	//Create an initilized struct here to return
	logs := &Logger{
		localWriter: logWriter,
		logLevel:    LoglevelDebug,
		messageChan: make(chan logMessanger),
	}

	go logs.startListening()

	return logs, nil
}

//LogMessage is a struct that will hold all of the information on a log message
type LogMessage struct {
	fileName string
	level    int
	message  []interface{}
	line     int
	time     string
	format   string
}

//GetFormattedMessage allows LogMessage to satisfy the interface to the channel
func (m *LogMessage) GetFormattedMessage() (logString string) {

	var messageString string
	var messages []string
	var message string

	//Messages are going to be an array of interfaces so lets loop over them and cast them
	for _, s := range m.message {
		switch s.(type) {
		case float32, float64:
			message = fmt.Sprintf("%v", s)

		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
			message = fmt.Sprintf("%v", s)

		case []byte:
			message = fmt.Sprintf("%s", s)

		case string:
			message = s.(string)

		default:
			message = fmt.Sprintf("(%s) %v", reflect.TypeOf(s), s)
		}

		messages = append(messages, message)
	}
	//Add a separator to each of the log messages
	messageString = strings.Join(messages, " ")

	//Build the beginning of the string for the file
	logString = fmt.Sprintf(FormatStringDisplay, FormatPrefix, m.time, m.fileName,
		m.line, FormatOff, m.format, messageString, FormatOff)

	//Implicitly return logString
	return
}
