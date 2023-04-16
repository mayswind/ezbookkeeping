import { autoChangeTextareaSize } from '../../lib/mobile/ui.js';

export default {
    mounted(el) {
        autoChangeTextareaSize(el);
    },
    updated(el) {
        autoChangeTextareaSize(el);
    }
}
