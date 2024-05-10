import { type ColumnDef, getCoreRowModel, useVueTable } from '@tanstack/vue-table'
import { createGlobalState } from '@vueuse/core'
import { computed, h } from 'vue'
import { useI18n } from 'vue-i18n'

import GreetingsTableActions from '../ui/greetings-table-actions.vue'

import { type Greetings, useGreetingsApi } from '@/api/greetings.js'
import UsersTableCellUser from '@/features/admin-panel/manage-users/components/users-table-cell-user.vue'

export const useGreetingsTable = createGlobalState(() => {
	const { t } = useI18n()
	const greetingsApi = useGreetingsApi()

	const { data, fetching } = greetingsApi.useQueryGreetings()
	const greetings = computed<Greetings[]>(() => {
		if (!data.value) return []
		return data.value.greetings
	})

	const tableColumns = computed<ColumnDef<Greetings>[]>(() => [
		{
			accessorKey: 'user',
			size: 60,
			header: () => h('div', {}, t('sharedTexts.user')),
			cell: ({ row }) => {
				return h('a', {
					class: 'flex flex-col',
					href: `https://twitch.tv/${row.original.twitchProfile.login}`,
					target: '_blank',
				}, h(UsersTableCellUser, {
					avatar: row.original.twitchProfile.profileImageUrl,
					name: row.original.twitchProfile.login,
					displayName: row.original.twitchProfile.displayName,
				}))
			},
		},
		{
			accessorKey: 'text',
			size: 30,
			header: () => h('div', {}, t('sharedTexts.response')),
			cell: ({ row }) => h('span', row.original.text),
		},
		{
			accessorKey: 'actions',
			size: 10,
			header: () => '',
			cell: ({ row }) => h(GreetingsTableActions, { greetings: row.original }),
		},
	])

	const table = useVueTable({
		get data() {
			return greetings.value
		},
		get columns() {
			return tableColumns.value
		},
		getCoreRowModel: getCoreRowModel(),
	})

	return {
		isLoading: fetching,
		table,
	}
})
