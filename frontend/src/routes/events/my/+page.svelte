<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { Button } from '$lib/components/ui/button';
	import * as Card from '$lib/components/ui/card';
	import { Badge } from '$lib/components/ui/badge';
	import * as Alert from '$lib/components/ui/alert';
	import { Separator } from '$lib/components/ui/separator';
	import { user } from '$lib/stores/user';
	import {
		IconCalendarEvent,
		IconTicket,
		IconSpeakerphone,
		IconAlertCircle,
		IconArrowLeft
	} from '@tabler/icons-svelte';

	interface Event {
		event_id: string;
		name: string;
		description: string;
		date: string;
		fee: number;
		tag?: string[];
	}

	interface MyEventsResponse {
		hosting: Event[];
		attending: Event[];
	}

	const api = import.meta.env.VITE_API_URL;
	let hosting = $state<Event[]>([]);
	let attending = $state<Event[]>([]);
	let loading = $state(true);
	let error = $state('');
	let activeTab = $state<'attending' | 'hosting'>('attending');

	onMount(async () => {
		await fetchMyEvents();
	});

	async function fetchMyEvents() {
		loading = true;
		error = '';

		try {
			const response = await fetch(`${api}/api/my-events`, {
				credentials: 'include'
			});
			if (response.ok) {
				const data: MyEventsResponse = await response.json();
				hosting = data.hosting || [];
				attending = data.attending || [];
			} else if (response.status === 401) {
				error = 'Please login to view your events';
				setTimeout(() => goto('/login'), 2000);
			} else {
				error = 'Failed to load your events';
			}
		} catch (err) {
			error = 'An error occurred while loading your events';
		} finally {
			loading = false;
		}
	}

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
</script>

<div class="container mx-auto px-4 py-8">
	<div class="mb-6">
		<Button href="/events" variant="ghost" size="sm" class="gap-2">
			<IconArrowLeft class="w-4 h-4" />
			Back to All Events
		</Button>
	</div>

	<div class="mb-6 space-y-3">
		<h1 class="text-4xl font-bold tracking-tight flex items-center gap-3">
			<IconCalendarEvent class="w-10 h-10 text-primary" />
			My Events
		</h1>
		<p class="text-muted-foreground">View events you're hosting and attending</p>
	</div>

	{#if loading}
		<div class="flex flex-col items-center justify-center h-64 space-y-4">
			<div class="animate-spin rounded-full h-12 w-12 border-b-2 border-primary"></div>
			<p class="text-muted-foreground">Loading your events...</p>
		</div>
	{:else if error}
		<Alert.Root variant="destructive">
			<IconAlertCircle class="h-4 w-4" />
			<Alert.Description>{error}</Alert.Description>
		</Alert.Root>
	{:else}
		<div class="mb-6">
			<div class="inline-flex rounded-lg border bg-muted p-1">
				<button
					class="inline-flex items-center gap-2 px-4 py-2 text-sm font-medium rounded-md transition-colors {activeTab ===
					'attending'
						? 'bg-background shadow-sm'
						: 'hover:bg-background/50'}"
					onclick={() => (activeTab = 'attending')}
				>
					<IconTicket class="w-4 h-4" />
					Attending ({attending.length})
				</button>
				<button
					class="inline-flex items-center gap-2 px-4 py-2 text-sm font-medium rounded-md transition-colors {activeTab ===
					'hosting'
						? 'bg-background shadow-sm'
						: 'hover:bg-background/50'}"
					onclick={() => (activeTab = 'hosting')}
				>
					<IconSpeakerphone class="w-4 h-4" />
					Hosting ({hosting.length})
				</button>
			</div>
		</div>

		{#if activeTab === 'attending'}
			{#if attending.length === 0}
				<Card.Root class="border-dashed">
					<Card.Content class="flex flex-col items-center justify-center py-16">
						<IconTicket class="w-12 h-12 text-muted-foreground mb-4" />
						<Card.Title class="mb-2">No events yet</Card.Title>
						<Card.Description class="text-center mb-4">
							You haven't registered for any events yet. Browse events to find something
							interesting!
						</Card.Description>
						<Button href="/events">Browse Events</Button>
					</Card.Content>
				</Card.Root>
			{:else}
				<div class="grid md:grid-cols-2 lg:grid-cols-3 gap-6">
					{#each attending as event}
						<Card.Root class="hover:shadow-lg transition-shadow duration-200 flex flex-col">
							<Card.Header>
								<div class="flex items-start justify-between gap-2">
									<Card.Title class="text-xl line-clamp-2">
										{event.name}
									</Card.Title>
									<Badge variant="secondary">Attending</Badge>
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
			{/if}
		{:else if hosting.length === 0}
			<Card.Root class="border-dashed">
				<Card.Content class="flex flex-col items-center justify-center py-16">
					<IconSpeakerphone class="w-12 h-12 text-muted-foreground mb-4" />
					<Card.Title class="mb-2">No events hosted</Card.Title>
					<Card.Description class="text-center mb-4">
						{#if $user && $user.role === 1}
							You haven't created any events yet. Start hosting your first event!
						{:else}
							Become a host to start creating and managing your own events.
						{/if}
					</Card.Description>
					{#if $user && $user.role === 1}
						<Button href="/events/create">Create Event</Button>
					{:else}
						<Button href="/profile">Become a Host</Button>
					{/if}
				</Card.Content>
			</Card.Root>
		{:else}
			<div class="grid md:grid-cols-2 lg:grid-cols-3 gap-6">
				{#each hosting as event}
					<Card.Root class="hover:shadow-lg transition-shadow duration-200 flex flex-col">
						<Card.Header>
							<div class="flex items-start justify-between gap-2">
								<Card.Title class="text-xl line-clamp-2">
									{event.name}
								</Card.Title>
								<Badge variant="default">Hosting</Badge>
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
		{/if}
	{/if}
</div>
