package localEmail

import (
	"fmt"
	"io"
	"io/ioutil"
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
	mailboxes := make(chan *imap.MailboxInfo, number)
	mailBoxeDone := make(chan error, 1)
	go func() {
		mailBoxeDone <- c.List("", "*", mailboxes)
	}()
	for box := range mailboxes {
		if box.Name != "INBOX" {
			continue
		}
		//fmt.Println("切换目录:", box.Name)
		mbox, err := c.Select(box.Name, false)
		// 选择收件箱
		if err != nil {
			//fmt.Println("select inbox err: ", err)
			continue
		}
		if mbox.Messages == 0 {
			continue
		}

		// 选择收取邮件的时间段
		criteria := imap.NewSearchCriteria()
		// 收取7天之内的邮件
		criteria.Since = time.Now().Add(-7 * time.Hour * 24)
		// 按条件查询邮件
		ids, err := c.UidSearch(criteria)
		//fmt.Println("邮件数：", len(ids))
		if err != nil || len(ids) == 0 {
			continue
		}
		seqset := new(imap.SeqSet)
		seqset.AddNum(ids...)
		sect := &imap.BodySectionName{Peek: true}

		messages := make(chan *imap.Message, 100)
		messageDone := make(chan error, 1)

		go func() {
			messageDone <- c.UidFetch(seqset, []imap.FetchItem{sect.FetchItem()}, messages)
		}()
		for msg := range messages {
			r := msg.GetBody(sect)
			mr, err := mail.CreateReader(r)
			if err != nil {
				//fmt.Println(err)
				continue
			}
			header := mr.Header

			tmp := ImapEmail{}
			if date, err := header.Date(); err == nil {
				tmp.TimeStamp = date.Unix()
				//log.Println("邮件时间 Date:", date)
			}
			if from, err := header.AddressList("From"); err == nil {
				if len(from) > 0 {
					tmp.To = from[0].Address
				}
			}
			if to, err := header.AddressList("To"); err == nil {
				if len(to) > 0 {
					tmp.To = to[0].Address
				}
			}
			if subject, err := header.Subject(); err == nil {
				tmp.Subject = subject
				//log.Println("邮件主题 Subject:", subject)
			}
			body, _ := parseEmail(mr)
			tmp.Body = string(body)
			//fmt.Println("邮件正文:", string(body))
			res = append(res, tmp)
			//for k, _ := range fileName {
			//	fmt.Println("收取到附件:", k)
			//}
		}
		return res
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
