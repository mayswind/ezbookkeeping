<template>
    <main-page-layout>
        <template #content>
            <v-row class="match-height">
                <v-col cols="12" v-if="!canManage">
                    <v-alert type="info" variant="tonal">{{ tt('Reports are for managers and the owner.') }}</v-alert>
                </v-col>

                <v-col cols="12" v-else>
                    <v-card>
                        <template #title>
                            <div class="d-flex align-center flex-wrap ga-3">
                                <span class="me-2">{{ tt('Reports') }}</span>
                                <v-select style="max-width: 220px" density="compact" hide-details variant="outlined"
                                          :label="tt('Location')" :items="locationFilterOptions" item-title="title" item-value="value"
                                          :disabled="loading" v-model="locationFilter"
                                          v-if="locations.length > 1 && tab !== 'receivables'" />
                                <v-spacer />
                                <v-btn variant="tonal" :disabled="loading" @click="exportCsv">{{ tt('Download CSV') }}</v-btn>
                                <v-btn variant="text" :disabled="loading" @click="loadTab">{{ tt('Refresh') }}</v-btn>
                            </div>
                        </template>
                        <v-tabs color="primary" v-model="tab">
                            <v-tab value="stock">{{ tt('Stock value') }}</v-tab>
                            <v-tab value="low">
                                {{ tt('Low stock') }}
                                <v-chip class="ms-2" size="x-small" color="warning" v-if="lowStock.length > 0">{{ lowStock.length }}</v-chip>
                            </v-tab>
                            <v-tab value="receivables">{{ tt('Who owes what') }}</v-tab>
                        </v-tabs>
                        <v-progress-linear indeterminate v-if="loading" />

                        <!-- stock value -->
                        <v-card-text v-if="tab === 'stock' && stockValue">
                            <v-row>
                                <v-col cols="12" sm="4">
                                    <div class="text-caption text-medium-emphasis">{{ tt('Stock at cost') }}</div>
                                    <div class="text-h5">{{ money(stockValue.totalCostValue) }}</div>
                                </v-col>
                                <v-col cols="12" sm="4">
                                    <div class="text-caption text-medium-emphasis">{{ tt('Stock at selling price') }}</div>
                                    <div class="text-h5">{{ money(stockValue.totalRetailValue) }}</div>
                                </v-col>
                                <v-col cols="12" sm="4">
                                    <div class="text-caption text-medium-emphasis">{{ tt('Profit if all of it sells') }}</div>
                                    <div class="text-h5">{{ money(stockValue.totalRetailValue - stockValue.totalCostValue) }}</div>
                                </v-col>
                            </v-row>
                            <div class="text-body-2 text-medium-emphasis mt-2">
                                {{ tt('Values use each item\'s current cost price and sale price from the Inventory page.') }}
                            </div>
                        </v-card-text>
                        <v-table v-if="tab === 'stock' && stockValue">
                            <thead>
                                <tr>
                                    <th>{{ tt('Item') }}</th>
                                    <th class="text-end">{{ tt('In stock') }}</th>
                                    <th class="text-end">{{ tt('At cost') }}</th>
                                    <th class="text-end">{{ tt('At selling price') }}</th>
                                </tr>
                            </thead>
                            <tbody>
                                <tr :key="row.item.id" v-for="row in stockValue.rows">
                                    <td>
                                        {{ row.item.name }}
                                        <div class="text-caption text-medium-emphasis">{{ row.item.sku }}</div>
                                    </td>
                                    <td class="text-end">
                                        {{ formatQty(row.qty) }}<span class="text-medium-emphasis" v-if="row.item.unit"> {{ row.item.unit }}</span>
                                        <div class="text-caption text-medium-emphasis" v-if="row.locations && row.locations.length > 1">
                                            {{ locationBreakdown(row.locations) }}
                                        </div>
                                    </td>
                                    <td class="text-end">{{ money(row.costValue) }}</td>
                                    <td class="text-end">{{ money(row.retailValue) }}</td>
                                </tr>
                                <tr v-if="!loading && stockValue.rows.length < 1">
                                    <td colspan="4" class="text-center text-medium-emphasis py-6">{{ tt('There is no stock on hand.') }}</td>
                                </tr>
                            </tbody>
                        </v-table>

                        <!-- low stock -->
                        <template v-if="tab === 'low'">
                            <v-card-text v-if="!loading && lowStock.length < 1">
                                <v-alert type="success" variant="tonal" density="compact">
                                    {{ tt('Nothing is running low. Set a reorder level on an item to be warned when it does.') }}
                                </v-alert>
                            </v-card-text>
                            <v-table v-else>
                                <thead>
                                    <tr>
                                        <th>{{ tt('Item') }}</th>
                                        <th class="text-end">{{ tt('In stock') }}</th>
                                        <th class="text-end">{{ tt('Reorder level') }}</th>
                                        <th class="text-end">{{ tt('Short by') }}</th>
                                    </tr>
                                </thead>
                                <tbody>
                                    <tr :key="row.item.id" v-for="row in lowStock">
                                        <td>
                                            {{ row.item.name }}
                                            <div class="text-caption text-medium-emphasis">{{ row.item.sku }}</div>
                                        </td>
                                        <td class="text-end">
                                            <v-chip size="x-small" color="error" class="me-2" v-if="row.qty <= 0">{{ tt('Out of stock') }}</v-chip>
                                            {{ formatQty(row.qty) }}<span class="text-medium-emphasis" v-if="row.item.unit"> {{ row.item.unit }}</span>
                                        </td>
                                        <td class="text-end">{{ formatQty(row.item.reorderLevel) }}</td>
                                        <td class="text-end font-weight-medium">{{ row.shortfall > 0 ? formatQty(row.shortfall) : '–' }}</td>
                                    </tr>
                                </tbody>
                            </v-table>
                        </template>

                        <!-- who owes what -->
                        <template v-if="tab === 'receivables' && receivables">
                            <v-card-text>
                                <v-row>
                                    <v-col cols="12" sm="4" lg="2">
                                        <div class="text-caption text-medium-emphasis">{{ tt('Total owed to you') }}</div>
                                        <div class="text-h5">{{ money(receivables.totalOutstanding) }}</div>
                                    </v-col>
                                    <v-col cols="6" sm="2" lg="2" :key="bucket.label" v-for="bucket in buckets">
                                        <div class="text-caption text-medium-emphasis">{{ bucket.label }}</div>
                                        <div class="text-subtitle-1" :class="bucket.warn && bucket.value > 0 ? 'text-error' : ''">{{ money(bucket.value) }}</div>
                                    </v-col>
                                </v-row>
                                <div class="text-body-2 text-medium-emphasis mt-2">{{ tt('Debts are aged from the date of each sale.') }}</div>
                            </v-card-text>
                            <v-table>
                                <thead>
                                    <tr>
                                        <th>{{ tt('Customer') }}</th>
                                        <th class="text-end">{{ tt('Owes') }}</th>
                                        <th class="text-end">{{ tt('Up to 30 days') }}</th>
                                        <th class="text-end">{{ tt('31 to 60 days') }}</th>
                                        <th class="text-end">{{ tt('61 to 90 days') }}</th>
                                        <th class="text-end">{{ tt('Over 90 days') }}</th>
                                        <th>{{ tt('Oldest unpaid sale') }}</th>
                                    </tr>
                                </thead>
                                <tbody>
                                    <tr :key="row.customer.id" v-for="row in receivables.rows">
                                        <td>
                                            {{ row.customer.name }}
                                            <div class="text-caption text-medium-emphasis">{{ tt('{count} unpaid sales', { count: row.openSales }) }}</div>
                                        </td>
                                        <td class="text-end font-weight-medium">{{ money(row.outstanding) }}</td>
                                        <td class="text-end">{{ cellMoney(row.current) }}</td>
                                        <td class="text-end">{{ cellMoney(row.days31To60) }}</td>
                                        <td class="text-end">{{ cellMoney(row.days61To90) }}</td>
                                        <td class="text-end" :class="{ 'text-error': row.over90 > 0 }">{{ cellMoney(row.over90) }}</td>
                                        <td>{{ formatDate(row.oldestSaleTime) }}</td>
                                    </tr>
                                    <tr v-if="!loading && receivables.rows.length < 1">
                                        <td colspan="7" class="text-center text-medium-emphasis py-6">{{ tt('Nobody owes you anything right now.') }}</td>
                                    </tr>
                                </tbody>
                            </v-table>
                        </template>
                    </v-card>
                </v-col>
            </v-row>

            <ext-snack-bar ref="snackbar" />
        </template>
    </main-page-layout>
