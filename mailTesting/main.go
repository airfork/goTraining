package main

import (
	"bufio"
	"fmt"
	"os"
)

// Send text via email
/** This part is written by me
    the code in mail.go was taken
    from https://github.com/tangingw/go_smtp
**/

/** This sends via gmail, not sure if you can send via
    a different email provider. When inputing your password
    You need to get an app password, NOT the password for your email account
    Putting in your email account password will result in an error
**/

func main() {
	var c int
	var m string
	carriers := []string{"@txt.att.net", "@tmomail.net", "@vtext.com", "@messaging.sprintpcs.com", "number@vmobl.com"}
	sender := NewSender("<YOUR EMAIL>", "<YOUR PASSWORD>")

	if len(os.Args) != 2 {
		fmt.Println("Please put only the number you want to send to as a parameter.")
		os.Exit(1)
	}
	if len(os.Args[1]) != 10 {
		fmt.Println("Please provide a valid ten digit U.S. number.")
		os.Exit(2)
	}
	num := os.Args[1]

	fmt.Println("Input the number corresponding to your carrier:\n\t0: AT&T\n\t1: T-Mobile\n\t2: Verizon\n\t3: Sprint\n\t4: Virgin Mobile")

	_, err := fmt.Scan(&c)
	if err != nil {
		fmt.Println("Please type in a valid number")
		os.Exit(3)
	}
	if c > 4 || c < 0 {
		fmt.Println("Please pick a number corresponding to a carrier")
		os.Exit(4)
	}

	fmt.Printf("Please type in your message (160 characters):")
	r := bufio.NewReader(os.Stdin)
	m, err = r.ReadString('\n')
	if err != nil {
		fmt.Println("I'm not sure what went wrong...sorry about that :(")
		os.Exit(5)
	}
	if len(m) > 160 {
		fmt.Println("Please keep your message below 160 characters.")
		os.Exit(6)
	}

	receiver := []string{num + carriers[c]}

	bodyMessage := sender.WriteText(receiver, m)

	sender.SendText(receiver, bodyMessage)
}
