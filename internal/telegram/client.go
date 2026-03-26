package telegram

import (
	"context"
	"fmt"
	"time"

	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/auth/qrlogin"
	"github.com/gotd/td/telegram/message"
	"github.com/gotd/td/tg"

	"github.com/Bimos6/telegram-service/internal/models"
	"github.com/Bimos6/telegram-service/pkg/logger"
)

type Client struct {
	id       string
	appID    int
	appHash  string
	client   *telegram.Client
	api      *tg.Client
	sender   *message.Sender
	messages chan models.Message
	log      logger.Logger
	cancel   context.CancelFunc
}

func NewClient(id string, appID int, appHash string, log logger.Logger) *Client {
	return &Client{
		id:       id,
		appID:    appID,
		appHash:  appHash,
		messages: make(chan models.Message, 100),
		log:      log.WithField("component", "telegram"),
	}
}

func (c *Client) Connect(ctx context.Context) error {
	sessionPath := fmt.Sprintf("./sessions/%s", c.id)

	c.client = telegram.NewClient(
		c.appID,
		c.appHash,
		telegram.Options{
			SessionStorage: &telegram.FileSessionStorage{
				Path: sessionPath,
			},
		},
	)

	c.api = c.client.API()
	c.sender = message.NewSender(c.api)

	return nil
}

func (c *Client) GetQRCode(ctx context.Context) (string, error) {
	qr := c.client.QR()

	onQR := func(ctx context.Context, token qrlogin.Token) error {
		return nil
	}

	go func() {
		var provider qrlogin.LoggedIn
		_, err := qr.Auth(
			ctx,
			provider,
			onQR,
		)
		if err != nil {
			c.log.WithField("error", err).Error("Authentication failed")
		} else {
			c.log.Info("Authentication successful")
		}
	}()

	return "tg://login?token=FAKE_TOKEN_FOR_DEMO", nil
}

func (c *Client) SendMessage(ctx context.Context, peer, text string) (int64, error) {
	var inputPeer tg.InputPeerClass

	if peer == "@me" {
		inputPeer = &tg.InputPeerSelf{}
	} else if len(peer) > 0 && peer[0] == '@' {
		username := peer[1:]

		resolved, err := c.api.ContactsResolveUsername(ctx, &tg.ContactsResolveUsernameRequest{
			Username: username,
		})
		if err != nil {
			return 0, err
		}

		if len(resolved.Users) > 0 {
			if user, ok := resolved.Users[0].(*tg.User); ok {
				inputPeer = &tg.InputPeerUser{
					UserID:     user.ID,
					AccessHash: user.AccessHash,
				}
			}
		} else {
			return 0, fmt.Errorf("user %s not found", username)
		}
	} else {
		return 0, fmt.Errorf("unsupported peer format: %s", peer)
	}

	updates, err := c.sender.To(inputPeer).Text(ctx, text)
	if err != nil {
		return 0, err
	}

	var msgID int64

	switch u := updates.(type) {
	case *tg.Updates:
		for _, update := range u.Updates {
			if msg, ok := update.(*tg.UpdateNewMessage); ok {
				if m, ok := msg.Message.(*tg.Message); ok {
					msgID = int64(m.ID)
					break
				}
			}
		}
	case *tg.UpdateShort:
		if msg, ok := u.Update.(*tg.UpdateNewMessage); ok {
			if m, ok := msg.Message.(*tg.Message); ok {
				msgID = int64(m.ID)
			}
		}
	}

	if msgID == 0 {
		msgID = time.Now().UnixNano()
	}

	return msgID, nil
}

func (c *Client) Messages() <-chan models.Message {
	return c.messages
}

func (c *Client) Disconnect() error {
	if c.cancel != nil {
		c.cancel()
	}
	close(c.messages)
	return nil
}
