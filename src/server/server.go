package server

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net"
	"protector"
	"regexp"
	"time"
	"util"
)

type Session struct {
	In        *bufio.Reader
	Out       *bufio.Writer
	Protector protector.SessionProtector
}

func Serve(port int, maxConnections int) {
	listener, _ := net.Listen("tcp", fmt.Sprintf(":%d", port))
	currentConnections := 0

	for {
		if currentConnections < maxConnections {
			connection, _ := listener.Accept()
			currentConnections++
			log.Printf("New connection: %s. Number of current connections: %d", connection.RemoteAddr(), currentConnections)
			go func() {
				connection := connection
				defer closeConnection(connection, &currentConnections)
				session := handShake(connection)
				if session != nil {
					maintainEchoService(*session)
				}
			}()
		} else {
			time.Sleep(2 * time.Second)
		}
	}

}

func ParseMessage(message string) (key *string, resultMessage *string, err error) {
	re := regexp.MustCompile(`KEY:(\d+) MSG:(.*)`)
	matches := re.FindStringSubmatch(message)
	if matches == nil {
		return nil, nil, errors.New("Bad message")
	}
	key = &matches[1]
	resultMessage = &matches[2]
	return key, resultMessage, nil
}

func handShake(connection net.Conn) *Session {
	in := bufio.NewReader(connection)
	out := bufio.NewWriter(connection)
	message, err := util.ReadMessage(in)
	if err != nil {
		return nil
	}
	re := regexp.MustCompile(`HASH:(\d+) KEY:(\d+)`)
	matches := re.FindStringSubmatch(*message)
	if matches == nil {
		util.SendMessage(out, "ERROR: BADHANDSHAKE DETAILS: SYNTAX 'HASH:<DIGITAL_HASH> KEY:<DIGITAL_KEY>'\n")
		return nil
	}
	hash := matches[1]
	initialKey := matches[2]
	log.Printf("Got hash '%s' and initial key '%s'", hash, initialKey)
	protector := protector.SessionProtector{HashString: hash}
	key, _ := protector.NextSessionKey(initialKey)
	util.SendMessage(out, fmt.Sprintf("KEY:%s", *key))
	return &Session{
		In:        in,
		Out:       out,
		Protector: protector,
	}
}

func maintainEchoService(session Session) {
	for {
		message, err := util.ReadMessage(session.In)
		if err != nil {
			break
		}
		key, message, err := ParseMessage(*message)
		if err != nil {
			util.SendMessage(session.Out, "ERROR: BADMSG DETAILS: SYNTAX 'KEY:<DIGITAL_KEY> MSG:<YOUR MESSAGE>'\n")
			break
		}
		newKey, _ := session.Protector.NextSessionKey(*key)
		util.SendMessage(session.Out, fmt.Sprintf("KEY:%s MSG:%s", *newKey, *message))
	}
}

func closeConnection(connection net.Conn, currentConnections *int) {
	log.Printf("Closing connection: %s", connection.RemoteAddr())
	connection.Close()
	*currentConnections--
	log.Printf("Closed connection: %s. Number of current connections: %d", connection.RemoteAddr(), *currentConnections)
}
