<template>
    <v-text-field type="number" :class="extraClass"
                  :density="density" :disabled="disabled"
                  :label="label"
                  :placeholder="placeholder"
                  :persistent-placeholder="persistentPlaceholder"
                  :rules="enableRules ? rules : []" v-model="value"
                  @keydown="onKeyUpDown" @keyup="onKeyUpDown">
        <template #prepend-inner v-if="currency && prependText">
            <div style="margin-top: 2px">{{ prependText }}</div>
        </template>
        <template #append-inner v-if="currency && appendText">
            <div class="text-no-wrap">{{ appendText }}</div>
        </template>
    </v-text-field>
</template>

<script>
import { mapStores } from 'pinia';
import { useSettingsStore } from '@/stores/setting.js';

import transactionConstants from '@/consts/transaction.js';

export default {
    props: [
        'class',
        'density',
        'currency',
        'label',
        'placeholder',
        'persistentPlaceholder',
        'disabled',
        'enableRules',
        'modelValue'
    ],
    emits: [
        'update:modelValue'
    ],
    data() {
        const self = this;

        return {
            rules: [
                (v) => {
                    if (v === '') {
                        return self.$t('Amount value is not number');
                    }

                    try {
                        const val = parseFloat(v);
                        return (val >= transactionConstants.minAmount && val <= transactionConstants.maxAmount) || self.$t('Amount value exceeds limitation');
                    } catch (e) {
                        return self.$t('Amount value is not number');
                    }
                }
            ]
        }
    },
    computed: {
        ...mapStores(useSettingsStore),
        value: {
            get: function () {
                return this.modelValue;
            },
            set: function (value) {
                this.$emit('update:modelValue', value);
            }
        },
        extraClass() {
            return this.class;
        },
        prependText() {
            if (!this.currency) {
                return '';
            }

            const texts = this.getDisplayCurrencyPrependAndAppendText();

            if (!texts) {
                return '';
            }

            return texts.prependText;
        },
        appendText() {
            if (!this.currency) {
                return '';
            }

            const texts = this.getDisplayCurrencyPrependAndAppendText();

            if (!texts) {
                return '';
            }

            return texts.appendText;
        }
    },
    methods: {
        onKeyUpDown(e) {
            if (e.target.value === '' || e.target.value === 0) {
                return;
            }

            let decimalLength = 0;

            if (e.target.value.indexOf('.') > 0) {
                decimalLength = e.target.value.length - e.target.value.indexOf('.') - 1;
            }

            if (decimalLength > 2) {
                e.target.value = e.target.value.substring(0, e.target.value.length - 1);
                e.preventDefault();
                return;
            }

            try {

                const val = parseFloat(e.target.value);
                let maxLength = transactionConstants.maxAmount.toString().length;

                if (val < 0) {
                    maxLength = transactionConstants.minAmount.toString().length;
                }

                if (val < transactionConstants.minAmount) {
                    e.target.value = transactionConstants.minAmount;
                    e.preventDefault();
                } else if (val > transactionConstants.maxAmount) {
                    e.target.value = transactionConstants.maxAmount;
                    e.preventDefault();
                } else if (e.target.value.length > maxLength) {
                    e.target.value = e.target.value.substring(0, maxLength);
                    e.preventDefault();
                } else if (e.target.value.charAt(0) === '0' && e.target.value.length > 1) {
                    e.target.value = e.target.value.substring(1);
                    e.preventDefault();
                }
            } catch (e) {
                e.target.value = 0;
            }
        },
        getDisplayCurrencyPrependAndAppendText() {
            return this.$locale.getDisplayCurrencyPrependAndAppendText(this.currency, this.settingsStore.appSettings.currencyDisplayMode);
        }
    }
}
</script>
