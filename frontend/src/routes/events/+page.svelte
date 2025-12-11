<script lang="ts">
	import { onMount } from 'svelte';
	import { Button } from '$lib/components/ui/button';
	import * as Card from '$lib/components/ui/card';
	import { Badge } from '$lib/components/ui/badge';
	import * as Alert from '$lib/components/ui/alert';
	import { Separator } from '$lib/components/ui/separator';
	import { Input } from '$lib/components/ui/input';
	import { user } from '$lib/stores/user';
	import { AVAILABLE_TAGS } from '$lib/constants/tags';
	import { IconSearch, IconX, IconFilter } from '@tabler/icons-svelte';

	interface Event {
		event_id: string;
		name: string;
		description: string;
		date: string;
		tag?: string[];
	}

	const api = import.meta.env.VITE_API_URL;
	let events = $state<Event[]>([]);
	let loading = $state(true);
	let error = $state('');

	let searchQuery = $state('');
	let selectedTags = $state<string[]>([]);
	let showFilters = $state(false);

	let filteredEvents = $derived(() => {
		let result = events;

		if (searchQuery.trim()) {
			const query = searchQuery.toLowerCase();
			result = result.filter(
				(event) =>
					event.name.toLowerCase().includes(query) ||
					event.description.toLowerCase().includes(query)
			);
		}

		if (selectedTags.length > 0) {
			result = result.filter(
				(event) => event.tag && event.tag.some((t) => selectedTags.includes(t))
			);
		}

		return result;
	});

	function toggleTagFilter(tag: string) {
		if (selectedTags.includes(tag)) {
			selectedTags = selectedTags.filter((t) => t !== tag);
		} else {
			selectedTags = [...selectedTags, tag];
		}
	}

	function clearFilters() {
		searchQuery = '';
		selectedTags = [];
	}

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
	{#if $user && $user.role === 1}
		<div class="mb-6">
			<Button href="/events/create" size="lg">Create New Event</Button>
		</div>
	{/if}

	<div class="mb-6 space-y-3">
		<h1 class="text-4xl font-bold tracking-tight">Upcoming Events</h1>
	</div>

	<div class="mb-6 space-y-4">
		<div class="flex gap-3 flex-wrap">
			<div class="relative flex-1 min-w-[200px] max-w-md">
				<IconSearch
					class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-muted-foreground"
				/>
				<Input
					type="text"
					placeholder="Search events..."
					class="pl-9 pr-9"
					bind:value={searchQuery}
				/>
				{#if searchQuery}
					<button
						onclick={() => (searchQuery = '')}
						class="absolute right-3 top-1/2 -translate-y-1/2 text-muted-foreground hover:text-foreground"
					>
						<IconX class="w-4 h-4" />
					</button>
				{/if}
			</div>
			<Button
				variant={showFilters ? 'default' : 'outline'}
				onclick={() => (showFilters = !showFilters)}
				class="gap-2"
			>
				<IconFilter class="w-4 h-4" />
				Filters
				{#if selectedTags.length > 0}
					<Badge variant="secondary" class="ml-1">{selectedTags.length}</Badge>
				{/if}
			</Button>
			{#if searchQuery || selectedTags.length > 0}
				<Button variant="ghost" onclick={clearFilters} class="gap-2">
					<IconX class="w-4 h-4" />
					Clear
				</Button>
			{/if}
		</div>

		{#if showFilters}
			<Card.Root class="p-4">
				<div class="space-y-3">
					<div class="text-sm font-medium">Filter by tags</div>
					<div class="flex flex-wrap gap-2">
						{#each AVAILABLE_TAGS as tag}
							<button
								onclick={() => toggleTagFilter(tag)}
								class="px-3 py-1.5 text-sm rounded-full border transition-colors {selectedTags.includes(
									tag
								)
									? 'bg-primary text-primary-foreground border-primary'
									: 'bg-background hover:bg-muted border-input'}"
							>
								{tag}
							</button>
						{/each}
					</div>
				</div>
			</Card.Root>
		{/if}

		{#if searchQuery || selectedTags.length > 0}
			<div class="text-sm text-muted-foreground">
				Showing {filteredEvents().length} of {events.length} events
			</div>
		{/if}
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
	{:else if filteredEvents().length === 0}
		<Card.Root class="border-dashed">
			<Card.Content class="flex flex-col items-center justify-center py-16">
				<Card.Title class="mb-2">
					{events.length === 0 ? 'No events available' : 'No events match your filters'}
				</Card.Title>
				<Card.Description>
					{events.length === 0
						? 'Check back later for upcoming events!'
						: 'Try adjusting your search or filters'}
				</Card.Description>
				{#if events.length > 0}
					<Button variant="outline" onclick={clearFilters} class="mt-4">Clear Filters</Button>
				{/if}
			</Card.Content>
		</Card.Root>
	{:else}
		<div class="grid md:grid-cols-2 lg:grid-cols-3 gap-6">
			{#each filteredEvents() as event}
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
						{#if event.tag && event.tag.length > 0}
							<div class="flex flex-wrap gap-1.5">
								{#each event.tag as tag}
									<Badge variant="outline" class="text-xs">{tag}</Badge>
								{/each}
							</div>
						{/if}
						<Separator />
						<div class="space-y-2.5">
							<div class="flex items-center gap-3 text-sm">
								<span class="text-xl">🕐</span>
								<span class="text-muted-foreground">{formatDate(event.date)}</span>
							</div>
						</div>
					</Card.Content>
					<Card.Footer>
						<Button class="w-full" href={`/events/${event.event_id}`}>View Details</Button>
					</Card.Footer>
				</Card.Root>
			{/each}
		</div>
	{/if}
</div>
