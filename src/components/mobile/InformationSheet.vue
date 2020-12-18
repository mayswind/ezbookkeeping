<template>
    <f7-sheet style="height:auto" :opened="show" @sheet:closed="onSheetClosed">
        <f7-page-content>
            <div class="display-flex padding justify-content-space-between align-items-center">
                <div style="font-size: 18px" v-if="title"><b>{{ title }}</b></div>
            </div>
            <div class="padding-horizontal padding-bottom">
                <p class="no-margin-top margin-bottom-half" v-if="hint">
                    <span>{{ hint }}</span>
                    <f7-link class="icon-after-text"
                             icon-only icon-f7="doc_on_doc" icon-size="16px"
                             v-if="enableCopy"
                             v-clipboard:copy="information" v-clipboard:success="onCopied"></f7-link>
                </p>
                <textarea class="information-content full-line" :rows="rowCount" readonly="readonly" v-model="information"></textarea>
                <div class="margin-top text-align-center">
                    <f7-link @click="cancel" :text="$t('Close')"></f7-link>
                </div>
            </div>
        </f7-page-content>
    </f7-sheet>
</template>

<script>
export default {
    props: [
        'title',
        'hint',
        'information',
        'rowCount',
        'enableCopy',
        'show'
    ],
    methods: {
        onSheetClosed() {
            this.$emit('update:show', false);
        },
        onCopied() {
            this.$emit('info:copied');
        },
        cancel() {
            this.$emit('update:show', false);
        }
    }
}
</script>
