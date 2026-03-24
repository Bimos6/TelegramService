package telegram

import (
	"context"

	"github.com/Bimos6/telegram-service/internal/models"
	"github.com/Bimos6/telegram-service/pkg/logger"
)

type Client struct {
	id       string
	appID    int
	appHash  string
	messages chan models.Message
	log      logger.Logger
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

// Connect - устанавливает соединение с Telegram
// Должен инициализировать клиент и подготовить к работе
func (c *Client) Connect(ctx context.Context) error {
	// TODO: реализовать
	// 1. Создать telegram.NewClient с appID, appHash
	// 2. Настроить SessionStorage для сохранения сессии
	// 3. Сохранить client и api для дальнейшего использования
	return nil
}

// GetQRCode - возвращает URL для QR кода авторизации
// Должен вернуть строку формата "tg://login?token=..."
func (c *Client) GetQRCode(ctx context.Context) (string, error) {
	// TODO: реализовать
	// 1. Запросить QR токен через client.Auth().QR()
	// 2. Сформировать URL: fmt.Sprintf("tg://login?token=%s", base64.URLEncoding.EncodeToString(token.Token))
	// 3. Запустить горутину для ожидания авторизации
	return "", nil
}

// SendMessage - отправляет текстовое сообщение
// peer - получатель ("@username" или "@me")
// text - текст сообщения
// возвращает ID сообщения
func (c *Client) SendMessage(ctx context.Context, peer, text string) (int64, error) {
	// TODO: реализовать
	// 1. Преобразовать peer в tg.InputPeerClass
	// 2. Отправить сообщение через sender.To(peer).Text()
	// 3. Получить ID сообщения из ответа
	return 0, nil
}

// Messages - возвращает канал для получения входящих сообщений
func (c *Client) Messages() <-chan models.Message {
	return c.messages
}

// Disconnect - закрывает соединение
func (c *Client) Disconnect() error {
	// TODO: реализовать
	// 1. Закрыть канал messages
	// 2. Закрыть клиент если нужно
	return nil
}
