package main

import (
	"bufio"
	"fmt"
	"github.com/google/logger"
	"io"
	"net"
	"os"
	"strings"
)

const logPath = "./dynatrace-listener.log"

func handleMessages(conn net.Conn) {
	rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
	defer conn.Close()

	for {
		msg, err := rw.ReadString('\n')

		switch {
		case err == io.EOF:
			logger.Infof("Reached EOF - close this connection.")
			return
		case err != nil:
			logger.Errorf("Error reading message. Got: %s: %s", msg, err)
			return
		}

		msg = strings.Trim(msg, "\n ")
		logger.Infof("Message: %s", msg)

	}

}

func Listen() error {
	var err error

	addr := fmt.Sprintf(":%s", os.Args[1])
	listener, err := net.Listen("tcp", addr)

	if err != nil {
		logger.Fatalf("Could not listen on %s.", addr)
	}
	logger.Infof("Listening on %s", listener.Addr().String())

	for {
		logger.Infof("Waiting for a connection request.")
		conn, err := listener.Accept()
		if err != nil {
			logger.Infof("Failed accepting a connection request:", err)
			continue
		}
		logger.Infof("Handle incoming messages.")
		go handleMessages(conn)

	}
}

func main() {

	lf, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
	defer lf.Close()

	loggerOne := logger.Init("LoggerDynatrace", true, false, lf)

	defer loggerOne.Close()
	if err != nil {
		logger.Fatalf("Failed to open log file: %v", err)
	}

	err = Listen()
	if err != nil {
		logger.Errorf("Error: %s", err)
	}
}
