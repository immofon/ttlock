package ttlock

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// Lock represents a lock object returned by the API
type Lock struct {
	LockID           int    `json:"lockId"`           // 锁ID，由锁初始化接口生成
	LockName         string `json:"lockName"`         // 锁的蓝牙名称
	LockAlias        string `json:"lockAlias"`        // 锁别名
	LockMac          string `json:"lockMac"`          // 锁MAC地址
	ElectricQuantity int    `json:"electricQuantity"` // 锁电量
	FeatureValue     string `json:"featureValue"`     // 锁特征值，用于表示锁支持的功能
	HasGateway       int    `json:"hasGateway"`       // 是否有连接网关：1-是、0-否
	LockData         string `json:"lockData"`         // 锁数据，用于操作锁
	GroupID          int    `json:"groupId"`          // 分组ID
	GroupName        string `json:"groupName"`        // 分组名称
	Date             int64  `json:"date"`             // 锁初始化时间（毫秒时间戳）
}

// LockListResponse represents the response for the lock list API
type LockListResponse struct {
	List     []Lock `json:"list"`
	PageNo   int    `json:"pageNo"`
	PageSize int    `json:"pageSize"`
	Pages    int    `json:"pages"`
	Total    int    `json:"total"`
	Errcode  int    `json:"errcode"`
	Errmsg   string `json:"errmsg"`
}

// LockDetail represents the detailed information of a lock
type LockDetail struct {
	LockID                int    `json:"lockId"`                // 锁ID
	LockName              string `json:"lockName"`              // 锁的蓝牙名称
	LockAlias             string `json:"lockAlias"`             // 锁别名
	LockMac               string `json:"lockMac"`               // 锁MAC地址
	NoKeyPwd              string `json:"noKeyPwd"`              // 管理员钥匙及键盘密码，管理员用该密码开门
	ElectricQuantity      int    `json:"electricQuantity"`      // 锁电量
	FeatureValue          string `json:"featureValue"`          // 锁特征值，表示锁支持的功能
	TimezoneRawOffset     int64  `json:"timezoneRawOffset"`     // 锁所在时区与UTC的毫秒差
	ModelNum              string `json:"modelNum"`              // 产品型号（用于固件升级）
	HardwareRevision      string `json:"hardwareRevision"`      // 硬件版本号（用于固件升级）
	FirmwareRevision      string `json:"firmwareRevision"`      // 固件版本号（用于固件升级）
	AutoLockTime          int    `json:"autoLockTime"`          // 自动闭锁时间（秒），-1 表示关闭自动闭锁
	LockSound             int    `json:"lockSound"`             // 锁声音开关：0-未知、1-开启、2-关闭
	PrivacyLock           int    `json:"privacyLock"`           // 反锁开关：0-未知、1-开启、2-关闭
	TamperAlert           int    `json:"tamperAlert"`           // 防撬开关：0-未知、1-开启、2-关闭
	ResetButton           int    `json:"resetButton"`           // 重置按键开关：0-未知、1-开启、2-关闭
	OpenDirection         int    `json:"openDirection"`         // 开门方向：0-未知、1-左开、2-右开
	PassageMode           int    `json:"passageMode"`           // 常开模式：1-开启、2-关闭
	PassageModeAutoUnlock int    `json:"passageModeAutoUnlock"` // 常开模式自动开锁：1-开启、2-关闭
	Date                  int64  `json:"date"`                  // 锁初始化时间（时间戳，毫秒）
	Errcode               int    `json:"errcode"`               // 接口错误码，0 表示成功
	Errmsg                string `json:"errmsg"`                // 接口错误描述
}

// GetLockList retrieves the list of locks for the account.
// lockAlias and groupId are optional filters. Pass empty string/0 to ignore.
func (c *Client) GetLockList(pageNo, pageSize int, lockAlias string, groupId int) (*LockListResponse, error) {
	accessToken := c.AccessToken()
	endpoint := c.BaseURL + "/v3/lock/list"

	params := url.Values{}
	params.Set("clientId", c.ClientID)
	params.Set("accessToken", accessToken)
	params.Set("pageNo", strconv.Itoa(pageNo))
	params.Set("pageSize", strconv.Itoa(pageSize))
	params.Set("date", strconv.FormatInt(time.Now().UnixMilli(), 10))

	if lockAlias != "" {
		params.Set("lockAlias", lockAlias)
	}
	if groupId != 0 {
		params.Set("groupId", strconv.Itoa(groupId))
	}

	reqURL := fmt.Sprintf("%s?%s", endpoint, params.Encode())
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var listResp LockListResponse
	if err := json.Unmarshal(body, &listResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if listResp.Errcode != 0 {
		return nil, NewError(ErrorCode(listResp.Errcode))
	}

	return &listResp, nil
}

// GetLockDetail retrieves the detailed information of a lock.
func (c *Client) GetLockDetail(lockId int) (*LockDetail, error) {
	accessToken := c.AccessToken()
	endpoint := c.BaseURL + "/v3/lock/detail"

	params := url.Values{}
	params.Set("clientId", c.ClientID)
	params.Set("accessToken", accessToken)
	params.Set("lockId", strconv.Itoa(lockId))
	params.Set("date", strconv.FormatInt(time.Now().UnixMilli(), 10))

	reqURL := fmt.Sprintf("%s?%s", endpoint, params.Encode())
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var detail LockDetail
	if err := json.Unmarshal(body, &detail); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if detail.Errcode != 0 {
		return nil, NewError(ErrorCode(detail.Errcode))
	}

	return &detail, nil
}

// LockIterator allows iterating over locks without manually handling pagination
type LockIterator struct {
	client      *Client
	accessToken string
	lockAlias   string
	groupId     int
	pageNo      int
	pageSize    int
	items       []Lock
	index       int
	done        bool
}

// Next returns the next lock in the iterator.
// It returns nil, nil when there are no more locks.
func (it *LockIterator) Next() (*Lock, error) {
	if it.index >= len(it.items) {
		if it.done {
			return nil, nil
		}
		it.pageNo++
		resp, err := it.client.GetLockList(it.pageNo, it.pageSize, it.lockAlias, it.groupId)
		if err != nil {
			return nil, err
		}
		if len(resp.List) == 0 {
			it.done = true
			return nil, nil
		}
		it.items = resp.List
		it.index = 0
		if it.pageNo >= resp.Pages {
			it.done = true
		}
	}

	item := &it.items[it.index]
	it.index++
	return item, nil
}

// IterateLocks creates a new iterator for locks.
func (c *Client) IterateLocks(lockAlias string, groupId int) *LockIterator {
	accessToken := c.AccessToken()
	return &LockIterator{
		client:      c,
		accessToken: accessToken,
		lockAlias:   lockAlias,
		groupId:     groupId,
		pageNo:      0,
		pageSize:    200, // Default page size
	}
}
