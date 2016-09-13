package main

import (
	"bufio"
	"os/exec"
)

type RemoteLog struct {
	LogPath string
	Host    string
	LiveLog chan string
}

func NewRemoteLog(logPath string, host string) RemoteLog {
	return RemoteLog{
		LogPath: logPath,
		Host:    host,
		LiveLog: make(chan string),
	}
}

func (remote *RemoteLog) StartRead() {
	cmd := exec.Command("ssh", "-t", remote.Host, "tail", "-f", remote.LogPath)
	stdout, stdoutErr := cmd.StdoutPipe()
	if stdoutErr != nil {
		Error.Println(stdoutErr)
	}
	startErr := cmd.Start()
	if startErr != nil {
		Error.Println(startErr)
	}
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		remote.LiveLog <- scanner.Scan()
	}
	scannerErr := scanner.Err()
	if scannerErr != nil {
		Error.Println(scannerErr)
	}
}
