package main

import "flag"

import "os"

var (
	debugOut   = os.Stdout
	infoOut    = os.Stdout
	warningOut = os.Stdout
	errorOut   = os.Stdout

	// Address is the address that we will try to read the
	// logs from.
	Address string
	// Port is the SSH port of the remote server
	Port string
	// User is the user used to connect to the remote server
	User string
	// KeyPath is the path of the keyfile used for authentication
	KeyPath string
	// Password is the password used for authentication
	Password string
	// LogPath is the path of the log file on the remote server
	LogPath string
	// Listen is the port that the WebKitten server listens on
	Listen string
)

func main() {
	InitLog(debugOut, infoOut, warningOut, errorOut)

	Info.Println("WebKitten initializing...")

	flag.StringVar(&Address, "address", "10.49.4.2", "Address of the remote server")
	flag.StringVar(&Port, "port", "22", "SSH port of teh remote server")
	flag.StringVar(&User, "user", "lvuser", "User used to connect to the remote server")
	flag.StringVar(&KeyPath, "key", "", "Path of the keyfile used for authentication")
	flag.StringVar(&Password, "password", "", "Password used for authentication")
	flag.StringVar(&LogPath, "log", "/home/lvuser/logs/recent.log", "Path of the log file on the remote server")
	flag.StringVar(&Listen, "listen", "8080", "Port that the WebKitten server listens on")

	flag.Parse()

	Info.Println("Initializing remote connection...")
	session, sessionErr := NewSession(Address, Port, User, Password, KeyPath)
	if sessionErr != nil {
		panic(sessionErr)
	}
	Info.Println("Connection established.")

	Info.Println("Initialized Logger")
	remoteLog := NewRemoteLog(LogPath)
	Info.Println("Initialized Remote Log")
	go remoteLog.StartRead(session)
	Info.Println("Started Remote Read")
	Info.Println(<-remoteLog.LiveLog)
	hub := NewHub()
	hub.addInput(remoteLog.LiveLog)
	hub.merge()
	go hub.run()
	StartServer(Listen, hub)
}
