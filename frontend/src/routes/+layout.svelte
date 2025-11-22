<script lang="ts">
	import '../app.css';
	import favicon from '$lib/assets/favicon.svg';
	import { Button } from "$lib/components/ui/button";
	import { Separator } from "$lib/components/ui/separator";
	import { ModeWatcher } from 'mode-watcher';
	import ThemeToggle from '$lib/components/ui/theme-toggle/theme-toggle.svelte';
	
	let { children } = $props();
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
	<title>Convenly - Event Management Platform</title>
</svelte:head>

<div class="min-h-screen bg-background flex flex-col">
	<header class="sticky top-0 z-50 w-full border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
		<div class="container mx-auto flex h-16 items-center justify-between">
			<div class="flex items-center gap-8">
				<a href="/" class="flex items-center space-x-2">
					<span class="text-2xl font-bold bg-gradient-to-r from-primary to-primary/60 bg-clip-text text-transparent">
						Convenly
					</span>
				</a>
				<nav class="hidden md:flex items-center space-x-6 text-sm font-medium">
					<a 
						href="/events" 
						class="transition-colors hover:text-foreground/80 text-foreground/60"
					>
						Events
					</a>
					<a 
						href="/" 
						class="transition-colors hover:text-foreground/80 text-foreground/60"
					>
						About
					</a>
				</nav>
			</div>
			<div class="flex items-center gap-3">
				<ThemeToggle />
				<Button href="/login" variant="ghost">
					Log in
				</Button>
				<Button href="/register">
					Sign up
				</Button>
			</div>
		</div>
	</header>

	<main class="flex-1">
		<ModeWatcher />
		{@render children()}
	</main>

	<footer class="border-t bg-background">
		<div class="container mx-auto py-8">
			<div class="grid grid-cols-1 md:grid-cols-3 gap-8">
				<div>
					<h3 class="font-bold text-lg mb-3">Convenly</h3>
				</div>
				<div>
					<h4 class="font-semibold mb-3">Quick Links</h4>
					<ul class="space-y-2 text-sm text-muted-foreground">
						<li><a href="/events" class="hover:text-foreground transition-colors">Browse Events</a></li>
						<li><a href="/register" class="hover:text-foreground transition-colors">Sign Up</a></li>
						<li><a href="/login" class="hover:text-foreground transition-colors">Log In</a></li>
					</ul>
				</div>
			</div>
			<Separator class="my-6" />
			<div class="text-center text-sm text-muted-foreground">
				© {new Date().getFullYear()} Convenly. All rights reserved.
			</div>
		</div>
	</footer>
</div>
