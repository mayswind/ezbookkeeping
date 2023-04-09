<template>
    <f7-sheet style="height:auto" :opened="show"
              @sheet:open="onSheetOpen" @sheet:closed="onSheetClosed">
        <f7-page-content>
            <div class="display-flex padding justify-content-space-between align-items-center">
                <div style="font-size: 18px" v-if="title"><b>{{ title }}</b></div>
            </div>
            <div class="padding-horizontal padding-bottom">
                <p class="no-margin-top margin-bottom-half" v-if="hint">{{ hint }}</p>
                <f7-list no-hairlines strong class="no-margin-top margin-bottom">
                    <f7-list-input
                        type="password"
                        autocomplete="current-password"
                        outline
                        floating-label
                        clear-button
                        class="no-margin"
                        :label="$t('Password')"
                        :placeholder="$t('Password')"
                        v-model:value="currentPassword"
                        @keyup.enter="confirm()"
                    ></f7-list-input>
                </f7-list>
                <f7-button large fill
                           :class="{ 'disabled': !currentPassword || confirmDisabled }"
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
            currentPassword: ''
        }
    },
    methods: {
        onSheetOpen() {
            this.currentPassword = '';
        },
        onSheetClosed() {
            this.$emit('update:show', false);
        },
        confirm() {
            if (!this.currentPassword || this.confirmDisabled) {
                return;
            }

            this.$emit('update:modelValue', this.currentPassword);
            this.$emit('password:confirm', this.currentPassword);
        },
        cancel() {
            this.$emit('update:show', false);
        }
    }
}
</script>
