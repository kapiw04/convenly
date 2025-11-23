<script lang="ts">
	import { goto } from '$app/navigation';
	import { Button } from '$lib/components/ui/button';
	import * as Card from '$lib/components/ui/card';
	import { Input } from '$lib/components/ui/input';
	import { Textarea } from '$lib/components/ui/textarea';
	import { Label } from '$lib/components/ui/label';
	import { MapPicker } from '$lib/components/ui/map-picker';
	import * as Alert from '$lib/components/ui/alert';
	import { Separator } from '$lib/components/ui/separator';
	import Calendar from '$lib/components/ui/calendar/calendar.svelte';
	import { CalendarDate, type DateValue } from '@internationalized/date';
	import {
		IconCalendarEvent,
		IconMapPin,
		IconFileDescription,
		IconCurrencyDollar,
		IconClock,
		IconCheck,
		IconAlertCircle,
		IconCalendarPlus
	} from '@tabler/icons-svelte';

	const api = import.meta.env.VITE_API_URL;

	let eventName = $state('');
	let eventDescription = $state('');
	let eventFee = $state('');
	let selectedLatitude = $state(51.505);
	let selectedLongitude = $state(-0.09);
	let isSubmitting = $state(false);
	let successMessage = $state('');
	let errorMessage = $state('');

	let calendarValue = $state<DateValue | undefined>(undefined);
	let startTime = $state('10:00:00');
	let endTime = $state('12:00:00');

	function handleLocationChange(lat: number, lng: number) {
		selectedLatitude = lat;
		selectedLongitude = lng;
	}

	async function handleSubmit() {
		successMessage = '';
		errorMessage = '';

		if (!eventName.trim()) {
			errorMessage = 'Please enter an event name';
			return;
		}
		if (!eventDescription.trim()) {
			errorMessage = 'Please enter an event description';
			return;
		}
		if (!calendarValue) {
			errorMessage = 'Please select an event date';
			return;
		}
		if (!startTime) {
			errorMessage = 'Please select a start time';
			return;
		}
		if (eventFee === '' || isNaN(parseFloat(eventFee)) || parseFloat(eventFee) < 0) {
			errorMessage = 'Please enter a valid event fee (0 or greater)';
			return;
		}

		const eventDateTime = `${calendarValue.year}-${String(calendarValue.month).padStart(2, '0')}-${String(calendarValue.day).padStart(2, '0')}T${startTime}Z`;

		isSubmitting = true;

		try {
			const response = await fetch(`${api}/api/events/add`, {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json'
				},
				credentials: 'include',
				body: JSON.stringify({
					name: eventName,
					description: eventDescription,
					date: eventDateTime,
					latitude: selectedLatitude,
					longitude: selectedLongitude,
					fee: parseFloat(eventFee)
				})
			});

			if (response.ok) {
				successMessage = 'Event created successfully! Redirecting...';
				eventName = '';
				eventDescription = '';
				calendarValue = undefined;
				startTime = '10:00:00';
				endTime = '12:00:00';
				eventFee = '';
				selectedLatitude = 51.505;
				selectedLongitude = -0.09;

				setTimeout(() => {
					goto('/events');
				}, 1500);
			} else {
				const error = await response.json();
				errorMessage = error.message || 'Failed to create event. Please try again.';
			}
		} catch (err) {
			errorMessage = 'An error occurred. Please check your connection and try again.';
		} finally {
			isSubmitting = false;
		}
	}
</script>

