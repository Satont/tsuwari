import {
	type ColumnDef,
	getCoreRowModel,
	getFacetedRowModel,
	getFacetedUniqueValues,
	getPaginationRowModel,
	getSortedRowModel,
	useVueTable,
} from '@tanstack/vue-table'
import { defineStore, storeToRefs } from 'pinia'
import { computed, h } from 'vue'
import { useI18n } from 'vue-i18n'

import { useCommunityEmotesStatisticFilters } from './use-community-emotes-statistic-filters.js'
import CommunityEmotesTableColumnChartRange
	from '../components/community-emotes-table-column-chart-range.vue'
import CommunityEmotesTableColumnChart from '../components/community-emotes-table-column-chart.vue'
import CommunityEmotesTableColumn from '../components/community-emotes-table-column.vue'

import { type EmotesStatistics, useEmotesStatisticQuery } from '@/api/emotes-statistic.js'
import { usePagination } from '@/composables/use-pagination.js'
import CommunityEmotesTableColumnActions
	from '@/features/community-emotes-statistic/components/community-emotes-table-column-actions.vue'
import { EmoteStatisticRange, type EmotesStatisticsOpts } from '@/gql/graphql'
import { valueUpdater } from '@/helpers/value-updater.js'

export const useCommunityEmotesStatisticTable = defineStore('features/community-emotes-statistic-table', () => {
	const { t } = useI18n()
	const { pagination } = usePagination()
	const {
		debouncedSearchInput,
		tableRange,
		sortingState,
		tableOrder,
	} = storeToRefs(useCommunityEmotesStatisticFilters())

	const emotesQueryOptions = computed<EmotesStatisticsOpts>((prevParams) => {
		if (prevParams?.search !== debouncedSearchInput.value) {
			pagination.value.pageIndex = 0
		}

		return {
			search: debouncedSearchInput.value,
			perPage: pagination.value.pageSize,
			page: pagination.value.pageIndex,
			graphicRange: tableRange.value,
			order: tableOrder.value,
		}
	})
	const { data, fetching } = useEmotesStatisticQuery(emotesQueryOptions)

	const emotes = computed<EmotesStatistics>(() => {
		if (!data.value) return []
		return data.value.emotesStatistics.emotes
	})
	const totalEmotes = computed(() => data.value?.emotesStatistics.total ?? 0)
	const pageCount = computed(() => {
		return Math.ceil(totalEmotes.value / pagination.value.pageSize)
	})

	const statsColumn = computed<ColumnDef<EmotesStatistics[0]>[]>(() => [
		{
			accessorKey: 'name',
			size: 5,
			header: () => h('div', {}, t('community.emotesStatistic.table.emote')),
			cell: ({ row }) => {
				return h('div', { class: 'break-words max-w-[450px]', innerHTML: row.original.emoteName })
			},
		},
		{
			accessorKey: 'usages',
			size: 5,
			header: ({ column }) => {
				return h(CommunityEmotesTableColumn, {
					column,
					title: t('community.emotesStatistic.table.usages'),
				})
			},
			cell: ({ row }) => {
				return h('div', `${row.original.totalUsages}`)
			},
		},
		{
			accessorKey: 'chart',
			size: 80,
			header: () => h(CommunityEmotesTableColumnChartRange),
			cell: ({ row }) => {
				return h(CommunityEmotesTableColumnChart, {
					isDayRange: tableRange.value === EmoteStatisticRange.LastDay,
					usages: row.original.graphicUsages,
				})
			},
		},
		{
			accessorKey: 'actions',
			size: 10,
			header: () => '',
			cell: ({ row }) => {
				return h(CommunityEmotesTableColumnActions, { emoteName: row.original.emoteName })
			},
		},
	])

	const table = useVueTable({
		get pageCount() {
			return pageCount.value
		},
		get data() {
			return emotes.value
		},
		get columns() {
			return statsColumn.value
		},
		state: {
			get sorting() {
				return sortingState.value
			},
			get pagination() {
				return pagination.value
			},
		},
		manualPagination: true,
		onPaginationChange: (updaterOrValue) => valueUpdater(updaterOrValue, pagination),
		getCoreRowModel: getCoreRowModel(),
		getPaginationRowModel: getPaginationRowModel(),
		onSortingChange: updaterOrValue => valueUpdater(updaterOrValue, sortingState),
		getSortedRowModel: getSortedRowModel(),
		getFacetedRowModel: getFacetedRowModel(),
		getFacetedUniqueValues: getFacetedUniqueValues(),
	})

	return {
		isLoading: fetching,
		table,
		totalEmotes,
		pageCount,
		pagination,
	}
})
