package mail

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"github.com/ProjectOrangeJuice/mail-filter/filter"
	"github.com/emersion/go-imap"
	move "github.com/emersion/go-imap-move"
	"github.com/emersion/go-imap/client"
)

var lastID uint32

func CheckInbox() {
	// restore last uid
	readUID()
	log.Println("Connecting to server...")

	// Connect to server
	c, err := client.DialTLS("mail.harriso.co.uk:993", nil)

	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected")

	// Don't forget to logout
	defer c.Logout()

	// Login
	if err := c.Login(readAddress(), readPass()); err != nil {
		log.Fatal(err)
	}
	log.Println("Logged in")

	// Select INBOX
	mbox, err := c.Select("INBOX", false)
	if err != nil {
		log.Fatal(err)
	}

	from := lastID
	if lastID == 0 {
		from = uint32(1)
	}
	to := mbox.Messages
	if mbox.Messages > 50 { // only check the last 50 messages
		// get last 50 messages
		if mbox.Messages > 49 {
			to = from + 50
		}
	}

	seqset := new(imap.SeqSet)
	fmt.Printf("from %v, to %v\n", from, to)
	seqset.AddRange(from, to)

	messages := make(chan *imap.Message, 10)
	done := make(chan error, 1)
	go func() {
		done <- c.Fetch(seqset, []imap.FetchItem{imap.FetchEnvelope}, messages)
	}()
	for msg := range messages {
		fmt.Printf("Reading message %s with id %v\n", msg.Envelope.Subject, msg.SeqNum)
		// For each of the messages, take a look at the from and add it to the whitelist
		for _, f := range msg.Envelope.From {
			filter.Add(f.HostName)
		}
		fmt.Printf("setting last id as.. %v\n", msg.SeqNum)
		lastID = msg.SeqNum
	}

	if err := <-done; err != nil {
		log.Fatal(err)
	}
	saveLastUID()
	shakeSpam(c)
	log.Println("Done!")
}

func shakeSpam(c *client.Client) {
	c.SetDebug(os.Stdout)
	fmt.Println("** shakin ** ")
	// Select INBOX
	mbox, err := c.Select("INBOX.possible spam", false)
	if err != nil {
		log.Fatal(err)
	}

	// scan all messages
	from := uint32(1)
	to := mbox.Messages

	seqset := new(imap.SeqSet)
	fmt.Printf("from %v, to %v\n", from, to)
	seqset.AddRange(from, to)

	messages := make(chan *imap.Message, to)

	err = c.Fetch(seqset, []imap.FetchItem{imap.FetchEnvelope}, messages)
	if err != nil {
		log.Fatal(err)
	}
	var toMove []uint32
	for msg := range messages {
		fmt.Printf("Reading message %s\n", msg.Envelope.Subject)

		for _, f := range msg.Envelope.From {
			if filter.IsSafe(f.HostName) {
				toMove = append(toMove, msg.SeqNum)
				break
			}

		}

		fmt.Println("Next message")
	}

	fmt.Printf("These are to be moved.. %v\n", toMove)
	if len(toMove) > 0 {
		mover := move.NewClient(c)
		msgV := new(imap.SeqSet)
		for _, s := range toMove {
			msgV.AddNum(s)
		}
		fmt.Printf("Nums - %+v\n", msgV)

		err := mover.MoveWithFallback(msgV, "INBOX")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Moved :D")
	}

	fmt.Println("Shake done")

}

func saveLastUID() {
	fmt.Printf("saving %v\n", lastID)
	f, err := os.OpenFile("last.txt", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()
	if _, err := f.WriteString(fmt.Sprintf("%d", lastID)); err != nil {
		log.Println(err)
	}
}

func readUID() {
	dat, err := ioutil.ReadFile("last.txt")
	if err != nil {
		log.Printf("Couldn't read last uid, %s", err)
		return
	}
	last, err := strconv.ParseUint(string(dat), 10, 32)
	if err != nil {
		log.Printf("Couldn't read uid, %s", err)
		return
	}
	lastID = uint32(last)
	fmt.Printf("Last uid restored as %d\n", lastID)
}
