<template>
    <v-dialog width="420" v-model="showState">
        <one-column-dialog-layout :title="tt('Use on Mobile Device')" :cancel-button-title="tt('Close')"
                                  @cancel="showState = false">
            <template #content>
                <div class="text-body-large text-wrap">{{ tt('You can scan the QR code below on your mobile device.') }}</div>
                <div class="d-flex justify-center w-100 mt-4">
                    <v-img alt="qrcode" class="img-url-qrcode" :src="mobileUrlQrCodePath">
                        <template #placeholder>
                            <div class="d-flex align-center justify-center">
                                <v-progress-circular color="grey-500" indeterminate size="48"></v-progress-circular>
                            </div>
                        </template>
                        <template #error>
                            <div class="d-flex align-center justify-center">
                                <span class="text-body-large">{{ tt('Failed to load QR code') }}</span>
                            </div>
                        </template>
                    </v-img>
                </div>
            </template>

            <template #footer>
                <v-spacer />
                <v-btn :href="mobileVersionPath">{{ tt('Switch to Mobile Version') }}</v-btn>
            </template>
        </one-column-dialog-layout>
    </v-dialog>
</template>

<script setup lang="ts">
import { computed } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import { getMobileUrlQrCodePath } from '@/lib/qrcode.ts';
import { getMobileVersionPath } from '@/lib/version.ts';

const props = defineProps<{
    show: boolean;
}>();

const emit = defineEmits<{
    (e: 'update:show', value: boolean): void;
}>();

const { tt } = useI18n();

const mobileUrlQrCodePath = getMobileUrlQrCodePath();
const mobileVersionPath = getMobileVersionPath();

const showState = computed<boolean>({
    get: () => {
        return props.show;
    },
    set: value => {
        emit('update:show', value);
    }
});
</script>

<style>
.img-url-qrcode {
    width: 320px;
    height: 320px
}
</style>
