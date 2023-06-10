import { autoChangeTextareaSize } from '@/lib/ui.mobile.js';

export default {
    mounted(el) {
        autoChangeTextareaSize(el);
    },
    updated(el) {
        autoChangeTextareaSize(el);
    }
}
