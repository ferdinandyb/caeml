package main

import (
	"bufio"
	"fmt"
	"io"
	"mime"
	"os"
	"strings"

	"git.sr.ht/~sircmpwn/getopt"
	_ "github.com/emersion/go-message/charset"
	"github.com/emersion/go-message/mail"
)

func main() {
	var r *mail.Reader

	args := make([]string, len(os.Args))
	copy(args, os.Args)

	opts, optind, err := getopt.Getopts(args, "H:O")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	headers := []string{"From", "To", "Cc", "Bcc", "Date", "Subject"}
	printBody := true
	for _, opt := range opts {
		if opt.Option == 'H' {
			headers = strings.Split(opt.Value, ",")
			for i, _ := range headers {
				headers[i] = strings.TrimSpace(headers[i])
			}
		} else if opt.Option == 'O' {
			printBody = false
		}

	}

	if len(args[optind:]) > 0 {
		path := strings.Join(args[optind:], " ")
		f, err := os.Open(path)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		var readerr error
		r, readerr = mail.CreateReader(f)
		if readerr != nil {
			fmt.Println(readerr)
			f.Close()
			os.Exit(1)
		}
	} else {

		var readerr error
		reader := bufio.NewReader(os.Stdin)
		r, readerr = mail.CreateReader(reader)
		if readerr != nil {
			fmt.Println(readerr)
			os.Exit(1)
		}
	}

	h := &mail.Header{Header: r.Header.Header}

	for _, header := range headers {
		switch strings.ToLower(header) {
		case "from", "to", "cc", "bcc":
			addrlist := []string{}
			headerContent, err := h.AddressList(header)
			if err != nil {
				continue
			}
			for _, address := range headerContent {
				dec := new(mime.WordDecoder)
				name, err := dec.DecodeHeader(address.Name)
				if err != nil {
					continue
				}
				addrlist = append(addrlist, name+" "+"<"+address.Address+">")
			}
			if len(addrlist) > 0 {
				fmt.Print(header + ": ")
				for i := 0; i < len(addrlist); i++ {
					if i == len(addrlist)-1 {
						fmt.Print(addrlist[i] + "\n")
					} else {
						fmt.Print(addrlist[i] + ", ")
					}
				}
			}
		case "date":
			time, err := h.Date()
			if err == nil {
				fmt.Println("Date: ", time)
			}
		case "subject":
			subject, err := h.Subject()
			if err == nil {
				fmt.Println("Subject: ", subject)
			}
		default:
			text, err := h.Text(header)
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Println(header + ": " + text)
		}
	}
	if !printBody {
		return
	}
	fmt.Print("\n-----------------------------------------\n\n")
	for true {
		part, err := r.NextPart()
		if err != nil {
			break
		}
		contenttype := part.Header.Get("Content-Type")
		if strings.Split(contenttype, ";")[0] == "text/plain" {
			buf := new(strings.Builder)
			_, err := io.Copy(buf, part.Body)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(buf.String())
		}
	}

}
