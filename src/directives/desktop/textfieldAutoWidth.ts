import type { DirectiveBinding } from 'vue';

interface AutoWidthOptions {
    auxSpanId: string;
    minWidth?: number;
    maxWidth?: number;
}

function updateWidth(el: HTMLElement, options: AutoWidthOptions, initStyle: boolean): void {
    const input = el.querySelector('input');

    if (!input) {
        return;
    }

    const auxEl = el.parentElement?.querySelector(`span#${options.auxSpanId}`);

    if (!auxEl) {
        return;
    }

    const span = auxEl as HTMLSpanElement;

    if (initStyle) {
        const inputStyle = window.getComputedStyle(input);
        span.style.position = 'absolute';
        span.style.visibility = 'hidden';
        span.style.whiteSpace = 'pre';
        span.style.font = inputStyle.font;
        span.style.letterSpacing = inputStyle.letterSpacing;
        span.style.padding = '0';
        span.style.margin = '0';
    }

    span.textContent = input.value || input.placeholder || '';

    let width: number = span.offsetWidth;

    if (options.minWidth) {
        width = Math.max(width, options.minWidth);
    }

    if (options.maxWidth) {
        width = Math.min(width, options.maxWidth);
    }

    el.style.width = `${width}px`;
}

export default {
    mounted(el: HTMLElement, binding: DirectiveBinding<AutoWidthOptions>): void {
        updateWidth(el, binding.value, true);
    },
    updated(el: HTMLElement, binding: DirectiveBinding<AutoWidthOptions>): void {
        updateWidth(el, binding.value, false);
    }
}
