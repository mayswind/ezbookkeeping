<template>
    <v-dialog width="440" persistent :model-value="show" @update:model-value="close">
        <v-card>
            <v-card-title>{{ tt('Add customer') }}</v-card-title>
            <v-card-text>
                <v-form @submit.prevent="save">
                    <v-text-field density="compact" autofocus :label="tt('Name')" :disabled="saving"
                                  :error-messages="nameError" v-model="name" />
                    <v-text-field density="compact" :label="tt('Phone (optional)')" :disabled="saving" v-model="phone" />
                    <v-text-field density="compact" type="email" :label="tt('Email (optional)')" :disabled="saving" v-model="email" />
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
import { ref, watch } from 'vue';

import { useExtI18n } from '@/ext/i18n.ts';

import api from '@/ext/api.ts';
import type { CustomerInfo } from '@/ext/types.ts';

const props = defineProps<{
    show: boolean;
}>();

const emit = defineEmits<{
    (e: 'update:show', value: boolean): void;
    (e: 'saved', customer: CustomerInfo): void;
    (e: 'error', error: unknown): void;
}>();

const { tt } = useExtI18n();

const saving = ref<boolean>(false);
const name = ref<string>('');
const phone = ref<string>('');
const email = ref<string>('');
const nameError = ref<string>('');

watch(() => props.show, open => {
    if (open) {
        name.value = phone.value = email.value = nameError.value = '';
    }
});

function close(value: boolean = false): void {
    emit('update:show', value);
}

async function save(): Promise<void> {
    if (!name.value.trim()) {
        nameError.value = tt('Name is required');
        return;
    }

    nameError.value = '';
    saving.value = true;

    try {
        const customer = await api.addCustomer({ name: name.value.trim(), phone: phone.value.trim(), email: email.value.trim() });
        emit('saved', customer);
        close(false);
    } catch (error) {
        emit('error', error);
    } finally {
        saving.value = false;
    }
}
</script>
