<script lang="ts">
	import { Button } from "$lib/components/ui/button";
	import { Badge } from "$lib/components/ui/badge";
</script>

<div class="container mx-auto px-4 py-12">
	<div class="text-center mb-20 space-y-6">
		<Badge variant="secondary" class="mb-4">Event Management Platform</Badge>
		<h1 class="text-6xl font-bold tracking-tight">
			Welcome to <span class="text-primary">Convenly</span>
		</h1>
		<div class="flex justify-center gap-4 pt-6">
			<Button href="/events" size="lg" class="text-lg px-8">
				Browse Events
			</Button>
			<Button href="/register" variant="outline" size="lg" class="text-lg px-8">
				Get Started
			</Button>
		</div>
	</div>
</div>
