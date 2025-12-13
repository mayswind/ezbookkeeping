<template>
    <v-dialog width="640" :persistent="true" v-model="showState">
        <v-card class="pa-sm-1 pa-md-2">
            <template #title>
                <h4 class="text-h4 text-error text-wrap">{{ tt('Are you sure you want to unlink this login method?') }}</h4>
            </template>
            <v-card-text class="w-100 d-flex justify-center">
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
            </v-card-text>
            <v-card-text>
                <div class="w-100 d-flex justify-center flex-wrap mt-sm-1 mt-md-2 gap-4">
                    <v-btn color="error" :disabled="!currentPassword || unlinking" @click="confirm">
                        {{ tt('Confirm') }}
                        <v-progress-circular indeterminate size="22" class="ms-2" v-if="unlinking"></v-progress-circular>
                    </v-btn>
                    <v-btn color="secondary" variant="tonal" :disabled="unlinking" @click="cancel">
                        {{ tt('Cancel') }}
                    </v-btn>
                </div>
            </v-card-text>
        </v-card>
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
