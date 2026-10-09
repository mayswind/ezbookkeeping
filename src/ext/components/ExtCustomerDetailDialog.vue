<template>
    <v-dialog width="760" scrollable :model-value="show" @update:model-value="close">
        <v-card v-if="customer">
            <v-card-title>{{ customer.name }}</v-card-title>
            <v-card-subtitle>
                <span v-if="customer.phone">{{ customer.phone }}</span>
                <span v-if="customer.phone && customer.email"> · </span>
                <span v-if="customer.email">{{ customer.email }}</span>
            </v-card-subtitle>
            <v-card-text>
                <div class="d-flex align-center flex-wrap ga-3 mb-4">
                    <div>
                        <div class="text-caption text-medium-emphasis">{{ tt('Currently owes') }}</div>
                        <div class="text-h5" :class="{ 'text-error': customer.outstanding > 0 }">{{ money(customer.outstanding) }}</div>
                    </div>
                    <v-spacer />
                    <v-btn color="primary" :disabled="customer.outstanding <= 0" @click="emit('repay')">{{ tt('Record repayment') }}</v-btn>
                </div>

                <v-progress-linear indeterminate v-if="loading" />

                <div class="text-subtitle-1 mb-1">{{ tt('Unpaid sales') }}</div>
                <v-table density="compact" class="mb-6">
                    <thead>
                        <tr>
                            <th>{{ tt('Date') }}</th>
                            <th class="text-end">{{ tt('Total') }}</th>
                            <th class="text-end">{{ tt('Paid') }}</th>
                            <th class="text-end">{{ tt('Owed') }}</th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr :key="sale.id" v-for="sale in openSales">
                            <td>{{ formatTime(sale.time) }}</td>
                            <td class="text-end">{{ money(sale.total) }}</td>
                            <td class="text-end">{{ money(sale.paid) }}</td>
                            <td class="text-end">{{ money(sale.outstanding) }}</td>
                        </tr>
                        <tr v-if="!loading && openSales.length < 1">
                            <td colspan="4" class="text-medium-emphasis">{{ tt('Nothing is owed.') }}</td>
                        </tr>
                    </tbody>
                </v-table>

                <div class="text-subtitle-1 mb-1">{{ tt('Repayments received') }}</div>
                <v-table density="compact">
                    <thead>
                        <tr>
                            <th>{{ tt('Date') }}</th>
                            <th class="text-end">{{ tt('Amount') }}</th>
                            <th>{{ tt('Received by') }}</th>
                            <th>{{ tt('Note') }}</th>
                            <th></th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr :key="repayment.id" v-for="repayment in repayments">
                            <td>{{ formatTime(repayment.time) }}</td>
                            <td class="text-end">{{ money(repayment.amount) }}</td>
                            <td>{{ nameOf(repayment.actorUid) || '–' }}</td>
                            <td>{{ repayment.note }}</td>
                            <td class="text-end"><v-btn size="small" variant="text" @click="emit('receipt', repayment)">{{ tt('Receipt') }}</v-btn></td>
                        </tr>
                        <tr v-if="!loading && repayments.length < 1">
                            <td colspan="5" class="text-medium-emphasis">{{ tt('No repayments yet.') }}</td>
                        </tr>
                    </tbody>
                </v-table>
            </v-card-text>
            <v-card-actions>
                <v-spacer />
                <v-btn variant="text" @click="close(false)">{{ tt('Close') }}</v-btn>
            </v-card-actions>
        </v-card>
    </v-dialog>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue';

import { useExtI18n } from '@/ext/i18n.ts';

import { parseBigDecimal } from '@/lib/numeral.ts';
import { parseDateTimeFromUnixTime } from '@/lib/datetime.ts';

import api from '@/ext/api.ts';
import { usePeople } from '@/ext/people.ts';
import type { CustomerInfo, RepaymentInfo, SaleInfo } from '@/ext/types.ts';

const props = defineProps<{
    show: boolean;
    customer: CustomerInfo | null;
    currency: string;
}>();

const emit = defineEmits<{
    (e: 'update:show', value: boolean): void;
    (e: 'repay'): void;
    (e: 'receipt', repayment: RepaymentInfo): void;
    (e: 'error', error: unknown): void;
}>();

const { tt, formatAmountToLocalizedNumeralsWithCurrency, formatDateTimeToLongDateTime } = useExtI18n();
const { load: loadPeople, nameOf } = usePeople();

const loading = ref<boolean>(false);
const openSales = ref<SaleInfo[]>([]);
const repayments = ref<RepaymentInfo[]>([]);

function money(minorUnits: number): string {
    return formatAmountToLocalizedNumeralsWithCurrency(parseBigDecimal(minorUnits), props.currency);
}

function formatTime(unixTime: number): string {
    return formatDateTimeToLongDateTime(parseDateTimeFromUnixTime(unixTime));
}

function close(value: boolean = false): void {
    emit('update:show', value);
}

async function load(): Promise<void> {
    if (!props.customer) {
        return;
    }

    loading.value = true;

    try {
        [openSales.value, repayments.value] = await Promise.all([
            api.listSales({ customerId: props.customer.id, onlyOpen: true }),
            api.listRepayments(props.customer.id),
            loadPeople()
        ]);
    } catch (error) {
        emit('error', error);
    } finally {
        loading.value = false;
    }
}

// reload when opened, and again when the customer's balance changes (a repayment was recorded)
watch(() => [props.show, props.customer?.outstanding], () => {
    if (props.show) {
        load();
    }
});
</script>
