<template>
    <f7-sheet swipe-to-close swipe-handler=".swipe-handler" style="height:auto"
              :opened="show" @sheet:closed="onSheetClosed">
        <div class="swipe-handler"></div>
        <f7-page-content class="margin-top no-padding-top">
            <div class="display-flex padding justify-content-space-between align-items-center">
                <div style="font-size: 18px" v-if="title"><b>{{ title }}</b></div>
            </div>
            <div class="padding-horizontal padding-bottom">
                <p class="no-margin-top margin-bottom-half" v-if="hint">
                    <span>{{ hint }}</span>
                    <f7-link id="copy-to-clipboard-icon" ref="copyToClipboardIcon"
                             class="icon-after-text"
                             icon-only icon-f7="doc_on_doc" icon-size="16px"
                             v-if="enableCopy"
                    ></f7-link>
                </p>
                <textarea class="information-content full-line" :rows="rowCount" :value="information"></textarea>
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
    data() {
        return {
            clipboardHolder: null
        }
    },
    mounted() {
        this.makeCopyToClipboardClickable();
    },
    updated() {
        this.makeCopyToClipboardClickable();
    },
    watch: {
        'information': function (newValue) {
            if (this.clipboardHolder) {
                this.$utilities.changeClipboardObjectTxet(this.clipboardHolder, newValue);
            }
        }
    },
    methods: {
        onSheetClosed() {
            this.close();
        },
        cancel() {
            this.close();
        },
        makeCopyToClipboardClickable() {
            const self = this;

            if (self.clipboardHolder) {
                return;
            }

            if (self.$refs.copyToClipboardIcon) {
                self.clipboardHolder = self.$utilities.makeButtonCopyToClipboard({
                    el: '#copy-to-clipboard-icon',
                    text: self.information,
                    successCallback: function () {
                        self.$emit('info:copied');
                    }
                });
            }
        },
        close() {
            this.$emit('update:show', false);
        }
    }
}
</script>
