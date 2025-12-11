export const AVAILABLE_TAGS = [
	'Music',
	'Sports',
	'Food & Drink',
	'Networking',
	'Workshop',
	'Party',
	'Conference',
	'Meetup',
	'Art',
	'Charity',
	'Outdoor',
	'Gaming',
	'Tech',
	'Health & Wellness',
	'Education'
] as const;

export type Tag = (typeof AVAILABLE_TAGS)[number];
