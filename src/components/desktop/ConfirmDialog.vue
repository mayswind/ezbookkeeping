<template>
    <v-dialog persistent min-width="320" width="auto" v-model="showState">
        <v-card>
            <v-toolbar :color="finalColor">
                <v-toolbar-title>{{ titleContent }}</v-toolbar-title>
            </v-toolbar>
            <v-card-text v-if="textContent" class="pa-4">{{ textContent }}</v-card-text>
            <v-card-actions>
                <v-spacer></v-spacer>
                <v-btn color="gray" @click="cancel">{{ $t('Cancel') }}</v-btn>
                <v-btn :color="finalColor" @click="confirm">{{ $t('OK') }}</v-btn>
            </v-card-actions>
        </v-card>
    </v-dialog>
</template>

<script>
import { isString } from '@/lib/common.js';

export default {
    props: [
        'show',
        'color',
        'title',
        'text'
    ],
    emits: [
        'update:show'
    ],
    expose: [
        'open'
    ],
    data() {
        const self = this;

        return {
            showState: self.show,
            titleContent: self.title || self.$t('global.app.title'),
            textContent: self.text || '',
            finalColor: self.color || 'primary',
            resolve: null,
            reject: null
        }
    },
    watch: {
        'showState': function (newValue) {
            this.$emit('update:show', newValue);
        }
    },
    methods: {
        open(title, text, options) {
            this.showState = true;

            if (isString(text)) {
                this.titleContent = this.$t(title);
                this.textContent = this.$t(text);
            } else {
                this.titleContent = this.$t('global.app.title');
                this.textContent = this.$t(title);
                options = text;
            }

            if (options && options.color) {
                this.finalColor = options.color || 'primary';
            }

            return new Promise((resolve, reject) => {
                this.resolve = resolve;
                this.reject = reject;
            });
        },
        confirm() {
            if (this.resolve) {
                this.resolve();
            }

            this.showState = false;
            this.$emit('update:show', false);
        },
        cancel() {
            if (this.reject) {
                this.reject();
            }

            this.showState = false;
            this.$emit('update:show', false);
        }
    }
}
</script>
