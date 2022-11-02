package main

import (
	"bufio"
	"fmt"
	_ "github.com/emersion/go-message/charset"
	"github.com/emersion/go-message/mail"
	"io"
	"mime"
	"os"
	"strings"
)

func main() {
	var r *mail.Reader

	if len(os.Args) > 1 {
		path := os.Args[1]
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
	fields := [4]string{"From", "To", "Cc", "Bcc"}
	for _, field := range fields {
		addrlist := []string{}
		header, err := h.AddressList(field)
		if err != nil {
			continue
		}
		for _, address := range header {

			dec := new(mime.WordDecoder)
			name, err := dec.DecodeHeader(address.Name)
			if err != nil {
				continue
			}
			addrlist = append(addrlist, name+" "+"<"+address.Address+">")
		}
		if len(addrlist) > 0 {
			fmt.Print(field + ": ")
			for i := 0; i < len(addrlist); i++ {
				if i == len(addrlist)-1 {
					fmt.Print(addrlist[i] + "\n")
				} else {
					fmt.Print(addrlist[i] + ", ")
				}
			}
		}
	}
	time, err := h.Date()
	if err == nil {
		fmt.Println("Date: ", time)
	}
	subject, err := h.Subject()
	if err == nil {

		fmt.Println("Subject: ", subject)
	}
	fmt.Print("\n--------------------------------------------------------\n\n")
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
