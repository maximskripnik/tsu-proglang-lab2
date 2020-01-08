package client

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"protector"
	"regexp"
	"server"
	"time"
	"util"
)

func Connect(host string, port int) {
	connection, _ := net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
	session, key := handShake(connection)
	reader := bufio.NewReader(os.Stdin)
	var newKey *string
	for {
		fmt.Print("$")
		text, _ := reader.ReadString('\n')
		log.Printf("TEXT %s", text)
		text = text[:len(text)-1]
		newKey, _ = session.Protector.NextSessionKey(key)
		util.SendMessage(session.Out, fmt.Sprintf("KEY:%s MSG:%s", key, text))
		response, _ := util.ReadMessage(session.In)
		receivedKey, _, _ := server.ParseMessage(*response)
		newKey, _ = session.Protector.NextSessionKey(key)
		key = *newKey
		validateKey(key, *receivedKey)
	}
}

func handShake(connection net.Conn) (session server.Session, key string) {
	in := bufio.NewReader(connection)
	out := bufio.NewWriter(connection)
	source := rand.NewSource(time.Now().Unix())
	hash := protector.GetHashStr(source)
	initialKey := protector.GetSessionKey(source)
	log.Printf("Generated hash '%s' and initial key '%s'", hash, initialKey)
	protector := protector.SessionProtector{HashString: hash}
	util.SendMessage(out, fmt.Sprintf("HASH:%s KEY:%s", hash, initialKey))
	re := regexp.MustCompile(`KEY:(\d+)`)
	response, _ := util.ReadMessage(in)
	receivedKey := re.FindStringSubmatch(*response)[1]
	newKey, _ := protector.NextSessionKey(initialKey)
	validateKey(*newKey, receivedKey)
	return server.Session{
		In:        in,
		Out:       out,
		Protector: protector,
	}, *newKey
}

func validateKey(expectedKey string, receivedKey string) {
	log.Printf("Validating received key '%s' against expected key '%s'", receivedKey, expectedKey)
	if receivedKey == expectedKey {
		log.Print("Ok. Key is correct")
	} else {
		log.Fatalf("!BAD KEY!")
	}
}
