import { Fragment, useCallback, useEffect, useRef, useState } from 'react';
import { useParams } from 'react-router-dom';
import useWebSocket from 'react-use-websocket';

declare global {
	interface Window {
		webkitAudioContext: typeof AudioContext
	}
}

interface Layer {
	id: string
	type: 'HTML'
	settings: LayerSettings
	overlay_id: string
	pos_x: number
	pos_y: number
	width: number
	height: number
	createdAt: string
	updatedAt: string
	overlay: any

	htmlContent?: string;
}

export interface LayerSettings {
	htmlOverlayDataPollSecondsInterval: number
	htmlOverlayHtml: string
	htmlOverlayCss: string
	htmlOverlayJs: string
}

function fromBase64(str: string) {
	return decodeURIComponent(atob(str).split('').map(function(c) {
		return '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2);
	}).join(''));
}

export const OverlaysRegistry: React.FC = () => {
	const [url, setUrl] = useState<string | null>(null);
	const { apiKey, overlayId } = useParams();
	const contentRef = useRef<HTMLDivElement>(null);

	const { lastMessage, sendMessage } = useWebSocket(url, {
		shouldReconnect: () => true,
		onOpen: () => {
			sendMessage(JSON.stringify({
				eventName: 'getLayers',
				data: {
					overlayId,
				},
			}));
		},
		reconnectInterval: 500,
	});

	const [layers, setLayers] = useState<Layer[]>([]);

	useEffect(() => {
		if (!apiKey) return;

		const protocol = window.location.protocol === 'https:' ? 'wss' : 'ws';
		const host = window.location.host;

		setUrl(`${protocol}://${host}/socket/registry/overlays?apiKey=${apiKey}`);
	}, [apiKey]);

	useEffect(() => {
		if (!lastMessage) return;
		try {
			const parsedData = JSON.parse(lastMessage.data);

			if (parsedData.eventName === 'refreshOverlays') {
				window.location.reload();
			}

			if (parsedData.eventName === 'layers') {
				setLayers(parsedData.layers);
				for (const layer of parsedData.layers) {
					if (layer.type === 'HTML') {
						pollHtmlOverlayData(layer as Layer);
					}
				}
			}

			if (parsedData.eventName === 'parsedLayerVariables') {
				processParsedLayerVariables(parsedData);
			}
		} catch (e) {
			console.error('cannot parse message', lastMessage.data);
		}
	}, [lastMessage]);

	const processParsedLayerVariables = useCallback((parsedData: any) => {
		const layer = layers.find((l) => l.id === parsedData.layerId);
		if (!layer) return;

		setLayers((prevLayers) => {
			return prevLayers.map((l) => {
				if (l.id !== parsedData.layerId) return l;
				return {
					...l,
					htmlContent: parsedData.data,
				};
			});
		});
	}, [layers]);

	const pollHtmlOverlayData = useCallback((l: Layer) => {
		if (l.type !== 'HTML') return;
		if (l.settings.htmlOverlayDataPollSecondsInterval <= 0) return;

		const getInfo = () => sendMessage(JSON.stringify({
			eventName: 'parseLayerVariables',
			data: {
				layerId: l.id,
			},
		}));
		getInfo();

		const interval = setInterval(() => {
			getInfo();
		}, l.settings.htmlOverlayDataPollSecondsInterval * 1000);

		return () => {
			clearInterval(interval);
		};
	}, [layers]);

	return <div ref={contentRef} style={{
		// aspectRatio: '16 / 9',
		width: '100%',
		height: '100%',
		// overflow: 'hidden',
	}}>
		{layers.filter(l => l.type === 'HTML').map((layer) => {
			return <Fragment key={layer.id}>
				<style>
					{`.layer-${layer.id} {
						${fromBase64(layer.settings.htmlOverlayCss)}
					}`}
				</style>
				<div
					key={layer.id}
					style={{
						position: 'absolute',
						top: layer.pos_y,
						left: layer.pos_x,
						width: layer.width,
						height: layer.height,
						overflow: 'hidden',
						// eslint-disable-next-line @typescript-eslint/ban-ts-comment
						// @ts-ignore
						textWrap: 'nowrap',
					}}
					className={'layer-' + layer.id}
					dangerouslySetInnerHTML={{ __html: layer.htmlContent ? fromBase64(layer.htmlContent) : '' }}
				/>
			</Fragment>;
		})}
	</div>;
};
