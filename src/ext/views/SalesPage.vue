<template>
    <main-page-layout>
        <template #content>
            <v-row class="match-height">
                <!-- pick items -->
                <v-col cols="12" md="7">
                    <v-card>
                        <template #title>
                            <div class="d-flex align-center flex-wrap ga-3">
                                <span class="me-2">{{ tt('New sale') }}</span>
                                <v-text-field style="max-width: 280px" density="compact" hide-details variant="outlined"
                                              :label="tt('Search by name or SKU')" clearable v-model="search" />
                            </div>
                        </template>
                        <v-progress-linear indeterminate v-if="loading" />
                        <v-table density="comfortable">
                            <thead>
                                <tr>
                                    <th>{{ tt('Item') }}</th>
                                    <th class="text-end">{{ tt('Available') }}</th>
                                    <th class="text-end">{{ tt('Price') }}</th>
                                    <th></th>
                                </tr>
                            </thead>
                            <tbody>
                                <tr :key="item.id" v-for="item in visibleItems">
                                    <td>
                                        {{ item.name }}
                                        <div class="text-caption text-medium-emphasis">{{ item.sku }}</div>
                                    </td>
                                    <td class="text-end">
                                        <template v-if="item.trackStock">
                                            <span :class="{ 'text-error': availableOf(item) <= 0 }">{{ formatQty(availableOf(item)) }}</span>
                                            <span class="text-medium-emphasis" v-if="item.unit"> {{ item.unit }}</span>
                                        </template>
                                        <span class="text-medium-emphasis" v-else>–</span>
                                    </td>
                                    <td class="text-end">{{ money(item.salePrice) }}</td>
                                    <td class="text-end">
                                        <v-btn size="small" color="primary" variant="tonal"
                                               :disabled="item.trackStock && availableOf(item) <= 0"
                                               @click="addToCart(item)">{{ tt('Add') }}</v-btn>
                                    </td>
                                </tr>
                                <tr v-if="!loading && visibleItems.length < 1">
                                    <td colspan="4" class="text-center text-medium-emphasis py-6">
                                        {{ items.length < 1 ? tt('No items yet. Add items on the Inventory page first.') : tt('No items match your search.') }}
                                    </td>
                                </tr>
                            </tbody>
                        </v-table>
                    </v-card>
                </v-col>

                <!-- cart and payment -->
                <v-col cols="12" md="5">
                    <v-card :title="tt('Cart')">
                        <v-card-text>
                            <v-select density="compact" :label="tt('Sell from')" :items="locationOptions" item-title="title" item-value="value"
                                      :disabled="submitting" v-model="locationId" v-if="locations.length > 1" />

                            <div class="text-medium-emphasis py-4" v-if="cart.length < 1">{{ tt('Add items from the list to start a sale.') }}</div>

                            <div class="mb-3" :key="line.itemId" v-for="line in cart">
                                <div class="d-flex align-center ga-2">
                                    <div class="flex-grow-1">
                                        <div class="font-weight-medium">{{ line.item.name }}</div>
                                        <div class="text-caption text-medium-emphasis">
                                            {{ money(line.unitPrice) }}<span v-if="line.item.unit"> / {{ line.item.unit }}</span>
                                        </div>
                                    </div>
                                    <v-text-field style="max-width: 96px" density="compact" hide-details variant="outlined" inputmode="decimal"
                                                  :aria-label="tt('Quantity')" :error="line.qty === null || line.qtyProblem !== ''"
                                                  :disabled="submitting" :model-value="line.qtyText"
                                                  @update:model-value="setQty(line.itemId, $event)" />
                                    <div class="text-end" style="min-width: 84px">{{ money(line.total) }}</div>
                                    <v-btn size="small" variant="text" :icon="true" :aria-label="tt('Remove')"
                                           :disabled="submitting" @click="removeLine(line.itemId)">
                                        <v-icon :icon="mdiClose" size="18" />
                                    </v-btn>
                                </div>
                                <div class="text-caption text-error" v-if="line.qtyProblem">{{ line.qtyProblem }}</div>
                                <div class="mt-1" v-if="canManage">
                                    <amount-input density="compact" :currency="currency" :label="tt('Price override')"
                                                  :disabled="submitting" :model-value="line.unitPrice"
                                                  @update:model-value="setPrice(line.itemId, $event)" />
                                </div>
                            </div>

                            <template v-if="cart.length > 0">
                                <v-divider class="my-3" />
                                <div class="d-flex justify-space-between"><span>{{ tt('Subtotal') }}</span><span>{{ money(totals.subtotal) }}</span></div>
                                <amount-input class="mt-2" density="compact" :currency="currency" :label="tt('Discount')"
                                              :disabled="submitting" v-model="discount" v-if="canManage" />
                                <div class="d-flex justify-space-between text-h6 mt-2"><span>{{ tt('Total') }}</span><span>{{ money(totals.total) }}</span></div>

                                <v-divider class="my-3" />
                                <div class="d-flex flex-column ga-4">
                                <v-btn-toggle class="w-100" mandatory density="comfortable" color="primary" variant="outlined" divided
                                              :disabled="submitting" v-model="payMode">
                                    <v-btn class="flex-grow-1" value="full">{{ tt('Paid in full') }}</v-btn>
                                    <v-btn class="flex-grow-1" value="partial">{{ tt('Part payment') }}</v-btn>
                                    <v-btn class="flex-grow-1" value="credit">{{ tt('On credit') }}</v-btn>
                                </v-btn-toggle>

                                <amount-input density="compact" :currency="currency" :label="tt('Amount paid now')"
                                              :disabled="submitting" v-model="partialPaid" v-if="payMode === 'partial'" />

                                <v-select density="compact" hide-details="auto" :label="tt('Money goes into')" :items="paymentAccountOptions" item-title="title" item-value="value"
                                          :disabled="submitting" :error-messages="problems.payment" v-model="paymentAccountId" v-if="totals.paid > 0" />

                                <template v-if="totals.credit > 0">
                                    <v-alert type="warning" variant="tonal" density="compact">
                                        {{ tt('{amount} will be recorded as owed by the customer.', { amount: money(totals.credit) }) }}
                                    </v-alert>
                                    <div class="d-flex align-start ga-2">
                                        <v-autocomplete class="flex-grow-1" density="compact" hide-details="auto" :label="tt('Customer')"
                                                        :items="customerOptions" item-title="title" item-value="value" clearable
                                                        :disabled="submitting" :error-messages="problems.customer" v-model="customerId" />
                                        <v-btn size="small" variant="tonal" :disabled="submitting" @click="showCustomerDialog = true">{{ tt('New') }}</v-btn>
                                    </div>
                                    <v-select density="compact" hide-details="auto" :label="tt('Owed amount is tracked in')" :items="receivableAccountOptions" item-title="title" item-value="value"
                                              :disabled="submitting" :error-messages="problems.receivable" v-model="receivableAccountId" />
                                </template>
                                <div class="d-flex align-start ga-2" v-else>
                                    <v-autocomplete class="flex-grow-1" density="compact" hide-details="auto" :label="tt('Customer (optional)')"
                                                    :items="customerOptions" item-title="title" item-value="value" clearable
                                                    :disabled="submitting" v-model="customerId" />
                                    <v-btn size="small" variant="tonal" :disabled="submitting" @click="showCustomerDialog = true">{{ tt('New') }}</v-btn>
                                </div>

                                <v-alert type="info" variant="tonal" density="compact" v-if="categoryOptions.length < 1">
                                    {{ tt('You have no income categories yet. A sale must be recorded under one, such as "Product sales".') }}
                                    <router-link class="ms-1" to="/category/list">{{ tt('Add one on the Categories page') }}</router-link>
                                </v-alert>
                                <v-select density="compact" hide-details="auto" :label="tt('Income category')" :items="categoryOptions" item-title="title" item-value="value"
                                          :hint="tt('Where this sale appears in your income reports. Manage categories on the Categories page.')"
                                          :disabled="submitting" :error-messages="problems.category" v-model="categoryId" v-else />
                                <v-text-field density="compact" hide-details="auto" :label="tt('Note (optional)')" :disabled="submitting" v-model="note" />
                                </div>

                                <v-btn block class="mt-4" size="large" color="primary" :loading="submitting" :disabled="!canSubmit" @click="submit">
                                    {{ tt('Complete sale') }} · {{ money(totals.total) }}
                                </v-btn>
                                <v-btn block class="mt-2" variant="text" :disabled="submitting" @click="clearCart">{{ tt('Clear cart') }}</v-btn>
                            </template>
                        </v-card-text>
                    </v-card>
                </v-col>

                <!-- history -->
                <v-col cols="12">
                    <v-card :title="tt('Recent sales')">
                        <v-table density="comfortable">
                            <thead>
                                <tr>
                                    <th>{{ tt('When') }}</th>
                                    <th>{{ tt('Customer') }}</th>
                                    <th class="text-end">{{ tt('Total') }}</th>
                                    <th class="text-end">{{ tt('Owed') }}</th>
                                    <th>{{ tt('Status') }}</th>
                                    <th>{{ tt('Recorded by') }}</th>
                                    <th v-if="canManage"></th>
                                </tr>
                            </thead>
                            <tbody>
                                <tr :key="sale.id" v-for="sale in sales">
                                    <td>{{ formatTime(sale.time) }}</td>
                                    <td>{{ customerName(sale.customerId) }}</td>
                                    <td class="text-end">{{ money(sale.total) }}</td>
                                    <td class="text-end">{{ sale.outstanding > 0 ? money(sale.outstanding) : '–' }}</td>
                                    <td>
                                        <v-chip size="small" color="error" v-if="sale.voided">{{ tt('Voided') }}</v-chip>
                                        <v-chip size="small" color="warning" v-else-if="sale.outstanding > 0">{{ tt('On credit') }}</v-chip>
                                        <v-chip size="small" color="success" v-else>{{ tt('Paid') }}</v-chip>
                                    </td>
                                    <td>{{ nameOf(sale.actorUid) || '–' }}</td>
                                    <td class="text-end" v-if="canManage">
                                        <v-btn size="small" variant="text" color="error" :disabled="sale.voided"
                                               @click="voidSale(sale)">{{ tt('Void') }}</v-btn>
                                    </td>
                                </tr>
                                <tr v-if="sales.length < 1">
                                    <td :colspan="canManage ? 7 : 6" class="text-center text-medium-emphasis py-6">{{ tt('No sales yet.') }}</td>
                                </tr>
                            </tbody>
                        </v-table>
                    </v-card>
                </v-col>
            </v-row>

            <ext-customer-dialog v-model:show="showCustomerDialog" @saved="onCustomerSaved" @error="onError" />
            <confirm-dialog ref="confirmDialog" />
            <ext-snack-bar ref="snackbar" />
        </template>
    </main-page-layout>
