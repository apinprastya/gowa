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

func (c *Cb) OnQRChangedEvent(qrCode string) {
	fmt.Println("NEW QR CODE: ", time.Now(), qrCode)
}

func (c *Cb) OnLogoutEvent() {
	fmt.Println("LOGOUT EVENT: ", time.Now())
}

func (c *Cb) OnReady() {
	fmt.Println("READY EVENT: ", time.Now())
	time.Sleep(10 * time.Second)
	data, err := ReadFileToBase64("/home/apin/billingv2.pdf")
	if err != nil {
		fmt.Println(err)
		return
	}
	id, err := c.client.SendMessage("@c.us", &gowa.MessageMedia{MimeType: "application/pdf",
		Filename: "MyFilename.pdf", Caption: "This is the caption bosku", Data: data})
	fmt.Println(id, err)
}

func main() {
	cb := &Cb{}
	client := gowa.New(gowa.BrowserTypeFirefox, cb)
	cb.client = client
	err := client.Init()
	if err != nil {
		panic(err)
	}

	/*code, err := client.RequestPairingCode("")
	if err != nil {
		panic(err)
	}
	fmt.Println("PAIRING CODE: ", code)*/

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	fmt.Println("closing client")
	client.Close()
}
