import { ref, computed } from 'vue';

import { getCurrentUserInfo } from '@/lib/userstate.ts';

import api from './api.ts';
import { useBusiness } from './business.ts';

// Inventory, sales and customers are opt-in: a person turns them on in Settings. The choice is kept on the server, so
// it follows the person to every device. People invited to work in somebody else's business get the features
// automatically, because being invited is the opt-in.
//
// The browser keeps a copy of the last answer only so the menus do not flicker while the page loads;
// the server is the source of truth.
const CACHE_KEY = 'ebk_ext_features';

function currentUsername(): string {
    return getCurrentUserInfo()?.username ?? '';
}

function readCache(): boolean {
    try {
        const stored = JSON.parse(localStorage.getItem(CACHE_KEY) ?? '{}') as Record<string, boolean>;
        const username = currentUsername();

        return !!username && stored[username] === true;
    } catch {
        return false;
    }
}

function writeCache(enabled: boolean): void {
    try {
        const stored = JSON.parse(localStorage.getItem(CACHE_KEY) ?? '{}') as Record<string, boolean>;
        stored[currentUsername()] = enabled;
        localStorage.setItem(CACHE_KEY, JSON.stringify(stored));
    } catch {
        // the copy is only a convenience
    }
}

const enabled = ref<boolean>(readCache());
const loadedFor = ref<string>(''); // the user name the server answer belongs to
let loadingPromise: Promise<void> | null = null;

/** Asks the server once per person whether the features are on. */
export function loadFeatureSetting(): Promise<void> {
    const user = currentUsername();

    if (loadedFor.value === user) {
        return Promise.resolve();
    }

    if (!loadingPromise) {
        loadingPromise = (async () => {
            try {
                let settings = await api.getMySettings();

                // a choice made before settings were kept on the server (stored in this browser) is carried over once
                if (!settings.configured && readCache()) {
                    settings = await api.updateMySettings(true);
                }

                enabled.value = settings.businessFeatures;
                writeCache(settings.businessFeatures);
                loadedFor.value = user;
            } finally {
                loadingPromise = null;
            }
        })();
    }

    return loadingPromise;
}

export function useBusinessFeatures() {
    if (loadedFor.value !== currentUsername()) {
        enabled.value = readCache(); // a different person may have logged in since this module was loaded
    }

    const { working, invitations } = useBusiness();

    const worksForSomeoneElse = computed<boolean>(() => working.value.some(b => b.status === 'active'));

    /** Inventory, sales and customers are shown. */
    const available = computed<boolean>(() => enabled.value || worksForSomeoneElse.value);

    /** The Team entry is shown: with the features, or when there is an invitation to answer. */
    const teamAvailable = computed<boolean>(() => available.value || invitations.value.length > 0);

    /** Saves the choice on the server. Throws, and leaves the switch where it was, if that fails. */
    async function setEnabled(value: boolean): Promise<void> {
        const previous = enabled.value;
        enabled.value = value;

        try {
            const saved = await api.updateMySettings(value);
            enabled.value = saved.businessFeatures;
            writeCache(saved.businessFeatures);
            loadedFor.value = currentUsername();
        } catch (error) {
            enabled.value = previous;
            throw error;
        }
    }

    return { enabled, available, teamAvailable, worksForSomeoneElse, loadFeatureSetting, setEnabled };
}

/** For route guards: waits for the business list and the setting, then says whether the features may be used. */
export async function canUseBusinessFeatures(allowInvitations: boolean = false): Promise<boolean> {
    const { ensureLoaded } = useBusiness();

    await Promise.all([ensureLoaded(), loadFeatureSetting()].map(promise => promise.catch(() => {
        // if either answer is unavailable only what is already known counts
    })));

    const { available, teamAvailable } = useBusinessFeatures();

    return allowInvitations ? teamAvailable.value : available.value;
}