</template>

<script setup lang="ts">
import ConfirmDialog from '@/components/desktop/ConfirmDialog.vue';
import ExtSnackBar from '@/ext/components/ExtSnackBar.vue';
import ExtCustomerDialog from '@/ext/components/ExtCustomerDialog.vue';

import { ref, reactive, computed, watch, onMounted, useTemplateRef } from 'vue';

import { useExtI18n } from '@/ext/i18n.ts';
import { parseBigDecimal } from '@/lib/numeral.ts';
import { parseDateTimeFromUnixTime } from '@/lib/datetime.ts';

import api from '@/ext/api.ts';
import { useBusiness } from '@/ext/business.ts';
import { useBusinessAccounts } from '@/ext/accounts.ts';
import { usePeople } from '@/ext/people.ts';
import { loadFormDefaults, saveFormDefaults } from '@/ext/defaults.ts';
import { formatQty, parseQty } from '@/ext/qty.ts';
import { computeTotals, lineTotal, type PayMode } from '@/ext/money.ts';
import type { CustomerInfo, ItemInfo, LocationInfo, SaleInfo, StockLevel } from '@/ext/types.ts';

import { mdiClose } from '@mdi/js';

type ConfirmDialogType = InstanceType<typeof ConfirmDialog>;
type SnackBarType = InstanceType<typeof ExtSnackBar>;

