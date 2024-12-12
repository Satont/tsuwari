package pgx

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/twirapp/twir/libs/repositories/timers"
	"github.com/twirapp/twir/libs/repositories/timers/model"
)

type Opts struct {
	Pgx *pgxpool.Pool
}

func New(opts Opts) *Pgx {
	return &Pgx{
		pool: opts.Pgx,
	}
}

func NewFx(pgxpool *pgxpool.Pool) *Pgx {
	return New(
		Opts{
			Pgx: pgxpool,
		},
	)
}

var _ timers.Repository = (*Pgx)(nil)
var sq = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

type Pgx struct {
	pool *pgxpool.Pool
}

func (c *Pgx) GetByID(ctx context.Context, id string) (model.Timer, error) {
	query := `
SELECT
    t."id", t."channelId", t."name", t."enabled", t."timeInterval", t."messageInterval", t."lastTriggerMessageNumber",
    (select array_agg(row(id, text, "isAnnounce") order by r.id) from channels_timers_responses r where r."timerId" = t.id) responses
FROM
    "channels_timers" t
WHERE
   "channelId" = $1
ORDER BY "id";
`
	rows, err := c.pool.Query(ctx, query, id)
	if err != nil {
		return model.Nil, err
	}

	timer, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.Timer])
	if err != nil {
		return model.Nil, err
	}
	return timer, nil
}

func (c *Pgx) CountByChannelID(ctx context.Context, channelID string) (int, error) {
	query := `SELECT count(*) from "channels_timers" where "channelId" = $1`

	var count int
	err := c.pool.QueryRow(ctx, query, channelID).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (c *Pgx) Create(ctx context.Context, data timers.CreateInput) (model.Timer, error) {
	createQuery := `
INSERT INTO "channels_timers" ("channelId", "name", "enabled", "timeInterval", "messageInterval")
VALUES ($1, $2, $3, $4, $5)
RETURNING "id", "channelId", "name", "enabled", "timeInterval", "messageInterval"
`
	createResponseQuery := `
INSERT INTO "channels_timers_responses" ("id", "text", "isAnnounce", "timerId")
VALUES ($1, $2, $3, $4)
RETURNING "id", "text", "isAnnounce", "timerId"
`
	tx, err := c.pool.Begin(ctx)
	if err != nil {
		return model.Nil, err
	}

	defer tx.Rollback(ctx)

	var newTimer model.Timer

	if err := tx.QueryRow(
		ctx,
		createQuery,
		data.ChannelID,
		data.Name,
		data.Enabled,
		data.TimeInterval,
		data.MessageInterval,
	).Scan(
		&newTimer.ID,
		&newTimer.ChannelID,
		&newTimer.Name,
		&newTimer.Enabled,
		&newTimer.TimeInterval,
		&newTimer.MessageInterval,
	); err != nil {
		return model.Nil, err
	}

	for _, r := range data.Responses {
		var newResponse model.Response

		if err := tx.QueryRow(
			ctx,
			createResponseQuery,
			uuid.New(),
			r.Text,
			r.IsAnnounce,
			newTimer.ID,
		).Scan(
			&newResponse.ID,
			&newResponse.Text,
			&newResponse.IsAnnounce,
			&newResponse.TimerID,
		); err != nil {
			return model.Nil, err
		}
		newTimer.Responses = append(newTimer.Responses, newResponse)
	}

	if err := tx.Commit(ctx); err != nil {
		return model.Nil, err
	}

	return newTimer, nil
}

func (c *Pgx) GetAllByChannelID(ctx context.Context, channelID string) ([]model.Timer, error) {
	query := `
SELECT
    t."id", t."channelId", t."name", t."enabled", t."timeInterval", t."messageInterval", t."lastTriggerMessageNumber",
    (select array_agg(row(id, text, "isAnnounce") order by r.id) from channels_timers_responses r where r."timerId" = t.id) responses
FROM
    "channels_timers" t
WHERE
   "channelId" = $1
ORDER BY "id";
`
	rows, err := c.pool.Query(ctx, query, channelID)
	if err != nil {
		return nil, err
	}

	models, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.Timer])
	if err != nil {
		return nil, err
	}
	return models, nil
}