</template>

<script setup lang="ts">
import ExtSnackBar from '@/ext/components/ExtSnackBar.vue';

import { ref, computed, watch, onMounted, useTemplateRef } from 'vue';

import { useExtI18n } from '@/ext/i18n.ts';
import { useUserStore } from '@/stores/user.ts';

import { parseBigDecimal } from '@/lib/numeral.ts';
import { parseDateTimeFromUnixTime } from '@/lib/datetime.ts';

import api from '@/ext/api.ts';
import { useBusiness } from '@/ext/business.ts';
import { downloadCsv } from '@/ext/csv.ts';
import { formatQty } from '@/ext/qty.ts';
import type { LocationInfo, LowStockRow, ReceivablesReport, StockValueReport } from '@/ext/types.ts';

type SnackBarType = InstanceType<typeof ExtSnackBar>;
type ReportTab = 'stock' | 'low' | 'receivables';

const ALL = 'all';

const { tt, formatAmountToLocalizedNumeralsWithCurrency, formatDateTimeToLongDate } = useExtI18n();
const userStore = useUserStore();
const { canManage, ensureLoaded } = useBusiness();

const snackbar = useTemplateRef<SnackBarType>('snackbar');

const loading = ref<boolean>(true);
const tab = ref<ReportTab>('stock');
const locations = ref<LocationInfo[]>([]);
const locationFilter = ref<string>(ALL);
const stockValue = ref<StockValueReport | null>(null);
const lowStock = ref<LowStockRow[]>([]);
const receivables = ref<ReceivablesReport | null>(null);

