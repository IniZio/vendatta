package github

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type UserInfo struct {
	Username  string `json:"username" yaml:"username"`
	UserID    int64  `json:"user_id" yaml:"user_id"`
	AvatarURL string `json:"avatar_url" yaml:"avatar_url"`
}

func (u *UserInfo) Validate() error {
	if u.Username == "" {
		return fmt.Errorf("user_info validation failed: username is required")
	}
	if u.UserID == 0 {
		return fmt.Errorf("user_info validation failed: user_id is required")
	}
	if u.AvatarURL == "" {
		return fmt.Errorf("user_info validation failed: avatar_url is required")
	}
	return nil
}

func ExtractUserInfo(ghCLIPath string) (*UserInfo, error) {
	username, err := ExecuteGHCommand(ghCLIPath, "api", "user", "--jq", ".login")
	if err != nil {
		return nil, fmt.Errorf("failed to extract username: %w", err)
	}

	userIDStr, err := ExecuteGHCommand(ghCLIPath, "api", "user", "--jq", ".id")
	if err != nil {
		return nil, fmt.Errorf("failed to extract user id: %w", err)
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid user id format: %w", err)
	}

	avatarURL, err := ExecuteGHCommand(ghCLIPath, "api", "user", "--jq", ".avatar_url")
	if err != nil {
		return nil, fmt.Errorf("failed to extract avatar url: %w", err)
	}

	userInfo := &UserInfo{
		Username:  username,
		UserID:    userID,
		AvatarURL: avatarURL,
	}

	if err := userInfo.Validate(); err != nil {
		return nil, err
	}

	return userInfo, nil
}

func MergeTempAndPersistentUserInfo(persistent, temp *UserInfo) *UserInfo {
	if persistent == nil {
		return temp
	}
	if temp == nil {
		return persistent
	}

	merged := &UserInfo{
		Username:  persistent.Username,
		UserID:    persistent.UserID,
		AvatarURL: persistent.AvatarURL,
	}

	if temp.Username != "" {
		merged.Username = temp.Username
	}
	if temp.UserID != 0 {
		merged.UserID = temp.UserID
	}
	if temp.AvatarURL != "" {
		merged.AvatarURL = temp.AvatarURL
	}

	return merged
}

func ParseUserInfoJSON(data []byte) (*UserInfo, error) {
	var userInfo UserInfo
	if err := json.Unmarshal(data, &userInfo); err != nil {
		return nil, fmt.Errorf("failed to parse user info: %w", err)
	}
	return &userInfo, nil
}
