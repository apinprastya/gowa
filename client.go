package gowa

import (
	"errors"
	"fmt"
	"strings"

	"github.com/apinprastya/gowa/jscode"
	"github.com/playwright-community/playwright-go"
)

type BrowserType int

const (
	BrowserTypeFirefox BrowserType = iota
	BrowserTypeChromium
	BrowserTypeWebKit
)

type ClientCallback interface {
	OnReady()
	OnLogoutEvent()
	OnMessageAck(messageID string, messageAck MessageAck)
}

type Client struct {
	pw          *playwright.Playwright
	page        playwright.Page
	browserType BrowserType
	callback    ClientCallback
	loggedIn    bool
}

func New(browserType BrowserType, callback ClientCallback) *Client {
	return &Client{
		browserType: browserType,
		callback:    callback,
		loggedIn:    false,
	}
}

func (c *Client) IsLoggedIn() bool {
	return c.loggedIn
}

func (c *Client) Init() error {
	pw, err := c.runOrInstall()
	if err != nil {
		return err
	}
	c.pw = pw
	var browser playwright.BrowserContext

	switch c.browserType {
	case BrowserTypeFirefox:
		browser, err = pw.Firefox.LaunchPersistentContext("/home/apin/project/go/gowa/.local", playwright.BrowserTypeLaunchPersistentContextOptions{
			Headless: playwright.Bool(false),
			Timeout:  playwright.Float(30000),
		})
		if err != nil {
			return err
		}
	case BrowserTypeChromium:
		browser, err = pw.Chromium.LaunchPersistentContext("/home/apin/project/go/gowa/.local_chromium", playwright.BrowserTypeLaunchPersistentContextOptions{
			Headless: playwright.Bool(false),
			Timeout:  playwright.Float(30000),
			Args:     []string{"--disable-blink-features=AutomationControlled"},
		})
		if err != nil {
			return err
		}
	}
	if len(browser.Pages()) > 0 {
		c.page = browser.Pages()[0]
	}
	c.page.Goto("https://web.whatsapp.com")

	err = c.page.WaitForLoadState(playwright.PageWaitForLoadStateOptions{State: playwright.LoadStateDomcontentloaded, Timeout: playwright.Float(60000)})
	if err != nil {
		return err
	}

	_, err = c.page.Evaluate(jscode.InjectAuthJS)
	if err != nil {
		return err
	}

	needLogin, err := c.needLogin()
	if err != nil {
		return err
	}
	c.loggedIn = !needLogin
	/*if needLogin {
		err = c.page.ExposeFunction("onQRChangedEvent", func(args ...interface{}) interface{} {
			if len(args) > 0 {
				if str, ok := args[0].(string); ok {
					if c.callback != nil {
						c.callback.OnQRChangedEvent(str)
					}
				}
			}
			return nil
		})
		if err != nil {
			return err
		}
		_, err = c.page.Evaluate(QrJS)
		if err != nil {
			return err
		}
	}*/

	err = c.page.ExposeFunction("onLogoutEvent", func(args ...interface{}) interface{} {
		if c.callback != nil {
			c.callback.OnLogoutEvent()
		}
		return nil
	})
	if err != nil {
		return err
	}

	err = c.page.ExposeFunction("onAuthAppStateChangedEvent", func(args ...interface{}) interface{} {
		return nil
	})
	if err != nil {
		return err
	}

	err = c.page.ExposeFunction("onAppStateHasSyncedEvent", func(args ...interface{}) interface{} {
		c.loggedIn = true
		_, err := c.page.Evaluate(jscode.StoreJS)
		if err != nil {
			return nil
		}
		_, err = c.page.Evaluate(jscode.UtilJS)
		if err != nil {
			return nil
		}

		_, err = c.page.Evaluate(jscode.MessageEventJS)
		if err != nil {
			return nil
		}

		if c.callback != nil {
			c.callback.OnReady()
		}
		return nil
	})
	if err != nil {
		return err
	}

	_, err = c.page.Evaluate(jscode.StateChangeJS)
	if err != nil {
		return err
	}

	err = c.page.ExposeFunction("onMessageAckEvent", func(args ...interface{}) interface{} {
		if len(args) >= 2 {
			fmt.Println(args)
			if id, ok := args[0].(string); ok {
				if ack, ok := args[1].(int); ok {
					if c.callback != nil {
						c.callback.OnMessageAck(id, MessageAck(ack))
					}
				}
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) RequestPairingCode(phone string) (string, error) {
	val, err := c.page.Evaluate(jscode.RequestPairCodeJS, phone)
	if err != nil {
		return "", err
	}
	if valStr, ok := val.(string); ok {
		return valStr, nil
	}
	return "", errors.New("invalid return value")
}

func (c *Client) SendMessage(phone string, message Message) (string, error) {
	result, err := c.page.Evaluate(jscode.SendMessageJS, map[string]any{
		"chatId":  c.formatPhoneNumber(phone),
		"message": message.Content(),
		"options": message.Option(),
	})
	if err != nil {
		return "", err
	}
	if id, ok := result.(string); ok {
		return id, nil
	}

	return "", errors.New("invalid return value")
}

func (c *Client) runOrInstall() (*playwright.Playwright, error) {
	pw, err := playwright.Run()
	if err != nil && strings.Contains(err.Error(), "please install the driver") {
		err = playwright.Install()
		if err != nil {
			return nil, err
		}
		pw, err = playwright.Run()
		if err != nil {
			return nil, err
		}
	}
	return pw, nil
}

func (c *Client) Close() {
	if c.page != nil {
		c.page.Close()
	}
	if c.pw != nil {
		c.pw.Stop()
	}
}

func (c *Client) needLogin() (bool, error) {
	needAuthInt, err := c.page.Evaluate(jscode.NeedAuthJS)
	if err != nil {
		return false, err
	}
	if val, ok := needAuthInt.(bool); ok {
		return val, nil
	}
	return true, nil
}

func (c *Client) formatPhoneNumber(phone string) string {
	if !strings.HasSuffix(phone, "@s.whatsapp.net") {
		phone = strings.ReplaceAll(phone, "c.us", "s.whatsapp.net")
	}
	if !strings.HasSuffix(phone, "@s.whatsapp.net") {
		phone = phone + "s.whatsapp.net"
	}
	return phone
}
