import { ref } from 'vue';

import { AccountCategory } from '@/core/account.ts';

import { useSettingsStore } from '@/stores/setting.ts';

export function useAccountCategoryDisplayOrderSettingsPageBase() {
    const settingsStore = useSettingsStore();

    const accountCategories = ref<AccountCategory[]>(AccountCategory.values(settingsStore.appSettings.accountCategoryOrders));

    function isDisplayOrderModified(): boolean {
        const currentOrders = AccountCategory.values(settingsStore.appSettings.accountCategoryOrders);

        if (currentOrders.length !== accountCategories.value.length) {
            return true;
        }

        for (let i = 0; i < currentOrders.length; i++) {
            const accountCategory = accountCategories.value[i];
            const currentCategory = currentOrders[i];

            if (!accountCategory || !currentCategory) {
                return true;
            }

            if (accountCategory.type !== currentCategory.type) {
                return true;
            }
        }

        return false;
    }

    function loadDisplayOrderFromSettings(): void {
        accountCategories.value = AccountCategory.values(settingsStore.appSettings.accountCategoryOrders);
    }

    function saveDisplayOrderToSettings(): void {
        const displayOrders = accountCategories.value.map(category => category.type).join(',');
        const defaultOrders = AccountCategory.values('').map(category => category.type).join(',');

        if (displayOrders === defaultOrders) {
            settingsStore.setAccountCategoryOrders('');
        } else {
            settingsStore.setAccountCategoryOrders(displayOrders);
        }
    }

    function resetDisplayOrderToDefault(): void {
        accountCategories.value = AccountCategory.values('');
    }

    return {
        // states
        accountCategories,
        // functions
        isDisplayOrderModified,
        loadDisplayOrderFromSettings,
        saveDisplayOrderToSettings,
        resetDisplayOrderToDefault
    };
}
