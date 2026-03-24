import { ref, computed } from 'vue';
import { defineStore } from 'pinia';

import {
    isVaultUnlocked,
    setupVault as setupVaultService,
    unlock as unlockVaultService,
    lock as lockVaultService,
    type VaultParamsDTO
} from '@/lib/vault-service.ts';

import logger from '@/lib/logger.ts';
import services from '@/lib/services.ts';

export const useVaultStore = defineStore('vault', () => {
    const hasVault = ref(false);
    const tier = ref('free');
    const accountBalances = ref<Map<number, number>>(new Map());

    const vaultUnlocked = computed(() => isVaultUnlocked());
    const canUseAI = computed(() => ['alfred', 'maurice', 'jared'].includes(tier.value));
    const canUseGroups = computed(() => ['maurice', 'jared'].includes(tier.value));

    function setFromAuthResponse(authHasVault: boolean, authTier: string): void {
        hasVault.value = authHasVault;
        tier.value = authTier || 'free';
    }

    async function initVault(passphrase: string): Promise<void> {
        const dto = await setupVaultService(passphrase);

        const response = await services.vaultInit(dto);
        const data = response.data;

        if (!data || !data.success) {
            throw new Error('Failed to initialize vault on server');
        }

        hasVault.value = true;
        logger.info('Vault initialized');
    }

    async function unlockVault(passphrase: string): Promise<void> {
        const response = await services.vaultGetParams();
        const data = response.data;

        if (!data || !data.success || !data.result) {
            throw new Error('Failed to fetch vault params');
        }

        const dto = data.result as VaultParamsDTO;
        await unlockVaultService(passphrase, dto);
        logger.info('Vault unlocked');
    }

    function lockVault(): void {
        lockVaultService();
        accountBalances.value = new Map();
        logger.info('Vault locked');
    }

    async function linkApiKey(apiKey: string): Promise<string> {
        const response = await services.linkApiKey({ apiKey });
        const data = response.data;

        if (!data || !data.success || !data.result) {
            throw new Error('Failed to link API key');
        }

        tier.value = data.result.tier;
        return data.result.tier;
    }

    function updateAccountBalance(accountId: number, balance: number): void {
        accountBalances.value.set(accountId, balance);
    }

    function resetVaultState(): void {
        lockVaultService();
        hasVault.value = false;
        tier.value = 'free';
        accountBalances.value = new Map();
    }

    return {
        hasVault,
        tier,
        accountBalances,
        vaultUnlocked,
        canUseAI,
        canUseGroups,
        setFromAuthResponse,
        initVault,
        unlockVault,
        lockVault,
        linkApiKey,
        updateAccountBalance,
        resetVaultState,
    };
});
