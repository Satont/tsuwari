package timers

import (
	"context"
	"github.com/samber/do"
	"github.com/samber/lo"
	"github.com/satont/twir/apps/api/internal/di"
	"github.com/satont/twir/apps/api/internal/interfaces"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/grpc/generated/timers"

	"github.com/satont/twir/apps/api/internal/types"
)

func handleGet(channelId string) []model.ChannelsTimers {
	timersService := do.MustInvoke[interfaces.TimersService](di.Provider)

	timers, err := timersService.FindManyByChannelId(channelId)
	if err != nil {
		return nil
	}

	return timers
}

func handlePost(
	channelId string,
	dto *timerDto,
	services types.Services,
) (*model.ChannelsTimers, error) {
	timersService := do.MustInvoke[interfaces.TimersService](di.Provider)
	timersGrpc := do.MustInvoke[timers.TimersClient](di.Provider)

	responses := lo.Map(
		dto.Responses, func(r responseDto, _ int) model.ChannelsTimersResponses {
			return model.ChannelsTimersResponses{
				Text:       r.Text,
				IsAnnounce: *r.IsAnnounce,
			}
		},
	)

	timer, err := timersService.Create(
		model.ChannelsTimers{
			ChannelID:                channelId,
			Name:                     dto.Name,
			Enabled:                  *dto.Enabled,
			TimeInterval:             int32(dto.TimeInterval),
			MessageInterval:          int32(dto.MessageInterval),
			LastTriggerMessageNumber: 0,
		}, responses,
	)

	if err != nil {
		return nil, err
	}

	timersGrpc.AddTimerToQueue(
		context.Background(), &timers.Request{
			TimerId: timer.ID,
		},
	)

	return timer, nil
}

func handleDelete(timerId string, services types.Services) error {
	timersService := do.MustInvoke[interfaces.TimersService](di.Provider)
	timersGrpc := do.MustInvoke[timers.TimersClient](di.Provider)

	err := timersService.Delete(timerId)

	if err != nil {
		return err
	}

	timersGrpc.RemoveTimerFromQueue(
		context.Background(), &timers.Request{
			TimerId: timerId,
		},
	)

	return nil
}

func handlePut(
	timerId string,
	dto *timerDto,
	services types.Services,
) (*model.ChannelsTimers, error) {
	timersService := do.MustInvoke[interfaces.TimersService](di.Provider)
	timersGrpc := do.MustInvoke[timers.TimersClient](di.Provider)

	responses := lo.Map(
		dto.Responses, func(r responseDto, _ int) model.ChannelsTimersResponses {
			return model.ChannelsTimersResponses{
				Text:       r.Text,
				IsAnnounce: *r.IsAnnounce,
			}
		},
	)

	timer, err := timersService.Update(
		timerId,
		model.ChannelsTimers{
			Name:            dto.Name,
			MessageInterval: int32(dto.MessageInterval),
			TimeInterval:    int32(dto.TimeInterval),
			Enabled:         *dto.Enabled,
		},
		responses,
	)

	if err != nil {
		return nil, err
	}

	if timer.Enabled {
		timersGrpc.AddTimerToQueue(
			context.Background(), &timers.Request{
				TimerId: timer.ID,
			},
		)
	} else {
		timersGrpc.RemoveTimerFromQueue(
			context.Background(), &timers.Request{
				TimerId: timer.ID,
			},
		)
	}

	return timer, nil
}

func handlePatch(
	timerId string,
	dto *timerPatchDto,
	services types.Services,
) (*model.ChannelsTimers, error) {
	timersService := do.MustInvoke[interfaces.TimersService](di.Provider)
	timersGrpc := do.MustInvoke[timers.TimersClient](di.Provider)

	updatedTimer, err := timersService.SetEnabled(timerId, *dto.Enabled)
	if err != nil {
		return nil, err
	}

	if updatedTimer.Enabled {
		timersGrpc.AddTimerToQueue(
			context.Background(), &timers.Request{
				TimerId: timerId,
			},
		)
	} else {
		timersGrpc.RemoveTimerFromQueue(
			context.Background(), &timers.Request{
				TimerId: timerId,
			},
		)
	}

	return updatedTimer, nil
}