interface CartLine {
    itemId: string;
    item: ItemInfo;
    qtyText: string;
    qty: number | null;
    priceOverride: number | null;
    unitPrice: number;
    total: number;
    qtyProblem: string;
}

const { tt, formatAmountToLocalizedNumeralsWithCurrency, formatDateTimeToLongDateTime } = useExtI18n();
const { canManage, current, ensureLoaded } = useBusiness();
const { load: loadPeople, nameOf } = usePeople();

const confirmDialog = useTemplateRef<ConfirmDialogType>('confirmDialog');
const snackbar = useTemplateRef<SnackBarType>('snackbar');

const loading = ref<boolean>(true);
const submitting = ref<boolean>(false);
const items = ref<ItemInfo[]>([]);
const locations = ref<LocationInfo[]>([]);
const levels = ref<StockLevel[]>([]);
const customers = ref<CustomerInfo[]>([]);
const sales = ref<SaleInfo[]>([]);
const search = ref<string | null>('');

const cartInput = ref<{ itemId: string, qtyText: string, priceOverride: number | null }[]>([]);
const locationId = ref<string>('');
const payMode = ref<PayMode>('full');
const partialPaid = ref<number>(0);
const discount = ref<number>(0);
const customerId = ref<string | null>(null);
const paymentAccountId = ref<string>('');
const receivableAccountId = ref<string>('');
const categoryId = ref<string>('');
const note = ref<string>('');
const showCustomerDialog = ref<boolean>(false);
const showProblems = ref<boolean>(false);

