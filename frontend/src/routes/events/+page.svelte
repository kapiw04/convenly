<script lang="ts">
	import { onMount } from 'svelte';
	import { Button } from "$lib/components/ui/button";
	import * as Card from "$lib/components/ui/card";
	import { Badge } from "$lib/components/ui/badge";
	import * as Alert from "$lib/components/ui/alert";
	import { Separator } from "$lib/components/ui/separator";

	interface Event {
		name: string;
		description: string;
		date: string;
	}
	
	const api = import.meta.env.VITE_API_URL;
	let events = $state<Event[]>([]);
	let loading = $state(true);
	let error = $state('');

	onMount(async () => {
		try {
			const response = await fetch(`${api}/api/events`, {
				credentials: 'include'
			});
			if (response.ok) {
				events = await response.json();
				console.log(events);
			} else {
				error = 'Failed to load events';
			}
		} catch (err) {
			error = 'An error occurred while loading events';
		} finally {
			loading = false;
		}
	});

	function formatDate(dateString: string): string {
		const date = new Date(dateString);
		return date.toLocaleDateString('en-US', { 
			year: 'numeric', 
			month: 'short', 
			day: 'numeric',
			hour: '2-digit',
			minute: '2-digit'
		});
	}
</script>

<div class="container mx-auto px-4 py-8">
	<div class="mb-10 space-y-3">
		<h1 class="text-4xl font-bold tracking-tight">Upcoming Events</h1>
	</div>

	{#if loading}
		<div class="flex flex-col items-center justify-center h-64 space-y-4">
			<div class="animate-spin rounded-full h-12 w-12 border-b-2 border-primary"></div>
			<p class="text-muted-foreground">Loading events...</p>
		</div>
	{:else if error}
		<Alert.Root variant="destructive">
			<Alert.Description>{error}</Alert.Description>
		</Alert.Root>
	{:else if events.length === 0}
		<Card.Root class="border-dashed">
			<Card.Content class="flex flex-col items-center justify-center py-16">
				<Card.Title class="mb-2">No events available</Card.Title>
				<Card.Description>
					Check back later for upcoming events!
				</Card.Description>
			</Card.Content>
		</Card.Root>
	{:else}
		<div class="grid md:grid-cols-2 lg:grid-cols-3 gap-6">
			{#each events as event}
				<Card.Root class="hover:shadow-lg transition-shadow duration-200 flex flex-col">
					<Card.Header>
						<div class="flex items-start justify-between gap-2">
							<Card.Title class="text-xl line-clamp-2">
								{event.name}
							</Card.Title>
							<Badge variant="secondary">Active</Badge>
						</div>
						<Card.Description class="line-clamp-3 text-base">
							{event.description}
						</Card.Description>
					</Card.Header>
					<Card.Content class="flex-grow space-y-3">
						<Separator />
						<div class="space-y-2.5">
							<div class="flex items-center gap-3 text-sm">
								<span class="text-xl">🕐</span>
								<span class="text-muted-foreground">{formatDate(event.date)}</span>
							</div>
						</div>
					</Card.Content>
					<Card.Footer>
						<Button class="w-full">
							View Details
						</Button>
					</Card.Footer>
				</Card.Root>
			{/each}
		</div>
	{/if}
</div>
