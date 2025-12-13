<template>
    <v-dialog width="420" v-model="showState">
        <v-card class="pa-sm-1 pa-md-2">
            <template #title>
                <h4 class="text-h4">{{ tt('Use on Mobile Device') }}</h4>
            </template>
            <template #subtitle>
                <div class="text-body-1 text-wrap mt-4">{{ tt('You can scan the QR code below on your mobile device.') }}</div>
            </template>
            <v-card-text>
                <v-row>
                    <v-col cols="12" md="12">
                        <div class="w-100 d-flex justify-center">
                            <img alt="qrcode" class="img-url-qrcode" :src="mobileUrlQrCodePath" />
                        </div>
                    </v-col>
                </v-row>
            </v-card-text>
            <v-card-text>
                <div class="w-100 d-flex justify-center flex-wrap mt-sm-1 mt-md-2 gap-4">
                    <v-btn :href="mobileVersionPath">{{ tt('Switch to Mobile Version') }}</v-btn>
                    <v-btn color="secondary" variant="tonal" @click="showState = false">{{ tt('Close') }}</v-btn>
                </div>
            </v-card-text>
        </v-card>
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
