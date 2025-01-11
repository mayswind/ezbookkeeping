<template>
    <div class="pin-codes-input" :style="`grid-template-columns: repeat(${length}, minmax(0, 1fr))`">
        <div class="pin-code-input pin-code-input-outline"
             :class="{ 'pin-code-input-focued': codes[index].focused }" :key="index"
             v-for="(code, index) in codes">
            <input ref="pin-code-input" min="0" maxlength="1" pattern="[0-9]*"
                   :value="codes[index].value"
                   :type="codes[index].inputType"
                   :disabled="disabled ? 'disabled' : undefined"
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

<script setup lang="ts">
import { ref, computed, watch, useTemplateRef } from 'vue';

interface PinCode {
    value: string;
    inputType: string;
    inputTimer: unknown | null;
    focused: boolean;
}

const props = defineProps<{
    modelValue: string;
    length: number;
    disabled?: boolean;
    autofocus?: boolean;
    autoConfirm?: boolean;
    secure?: boolean;
}>();

const emit = defineEmits<{
    (e: 'update:modelValue', value: string): void;
    (e: 'pincode:confirm', value: string): void;
}>();

const codes = ref<PinCode[]>([]);
const pinCodeInputs = useTemplateRef<HTMLInputElement[]>('pin-code-input');

const finalPinCode = computed<string>(() => {
    let ret = '';

    for (let i = 0; i < codes.value.length; i++) {
        if (codes.value[i].value) {
            ret += codes.value[i].value;
        } else {
            break;
        }
    }

    return ret;
});

function init(length: number, value: string): void {
    codes.value.length = 0;

    for (let i = 0; i < length; i++) {
        const code: PinCode = {
            value: '',
            inputType: 'tel',
            inputTimer: null,
            focused: false
        };

        if (value && value[i]) {
            code.value = value[i];

            if (props.secure) {
                code.inputType = 'password';
            }
        }

        codes.value.push(code);
    }
}

function autoFillText(index: number, text: string): void {
    let lastIndex = index;

    for (let i = index, j = 0; i < codes.value.length && j < text.length; i++, j++) {
        if (text[j] < '0' || text[j] > '9') {
            codes.value[i].value = '';
            break;
        }

        codes.value[i].value = text[j];
        setInputType(i);
        lastIndex = i;
    }

    setFocus(lastIndex);

    if (finalPinCode.value.length === length) {
        emit('pincode:confirm', finalPinCode.value);
    }
}

function setInputType(index: number): void {
    if (!props.secure) {
        return;
    }

    if (!codes.value[index].value) {
        codes.value[index].inputType = 'tel';
        return;
    }

    if (codes.value[index].inputTimer) {
        return;
    }

    codes.value[index].inputTimer = setTimeout(() => {
        if (codes.value[index].value) {
            codes.value[index].inputType = 'password';
        } else {
            codes.value[index].inputType = 'tel';
        }

        codes.value[index].inputTimer = null;
    }, 300);
}

function setFocus(index: number): void {
    if (pinCodeInputs.value && pinCodeInputs.value[index]) {
        pinCodeInputs.value[index].focus();
        pinCodeInputs.value[index].select();
    }
}

function setPreviousFocus(index: number): void {
    if (index > 0) {
        setFocus(index - 1);
    }
}

function setNextFocus(index: number): void {
    if (index < props.length - 1) {
        setFocus(index + 1);
    }
}

function onKeydown(index: number, event: KeyboardEvent): void {
    if (event.altKey || (event.key.indexOf('F') === 0 && (event.key.length === 2 || event.key.length === 3))) {
        return;
    }

    if (index <= 0 && (event.shiftKey && event.key === 'Tab')) {
        return;
    }

    if (index >= props.length - 1 && (!event.shiftKey && event.key === 'Tab')) {
        return;
    }

    if (event.key === 'Enter' && finalPinCode.value.length === props.length) {
        emit('pincode:confirm', finalPinCode.value);
        event.preventDefault();
        return;
    }

    if (event.key === 'ArrowLeft' || (event.shiftKey && event.key === 'Tab')) {
        setPreviousFocus(index);
        event.preventDefault();
        return;
    }

    if (event.key === 'ArrowRight' || (!event.shiftKey && event.key === 'Tab')) {
        setNextFocus(index);
        event.preventDefault();
        return;
    }

    if (event.key === 'Home') {
        setFocus(0);
        event.preventDefault();
        return;
    }

    if (event.key === 'End') {
        setFocus(props.length - 1);
        event.preventDefault();
        return;
    }

    if (((event.ctrlKey || event.metaKey) && event.key === 'v') || event.key === 'Paste') {
        return;
    }

    if (event.key === 'Backspace' || event.key === 'Delete' || event.key === 'Del') {
        for (let i = index; i < codes.value.length; i++) {
            codes.value[i].value = '';
            setInputType(i);
        }

        if (event.code === 'Backspace') {
            setPreviousFocus(index);
        }

        event.preventDefault();
        return;
    }

    if (event.key.length === 1 && '0' <= event.key && event.key <= '9') {
        codes.value[index].value = event.key;
        setInputType(index);
        setNextFocus(index);

        if (props.autoConfirm && finalPinCode.value.length === props.length) {
            emit('pincode:confirm', finalPinCode.value);
        }
    }

    event.preventDefault();
}

function onPaste(index: number, event: ClipboardEvent): void {
    if (!event.clipboardData) {
        event.preventDefault();
        return;
    }

    const text = event.clipboardData.getData('Text');

    if (!text) {
        event.preventDefault();
        return;
    }

    autoFillText(index, text);

    event.preventDefault();
}

function onInput(index: number, event: Event | { target: { value: string }, preventDefault: () => void }): void {
    if (!event.target || !(event.target as { value: string }).value) {
        event.preventDefault();
        return;
    }

    autoFillText(index, (event.target as { value: string }).value);

    event.preventDefault();
}

watch(() => props.length, newValue => {
    init(newValue, props.modelValue);
});

watch(() => props.modelValue, newValue => {
    if (newValue === finalPinCode.value) {
        return;
    }

    init(props.length, newValue);
});

watch(codes, () => {
    emit('update:modelValue', finalPinCode.value);
}, {
    deep: true
});

init(props.length, props.modelValue);
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
