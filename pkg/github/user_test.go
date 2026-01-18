package github

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExtractUserInfo(t *testing.T) {
	tests := []struct {
		name      string
		ghCLIPath string
		wantErr   bool
		validate  func(*UserInfo)
	}{
		{
			name:      "extract valid user info",
			ghCLIPath: "gh",
			wantErr:   false,
			validate: func(ui *UserInfo) {
				if ui != nil {
					assert.NotEmpty(t, ui.Username)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userInfo, err := ExtractUserInfo(tt.ghCLIPath)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, userInfo)
				if tt.validate != nil {
					tt.validate(userInfo)
				}
			}
		})
	}
}

func TestUserInfoValidate(t *testing.T) {
	tests := []struct {
		name     string
		userInfo *UserInfo
		wantErr  bool
		errMsg   string
	}{
		{
			name: "valid user info",
			userInfo: &UserInfo{
				Username:  "testuser",
				UserID:    12345,
				AvatarURL: "https://avatars.githubusercontent.com/u/12345",
			},
			wantErr: false,
		},
		{
			name: "missing username",
			userInfo: &UserInfo{
				UserID:    12345,
				AvatarURL: "https://avatars.githubusercontent.com/u/12345",
			},
			wantErr: true,
			errMsg:  "username",
		},
		{
			name: "zero user id",
			userInfo: &UserInfo{
				Username:  "testuser",
				AvatarURL: "https://avatars.githubusercontent.com/u/12345",
			},
			wantErr: true,
			errMsg:  "user_id",
		},
		{
			name: "missing avatar url",
			userInfo: &UserInfo{
				Username: "testuser",
				UserID:   12345,
			},
			wantErr: true,
			errMsg:  "avatar_url",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.userInfo.Validate()
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestMergeTempAndPersistentUserInfo(t *testing.T) {
	tests := []struct {
		name       string
		persistent *UserInfo
		temp       *UserInfo
		want       *UserInfo
	}{
		{
			name: "merge overlapping user info",
			persistent: &UserInfo{
				Username: "olduser",
				UserID:   123,
			},
			temp: &UserInfo{
				Username:  "newuser",
				UserID:    456,
				AvatarURL: "https://example.com/avatar.jpg",
			},
			want: &UserInfo{
				Username:  "newuser",
				UserID:    456,
				AvatarURL: "https://example.com/avatar.jpg",
			},
		},
		{
			name: "preserve persistent fields not in temp",
			persistent: &UserInfo{
				Username:  "user",
				UserID:    123,
				AvatarURL: "https://example.com/avatar.jpg",
			},
			temp: &UserInfo{
				Username: "user",
			},
			want: &UserInfo{
				Username:  "user",
				UserID:    123,
				AvatarURL: "https://example.com/avatar.jpg",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			merged := MergeTempAndPersistentUserInfo(tt.persistent, tt.temp)
			require.NotNil(t, merged)
			assert.Equal(t, tt.want.Username, merged.Username)
			assert.Equal(t, tt.want.UserID, merged.UserID)
			assert.Equal(t, tt.want.AvatarURL, merged.AvatarURL)
		})
	}
}
