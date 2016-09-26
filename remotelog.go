package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"strings"

	"golang.org/x/crypto/ssh"
)

type RemoteLog struct {
	LogPath string
	LiveLog chan string
}

// 2016-11-16_11:46:20 WTF: getStuff: blah
type LogItem struct {
	Timestamp string `json:"timestamp"`
	LogLevel  string `json:"loglevel"`
	Source    string `json:"source"`
	Message   string `json:"message"`
}

func NewRemoteLog(logPath string) RemoteLog {
	return RemoteLog{
		LogPath: logPath,
		LiveLog: make(chan string),
	}
}

func (remote *RemoteLog) StartRead(session *ssh.Session) {
	stdout, stdoutErr := session.StdoutPipe()
	if stdoutErr != nil {
		Error.Println(stdoutErr)
	}
	startErr := session.Start("tail -f " + remote.LogPath)
	if startErr != nil {
		Error.Println(stdoutErr)
	}
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		logItem, logItemErr := ParseLogItem(scanner.Text())
		if logItemErr != nil {
			Error.Println(logItemErr)
			continue
		}
		marshal, marshalErr := json.Marshal(logItem)
		if marshalErr != nil {
			Error.Println(marshalErr)
			continue
		}
		remote.LiveLog <- string(marshal)
	}
	scannerErr := scanner.Err()
	if scannerErr != nil {
		close(remote.LiveLog)
		Error.Println(scannerErr.Error())
	}
}

func ParseLogItem(rawitem string) (LogItem, error) {
	logItem := LogItem{}
	parts := strings.SplitN(rawitem, " ", 4)
	if len(parts) < 4 {
		return &LogItem{}, errors.New("logitem: couldn't parse log item: malformatted input \"" + rawitem + "\"")
	}
	logItem.Timestamp = parts[0]
	logItem.LogLevel = strings.TrimSuffix(parts[1])
	logItem.Source = strings.TrimSuffix(parts[2])
	logItem.Message = parts[3]
	return logItem, nil
}
