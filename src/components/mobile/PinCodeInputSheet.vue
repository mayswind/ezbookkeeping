<template>
    <f7-sheet style="height:auto" :opened="show"
              @sheet:open="onSheetOpen" @sheet:closed="onSheetClosed">
        <f7-page-content>
            <div class="display-flex padding justify-content-space-between align-items-center">
                <div style="font-size: 18px"><b>{{ title }}</b></div>
            </div>
            <div class="padding-horizontal padding-bottom">
                <p class="no-margin-top margin-bottom-half">{{ hint }}</p>
                <f7-list no-hairlines class="no-margin-top margin-bottom">
                    <f7-list-item class="list-item-pincode-input">
                        <PincodeInput secure :length="6" v-model="currentPinCode" />
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
        'value',
        'title',
        'hint',
        'confirmDisabled',
        'cancelDisabled',
        'show'
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

            this.$emit('input', this.currentPinCode);
            this.$emit('pincode:confirm', this.currentPinCode);
        },
        cancel() {
            this.$emit('update:show', false);
        }
    }
}
</script>

<style>
.list-item-pincode-input .item-inner {
    justify-content: center;
}
</style>
