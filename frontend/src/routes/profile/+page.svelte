<script lang="ts">
	import { goto } from '$app/navigation';
	import { Button } from '$lib/components/ui/button';
	import * as Card from '$lib/components/ui/card';
	import * as Alert from '$lib/components/ui/alert';
	import * as Avatar from '$lib/components/ui/avatar';
	import { Badge } from '$lib/components/ui/badge';
	import { Separator } from '$lib/components/ui/separator';
	import { user, type User } from '$lib/stores/user';

	import {
		IconMail,
		IconLogout,
		IconSettings,
		IconCalendarEvent,
		IconTicket,
		IconAlertCircle,
		IconSpeakerphone
	} from '@tabler/icons-svelte';

	const api = import.meta.env.VITE_API_URL;

	let userProfile = $state<User | null>(null);
	let loading = $state(true);
	let error = $state<string | null>(null);
	let loggingOut = $state(false);

	const roleNames: Record<number, string> = {
		0: 'Attendee',
		1: 'Host'
	};

	const roleColors: Record<number, 'default' | 'secondary' | 'destructive' | 'outline'> = {
		0: 'default',
		1: 'secondary'
	};

	$effect(() => {
		async function fetchUserProfile() {
			try {
				loading = true;
				error = null;
				const response = await fetch(`${import.meta.env.VITE_API_URL}/api/me`, {
					credentials: 'include'
				});
				if (response.ok) {
					const data: User = await response.json();
					userProfile = data;
				} else if (response.status === 401) {
					error = 'Not authenticated. Please log in.';
					setTimeout(() => goto('/login'), 2000);
				} else {
					error = 'Failed to load user profile';
				}
			} catch (err) {
				error = 'An error occurred while loading user profile';
				console.error(err);
			} finally {
				loading = false;
			}
		}
		fetchUserProfile();
	});

	async function handleLogout(e: Event) {
		e.preventDefault();
		error = '';

		try {
			const response = await fetch(`${api}/api/logout`, {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json'
				},
				credentials: 'include'
			});

			if (response.ok) {
				user.set(null);
				window.location.href = '/';
			} else {
				error = 'Failed to log out.';
			}
		} catch (err) {
			error = 'An error occurred. Please try again.';
		}
	}

	async function handleBecomeHost() {
		try {
			const response = await fetch(`${api}/api/become-host`, {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json'
				},
				credentials: 'include'
			});

			if (response.ok) {
				alert('Your request to become a host has been submitted!');
			} else {
				alert('Failed to submit request. Please try again.');
			}
		} catch (err) {
			alert('An error occurred. Please try again.');
		}
	}

	function getInitials(name: string): string {
		return name
			.split(' ')
			.map((n) => n[0])
			.join('')
			.toUpperCase()
			.slice(0, 2);
	}
</script>

<div class="container mx-auto px-4 py-8 md:py-12 max-w-4xl">
	<div class="mb-8">
		<h1 class="text-3xl md:text-4xl font-bold mb-2">User Profile</h1>
		<p class="text-muted-foreground">Manage your account and view your information</p>
	</div>

	{#if error}
		<Alert.Root variant="destructive" class="mb-6">
			<IconAlertCircle class="size-4" />
			<Alert.Title>Error</Alert.Title>
			<Alert.Description>{error}</Alert.Description>
		</Alert.Root>
	{/if}

	{#if loading}
		<div class="flex items-center justify-center py-12">
			<div class="text-center">
				<div
					class="inline-block h-8 w-8 animate-spin rounded-full border-4 border-solid border-current border-r-transparent motion-reduce:animate-[spin_1.5s_linear_infinite]"
					role="status"
				>
					<span class="sr-only">Loading...</span>
				</div>
				<p class="mt-4 text-muted-foreground">Loading profile...</p>
			</div>
		</div>
	{:else if userProfile}
		<div class="grid gap-6 md:grid-cols-3">
			<Card.Root class="md:col-span-2">
				<Card.Header>
					<div class="flex items-start justify-between">
						<div class="flex items-center gap-4">
							<Avatar.Root class="size-16 md:size-20">
								<Avatar.Fallback class="text-lg md:text-xl font-semibold">
									{getInitials(userProfile.name)}
								</Avatar.Fallback>
							</Avatar.Root>
							<div>
								<Card.Title class="text-2xl">{userProfile.name}</Card.Title>
								<Card.Description class="flex items-center gap-2 mt-1">
									<IconMail class="size-4" />
									{userProfile.email}
								</Card.Description>
							</div>
						</div>
						<Badge variant={roleColors[userProfile.role]}>
							{#if userProfile.role === 0}
								<IconTicket class="size-3 mr-1" />
								<span>Attendee</span>
							{:else if userProfile.role === 1}
								<IconSpeakerphone class="size-3 mr-1" />
								<span>Host</span>
							{/if}
						</Badge>
					</div>
				</Card.Header>
				<Card.Content>
					<Separator class="mb-6" />
					<div class="space-y-4">
						<div class="flex items-center gap-3 p-3 rounded-lg bg-muted/50">
							<div class="p-2 rounded-md bg-primary/10">
								<IconMail class="size-5 text-primary" />
							</div>
							<div>
								<p class="text-sm text-muted-foreground">Email Address</p>
								<p class="font-medium break-all">{userProfile.email}</p>
							</div>
						</div>
						<div class="flex items-center gap-3 p-3 rounded-lg bg-muted/50">
							<div class="p-2 rounded-md bg-primary/10">
								<IconTicket class="size-5 text-primary" />
							</div>
							<div>
								<p class="text-sm text-muted-foreground">Account Role</p>
								<p class="font-medium">{roleNames[userProfile.role] || 'Unknown'}</p>
							</div>
						</div>
					</div>
				</Card.Content>
			</Card.Root>

			<Card.Root>
				<Card.Header>
					<Card.Title>Quick Actions</Card.Title>
					<Card.Description>Manage your account and settings</Card.Description>
				</Card.Header>
				<Card.Content>
					<div class="space-y-3">
						<Button
							variant="outline"
							class="w-full justify-start"
							onclick={() => goto('/events/my')}
						>
							<IconCalendarEvent class="size-4 mr-2" />
							My Events
						</Button>
						<Button
							variant="outline"
							class="w-full justify-start"
							onclick={() => alert('Settings coming soon!')}
						>
							<IconSettings class="size-4 mr-2" />
							Account Settings
						</Button>
						{#if userProfile.role === 0}
							<Button variant="outline" class="w-full justify-start" onclick={handleBecomeHost}>
								<IconSettings class="size-4 mr-2" />
								Become a Host
							</Button>
						{/if}
						<Separator />
						<Button
							variant="destructive"
							class="w-full justify-start"
							onclick={handleLogout}
							disabled={loggingOut}
						>
							<IconLogout class="size-4 mr-2" />
							{loggingOut ? 'Logging out...' : 'Logout'}
						</Button>
					</div>
				</Card.Content>
			</Card.Root>
		</div>
	{/if}
</div>

<style>
	@keyframes spin {
		to {
			transform: rotate(360deg);
		}
	}
</style>
