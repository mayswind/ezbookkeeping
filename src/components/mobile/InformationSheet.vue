<template>
    <f7-sheet swipe-to-close swipe-handler=".swipe-handler" style="height:auto"
              :opened="show" @sheet:closed="onSheetClosed">
        <div class="swipe-handler" style="z-index: 10"></div>
        <f7-page-content class="margin-top no-padding-top">
            <div class="display-flex padding justify-content-space-between align-items-center">
                <div class="ebk-sheet-title" v-if="title"><b>{{ title }}</b></div>
            </div>
            <div class="padding-horizontal padding-bottom">
                <p class="no-margin-top margin-bottom-half" v-if="hint">
                    <span>{{ hint }}</span>
                    <f7-link id="copy-to-clipboard-icon" ref="copyToClipboardIcon"
                             class="icon-after-text"
                             icon-only icon-f7="doc_on_doc"
                             v-if="enableCopy"
                    ></f7-link>
                </p>
                <textarea class="information-content full-line" readonly="readonly" :rows="rowCount" :value="information"></textarea>
                <div class="margin-top text-align-center">
                    <f7-link @click="cancel" :text="$t('Close')"></f7-link>
                </div>
            </div>
        </f7-page-content>
    </f7-sheet>
</template>

<script>
import { makeButtonCopyToClipboard, changeClipboardObjectText } from '@/lib/misc.js';

export default {
    props: [
        'title',
        'hint',
        'information',
        'rowCount',
        'enableCopy',
        'show'
    ],
    emits: [
        'update:show',
        'info:copied'
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
                changeClipboardObjectText(this.clipboardHolder, newValue);
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
                self.clipboardHolder = makeButtonCopyToClipboard({
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
