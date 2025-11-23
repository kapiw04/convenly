<script lang="ts">
	import '../app.css';
	import favicon from '$lib/assets/favicon.svg';
	import { Button } from '$lib/components/ui/button';
	import { Separator } from '$lib/components/ui/separator';
	import * as Avatar from '$lib/components/ui/avatar';
	import { ModeWatcher } from 'mode-watcher';
	import ThemeToggle from '$lib/components/ui/theme-toggle/theme-toggle.svelte';
	import { user } from '$lib/stores/user';

	const api = import.meta.env.VITE_API_URL;
	let error = $state('');

	let { children } = $props();

	$effect(() => {
		async function checkSession() {
			try {
				const response = await fetch(`${api}/api/me`, {
					credentials: 'include'
				});
				if (response.ok) {
					const data = await response.json();
					$user = data;
					console.log('User is logged in:', data);
				} else {
					$user = null;
				}
			} catch (err) {
				$user = null;
			}
		}

		checkSession();
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
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
	<title>Convenly - Event Management Platform</title>
</svelte:head>

<div class="min-h-screen bg-background flex flex-col">
	<header
		class="sticky top-0 z-50 w-full border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60"
	>
		<div class="container mx-auto flex h-16 items-center justify-between">
			<div class="flex items-center gap-8">
				<a href="/" class="flex items-center space-x-2">
					<span
						class="text-2xl font-bold bg-gradient-to-r from-primary to-primary/60 bg-clip-text text-transparent"
					>
						Convenly
					</span>
				</a>
				<nav class="hidden md:flex items-center space-x-6 text-sm font-medium">
					<a href="/events" class="transition-colors hover:text-foreground/80 text-foreground/60">
						Events
					</a>
				</nav>
			</div>
			<div class="flex items-center gap-3">
				<ThemeToggle />
				{#if $user === null}
					<Button href="/login" variant="ghost">Log in</Button>
					<Button href="/register">Sign up</Button>
				{:else}
					<Button onclick={handleLogout} variant="ghost">Log out</Button>
					<Button href="/profile" variant="ghost" class="p-0">
						<Avatar.Root class="size-8 rounded-lg">
							<Avatar.Fallback class="rounded-lg"
								>{$user.name.charAt(0).toUpperCase()}</Avatar.Fallback
							>
						</Avatar.Root>
					</Button>
				{/if}
			</div>
		</div>
	</header>

	<main class="flex-1">
		<ModeWatcher />
		{@render children()}
	</main>
</div>
