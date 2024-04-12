/* eslint-disable */
import type { TypedDocumentNode as DocumentNode } from '@graphql-typed-document-node/core';
export type Maybe<T> = T | null;
export type InputMaybe<T> = Maybe<T>;
export type Exact<T extends { [key: string]: unknown }> = { [K in keyof T]: T[K] };
export type MakeOptional<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]?: Maybe<T[SubKey]> };
export type MakeMaybe<T, K extends keyof T> = Omit<T, K> & { [SubKey in K]: Maybe<T[SubKey]> };
export type MakeEmpty<T extends { [key: string]: unknown }, K extends keyof T> = { [_ in K]?: never };
export type Incremental<T> = T | { [P in keyof T]?: P extends ' $fragmentName' | '__typename' ? T[P] : never };
/** All built-in and custom scalars, mapped to their actual values */
export type Scalars = {
  ID: { input: string; output: string; }
  String: { input: string; output: string; }
  Boolean: { input: boolean; output: boolean; }
  Int: { input: number; output: number; }
  Float: { input: number; output: number; }
  Time: { input: any; output: any; }
  Upload: { input: any; output: any; }
};

export type AdminNotification = Notification & {
  __typename?: 'AdminNotification';
  createdAt: Scalars['Time']['output'];
  id: Scalars['ID']['output'];
  text: Scalars['String']['output'];
  twitchProfile?: Maybe<TwirUserTwitchInfo>;
  userId?: Maybe<Scalars['ID']['output']>;
};

export type AdminNotificationsParams = {
  page?: InputMaybe<Scalars['Int']['input']>;
  perPage?: InputMaybe<Scalars['Int']['input']>;
  search?: InputMaybe<Scalars['String']['input']>;
  type?: InputMaybe<NotificationType>;
};

export type AdminNotificationsResponse = {
  __typename?: 'AdminNotificationsResponse';
  notifications: Array<AdminNotification>;
  total: Scalars['Int']['output'];
};

export type AuthenticatedUser = TwirUser & {
  __typename?: 'AuthenticatedUser';
  apiKey: Scalars['String']['output'];
  botId?: Maybe<Scalars['ID']['output']>;
  hideOnLandingPage: Scalars['Boolean']['output'];
  id: Scalars['ID']['output'];
  isBanned: Scalars['Boolean']['output'];
  isBotAdmin: Scalars['Boolean']['output'];
  isBotModerator?: Maybe<Scalars['Boolean']['output']>;
  isEnabled?: Maybe<Scalars['Boolean']['output']>;
  twitchProfile: TwirUserTwitchInfo;
};

export type Badge = {
  __typename?: 'Badge';
  createdAt: Scalars['String']['output'];
  enabled: Scalars['Boolean']['output'];
  fileUrl: Scalars['String']['output'];
  id: Scalars['ID']['output'];
  name: Scalars['String']['output'];
  /** IDS of users which has this badge */
  users?: Maybe<Array<Scalars['String']['output']>>;
};

export type Command = {
  __typename?: 'Command';
  aliases?: Maybe<Array<Scalars['String']['output']>>;
  allowedUsersIds?: Maybe<Array<Scalars['String']['output']>>;
  cooldown?: Maybe<Scalars['Int']['output']>;
  cooldownRolesIds?: Maybe<Array<Scalars['String']['output']>>;
  cooldownType: Scalars['String']['output'];
  default: Scalars['Boolean']['output'];
  defaultName?: Maybe<Scalars['String']['output']>;
  deniedUsersIds?: Maybe<Array<Scalars['String']['output']>>;
  description?: Maybe<Scalars['String']['output']>;
  enabled: Scalars['Boolean']['output'];
  enabledCategories?: Maybe<Array<Scalars['String']['output']>>;
  id: Scalars['ID']['output'];
  isReply: Scalars['Boolean']['output'];
  keepResponsesOrder: Scalars['Boolean']['output'];
  module: Scalars['String']['output'];
  name: Scalars['String']['output'];
  onlineOnly: Scalars['Boolean']['output'];
  requiredMessages: Scalars['Int']['output'];
  requiredUsedChannelPoints: Scalars['Int']['output'];
  requiredWatchTime: Scalars['Int']['output'];
  responses?: Maybe<Array<CommandResponse>>;
  rolesIds?: Maybe<Array<Scalars['String']['output']>>;
  visible: Scalars['Boolean']['output'];
};

