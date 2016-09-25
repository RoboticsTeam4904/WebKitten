package main

import (
	"bufio"

	"golang.org/x/crypto/ssh"
)

type RemoteLog struct {
	LogPath string
	LiveLog chan string
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
		remote.LiveLog <- scanner.Text()
	}
	scannerErr := scanner.Err()
	if scannerErr != nil {
		close(remote.LiveLog)
		Error.Println(scannerErr.Error())
	}
}
