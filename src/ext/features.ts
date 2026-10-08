import { ref, computed } from 'vue';

import { getCurrentUserInfo } from '@/lib/userstate.ts';

import { useBusiness } from './business.ts';

// Inventory, sales and customers are opt-in: a person turns them on in Settings. The choice is kept in this browser
// per user name. People who were invited to work in somebody else's business get the features automatically,
// because being invited is the opt-in.
const STORAGE_KEY = 'ebk_ext_features';

function currentUsername(): string {
    return getCurrentUserInfo()?.username ?? '';
}

function readEnabled(): boolean {
    try {
        const stored = JSON.parse(localStorage.getItem(STORAGE_KEY) ?? '{}') as Record<string, boolean>;
        const username = currentUsername();

        return !!username && stored[username] === true;
    } catch {
        return false;
    }
}

function writeEnabled(enabled: boolean): void {
    try {
        const stored = JSON.parse(localStorage.getItem(STORAGE_KEY) ?? '{}') as Record<string, boolean>;
        stored[currentUsername()] = enabled;
        localStorage.setItem(STORAGE_KEY, JSON.stringify(stored));
    } catch {
        // storage can be unavailable; the switch then only lasts until the page is reloaded
    }
}

const enabled = ref<boolean>(readEnabled());

export function useBusinessFeatures() {
    enabled.value = readEnabled(); // a different person may have logged in since this module was loaded

    const { working, invitations } = useBusiness();

    const worksForSomeoneElse = computed<boolean>(() => working.value.some(b => b.status === 'active'));

    /** Inventory, sales and customers are shown. */
    const available = computed<boolean>(() => enabled.value || worksForSomeoneElse.value);

    /** The Team entry is shown: with the features, or when there is an invitation to answer. */
    const teamAvailable = computed<boolean>(() => available.value || invitations.value.length > 0);

    function setEnabled(value: boolean): void {
        writeEnabled(value);
        enabled.value = value;
    }

    return { enabled, available, teamAvailable, worksForSomeoneElse, setEnabled };
}

/** For route guards: waits for the business list, then says whether the features may be used. */
export async function canUseBusinessFeatures(allowInvitations: boolean = false): Promise<boolean> {
    const { ensureLoaded } = useBusiness();

    try {
        await ensureLoaded();
    } catch {
        // without the list only the explicit opt-in counts
    }

    const { available, teamAvailable } = useBusinessFeatures();

    return allowInvitations ? teamAvailable.value : available.value;
}
