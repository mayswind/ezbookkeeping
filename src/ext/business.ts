import { ref, computed } from 'vue';
import axios from 'axios';

import { getCurrentUserInfo } from '@/lib/userstate.ts';

import api from './api.ts';
import type { BusinessInfo, BusinessRole } from './types.ts';

/** Request header that tells the server whose business a request is about (see pkg/ext/middleware). */
export const BUSINESS_HEADER = 'X-Business-Id';

const STORAGE_KEY = 'ebk_ext_business';

function currentUsername(): string {
    return getCurrentUserInfo()?.username ?? '';
}

// The choice is remembered per user name, so a different person logging in on the same browser never inherits it.
function readSelection(): string {
    try {
        const raw = localStorage.getItem(STORAGE_KEY);

        if (!raw) {
            return '';
        }

        const stored = JSON.parse(raw) as { user?: string, business?: string };
        return stored.user && stored.user === currentUsername() ? (stored.business ?? '') : '';
    } catch {
        return '';
    }
}

function writeSelection(ownerUid: string): void {
    try {
        if (ownerUid) {
            localStorage.setItem(STORAGE_KEY, JSON.stringify({ user: currentUsername(), business: ownerUid }));
        } else {
            localStorage.removeItem(STORAGE_KEY);
        }
    } catch {
        // storage can be unavailable (private windows); the selection then lasts until reload only
    }
}

/** The business the person is working in, '' meaning their own. */
const selectedBusinessId = ref<string>(readSelection());
const businesses = ref<BusinessInfo[]>([]);
const loading = ref<boolean>(false);
const loaded = ref<boolean>(false);

let headerInstalled = false;

/** Adds the business header to every request of the app once. Safe to call many times. */
export function installBusinessHeader(): void {
    if (headerInstalled) {
        return;
    }

    headerInstalled = true;

    axios.interceptors.request.use(config => {
        const id = readSelection();

        if (id) {
            config.headers[BUSINESS_HEADER] = id;
        }

        return config;
    });
}

/** Selects a business and reloads, because every screen holds data of the previous one. */
export function switchBusiness(ownerUid: string): void {
    const own = businesses.value.find(b => b.status === 'owner');
    writeSelection(own && own.ownerUid === ownerUid ? '' : ownerUid);
    window.location.reload();
}

export function useBusiness() {
    const own = computed<BusinessInfo | undefined>(() => businesses.value.find(b => b.status === 'owner'));
    const working = computed<BusinessInfo[]>(() => businesses.value.filter(b => b.status !== 'pending'));
    const invitations = computed<BusinessInfo[]>(() => businesses.value.filter(b => b.status === 'pending'));

    const current = computed<BusinessInfo | undefined>(() => {
        if (!selectedBusinessId.value) {
            return own.value;
        }

        return businesses.value.find(b => b.ownerUid === selectedBusinessId.value && b.status === 'active');
    });

    const role = computed<BusinessRole>(() => current.value?.role ?? 'owner');
    const isOwner = computed<boolean>(() => role.value === 'owner');
    const canManage = computed<boolean>(() => role.value === 'owner' || role.value === 'manager');
    const workingForSomeoneElse = computed<boolean>(() => !!selectedBusinessId.value && role.value !== 'owner');

    async function refresh(): Promise<void> {
        if (loading.value) {
            return;
        }

        loading.value = true;

        try {
            businesses.value = await api.getMyBusinesses();
            loaded.value = true;

            // the saved business may be gone (the owner removed this person): fall back to their own
            const stored = readSelection();

            if (stored && !businesses.value.some(b => b.ownerUid === stored && b.status === 'active')) {
                writeSelection('');
                window.location.reload();
            }
        } finally {
            loading.value = false;
        }
    }

    return { businesses, own, working, invitations, current, role, isOwner, canManage, workingForSomeoneElse, loading, loaded, selectedBusinessId, refresh, switchBusiness };
}