export type CommandResponse = {
  __typename?: 'CommandResponse';
  commandId: Scalars['ID']['output'];
  id: Scalars['ID']['output'];
  order: Scalars['Int']['output'];
  text: Scalars['String']['output'];
};

export type CreateCommandInput = {
  aliases?: InputMaybe<Array<Scalars['String']['input']>>;
  description?: InputMaybe<Scalars['String']['input']>;
  name: Scalars['String']['input'];
  responses?: InputMaybe<Array<CreateCommandResponseInput>>;
};

export type CreateCommandResponseInput = {
  order: Scalars['Int']['input'];
  text: Scalars['String']['input'];
};

export type DashboardStats = {
  __typename?: 'DashboardStats';
  categoryId: Scalars['ID']['output'];
  categoryName: Scalars['String']['output'];
  chatMessages: Scalars['Int']['output'];
  followers: Scalars['Int']['output'];
  requestedSongs: Scalars['Int']['output'];
  startedAt?: Maybe<Scalars['Time']['output']>;
  subs: Scalars['Int']['output'];
  title: Scalars['String']['output'];
  usedEmotes: Scalars['Int']['output'];
  viewers?: Maybe<Scalars['Int']['output']>;
};

export type Mutation = {
  __typename?: 'Mutation';
  badgesAddUser: Scalars['Boolean']['output'];
  badgesCreate: Badge;
  badgesDelete: Scalars['Boolean']['output'];
  badgesRemoveUser: Scalars['Boolean']['output'];
  badgesUpdate: Badge;
  createCommand: Command;
  notificationsCreate: AdminNotification;
  notificationsDelete: Scalars['Boolean']['output'];
  notificationsUpdate: AdminNotification;
  removeCommand: Scalars['Boolean']['output'];
  switchUserAdmin: Scalars['Boolean']['output'];
  switchUserBan: Scalars['Boolean']['output'];
  updateCommand: Command;
};


export type MutationBadgesAddUserArgs = {
  id: Scalars['ID']['input'];
  userId: Scalars['String']['input'];
};


export type MutationBadgesCreateArgs = {
  file: Scalars['Upload']['input'];
  name: Scalars['String']['input'];
};


export type MutationBadgesDeleteArgs = {
  id: Scalars['ID']['input'];
};


export type MutationBadgesRemoveUserArgs = {
  id: Scalars['ID']['input'];
  userId: Scalars['String']['input'];
};


export type MutationBadgesUpdateArgs = {
  id: Scalars['ID']['input'];
  opts: TwirBadgeUpdateOpts;
};


export type MutationCreateCommandArgs = {
  opts: CreateCommandInput;
};


export type MutationNotificationsCreateArgs = {
  text: Scalars['String']['input'];
  userId?: InputMaybe<Scalars['String']['input']>;
};


export type MutationNotificationsDeleteArgs = {
  id: Scalars['ID']['input'];
};


export type MutationNotificationsUpdateArgs = {
  id: Scalars['ID']['input'];
  opts: NotificationUpdateOpts;
};


export type MutationRemoveCommandArgs = {
  id: Scalars['String']['input'];
};


export type MutationSwitchUserAdminArgs = {
  userId: Scalars['ID']['input'];
};


export type MutationSwitchUserBanArgs = {
  userId: Scalars['ID']['input'];
};


export type MutationUpdateCommandArgs = {
  id: Scalars['String']['input'];
  opts: UpdateCommandOpts;
};

