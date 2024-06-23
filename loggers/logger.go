/*
 *@author  chengkenli
 *@project logger
 *@package logger
 *@file    logger
 *@date    2024/6/23 16:51
 */

package loggers

import (
    "github.com/natefinch/lumberjack"
    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
    "os"
)

type LoggersParms struct {
    LogPath      string // logPath 日志文件路径
    LogLevel     string // logLevel 日志级别 debug/info/warn/err
    MaxSize      int    // maxSize 单个文件大小,MB
    MaxBackups   int    // maxBackups 保存的文件个数
    MaxAge       int    // maxAge 保存的天数， 没有的话不删除
    Compress     bool   // compress 压缩
    JsonFormat   bool   // jsonFormat 是否输出为json格式
    ShowLine     bool   // shoowLine 显示代码行
    LogInConsole bool   // logInConsole 是否同时输出到控制台
}

func Loggers(logger LoggersParms) *zap.Logger {
    _, err := os.OpenFile(logger.LogPath, os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
    if err != nil {
        return nil
    }

    hook := lumberjack.Logger{
        Filename:   logger.LogPath,    // 日志文件路径
        MaxSize:    logger.MaxSize,    // megabytes
        MaxBackups: logger.MaxBackups, // 最多保留300个备份
        Compress:   logger.Compress,   // 是否压缩 disabled by default
        MaxAge:     logger.MaxAge,     // maxAge 保存的天数， 没有的话不删除
    }

    var syncer zapcore.WriteSyncer
    if logger.LogInConsole {
        syncer = zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&hook))
    } else {
        syncer = zapcore.AddSync(&hook)
    }

    encoderConfig := zapcore.EncoderConfig{
        TimeKey:        "time",
        LevelKey:       "level",
        NameKey:        "logger",
        CallerKey:      "line",
        MessageKey:     "msg",
        StacktraceKey:  "stacktrace",
        LineEnding:     zapcore.DefaultLineEnding,
        EncodeLevel:    zapcore.LowercaseLevelEncoder,  // 小写编码器
        EncodeTime:     zapcore.ISO8601TimeEncoder,     // ISO8601 UTC 时间格式
        EncodeDuration: zapcore.SecondsDurationEncoder, //
        EncodeCaller:   zapcore.ShortCallerEncoder,     // 全路径编码器
        EncodeName:     zapcore.FullNameEncoder,
    }

    var encoder zapcore.Encoder
    if logger.JsonFormat {
        encoder = zapcore.NewJSONEncoder(encoderConfig)
    } else {
        encoder = zapcore.NewConsoleEncoder(encoderConfig)
    }

    // 设置日志级别,debug可以打印出info,debug,warn；info级别可以打印warn，info；warn只能打印warn
    // debug->info->warn->error->fatal
    var level zapcore.Level
    switch logger.LogLevel {
    case "debug":
        level = zap.DebugLevel
    case "info":
        level = zap.InfoLevel
    case "warn":
        level = zap.WarnLevel
    case "error":
        level = zap.ErrorLevel
    case "fatal":
        level = zap.FatalLevel
    case "panic":
        level = zap.PanicLevel
    default:
        level = zap.InfoLevel
    }

    core := zapcore.NewCore(
        encoder,
        syncer,
        level,
    )

    logg := zap.New(core)
    if logger.ShowLine {
        logg = logg.WithOptions(zap.AddCaller())
    }
    return logg
}
