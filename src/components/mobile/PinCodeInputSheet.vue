<template>
    <f7-sheet swipe-to-close swipe-handler=".swipe-handler" style="height:auto"
              :opened="show" @sheet:open="onSheetOpen" @sheet:closed="onSheetClosed">
        <div class="swipe-handler" style="z-index: 10"></div>
        <f7-page-content class="margin-top no-padding-top">
            <div class="display-flex padding justify-content-space-between align-items-center">
                <div class="ebk-sheet-title"><b>{{ title }}</b></div>
            </div>
            <div class="padding-horizontal padding-bottom">
                <p class="no-margin">{{ hint }}</p>
                <f7-list class="no-margin">
                    <f7-list-item class="list-item-pincode-input padding-vertical-half">
                        <pin-code-input :secure="true" :length="6" v-model="currentPinCode" @pincode:confirm="confirm"/>
                    </f7-list-item>
                </f7-list>
                <f7-button large fill
                           :class="{ 'disabled': !currentPinCodeValid || confirmDisabled }"
                           :text="$t('Continue')"
                           @click="confirm">
                </f7-button>
                <div class="margin-top text-align-center">
                    <f7-link @click="cancel" :text="$t('Cancel')"></f7-link>
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
    emits: [
        'update:modelValue',
        'update:show',
        'pincode:confirm'
    ],
    data() {
        return {
            currentPinCode: ''
        }
    },
    computed: {
        currentPinCodeValid() {
            return this.currentPinCode && this.currentPinCode.length === 6;
        },
    },
    methods: {
        onSheetOpen() {
            this.currentPinCode = '';
        },
        onSheetClosed() {
            this.$emit('update:show', false);
        },
        confirm() {
            if (!this.currentPinCodeValid || this.confirmDisabled) {
                return;
            }

            this.$emit('update:modelValue', this.currentPinCode);
            this.$emit('pincode:confirm', this.currentPinCode);
        },
        cancel() {
            this.$emit('update:show', false);
        }
    }
}
</script>

<style>
.list-item-pincode-input .item-content {
    padding-left: 0;
    padding-right: 0;
}

.list-item-pincode-input .item-content .item-inner {
    padding-left: 0;
    padding-right: 0;
    justify-content: center;
}
</style>
