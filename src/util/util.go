package util

import (
	"bufio"
	"log"
)

func SendMessage(out *bufio.Writer, message string) {
	out.Write([]byte(message + "\n"))
	out.Flush()
	log.Printf("Sent message: '%s'", message)
}

func ReadMessage(in *bufio.Reader) (*string, error) {
	message, err := in.ReadString('\n')
	if err == nil {
		log.Printf("Recieved message: '%s'", message[:len(message)-1])
		return &message, nil
	}
	return nil, err
}
