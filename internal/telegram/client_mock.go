package telegram

import (
	"context"
	"fmt"
	"time"

	"github.com/Bimos6/telegram-service/internal/models"
	"github.com/Bimos6/telegram-service/pkg/logger"
)

type MockClient struct {
	id       string
	messages chan models.Message
	log      logger.Logger
	cancel   context.CancelFunc
}

func NewMockClient(id string, log logger.Logger) *MockClient {
	return &MockClient{
		id:       id,
		messages: make(chan models.Message, 100),
		log:      log.WithField("component", "telegram_mock"),
	}
}

func (c *MockClient) Connect(ctx context.Context) error {
	c.log.Info("Mock client connecting...")

	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		msgID := int64(1)
		for {
			select {
			case <-ctx.Done():
				c.log.Debug("Message generator stopped")
				return
			case <-ticker.C:
				msg := models.Message{
					ID:        msgID,
					From:      "test_bot",
					Text:      fmt.Sprintf("Test message #%d from session %s", msgID, c.id),
					Timestamp: time.Now(),
					SessionID: c.id,
				}
				select {
				case c.messages <- msg:
					c.log.WithField("msg_id", msgID).Debug("Test message sent")
				default:
				}
				msgID++
			}
		}
	}()

	c.log.Info("Mock client connected")
	return nil
}

func (c *MockClient) GetQRCode(ctx context.Context) (string, error) {
	return "tg://login?token=mock_token_for_testing", nil
}

func (c *MockClient) SendMessage(ctx context.Context, peer, text string) (int64, error) {
	c.log.WithField("peer", peer).WithField("text", text).Info("Mock: message sent")
	return time.Now().UnixNano(), nil
}

func (c *MockClient) Messages() <-chan models.Message {
	return c.messages
}

func (c *MockClient) Disconnect() error {
	if c.cancel != nil {
		c.cancel()
	}
	close(c.messages)
	return nil
}
