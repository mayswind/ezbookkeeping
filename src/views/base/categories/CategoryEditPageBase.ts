import { ref, computed } from 'vue';

import { useTransactionCategoriesStore } from '@/stores/transactionCategory.ts';

import { CategoryType } from '@/core/category.ts';
import { TransactionCategory } from '@/models/transaction_category.ts';

import { allVisiblePrimaryTransactionCategoriesByType } from '@/lib/category.ts';

export function useCategoryEditPageBase(type?: CategoryType, parentId?: string) {
    const transactionCategoriesStore = useTransactionCategoriesStore();

    const editCategoryId = ref<string | null>(null);
    const clientSessionId = ref<string>('');
    const loading = ref<boolean>(false);
    const submitting = ref<boolean>(false);
    const category = ref<TransactionCategory>(TransactionCategory.createNewCategory(type, parentId));

    const allAvailableCategories = computed<TransactionCategory[]>(() => allVisiblePrimaryTransactionCategoriesByType(transactionCategoriesStore.allTransactionCategories, category.value.type));

    const title = computed<string>(() => {
        if (!editCategoryId.value) {
            if (category.value.parentId === '0') {
                return 'Add Primary Category';
            } else {
                return 'Add Secondary Category';
            }
        } else {
            return 'Edit Category';
        }
    });

    const saveButtonTitle = computed<string>(() => {
        if (!editCategoryId.value) {
            return 'Add';
        } else {
            return 'Save';
        }
    });

    const inputEmptyProblemMessage = computed<string | null>(() => {
        if (!category.value.name) {
            return 'Category name cannot be blank';
        } else {
            return null;
        }
    });

    const inputIsEmpty = computed<boolean>(() => !!inputEmptyProblemMessage.value);

    return {
        // states
        editCategoryId,
        clientSessionId,
        loading,
        submitting,
        category,
        // computed states
        allAvailableCategories,
        title,
        saveButtonTitle,
        inputEmptyProblemMessage,
        inputIsEmpty
    };
}
