<script lang="ts">
	import { cn } from '$lib/utils';
	import type { HTMLAttributes } from 'svelte/elements';

	interface Props extends HTMLAttributes<HTMLDivElement> {
		latitude: number;
		longitude: number;
		height?: string;
		zoom?: number;
	}

	let {
		class: className,
		latitude,
		longitude,
		height = '300px',
		zoom = 13,
		...restProps
	}: Props = $props();

	let mapContainer = $state<HTMLDivElement | null>(null);
	let map = $state<any>(null);
	let marker = $state<any>(null);
	let isInitialized = $state(false);

	$effect(() => {
		if (!mapContainer || isInitialized) return;

		(async () => {
			const L = await import('leaflet');

			if (!document.querySelector('link[href*="leaflet.css"]')) {
				const link = document.createElement('link');
				link.rel = 'stylesheet';
				link.href = 'https://unpkg.com/leaflet@1.9.4/dist/leaflet.css';
				document.head.appendChild(link);
			}

			await new Promise((resolve) => setTimeout(resolve, 100));

			map = L.map(mapContainer, {
				dragging: true,
				touchZoom: true,
				scrollWheelZoom: true,
				doubleClickZoom: true,
				boxZoom: true
			}).setView([latitude, longitude], zoom);

			L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
				attribution:
					'&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors',
				maxZoom: 19
			}).addTo(map);

			const customIcon = L.icon({
				iconUrl: 'https://unpkg.com/leaflet@1.9.4/dist/images/marker-icon.png',
				iconRetinaUrl: 'https://unpkg.com/leaflet@1.9.4/dist/images/marker-icon-2x.png',
				shadowUrl: 'https://unpkg.com/leaflet@1.9.4/dist/images/marker-shadow.png',
				iconSize: [25, 41],
				iconAnchor: [12, 41],
				popupAnchor: [1, -34],
				shadowSize: [41, 41]
			});

			marker = L.marker([latitude, longitude], {
				draggable: false,
				icon: customIcon
			}).addTo(map);

			marker.bindPopup(
				`<b>Event Location</b><br>Lat: ${latitude.toFixed(4)}<br>Lng: ${longitude.toFixed(4)}`
			);

			setTimeout(() => {
				map.invalidateSize();
			}, 250);

			isInitialized = true;
		})();

		return () => {
			if (map) {
				map.remove();
				map = null;
				marker = null;
				isInitialized = false;
			}
		};
	});

	$effect(() => {
		if (map && marker && isInitialized) {
			marker.setLatLng([latitude, longitude]);
			map.setView([latitude, longitude], zoom);
			marker.setPopupContent(
				`<b>Event Location</b><br>Lat: ${latitude.toFixed(4)}<br>Lng: ${longitude.toFixed(4)}`
			);
		}
	});
</script>

<div class={cn('space-y-2', className)} {...restProps}>
	<div
		bind:this={mapContainer}
		style="height: {height}; width: 100%; min-height: {height};"
		class="rounded-lg border border-input overflow-hidden shadow-md"
	></div>
</div>
