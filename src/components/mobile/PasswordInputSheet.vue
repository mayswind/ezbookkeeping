<template>
    <f7-sheet swipe-to-close swipe-handler=".swipe-handler" style="height:auto"
              :opened="show" @sheet:open="onSheetOpen" @sheet:closed="onSheetClosed">
        <div class="swipe-handler" style="z-index: 10"></div>
        <f7-page-content class="margin-top no-padding-top">
            <div class="display-flex padding justify-content-space-between align-items-center">
                <div class="ebk-sheet-title" v-if="title"><b>{{ title }}</b></div>
            </div>
            <div class="padding-horizontal padding-bottom">
                <p class="no-margin" v-if="hint">{{ hint }}</p>
                <f7-list strong class="no-margin">
                    <f7-list-input
                        type="password"
                        autocomplete="current-password"
                        outline
                        floating-label
                        clear-button
                        class="no-margin no-padding-bottom"
                        :class="color ? 'color-' + color : ''"
                        :label="$t('Current Password')"
                        :placeholder="$t('Current Password')"
                        v-model:value="currentPassword"
                        @keyup.enter="confirm()"
                    ></f7-list-input>
                </f7-list>
                <f7-button large fill
                           :class="{ 'disabled': !currentPassword || confirmDisabled }"
                           :color="color || 'primary'"
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
        'color',
        'confirmDisabled',
        'cancelDisabled',
        'show'
    ],
    emits: [
        'update:modelValue',
        'update:show',
        'password:confirm'
    ],
    data() {
        return {
            currentPassword: ''
        }
    },
    methods: {
        onSheetOpen() {
            this.currentPassword = '';
        },
        onSheetClosed() {
            this.close();
        },
        confirm() {
            if (!this.currentPassword || this.confirmDisabled) {
                return;
            }

            this.$emit('update:modelValue', this.currentPassword);
            this.$emit('password:confirm', this.currentPassword);
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
