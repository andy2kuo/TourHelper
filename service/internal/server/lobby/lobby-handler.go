package lobby

import (
	"net/http"

	"github.com/andy2kuo/TourHelper/internal/logger"
	"github.com/gin-gonic/gin"
)

// LoginRequest 登入請求結構
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Platform string `json:"platform"` // 平台: web, line, telegram
}

// LoginResponse 登入回應結構
type LoginResponse struct {
	Success  bool   `json:"success"`
	Message  string `json:"message"`
	Token    string `json:"token,omitempty"`
	MemberID string `json:"member_id,omitempty"`
}

// handleLogin 處理玩家登入驗證
func (s *LobbyServer) handleLogin(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, LoginResponse{
			Success: false,
			Message: "無效的請求格式",
		})
		return
	}

	logger.WithFields(map[string]interface{}{
		"username": req.Username,
		"platform": req.Platform,
	}).Info("收到登入請求")

	// TODO: 實作登入驗證邏輯
	// 1. 驗證使用者帳號密碼
	// 2. 產生 JWT Token
	// 3. 更新會員線上狀態
	// 4. 通知 Tour Server 會員已登入

	// 暫時回傳成功(實際實作時需要真正的驗證)
	c.JSON(http.StatusOK, LoginResponse{
		Success:  true,
		Message:  "登入成功(功能待實作)",
		Token:    "dummy_token",
		MemberID: req.Username,
	})
}

// LineLoginRequest LINE 登入請求結構
type LineLoginRequest struct {
	Code  string `json:"code" binding:"required"`  // LINE 授權碼
	State string `json:"state" binding:"required"` // 防 CSRF 攻擊的 state
}

// handleLineLogin 處理 LINE 第三方登入驗證
func (s *LobbyServer) handleLineLogin(c *gin.Context) {
	var req LineLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, LoginResponse{
			Success: false,
			Message: "無效的請求格式",
		})
		return
	}

	logger.WithFields(map[string]interface{}{
		"code":  req.Code,
		"state": req.State,
	}).Info("收到 LINE 登入請求")

	// TODO: LINE 第三方登入驗證功能待補
	// 1. 驗證 state 參數防止 CSRF 攻擊
	// 2. 使用 code 向 LINE 伺服器換取 access_token
	// 3. 使用 access_token 取得使用者資訊
	// 4. 檢查使用者是否已註冊,未註冊則自動註冊
	// 5. 產生 JWT Token
	// 6. 通知 Tour Server 會員已登入

	// 參考文件:
	// https://developers.line.biz/en/docs/line-login/integrate-line-login/

	c.JSON(http.StatusNotImplemented, LoginResponse{
		Success: false,
		Message: "LINE 第三方登入功能尚未實作",
	})
}

// LogoutRequest 登出請求結構
type LogoutRequest struct {
	Token string `json:"token" binding:"required"`
}

// handleLogout 處理玩家登出
func (s *LobbyServer) handleLogout(c *gin.Context) {
	var req LogoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "無效的請求格式",
		})
		return
	}

	logger.WithFields(map[string]interface{}{
		"token": req.Token,
	}).Info("收到登出請求")

	// TODO: 實作登出邏輯
	// 1. 驗證 Token
	// 2. 更新會員離線狀態
	// 3. 清除 Session
	// 4. 通知 Tour Server 會員已登出

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "登出成功(功能待實作)",
	})
}

// VerifyTokenRequest 驗證 Token 請求結構
type VerifyTokenRequest struct {
	Token string `json:"token" binding:"required"`
}

// VerifyTokenResponse 驗證 Token 回應結構
type VerifyTokenResponse struct {
	Valid    bool   `json:"valid"`
	MemberID string `json:"member_id,omitempty"`
	Message  string `json:"message,omitempty"`
}

// handleVerifyToken 驗證 Token 是否有效
func (s *LobbyServer) handleVerifyToken(c *gin.Context) {
	var req VerifyTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, VerifyTokenResponse{
			Valid:   false,
			Message: "無效的請求格式",
		})
		return
	}

	// TODO: 實作 Token 驗證邏輯
	// 1. 解析 JWT Token
	// 2. 驗證簽章
	// 3. 檢查過期時間
	// 4. 回傳會員資訊

	c.JSON(http.StatusOK, VerifyTokenResponse{
		Valid:    true,
		MemberID: "dummy_member_id",
		Message:  "Token 驗證成功(功能待實作)",
	})
}

// handleGetMemberInfo 取得會員資訊
func (s *LobbyServer) handleGetMemberInfo(c *gin.Context) {
	memberID := c.Param("id")

	logger.WithFields(map[string]interface{}{
		"member_id": memberID,
	}).Info("查詢會員資訊")

	// TODO: 實作取得會員資訊邏輯
	// 1. 驗證請求者權限
	// 2. 從資料庫查詢會員資訊
	// 3. 回傳會員資料

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"member": gin.H{
			"id":       memberID,
			"username": "dummy_user",
			"platform": "web",
			"status":   "online",
		},
		"message": "功能待實作",
	})
}

// UpdateMemberInfoRequest 更新會員資訊請求結構
type UpdateMemberInfoRequest struct {
	Nickname string                 `json:"nickname"`
	Avatar   string                 `json:"avatar"`
	Metadata map[string]interface{} `json:"metadata"`
}

// handleUpdateMemberInfo 更新會員資訊
func (s *LobbyServer) handleUpdateMemberInfo(c *gin.Context) {
	memberID := c.Param("id")

	var req UpdateMemberInfoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "無效的請求格式",
		})
		return
	}

	logger.WithFields(map[string]interface{}{
		"member_id": memberID,
		"nickname":  req.Nickname,
	}).Info("更新會員資訊")

	// TODO: 實作更新會員資訊邏輯
	// 1. 驗證請求者權限
	// 2. 更新資料庫
	// 3. 通知 Tour Server 會員資訊已更新

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "會員資訊更新成功(功能待實作)",
	})
}