// Amounts are in the business currency; until the server says which, the user's default currency is shown
const currency = computed<string>(() => userStore.currentUserDefaultCurrency);

const locationFilterOptions = computed(() => [
    { title: tt('All locations'), value: ALL },
    ...locations.value.map(l => ({ title: l.name, value: l.id }))
]);

const buckets = computed(() => receivables.value ? [
    { label: tt('Up to 30 days'), value: receivables.value.current, warn: false },
    { label: tt('31 to 60 days'), value: receivables.value.days31To60, warn: false },
    { label: tt('61 to 90 days'), value: receivables.value.days61To90, warn: true },
    { label: tt('Over 90 days'), value: receivables.value.over90, warn: true }
] : []);

function money(minorUnits: number): string {
    return formatAmountToLocalizedNumeralsWithCurrency(parseBigDecimal(minorUnits), currency.value);
}

function cellMoney(minorUnits: number): string {
    return minorUnits > 0 ? money(minorUnits) : '–';
}

function formatDate(unixTime: number): string {
    return formatDateTimeToLongDate(parseDateTimeFromUnixTime(unixTime));
}

function locationName(id: string): string {
    return locations.value.find(l => l.id === id)?.name ?? '';
}

function locationBreakdown(parts: { locationId: string, qty: number }[]): string {
    return parts.map(part => `${locationName(part.locationId)}: ${formatQty(part.qty)}`).join(' · ');
}

const locationParam = computed<string | undefined>(() => locationFilter.value === ALL ? undefined : locationFilter.value);

async function loadTab(): Promise<void> {
    if (!canManage.value) {
        return;
    }

    loading.value = true;

    try {
        if (tab.value === 'stock') {
            stockValue.value = await api.getStockValueReport(locationParam.value);
        } else if (tab.value === 'low') {
            lowStock.value = await api.getLowStockReport(locationParam.value);
        } else {
            receivables.value = await api.getReceivablesReport();
        }
    } catch (error) {
        snackbar.value?.showError(error);
    } finally {
        loading.value = false;
    }
}

function exportCsv(): void {
    if (tab.value === 'stock' && stockValue.value) {
        downloadCsv('stock-value.csv', [
            [tt('SKU'), tt('Item'), tt('In stock'), tt('At cost'), tt('At selling price')],
            ...stockValue.value.rows.map(r => [r.item.sku, r.item.name, formatQty(r.qty), money(r.costValue), money(r.retailValue)]),
            ['', tt('Total'), '', money(stockValue.value.totalCostValue), money(stockValue.value.totalRetailValue)]
        ]);
    } else if (tab.value === 'low') {
        downloadCsv('low-stock.csv', [
            [tt('SKU'), tt('Item'), tt('In stock'), tt('Reorder level'), tt('Short by')],
            ...lowStock.value.map(r => [r.item.sku, r.item.name, formatQty(r.qty), formatQty(r.item.reorderLevel), formatQty(r.shortfall)])
        ]);
    } else if (tab.value === 'receivables' && receivables.value) {
        downloadCsv('who-owes-what.csv', [
            [tt('Customer'), tt('Owes'), tt('Up to 30 days'), tt('31 to 60 days'), tt('61 to 90 days'), tt('Over 90 days'), tt('Oldest unpaid sale')],
            ...receivables.value.rows.map(r => [r.customer.name, money(r.outstanding), money(r.current), money(r.days31To60), money(r.days61To90), money(r.over90), formatDate(r.oldestSaleTime)]),
            [tt('Total'), money(receivables.value.totalOutstanding), money(receivables.value.current), money(receivables.value.days31To60), money(receivables.value.days61To90), money(receivables.value.over90), '']
        ]);
    }
}

watch([tab, locationFilter], loadTab);

onMounted(async () => {
    try {
        await ensureLoaded();
        locations.value = canManage.value ? await api.listLocations() : [];
        // the low-stock count in the tab title needs its data even before that tab is opened
        lowStock.value = canManage.value ? await api.getLowStockReport() : [];
    } catch (error) {
        snackbar.value?.showError(error);
    }

    await loadTab();
});
</script>
