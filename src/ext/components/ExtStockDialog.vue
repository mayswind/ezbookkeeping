<template>
    <v-dialog width="480" persistent :model-value="show" @update:model-value="close">
        <v-card>
            <v-card-title>{{ title }}</v-card-title>
            <v-card-subtitle v-if="item">{{ item.name }} ({{ item.sku }})</v-card-subtitle>
            <v-card-text>
                <v-form @submit.prevent="save">
                    <v-select density="compact" :label="mode === 'transfer' ? tt('From') : tt('Location')"
                              :items="locationOptions" item-title="title" item-value="value"
                              :disabled="saving" v-model="locationId" v-if="locations.length > 1" />
                    <v-select density="compact" :label="tt('To')"
                              :items="destinationOptions" item-title="title" item-value="value"
                              :disabled="saving" v-model="toLocationId" v-if="mode === 'transfer'" />

                    <v-text-field density="compact" :label="quantityLabel" :disabled="saving"
                                  :error-messages="errors.qty" v-model="qtyText" />
                    <amount-input density="compact" :currency="currency" :label="tt('Cost per unit')"
                                  :disabled="saving" v-model="unitCost" v-if="mode === 'receive'" />
                    <v-switch density="compact" hide-details :label="tt('This is opening stock')"
                              :disabled="saving" v-model="opening" v-if="mode === 'receive'" />
                    <v-text-field density="compact" :label="tt('Note (optional)')" :disabled="saving" v-model="note" />

                    <div class="text-body-2 text-medium-emphasis mt-1" v-if="mode === 'adjust'">
                        {{ tt('Use a minus sign to remove stock, for example -2 for two damaged items.') }}
                    </div>
                </v-form>
            </v-card-text>
            <v-card-actions>
                <v-spacer />
                <v-btn variant="text" :disabled="saving" @click="close(false)">{{ tt('Cancel') }}</v-btn>
                <v-btn color="primary" :disabled="saving" @click="save">{{ tt('Save') }}</v-btn>
            </v-card-actions>
        </v-card>
    </v-dialog>
</template>

<script setup lang="ts">
import { ref, reactive, computed, watch } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import api from '@/ext/api.ts';
import { parseQty, parseSignedQty } from '@/ext/qty.ts';
import type { ItemInfo, LocationInfo } from '@/ext/types.ts';

export type StockDialogMode = 'receive' | 'adjust' | 'transfer';

const props = defineProps<{
    show: boolean;
    mode: StockDialogMode;
    item: ItemInfo | null;
    locations: LocationInfo[];
    currency: string;
}>();

const emit = defineEmits<{
    (e: 'update:show', value: boolean): void;
    (e: 'saved'): void;
    (e: 'error', error: unknown): void;
}>();

const { tt } = useI18n();

const saving = ref<boolean>(false);
const locationId = ref<string>('');
const toLocationId = ref<string>('');
const qtyText = ref<string>('');
const unitCost = ref<number>(0);
const opening = ref<boolean>(false);
const note = ref<string>('');
const errors = reactive<{ qty: string }>({ qty: '' });

const title = computed<string>(() => {
    switch (props.mode) {
        case 'receive': return tt('Receive stock');
        case 'adjust': return tt('Adjust stock');
        default: return tt('Transfer stock');
    }
});

const quantityLabel = computed<string>(() => props.mode === 'adjust' ? tt('Change in quantity') : tt('Quantity'));

const locationOptions = computed(() => props.locations.map(l => ({ title: l.name, value: l.id })));
const destinationOptions = computed(() => locationOptions.value.filter(o => o.value !== locationId.value));

watch(() => props.show, open => {
    if (!open) {
        return;
    }

    const defaultLocation = props.locations.find(l => l.isDefault) ?? props.locations[0];
    locationId.value = defaultLocation?.id ?? '';
    toLocationId.value = props.locations.find(l => l.id !== locationId.value)?.id ?? '';
    qtyText.value = '';
    unitCost.value = props.item?.costPrice ?? 0;
    opening.value = false;
    note.value = '';
    errors.qty = '';
});

function close(value: boolean = false): void {
    emit('update:show', value);
}

async function save(): Promise<void> {
    if (!props.item) {
        return;
    }

    const qty = props.mode === 'adjust' ? parseSignedQty(qtyText.value) : parseQty(qtyText.value);

    if (qty === null || qty === 0) {
        errors.qty = tt('Enter a number such as 5 or 2.5');
        return;
    }

    errors.qty = '';
    saving.value = true;

    try {
        if (props.mode === 'receive') {
            await api.receiveStock({ itemId: props.item.id, locationId: locationId.value || '0', qty, unitCost: unitCost.value, note: note.value, opening: opening.value });
        } else if (props.mode === 'adjust') {
            await api.adjustStock({ itemId: props.item.id, locationId: locationId.value || '0', qty, note: note.value });
        } else {
            await api.transferStock({ itemId: props.item.id, fromLocationId: locationId.value || '0', toLocationId: toLocationId.value, qty, note: note.value });
        }

        emit('saved');
        close(false);
    } catch (error) {
        emit('error', error);
    } finally {
        saving.value = false;
    }
}
</script>
