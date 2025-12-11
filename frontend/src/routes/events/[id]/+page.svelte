<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import { Button } from '$lib/components/ui/button';
	import * as Card from '$lib/components/ui/card';
	import { Badge } from '$lib/components/ui/badge';
	import * as Alert from '$lib/components/ui/alert';
	import { Separator } from '$lib/components/ui/separator';
	import { MapViewer } from '$lib/components/ui/map-viewer';
	import { user } from '$lib/stores/user';
	import {
		IconCalendarEvent,
		IconMapPin,
		IconCurrencyDollar,
		IconUser,
		IconClock,
		IconArrowLeft,
		IconCheck,
		IconAlertCircle,
		IconUsers,
		IconTag
	} from '@tabler/icons-svelte';

	interface Event {
		event_id: string;
		name: string;
		description: string;
		date: string;
		latitude: number;
		longitude: number;
		fee: number;
		organizer_id: string;
		attendees_count: number;
		user_registered: boolean;
		tag?: string[];
	}

	const api = import.meta.env.VITE_API_URL;
	let event = $state<Event | null>(null);
	let loading = $state(true);
	let error = $state('');
	let isRegistering = $state(false);
	let registrationSuccess = $state(false);
	let registrationError = $state('');
	let isRegistered = $state(false);

	const eventId = $page.params.id;

	onMount(async () => {
		await fetchEventDetails();
	});

	async function fetchEventDetails() {
		try {
			const response = await fetch(`${api}/api/events/${eventId}`, {
				credentials: 'include'
			});
			if (response.ok) {
				event = await response.json();
				if (event == null) {
					error = 'Event not found';
					return;
				}
				isRegistered = event.user_registered;
			} else if (response.status === 401) {
				error = 'Please login to view event details';
			} else {
				error = 'Failed to load event details';
			}
		} catch (err) {
			error = 'An error occurred while loading event details';
		} finally {
			loading = false;
		}
	}

	async function handleRegister() {
		if (!$user) {
			goto('/login');
			return;
		}

		if (isRegistered) {
			registrationError = 'You are already registered for this event';
			return;
		}

		isRegistering = true;
		registrationError = '';
		registrationSuccess = false;

		try {
			const response = await fetch(`${api}/api/events/${eventId}/register`, {
				method: 'POST',
				credentials: 'include'
			});

			if (response.ok) {
				registrationSuccess = true;
				isRegistered = true;
				await fetchEventDetails();
			} else {
				const data = await response.json();
				registrationError = data.error || 'Failed to register for event';
			}
		} catch (err) {
			registrationError = 'An error occurred while registering';
		} finally {
			isRegistering = false;
		}
	}

	async function handleUnregister() {
		if (!$user) {
			return;
		}

		if (!isRegistered) {
			registrationError = 'You are not registered for this event';
			return;
		}

		isRegistering = true;
		registrationError = '';

		try {
			const response = await fetch(`${api}/api/events/${eventId}/unregister`, {
				method: 'DELETE',
				credentials: 'include'
			});

			if (response.ok) {
				isRegistered = false;
				registrationSuccess = false;
				await fetchEventDetails();
			} else {
				const data = await response.json();
				registrationError = data.error || 'Failed to unregister from event';
			}
		} catch (err) {
			registrationError = 'An error occurred while unregistering';
		} finally {
			isRegistering = false;
		}
	}

	function formatDate(dateString: string): string {
		const date = new Date(dateString);
		return date.toLocaleDateString('en-US', {
			weekday: 'long',
			year: 'numeric',
			month: 'long',
			day: 'numeric'
		});
	}

	function formatTime(dateString: string): string {
		const date = new Date(dateString);
		return date.toLocaleTimeString('en-US', {
			hour: '2-digit',
			minute: '2-digit'
		});
	}

	function formatFee(fee: number): string {
		return fee === 0 ? 'Free' : `$${fee.toFixed(2)}`;
	}
</script>

