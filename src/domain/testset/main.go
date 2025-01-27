package main

import (
	"log"
	"net/smtp"
)

func main() {
	send("hello bro, This just test.")
}

func send(body string) {
	from := "noellimxaws@gmail.com"
	pass := ""
	to := "noellimx@gmail.com"

	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: Hello there\n\n" +
		body

	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
		from, []string{to}, []byte(msg))

	if err != nil {
		log.Printf("smtp error: %s", err)
		return
	}
	log.Println("Successfully sended to " + to)
}
