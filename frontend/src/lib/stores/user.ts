import { writable } from "svelte/store";

export interface User {
    uuid: string;
    name: string;
    email: string;
    role: number;
}

export const user = writable<User | null>(null);