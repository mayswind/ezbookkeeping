import { ref, computed } from 'vue';

import { useSettingsStore } from '@/stores/setting.ts';

import { isWebAuthnCompletelySupported } from '@/lib/webauthn.ts';

export function useAppLockPageBase() {
    const settingsStore = useSettingsStore();

    const isSupportedWebAuthn = ref<boolean>(false);

    const isEnableApplicationLock = computed<boolean>({
        get: () => settingsStore.appSettings.applicationLock,
        set: (value) => settingsStore.setEnableApplicationLock(value)
    });

    const isEnableApplicationLockWebAuthn = computed<boolean>({
        get: () => settingsStore.appSettings.applicationLockWebAuthn,
        set: (value) => settingsStore.setEnableApplicationLockWebAuthn(value)
    });

    isWebAuthnCompletelySupported().then(result => {
        isSupportedWebAuthn.value = result;
    });

    return {
        // states
        isSupportedWebAuthn,
        // computed states
        isEnableApplicationLock,
        isEnableApplicationLockWebAuthn
    };
}