// ---- accounts and categories come from the app's own stores, which already follow the selected business
const {
    paymentAccounts, receivableAccounts, selectedPaymentAccount, compatibleReceivables,
    paymentAccountOptions, receivableAccountOptions, incomeCategoryOptions: categoryOptions, currency, load: loadAccounts
} = useBusinessAccounts(paymentAccountId);

const locationOptions = computed(() => locations.value.map(l => ({ title: l.name, value: l.id })));
const customerOptions = computed(() => customers.value.map(c => ({ title: c.outstanding > 0 ? `${c.name} · ${tt('owes')} ${money(c.outstanding)}` : c.name, value: c.id })));

const visibleItems = computed<ItemInfo[]>(() => {
    const needle = (search.value ?? '').trim().toLowerCase();
    return items.value.filter(item => !needle || item.name.toLowerCase().includes(needle) || item.sku.toLowerCase().includes(needle));
});

function money(minorUnits: number): string {
    return formatAmountToLocalizedNumeralsWithCurrency(parseBigDecimal(minorUnits), currency.value);
}

function formatTime(unixTime: number): string {
    return formatDateTimeToLongDateTime(parseDateTimeFromUnixTime(unixTime));
}

function customerName(id: string): string {
    return id === '0' ? tt('Walk-in customer') : (customers.value.find(c => c.id === id)?.name ?? '…');
}

// stock at the chosen location, minus what is already in the cart
function stockAt(item: ItemInfo): number {
    return levels.value.filter(l => l.itemId === item.id && l.locationId === locationId.value).reduce((sum, l) => sum + l.qty, 0);
}

function availableOf(item: ItemInfo): number {
    return stockAt(item);
}

// ---- cart
const cart = computed<CartLine[]>(() => cartInput.value.flatMap(entry => {
    const item = items.value.find(i => i.id === entry.itemId);

    if (!item) {
        return [];
    }

    const qty = parseQty(entry.qtyText);
    const unitPrice = entry.priceOverride ?? item.salePrice;
    let qtyProblem = '';

    if (qty === null || qty <= 0) {
        qtyProblem = tt('Enter a quantity such as 2 or 0.5');
    } else if (item.trackStock && qty > stockAt(item)) {
        qtyProblem = tt('Only {available} available here', { available: formatQty(stockAt(item)) });
    }

    return [{
        itemId: item.id, item, qtyText: entry.qtyText, qty: qty !== null && qty > 0 ? qty : null,
        priceOverride: entry.priceOverride, unitPrice, qtyProblem,
        total: qty !== null && qty > 0 ? lineTotal(qty, unitPrice) : 0
    }];
}));

const totals = computed(() => computeTotals(cart.value.map(l => l.total), canManage.value ? discount.value : 0, payMode.value, partialPaid.value));

function addToCart(item: ItemInfo): void {
    const existing = cartInput.value.find(entry => entry.itemId === item.id);

    if (existing) {
        const qty = parseQty(existing.qtyText) ?? 0;
        existing.qtyText = formatQty(qty + 1000);
    } else {
        cartInput.value.push({ itemId: item.id, qtyText: '1', priceOverride: null });
    }
}

function removeLine(itemId: string): void {
    cartInput.value = cartInput.value.filter(entry => entry.itemId !== itemId);
}

// the cart lines shown on screen are computed from cartInput, so edits must go to cartInput
function setQty(itemId: string, text: string): void {
    const entry = cartInput.value.find(e => e.itemId === itemId);

    if (entry) {
        entry.qtyText = text;
    }
}

function setPrice(itemId: string, value: number): void {
    const entry = cartInput.value.find(e => e.itemId === itemId);
    const item = items.value.find(i => i.id === itemId);

    if (entry && item) {
        entry.priceOverride = value === item.salePrice ? null : value;
    }
}

function clearCart(): void {
    cartInput.value = [];
    discount.value = 0;
    partialPaid.value = 0;
    payMode.value = 'full';
    note.value = '';
    showProblems.value = false;
}

// ---- validation
const problems = reactive({ payment: '', customer: '', receivable: '', category: '' });

const problemList = computed(() => ({
    payment: totals.value.paid > 0 && !selectedPaymentAccount.value ? tt('Choose where the money goes') : '',
    customer: totals.value.credit > 0 && !customerId.value ? tt('A sale on credit needs a customer') : '',
    receivable: totals.value.credit > 0 && !compatibleReceivables.value.some(a => a.id === receivableAccountId.value)
        ? (compatibleReceivables.value.length < 1 ? tt('Create an account of type Receivables first') : tt('Choose where the owed amount is tracked'))
        : '',
    category: !categoryId.value ? tt('Choose an income category') : ''
}));

