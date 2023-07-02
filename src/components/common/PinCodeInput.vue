<template>
    <div class="pin-codes-input" :style="`grid-template-columns: repeat(${length}, minmax(0, 1fr))`">
        <div class="pin-code-input pin-code-input-outline"
             :class="{ 'pin-code-input-focued': codes[index].focused }" :key="index"
             v-for="(code, index) in codes">
            <input min="0" maxlength="1" pattern="[0-9]*"
                   :ref="`pin-code-input-${index}`"
                   :value="codes[index].value"
                   :type="codes[index].inputType"
                   :autofocus="autofocus && index === 0 ? 'autofocus' : undefined"
                   @focus="codes[index].focused = true"
                   @blur="codes[index].focused = false"
                   @keydown="onKeydown(index, $event)"
                   @paste="onPaste(index, $event)"
                   @change="onInput(index, $event)"
            />
        </div>
    </div>
</template>

<script>
export default {
    props: [
        'modelValue',
        'autofocus',
        'secure',
        'length'
    ],
    emits: [
        'update:modelValue',
        'pincode:confirm'
    ],
    data() {
        return {
            codes: []
        }
    },
    computed: {
        finalPinCode() {
            let finalPinCode = '';

            for (let i = 0; i < this.codes.length; i++) {
                if (this.codes[i].value) {
                    finalPinCode += this.codes[i].value;
                } else {
                    break;
                }
            }

            return finalPinCode;
        }
    },
    watch: {
        'length': function (newValue) {
            this.init(newValue, this.modelValue);
        },
        'modelValue': function (newValue) {
            if (newValue === this.finalPinCode) {
                return;
            }

            this.init(this.length, newValue);
        },
        'codes': {
            handler() {
                this.$emit('update:modelValue', this.finalPinCode);
            },
            deep: true
        }
    },
    created() {
        this.init(this.length, this.modelValue);
    },
    methods: {
        init(length, value) {
            this.codes.length = 0;

            for (let i = 0; i < length; i++) {
                const code = {
                    value: '',
                    inputType: 'tel',
                    inputTimer: null,
                    focused: false
                };

                if (value && value[i]) {
                    code.value = value[i];

                    if (this.secure) {
                        code.inputType = 'password';
                    }
                }

                this.codes.push(code);
            }
        },
        autoFillText(index, text) {
            let lastIndex = index;

            for (let i = index, j = 0; i < this.codes.length && j < text.length; i++, j++) {
                if (text[j] < '0' || text[j] > '9') {
                    this.codes[i].value = '';
                    this.$forceUpdate();
                    break;
                }

                this.codes[i].value = text[j];
                this.setInputType(i);
                lastIndex = i;
            }

            this.setFocus(lastIndex);

            if (this.finalPinCode.length === this.length) {
                this.$emit('pincode:confirm', this.finalPinCode);
            }
        },
        setInputType(index) {
            const self = this;

            if (!self.secure) {
                return;
            }

            if (!self.codes[index].value) {
                self.codes[index].inputType = 'tel';
                return;
            }

            if (self.codes[index].inputTimer) {
                return;
            }

            self.codes[index].inputTimer = setTimeout(() => {
                if (self.codes[index].value) {
                    self.codes[index].inputType = 'password';
                } else {
                    self.codes[index].inputType = 'tel';
                }

                self.codes[index].inputTimer = null;
            }, 300);
        },
        setFocus(index) {
            const refId = `pin-code-input-${index}`;
            const ref = this.$refs[refId];

            if (ref && ref[0]) {
                ref[0].focus();
                ref[0].select();
            }
        },
        setPreviousFocus(index) {
            if (index > 0) {
                this.setFocus(index - 1);
            }
        },
        setNextFocus(index) {
            if (index < this.length - 1) {
                this.setFocus(index + 1);
            }
        },
        onKeydown(index, event) {
            if (event.code === 'Enter' && this.finalPinCode.length === this.length) {
                this.$emit('pincode:confirm', this.finalPinCode);
                event.preventDefault();
                return;
            }

            if (event.code === 'ArrowLeft' || (event.shiftKey && event.code === 'Tab')) {
                this.setPreviousFocus(index);
                event.preventDefault();
                return;
            }

            if (event.code === 'ArrowRight' || (!event.shiftKey && event.code === 'Tab')) {
                this.setNextFocus(index);
                event.preventDefault();
                return;
            }

            if ((event.ctrlKey || event.metaKey) && event.code === 'KeyV') {
                return;
            }

            if (event.code === 'Backspace' || event.code === 'Delete' || event.code === 'Del') {
                for (let i = index; i < this.codes.length; i++) {
                    this.codes[i].value = '';
                    this.setInputType(i);
                }

                if (event.code === 'Backspace') {
                    this.setPreviousFocus(index);
                }

                event.preventDefault();
                return;
            }

            if (event.code.indexOf('Digit') === 0 && event.code.length === 6) {
                this.codes[index].value = event.key;
                this.setInputType(index);
                this.setNextFocus(index);

                if (this.finalPinCode.length === this.length) {
                    this.$emit('pincode:confirm', this.finalPinCode);
                }
            }

            event.preventDefault();
        },
        onPaste(index, event) {
            if (!event.clipboardData) {
                event.preventDefault();
                return;
            }

            const text = event.clipboardData.getData('Text');

            if (!text) {
                event.preventDefault();
                return;
            }

            this.autoFillText(index, text);

            event.preventDefault();
        },
        onInput(index, event) {
            if (!event.target.value) {
                event.preventDefault();
                return;
            }

            this.autoFillText(index, event.target.value);

            event.preventDefault();
        }
    }
}
</script>

<style>
.pin-codes-input {
    --ebk-pin-code-border-color: #bbb;
    --ebk-pin-code-focued-color: #c67e48;
    --ebk-pin-code-border-radius: 8px;
    --ebk-pin-code-input-height: 46px;
    --ebk-pin-code-input-gap: 8px;
    --ebk-pin-code-transition-duration: 200ms;
    display: grid;
    gap: var(--ebk-pin-code-input-gap);
}

.pin-code-input {
    position: relative;
}

.pin-code-input input {
    text-align: center;
    padding-left: 10px;
    padding-right: 10px;
    width: 100%;
    height: var(--ebk-pin-code-input-height) !important;
}

.pin-code-input input:focus {
    outline: none;
}

.pin-code-input-outline::after {
    content: '';
    position: absolute;
    left: 0;
    right: 0;
    top: 0;
    bottom: 0;
    border: 1px solid var(--ebk-pin-code-border-color);
    border-radius: var(--ebk-pin-code-border-radius);
    pointer-events: none;
    box-sizing: border-box;
    transition-duration: var(--ebk-pin-code-transition-duration);
}

.pin-code-input-outline.pin-code-input-focued::after {
    border-width: 2px;
    border-color: var(--ebk-pin-code-focued-color);
}
</style>
