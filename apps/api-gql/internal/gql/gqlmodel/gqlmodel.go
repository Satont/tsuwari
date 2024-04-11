// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package gqlmodel

import (
	"fmt"
	"io"
	"strconv"

	"github.com/99designs/gqlgen/graphql"
)

type Notification interface {
	IsNotification()
	GetID() string
	GetUserID() *string
	GetText() string
}

type TwirUser interface {
	IsTwirUser()
	GetID() string
	GetTwitchProfile() *TwirUserTwitchInfo
}

type AdminNotification struct {
	ID            string              `json:"id"`
	Text          string              `json:"text"`
	UserID        *string             `json:"userId,omitempty"`
	TwitchProfile *TwirUserTwitchInfo `json:"twitchProfile,omitempty"`
}

func (AdminNotification) IsNotification()         {}
func (this AdminNotification) GetID() string      { return this.ID }
func (this AdminNotification) GetUserID() *string { return this.UserID }
func (this AdminNotification) GetText() string    { return this.Text }

type AuthenticatedUser struct {
	ID                string              `json:"id"`
	IsBotAdmin        bool                `json:"isBotAdmin"`
	IsBanned          bool                `json:"isBanned"`
	IsEnabled         *bool               `json:"isEnabled,omitempty"`
	IsBotModerator    *bool               `json:"isBotModerator,omitempty"`
	APIKey            string              `json:"apiKey"`
	HideOnLandingPage bool                `json:"hideOnLandingPage"`
	BotID             *string             `json:"botId,omitempty"`
	TwitchProfile     *TwirUserTwitchInfo `json:"twitchProfile"`
}

func (AuthenticatedUser) IsTwirUser()                                {}
func (this AuthenticatedUser) GetID() string                         { return this.ID }
func (this AuthenticatedUser) GetTwitchProfile() *TwirUserTwitchInfo { return this.TwitchProfile }