watch([problemList, showProblems], () => {
    problems.payment = showProblems.value ? problemList.value.payment : '';
    problems.customer = showProblems.value ? problemList.value.customer : '';
    problems.receivable = showProblems.value ? problemList.value.receivable : '';
    problems.category = showProblems.value ? problemList.value.category : '';
}, { deep: true });

const canSubmit = computed<boolean>(() => cart.value.length > 0
    && cart.value.every(l => l.qty !== null && !l.qtyProblem)
    && totals.value.total > 0);

// ---- remembered choices
function businessKey(): string {
    return current.value?.ownerUid ?? 'own';
}

function loadDefaults(): void {
    const saved = loadFormDefaults(businessKey());

    paymentAccountId.value = paymentAccounts.value.some(a => a.id === saved.payment) ? saved.payment! : (paymentAccounts.value[0]?.id ?? '');
    receivableAccountId.value = receivableAccounts.value.some(a => a.id === saved.receivable) ? saved.receivable! : (receivableAccounts.value[0]?.id ?? '');
    categoryId.value = categoryOptions.value.some(c => c.value === saved.category) ? saved.category! : (categoryOptions.value[0]?.value ?? '');
}

function saveDefaults(): void {
    saveFormDefaults(businessKey(), { payment: paymentAccountId.value, receivable: receivableAccountId.value, category: categoryId.value });
}

// keep the owed-money account valid when the payment account (and so the currency) changes
watch(compatibleReceivables, accounts => {
    if (!accounts.some(a => a.id === receivableAccountId.value)) {
        receivableAccountId.value = accounts[0]?.id ?? '';
    }
});

// ---- loading
async function loadStockAndSales(): Promise<void> {
    [levels.value, sales.value] = await Promise.all([api.getStockLevels(), api.listSales()]);
}

async function load(): Promise<void> {
    loading.value = true;

    try {
        await ensureLoaded();
        [items.value, locations.value, customers.value] = await Promise.all([api.listItems(), api.listLocations(), api.listCustomers()]);
        await Promise.all([loadStockAndSales(), loadAccounts(), loadPeople()]);

        locationId.value = (locations.value.find(l => l.isDefault) ?? locations.value[0])?.id ?? '';
        loadDefaults();
    } catch (error) {
        snackbar.value?.showError(error);
    } finally {
        loading.value = false;
    }
}

// ---- actions
async function submit(): Promise<void> {
    showProblems.value = true;

    if (!canSubmit.value || Object.values(problemList.value).some(Boolean)) {
        return;
    }

    submitting.value = true;

    try {
        const sale = await api.createSale({
            locationId: locationId.value || '0',
            customerId: customerId.value || '0',
            time: Math.floor(Date.now() / 1000),
            utcOffset: -new Date().getTimezoneOffset(),
            lines: cart.value.map(l => ({ itemId: l.itemId, qty: l.qty as number, unitPrice: l.priceOverride ?? undefined })),
            discount: totals.value.discount,
            amountPaid: totals.value.paid,
            paymentAccountId: totals.value.paid > 0 ? paymentAccountId.value : '0',
            receivableAccountId: totals.value.credit > 0 ? receivableAccountId.value : '0',
            categoryId: categoryId.value,
            note: note.value.trim()
        });

        saveDefaults();
        snackbar.value?.showMessage(tt('Sale #{id} recorded: {total}', { id: sale.id, total: money(sale.total) }));
        clearCart();
        await Promise.all([loadStockAndSales(), api.listCustomers().then(list => { customers.value = list; }), loadAccounts(true)]);
    } catch (error) {
        snackbar.value?.showError(error);
    } finally {
        submitting.value = false;
    }
}

function voidSale(sale: SaleInfo): void {
    confirmDialog.value?.open('ext.confirmVoidSale', { total: money(sale.total) }).then(async () => {
        try {
            await api.voidSale(sale.id);
            snackbar.value?.showMessage(tt('Sale voided'));
            await Promise.all([loadStockAndSales(), api.listCustomers().then(list => { customers.value = list; }), loadAccounts(true)]);
        } catch (error) {
            snackbar.value?.showError(error);
        }
    }).catch(() => {
        // cancelled
    });
}

function onCustomerSaved(customer: CustomerInfo): void {
    customers.value = [...customers.value, customer];
    customerId.value = customer.id;
}

function onError(error: unknown): void {
    snackbar.value?.showError(error);
}

onMounted(load);
</script>
