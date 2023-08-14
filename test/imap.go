package test

import (
	"fmt"

	"github.com/Covsj/goTool/email"
)

func main() {
	c, err := email.LoginImapEmail("****", "****")
	if err != nil {
		panic(err)
	}
	defer c.Close()
	message := localEmail.GetImapEmailMessage(c, 10)
	for _, email := range message {
		fmt.Println("From: ", email.From, "To: ", email.To, "Subject: ", email.Subject, "Body: ")
	}
}
