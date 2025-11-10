package line

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/v8/linebot/messaging_api"
	"github.com/line/line-bot-sdk-go/v8/linebot/webhook"
)

// Bot Line Bot 結構
type Bot struct {
	channelSecret      string
	channelAccessToken string
	client             *messaging_api.MessagingApiAPI
}

// NewBot 建立新的 Line Bot
func NewBot(channelSecret, channelAccessToken string) *Bot {
	client, err := messaging_api.NewMessagingApiAPI(channelAccessToken)
	if err != nil {
		log.Printf("建立 Line Bot 客戶端錯誤: %v", err)
	}

	return &Bot{
		channelSecret:      channelSecret,
		channelAccessToken: channelAccessToken,
		client:             client,
	}
}

// HandleWebhook 處理 Line webhook
func (b *Bot) HandleWebhook(c *gin.Context) {
	cb, err := webhook.ParseRequest(b.channelSecret, c.Request)
	if err != nil {
		if err == webhook.ErrInvalidSignature {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid signature"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	for _, event := range cb.Events {
		switch e := event.(type) {
		case webhook.MessageEvent:
			switch message := e.Message.(type) {
			case webhook.TextMessageContent:
				b.handleTextMessage(e.ReplyToken, message.Text, e.Source)
			case webhook.LocationMessageContent:
				b.handleLocationMessage(e.ReplyToken, message, e.Source)
			}
		case webhook.FollowEvent:
			b.handleFollowEvent(e.ReplyToken, e.Source)
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

// handleTextMessage 處理文字訊息
func (b *Bot) handleTextMessage(replyToken, text string, source webhook.SourceInterface) {
	log.Printf("收到文字訊息: %s", text)

	var replyText string

	// TODO: 實作智慧回覆邏輯
	switch text {
	case "推薦":
		replyText = "請分享您的位置資訊，我會為您推薦適合的旅遊景點！"
	case "設定":
		replyText = "請告訴我您的偏好：\n1. 距離範圍（例如：50公里內）\n2. 景點類型（自然、文化、美食等）"
	default:
		replyText = fmt.Sprintf("您說：%s\n\n請輸入「推薦」來獲取旅遊建議，或輸入「設定」來調整偏好設定。", text)
	}

	if _, err := b.client.ReplyMessage(
		&messaging_api.ReplyMessageRequest{
			ReplyToken: replyToken,
			Messages: []messaging_api.MessageInterface{
				messaging_api.TextMessage{
					Text: replyText,
				},
			},
		},
	); err != nil {
		log.Printf("回覆訊息錯誤: %v", err)
	}
}

// handleLocationMessage 處理位置訊息
func (b *Bot) handleLocationMessage(replyToken string, location webhook.LocationMessageContent, source webhook.SourceInterface) {
	log.Printf("收到位置訊息: %s (%f, %f)", location.Address, location.Latitude, location.Longitude)

	// TODO: 根據位置資訊推薦景點
	replyText := fmt.Sprintf("已收到您的位置：%s\n正在為您尋找附近的景點...", location.Address)

	if _, err := b.client.ReplyMessage(
		&messaging_api.ReplyMessageRequest{
			ReplyToken: replyToken,
			Messages: []messaging_api.MessageInterface{
				messaging_api.TextMessage{
					Text: replyText,
				},
			},
		},
	); err != nil {
		log.Printf("回覆訊息錯誤: %v", err)
	}
}

// handleFollowEvent 處理加入好友事件
func (b *Bot) handleFollowEvent(replyToken string, source webhook.SourceInterface) {
	log.Println("新使用者加入")

	welcomeText := "歡迎使用 TourHelper！\n\n我可以根據您的位置、天氣和偏好，為您推薦適合的旅遊景點。\n\n請分享您的位置，或輸入「推薦」開始使用。"

	if _, err := b.client.ReplyMessage(
		&messaging_api.ReplyMessageRequest{
			ReplyToken: replyToken,
			Messages: []messaging_api.MessageInterface{
				messaging_api.TextMessage{
					Text: welcomeText,
				},
			},
		},
	); err != nil {
		log.Printf("回覆訊息錯誤: %v", err)
	}
}
