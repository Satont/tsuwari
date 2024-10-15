package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.45

import (
	"context"
	"log/slog"

	model "github.com/satont/twir/libs/gomodels"
	data_loader "github.com/twirapp/twir/apps/api-gql/internal/gql/data-loader"
	"github.com/twirapp/twir/apps/api-gql/internal/gql/gqlmodel"
	"github.com/twirapp/twir/apps/api-gql/internal/gql/graph"
	"github.com/twirapp/twir/apps/api-gql/internal/gql/mappers"
)

// User is the resolver for the user field.
func (r *adminAuditLogResolver) User(ctx context.Context, obj *gqlmodel.AdminAuditLog) (*gqlmodel.TwirUserTwitchInfo, error) {
	if obj.UserID == nil {
		return nil, nil
	}

	return data_loader.GetHelixUserById(ctx, *obj.UserID)
}

// Channel is the resolver for the channel field.
func (r *adminAuditLogResolver) Channel(ctx context.Context, obj *gqlmodel.AdminAuditLog) (*gqlmodel.TwirUserTwitchInfo, error) {
	if obj.ChannelID == nil {
		return nil, nil
	}

	return data_loader.GetHelixUserById(ctx, *obj.ChannelID)
}

// AdminAuditLogs is the resolver for the adminAuditLogs field.
func (r *queryResolver) AdminAuditLogs(ctx context.Context, input gqlmodel.AdminAuditLogsInput) ([]gqlmodel.AdminAuditLog, error) {
	var page int
	perPage := 20

	if input.Page.IsSet() {
		page = *input.Page.Value()
	}

	if input.PerPage.IsSet() {
		perPage = *input.PerPage.Value()
	}

	query := r.gorm.
		WithContext(ctx)

	if input.UserID.IsSet() {
		query = query.Where("user_id = ?", *input.UserID.Value())
	}

	if input.ChannelID.IsSet() {
		query = query.Where("channel_id = ?", *input.ChannelID.Value())
	}

	if input.ObjectID.IsSet() {
		query = query.Where("object_id = ?", *input.ObjectID.Value())
	}

	if input.Table.IsSet() {
		query = query.Where("table_name = ?", *input.Table.Value())
	}

	if input.OperationType.IsSet() {
		query = query.Where(
			"operation_type = ?",
			mappers.AuditTypeGqlToModel(*input.OperationType.Value()),
		)
	}

	var logs []model.AuditLog
	if err := query.
		Limit(perPage).
		Offset(page * perPage).
		Order("created_at DESC").
		Find(&logs).Error; err != nil {
		r.logger.Error("error in fetching audit logs", slog.Any("err", err))
		return nil, err
	}

	gqllogs := make([]gqlmodel.AdminAuditLog, 0, len(logs))
	for _, l := range logs {
		gqllogs = append(
			gqllogs, gqlmodel.AdminAuditLog{
				ID:            l.ID,
				Table:         l.Table,
				OperationType: mappers.AuditTypeModelToGql(l.OperationType),
				OldValue:      l.OldValue.Ptr(),
				NewValue:      l.NewValue.Ptr(),
				ObjectID:      l.ObjectID.Ptr(),
				UserID:        l.UserID.Ptr(),
				ChannelID:     l.ChannelID.Ptr(),
				CreatedAt:     l.CreatedAt,
			},
		)
	}

	return gqllogs, nil
}

// AdminAuditLog returns graph.AdminAuditLogResolver implementation.
func (r *Resolver) AdminAuditLog() graph.AdminAuditLogResolver { return &adminAuditLogResolver{r} }

type adminAuditLogResolver struct{ *Resolver }
