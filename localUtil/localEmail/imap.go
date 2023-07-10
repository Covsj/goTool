package localEmail

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"time"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-message/mail"
)

func LoginImapEmail(UserName, Password string) (*client.Client, error) {
	server := "imap.mail.ru:993"
	dial := new(net.Dialer)
	dial.Timeout = time.Duration(3) * time.Second
	c, err := client.DialWithDialerTLS(dial, server, nil)
	if err != nil {
		c, err = client.DialWithDialer(dial, server) // 非加密登录
	}
	if err != nil {
		return nil, err
	}
	//登陆
	if err = c.Login(UserName, Password); err != nil {
		return nil, err
	}
	return c, nil
}

func GetImapEmailMessage(c *client.Client, number int) []ImapEmail {
	res := []ImapEmail{}
	if number == 0 {
		number = 10
	}

	mbox, err := c.Select("INBOX", false)
	if err != nil {
		fmt.Println(err)
		return res
	}
	if mbox.Messages == 0 {
		return res
	}

	from := uint32(1)
	to := mbox.Messages
	if mbox.Messages > uint32(number) {
		from = mbox.Messages - uint32(number)
	}
	for i := from; i <= to; i++ {
		seqset := new(imap.SeqSet)
		seqset.AddNum(i)
		messages := make(chan *imap.Message, 1)
		done := make(chan error, 1)
		go func() {
			done <- c.Fetch(seqset, []imap.FetchItem{imap.FetchEnvelope, imap.FetchRFC822}, messages)
		}()

		for msg := range messages {
			tmp := ImapEmail{}
			tmp.TimeStamp = msg.Envelope.Date.Unix()
			if len(msg.Envelope.From) > 0 {
				from := msg.Envelope.From[0]
				tmp.From = from.Address()
			}
			if len(msg.Envelope.To) > 0 {
				To := msg.Envelope.To[0]
				tmp.From = To.Address()
			}
			tmp.Subject = msg.Envelope.Subject
			if body := msg.GetBody(&imap.BodySectionName{Peek: true}); body != nil {
				bytes, _ := ioutil.ReadAll(body)
				tmp.Body = string(bytes)
			}
			res = append(res, tmp)
		}

		if err := <-done; err != nil {
			log.Println("Failed to fetch message:", err)
			continue
		}
	}
	return res
}

func parseEmail(mr *mail.Reader) (body []byte, fileMap map[string][]byte) {
	for {
		p, err := mr.NextPart()
		if err == io.EOF {
			break
		} else if err != nil {
			break
		}
		if p != nil {
			switch h := p.Header.(type) {
			case *mail.InlineHeader:
				body, err = ioutil.ReadAll(p.Body)
				if err != nil {
					fmt.Println("read body err:", err.Error())
				}

			case *mail.AttachmentHeader:
				fileName, _ := h.Filename()
				fileContent, _ := ioutil.ReadAll(p.Body)
				fileMap[fileName] = fileContent
			}
		}
	}
	return
}
