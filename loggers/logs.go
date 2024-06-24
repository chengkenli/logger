/*
 *@author  chengkenli
 *@project uar
 *@package util
 *@file    logger
 *@date    2024/6/17 13:28
 */

package loggers

import (
	"fmt"
	"os"
	"reflect"
	"time"
)

// LogLevel 日志级别枚举
type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARNING
	ERROR
	FATAL
)

func (ll LogLevel) String() string {
	switch ll {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARNING:
		return "WARNING"
	case ERROR:
		return "ERROR"
	case FATAL:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

// Logger 是我们的日志结构体
type Logger struct {
	level LogLevel
}

// LoggersConsole 创建一个新的Logger实例
// logger := LoggersConsole(logger.INFO)
func LoggersConsole(level LogLevel) *Logger {
	return &Logger{
		level: level,
	}
}

// Log 打印一条具有时间戳和日志级别的消息，支持任意类型的参数
func (l *Logger) Log(level LogLevel, v ...interface{}) {
	if level < l.level {
		return // 如果日志级别低于当前设置的级别，则不打印
	}

	// 获取当前时间并格式化
	now := time.Now().Format("2006-01-02 15:04:05")

	// 打印时间、日志级别、文件名、行号和函数名
	fmt.Printf("[%s] [%s] : ", now, level)

	// 打印参数，这里使用反射来处理任意类型的参数
	for i, arg := range v {
		if i > 0 {
			fmt.Print(" ")
		}
		// 使用反射获取参数类型和值
		switch val := reflect.ValueOf(arg); val.Kind() {
		case reflect.Bool:
			fmt.Printf("%t", val.Bool())
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			fmt.Printf("%d", val.Int())
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			fmt.Printf("%d", val.Uint())
		case reflect.Float32, reflect.Float64:
			fmt.Printf("%v", val.Float())
		case reflect.String:
			fmt.Printf("%s", val.String())
		case reflect.Slice, reflect.Array:
			fmt.Printf("%v", val.Interface())
		case reflect.Map:
			fmt.Printf("%v", val.Interface())
		case reflect.Struct:
			fmt.Printf("%+v", val.Interface())
		default:
			fmt.Printf("%v", val.Interface())
		}
	}
	fmt.Println()
}

// Debug 打印一条Debug级别的日志
func (l *Logger) Debug(v ...interface{}) {
	l.Log(DEBUG, v...)
}

// Info 打印一条Info级别的日志
func (l *Logger) Info(v ...interface{}) {
	l.Log(INFO, v...)
}

// Warning 打印一条Warning级别的日志
func (l *Logger) Warning(v ...interface{}) {
	l.Log(WARNING, v...)
}

// Error 打印一条Error级别的日志
func (l *Logger) Error(v ...interface{}) {
	l.Log(ERROR, v...)
}

// Fatal 打印一条Fatal级别的日志，并退出程序
func (l *Logger) Fatal(v ...interface{}) {
	l.Log(FATAL, v...)
	os.Exit(1) // 立即退出程序
}