<div class="container mx-auto px-4 py-8 max-w-5xl">
	<div class="mb-6">
		<Button href="/events" variant="ghost" size="sm" class="gap-2">
			<IconArrowLeft class="w-4 h-4" />
			Back to Events
		</Button>
	</div>

	{#if loading}
		<div class="flex flex-col items-center justify-center h-96 space-y-4">
			<div class="animate-spin rounded-full h-16 w-16 border-b-2 border-primary"></div>
			<p class="text-muted-foreground text-lg">Loading event details...</p>
		</div>
	{:else if error}
		<Alert.Root variant="destructive">
			<IconAlertCircle class="h-4 w-4" />
			<Alert.Title>Error</Alert.Title>
			<Alert.Description>{error}</Alert.Description>
		</Alert.Root>
		{#if error == 'Please login to view event details'}
			<div class="mt-6">
				<Button href="/login">Login to Continue</Button>
			</div>
		{/if}
	{:else if event}
		{#if registrationSuccess}
			<Alert.Root class="mb-6 border-green-500 bg-green-50 dark:bg-green-950">
				<IconCheck class="h-4 w-4 text-green-600" />
				<Alert.Title class="text-green-800 dark:text-green-200"
					>Successfully Registered!</Alert.Title
				>
				<Alert.Description class="text-green-700 dark:text-green-300">
					You've been registered for this event. See you there!
				</Alert.Description>
			</Alert.Root>
		{/if}

		{#if registrationError}
			<Alert.Root variant="destructive" class="mb-6">
				<IconAlertCircle class="h-4 w-4" />
				<Alert.Title>Registration Failed</Alert.Title>
				<Alert.Description>{registrationError}</Alert.Description>
			</Alert.Root>
		{/if}

		<div class="mb-8 space-y-4">
			<div class="flex items-start justify-between gap-4 flex-wrap">
				<div class="space-y-2 flex-1">
					<h1 class="text-4xl md:text-5xl font-bold tracking-tight">{event.name}</h1>
					<div class="flex items-center gap-2 flex-wrap">
						<Badge variant="secondary" class="text-sm">
							<IconCalendarEvent class="w-3 h-3 mr-1" />
							Upcoming Event
						</Badge>
						{#if event.fee === 0}
							<Badge
								variant="outline"
								class="text-sm bg-green-50 dark:bg-green-950 border-green-500"
							>
								<IconCheck class="w-3 h-3 mr-1" />
								Free Entry
							</Badge>
						{/if}
					</div>
					{#if event.tag && event.tag.length > 0}
						<div class="flex items-center gap-2 flex-wrap pt-2">
							<IconTag class="w-4 h-4 text-muted-foreground" />
							{#each event.tag as tag}
								<Badge variant="outline" class="text-sm">{tag}</Badge>
							{/each}
						</div>
					{/if}
				</div>
			</div>
		</div>

		<div class="grid lg:grid-cols-3 gap-6">
			<div class="lg:col-span-2 space-y-6">
				<Card.Root>
					<Card.Header>
						<Card.Title class="text-2xl flex items-center gap-2">
							<IconCalendarEvent class="w-6 h-6 text-primary" />
							About This Event
						</Card.Title>
					</Card.Header>
					<Card.Content>
						<p class="text-base text-muted-foreground leading-relaxed whitespace-pre-wrap">
							{event.description}
						</p>
					</Card.Content>
				</Card.Root>

				<Card.Root>
					<Card.Header>
						<Card.Title class="text-2xl flex items-center gap-2">
							<IconMapPin class="w-6 h-6 text-primary" />
							Event Location
						</Card.Title>
					</Card.Header>
					<Card.Content>
						<div class="space-y-4">
							<MapViewer
								latitude={event.latitude}
								longitude={event.longitude}
								height="350px"
								zoom={14}
							/>

							<div class="flex items-center justify-between text-sm text-muted-foreground">
								<div class="flex items-center gap-2">
									<IconMapPin class="w-4 h-4" />
									<span>
										{event.latitude.toFixed(4)}, {event.longitude.toFixed(4)}
									</span>
								</div>
							</div>

							<div class="flex gap-2">
								<Button
									variant="outline"
									class="flex-1"
									onclick={() => {
										if (event) {
											window.open(
												`https://www.google.com/maps?q=${event.latitude},${event.longitude}`,
												'_blank'
											);
										}
									}}
								>
									<IconMapPin class="w-4 h-4 mr-2" />
									Open in Google Maps
								</Button>
							</div>
						</div>
					</Card.Content>
				</Card.Root>
			</div>

			<div class="space-y-6">
				<Card.Root class="sticky top-4">
					<Card.Header>
						<Card.Title class="text-xl">Event Details</Card.Title>
					</Card.Header>
					<Card.Content class="space-y-4">
						<div class="space-y-3">
							<div class="flex items-start gap-3">
								<div class="p-2 rounded-lg bg-primary/10 text-primary flex-shrink-0">
									<IconCalendarEvent class="w-5 h-5" />
								</div>
								<div class="flex-1 min-w-0">
									<p class="text-sm font-medium text-muted-foreground">Date</p>
									<p class="text-base font-semibold">{formatDate(event.date)}</p>
								</div>
							</div>

							<div class="flex items-start gap-3">
								<div class="p-2 rounded-lg bg-primary/10 text-primary flex-shrink-0">
									<IconClock class="w-5 h-5" />
								</div>
								<div class="flex-1 min-w-0">
									<p class="text-sm font-medium text-muted-foreground">Time</p>
									<p class="text-base font-semibold">{formatTime(event.date)}</p>
								</div>
							</div>

							<Separator />

							<div class="flex items-start gap-3">
								<div class="p-2 rounded-lg bg-primary/10 text-primary flex-shrink-0">
									<IconCurrencyDollar class="w-5 h-5" />
								</div>
								<div class="flex-1 min-w-0">
									<p class="text-sm font-medium text-muted-foreground">Entry Fee</p>
									<p class="text-base font-semibold">{formatFee(event.fee)}</p>
								</div>
							</div>

							<Separator />

							<div class="flex items-start gap-3">
								<div class="p-2 rounded-lg bg-primary/10 text-primary flex-shrink-0">
									<IconUsers class="w-5 h-5" />
								</div>
								<div class="flex-1 min-w-0">
									<p class="text-sm font-medium text-muted-foreground">Attendees</p>
									<p class="text-base font-semibold">{event.attendees_count} registered</p>
								</div>
							</div>
						</div>

						<Separator />

						<div class="space-y-3">
							{#if !$user}
								<Button class="w-full gap-2" size="lg" href="/login">
									<IconUser class="w-4 h-4" />
									Login to Register
								</Button>
							{:else if event && event.organizer_id === $user.uuid}
								<div class="text-center p-4 bg-muted rounded-lg">
									<p class="text-sm font-medium">You are the organizer of this event</p>
								</div>
							{:else if isRegistered}
								<div class="space-y-2">
									<Button variant="outline" class="w-full gap-2" size="lg" disabled>
										<IconCheck class="w-4 h-4" />
										Registered
									</Button>
									<Button
										variant="destructive"
										class="w-full"
										disabled={isRegistering}
										onclick={handleUnregister}
									>
										{#if isRegistering}
											<div
												class="animate-spin rounded-full h-4 w-4 border-b-2 border-current"
											></div>
											Unregistering...
										{:else}
											Unregister
										{/if}
									</Button>
								</div>
							{:else}
								<Button
									class="w-full gap-2"
									size="lg"
									disabled={isRegistering}
									onclick={handleRegister}
								>
									{#if isRegistering}
										<div class="animate-spin rounded-full h-4 w-4 border-b-2 border-current"></div>
										Registering...
									{:else}
										<IconCheck class="w-4 h-4" />
										Register for Event
									{/if}
								</Button>
							{/if}
						</div>
					</Card.Content>
				</Card.Root>
			</div>
		</div>
	{/if}
</div>
