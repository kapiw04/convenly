<script lang="ts">
	import { Button } from "$lib/components/ui/button";
	import * as Card from "$lib/components/ui/card";
	import { Input } from "$lib/components/ui/input";
	import { Label } from "$lib/components/ui/label";
	import * as Alert from "$lib/components/ui/alert";

	const api = import.meta.env.VITE_API_URL;
	if (!api) {
		throw new Error('VITE_API_URL is not defined in environment variables');
	}
	let email = $state('');
	let password = $state('');
	let confirmPassword = $state('');
	let username = $state('');
	let error = $state('');
	let loading = $state(false);

	async function handleRegister(e: Event) {
		e.preventDefault();
		loading = true;
		error = '';

		if (password !== confirmPassword) {
			error = 'Passwords do not match';
			loading = false;
			return;
		}

		try {
			const data = JSON.stringify({ email, password, name: username });
			console.log(data);
			const response = await fetch(`${api}/api/register`, {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json',
				},
				body: data,
			});

			if (response.ok) {
				window.location.href = '/login';
			} else {
				const data = await response.json();
				error = data.message || 'Registration failed. Please try again.';
			}
		} catch (err) {
			error = 'An error occurred. Please try again.';
		} finally {
			loading = false;
		}
	}
</script>

<div class="container mx-auto flex items-center justify-center min-h-[80vh] px-4">
	<Card.Root class="w-full max-w-md">
		<Card.Header class="space-y-1">
			<Card.Title class="text-3xl font-bold text-center">Create account</Card.Title>
			<Card.Description class="text-center">
				Enter your information to get started
			</Card.Description>
		</Card.Header>
		<Card.Content>
			<form onsubmit={handleRegister} class="space-y-4">
				{#if error}
					<Alert.Root variant="destructive">
						<Alert.Description>{error}</Alert.Description>
					</Alert.Root>
				{/if}

				<div class="space-y-2">
					<Label for="username">Username</Label>
					<Input
						id="username"
						type="text"
						placeholder="johndoe"
						bind:value={username}
						required
					/>
				</div>

				<div class="space-y-2">
					<Label for="email">Email address</Label>
					<Input
						id="email"
						type="email"
						placeholder="name@example.com"
						bind:value={email}
						required
					/>
				</div>

				<div class="space-y-2">
					<Label for="password">Password</Label>
					<Input
						id="password"
						type="password"
						placeholder="Create a strong password"
						bind:value={password}
						required
					/>
				</div>

				<div class="space-y-2">
					<Label for="confirmPassword">Confirm Password</Label>
					<Input
						id="confirmPassword"
						type="password"
						placeholder="Confirm your password"
						bind:value={confirmPassword}
						required
					/>
				</div>

				<Button type="submit" class="w-full" disabled={loading}>
					{loading ? 'Creating account...' : 'Sign up'}
				</Button>
			</form>
		</Card.Content>
		<Card.Footer class="flex flex-col space-y-4">
			<div class="text-sm text-center text-muted-foreground">
				Already have an account?
				<a href="/login" class="text-primary hover:underline font-medium">
					Sign in
				</a>
			</div>
		</Card.Footer>
	</Card.Root>
</div>
