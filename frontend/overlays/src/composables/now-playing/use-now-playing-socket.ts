import { useSubscription } from '@urql/vue'
import { ref, watch } from 'vue'

import type { Settings, Track } from '@twir/frontend-now-playing'

import { graphql } from '@/gql'

interface Options {
	apiKey: string
	overlayId: string
}

export function useNowPlayingSocket(options: Options) {
	const settingsSub = useSubscription({
		query: graphql(`
			subscription NowPlayingOverlaySettings($overlayId: String!, $apiKey: String!) {
				nowPlayingOverlaySettings(id: $overlayId, apiKey: $apiKey) {
					id
					backgroundColor
					channelId
					fontFamily
					fontWeight
					hideTimeout
					preset
					showImage
				}
			}
		`),
		variables: {
			apiKey: options.apiKey,
			overlayId: options.overlayId,
		},
	})

	const currentTrackSub = useSubscription({
		query: graphql(`
			subscription NowPlayingOverlayNowPlaying($apiKey: String!) {
				nowPlayingCurrentTrack(apiKey: $apiKey) {
					title
					artist
					imageUrl
					progressMs
					durationMs
				}
			}
		`),
		variables: {
			apiKey: options.apiKey,
		},
	})

	const currentTrack = ref<Track | null | undefined>()
	const settings = ref<Settings>()

	watch(currentTrackSub.data, (data) => {
		currentTrack.value = data?.nowPlayingCurrentTrack
	}, { immediate: true })

	watch(settingsSub.data, (data) => {
		if (!data?.nowPlayingOverlaySettings) return

		settings.value = data.nowPlayingOverlaySettings
	}, { immediate: true })

	return {
		currentTrack,
		settings,
	}
}
