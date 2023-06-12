<template>
    <f7-page>
        <f7-navbar>
            <f7-nav-left :back-link="$t('Back')"></f7-nav-left>
            <f7-nav-title :title="$t('Font Size')"></f7-nav-title>
            <f7-nav-right>
                <f7-link :text="$t('Done')" @click="setFontSize"></f7-link>
            </f7-nav-right>
        </f7-navbar>

        <f7-list strong inset>
            <f7-list-item>
                <div class="full-line padding-bottom padding-top-half">
                    <div class="display-flex justify-content-space-between">
                        <div class="fontsize-minimum">A</div>
                        <div class="fontsize-maximum">A</div>
                        <div class="fontsize-default" :style="'left: ' + (100 / maxFontSizeType - 6) + '%'">{{ $t('Default') }}</div>
                    </div>
                    <f7-range
                        :min="minFontSizeType"
                        :max="maxFontSizeType"
                        :step="1"
                        :scale="true"
                        :scale-steps="maxFontSizeType"
                        :scale-sub-steps="1"
                        :format-scale-label="getFontSizeName"
                        v-model:value="fontSize"
                    />
                </div>
            </f7-list-item>
        </f7-list>
    </f7-page>
</template>

<script>
import fontConstants from '@/consts/font.js';
import { setAppFontSize } from '@/lib/ui.mobile.js';

export default {
    props: [
        'f7router'
    ],
    data() {
        return {
            fontSize: this.$settings.getFontSize()
        }
    },
    computed: {
        minFontSizeType() {
            return 0;
        },
        maxFontSizeType() {
            return fontConstants.allFontSizeArray.length - 1;
        }
    },
    methods: {
        setFontSize() {
            const router = this.f7router;

            if (this.fontSize !== this.$settings.getFontSize()) {
                this.$settings.setFontSize(this.fontSize);
                setAppFontSize(this.fontSize);
            }

            router.back();
        },
        getFontSizeName() {
            return '';
        }
    }
}
</script>

<style>
.fontsize-minimum {
    font-size: 15px;
    align-self: end;
}

.fontsize-maximum {
    font-size: 24px;
    align-self: end;
}

.fontsize-default {
    font-size: 17px;
    position: absolute;
    align-self: end;
}
</style>
