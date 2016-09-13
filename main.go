package main

import "os"

var (
	debugOut   = os.Stdout
	infoOut    = os.Stdout
	warningOut = os.Stdout
	errorOut   = os.Stdout

	// RemoteHost is the host that we will try to read the
	// logs from.
	RemoteHost = "localhost"
	RemotePath = "/home/lvuser/logs/recent.log"
)

func main() {
	InitLog(debugOut, infoOut, warningOut, errorOut)
	Info.Println("Initialized Logger")
	remoteLog := NewRemoteLog(RemotePath, RemoteHost)
	Info.Println("Initialized Remote Log")
	go remoteLog.StartRead()
	Info.Println("Started Remote Read")
	Info.Println(<-remoteLog.LiveLog)
}
