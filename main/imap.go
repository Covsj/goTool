package main

import (
	"fmt"

	"github.com/Covsj/goTool/localUtil/localEmail"
)

func main() {
	c, err := localEmail.LoginImapEmail("default@covsj.top", "")
	if err != nil {
		panic(err)
	}
	defer c.Close()
	message := localEmail.GetImapEmailMessage(c, 10)
	for _, email := range message {
		fmt.Println(email.From, email.To, email.Subject, email.Body)
	}
}
