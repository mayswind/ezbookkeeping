<template>
    <main-page-layout>
        <template #content>
            <v-row class="match-height">
                <v-col cols="12">
                    <v-card>
                        <template #title>
                            <div class="d-flex align-center flex-wrap ga-3">
                                <span class="me-2">{{ tt('Inventory') }}</span>
                                <v-select style="max-width: 220px" density="compact" hide-details variant="outlined"
                                          :label="tt('Location')" :items="locationFilterOptions" item-title="title" item-value="value"
                                          :disabled="loading" v-model="locationFilter" v-if="locations.length > 1" />
                                <v-text-field style="max-width: 260px" density="compact" hide-details variant="outlined"
                                              :label="tt('Search by name or SKU')" clearable v-model="search" />
                                <v-spacer />
                                <template v-if="canManage">
                                    <v-btn variant="tonal" :disabled="loading" @click="showLocations = true">{{ tt('Locations') }}</v-btn>
                                    <v-btn color="primary" :disabled="loading" @click="openItemDialog(null)">{{ tt('Add item') }}</v-btn>
                                </template>
                            </div>
                        </template>

                        <v-card-text class="pt-0" v-if="!canManage">
                            <v-alert type="info" variant="tonal" density="compact">
                                {{ tt('You can look up items and stock. Ask a manager or the owner to change them.') }}
                            </v-alert>
                        </v-card-text>

                        <v-progress-linear indeterminate v-if="loading" />

                        <v-table>
                            <thead>
                                <tr>
                                    <th>{{ tt('SKU') }}</th>
                                    <th>{{ tt('Name') }}</th>
                                    <th class="text-end">{{ tt('In stock') }}</th>
                                    <th class="text-end">{{ tt('Sale price') }}</th>
                                    <th v-if="canManage"></th>
                                </tr>
                            </thead>
                            <tbody>
                                <tr :key="item.id" v-for="item in visibleItems">
                                    <td>{{ item.sku }}</td>
                                    <td>{{ item.name }}</td>
                                    <td class="text-end">
                                        <template v-if="item.trackStock">
                                            <v-chip size="x-small" color="warning" class="me-2" v-if="isLow(item)">{{ tt('Low') }}</v-chip>
                                            <span>{{ formatQty(stockOf(item)) }}</span>
                                            <span class="text-medium-emphasis" v-if="item.unit"> {{ item.unit }}</span>
                                            <div class="text-caption text-medium-emphasis" v-if="locationFilter === ALL && locations.length > 1">
                                                {{ breakdown(item) }}
                                            </div>
                                        </template>
                                        <span class="text-medium-emphasis" v-else>{{ tt('Not tracked') }}</span>
                                    </td>
                                    <td class="text-end">{{ money(item.salePrice) }}</td>
                                    <td class="text-end text-no-wrap" v-if="canManage">
                                        <template v-if="item.trackStock">
                                            <v-btn size="small" variant="text" color="primary" @click="openStockDialog('receive', item)">{{ tt('Receive') }}</v-btn>
                                            <v-btn size="small" variant="text" @click="openStockDialog('adjust', item)">{{ tt('Adjust') }}</v-btn>
                                            <v-btn size="small" variant="text" v-if="locations.length > 1"
                                                   @click="openStockDialog('transfer', item)">{{ tt('Transfer') }}</v-btn>
                                        </template>
                                        <v-btn size="small" variant="text" v-if="item.trackStock" @click="openHistory(item)">{{ tt('History') }}</v-btn>
                                        <v-btn size="small" variant="text" @click="openItemDialog(item)">{{ tt('Edit') }}</v-btn>
                                        <v-btn size="small" variant="text" color="error" @click="removeItem(item)">{{ tt('Delete') }}</v-btn>
                                    </td>
                                </tr>
                                <tr v-if="!loading && visibleItems.length < 1">
                                    <td :colspan="canManage ? 5 : 4" class="text-center text-medium-emphasis py-6">
                                        {{ items.length < 1 ? tt('No items yet.') : tt('No items match your search.') }}
                                    </td>
                                </tr>
                            </tbody>
                        </v-table>
                    </v-card>
                </v-col>
            </v-row>

            <ext-item-dialog :item="editingItem" :currency="currency" v-model:show="showItemDialog"
                             @saved="onSaved(tt('Item saved'))" @error="onError" />
            <ext-stock-dialog :mode="stockMode" :item="stockItem" :locations="locations" :currency="currency"
                              v-model:show="showStockDialog" @saved="onSaved(tt('Stock updated'))" @error="onError" />
            <ext-stock-history-dialog :item="historyItem" :locations="locations" v-model:show="showHistory" @error="onError" />
            <ext-locations-dialog :locations="locations" v-model:show="showLocations" @changed="load" @error="onError" />
            <confirm-dialog ref="confirmDialog" />
            <ext-snack-bar ref="snackbar" />
        </template>
    </main-page-layout>
