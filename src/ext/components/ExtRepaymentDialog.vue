<template>
    <v-dialog width="520" persistent :model-value="show" @update:model-value="close">
        <v-card>
            <v-card-title>{{ tt('Record repayment') }}</v-card-title>
            <v-card-subtitle v-if="customer">
                {{ tt('{name} owes {amount}', { name: customer.name, amount: money(customer.outstanding) }) }}
            </v-card-subtitle>
            <v-card-text>
                <v-progress-linear indeterminate v-if="loading" />
                <v-form class="d-flex flex-column ga-4 pt-2" @submit.prevent="save">
                    <v-select density="compact" hide-details="auto" :label="tt('Apply to')"
                              :items="saleOptions" item-title="title" item-value="value"
                              :disabled="saving || loading" v-model="saleId" />

                    <amount-input density="compact" :currency="currency" :label="tt('Amount received')"
                                  :disabled="saving || loading" v-model="amount" />
                    <div class="text-caption text-error" style="margin-top: -12px" v-if="errors.amount">{{ errors.amount }}</div>

                    <v-select density="compact" hide-details="auto" :label="tt('Money goes into')"
                              :items="paymentAccountOptions" item-title="title" item-value="value"
                              :error-messages="errors.payment" :disabled="saving || loading" v-model="paymentAccountId" />
                    <v-select density="compact" hide-details="auto" :label="tt('Owed amount was tracked in')"
                              :items="receivableAccountOptions" item-title="title" item-value="value"
                              :error-messages="errors.receivable" :disabled="saving || loading" v-model="receivableAccountId" />

                    <v-alert type="info" variant="tonal" density="compact" v-if="transferCategoryOptions.length < 1 && !loading">
                        {{ tt('You have no transfer categories yet. A repayment is recorded as a transfer and needs one.') }}
                        <router-link class="ms-1" to="/category/list">{{ tt('Add one on the Categories page') }}</router-link>
                    </v-alert>
                    <v-select density="compact" hide-details="auto" :label="tt('Transfer category')"
                              :items="transferCategoryOptions" item-title="title" item-value="value"
                              :hint="tt('Repayments are recorded as a transfer from the owed-money account to where the money arrived.')"
                              :error-messages="errors.category" :disabled="saving || loading" v-model="categoryId" v-else />

                    <v-text-field density="compact" hide-details="auto" :label="tt('Note (optional)')" :disabled="saving || loading" v-model="note" />
                </v-form>
            </v-card-text>
            <v-card-actions>
                <v-spacer />
                <v-btn variant="text" :disabled="saving" @click="close(false)">{{ tt('Cancel') }}</v-btn>
                <v-btn color="primary" :loading="saving" :disabled="loading" @click="save">{{ tt('Record repayment') }}</v-btn>
            </v-card-actions>
        </v-card>
    </v-dialog>
</template>

<script setup lang="ts">
import { ref, reactive, computed, watch } from 'vue';

import { useExtI18n } from '@/ext/i18n.ts';

import { parseBigDecimal } from '@/lib/numeral.ts';
import { parseDateTimeFromUnixTime } from '@/lib/datetime.ts';

import api from '@/ext/api.ts';
import { useBusiness } from '@/ext/business.ts';
import { useBusinessAccounts } from '@/ext/accounts.ts';
import { loadFormDefaults, saveFormDefaults } from '@/ext/defaults.ts';
import type { CustomerInfo, RepaymentInfo, SaleInfo } from '@/ext/types.ts';

const props = defineProps<{
    show: boolean;
    customer: CustomerInfo | null;
}>();

const emit = defineEmits<{
    (e: 'update:show', value: boolean): void;
    (e: 'saved', repayment: RepaymentInfo): void;
    (e: 'error', error: unknown): void;
}>();

const { tt, formatAmountToLocalizedNumeralsWithCurrency, formatDateTimeToLongDate } = useExtI18n();
const { current } = useBusiness();

const loading = ref<boolean>(false);
const saving = ref<boolean>(false);
const openSales = ref<SaleInfo[]>([]);
const saleId = ref<string>('0');
const amount = ref<number>(0);
const paymentAccountId = ref<string>('');
const receivableAccountId = ref<string>('');
const categoryId = ref<string>('');
const note = ref<string>('');
const errors = reactive<{ amount: string, payment: string, receivable: string, category: string }>({ amount: '', payment: '', receivable: '', category: '' });

