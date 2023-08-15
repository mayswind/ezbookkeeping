<template>
    <f7-sheet swipe-to-close swipe-handler=".swipe-handler" style="height:auto"
              :opened="show" @sheet:open="onSheetOpen" @sheet:closed="onSheetClosed">
        <f7-toolbar>
            <div class="swipe-handler"></div>
            <div class="left"></div>
            <div class="right">
                <f7-link :text="$t('Done')" @click="save"></f7-link>
            </div>
        </f7-toolbar>
        <f7-page-content class="no-margin-vertical no-padding-vertical">
            <map-view ref="map" height="400px" :geo-location="geoLocation">
                <template #error-title="{ mapSupported, mapDependencyLoaded }">
                    <div class="display-flex padding justify-content-space-between align-items-center">
                        <div class="ebk-sheet-title" v-if="!mapSupported"><b>{{ $t('Unsupported Map Provider') }}</b></div>
                        <div class="ebk-sheet-title" v-else-if="!mapDependencyLoaded"><b>{{ $t('Cannot Initialize Map') }}</b></div>
                        <div class="ebk-sheet-title" v-else></div>
                    </div>
                </template>
                <template #error-content>
                    <div class="padding-horizontal padding-bottom">
                        <p class="no-margin">{{ $t('Please refresh the page and try again. If the error is still displayed, make sure that server map settings are set correctly.') }}</p>
                        <div class="margin-top text-align-center">
                            <f7-link @click="close" :text="$t('Close')"></f7-link>
                        </div>
                    </div>
                </template>
            </map-view>
        </f7-page-content>
    </f7-sheet>
</template>

<script>
export default {
    props: [
        'modelValue',
        'show'
    ],
    emits: [
        'update:modelValue',
        'update:show'
    ],
    computed: {
        geoLocation: {
            get: function () {
                return this.modelValue;
            },
            set: function (value) {
                this.$emit('update:modelValue', value);
            }
        }
    },
    methods: {
        save() {
            this.$emit('update:show', false);
        },
        onSheetOpen() {
            this.$refs.map.init();
        },
        onSheetClosed() {
            this.close();
        },
        close() {
            this.$emit('update:show', false);
        }
    }
}
</script>
