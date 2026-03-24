import { ref } from 'vue';

import { useVaultStore } from '@/stores/vault.ts';

import logger from '@/lib/logger.ts';

export function useVaultUnlockPageBase() {
    const vaultStore = useVaultStore();

    const passphrase = ref('');
    const unlocking = ref(false);
    const errorMessage = ref('');

    async function doUnlock(): Promise<boolean> {
        if (!passphrase.value || unlocking.value) return false;

        unlocking.value = true;
        errorMessage.value = '';

        try {
            await vaultStore.unlockVault(passphrase.value);
            logger.info('Vault unlocked');
            return true;
        } catch (e: unknown) {
            if (e instanceof Error && e.message.includes('tag')) {
                errorMessage.value = 'Wrong passphrase';
            } else {
                const msg = e instanceof Error ? e.message : 'Failed to unlock vault';
                errorMessage.value = msg;
            }
            logger.warn('Vault unlock failed: ' + errorMessage.value);
            return false;
        } finally {
            unlocking.value = false;
        }
    }

    return {
        passphrase,
        unlocking,
        errorMessage,
        doUnlock,
    };
}
