<template>
    <v-dialog width="440" v-model="showState">
        <v-card class="pa-2 pa-sm-4 pa-md-4">
            <template #title>
                <div class="d-flex align-center justify-center">
                    <h5 class="text-h5">{{ $t('Use on Mobile Device') }}</h5>
                </div>
            </template>
            <template #subtitle>
                <div class="text-body-1 text-center text-wrap mt-6">{{ $t('You can scan the below QR code on your mobile device.') }}</div>
            </template>
            <v-card-text class="mb-md-4">
                <v-row>
                    <v-col cols="12" md="12">
                        <div class="w-100 d-flex justify-center">
                            <img alt="qrcode" class="img-url-qrcode" :src="mobileUrlQrCodePath" />
                        </div>
                    </v-col>
                </v-row>
            </v-card-text>
            <v-card-text class="overflow-y-visible">
                <div class="w-100 d-flex justify-center gap-4">
                    <v-btn :href="mobileVersionPath">{{$t('Switch to Mobile Version') }}</v-btn>
                    <v-btn color="secondary" variant="tonal" @click="showState = false">{{ $t('Close') }}</v-btn>
                </div>
            </v-card-text>
        </v-card>
    </v-dialog>
</template>

<script>
import { getMobileUrlQrCodePath } from '@/lib/qrcode.js';
import { getMobileVersionPath } from '@/lib/version.js';

export default {
    props: [
        'show'
    ],
    emits: [
        'update:show'
    ],
    data() {
        return {
            mobileUrlQrCodePath: getMobileUrlQrCodePath(),
            mobileVersionPath: getMobileVersionPath(),
        }
    },
    computed: {
        showState: {
            get: function () {
                return this.show;
            },
            set: function (value) {
                this.$emit('update:show', value);
            }
        }
    }
}
</script>

<style>
.img-url-qrcode {
    width: 320px;
    height: 320px
}
</style>
