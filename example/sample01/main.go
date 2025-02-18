package main

import (
	"encoding/base64"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/apinprastya/gowa"
)

func ReadFileToBase64(filePath string) (string, error) {
	// Read file contents
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	// Encode to Base64
	base64Str := base64.StdEncoding.EncodeToString(data)
	return base64Str, nil
}

type Cb struct {
	client *gowa.Client
}

func (c *Cb) OnLogoutEvent() {
	fmt.Println("got logout event: ", time.Now())
	time.Sleep(10 * time.Second)
	c.client.Close()
	time.Sleep(5 * time.Second)
	err := c.client.Init()
	if err != nil {
		fmt.Println("onLogout error", err)
	}
}

func (c *Cb) OnMessageAck(messageID string, messageAck gowa.MessageAck) {
	fmt.Println("got message ack: ", time.Now(), ", message id:", messageID, ", ack:", messageAck)
}

func (c *Cb) OnReady() {
	fmt.Println("Client is Ready: ", time.Now())
	go c.doInput()
}

func (c *Cb) OnNeedLogin() {
	go c.doLogin()
}

func (c *Cb) doLogin() {
	phoneNumber := ""
	fmt.Println("Need to connect to whatsapp web. Please input the phone number what to connect.")
	fmt.Print("Phone number (start with country code): ")
	fmt.Scanln(&phoneNumber)
	fmt.Println("Please wait, your phone will get notification from whatsapp to input the pairing code.")
	code, err := c.client.RequestPairingCode(phoneNumber)
	if err != nil {
		panic(err)
	}
	fmt.Println("Your pairing code: ", code)
}

func (c *Cb) doInput() {
	phone := ""
	action := ""
	time.Sleep(100 * time.Millisecond)
	fmt.Println("Send Whatsapp message. Please input phone number.")
	fmt.Print("Target phone number: ")
	fmt.Scanln(&phone)
	fmt.Print("Please select action:\n1. Send text message\n2. Send media file.\nAction: ")
	fmt.Scanln(&action)
	switch action {
	case "1":
		content := ""
		fmt.Print("Message: ")
		fmt.Scanln(&content)
		fmt.Println("Sending message to", phone, "with content:", content)
		messageID, err := c.client.SendMessage(phone, &gowa.MessageText{MessageContent: content})
		if err != nil {
			fmt.Println("Error send message", err)
		} else {
			fmt.Println("Message sent with ID:", messageID)
		}
	case "2":
		filepath := ""
		mimetype := ""
		caption := ""
		fmt.Print("File path: ")
		fmt.Scanln(&filepath)
		base64Str, err := ReadFileToBase64(filepath)
		if err != nil {
			fmt.Println("Error read file: ", err)
			return
		}
		fmt.Print("Mime type: ")
		fmt.Scanln(&mimetype)
		fmt.Print("Caption: ")
		fmt.Scanln(&caption)
		messageID, err := c.client.SendMessage(phone, &gowa.MessageMedia{
			MimeType: mimetype,
			Caption:  caption,
			Data:     base64Str,
		})
		if err != nil {
			fmt.Println("Error send message", err)
		} else {
			fmt.Println("Message sent with ID:", messageID)
		}
	default:
		fmt.Println("Invalid action")
	}
}

func main() {
	cb := &Cb{}
	client := gowa.New("sample", gowa.BrowserTypeFirefox, cb)
	cb.client = client
	err := client.Init()
	if err != nil {
		panic(err)
	}

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	fmt.Println("closing client")
	client.Close()
}
