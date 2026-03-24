import { ref, computed, onMounted } from 'vue';

import { useVaultStore } from '@/stores/vault.ts';

import logger from '@/lib/logger.ts';

export function useVaultSetupPageBase() {
    const vaultStore = useVaultStore();

    const passphrase = ref('');
    const confirmPassphrase = ref('');
    const warningAccepted = ref(false);
    const creating = ref(false);
    const errorMessage = ref('');
    const strengthScore = ref(-1);
    const strengthFeedback = ref('');

    let zxcvbnFn: ((password: string, userInputs?: string[]) => import('zxcvbn').ZXCVBNResult) | null = null;

    onMounted(async () => {
        try {
            const mod = await import('zxcvbn');
            zxcvbnFn = mod.default || mod;
        } catch (e) {
            logger.warn('Failed to load zxcvbn, passphrase strength check disabled');
        }
    });

    function updateStrength(): void {
        if (!passphrase.value || !zxcvbnFn) {
            strengthScore.value = -1;
            strengthFeedback.value = '';
            return;
        }

        const result = zxcvbnFn(passphrase.value);
        strengthScore.value = result.score;

        if (result.feedback.warning) {
            strengthFeedback.value = result.feedback.warning;
        } else if (result.feedback.suggestions.length > 0) {
            strengthFeedback.value = result.feedback.suggestions[0] || '';
        } else {
            strengthFeedback.value = '';
        }
    }

    const passphraseMatch = computed(() => {
        if (!confirmPassphrase.value) return true;
        return passphrase.value === confirmPassphrase.value;
    });

    const canSubmit = computed(() => {
        return passphrase.value.length >= 8
            && confirmPassphrase.value === passphrase.value
            && warningAccepted.value
            && strengthScore.value >= 3
            && !creating.value;
    });

    async function doSetup(): Promise<boolean> {
        if (!canSubmit.value) return false;

        creating.value = true;
        errorMessage.value = '';

        try {
            await vaultStore.initVault(passphrase.value);
            logger.info('Vault setup complete');
            return true;
        } catch (e: unknown) {
            const msg = e instanceof Error ? e.message : 'Failed to create vault';
            errorMessage.value = msg;
            logger.error('Vault setup failed: ' + msg);
            return false;
        } finally {
            creating.value = false;
        }
    }

    return {
        passphrase,
        confirmPassphrase,
        warningAccepted,
        creating,
        errorMessage,
        strengthScore,
        strengthFeedback,
        passphraseMatch,
        canSubmit,
        updateStrength,
        doSetup,
    };
}
