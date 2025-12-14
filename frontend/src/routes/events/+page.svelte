<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import { Button } from '$lib/components/ui/button';
	import * as Card from '$lib/components/ui/card';
	import { Badge } from '$lib/components/ui/badge';
	import * as Alert from '$lib/components/ui/alert';
	import { Separator } from '$lib/components/ui/separator';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import { user } from '$lib/stores/user';
	import { AVAILABLE_TAGS } from '$lib/constants/tags';
	import {
		IconSearch,
		IconX,
		IconFilter,
		IconChevronLeft,
		IconChevronRight
	} from '@tabler/icons-svelte';

	interface Event {
		event_id: string;
		name: string;
		description: string;
		date: string;
		fee: number;
		tag?: string[];
	}

	const api = import.meta.env.VITE_API_URL;
	const PAGE_SIZE = 12;

	let events = $state<Event[]>([]);
	let loading = $state(true);
	let error = $state('');

	let searchQuery = $state('');
	let selectedTags = $state<string[]>([]);
	let showFilters = $state(false);

	let dateFrom = $state('');
	let dateTo = $state('');
	let minFee = $state('');
	let maxFee = $state('');
	let freeOnly = $state(false);

	let currentPage = $state(1);
	let hasMorePages = $state(true);

	function initFiltersFromUrl() {
		const url = $page.url;
		searchQuery = url.searchParams.get('q') || '';
		const tagsParam = url.searchParams.get('tags');
		selectedTags = tagsParam ? tagsParam.split(',') : [];
		dateFrom = url.searchParams.get('date_from') || '';
		dateTo = url.searchParams.get('date_to') || '';
		minFee = url.searchParams.get('min_fee') || '';
		maxFee = url.searchParams.get('max_fee') || '';
		freeOnly = url.searchParams.get('free') === 'true';
		currentPage = parseInt(url.searchParams.get('page') || '1', 10);

		if (selectedTags.length > 0 || dateFrom || dateTo || minFee || maxFee || freeOnly) {
			showFilters = true;
		}
	}

	function updateUrl() {
		const params = new URLSearchParams();

		if (searchQuery.trim()) {
			params.set('q', searchQuery.trim());
		}
		if (selectedTags.length > 0) {
			params.set('tags', selectedTags.join(','));
		}
		if (dateFrom) {
			params.set('date_from', dateFrom);
		}
		if (dateTo) {
			params.set('date_to', dateTo);
		}
		if (freeOnly) {
			params.set('free', 'true');
		} else {
			if (minFee) {
				params.set('min_fee', minFee);
			}
			if (maxFee) {
				params.set('max_fee', maxFee);
			}
		}
		if (currentPage > 1) {
			params.set('page', currentPage.toString());
		}

		const queryString = params.toString();
		const newUrl = queryString ? `/events?${queryString}` : '/events';
		goto(newUrl, { replaceState: true, keepFocus: true });
	}

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
		dateFrom = '';
		dateTo = '';
		minFee = '';
		maxFee = '';
		freeOnly = false;
		currentPage = 1;
		updateUrl();
		fetchEvents();
	}

	async function fetchEvents() {
		loading = true;
		error = '';

		try {
			const params = new URLSearchParams();

			if (selectedTags.length > 0) {
				params.set('tags', selectedTags.join(','));
			}
			if (dateFrom) {
				params.set('date_from', dateFrom);
			}
			if (dateTo) {
				params.set('date_to', dateTo);
			}
			if (freeOnly) {
				params.set('min_fee', '0');
				params.set('max_fee', '0');
			} else {
				if (minFee) {
					params.set('min_fee', minFee);
				}
				if (maxFee) {
					params.set('max_fee', maxFee);
				}
			}

			params.set('page', currentPage.toString());
			params.set('page_size', PAGE_SIZE.toString());

			const queryString = params.toString();
			const url = `${api}/api/events?${queryString}`;

			const response = await fetch(url, {
				credentials: 'include'
			});
			if (response.ok) {
				const data = await response.json();
				events = data;
				hasMorePages = data.length === PAGE_SIZE;
			} else {
				error = 'Failed to load events';
			}
		} catch (err) {
			error = 'An error occurred while loading events';
		} finally {
			loading = false;
		}
	}

	async function goToPage(pageNum: number) {
		currentPage = pageNum;
		updateUrl();
		await fetchEvents();
		window.scrollTo({ top: 0, behavior: 'smooth' });
	}

	async function applyFilters() {
		currentPage = 1;
		updateUrl();
		await fetchEvents();
	}

	onMount(async () => {
		initFiltersFromUrl();
		await fetchEvents();
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

	function formatFee(fee: number): string {
		if (fee === 0) return 'Free';
		return `$${fee.toFixed(2)}`;
	}

	let hasActiveFilters = $derived(
		selectedTags.length > 0 || dateFrom || dateTo || minFee || maxFee || freeOnly
	);
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
				{#if hasActiveFilters}
					<Badge variant="secondary" class="ml-1">●</Badge>
				{/if}
			</Button>
			{#if searchQuery || hasActiveFilters}
				<Button variant="ghost" onclick={clearFilters} class="gap-2">
					<IconX class="w-4 h-4" />
					Clear
				</Button>
			{/if}
		</div>

		{#if showFilters}
			<Card.Root class="p-4">
				<div class="space-y-4">
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

					<Separator />

					<div class="space-y-3">
						<div class="text-sm font-medium">Filter by date</div>
						<div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
							<div class="space-y-2">
								<Label for="date-from">From</Label>
								<Input type="date" id="date-from" bind:value={dateFrom} />
							</div>
							<div class="space-y-2">
								<Label for="date-to">To</Label>
								<Input type="date" id="date-to" bind:value={dateTo} />
							</div>
						</div>
					</div>

					<Separator />

					<div class="space-y-3">
						<div class="text-sm font-medium">Filter by fee</div>

						<!-- Free only checkbox -->
						<label class="flex items-center gap-2 cursor-pointer">
							<input
								type="checkbox"
								bind:checked={freeOnly}
								class="w-4 h-4 rounded border-input accent-primary"
							/>
							<span class="text-sm">Free events only</span>
						</label>

						<!-- Fee range inputs (disabled when freeOnly is checked) -->
						<div class="grid grid-cols-1 sm:grid-cols-2 gap-4" class:opacity-50={freeOnly}>
							<div class="space-y-2">
								<Label for="min-fee">Min ($)</Label>
								<Input
									type="number"
									id="min-fee"
									placeholder="0"
									min="0"
									step="0.01"
									bind:value={minFee}
									disabled={freeOnly}
								/>
							</div>
							<div class="space-y-2">
								<Label for="max-fee">Max ($)</Label>
								<Input
									type="number"
									id="max-fee"
									placeholder="Any"
									min="0"
									step="0.01"
									bind:value={maxFee}
									disabled={freeOnly}
								/>
							</div>
						</div>
					</div>

					<Separator />

					<div class="flex justify-end gap-2">
						<Button variant="outline" onclick={clearFilters}>Clear All</Button>
						<Button onclick={applyFilters}>Apply Filters</Button>
					</div>
				</div>
			</Card.Root>
		{/if}

		{#if searchQuery || hasActiveFilters}
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
							<div class="flex items-center gap-3 text-sm">
								<span class="text-xl">💰</span>
								<span class="text-muted-foreground font-medium">{formatFee(event.fee)}</span>
							</div>
						</div>
					</Card.Content>
					<Card.Footer>
						<Button class="w-full" href={`/events/${event.event_id}`}>View Details</Button>
					</Card.Footer>
				</Card.Root>
			{/each}
		</div>

		<!-- Pagination Controls -->
		<div class="flex items-center justify-center gap-4 mt-8">
			<Button
				variant="outline"
				disabled={currentPage === 1}
				onclick={() => goToPage(currentPage - 1)}
				class="gap-2"
			>
				<IconChevronLeft class="w-4 h-4" />
				Previous
			</Button>
			<span class="text-sm text-muted-foreground">
				Page {currentPage}
			</span>
			<Button
				variant="outline"
				disabled={!hasMorePages}
				onclick={() => goToPage(currentPage + 1)}
				class="gap-2"
			>
				Next
				<IconChevronRight class="w-4 h-4" />
			</Button>
		</div>
	{/if}
</div>
