<template>
    <f7-sheet swipe-to-close swipe-handler=".swipe-handler" style="height:auto"
              :opened="show" @sheet:open="onSheetOpen" @sheet:closed="onSheetClosed">
        <div class="swipe-handler"></div>
        <f7-page-content class="margin-top no-padding-top">
            <div class="display-flex padding justify-content-space-between align-items-center">
                <div style="font-size: 18px" v-if="title"><b>{{ title }}</b></div>
            </div>
            <div class="padding-horizontal padding-bottom">
                <p class="no-margin" v-if="hint">{{ hint }}</p>
                <slot></slot>
                <f7-list no-hairlines strong class="no-margin">
                    <f7-list-input
                        type="number"
                        autocomplete="one-time-code"
                        outline
                        floating-label
                        clear-button
                        class="no-margin no-padding-bottom"
                        :label="$t('Password')"
                        :placeholder="$t('Passcode')"
                        v-model:value="currentPasscode"
                        @keyup.enter="confirm()"
                    ></f7-list-input>
                </f7-list>
                <f7-button large fill
                           :class="{ 'disabled': !currentPasscode || confirmDisabled }"
                           :text="$t('Continue')"
                           @click="confirm">
                </f7-button>
                <div class="margin-top text-align-center">
                    <f7-link :class="{ 'disabled': cancelDisabled }" @click="cancel" :text="$t('Cancel')"></f7-link>
                </div>
            </div>
        </f7-page-content>
    </f7-sheet>
</template>

<script>
export default {
    props: [
        'modelValue',
        'title',
        'hint',
        'confirmDisabled',
        'cancelDisabled',
        'show'
    ],
    data() {
        return {
            currentPasscode: ''
        }
    },
    methods: {
        onSheetOpen() {
            this.currentPasscode = '';
        },
        onSheetClosed() {
            this.close();
        },
        confirm() {
            if (!this.currentPasscode || this.confirmDisabled) {
                return;
            }

            this.$emit('update:modelValue', this.currentPasscode);
            this.$emit('passcode:confirm', this.currentPasscode);
        },
        cancel() {
            this.close();
        },
        close() {
            this.$emit('update:show', false);
        }
    }
}
</script>
