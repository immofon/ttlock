package ttlock

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// SendKeyResponse represents the response for sending an eKey
type SendKeyResponse struct {
	KeyID   int    `json:"keyId"`
	Errcode int    `json:"errcode"`
	Errmsg  string `json:"errmsg"`
}

// SendKeyOptions contains optional parameters for SendKey
type SendKeyOptions struct {
	Remarks      string // 备注，留言
	RemoteEnable int    // 是否支持远程开锁：1-是、2-否
	KeyRight     int    // 是否授权管理员钥匙：1-是、0-否，默认0不授权
	CreateUser   int    // 是否自动创建通通锁账号：1-是、2-否(默认2)，仅receiverUsername为邮箱或手机号时生效
}

// SendKey sends an eKey to a user.
// 参数说明：
// - clientId: 创建应用分配的client_id（由客户端自动填写）
// - accessToken: 访问令牌，通过获取访问令牌接口获取（由客户端自动填写）
// - lockID: 锁ID，由锁初始化接口生成
// - receiverUsername: 接收方用户名
// - keyName: 钥匙名
// - startDate: 有效期开始时间，时间戳(毫秒)
// - endDate: 有效期结束时间，时间戳(毫秒)
// - options: 选填参数，包含remarks/remoteEnable/keyRight/createUser
// - date: 当前时间(毫秒时间戳，由方法内部自动添加)
func (c *Client) SendKey(lockID int, receiverUsername, keyName string, startDate, endDate int64, options *SendKeyOptions) (*SendKeyResponse, error) {
	endpoint := c.BaseURL + "/v3/key/send"

	data := url.Values{}
	data.Set("clientId", c.ClientID)
	data.Set("accessToken", c.AccessToken())
	data.Set("lockId", strconv.Itoa(lockID))
	data.Set("receiverUsername", receiverUsername)
	data.Set("keyName", keyName)
	data.Set("startDate", strconv.FormatInt(startDate, 10))
	data.Set("endDate", strconv.FormatInt(endDate, 10))
	data.Set("date", strconv.FormatInt(time.Now().UnixMilli(), 10))

	if options != nil {
		if options.Remarks != "" {
			data.Set("remarks", options.Remarks)
		}
		if options.RemoteEnable != 0 {
			data.Set("remoteEnable", strconv.Itoa(options.RemoteEnable))
		}
		// KeyRight default is 0, so we only send if it's explicitly set to 1 (or if we want to support other values)
		// However, since 0 is the default, skipping it is fine if it's 0.
		// If the user sets it to 1, it will be sent.
		if options.KeyRight != 0 {
			data.Set("keyRight", strconv.Itoa(options.KeyRight))
		}
		if options.CreateUser != 0 {
			data.Set("createUser", strconv.Itoa(options.CreateUser))
		}
	}

	req, err := http.NewRequest("POST", endpoint, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var result SendKeyResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if result.Errcode != 0 {
		return nil, NewError(ErrorCode(result.Errcode))
	}

	return &result, nil
}