// func (c *Pgx) GetAllByChannelID(ctx context.Context, channelID string) ([]model.Timer, error) {
// 	query := `
// SELECT t."id", t."channelId", t."name", t."enabled", t."timeInterval", t."messageInterval", t."lastTriggerMessageNumber",
// 			 r."id", r."text", r."isAnnounce", r."timerId"
// FROM "channels_timers" t
// LEFT JOIN "channels_timers_responses" r ON t."id" = r."timerId"
// WHERE t."channelId" = $1
// ORDER BY t."id";
// `
//
// 	rows, err := c.pool.Query(ctx, query, channelID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()
//
// 	var timersMap = make(map[uuid.UUID]*model.Timer)
// 	for rows.Next() {
// 		var timer model.Timer
// 		var response model.Response
// 		if err := rows.Scan(
// 			&timer.ID,
// 			&timer.ChannelID,
// 			&timer.Name,
// 			&timer.Enabled,
// 			&timer.TimeInterval,
// 			&timer.MessageInterval,
// 			&timer.LastTriggerMessageNumber,
// 			&response.ID,
// 			&response.Text,
// 			&response.IsAnnounce,
// 			&response.TimerID,
// 		); err != nil {
// 			return nil, err
// 		}
// 		if _, ok := timersMap[timer.ID]; !ok {
// 			timersMap[timer.ID] = &timer
// 		}
// 		timersMap[timer.ID].Responses = append(timersMap[timer.ID].Responses, response)
// 	}
//
// 	var timers []model.Timer
// 	for _, timer := range timersMap {
// 		timers = append(timers, *timer)
// 	}
//
// 	slices.SortFunc(
// 		timers, func(i, j model.Timer) int {
// 			return cmp.Compare(i.ID.String(), j.ID.String())
// 		},
// 	)
//
// 	return timers, nil
// }

func (c *Pgx) UpdateByID(ctx context.Context, id string, data timers.UpdateInput) (
	model.Timer,
	error,
) {
	updateBuilder := sq.Update("channels_timers")

	if data.Name != nil {
		updateBuilder = updateBuilder.Set("name", *data.Name)
	}

	if data.Enabled != nil {
		updateBuilder = updateBuilder.Set("enabled", *data.Enabled)
	}

	if data.TimeInterval != nil {
		updateBuilder = updateBuilder.Set("timeInterval", *data.TimeInterval)
	}

	if data.MessageInterval != nil {
		updateBuilder = updateBuilder.Set("messageInterval", *data.MessageInterval)
	}

	updateBuilder = updateBuilder.Where(squirrel.Eq{"id": id})

	query, args, err := updateBuilder.ToSql()
	if err != nil {
		return model.Nil, err
	}

	tx, err := c.pool.Begin(ctx)
	if err != nil {
		return model.Nil, err
	}
	defer tx.Rollback(ctx)

	result, err := tx.Exec(ctx, query, args...)
	if err != nil {
		return model.Nil, err
	}
	if result.RowsAffected() == 0 {
		return model.Nil, timers.ErrTimerNotFound
	}

	_, err = tx.Exec(ctx, `DELETE FROM channels_timers_responses WHERE "timerId" = $1`, id)
	if err != nil {
		return model.Nil, err
	}

	for _, r := range data.Responses {
		_, err := tx.Exec(
			ctx,
			`INSERT INTO "channels_timers_responses" ("id", "text", "isAnnounce", "timerId") VALUES ($1, $2, $3, $4)`,
			uuid.New(),
			r.Text,
			r.IsAnnounce,
			id,
		)
		if err != nil {
			return model.Nil, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return model.Nil, err
	}

	return c.GetByID(ctx, id)
}

func (c *Pgx) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM "channels_timers" WHERE "id" = $1`

	rows, err := c.pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if rows.RowsAffected() == 0 {
		return timers.ErrTimerNotFound
	}

	return nil
}
