package workflow

import (
	"context"
	"time"

	"github.com/satont/twir/apps/timers/internal/shared"
	"go.temporal.io/sdk/client"
)

func (c *Workflow) AddTimer(ctx context.Context, timerId string) error {
	timer, err := c.timersRepository.GetById(timerId)
	if err != nil {
		return err
	}

	scheduleID := timer.ID
	workflowID := "timers-" + timer.ID

	var every time.Duration
	if c.config.AppEnv == "development" {
		every = time.Duration(timer.Interval) * time.Second
	} else {
		every = time.Duration(timer.Interval) * time.Minute
	}

	_, err = c.cl.ScheduleClient().Create(
		ctx,
		client.ScheduleOptions{
			ID: scheduleID,
			Spec: client.ScheduleSpec{
				Intervals: []client.ScheduleIntervalSpec{
					{
						Every: every,
					},
				},
			},
			Action: &client.ScheduleWorkflowAction{
				ID:        workflowID,
				Workflow:  c.Flow,
				TaskQueue: shared.TimersWorkerTaskQueueName,
				Args:      []interface{}{timer},
				Memo:      map[string]interface{}{"lastResponse": 0},
			},
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func (c *Workflow) RemoveTimer(ctx context.Context, timerId string) error {
	timer, err := c.timersRepository.GetById(timerId)
	if err != nil {
		return err
	}

	scheduleID := timer.ID

	handle := c.cl.ScheduleClient().GetHandle(ctx, scheduleID)
	return handle.Delete(ctx)
}
