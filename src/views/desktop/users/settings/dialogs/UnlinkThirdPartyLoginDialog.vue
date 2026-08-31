<template>
    <v-dialog width="640" :persistent="true" v-model="showState">
        <one-column-dialog-layout :disabled="unlinking"
                                  :title="tt('Are you sure you want to unlink this login method?')"
                                  :cancel-button-title="tt('Cancel')"
                                  @cancel="cancel">
            <template #content>
                <div class="w-100">
                    <v-text-field
                        autocomplete="current-password"
                        type="password"
                        variant="underlined"
                        color="error"
                        :disabled="unlinking"
                        :placeholder="tt('Current Password')"
                        v-model="currentPassword"
                    />
                </div>
            </template>

            <template #footer>
                <v-btn color="secondary" variant="tonal" :disabled="unlinking" @click="cancel">{{ tt('Cancel') }}</v-btn>
                <v-spacer/>
                <v-btn color="error" :disabled="!currentPassword || unlinking" @click="confirm">
                    {{ tt('Confirm') }}
                    <v-progress-circular indeterminate size="22" class="ms-2" v-if="unlinking"></v-progress-circular>
                </v-btn>
            </template>
        </one-column-dialog-layout>
    </v-dialog>

    <snack-bar ref="snackbar" />
</template>

<script setup lang="ts">
import SnackBar from '@/components/desktop/SnackBar.vue';

import { ref, useTemplateRef } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import { useUserExternalAuthStore } from '@/stores/userExternalAuth.ts';

type SnackBarType = InstanceType<typeof SnackBar>;

const { tt } = useI18n();

const userExternalAuthStore = useUserExternalAuthStore();

let resolveFunc: (() => void) | null = null;
let rejectFunc: ((reason?: unknown) => void) | null = null;

const snackbar = useTemplateRef<SnackBarType>('snackbar');

const showState = ref<boolean>(false);
const unlinking = ref<boolean>(false);
const currentPassword = ref<string>('');
const currentExternalAuthType = ref<string | undefined>(undefined);

function open(externalAuthType: string): Promise<void> {
    showState.value = true;
    unlinking.value = false;
    currentPassword.value = '';
    currentExternalAuthType.value = externalAuthType;

    return new Promise((resolve, reject) => {
        resolveFunc = resolve;
        rejectFunc = reject;
    });
}

function confirm(): void {
    if (!currentExternalAuthType.value || !currentPassword.value) {
        return;
    }

    unlinking.value = true;

    userExternalAuthStore.unlinkExternalAuth({
        externalAuthType: currentExternalAuthType.value,
        password: currentPassword.value
    }).then(() => {
        unlinking.value = false;

        resolveFunc?.();
        showState.value = false;
    }).catch(error => {
        unlinking.value = false;

        if (!error.processed) {
            snackbar.value?.showError(error);
        }
    });
}

function cancel(): void {
    rejectFunc?.();
    showState.value = false;
}

defineExpose({
    open
});
</script>
