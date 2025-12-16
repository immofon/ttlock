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

// PasscodeType represents the type of keyboard password
type PasscodeType int

const (
	PasscodeTypeOneTime         PasscodeType = 1  // 单次：只能在开始时间后6小时内使用一次
	PasscodeTypePermanent       PasscodeType = 2  // 永久：从开始时间开始永久有效；必需在开始时间24小时内使用一次，否则将失效
	PasscodeTypePeriod          PasscodeType = 3  // 限期：在开始和结束时间内有效；必需在开始时间24小时内使用一次，否则将失效
	PasscodeTypeDelete          PasscodeType = 4  // 删除：在锁上使用后会删除之前在锁上使用过的密码
	PasscodeTypeWeekendCyclic   PasscodeType = 5  // 周末循环：在周末开始和结束时间指定时间段内有效
	PasscodeTypeDailyCyclic     PasscodeType = 6  // 每日循环：每天开始和结束时间指定时间段内有效
	PasscodeTypeWorkdayCyclic   PasscodeType = 7  // 工作日循环：工作日开始和结束时间指定的时间段内有效
	PasscodeTypeMondayCyclic    PasscodeType = 8  // 周一循环：每周一开始和结束时间指定时间段内有效
	PasscodeTypeTuesdayCyclic   PasscodeType = 9  // 周二循环：每周二开始和结束时间指定时间段内有效
	PasscodeTypeWednesdayCyclic PasscodeType = 10 // 周三循环：每周三开始和结束时间指定时间段内有效
	PasscodeTypeThursdayCyclic  PasscodeType = 11 // 周四循环：每周四开始和结束时间指定时间段内有效
	PasscodeTypeFridayCyclic    PasscodeType = 12 // 周五循环：每周五开始和结束时间指定时间段内有效
	PasscodeTypeSaturdayCyclic  PasscodeType = 13 // 周六循环：每周六开始和结束时间指定时间段内有效
	PasscodeTypeSundayCyclic    PasscodeType = 14 // 周日循环：每周日开始和结束时间指定时间段内有效
)

// Passcode represents a keyboard password object
type Passcode struct {
	KeyboardPwdID   int    `json:"keyboardPwdId"`
	LockID          int    `json:"lockId"`
	KeyboardPwd     string `json:"keyboardPwd"`
	KeyboardPwdName string `json:"keyboardPwdName"`
	KeyboardPwdType int    `json:"keyboardPwdType"`
	StartDate       int64  `json:"startDate"`
	EndDate         int64  `json:"endDate"`
	SendDate        int64  `json:"sendDate"`
	IsCustom        int    `json:"isCustom"`
	SenderUsername  string `json:"senderUsername"`
}

// RandomPasscodeResponse represents the response for getting a random passcode
type RandomPasscodeResponse struct {
	KeyboardPwd   string `json:"keyboardPwd"`
	KeyboardPwdID int    `json:"keyboardPwdId"`
	Errcode       int    `json:"errcode"`
	Errmsg        string `json:"errmsg"`
}

// PasscodeListResponse represents the response for getting the passcode list
type PasscodeListResponse struct {
	List     []Passcode `json:"list"`
	PageNo   int        `json:"pageNo"`
	PageSize int        `json:"pageSize"`
	Pages    int        `json:"pages"`
	Total    int        `json:"total"`
	Errcode  int        `json:"errcode"`
	Errmsg   string     `json:"errmsg"`
}

// GetRandomPasscode retrieves a random passcode from the cloud.
// startDate and endDate are timestamps in milliseconds.
// endDate is optional (pass 0 if not needed).
// Note: The validity period of the passcode is precise to the hour.
// It is recommended to pass the timestamp of the hour (e.g., 19:00:00).
func (c *Client) GetRandomPasscode(lockID int, pwdType PasscodeType, pwdName string, startDate, endDate int64) (*RandomPasscodeResponse, error) {
	accessToken := c.AccessToken()
	endpoint := c.BaseURL + "/v3/keyboardPwd/get"

	data := url.Values{}
	data.Set("clientId", c.ClientID)
	data.Set("accessToken", accessToken)
	data.Set("lockId", strconv.Itoa(lockID))
	data.Set("keyboardPwdType", strconv.Itoa(int(pwdType)))
	data.Set("date", strconv.FormatInt(time.Now().UnixMilli(), 10))
	data.Set("startDate", strconv.FormatInt(startDate, 10))

	if pwdName != "" {
		data.Set("keyboardPwdName", pwdName)
	}
	if endDate != 0 {
		data.Set("endDate", strconv.FormatInt(endDate, 10))
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

	var result RandomPasscodeResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if result.Errcode != 0 {
		return nil, NewError(ErrorCode(result.Errcode))
	}

	return &result, nil
}

// GetPasscodeList retrieves the list of passcodes for a lock.
// orderBy: 0-Ascending by name, 1-Descending by creation time, 2-Descending by name
func (c *Client) GetPasscodeList(lockID, pageNo, pageSize, orderBy int, searchStr string) (*PasscodeListResponse, error) {
	accessToken := c.AccessToken()
	endpoint := c.BaseURL + "/v3/lock/listKeyboardPwd"

	params := url.Values{}
	params.Set("clientId", c.ClientID)
	params.Set("accessToken", accessToken)
	params.Set("lockId", strconv.Itoa(lockID))
	params.Set("pageNo", strconv.Itoa(pageNo))
	params.Set("pageSize", strconv.Itoa(pageSize))
	params.Set("orderBy", strconv.Itoa(orderBy))
	params.Set("date", strconv.FormatInt(time.Now().UnixMilli(), 10))

	if searchStr != "" {
		params.Set("searchStr", searchStr)
	}

	req, err := http.NewRequest("GET", endpoint+"?"+params.Encode(), nil)
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

	var result PasscodeListResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if result.Errcode != 0 {
		return nil, NewError(ErrorCode(result.Errcode))
	}

	return &result, nil
}

// PasscodeIterator allows iterating over passcodes without manually handling pagination
type PasscodeIterator struct {
	client      *Client
	accessToken string
	lockID      int
	orderBy     int
	searchStr   string
	pageNo      int
	pageSize    int
	items       []Passcode
	index       int
	done        bool
}

// Next returns the next passcode in the iterator.
// It returns nil, nil when there are no more passcodes.
func (it *PasscodeIterator) Next() (*Passcode, error) {
	if it.index >= len(it.items) {
		if it.done {
			return nil, nil
		}
		it.pageNo++
		resp, err := it.client.GetPasscodeList(it.lockID, it.pageNo, it.pageSize, it.orderBy, it.searchStr)
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

// IteratePasscodes creates a new iterator for passcodes.
func (c *Client) IteratePasscodes(lockID int, orderBy int, searchStr string) *PasscodeIterator {
	accessToken := c.AccessToken()
	return &PasscodeIterator{
		client:      c,
		accessToken: accessToken,
		lockID:      lockID,
		orderBy:     orderBy,
		searchStr:   searchStr,
		pageNo:      0,
		pageSize:    200, // Default page size
	}
}
