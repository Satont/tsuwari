package types

// export type YoutubeSettings = {
//   maxRequests?: number;
//   acceptOnlyWhenOnline?: boolean;
//   channelPointsRewardName?: string;
//   user?: {
//     maxRequests?: number;
//     minWatchTime?: number;
//     minMessages?: number;
//     minFollowTime?: number;
//   };
//   song?: {
//     maxLength?: number;
//     minViews?: number;
//     acceptedCategories?: string[];
//   };
//   blackList?: {
//     users?: Array<{ userId: string, userName: string }>;
//     songs?: Array<{ id: string, title: string, thumbnail: string }>;
//     channels?: Array<{ id: string, title: string, thumbnail: string }>;
//     artistsNames?: string[];
//   };
// };

type YoutubeUserSettings struct {
	MaxRequests   *int `json:"maxRequests"`
	MinWatchTime  *int `json:"minWatchTime"`
	MinMessages   *int `json:"minMessages"`
	MinFollowTime *int `json:"minFollowTime"`
}

type YotubeSongSettings struct {
	MaxLength          *int     `validate:"lte=86400"          json:"maxLength"`
	MinViews           *int     `validate:"lte=10000000000000" json:"minViews"`
	AcceptedCategories []string `validate:"dive,max=300"       json:"acceptedCategories"`
}

type YoutubeBlacklistSettingsUsers struct {
	UserID   string `json:"userId"   validate:"max=50"`
	UserName string `json:"userName" validate:"required"`
}

type YoutubeBlacklistSettingsSongs struct {
	ID        string `validate:"max=300" json:"id"`
	Title     string `validate:"max=300" json:"title"`
	ThumbNail string `validate:"max=300" json:"thumbNail"`
}

type YoutubeBlacklistSettingsChannels struct {
	ID        string `json:"id"`
	Title     string `json:"title"     validate:"max=300"`
	ThumbNail string `json:"thumbNail" validate:"max=300"`
}

type YoutubeBlacklistSettings struct {
	Users        []YoutubeBlacklistSettingsUsers    `validate:"dive"         json:"users"`
	Songs        []YoutubeBlacklistSettingsSongs    `validate:"dive"         json:"songs"`
	Channels     []YoutubeBlacklistSettingsChannels `validate:"dive"         json:"channels"`
	ArtistsNames []string                           `validate:"dive,max=300" json:"artistsNames"`
}

type YoutubeSettings struct {
	MaxRequests             *int                      `validate:"lte=500" json:"maxRequests"`
	AcceptOnlyWhenOnline    *bool                     `                   json:"acceptOnlyWhenOnline"`
	ChannelPointsRewardName *string                   `validate:"max=100" json:"channelPointsRewardName"`
	User                    *YoutubeUserSettings      `validate:"dive"    json:"user"`
	Song                    *YotubeSongSettings       `validate:"dive"    json:"song"`
	BlackList               *YoutubeBlacklistSettings `validate:"dive"    json:"blacklist"`
}
