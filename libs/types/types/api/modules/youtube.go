package modules

type YouTubeUserSettings struct {
	MaxRequests   int   `json:"maxRequests"`
	MinWatchTime  int64 `json:"minWatchTime"`
	MinMessages   int   `json:"minMessages"`
	MinFollowTime int   `json:"minFollowTime"`
}

type YouTubeSongSettings struct {
	MinLength          int      `validate:"gte=0,lte=86399" json:"minLength"`
	MaxLength          int      `validate:"lte=86400"          json:"maxLength"`
	MinViews           int      `validate:"lte=10000000000000" json:"minViews"`
	AcceptedCategories []string `validate:"dive,max=300"       json:"acceptedCategories"`
}

type YouTubeDenySettingsUsers struct {
	UserID   string `json:"userId"   validate:"max=50"`
	UserName string `json:"userName" validate:"required"`
}

type YouTubeDenySettingsSongs struct {
	ID        string `validate:"required,min=1,max=300" json:"id"`
	Title     string `validate:"required,min=1,max=300" json:"title"`
	ThumbNail string `validate:"required,min=1,max=300" json:"thumbNail"`
}

type YouTubeDenySettingsChannels struct {
	ID        string `validate:"required,min=1"         json:"id"`
	Title     string `validate:"required,min=1,max=300" json:"title"`
	ThumbNail string `validate:"required,min=1,max=300" json:"thumbNail"`
}

type YouTubeDenyList struct {
	Users        []YouTubeDenySettingsUsers    `validate:"required,dive"         json:"users"`
	Songs        []YouTubeDenySettingsSongs    `validate:"required,dive"         json:"songs"`
	Channels     []YouTubeDenySettingsChannels `validate:"required,dive"         json:"channels"`
	ArtistsNames []string                      `validate:"required,dive,max=300" json:"artistsNames"`
}

type YouTubeTranslations struct {
	NowPlaying             string `json:"nowPlaying"`
	NotEnabled             string `json:"notEnabled"`
	NoText                 string `json:"noText"`
	AcceptOnlineWhenOnline string `json:"acceptOnlyWhenOnline"`
	User                   struct {
		Denied            string `json:"denied"`
		MaxRequests       string `json:"maxRequests"`
		NotEnoughMessages string `json:"notEnoughMessages"`
		NotEnoughWatched  string `json:"notEnoughWatched"`
		NotEnoughFollow   string `json:"notEnoughFollow"`
	} `json:"user"`
	Song struct {
		Denied               string `json:"denied"`
		NotFound             string `json:"notFound"`
		AlreadyInQueue       string `json:"alreadyInQueue"`
		AgeRestrictions      string `json:"ageRestrictions"`
		CannotGetInformation string `json:"cannotGetInformation"`
		Live                 string `json:"live"`
		MaxNumberRequested   string `json:"maxNumberRequested"`
		NotEnoughViews       string `json:"notEnoughViews"`
		MaxLength            string `json:"maxLength"`
		MinLength            string `json:"minLength"`
	} `json:"song"`
	Channel struct {
		Denied string `json:"denied"`
	}
}

type YouTubeSettings struct {
	Enabled               *bool               `validate:"required" json:"enabled"`
	AcceptOnlyWhenOnline  *bool               `validate:"required" json:"acceptOnlyWhenOnline"`
	MaxRequests           int                 `validate:"lte=500"  json:"maxRequests"`
	ChannelPointsRewardId string              `validate:"max=100"  json:"channelPointsRewardId"`
	AnnouncePlay          *bool               `validate:"required" json:"announcePlay"`
	User                  YouTubeUserSettings `validate:"required" json:"user"`
	Song                  YouTubeSongSettings `validate:"required" json:"song"`
	DenyList              YouTubeDenyList     `validate:"required" json:"denyList"`
	Translations          YouTubeTranslations `validate:"required" json:"translations"`
}

type SearchResult struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	ThumbNail string `json:"thumbNail"`
}

type YouTube struct {
	POST   YouTubeSettings
	GET    YouTubeSettings
	SEARCH []SearchResult
}
