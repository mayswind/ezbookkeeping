import { computed, type Ref } from 'vue';

import { useAccountsStore } from '@/stores/account.ts';
import { useTransactionCategoriesStore } from '@/stores/transactionCategory.ts';
import { useUserStore } from '@/stores/user.ts';

import { CategoryType } from '@/core/category.ts';
import { AccountCategory } from '@/core/account.ts';
import type { TransactionCategory } from '@/models/transaction_category.ts';

export interface PickerOption {
    readonly title: string;
    readonly value: string;
}

function flattenCategories(primaries: TransactionCategory[] | undefined): PickerOption[] {
    const options: PickerOption[] = [];

    for (const primary of primaries ?? []) {
        if (primary.hidden) {
            continue;
        }

        for (const secondary of primary.subCategories ?? []) {
            if (!secondary.hidden) {
                options.push({ title: `${primary.name} › ${secondary.name}`, value: secondary.id });
            }
        }
    }

    return options;
}

/**
 * Account and category pickers for recording money, built on the app's own stores (which already follow the
 * selected business). Sales and repayments book into one currency, so the receivables choices follow the
 * payment account's currency.
 */
export function useBusinessAccounts(paymentAccountId: Ref<string>) {
    const userStore = useUserStore();
    const accountsStore = useAccountsStore();
    const categoriesStore = useTransactionCategoriesStore();

    const paymentAccounts = computed(() => accountsStore.allVisiblePlainAccounts.filter(a => a.isAsset && a.category !== AccountCategory.Receivables.type));
    const receivableAccounts = computed(() => accountsStore.allVisiblePlainAccounts.filter(a => a.category === AccountCategory.Receivables.type));
    const selectedPaymentAccount = computed(() => paymentAccounts.value.find(a => a.id === paymentAccountId.value));

    const compatibleReceivables = computed(() => {
        const currency = selectedPaymentAccount.value?.currency;
        return receivableAccounts.value.filter(a => !currency || a.currency === currency);
    });

    const paymentAccountOptions = computed<PickerOption[]>(() => paymentAccounts.value.map(a => ({ title: `${a.name} (${a.currency})`, value: a.id })));
    const receivableAccountOptions = computed<PickerOption[]>(() => compatibleReceivables.value.map(a => ({ title: `${a.name} (${a.currency})`, value: a.id })));
    const incomeCategoryOptions = computed<PickerOption[]>(() => flattenCategories(categoriesStore.allTransactionCategories[CategoryType.Income]));
    const transferCategoryOptions = computed<PickerOption[]>(() => flattenCategories(categoriesStore.allTransactionCategories[CategoryType.Transfer]));

    const currency = computed<string>(() => selectedPaymentAccount.value?.currency
        ?? compatibleReceivables.value[0]?.currency
        ?? userStore.currentUserDefaultCurrency);

    async function load(force: boolean = false): Promise<void> {
        await Promise.all([accountsStore.loadAllAccounts({ force }), categoriesStore.loadAllCategories({ force })]);
    }

    return {
        paymentAccounts, receivableAccounts, selectedPaymentAccount, compatibleReceivables,
        paymentAccountOptions, receivableAccountOptions, incomeCategoryOptions, transferCategoryOptions, currency, load
    };
}
