<template>
    <f7-sheet style="height:auto" :opened="show"
              @sheet:open="onSheetOpen" @sheet:closed="onSheetClosed">
        <f7-page-content>
            <div class="display-flex padding justify-content-space-between align-items-center">
                <div style="font-size: 18px" v-if="title"><b>{{ title }}</b></div>
            </div>
            <div class="padding-horizontal padding-bottom">
                <p class="no-margin-top margin-bottom-half" v-if="hint">{{ hint }}</p>
                <slot></slot>
                <f7-list no-hairlines class="no-margin-top margin-bottom">
                    <f7-list-input
                        type="number"
                        autocomplete="one-time-code"
                        outline
                        clear-button
                        :placeholder="$t('Passcode')"
                        :value="currentPasscode"
                        @input="currentPasscode = $event.target.value"
                        @keyup.enter.native="confirm()"
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
        'value',
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
            this.$emit('update:show', false);
        },
        confirm() {
            if (!this.currentPasscode || this.confirmDisabled) {
                return;
            }

            this.$emit('input', this.currentPasscode);
            this.$emit('passcode:confirm', this.currentPasscode);
        },
        cancel() {
            this.$emit('update:show', false);
        }
    }
}
</script>
