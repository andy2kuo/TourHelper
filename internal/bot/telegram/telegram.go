package telegram

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Bot Telegram Bot 結構
type Bot struct {
	token string
	api   *tgbotapi.BotAPI
}

// NewBot 建立新的 Telegram Bot
func NewBot(token string) *Bot {
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Printf("建立 Telegram Bot 客戶端錯誤: %v", err)
		return &Bot{token: token}
	}

	api.Debug = false
	log.Printf("已授權 Telegram Bot 帳號: %s", api.Self.UserName)

	return &Bot{
		token: token,
		api:   api,
	}
}

// HandleWebhook 處理 Telegram webhook
func (b *Bot) HandleWebhook(c *gin.Context) {
	var update tgbotapi.Update

	if err := c.ShouldBindJSON(&update); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無效的請求格式"})
		return
	}

	if update.Message == nil {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
		return
	}

	// 處理不同類型的訊息
	if update.Message.Text != "" {
		b.handleTextMessage(update.Message)
	} else if update.Message.Location != nil {
		b.handleLocationMessage(update.Message)
	}

	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

// handleTextMessage 處理文字訊息
func (b *Bot) handleTextMessage(message *tgbotapi.Message) {
	log.Printf("[%s] %s", message.From.UserName, message.Text)

	var replyText string

	// 處理指令
	switch message.Text {
	case "/start":
		replyText = "歡迎使用 TourHelper！\n\n我可以根據您的位置、天氣和偏好，為您推薦適合的旅遊景點。\n\n請使用以下指令：\n/recommend - 取得推薦\n/settings - 設定偏好\n/help - 查看說明"
	case "/recommend":
		replyText = "請分享您的位置資訊，我會為您推薦適合的旅遊景點！"
		// 建立請求位置的鍵盤
		keyboard := tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButtonLocation("分享我的位置"),
			),
		)
		msg := tgbotapi.NewMessage(message.Chat.ID, replyText)
		msg.ReplyMarkup = keyboard
		if _, err := b.api.Send(msg); err != nil {
			log.Printf("發送訊息錯誤: %v", err)
		}
		return
	case "/settings":
		replyText = "請告訴我您的偏好：\n\n1. 距離範圍（例如：50公里內）\n2. 景點類型（自然、文化、美食等）\n3. 預算（低、中、高）"
	case "/help":
		replyText = "TourHelper 使用說明：\n\n/recommend - 取得旅遊推薦\n/settings - 設定偏好\n/history - 查看歷史記錄\n\n您也可以直接分享位置，我會立即為您推薦景點！"
	default:
		replyText = fmt.Sprintf("您說：%s\n\n請使用 /recommend 來獲取旅遊建議。", message.Text)
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, replyText)
	if _, err := b.api.Send(msg); err != nil {
		log.Printf("發送訊息錯誤: %v", err)
	}
}

// handleLocationMessage 處理位置訊息
func (b *Bot) handleLocationMessage(message *tgbotapi.Message) {
	location := message.Location
	log.Printf("收到位置訊息: (%f, %f)", location.Latitude, location.Longitude)

	// TODO: 根據位置資訊推薦景點
	replyText := fmt.Sprintf("已收到您的位置：\n緯度：%f\n經度：%f\n\n正在為您尋找附近的景點...",
		location.Latitude, location.Longitude)

	msg := tgbotapi.NewMessage(message.Chat.ID, replyText)

	// 移除位置分享鍵盤
	msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)

	if _, err := b.api.Send(msg); err != nil {
		log.Printf("發送訊息錯誤: %v", err)
	}

	// TODO: 實際呼叫推薦服務並回傳結果
}

// SetWebhook 設定 webhook
func (b *Bot) SetWebhook(webhookURL string) error {
	webhook, err := tgbotapi.NewWebhook(webhookURL)
	if err != nil {
		return fmt.Errorf("建立 webhook 錯誤: %w", err)
	}

	_, err = b.api.Request(webhook)
	if err != nil {
		return fmt.Errorf("設定 webhook 錯誤: %w", err)
	}

	log.Printf("Telegram webhook 已設定為: %s", webhookURL)
	return nil
}