export type Notification = {
  createdAt: Scalars['Time']['output'];
  id: Scalars['ID']['output'];
  text: Scalars['String']['output'];
  userId?: Maybe<Scalars['ID']['output']>;
};

export enum NotificationType {
  Global = 'GLOBAL',
  User = 'USER'
}

export type NotificationUpdateOpts = {
  text?: InputMaybe<Scalars['String']['input']>;
};

export type Query = {
  __typename?: 'Query';
  authenticatedUser: AuthenticatedUser;
  commands: Array<Command>;
  notificationsByAdmin: AdminNotificationsResponse;
  notificationsByUser: Array<UserNotification>;
  /** Twir badges */
  twirBadges: Array<Badge>;
  /** finding users on twitch with filter does they exists in database */
  twirUsers: TwirUsersResponse;
};


export type QueryNotificationsByAdminArgs = {
  opts: AdminNotificationsParams;
};


export type QueryTwirUsersArgs = {
  opts: TwirUsersSearchParams;
};

export type Subscription = {
  __typename?: 'Subscription';
  dashboardStats: DashboardStats;
  /** `newNotification` will return a stream of `Notification` objects. */
  newNotification: UserNotification;
};

export type TwirAdminUser = TwirUser & {
  __typename?: 'TwirAdminUser';
  apiKey: Scalars['String']['output'];
  id: Scalars['ID']['output'];
  isBanned: Scalars['Boolean']['output'];
  isBotAdmin: Scalars['Boolean']['output'];
  isBotEnabled: Scalars['Boolean']['output'];
  isBotModerator: Scalars['Boolean']['output'];
  twitchProfile: TwirUserTwitchInfo;
};

export type TwirBadgeUpdateOpts = {
  enabled?: InputMaybe<Scalars['Boolean']['input']>;
  file?: InputMaybe<Scalars['Upload']['input']>;
  name?: InputMaybe<Scalars['String']['input']>;
};

export type TwirUser = {
  id: Scalars['ID']['output'];
  twitchProfile: TwirUserTwitchInfo;
};

export type TwirUserTwitchInfo = {
  __typename?: 'TwirUserTwitchInfo';
  description: Scalars['String']['output'];
  displayName: Scalars['String']['output'];
  login: Scalars['String']['output'];
  profileImageUrl: Scalars['String']['output'];
};

export type TwirUsersResponse = {
  __typename?: 'TwirUsersResponse';
  total: Scalars['Int']['output'];
  users: Array<TwirAdminUser>;
};

export type TwirUsersSearchParams = {
  badges?: InputMaybe<Array<Scalars['String']['input']>>;
  isBanned?: InputMaybe<Scalars['Boolean']['input']>;
  isBotAdmin?: InputMaybe<Scalars['Boolean']['input']>;
  isBotEnabled?: InputMaybe<Scalars['Boolean']['input']>;
  page?: InputMaybe<Scalars['Int']['input']>;
  perPage?: InputMaybe<Scalars['Int']['input']>;
  search?: InputMaybe<Scalars['String']['input']>;
};

export type UpdateCommandOpts = {
  aliases?: InputMaybe<Array<Scalars['String']['input']>>;
  allowedUsersIds?: InputMaybe<Array<Scalars['String']['input']>>;
  cooldown?: InputMaybe<Scalars['Int']['input']>;
  cooldownRolesIds?: InputMaybe<Array<Scalars['String']['input']>>;
  cooldownType?: InputMaybe<Scalars['String']['input']>;
  deniedUsersIds?: InputMaybe<Array<Scalars['String']['input']>>;
  description?: InputMaybe<Scalars['String']['input']>;
  enabled?: InputMaybe<Scalars['Boolean']['input']>;
  enabledCategories?: InputMaybe<Array<Scalars['String']['input']>>;
  isReply?: InputMaybe<Scalars['Boolean']['input']>;
  keepResponsesOrder?: InputMaybe<Scalars['Boolean']['input']>;
  name?: InputMaybe<Scalars['String']['input']>;
  onlineOnly?: InputMaybe<Scalars['Boolean']['input']>;
  requiredMessages?: InputMaybe<Scalars['Int']['input']>;
  requiredUsedChannelPoints?: InputMaybe<Scalars['Int']['input']>;
  requiredWatchTime?: InputMaybe<Scalars['Int']['input']>;
  responses?: InputMaybe<Array<CreateCommandResponseInput>>;
  rolesIds?: InputMaybe<Array<Scalars['String']['input']>>;
  visible?: InputMaybe<Scalars['Boolean']['input']>;
};

