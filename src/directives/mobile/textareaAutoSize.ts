import { autoChangeTextareaSize } from '@/lib/ui/mobile.ts';

export default {
    mounted(el: HTMLElement): void {
        autoChangeTextareaSize(el);
    },
    updated(el: HTMLElement): void {
        autoChangeTextareaSize(el);
    }
}
