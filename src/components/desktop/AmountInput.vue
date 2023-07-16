<template>
    <v-text-field type="number" :class="extraClass" :density="density" :disabled="disabled"
                  :rules="enableRules ? rules : []" v-model="value"
                  @keydown="onKeyUpDown" @keyup="onKeyUpDown"
    ></v-text-field>
</template>

<script>
import transactionConstants from '@/consts/transaction.js';

export default {
    props: [
        'class',
        'density',
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
        }
    }
}
</script>
