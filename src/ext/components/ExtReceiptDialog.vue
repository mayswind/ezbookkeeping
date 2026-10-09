<template>
    <v-dialog width="520" scrollable :model-value="show" @update:model-value="close">
        <v-card v-if="receipt">
            <v-card-text class="pa-0">
                <div class="d-flex align-center flex-wrap ga-3 pa-4 pb-2">
                    <v-btn-toggle density="compact" mandatory color="primary" variant="outlined" divided v-model="paper">
                        <v-btn value="narrow" size="small">{{ tt('Small receipt (80 mm)') }}</v-btn>
                        <v-btn value="a4" size="small">{{ tt('Full page (A4)') }}</v-btn>
                    </v-btn-toggle>
                </div>

                <div class="pa-4 pt-2 ext-receipt-stage">
                    <div ref="paperElement" class="ext-receipt-paper" :class="paper === 'a4' ? 'ext-receipt-a4' : 'ext-receipt-narrow'">
                        <div class="ext-receipt-void" v-if="receipt.voided">{{ tt('VOID') }}</div>

                        <div class="ext-receipt-center">
                            <div class="ext-receipt-business">{{ receipt.business.name }}</div>
                            <div v-if="receipt.business.address">{{ receipt.business.address }}</div>
                            <div v-if="receipt.business.phone">{{ receipt.business.phone }}</div>
                        </div>

                        <hr />
                        <div class="ext-receipt-row ext-receipt-strong">
                            <span>{{ receipt.title }}</span>
                            <span>{{ receipt.number }}</span>
                        </div>
                        <div>{{ receipt.date }}</div>
                        <div class="ext-receipt-row" :key="fact.label" v-for="fact in receipt.facts">
                            <span>{{ fact.label }}</span>
                            <span class="ext-receipt-value">{{ fact.value }}</span>
                        </div>

                        <template v-if="receipt.lines.length > 0">
                            <hr />
                            <div class="ext-receipt-line" :key="index" v-for="(line, index) in receipt.lines">
                                <div>{{ line.name }}</div>
                                <div class="ext-receipt-row ext-receipt-muted">
                                    <span>{{ line.detail }}</span>
                                    <span class="ext-receipt-value">{{ line.total }}</span>
                                </div>
                            </div>
                        </template>

                        <hr />
                        <div class="ext-receipt-row" :class="{ 'ext-receipt-strong': row.strong }" :key="row.label" v-for="row in receipt.totals">
                            <span>{{ row.label }}</span>
                            <span class="ext-receipt-value">{{ row.value }}</span>
                        </div>

                        <template v-if="receipt.notes.length > 0">
                            <hr />
                            <div class="ext-receipt-muted" :key="note" v-for="note in receipt.notes">{{ note }}</div>
                        </template>

                        <template v-if="receipt.business.footer">
                            <hr />
                            <div class="ext-receipt-center">{{ receipt.business.footer }}</div>
                        </template>
                    </div>
                </div>
            </v-card-text>
            <v-card-actions>
                <v-btn variant="text" @click="copyText">{{ tt('Copy as text') }}</v-btn>
                <v-spacer />
                <v-btn variant="text" @click="close(false)">{{ tt('Close') }}</v-btn>
                <v-btn color="primary" @click="print">{{ tt('Print') }}</v-btn>
            </v-card-actions>
        </v-card>
        <ext-snack-bar ref="snackbar" />
    </v-dialog>
</template>

<script setup lang="ts">
import ExtSnackBar from '@/ext/components/ExtSnackBar.vue';

import { ref, watch, useTemplateRef } from 'vue';

import { useExtI18n } from '@/ext/i18n.ts';

import { receiptToText, type Receipt } from '@/ext/receipt.ts';

type SnackBarType = InstanceType<typeof ExtSnackBar>;

const props = defineProps<{
    show: boolean;
    receipt: Receipt | null;
}>();

const emit = defineEmits<{
    (e: 'update:show', value: boolean): void;
}>();

const PAPER_KEY = 'ebk_ext_receipt_paper';
const PRINT_ROOT_ID = 'ext-print-root';
const PRINT_STYLE_ID = 'ext-print-style';

const { tt } = useExtI18n();

const paperElement = useTemplateRef<HTMLElement>('paperElement');
const snackbar = useTemplateRef<SnackBarType>('snackbar');

