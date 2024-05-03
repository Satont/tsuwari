package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.45

import (
	"context"

	model "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twir/apps/api-gql/internal/gql/gqlmodel"
)

// EmotesStatistics is the resolver for the emotesStatistics field.
func (r *queryResolver) EmotesStatistics(
	ctx context.Context,
	opts gqlmodel.EmotesStatisticsOpts,
) (*gqlmodel.EmotesStatisticResponse, error) {
	dashboardId, err := r.sessions.GetSelectedDashboard(ctx)
	if err != nil {
		return nil, err
	}

	var page int
	perPage := 10

	if opts.Page.IsSet() {
		page = *opts.Page.Value()
	}

	if opts.PerPage.IsSet() {
		perPage = *opts.PerPage.Value()
	}

	query := r.gorm.WithContext(ctx).
		Where(`"channelId" = ?`, dashboardId).
		Limit(perPage).
		Offset(page * perPage)

	if opts.Search.IsSet() {
		query = query.Where(`"emote" LIKE ?`, "%"+*opts.Search.Value()+"%")
	}
	var entities []emoteEntityModelWithCount
	if err :=
		query.
			Debug().
			Select(`"emote", COUNT(emote) as count`).
			Group("emote").
			Order("count DESC").
			Find(&entities).
			Error; err != nil {
		return nil, err
	}

	var totalCount int64
	if err := query.
		WithContext(ctx).
		Raw(
			`
	WITH emote_counts AS (
    SELECT
        emote,
        COUNT(emote) AS count
    FROM
        channels_emotes_usages
    WHERE
        "channelId" = '995463913'
    GROUP BY
        emote
    ORDER BY
        count DESC
    LIMIT 10
)
SELECT
    COUNT(*) AS count
FROM
    emote_counts;
`,
		).
		Scan(&totalCount).
		Error; err != nil {
		return nil, err
	}

	models := make([]gqlmodel.EmotesStatistic, 0, len(entities))
	for _, entity := range entities {
		lastUsedEntity := &model.ChannelEmoteUsage{}
		if err := r.gorm.
			WithContext(ctx).
			Where(`"channelId" = ? AND "emote" = ?`, dashboardId, entity.Emote).
			Order(`"createdAt" DESC`).
			First(lastUsedEntity).Error; err != nil {
			return nil, err
		}

		var rangeType gqlmodel.EmoteStatisticRange
		if opts.Range.IsSet() {
			rangeType = *opts.Range.Value()
		} else {
			rangeType = gqlmodel.EmoteStatisticRangeLast24Hours
		}

		usagesForLastDay, err := r.getEmoteStatisticUsagesForRange(
			ctx,
			entity.Emote,
			rangeType,
		)
		if err != nil {
			return nil, err
		}

		models = append(
			models, gqlmodel.EmotesStatistic{
				EmoteName:        entity.Emote,
				Usages:           entity.Count,
				LastUsedAt:       lastUsedEntity.CreatedAt,
				Last24HourUsages: usagesForLastDay,
			},
		)
	}

	return &gqlmodel.EmotesStatisticResponse{
		Emotes: models,
		Total:  int(totalCount),
	}, nil
}

// EmotesStatisticEmote is the resolver for the emotesStatisticEmote field.
func (r *queryResolver) EmotesStatisticEmote(
	ctx context.Context,
	opts gqlmodel.EmotesStatisticEmoteOpts,
) ([]gqlmodel.EmoteStatisticUsage, error) {
	return r.getEmoteStatisticUsagesForRange(ctx, opts.EmoteName, opts.Range)
}
