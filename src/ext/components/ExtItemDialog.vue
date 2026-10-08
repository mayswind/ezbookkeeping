<template>
    <v-dialog width="560" persistent :model-value="show" @update:model-value="close">
        <v-card>
            <v-card-title>{{ item ? tt('Edit item') : tt('Add item') }}</v-card-title>
            <v-card-text>
                <v-form @submit.prevent="save">
                    <v-row>
                        <v-col cols="12" sm="5">
                            <v-text-field density="compact" :label="tt('SKU')" :disabled="saving"
                                          :error-messages="errors.sku" v-model="sku" />
                        </v-col>
                        <v-col cols="12" sm="7">
                            <v-text-field density="compact" :label="tt('Name')" :disabled="saving"
                                          :error-messages="errors.name" v-model="name" />
                        </v-col>
                        <v-col cols="12" sm="4">
                            <v-text-field density="compact" :label="tt('Unit (pcs, kg, ...)')" :disabled="saving" v-model="unit" />
                        </v-col>
                        <v-col cols="12" sm="4">
                            <amount-input density="compact" :currency="currency" :label="tt('Cost price')"
                                          :disabled="saving" v-model="costPrice" />
                        </v-col>
                        <v-col cols="12" sm="4">
                            <amount-input density="compact" :currency="currency" :label="tt('Sale price')"
                                          :disabled="saving" v-model="salePrice" />
                        </v-col>
                        <v-col cols="12" sm="6">
                            <v-switch density="compact" hide-details :label="tt('Track stock')"
                                      :disabled="saving" v-model="trackStock" />
                        </v-col>
                        <v-col cols="12" sm="6" v-if="trackStock">
                            <v-text-field density="compact" :label="tt('Reorder level')" :disabled="saving"
                                          :error-messages="errors.reorderLevel" v-model="reorderLevel" />
                        </v-col>
                    </v-row>
                    <div class="text-body-2 text-medium-emphasis">
                        {{ tt('Prices are per whole unit. Turn off stock tracking for services.') }}
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
import { ref, reactive, watch } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import api from '@/ext/api.ts';
import { formatQty, parseQty } from '@/ext/qty.ts';
import type { ItemInfo } from '@/ext/types.ts';

const props = defineProps<{
    show: boolean;
    item: ItemInfo | null;
    currency: string;
}>();

const emit = defineEmits<{
    (e: 'update:show', value: boolean): void;
    (e: 'saved'): void;
    (e: 'error', error: unknown): void;
}>();

const { tt } = useI18n();

const saving = ref<boolean>(false);
const sku = ref<string>('');
const name = ref<string>('');
const unit = ref<string>('');
const costPrice = ref<number>(0);
const salePrice = ref<number>(0);
const trackStock = ref<boolean>(true);
const reorderLevel = ref<string>('0');
const errors = reactive<{ sku: string, name: string, reorderLevel: string }>({ sku: '', name: '', reorderLevel: '' });

// start every opening from the item being edited, or from an empty form
watch(() => props.show, open => {
    if (!open) {
        return;
    }

    sku.value = props.item?.sku ?? '';
    name.value = props.item?.name ?? '';
    unit.value = props.item?.unit ?? '';
    costPrice.value = props.item?.costPrice ?? 0;
    salePrice.value = props.item?.salePrice ?? 0;
    trackStock.value = props.item?.trackStock ?? true;
    reorderLevel.value = formatQty(props.item?.reorderLevel ?? 0);
    errors.sku = errors.name = errors.reorderLevel = '';
});

function close(value: boolean = false): void {
    emit('update:show', value);
}

async function save(): Promise<void> {
    errors.sku = sku.value.trim() ? '' : tt('SKU is required');
    errors.name = name.value.trim() ? '' : tt('Name is required');

    const reorder = parseQty(reorderLevel.value);
    errors.reorderLevel = reorder === null ? tt('Enter a number such as 5 or 2.5') : '';

    if (errors.sku || errors.name || errors.reorderLevel || reorder === null) {
        return;
    }

    saving.value = true;

    try {
        const request = {
            sku: sku.value.trim(),
            name: name.value.trim(),
            unit: unit.value.trim(),
            costPrice: costPrice.value,
            salePrice: salePrice.value,
            reorderLevel: trackStock.value ? reorder : 0,
            trackStock: trackStock.value
        };

        if (props.item) {
            await api.modifyItem({ ...request, id: props.item.id });
        } else {
            await api.addItem(request);
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
