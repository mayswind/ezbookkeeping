<template>
    <div :class="imageBoxClass" :style="style">
        <img class="image-with-placeholder" :class="{ 'image-loading': loading }"
             :src="src" :alt="alt" v-if="!link && !loadError"
             @load="onLoad" @error="onError"/>
        <f7-link class="image-link" :class="{ 'image-loading': loading }"
                 :href="link" v-if="link && !loadError">
            <img class="image-with-placeholder" :src="src" :alt="alt"
                 @load="onLoad" @error="onError"/>
        </f7-link>
        <div class="image-loading-hint" v-if="loading && !loadError">
            <f7-preloader size="28" />
        </div>
        <div class="image-error-hint" v-if="!link && !loading && loadError">
            <slot name="error"></slot>
        </div>
        <f7-link class="image-error-hint" :href="link" v-if="link && !loading && loadError">
            <slot name="error"></slot>
        </f7-link>
    </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue';

const props = defineProps<{
    src: string;
    class?: string;
    style?: Record<string, string>;
    alt?: string;
    link?: string;
}>();

const loading = ref<boolean>(true);
const loadError = ref<boolean>(false);

const imageBoxClass = computed<string>(() => {
    let classes = 'image-box';

    if (props.class) {
        classes += ` ${props.class}`;
    }

    return classes;
});

function onLoad(): void {
    loading.value = false;
}

function onError(): void {
    loading.value = false;
    loadError.value = true;
}

watch(() => props.src, () => {
    loading.value = true;
    loadError.value = false;
});
</script>

<style scoped>
.image-box > .image-with-placeholder,
.image-box > .image-link,
.image-box > .image-link > .image-with-placeholder {
    width: 100%;
    height: 100%;
    display: block;
    object-fit: cover;
}

.image-box > .image-with-placeholder.image-loading,
.image-box > .image-link.image-loading {
    display: none !important;
}

.image-box > .image-loading-hint {
    width: 100%;
    height: 100%;
    display: flex;
    justify-content: center;
    align-items: center;
}

.image-box > .image-error-hint {
    display: flex;
    position: absolute;
    inset: 0;
    align-items: center;
    justify-content: center;
    text-align: center;
    font-size: var(--f7-list-item-footer-font-size);
    color: var(--f7-list-item-footer-text-color);
    padding: 4px;
    overflow: hidden;
}
</style>
