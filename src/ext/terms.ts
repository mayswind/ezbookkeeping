import { ref, computed } from 'vue';

import { getCurrentUserInfo, isUserLogined } from '@/lib/userstate.ts';

import api from './api.ts';
import { LEGAL_VERSION } from './legalVersion.ts';

// Everyone must accept the current Terms of Service and Privacy Policy. The acceptance is stored on the server (who,
// which version, when), and the app asks again whenever the published version changes (legal/details.json).
const acceptedVersion = ref<string | null>(null); // null: not known yet
const loadedFor = ref<string>('');
let loadingPromise: Promise<void> | null = null;

function currentUsername(): string {
    return getCurrentUserInfo()?.username ?? '';
}

export function useTerms() {
    const user = currentUsername();

    if (loadedFor.value !== user) {
        acceptedVersion.value = null; // a different person may have logged in since
    }

    /** True once the server has answered and the person has not accepted the current version. */
    const needsAcceptance = computed<boolean>(() => acceptedVersion.value !== null && acceptedVersion.value !== LEGAL_VERSION);

    function load(): Promise<void> {
        if (!isUserLogined() || loadedFor.value === user) {
            return Promise.resolve();
        }

        if (!loadingPromise) {
            loadingPromise = api.getMySettings().then(settings => {
                acceptedVersion.value = settings.acceptedTermsVersion;
                loadedFor.value = user;
            }).finally(() => {
                loadingPromise = null;
            });
        }

        return loadingPromise;
    }

    async function accept(): Promise<void> {
        await api.acceptTerms(LEGAL_VERSION);
        acceptedVersion.value = LEGAL_VERSION;
        loadedFor.value = user;
    }

    return { needsAcceptance, load, accept };
}
