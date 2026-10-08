<template>
    <main-page-layout>
        <template #content>
            <v-row class="match-height">
                <v-col cols="12">
                    <v-card>
                        <template #title>
                            <div class="d-flex align-center flex-wrap ga-3">
                                <span class="me-2">{{ tt('Customers') }}</span>
                                <v-text-field style="max-width: 260px" density="compact" hide-details variant="outlined"
                                              :label="tt('Search by name, phone or email')" clearable v-model="search" />
                                <v-switch density="compact" hide-details color="primary" :label="tt('Only people who owe me')" v-model="onlyOwing" />
                                <v-spacer />
                                <v-btn color="primary" :disabled="loading" @click="openCustomerDialog(null)">{{ tt('Add customer') }}</v-btn>
                            </div>
                        </template>
                        <template #subtitle>
                            <span v-if="totalOwed > 0">{{ tt('Customers owe you {amount} in total', { amount: money(totalOwed) }) }}</span>
                            <span v-else>{{ tt('Nobody owes you anything right now.') }}</span>
                        </template>

                        <v-progress-linear indeterminate v-if="loading" />

                        <v-table>
                            <thead>
                                <tr>
                                    <th>{{ tt('Name') }}</th>
                                    <th class="text-end">{{ tt('Owes') }}</th>
                                    <th></th>
                                </tr>
                            </thead>
                            <tbody>
                                <tr class="cursor-pointer" :key="customer.id" v-for="customer in visibleCustomers" @click="openDetail(customer)">
                                    <td>
                                        {{ customer.name }}
                                        <div class="text-caption text-medium-emphasis" v-if="customer.phone || customer.email">
                                            {{ [customer.phone, customer.email].filter(Boolean).join(' · ') }}
                                        </div>
                                    </td>
                                    <td class="text-end" :class="{ 'text-error font-weight-medium': customer.outstanding > 0 }">
                                        {{ customer.outstanding > 0 ? money(customer.outstanding) : '–' }}
                                    </td>
                                    <td class="text-end text-no-wrap" @click.stop>
                                        <v-btn size="small" variant="tonal" color="primary" v-if="customer.outstanding > 0"
                                               @click="openRepayment(customer)">{{ tt('Record repayment') }}</v-btn>
                                        <v-btn size="small" variant="text" @click="openCustomerDialog(customer)" v-if="canManage">{{ tt('Edit') }}</v-btn>
                                        <v-btn size="small" variant="text" color="error" v-if="canManage"
                                               :disabled="customer.outstanding > 0" @click="removeCustomer(customer)">{{ tt('Delete') }}</v-btn>
                                    </td>
                                </tr>
                                <tr v-if="!loading && visibleCustomers.length < 1">
                                    <td colspan="3" class="text-center text-medium-emphasis py-6">
                                        {{ customers.length < 1 ? tt('No customers yet.') : tt('No customers match.') }}
                                    </td>
                                </tr>
                            </tbody>
                        </v-table>
                    </v-card>
                </v-col>
            </v-row>

            <ext-customer-dialog :customer="editingCustomer" v-model:show="showCustomerDialog"
                                 @saved="onSaved(tt('Customer saved'))" @error="onError" />
            <ext-customer-detail-dialog :customer="detailCustomer" :currency="currency" v-model:show="showDetail"
                                        @repay="repayFromDetail" @error="onError" />
            <ext-repayment-dialog :customer="repaymentCustomer" v-model:show="showRepayment"
                                  @saved="onSaved(tt('Repayment recorded'))" @error="onError" />
            <confirm-dialog ref="confirmDialog" />
            <ext-snack-bar ref="snackbar" />
        </template>
    </main-page-layout>
</template>

<script setup lang="ts">
import ConfirmDialog from '@/components/desktop/ConfirmDialog.vue';
import ExtSnackBar from '@/ext/components/ExtSnackBar.vue';
import ExtCustomerDialog from '@/ext/components/ExtCustomerDialog.vue';
import ExtCustomerDetailDialog from '@/ext/components/ExtCustomerDetailDialog.vue';
import ExtRepaymentDialog from '@/ext/components/ExtRepaymentDialog.vue';