export type UserNotification = Notification & {
  __typename?: 'UserNotification';
  createdAt: Scalars['Time']['output'];
  id: Scalars['ID']['output'];
  text: Scalars['String']['output'];
  userId?: Maybe<Scalars['ID']['output']>;
};

export type NotificationsGetAllQueryVariables = Exact<{ [key: string]: never; }>;


export type NotificationsGetAllQuery = { __typename?: 'Query', notificationsByUser: Array<{ __typename?: 'UserNotification', id: string, text: string, createdAt: any }> };

export type NotificationsSubscriptionSubscriptionVariables = Exact<{ [key: string]: never; }>;


export type NotificationsSubscriptionSubscription = { __typename?: 'Subscription', newNotification: { __typename?: 'UserNotification', id: string, text: string, createdAt: any } };

export type NotificationsByAdminQueryVariables = Exact<{
  opts: AdminNotificationsParams;
}>;


export type NotificationsByAdminQuery = { __typename?: 'Query', notificationsByAdmin: { __typename?: 'AdminNotificationsResponse', total: number, notifications: Array<{ __typename?: 'AdminNotification', id: string, text: string, userId?: string | null, createdAt: any, twitchProfile?: { __typename?: 'TwirUserTwitchInfo', displayName: string, profileImageUrl: string } | null }> } };

export type CreateNotificationMutationVariables = Exact<{
  text: Scalars['String']['input'];
  userId?: InputMaybe<Scalars['String']['input']>;
}>;


export type CreateNotificationMutation = { __typename?: 'Mutation', notificationsCreate: { __typename?: 'AdminNotification', text: string, userId?: string | null } };

export type DeleteNotificationMutationVariables = Exact<{
  id: Scalars['ID']['input'];
}>;


export type DeleteNotificationMutation = { __typename?: 'Mutation', notificationsDelete: boolean };

export type UpdateNotificationsMutationVariables = Exact<{
  id: Scalars['ID']['input'];
  opts: NotificationUpdateOpts;
}>;


export type UpdateNotificationsMutation = { __typename?: 'Mutation', notificationsUpdate: { __typename?: 'AdminNotification', id: string, text: string } };

export type DashboardStatsSubscriptionVariables = Exact<{ [key: string]: never; }>;


export type DashboardStatsSubscription = { __typename?: 'Subscription', dashboardStats: { __typename?: 'DashboardStats', categoryId: string, categoryName: string, viewers?: number | null, startedAt?: any | null, title: string, chatMessages: number, followers: number, usedEmotes: number, requestedSongs: number, subs: number } };


