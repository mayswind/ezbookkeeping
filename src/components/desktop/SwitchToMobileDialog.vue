<template>
    <v-dialog width="400" v-model="showState">
        <v-card>
            <v-toolbar color="primary">
                <v-toolbar-title>{{ $t('global.app.title') }}</v-toolbar-title>
            </v-toolbar>
            <v-card-text class="pa-4">
                <p>{{ $t('You can scan the below QR code on your mobile device.') }}</p>
            </v-card-text>
            <v-card-text class="pa-4 w-100 d-flex justify-center">
                <img alt="qrcode" class="img-url-qrcode" :src="mobileUrlQrCodePath" />
            </v-card-text>
            <v-card-actions>
                <v-btn :href="mobileVersionPath">{{$t('Switch to Mobile Version') }}</v-btn>
                <v-spacer></v-spacer>
                <v-btn @click="showState = false">{{ $t('Close') }}</v-btn>
            </v-card-actions>
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