</template>

<script setup lang="ts">
import ConfirmDialog from '@/components/desktop/ConfirmDialog.vue';
import ExtSnackBar from '@/ext/components/ExtSnackBar.vue';
import ExtItemDialog from '@/ext/components/ExtItemDialog.vue';
import ExtStockDialog, { type StockDialogMode } from '@/ext/components/ExtStockDialog.vue';
import ExtLocationsDialog from '@/ext/components/ExtLocationsDialog.vue';
import ExtStockHistoryDialog from '@/ext/components/ExtStockHistoryDialog.vue';

import { ref, computed, onMounted, useTemplateRef } from 'vue';

import { useExtI18n } from '@/ext/i18n.ts';
import { useUserStore } from '@/stores/user.ts';

import { parseBigDecimal } from '@/lib/numeral.ts';

import api from '@/ext/api.ts';
import { useBusiness } from '@/ext/business.ts';
import { formatQty } from '@/ext/qty.ts';
import type { ItemInfo, LocationInfo, StockLevel } from '@/ext/types.ts';

type ConfirmDialogType = InstanceType<typeof ConfirmDialog>;
type SnackBarType = InstanceType<typeof ExtSnackBar>;

const ALL = 'all';

const { tt, formatAmountToLocalizedNumeralsWithCurrency } = useExtI18n();
const userStore = useUserStore();
const { canManage, refresh } = useBusiness();

const confirmDialog = useTemplateRef<ConfirmDialogType>('confirmDialog');
const snackbar = useTemplateRef<SnackBarType>('snackbar');

const loading = ref<boolean>(true);
const items = ref<ItemInfo[]>([]);
const locations = ref<LocationInfo[]>([]);
const levels = ref<StockLevel[]>([]);
const search = ref<string | null>('');
const locationFilter = ref<string>(ALL);

const showItemDialog = ref<boolean>(false);
const editingItem = ref<ItemInfo | null>(null);
const showStockDialog = ref<boolean>(false);
const stockMode = ref<StockDialogMode>('receive');
const stockItem = ref<ItemInfo | null>(null);
const showLocations = ref<boolean>(false);
const showHistory = ref<boolean>(false);
const historyItem = ref<ItemInfo | null>(null);

// Prices are in the business currency; until the server tells us otherwise we show them in the user's default currency
const currency = computed<string>(() => userStore.currentUserDefaultCurrency);

const locationFilterOptions = computed(() => [
    { title: tt('All locations'), value: ALL },
    ...locations.value.map(l => ({ title: l.name, value: l.id }))
]);

const visibleItems = computed<ItemInfo[]>(() => {
    const needle = (search.value ?? '').trim().toLowerCase();

    return items.value.filter(item => !needle || item.name.toLowerCase().includes(needle) || item.sku.toLowerCase().includes(needle));
});

function money(minorUnits: number): string {
    return formatAmountToLocalizedNumeralsWithCurrency(parseBigDecimal(minorUnits), currency.value);
}

function stockOf(item: ItemInfo): number {
    return levels.value
        .filter(l => l.itemId === item.id && (locationFilter.value === ALL || l.locationId === locationFilter.value))
        .reduce((sum, l) => sum + l.qty, 0);
}

function isLow(item: ItemInfo): boolean {
    return item.reorderLevel > 0 && stockOf(item) <= item.reorderLevel;
}

function breakdown(item: ItemInfo): string {
    return locations.value
        .map(location => ({ name: location.name, qty: levels.value.filter(l => l.itemId === item.id && l.locationId === location.id).reduce((sum, l) => sum + l.qty, 0) }))
        .filter(entry => entry.qty !== 0)
        .map(entry => `${entry.name}: ${formatQty(entry.qty)}`)
        .join(' · ');
}

async function load(): Promise<void> {
    loading.value = true;

    try {
        await refresh();
        [items.value, locations.value, levels.value] = await Promise.all([api.listItems(), api.listLocations(), api.getStockLevels()]);

        if (locationFilter.value !== ALL && !locations.value.some(l => l.id === locationFilter.value)) {
            locationFilter.value = ALL;
        }
    } catch (error) {
        snackbar.value?.showError(error);
    } finally {
        loading.value = false;
    }
}

function openItemDialog(item: ItemInfo | null): void {
    editingItem.value = item;
    showItemDialog.value = true;
}

function openStockDialog(mode: StockDialogMode, item: ItemInfo): void {
    stockMode.value = mode;
    stockItem.value = item;
    showStockDialog.value = true;
}

function openHistory(item: ItemInfo): void {
    historyItem.value = item;
    showHistory.value = true;
}

function removeItem(item: ItemInfo): void {
    confirmDialog.value?.open('ext.confirmDeleteItem', { name: item.name }).then(async () => {
        try {
            await api.deleteItem(item.id);
            snackbar.value?.showMessage(tt('Item deleted'));
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