function readPaper(): 'narrow' | 'a4' {
    try {
        return localStorage.getItem(PAPER_KEY) === 'a4' ? 'a4' : 'narrow';
    } catch {
        return 'narrow';
    }
}

const paper = ref<'narrow' | 'a4'>(readPaper());

watch(paper, value => {
    try {
        localStorage.setItem(PAPER_KEY, value);
    } catch {
        // remembering the paper is a convenience only
    }
});

function close(value: boolean = false): void {
    emit('update:show', value);
}

// Printing copies the receipt into its own element on the page and hides everything else while the browser's print
// dialog is open, so only the receipt comes out, on the paper size chosen. "Save as PDF" is in the same dialog.
function print(): void {
    if (!paperElement.value) {
        return;
    }

    cleanup();

    const root = document.createElement('div');
    root.id = PRINT_ROOT_ID;
    root.innerHTML = paperElement.value.outerHTML;
    document.body.appendChild(root);

    const style = document.createElement('style');
    style.id = PRINT_STYLE_ID;
    style.textContent = `
        #${PRINT_ROOT_ID} { display: none; }
        @media print {
            body > *:not(#${PRINT_ROOT_ID}) { display: none !important; }
            #${PRINT_ROOT_ID} { display: block !important; }
            #${PRINT_ROOT_ID} .ext-receipt-paper { box-shadow: none !important; margin: 0 !important; }
        }
        @page { size: ${paper.value === 'a4' ? 'A4' : '80mm auto'}; margin: ${paper.value === 'a4' ? '15mm' : '0'}; }
    `;
    document.head.appendChild(style);

    window.addEventListener('afterprint', cleanup, { once: true });
    window.print();
}

function cleanup(): void {
    document.getElementById(PRINT_ROOT_ID)?.remove();
    document.getElementById(PRINT_STYLE_ID)?.remove();
}

// Plain text for WhatsApp, SMS or email. Falls back to selecting a hidden box where the clipboard API is unavailable
// (it needs a secure page).
async function copyText(): Promise<void> {
    if (!props.receipt) {
        return;
    }

    const text = receiptToText(props.receipt, paper.value === 'a4' ? 42 : 32);

    try {
        await navigator.clipboard.writeText(text);
    } catch {
        const box = document.createElement('textarea');
        box.value = text;
        box.style.position = 'fixed';
        box.style.opacity = '0';
        document.body.appendChild(box);
        box.select();
        const copied = document.execCommand('copy');
        document.body.removeChild(box);

        if (!copied) {
            snackbar.value?.showError(tt('Could not copy. Use Print and choose Save as PDF instead.'));
            return;
        }
    }

    snackbar.value?.showMessage(tt('Receipt copied. Paste it into a message.'));
}

defineExpose({ print });
</script>

<style>
.ext-receipt-stage {
    background: rgba(0, 0, 0, 0.04);
    display: flex;
    justify-content: center;
    overflow-x: auto;
}

.ext-receipt-paper {
    position: relative;
    background: #fff;
    color: #000;
    font-family: 'Courier New', Courier, monospace;
    font-size: 13px;
    line-height: 1.4;
    padding: 12px;
    box-shadow: 0 1px 6px rgba(0, 0, 0, 0.25);
    box-sizing: border-box;
}

.ext-receipt-narrow { width: 80mm; max-width: 100%; }
.ext-receipt-a4 { width: 100%; max-width: 190mm; font-size: 15px; padding: 24px; }
.ext-receipt-paper hr { border: 0; border-top: 1px dashed #000; margin: 8px 0; }
.ext-receipt-center { text-align: center; }
.ext-receipt-business { font-weight: bold; font-size: 1.2em; }
.ext-receipt-row { display: flex; justify-content: space-between; gap: 12px; }
.ext-receipt-value { text-align: right; white-space: nowrap; }
.ext-receipt-strong { font-weight: bold; }
.ext-receipt-muted { opacity: 0.8; }
.ext-receipt-line { margin-bottom: 4px; }

.ext-receipt-void {
    position: absolute;
    top: 38%;
    left: 50%;
    transform: translate(-50%, -50%) rotate(-18deg);
    font-size: 3em;
    font-weight: bold;
    letter-spacing: 0.2em;
    color: rgba(200, 0, 0, 0.35);
    border: 0.12em solid rgba(200, 0, 0, 0.35);
    padding: 0 0.3em;
    pointer-events: none;
}
</style>