const {
    paymentAccounts, receivableAccounts, compatibleReceivables,
    paymentAccountOptions, receivableAccountOptions, transferCategoryOptions, currency, load: loadAccounts
} = useBusinessAccounts(paymentAccountId);

function money(minorUnits: number): string {
    return formatAmountToLocalizedNumeralsWithCurrency(parseBigDecimal(minorUnits), currency.value);
}

const saleOptions = computed(() => [
    { title: tt('Oldest sales first'), value: '0' },
    ...openSales.value.map(sale => ({
        title: tt('Sale #{id} on {date}: owes {amount}', { id: sale.id, date: formatDateTimeToLongDate(parseDateTimeFromUnixTime(sale.time)), amount: money(sale.outstanding) }),
        value: sale.id
    }))
]);

// the most that can be received now: everything the customer owes, or what is left on the chosen sale
const maxAmount = computed<number>(() => {
    if (saleId.value !== '0') {
        return openSales.value.find(s => s.id === saleId.value)?.outstanding ?? 0;
    }

    return props.customer?.outstanding ?? 0;
});

watch(saleId, () => {
    amount.value = maxAmount.value;
});

watch(compatibleReceivables, accounts => {
    if (!accounts.some(a => a.id === receivableAccountId.value)) {
        receivableAccountId.value = accounts[0]?.id ?? '';
    }
});

watch(() => props.show, async open => {
    if (!open || !props.customer) {
        return;
    }

    errors.amount = errors.payment = errors.receivable = errors.category = '';
    saleId.value = '0';
    amount.value = props.customer.outstanding;
    note.value = '';
    loading.value = true;

    try {
        [openSales.value] = await Promise.all([api.listSales({ customerId: props.customer.id, onlyOpen: true }), loadAccounts()]);

        const saved = loadFormDefaults(current.value?.ownerUid ?? 'own');
        paymentAccountId.value = paymentAccounts.value.some(a => a.id === saved.payment) ? saved.payment! : (paymentAccounts.value[0]?.id ?? '');
        receivableAccountId.value = receivableAccounts.value.some(a => a.id === saved.receivable) ? saved.receivable! : (compatibleReceivables.value[0]?.id ?? '');
        categoryId.value = transferCategoryOptions.value.some(c => c.value === saved.transfer) ? saved.transfer! : (transferCategoryOptions.value[0]?.value ?? '');
    } catch (error) {
        emit('error', error);
    } finally {
        loading.value = false;
    }
});

function close(value: boolean = false): void {
    emit('update:show', value);
}

async function save(): Promise<void> {
    if (!props.customer) {
        return;
    }

    errors.amount = amount.value <= 0 ? tt('Enter the amount received')
        : (amount.value > maxAmount.value ? tt('That is more than is owed: {amount}', { amount: money(maxAmount.value) }) : '');
    errors.payment = paymentAccountId.value ? '' : tt('Choose where the money goes');
    errors.receivable = compatibleReceivables.value.some(a => a.id === receivableAccountId.value) ? '' : tt('Choose the owed-money account');
    errors.category = categoryId.value ? '' : tt('Choose a transfer category');

    if (errors.amount || errors.payment || errors.receivable || errors.category) {
        return;
    }

    saving.value = true;

    try {
        const repayment = await api.addRepayment({
            customerId: props.customer.id,
            amount: amount.value,
            saleId: saleId.value,
            time: Math.floor(Date.now() / 1000),
            utcOffset: -new Date().getTimezoneOffset(),
            paymentAccountId: paymentAccountId.value,
            receivableAccountId: receivableAccountId.value,
            categoryId: categoryId.value,
            note: note.value.trim()
        });

        saveFormDefaults(current.value?.ownerUid ?? 'own', { payment: paymentAccountId.value, receivable: receivableAccountId.value, transfer: categoryId.value });
        emit('saved', repayment);
        close(false);
    } catch (error) {
        emit('error', error);
    } finally {
        saving.value = false;
    }
}
</script>
