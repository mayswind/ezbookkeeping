<template>
    <div class="btn-vertical-group d-flex flex-column">
        <v-btn border :key="idx"
               :color="value === button.value ? 'primary' : 'default'"
               :variant="value === button.value ? 'tonal' : 'outlined'" :disabled="disabled"
               v-for="(button, idx) in buttons"
               @click="value = button.value">
            {{ button.name }}
        </v-btn>
    </div>
</template>

<script>
export default {
    props: [
        'disabled',
        'buttons',
        'modelValue'
    ],
    emits: [
        'update:modelValue'
    ],
    computed: {
        value: {
            get: function () {
                return this.modelValue;
            },
            set: function (value) {
                if (value === this.modelValue) {
                    return;
                }

                this.$emit('update:modelValue', value);
            }
        }
    }
}
</script>

<style>
.btn-vertical-group .v-btn:not(:first-child) {
    border-top-left-radius: inherit;
    border-top-right-radius: inherit;
}

.btn-vertical-group .v-btn:not(:last-child) {
    border-bottom: 0;
    border-bottom-left-radius: inherit;
    border-bottom-right-radius: inherit;
}
</style>
