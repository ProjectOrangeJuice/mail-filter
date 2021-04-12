package mail

import (
	"fmt"
	"log"

	"github.com/ProjectOrangeJuice/mail-filter/filter"
	"github.com/emersion/go-imap"
	move "github.com/emersion/go-imap-move"
	"github.com/emersion/go-imap/client"
)

func CheckInbox() {
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

	// Get the last 4 messages
	from := uint32(1)
	to := mbox.Messages
	if mbox.Messages > 50 {
		// We're using unsigned integers here, only subtract if the result is > 0
		from = mbox.Messages - 50
	}
	seqset := new(imap.SeqSet)
	seqset.AddRange(from, to)

	messages := make(chan *imap.Message, 10)
	done := make(chan error, 1)
	go func() {
		done <- c.Fetch(seqset, []imap.FetchItem{imap.FetchEnvelope}, messages)
	}()
	for msg := range messages {
		// For each of the messages, take a look at the from and add it to the whitelist

		for _, f := range msg.Envelope.From {
			if filter.IsBlacklist(f.HostName) {
				// Move to spam!
				break
			}
			filter.Add(f.HostName)
		}

	}

	if err := <-done; err != nil {
		log.Fatal(err)
	}
	shakeSpam(c)
	log.Println("Done!")
}

func shakeSpam(c *client.Client) {
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
	re := 0
	for msg := range messages {
		re++
		fmt.Printf("Reading message %s ", msg.Envelope.Subject)

		for _, f := range msg.Envelope.From {
			fmt.Printf("with a host of %s, safe? %v\n", f.HostName, filter.IsSafe(f.HostName))
			if filter.IsSafe(f.HostName) {
				toMove = append(toMove, msg.SeqNum)
				break
			}

		}
	}
	fmt.Printf("I was going to read %d but I actually read %d\n", to, re)

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
