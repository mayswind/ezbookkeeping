<template>
    <v-dialog width="480" persistent :model-value="show" @update:model-value="close">
        <v-card>
            <v-card-title>{{ tt('Locations') }}</v-card-title>
            <v-card-subtitle>{{ tt('Stores or warehouses where you keep stock.') }}</v-card-subtitle>
            <v-card-text>
                <v-list density="compact">
                    <v-list-item :key="location.id" v-for="location in locations">
                        <template v-if="editingId === location.id">
                            <v-text-field density="compact" hide-details autofocus :disabled="busy"
                                          v-model="editingName" @keyup.enter="rename(location)" />
                        </template>
                        <v-list-item-title v-else>
                            {{ location.name }}
                            <v-chip class="ms-2" size="x-small" v-if="location.isDefault">{{ tt('Default') }}</v-chip>
                        </v-list-item-title>
                        <template #append>
                            <template v-if="editingId === location.id">
                                <v-btn size="small" variant="text" color="primary" :disabled="busy" @click="rename(location)">{{ tt('Save') }}</v-btn>
                                <v-btn size="small" variant="text" :disabled="busy" @click="editingId = ''">{{ tt('Cancel') }}</v-btn>
                            </template>
                            <template v-else>
                                <v-btn size="small" variant="text" :disabled="busy" @click="startEdit(location)">{{ tt('Rename') }}</v-btn>
                                <v-btn size="small" variant="text" color="error" :disabled="busy || locations.length < 2"
                                       @click="remove(location)">{{ tt('Delete') }}</v-btn>
                            </template>
                        </template>
                    </v-list-item>
                </v-list>

                <v-form class="d-flex align-center ga-2 mt-2" @submit.prevent="add">
                    <v-text-field density="compact" hide-details :label="tt('New location name')" :disabled="busy" v-model="newName" />
                    <v-btn color="primary" type="submit" :disabled="busy || !newName.trim()">{{ tt('Add') }}</v-btn>
                </v-form>
                <div class="text-body-2 text-medium-emphasis mt-2">
                    {{ tt('A location can only be deleted when it holds no stock. Move stock out with a transfer first.') }}
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

import { useI18n } from '@/locales/helpers.ts';

import api from '@/ext/api.ts';
import type { LocationInfo } from '@/ext/types.ts';

const props = defineProps<{
    show: boolean;
    locations: LocationInfo[];
}>();

const emit = defineEmits<{
    (e: 'update:show', value: boolean): void;
    (e: 'changed'): void;
    (e: 'error', error: unknown): void;
}>();

const { tt } = useI18n();

const busy = ref<boolean>(false);
const newName = ref<string>('');
const editingId = ref<string>('');
const editingName = ref<string>('');

watch(() => props.show, open => {
    if (open) {
        newName.value = '';
        editingId.value = '';
    }
});

function close(value: boolean = false): void {
    emit('update:show', value);
}

function startEdit(location: LocationInfo): void {
    editingId.value = location.id;
    editingName.value = location.name;
}

async function run(action: () => Promise<unknown>): Promise<void> {
    busy.value = true;

    try {
        await action();
        emit('changed');
    } catch (error) {
        emit('error', error);
    } finally {
        busy.value = false;
    }
}

async function add(): Promise<void> {
    const name = newName.value.trim();

    if (name) {
        await run(async () => {
            await api.addLocation(name);
            newName.value = '';
        });
    }
}

async function rename(location: LocationInfo): Promise<void> {
    const name = editingName.value.trim();

    if (name) {
        await run(async () => {
            await api.renameLocation(location.id, name);
            editingId.value = '';
        });
    }
}

async function remove(location: LocationInfo): Promise<void> {
    await run(() => api.deleteLocation(location.id));
}
</script>