import { ref, computed, onMounted, useTemplateRef } from 'vue';

import { useExtI18n } from '@/ext/i18n.ts';
import { useUserStore } from '@/stores/user.ts';

import { parseBigDecimal } from '@/lib/numeral.ts';

import api from '@/ext/api.ts';
import { useBusiness } from '@/ext/business.ts';
import type { CustomerInfo } from '@/ext/types.ts';

type ConfirmDialogType = InstanceType<typeof ConfirmDialog>;
type SnackBarType = InstanceType<typeof ExtSnackBar>;

const { tt, formatAmountToLocalizedNumeralsWithCurrency } = useExtI18n();
const userStore = useUserStore();
const { canManage, ensureLoaded } = useBusiness();

const confirmDialog = useTemplateRef<ConfirmDialogType>('confirmDialog');
const snackbar = useTemplateRef<SnackBarType>('snackbar');

const loading = ref<boolean>(true);
const customers = ref<CustomerInfo[]>([]);
const search = ref<string | null>('');
const onlyOwing = ref<boolean>(false);

const showCustomerDialog = ref<boolean>(false);
const editingCustomer = ref<CustomerInfo | null>(null);
const showDetail = ref<boolean>(false);
const detailCustomer = ref<CustomerInfo | null>(null);
const showRepayment = ref<boolean>(false);
const repaymentCustomer = ref<CustomerInfo | null>(null);

// Amounts are in the business currency; until the server says which, the user's default currency is shown
const currency = computed<string>(() => userStore.currentUserDefaultCurrency);

const totalOwed = computed<number>(() => customers.value.reduce((sum, c) => sum + c.outstanding, 0));

const visibleCustomers = computed<CustomerInfo[]>(() => {
    const needle = (search.value ?? '').trim().toLowerCase();

    return customers.value
        .filter(c => (!onlyOwing.value || c.outstanding > 0)
            && (!needle || [c.name, c.phone, c.email].some(text => text.toLowerCase().includes(needle))))
        .sort((a, b) => b.outstanding - a.outstanding || a.name.localeCompare(b.name));
});

function money(minorUnits: number): string {
    return formatAmountToLocalizedNumeralsWithCurrency(parseBigDecimal(minorUnits), currency.value);
}

async function load(): Promise<void> {
    loading.value = true;

    try {
        await ensureLoaded();
        customers.value = await api.listCustomers();

        // keep an open dialog pointing at the refreshed balance
        if (detailCustomer.value) {
            detailCustomer.value = customers.value.find(c => c.id === detailCustomer.value!.id) ?? detailCustomer.value;
        }

        if (repaymentCustomer.value) {
            repaymentCustomer.value = customers.value.find(c => c.id === repaymentCustomer.value!.id) ?? repaymentCustomer.value;
        }
    } catch (error) {
        snackbar.value?.showError(error);
    } finally {
        loading.value = false;
    }
}

function openCustomerDialog(customer: CustomerInfo | null): void {
    editingCustomer.value = customer;
    showCustomerDialog.value = true;
}

function openDetail(customer: CustomerInfo): void {
    detailCustomer.value = customer;
    showDetail.value = true;
}

function openRepayment(customer: CustomerInfo): void {
    repaymentCustomer.value = customer;
    showRepayment.value = true;
}

function repayFromDetail(): void {
    if (detailCustomer.value) {
        showDetail.value = false;
        openRepayment(detailCustomer.value);
    }
}

function removeCustomer(customer: CustomerInfo): void {
    confirmDialog.value?.open('ext.confirmDeleteCustomer', { name: customer.name }).then(async () => {
        try {
            await api.deleteCustomer(customer.id);
            snackbar.value?.showMessage(tt('Customer deleted'));
            await load();
        } catch (error) {
            snackbar.value?.showError(error);
        }
    }).catch(() => {
        // cancelled
    });
}

function onSaved(message: string): void {
    snackbar.value?.showMessage(message);
    load();
}

function onError(error: unknown): void {
    snackbar.value?.showError(error);
}

onMounted(load);
</script>
