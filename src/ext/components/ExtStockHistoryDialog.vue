<template>
    <v-dialog width="820" scrollable :model-value="show" @update:model-value="close">
        <v-card v-if="item">
            <v-card-title>{{ tt('Stock history') }}</v-card-title>
            <v-card-subtitle>{{ item.name }} ({{ item.sku }})</v-card-subtitle>
            <v-card-text>
                <v-progress-linear indeterminate v-if="loading" />
                <v-table density="compact">
                    <thead>
                        <tr>
                            <th>{{ tt('When') }}</th>
                            <th class="text-end">{{ tt('Change') }}</th>
                            <th>{{ tt('Why') }}</th>
                            <th v-if="locations.length > 1">{{ tt('Location') }}</th>
                            <th>{{ tt('Recorded by') }}</th>
                            <th>{{ tt('Note') }}</th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr :key="movement.id" v-for="movement in movements">
                            <td class="text-no-wrap">{{ formatTime(movement.time) }}</td>
                            <td class="text-end text-no-wrap" :class="movement.qtyChange < 0 ? 'text-error' : 'text-success'">
                                {{ movement.qtyChange > 0 ? '+' : '' }}{{ formatQty(movement.qtyChange) }}
                            </td>
                            <td>{{ reasonLabel(movement.reason) }}</td>
                            <td v-if="locations.length > 1">{{ locationName(movement.locationId) }}</td>
                            <td>{{ nameOf(movement.actorUid) || '–' }}</td>
                            <td>{{ movement.note }}</td>
                        </tr>
                        <tr v-if="!loading && movements.length < 1">
                            <td :colspan="locations.length > 1 ? 6 : 5" class="text-medium-emphasis">{{ tt('No stock movements yet.') }}</td>
                        </tr>
                    </tbody>
                </v-table>
                <div class="text-caption text-medium-emphasis mt-2" v-if="movements.length >= 50">
                    {{ tt('Showing the latest 50 movements.') }}
                </div>
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

import { parseDateTimeFromUnixTime } from '@/lib/datetime.ts';

import api from '@/ext/api.ts';
import { usePeople } from '@/ext/people.ts';
import { formatQty } from '@/ext/qty.ts';
import type { ItemInfo, LocationInfo, StockMovement } from '@/ext/types.ts';

const props = defineProps<{
    show: boolean;
    item: ItemInfo | null;
    locations: LocationInfo[];
}>();

const emit = defineEmits<{
    (e: 'update:show', value: boolean): void;
    (e: 'error', error: unknown): void;
}>();

const { tt, formatDateTimeToLongDateTime } = useExtI18n();
const { load: loadPeople, nameOf } = usePeople();

const loading = ref<boolean>(false);
const movements = ref<StockMovement[]>([]);

// matches StockReason in pkg/ext/models
function reasonLabel(reason: number): string {
    switch (reason) {
        case 1: return tt('Opening stock');
        case 2: return tt('Purchase');
        case 3: return tt('Sale');
        case 4: return tt('Adjustment');
        case 5: return tt('Transferred out');
        case 6: return tt('Transferred in');
        case 7: return tt('Sale voided');
        default: return String(reason);
    }
}

function locationName(id: string): string {
    return props.locations.find(l => l.id === id)?.name ?? '';
}

function formatTime(unixTime: number): string {
    return formatDateTimeToLongDateTime(parseDateTimeFromUnixTime(unixTime));
}

function close(value: boolean = false): void {
    emit('update:show', value);
}

watch(() => props.show, async open => {
    if (!open || !props.item) {
        return;
    }

    loading.value = true;
    movements.value = [];

    try {
        [movements.value] = await Promise.all([api.listStockMovements(props.item.id), loadPeople()]);
    } catch (error) {
        emit('error', error);
    } finally {
        loading.value = false;
    }
});
</script>