type Badge struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"createdAt"`
	FileURL   string `json:"fileUrl"`
	Enabled   bool   `json:"enabled"`
	// IDS of users which has this badge
	Users []string `json:"users,omitempty"`
}

type Command struct {
	ID                        string            `json:"id"`
	Name                      string            `json:"name"`
	Description               *string           `json:"description,omitempty"`
	Aliases                   []string          `json:"aliases,omitempty"`
	Responses                 []CommandResponse `json:"responses,omitempty"`
	Cooldown                  *int              `json:"cooldown,omitempty"`
	CooldownType              string            `json:"cooldownType"`
	Enabled                   bool              `json:"enabled"`
	Visible                   bool              `json:"visible"`
	Default                   bool              `json:"default"`
	DefaultName               *string           `json:"defaultName,omitempty"`
	Module                    string            `json:"module"`
	IsReply                   bool              `json:"isReply"`
	KeepResponsesOrder        bool              `json:"keepResponsesOrder"`
	DeniedUsersIds            []string          `json:"deniedUsersIds,omitempty"`
	AllowedUsersIds           []string          `json:"allowedUsersIds,omitempty"`
	RolesIds                  []string          `json:"rolesIds,omitempty"`
	OnlineOnly                bool              `json:"onlineOnly"`
	CooldownRolesIds          []string          `json:"cooldownRolesIds,omitempty"`
	EnabledCategories         []string          `json:"enabledCategories,omitempty"`
	RequiredWatchTime         int               `json:"requiredWatchTime"`
	RequiredMessages          int               `json:"requiredMessages"`
	RequiredUsedChannelPoints int               `json:"requiredUsedChannelPoints"`
}

type CommandResponse struct {
	ID        string `json:"id"`
	CommandID string `json:"commandId"`
	Text      string `json:"text"`
	Order     int    `json:"order"`
}

type CreateCommandInput struct {
	Name        string                                          `json:"name"`
	Description graphql.Omittable[*string]                      `json:"description,omitempty"`
	Aliases     graphql.Omittable[[]string]                     `json:"aliases,omitempty"`
	Responses   graphql.Omittable[[]CreateCommandResponseInput] `json:"responses,omitempty"`
}

type CreateCommandResponseInput struct {
	Text  string `json:"text"`
	Order int    `json:"order"`
}

type Mutation struct {
}

type NotificationUpdateOpts struct {
	Text graphql.Omittable[*string] `json:"text,omitempty"`
}

type Query struct {
}

type Subscription struct {
}

type TwirAdminUser struct {
	ID             string              `json:"id"`
	TwitchProfile  *TwirUserTwitchInfo `json:"twitchProfile"`
	IsBotAdmin     bool                `json:"isBotAdmin"`
	IsBanned       bool                `json:"isBanned"`
	IsBotModerator bool                `json:"isBotModerator"`
	IsBotEnabled   bool                `json:"isBotEnabled"`
	APIKey         string              `json:"apiKey"`
}

func (TwirAdminUser) IsTwirUser()                                {}
func (this TwirAdminUser) GetID() string                         { return this.ID }
func (this TwirAdminUser) GetTwitchProfile() *TwirUserTwitchInfo { return this.TwitchProfile }

type TwirBadgeUpdateOpts struct {
	Name    graphql.Omittable[*string]         `json:"name,omitempty"`
	File    graphql.Omittable[*graphql.Upload] `json:"file,omitempty"`
	Enabled graphql.Omittable[*bool]           `json:"enabled,omitempty"`
}

type TwirUserTwitchInfo struct {
	Login           string `json:"login"`
	DisplayName     string `json:"displayName"`
	ProfileImageURL string `json:"profileImageUrl"`
	Description     string `json:"description"`
}

type TwirUsersResponse struct {
	Users []TwirAdminUser `json:"users"`
	Total int             `json:"total"`
}

type TwirUsersSearchParams struct {
	Search       graphql.Omittable[*string]  `json:"search,omitempty"`
	Page         graphql.Omittable[*int]     `json:"page,omitempty"`
	PerPage      graphql.Omittable[*int]     `json:"perPage,omitempty"`
	IsBotAdmin   graphql.Omittable[*bool]    `json:"isBotAdmin,omitempty"`
	IsBanned     graphql.Omittable[*bool]    `json:"isBanned,omitempty"`
	IsBotEnabled graphql.Omittable[*bool]    `json:"isBotEnabled,omitempty"`
	Badges       graphql.Omittable[[]string] `json:"badges,omitempty"`
}

type UpdateCommandOpts struct {
	Name                      graphql.Omittable[*string]                      `json:"name,omitempty"`
	Description               graphql.Omittable[*string]                      `json:"description,omitempty"`
	Aliases                   graphql.Omittable[[]string]                     `json:"aliases,omitempty"`
	Responses                 graphql.Omittable[[]CreateCommandResponseInput] `json:"responses,omitempty"`
	Cooldown                  graphql.Omittable[*int]                         `json:"cooldown,omitempty"`
	CooldownType              graphql.Omittable[*string]                      `json:"cooldownType,omitempty"`
	Enabled                   graphql.Omittable[*bool]                        `json:"enabled,omitempty"`
	Visible                   graphql.Omittable[*bool]                        `json:"visible,omitempty"`
	IsReply                   graphql.Omittable[*bool]                        `json:"isReply,omitempty"`
	KeepResponsesOrder        graphql.Omittable[*bool]                        `json:"keepResponsesOrder,omitempty"`
	DeniedUsersIds            graphql.Omittable[[]string]                     `json:"deniedUsersIds,omitempty"`
	AllowedUsersIds           graphql.Omittable[[]string]                     `json:"allowedUsersIds,omitempty"`
	RolesIds                  graphql.Omittable[[]string]                     `json:"rolesIds,omitempty"`
	OnlineOnly                graphql.Omittable[*bool]                        `json:"onlineOnly,omitempty"`
	CooldownRolesIds          graphql.Omittable[[]string]                     `json:"cooldownRolesIds,omitempty"`
	EnabledCategories         graphql.Omittable[[]string]                     `json:"enabledCategories,omitempty"`
	RequiredWatchTime         graphql.Omittable[*int]                         `json:"requiredWatchTime,omitempty"`
	RequiredMessages          graphql.Omittable[*int]                         `json:"requiredMessages,omitempty"`
	RequiredUsedChannelPoints graphql.Omittable[*int]                         `json:"requiredUsedChannelPoints,omitempty"`
}

type UserNotification struct {
	ID     string  `json:"id"`
	UserID *string `json:"userId,omitempty"`
	Text   string  `json:"text"`
}

func (UserNotification) IsNotification()         {}
func (this UserNotification) GetID() string      { return this.ID }
func (this UserNotification) GetUserID() *string { return this.UserID }
func (this UserNotification) GetText() string    { return this.Text }

type NotificationType string

const (
	NotificationTypeGlobal NotificationType = "GLOBAL"
	NotificationTypeUser   NotificationType = "USER"
)

var AllNotificationType = []NotificationType{
	NotificationTypeGlobal,
	NotificationTypeUser,
}

func (e NotificationType) IsValid() bool {
	switch e {
	case NotificationTypeGlobal, NotificationTypeUser:
		return true
	}
	return false
}

func (e NotificationType) String() string {
	return string(e)
}

func (e *NotificationType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = NotificationType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid NotificationType", str)
	}
	return nil
}

func (e NotificationType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