export const NotificationsGetAllDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"NotificationsGetAll"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"notificationsByUser"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"text"}},{"kind":"Field","name":{"kind":"Name","value":"createdAt"}}]}}]}}]} as unknown as DocumentNode<NotificationsGetAllQuery, NotificationsGetAllQueryVariables>;
export const NotificationsSubscriptionDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"subscription","name":{"kind":"Name","value":"NotificationsSubscription"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"newNotification"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"text"}},{"kind":"Field","name":{"kind":"Name","value":"createdAt"}}]}}]}}]} as unknown as DocumentNode<NotificationsSubscriptionSubscription, NotificationsSubscriptionSubscriptionVariables>;
export const NotificationsByAdminDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"query","name":{"kind":"Name","value":"notificationsByAdmin"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"opts"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"AdminNotificationsParams"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"notificationsByAdmin"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"opts"},"value":{"kind":"Variable","name":{"kind":"Name","value":"opts"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"total"}},{"kind":"Field","name":{"kind":"Name","value":"notifications"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"text"}},{"kind":"Field","name":{"kind":"Name","value":"userId"}},{"kind":"Field","name":{"kind":"Name","value":"twitchProfile"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"displayName"}},{"kind":"Field","name":{"kind":"Name","value":"profileImageUrl"}}]}},{"kind":"Field","name":{"kind":"Name","value":"createdAt"}}]}}]}}]}}]} as unknown as DocumentNode<NotificationsByAdminQuery, NotificationsByAdminQueryVariables>;
export const CreateNotificationDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"CreateNotification"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"text"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"String"}}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"userId"}},"type":{"kind":"NamedType","name":{"kind":"Name","value":"String"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"notificationsCreate"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"text"},"value":{"kind":"Variable","name":{"kind":"Name","value":"text"}}},{"kind":"Argument","name":{"kind":"Name","value":"userId"},"value":{"kind":"Variable","name":{"kind":"Name","value":"userId"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"text"}},{"kind":"Field","name":{"kind":"Name","value":"userId"}}]}}]}}]} as unknown as DocumentNode<CreateNotificationMutation, CreateNotificationMutationVariables>;
export const DeleteNotificationDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"DeleteNotification"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"id"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ID"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"notificationsDelete"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"id"},"value":{"kind":"Variable","name":{"kind":"Name","value":"id"}}}]}]}}]} as unknown as DocumentNode<DeleteNotificationMutation, DeleteNotificationMutationVariables>;
export const UpdateNotificationsDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"mutation","name":{"kind":"Name","value":"UpdateNotifications"},"variableDefinitions":[{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"id"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"ID"}}}},{"kind":"VariableDefinition","variable":{"kind":"Variable","name":{"kind":"Name","value":"opts"}},"type":{"kind":"NonNullType","type":{"kind":"NamedType","name":{"kind":"Name","value":"NotificationUpdateOpts"}}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"notificationsUpdate"},"arguments":[{"kind":"Argument","name":{"kind":"Name","value":"id"},"value":{"kind":"Variable","name":{"kind":"Name","value":"id"}}},{"kind":"Argument","name":{"kind":"Name","value":"opts"},"value":{"kind":"Variable","name":{"kind":"Name","value":"opts"}}}],"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"id"}},{"kind":"Field","name":{"kind":"Name","value":"text"}}]}}]}}]} as unknown as DocumentNode<UpdateNotificationsMutation, UpdateNotificationsMutationVariables>;
export const DashboardStatsDocument = {"kind":"Document","definitions":[{"kind":"OperationDefinition","operation":"subscription","name":{"kind":"Name","value":"dashboardStats"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"dashboardStats"},"selectionSet":{"kind":"SelectionSet","selections":[{"kind":"Field","name":{"kind":"Name","value":"categoryId"}},{"kind":"Field","name":{"kind":"Name","value":"categoryName"}},{"kind":"Field","name":{"kind":"Name","value":"viewers"}},{"kind":"Field","name":{"kind":"Name","value":"startedAt"}},{"kind":"Field","name":{"kind":"Name","value":"title"}},{"kind":"Field","name":{"kind":"Name","value":"chatMessages"}},{"kind":"Field","name":{"kind":"Name","value":"followers"}},{"kind":"Field","name":{"kind":"Name","value":"usedEmotes"}},{"kind":"Field","name":{"kind":"Name","value":"requestedSongs"}},{"kind":"Field","name":{"kind":"Name","value":"subs"}}]}}]}}]} as unknown as DocumentNode<DashboardStatsSubscription, DashboardStatsSubscriptionVariables>;