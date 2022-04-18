package logging

import (
	"IMConnection/conf"
	"fmt"
	"log"
	"os"
	"time"
)

func getLogFilePath() string {
	return fmt.Sprintf("%s", conf.AppSetting.LogSavePath)
}

func getLogFileFullPath() string {
	prefixPath := getLogFilePath()
	suffixPath := fmt.Sprintf("%s%s.%s", conf.AppSetting.LogSaveName, time.Now().Format(conf.AppSetting.TimeFormat), conf.AppSetting.LogFileExt)

	return fmt.Sprintf("%s%s", prefixPath, suffixPath)
}

func openLogFile(filePath string) *os.File {
	_, err := os.Stat(filePath)
	switch {
	case os.IsNotExist(err):
		mkDir()
	case os.IsPermission(err):
		log.Fatalf("Permission :%v", err)
	}

	handle, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Fail to OpenFile :%v", err)
	}

	return handle
}

func mkDir() {
	dir, _ := os.Getwd()
	err := os.MkdirAll(dir+"/"+getLogFilePath(), os.ModePerm)
	if err != nil {
		println(err)
	}
}

func getLogFileName() string {
	return fmt.Sprintf("%s%s%s",
		conf.AppSetting.LogSaveName,
		time.Now().Format(conf.AppSetting.TimeFormat),
		conf.AppSetting.LogFileExt,
	)
}
