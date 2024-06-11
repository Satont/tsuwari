export interface Track {
	artist: string
	title: string
	imageUrl?: string | null
}

export const Preset = {
	TRANSPARENT: 'TRANSPARENT',
	AIDEN_REDESIGN: 'AIDEN_REDESIGN',
	SIMPLE_LINE: 'SIMPLE_LINE',
} as const

export interface Settings {
	id: string
	preset: keyof typeof Preset
	fontFamily: string
	fontWeight: number
	backgroundColor: string
	showImage: boolean
	hideTimeout?: number | null
}