<div class="container mx-auto px-2 py-8 max-w-6xl">
	<div class="mb-8 text-center">
		<div class="inline-flex items-center justify-center w-16 h-16 rounded-full bg-primary/10 mb-4">
			<IconCalendarPlus class="w-8 h-8 text-primary" />
		</div>
		<h1
			class="text-4xl font-bold mb-2 bg-gradient-to-r from-primary to-primary/60 bg-clip-text text-transparent"
		>
			Create New Event
		</h1>
		<p class="text-muted-foreground text-lg">Bring people together by creating an amazing event</p>
	</div>

	{#if successMessage}
		<Alert.Root class="mb-6 border-green-500 bg-green-50 dark:bg-green-950/20">
			<IconCheck class="h-5 w-5 text-green-600 dark:text-green-400" />
			<Alert.Title class="text-green-800 dark:text-green-300">Success</Alert.Title>
			<Alert.Description class="text-green-700 dark:text-green-400">
				{successMessage}
			</Alert.Description>
		</Alert.Root>
	{/if}

	{#if errorMessage}
		<Alert.Root class="mb-6 border-destructive bg-destructive/10">
			<IconAlertCircle class="h-5 w-5 text-destructive" />
			<Alert.Title class="text-destructive">Error</Alert.Title>
			<Alert.Description class="text-destructive/90">
				{errorMessage}
			</Alert.Description>
		</Alert.Root>
	{/if}

	<Card.Root class="shadow-lg">
		<Card.Header class="space-y-1 pb-6">
			<Card.Title class="text-2xl font-semibold flex items-center gap-2">
				<IconCalendarEvent class="w-6 h-6 text-primary" />
				Event Details
			</Card.Title>
			<Card.Description>Fill in the information below to create your event</Card.Description>
		</Card.Header>
		<Separator class="mb-6" />
		<Card.Content class="space-y-6">
			<div class="space-y-2">
				<Label for="event-name" class="text-base font-semibold flex items-center gap-2">
					<IconCalendarEvent class="w-4 h-4 text-primary" />
					Event Name
					<span class="text-destructive">*</span>
				</Label>
				<Input
					id="event-name"
					type="text"
					placeholder="Enter a catchy event name"
					class="w-full text-base"
					bind:value={eventName}
					disabled={isSubmitting}
				/>
			</div>

			<div class="space-y-2">
				<Label for="event-description" class="text-base font-semibold flex items-center gap-2">
					<IconFileDescription class="w-4 h-4 text-primary" />
					Event Description
					<span class="text-destructive">*</span>
				</Label>
				<Textarea
					id="event-description"
					placeholder="Describe what makes your event special..."
					class="w-full min-h-[120px] text-base resize-none"
					bind:value={eventDescription}
					disabled={isSubmitting}
				/>
				<p class="text-sm text-muted-foreground">
					Provide details about the event, activities, and what attendees can expect.
				</p>
			</div>

			<Separator />

			<div class="grid grid-cols-1 lg:grid-cols-5 gap-6">
				<div class="space-y-3 lg:col-span-2">
					<Label class="text-base font-semibold flex items-center gap-2">
						<IconClock class="w-4 h-4 text-primary" />
						Event Date & Time
						<span class="text-destructive">*</span>
					</Label>
					<p class="text-sm text-muted-foreground">Select the date and time range for your event</p>
					<Card.Root class="w-full py-4">
						<Card.Content class="px-4 flex justify-center">
							<Calendar
								type="single"
								bind:value={calendarValue}
								class="bg-transparent p-0"
								disabled={isSubmitting}
							/>
						</Card.Content>
						<Card.Footer class="*:[div]:w-full flex gap-2 border-t px-4 !pt-4">
							<div>
								<Label for="time-from" class="sr-only">Start Time</Label>
								<Input
									id="time-from"
									type="time"
									step="1"
									bind:value={startTime}
									disabled={isSubmitting}
									class="appearance-none [&::-webkit-calendar-picker-indicator]:hidden [&::-webkit-calendar-picker-indicator]:appearance-none"
								/>
							</div>
							<span class="self-center">-</span>
							<div>
								<Label for="time-to" class="sr-only">End Time</Label>
								<Input
									id="time-to"
									type="time"
									step="1"
									bind:value={endTime}
									disabled={isSubmitting}
									class="appearance-none [&::-webkit-calendar-picker-indicator]:hidden [&::-webkit-calendar-picker-indicator]:appearance-none"
								/>
							</div>
						</Card.Footer>
					</Card.Root>
				</div>

				<div class="space-y-3 lg:col-span-3">
					<Label class="text-base font-semibold flex items-center gap-2">
						<IconMapPin class="w-4 h-4 text-primary" />
						Event Location
						<span class="text-destructive">*</span>
					</Label>
					<MapPicker
						latitude={selectedLatitude}
						longitude={selectedLongitude}
						onLocationChange={handleLocationChange}
						height="400px"
					/>
				</div>
			</div>

			<Separator />

			<div class="space-y-2">
				<Label for="event-fee" class="text-base font-semibold flex items-center gap-2">
					<IconCurrencyDollar class="w-4 h-4 text-primary" />
					Event Fee
					<span class="text-destructive">*</span>
				</Label>
				<Input
					id="event-fee"
					type="number"
					placeholder="0.00"
					class="w-full text-base"
					bind:value={eventFee}
					disabled={isSubmitting}
					min="0"
					step="0.01"
				/>
				<p class="text-sm text-muted-foreground">Set to 0 for free events</p>
			</div>

			<Separator />

			<div class="flex gap-3 pt-4">
				<Button
					onclick={handleSubmit}
					class="flex-1 h-11 text-base font-semibold"
					disabled={isSubmitting}
				>
					{#if isSubmitting}
						<span class="flex items-center gap-2">
							<svg
								class="animate-spin h-5 w-5"
								xmlns="http://www.w3.org/2000/svg"
								fill="none"
								viewBox="0 0 24 24"
							>
								<circle
									class="opacity-25"
									cx="12"
									cy="12"
									r="10"
									stroke="currentColor"
									stroke-width="4"
								></circle>
								<path
									class="opacity-75"
									fill="currentColor"
									d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
								></path>
							</svg>
							Creating Event...
						</span>
					{:else}
						<span class="flex items-center gap-2">
							<IconCalendarPlus class="w-5 h-5" />
							Create Event
						</span>
					{/if}
				</Button>
				<Button
					variant="outline"
					onclick={() => goto('/events')}
					class="h-11 px-8"
					disabled={isSubmitting}
				>
					Cancel
				</Button>
			</div>
		</Card.Content>
	</Card.Root>

	<div class="mt-6 text-center">
		<p class="text-sm text-muted-foreground">
			<span class="text-destructive">*</span> Required fields
		</p>
	</div>
</div>
