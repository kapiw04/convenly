import { writable } from "svelte/store";
import { browser } from "$app/environment";

const LOCATION_STORAGE_KEY = "user_location";

export interface UserLocation {
    latitude: number;
    longitude: number;
}

const DEFAULT_LOCATION: UserLocation = {
    latitude: 51.505,
    longitude: -0.09
};

function getStoredLocation(): UserLocation {
    if (!browser) return DEFAULT_LOCATION;
    
    const stored = localStorage.getItem(LOCATION_STORAGE_KEY);
    if (!stored) return DEFAULT_LOCATION;
    
    try {
        const parsed = JSON.parse(stored);
        if (typeof parsed.latitude === 'number' && typeof parsed.longitude === 'number') {
            return parsed;
        }
    } catch {
        // Invalid JSON, return default
    }
    return DEFAULT_LOCATION;
}

function createLocationStore() {
    const initialValue = getStoredLocation();
    
    const { subscribe, set, update } = writable<UserLocation>(initialValue);

    return {
        subscribe,
        set: (value: UserLocation) => {
            if (browser) {
                localStorage.setItem(LOCATION_STORAGE_KEY, JSON.stringify(value));
            }
            set(value);
        },
        setCoords: (latitude: number, longitude: number) => {
            const value = { latitude, longitude };
            if (browser) {
                localStorage.setItem(LOCATION_STORAGE_KEY, JSON.stringify(value));
            }
            set(value);
        },
        clear: () => {
            if (browser) {
                localStorage.removeItem(LOCATION_STORAGE_KEY);
            }
            set(DEFAULT_LOCATION);
        }
    };
}

export const userLocation = createLocationStore();
