import { ref, computed } from 'vue';

import type { TransactionCategory } from '@/models/transaction_category.ts';

import { useTransactionCategoriesStore } from '@/stores/transactionCategory.ts';

export function useCategoryListPageBase() {
    const transactionCategoriesStore = useTransactionCategoriesStore();

    const loading = ref<boolean>(true);
    const primaryCategoryId = ref<string>('0');

    const currentPrimaryCategory = computed<TransactionCategory | undefined>(() => {
        if (!transactionCategoriesStore.allTransactionCategoriesMap || !transactionCategoriesStore.allTransactionCategoriesMap[primaryCategoryId.value]) {
            return undefined;
        }

        return transactionCategoriesStore.allTransactionCategoriesMap[primaryCategoryId.value];
    });

    return {
        // states
        loading,
        primaryCategoryId,
        // computed states
        currentPrimaryCategory
    };
}
