<script lang="ts">
	import { cn } from '$lib/utils';
	import type { HTMLAttributes } from 'svelte/elements';

	interface Props extends HTMLAttributes<HTMLDivElement> {
		latitude?: number;
		longitude?: number;
		onLocationChange?: (lat: number, lng: number) => void;
		height?: string;
	}

	let {
		class: className,
		latitude = 51.505,
		longitude = -0.09,
		onLocationChange,
		height = '400px',
		...restProps
	}: Props = $props();

	let mapContainer = $state<HTMLDivElement | null>(null);
	let map = $state<any>(null);
	let marker = $state<any>(null);
	let selectedLat = $state(latitude);
	let selectedLng = $state(longitude);
	let isInitialized = $state(false);

	$effect(() => {
		if (!mapContainer || isInitialized) return;

		(async () => {
			const L = await import('leaflet');

			const link = document.createElement('link');
			link.rel = 'stylesheet';
			link.href = 'https://unpkg.com/leaflet@1.9.4/dist/leaflet.css';
			document.head.appendChild(link);

			await new Promise((resolve) => setTimeout(resolve, 100));

			map = L.map(mapContainer).setView([latitude, longitude], 13);

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
				draggable: true,
				icon: customIcon
			}).addTo(map);

			marker.on('dragend', function (e: any) {
				const position = marker.getLatLng();
				selectedLat = position.lat;
				selectedLng = position.lng;
				onLocationChange?.(position.lat, position.lng);
			});

			map.on('click', function (e: any) {
				marker.setLatLng(e.latlng);
				selectedLat = e.latlng.lat;
				selectedLng = e.latlng.lng;
				onLocationChange?.(e.latlng.lat, e.latlng.lng);
			});

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
			map.setView([latitude, longitude], 13);
			selectedLat = latitude;
			selectedLng = longitude;
		}
	});
</script>

<div class={cn('space-y-2', className)} {...restProps}>
	<div
		bind:this={mapContainer}
		style="height: {height}; width: 100%; min-height: {height};"
		class="rounded-md border border-input overflow-hidden shadow-sm"
	></div>
</div>
