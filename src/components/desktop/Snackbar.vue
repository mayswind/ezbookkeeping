<template>
    <v-snackbar v-model="showState">
        {{ messageContent }}

        <template #actions>
            <v-btn color="primary" variant="text" @click="showState = false">{{ $t('Close') }}</v-btn>
        </template>
    </v-snackbar>
</template>

<script>
export default {
    props: [
        'show',
        'message'
    ],
    emits: [
        'update:show'
    ],
    expose: [
        'showMessage',
        'showError'
    ],
    data() {
        const self = this;

        return {
            showState: self.show,
            messageContent: self.message,
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
        showMessage(message) {
            this.showState = true;
            this.messageContent = this.$t(message);
        },
        showError(error) {
            this.showState = true;
            this.messageContent = this.$tError(error.message || error);
        }
    }
}
</script>
