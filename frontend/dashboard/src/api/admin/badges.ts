import { useQuery } from '@urql/vue'

import { useMutation } from '@/composables/use-mutation.js'
import { graphql } from '@/gql'

const invalidationKey = 'AdminBadgesInvalidateKey'

export function useAdminBadges() {
	const useMutationCreateBadge = () => useMutation(graphql(`
		mutation CreateBadge($opts: TwirBadgeCreateOpts!) {
			badgesCreate(opts: $opts) {
				id
			}
		}
	`), [invalidationKey])

	const useMutationDeleteBadge = () => useMutation(graphql(`
		mutation DeleteBadge($id: ID!) {
			badgesDelete(id: $id)
		}
	`), [invalidationKey])

	const useMutationUpdateBadge = () => useMutation(graphql(`
		mutation UpdateBadge($id: ID!, $opts: TwirBadgeUpdateOpts!) {
			badgesUpdate(id: $id, opts: $opts) {
				id
			}
		}
	`), [invalidationKey])

	const useMutationsAddUserBadge = () => useMutation(graphql(`
		mutation AddUserBadge($id: ID!, $userId: String!) {
			badgesAddUser(id: $id, userId: $userId)
		}
	`), [invalidationKey])

	const useMutationsRemoveUserBadge = () => useMutation(graphql(`
		mutation RemoveUserBadge($id: ID!, $userId: String!) {
			badgesRemoveUser(id: $id, userId: $userId)
		}
	`), [invalidationKey])

	return {
		useMutationCreateBadge,
		useMutationDeleteBadge,
		useMutationUpdateBadge,
		useMutationsAddUserBadge,
		useMutationsRemoveUserBadge
	}
}

export function useQueryBadges() {
	return useQuery({
		context: {
			additionalTypenames: [invalidationKey]
		},
		variables: {},
		query: graphql(`
			query BadgesGetAll {
				twirBadges {
					id
					name
					createdAt
					fileUrl
					enabled
					ffzSlot
					users
				}
			}
		`)
	})
}
