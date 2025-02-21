package spotify

import (
	"fmt"
	"strings"

	model "github.com/satont/twir/libs/gomodels"

	"github.com/guregu/null"
	"github.com/imroc/req/v3"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

type Spotify struct {
	integration *model.ChannelsIntegrations
	db          *gorm.DB
	isRetry     bool
}

func New(integration *model.ChannelsIntegrations, db *gorm.DB) *Spotify {
	if integration == nil || !integration.AccessToken.Valid || !integration.RefreshToken.Valid {
		return nil
	}

	service := Spotify{
		integration: integration,
		db:          db,
		isRetry:     false,
	}

	return &service
}

type SpotifyRefreshResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (c *Spotify) refreshToken() error {
	data := SpotifyRefreshResponse{}

	body := make(map[string]string, 2)
	body["grant_type"] = "refresh_token"
	body["refresh_token"] = c.integration.RefreshToken.String

	resp, err := req.R().
		SetFormData(body).
		SetSuccessResult(&data).
		SetBasicAuth(
			c.integration.Integration.ClientID.String,
			c.integration.Integration.ClientSecret.String,
		).
		Post("https://accounts.spotify.com/api/token")

	if err != nil {
		return err
	}
	if !resp.IsSuccessState() {
		return fmt.Errorf("cannot refresh spotify token: %s", resp.String())
	}

	c.integration.AccessToken = null.StringFrom(data.AccessToken)
	if len(data.RefreshToken) > 0 {
		c.integration.RefreshToken = null.StringFrom(data.RefreshToken)
	}
	if err := c.db.Save(&c.integration).Error; err != nil {
		return fmt.Errorf("cannot save spotify token: %w", err)
	}

	return nil
}

type SpotifyArtist struct {
	Name string `json:"name"`
}

type SpotifyImage struct {
	URL string `json:"url"`
}

type SpotifyAlbum struct {
	Images []SpotifyImage `json:"images"`
}

type SpotifyTrack struct {
	Artists    []SpotifyArtist `json:"artists"`
	Name       string          `json:"name"`
	Album      SpotifyAlbum    `json:"album"`
	DurationMs int             `json:"duration_ms"`
}

type SpotifyResponse struct {
	Track      *SpotifyTrack `json:"item"`
	IsPlaying  bool          `json:"is_playing"`
	ProgressMs int           `json:"progress_ms"`
}

type GetTrackResponse struct {
	Title      string `json:"title"`
	Artist     string `json:"artist"`
	Image      string `json:"image"`
	IsPlaying  bool   `json:"isPlaying"`
	DurationMs int    `json:"durationMs"`
	ProgressMs int    `json:"progressMs"`
}

func (c *Spotify) GetTrack() (*GetTrackResponse, error) {
	data := SpotifyResponse{}
	resp, err := req.R().
		SetBearerAuthToken(c.integration.AccessToken.String).
		SetSuccessResult(&data).
		Get("https://api.spotify.com/v1/me/player/currently-playing")

	if resp.StatusCode == 401 && !c.isRetry {
		c.isRetry = true
		defer func() {
			c.isRetry = false
		}()
		if err := c.refreshToken(); err != nil {
			return nil, err
		}

		return c.GetTrack()
	}
	if err != nil {
		return nil, err
	}
	if !resp.IsSuccessState() {
		return nil, fmt.Errorf("cannot get spotify track: %s", resp.String())
	}

	if data.Track == nil {
		return nil, nil
	}

	artistsMap := lo.Map(
		data.Track.Artists, func(artist SpotifyArtist, _ int) string {
			return artist.Name
		},
	)

	var imageUrl string
	if len(data.Track.Album.Images) > 0 {
		imageUrl = data.Track.Album.Images[0].URL
	}

	return &GetTrackResponse{
		Artist:     strings.Join(artistsMap, ", "),
		Title:      data.Track.Name,
		Image:      imageUrl,
		IsPlaying:  data.IsPlaying,
		ProgressMs: data.ProgressMs,
		DurationMs: data.Track.DurationMs,
	}, nil
}

type SpotifyProfile struct {
	Country      string `json:"country"`
	DisplayName  string `json:"display_name"`
	Email        string `json:"email"`
	ExternalUrls struct {
		Spotify string `json:"spotify"`
	} `json:"external_urls"`
	Href      string `json:"href"`
	ID        string `json:"id"`
	Product   string `json:"product"`
	Type      string `json:"type"`
	URI       string `json:"uri"`
	Followers struct {
		Href  string `json:"href"`
		Total int    `json:"total"`
	} `json:"followers"`
	Images []struct {
		URL    string `json:"url"`
		Height int    `json:"height"`
		Width  int    `json:"width"`
	} `json:"images"`
	ExplicitContent struct {
		FilterEnabled bool `json:"filter_enabled"`
		FilterLocked  bool `json:"filter_locked"`
	} `json:"explicit_content"`
}

func (c *Spotify) GetProfile() (*SpotifyProfile, error) {
	data := SpotifyProfile{}
	resp, err := req.R().
		SetBearerAuthToken(c.integration.AccessToken.String).
		SetSuccessResult(&data).
		Get("https://api.spotify.com/v1/me")

	if resp.StatusCode == 401 && !c.isRetry {
		c.isRetry = true
		c.refreshToken()
		return c.GetProfile()
	}

	if err != nil {
		return nil, err
	}

	if !resp.IsSuccessState() {
		return nil, fmt.Errorf("cannot get profile: %s", resp.String())
	}

	return &data, nil
}
