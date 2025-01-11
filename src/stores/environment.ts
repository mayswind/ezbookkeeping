import { ref } from 'vue';
import { defineStore } from 'pinia';

export const useEnvironmentsStore = defineStore('environments', () => {
    const framework7DarkMode = ref<boolean | undefined>(undefined);

    return {
        // states
        framework7DarkMode
    };
});
