<script setup lang="ts">
import { type Font, FontSelector } from '@twir/fontsource'
import { NButton, NColorPicker, NFormItem, NInputNumber, NSelect, NSwitch, useThemeVars } from 'naive-ui'
import { computed, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'

import { useNowPlayingForm } from './use-now-playing-form'

import {
	useNowPlayingOverlayApi,
	useProfile,
	useUserAccessFlagChecker,
} from '@/api'
import { useCopyOverlayLink } from '@/components/overlays/copyOverlayLink'
import { useNaiveDiscrete } from '@/composables/use-naive-discrete'
import { ChannelRolePermissionEnum } from '@/gql/graphql'

const { t } = useI18n()

const themeVars = useThemeVars()
const discrete = useNaiveDiscrete()
const { copyOverlayLink } = useCopyOverlayLink('now-playing')
const userCanEditOverlays = useUserAccessFlagChecker(ChannelRolePermissionEnum.ManageOverlays)

const { data: profile } = useProfile()
const { data: formValue } = useNowPlayingForm()

const canCopyLink = computed(() => {
	return profile?.value?.selectedDashboardId === profile.value?.id && userCanEditOverlays
})

const manager = useNowPlayingOverlayApi()
const updater = manager.useNowPlayingUpdate()

async function save() {
	if (!formValue.value?.id) return

	await updater.executeMutation({
		id: formValue.value.id,
		input: {
			preset: formValue.value.preset,
			fontFamily: formValue.value.fontFamily,
			fontWeight: formValue.value.fontWeight,
			backgroundColor: formValue.value.backgroundColor,
			showImage: formValue.value.showImage,
			hideTimeout: formValue.value.hideTimeout,
		},
	})

	discrete.notification.success({
		title: t('sharedTexts.saved'),
		duration: 1500,
	})
}

const fontData = ref<Font | null>(null)
watch(() => fontData.value, (font) => {
	console.log(font)
	if (!font) return
	formValue.value.fontFamily = font.id
}, { deep: true })

const fontWeightOptions = computed(() => {
	if (!fontData.value) return []
	return fontData.value.weights.map((weight) => ({ label: `${weight}`, value: weight }))
})
</script>

<template>
	<div v-if="formValue" class="card">
		<div class="card-header">
			<NButton
				secondary
				type="info"
				:disabled="!formValue.id || !canCopyLink"
				@click="copyOverlayLink({ id: formValue.id! })"
			>
				{{ t('overlays.copyOverlayLink') }}
			</NButton>
			<NButton secondary type="success" @click="save">
				{{ t('sharedButtons.save') }}
			</NButton>
		</div>

		<div class="card-body-column">
			<NFormItem label="Style">
				<NSelect
					v-model:value="formValue.preset"
					:options="[
						{ label: 'Aiden Redesign', value: 'AIDEN_REDESIGN' },
						{ label: 'Transparent', value: 'TRANSPARENT' },
						{ label: 'Simple line', value: 'SIMPLE_LINE' },
						{ label: 'Compact', value: 'COMPACT' },
					]"
				/>
			</NFormItem>

			<NFormItem label="Show image">
				<NSwitch v-model:value="formValue.showImage" />
			</NFormItem>

			<NFormItem label="Background color">
				<NColorPicker
					v-model:value="formValue.backgroundColor"
				/>
			</NFormItem>

			<NFormItem :label="t('overlays.chat.fontFamily')">
				<FontSelector
					v-model:font="fontData"
					:font-family="formValue.fontFamily"
					:font-weight="formValue.fontWeight"
					font-style="normal"
				/>
			</NFormItem>

			<NFormItem :label="t('overlays.chat.fontWeight')">
				<NSelect
					v-model:value="formValue.fontWeight"
					:options="fontWeightOptions"
				/>
			</NFormItem>

			<NFormItem :label="t('overlays.chat.hideTimeout')">
				<NInputNumber
					v-model:value="formValue.hideTimeout"
					:min="0"
					:max="600"
				/>
			</NFormItem>
		</div>
	</div>
</template>

<style scoped>
@import '../styles.css';

.card {
	background-color: v-bind('themeVars.cardColor');
}
</style>
