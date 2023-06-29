package processor

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/nicklaw5/helix/v2"
	"github.com/samber/lo"
	model "github.com/satont/twir/libs/gomodels"
)

func (c *Processor) getChannelVips() ([]helix.ChannelVips, error) {
	if c.cache.channelVips != nil {
		return c.cache.channelVips, nil
	}

	vips, err := c.streamerApiClient.GetChannelVips(
		&helix.GetChannelVipsParams{
			BroadcasterID: c.channelId,
		},
	)
	if err != nil {
		return nil, errors.New(vips.ErrorMessage)
	}

	if vips.ErrorMessage != "" {
		return nil, errors.New(vips.ErrorMessage)
	}

	c.cache.channelVips = vips.Data.ChannelsVips

	cursor := ""
	if vips.Data.Pagination.Cursor != "" {
		for {
			vips, err = c.streamerApiClient.GetChannelVips(
				&helix.GetChannelVipsParams{
					BroadcasterID: c.channelId,
					After:         cursor,
				},
			)

			if err != nil {
				return nil, errors.New(vips.ErrorMessage)
			}

			if vips.ErrorMessage != "" {
				return nil, errors.New(vips.ErrorMessage)
			}

			c.cache.channelVips = append(c.cache.channelVips, vips.Data.ChannelsVips...)

			if vips.Data.Pagination.Cursor == "" {
				break
			}

			cursor = vips.Data.Pagination.Cursor
		}
	}

	if len(c.cache.channelVips) == 0 {
		return nil, errors.New("cannot get vips")
	}

	return vips.Data.ChannelsVips, nil
}

func (c *Processor) VipOrUnvip(input string, operation model.EventOperationType) error {
	hydratedName, err := c.HydrateStringWithData(input, c.data)

	if err != nil || len(hydratedName) == 0 {
		return fmt.Errorf("cannot hydrate string %w", err)
	}

	hydratedName = strings.TrimSpace(strings.ReplaceAll(hydratedName, "@", ""))

	user, err := c.streamerApiClient.GetUsers(
		&helix.UsersParams{
			Logins: []string{hydratedName},
		},
	)

	if err != nil || len(user.Data.Users) == 0 {
		if err != nil {
			return err
		}
		return errors.New("cannot find user")
	}

	vips, err := c.getChannelVips()
	if err != nil {
		return err
	}

	dbChannel, err := c.getDbChannel()
	if err != nil {
		return err
	}

	if user.Data.Users[0].ID == dbChannel.BotID || user.Data.Users[0].ID == dbChannel.ID {
		return InternalError
	}

	mods, err := c.getChannelMods()
	if err != nil {
		return err
	}

	if lo.SomeBy(
		mods, func(item helix.Moderator) bool {
			return item.UserID == user.Data.Users[0].ID
		},
	) {
		return InternalError
	}

	isAlreadyVip := lo.SomeBy(
		vips, func(item helix.ChannelVips) bool {
			return item.UserID == user.Data.Users[0].ID
		},
	)

	if operation == "VIP" {
		if isAlreadyVip {
			return InternalError
		}

		resp, err := c.streamerApiClient.AddChannelVip(
			&helix.AddChannelVipParams{
				BroadcasterID: c.channelId,
				UserID:        user.Data.Users[0].ID,
			},
		)
		if resp.ErrorMessage != "" || err != nil {
			if err != nil {
				return err
			} else {
				return errors.New(resp.ErrorMessage)
			}
		} else {
			c.cache.channelVips = append(
				c.cache.channelVips, helix.ChannelVips{
					UserID:    user.Data.Users[0].ID,
					UserLogin: user.Data.Users[0].Login,
					UserName:  user.Data.Users[0].DisplayName,
				},
			)
		}
	} else {
		if !isAlreadyVip {
			return InternalError
		}
		resp, err := c.streamerApiClient.RemoveChannelVip(
			&helix.RemoveChannelVipParams{
				BroadcasterID: c.channelId,
				UserID:        user.Data.Users[0].ID,
			},
		)
		if resp.ErrorMessage != "" || err != nil {
			if err != nil {
				return err
			} else {
				return errors.New(resp.Error)
			}
		} else {
			c.cache.channelVips = lo.Filter(
				c.cache.channelVips, func(item helix.ChannelVips, index int) bool {
					return item.UserID != user.Data.Users[0].ID
				},
			)
		}
	}

	return nil
}

func (c *Processor) UnvipRandom(operation model.EventOperationType, slots string) error {
	vips, err := c.getChannelVips()
	if err != nil {
		return err
	}

	// if there is still slots available, we should skip unvip
	if operation == model.OperationUnvipRandomIfNoSlots {
		if slots == "" {
			return errors.New("input is empty")
		}

		slotsInt, err := strconv.Atoi(slots)
		if err != nil {
			return err
		}

		if len(vips) < slotsInt {
			return nil
		}
	}

	dbChannel, err := c.getDbChannel()
	if err != nil {
		return err
	}

	// choose random vip, but filter out bot account
	randomVip := lo.Sample(
		lo.Filter(
			vips, func(item helix.ChannelVips, index int) bool {
				return item.UserID != dbChannel.BotID
			},
		),
	)
	removeReq, err := c.streamerApiClient.RemoveChannelVip(
		&helix.RemoveChannelVipParams{
			BroadcasterID: c.channelId,
			UserID:        randomVip.UserID,
		},
	)
	if err != nil {
		return err
	}

	if removeReq.ErrorMessage != "" {
		return errors.New(removeReq.ErrorMessage)
	}

	c.cache.channelVips = lo.Filter(
		c.cache.channelVips, func(item helix.ChannelVips, index int) bool {
			return item.UserID != randomVip.UserID
		},
	)

	if len(c.data.PrevOperation.UnvipedUserName) > 0 {
		c.data.PrevOperation.UnvipedUserName += ", " + randomVip.UserName
	} else {
		c.data.PrevOperation.UnvipedUserName = randomVip.UserName
	}

	return nil
}
