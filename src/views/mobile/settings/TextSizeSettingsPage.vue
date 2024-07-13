<template>
    <f7-page>
        <f7-navbar>
            <f7-nav-left :back-link="$t('Back')"></f7-nav-left>
            <f7-nav-title :title="$t('Text Size')"></f7-nav-title>
            <f7-nav-right>
                <f7-link :text="$t('Done')" @click="setFontSize"></f7-link>
            </f7-nav-right>
        </f7-navbar>

        <f7-block class="no-padding no-margin readonly" :class="fontSizePreviewClassName">
            <f7-block class="combination-list-wrapper margin-vertical">
                <f7-accordion-item>
                    <f7-block-title>
                        <f7-accordion-toggle>
                            <f7-list strong inset dividers media-list
                                     class="transaction-amount-list combination-list-header combination-list-opened">
                                <f7-list-item>
                                    <template #title>
                                        <small>{{ currentLongYearMonth }}</small>
                                        <small class="transaction-amount-statistics">
                                            <span class="text-income">{{ `+${getDisplayAmount('12345')}` }}</span>
                                            <span class="text-expense">{{ `-${getDisplayAmount('67890')}` }}</span>
                                        </small>
                                        <f7-icon class="combination-list-chevron-icon" f7="chevron_up"></f7-icon>
                                    </template>
                                </f7-list-item>
                            </f7-list>
                        </f7-accordion-toggle>
                    </f7-block-title>
                    <f7-accordion-content style="height: auto">
                        <f7-list strong inset dividers media-list accordion-list class="transaction-info-list combination-list-content">
                            <f7-list-item chevron-center class="transaction-info" link="#">
                                <template #media>
                                    <div class="display-flex flex-direction-column transaction-date">
                                        <span class="transaction-day full-line flex-direction-column">{{ currentDayOfMonth }}</span>
                                        <span class="transaction-day-of-week full-line flex-direction-column">{{ currentDayOfWeek }}</span>
                                    </div>
                                </template>
                                <template #inner>
                                    <div class="display-flex no-padding-horizontal">
                                        <div class="item-media">
                                            <div class="transaction-icon display-flex align-items-center">
                                                <f7-icon f7="pencil_ellipsis_rectangle"></f7-icon>
                                            </div>
                                        </div>
                                        <div class="actual-item-inner">
                                            <div class="item-title-row">
                                                <div class="item-title">
                                                    <div class="transaction-category-name no-padding">
                                                        <span>{{ $t('Category Name') }}</span>
                                                    </div>
                                                </div>
                                                <div class="item-after">
                                                    <div class="transaction-amount">
                                                        <span>{{ getDisplayAmount('12345') }}</span>
                                                    </div>
                                                </div>
                                            </div>
                                            <div class="item-text">
                                                <div class="transaction-description">
                                                    <span>{{ $t('Description') }}</span>
                                                </div>
                                            </div>
                                            <div class="item-footer">
                                                <div class="transaction-tags">
                                                    <f7-chip media-bg-color="black" class="transaction-tag" :text="$t('Tag Title')">
                                                        <template #media>
                                                            <f7-icon f7="number"></f7-icon>
                                                        </template>
                                                    </f7-chip>
                                                </div>
                                                <div class="transaction-footer">
                                                    <span>{{ currentShortTime }}</span>
                                                    <span>·</span>
                                                    <span>{{ $t('Account Name') }}</span>
                                                </div>
                                            </div>
                                        </div>
                                    </div>
                                </template>
                            </f7-list-item>
                        </f7-list>
                    </f7-accordion-content>
                </f7-accordion-item>
            </f7-block>
        </f7-block>

        <f7-block class="fontsize-panel no-padding margin-bottom">
            <f7-block strong inset class="no-padding-bottom margin-bottom">
                <div class="full-line padding-bottom padding-top-half">
                    <div class="display-flex justify-content-space-between">
                        <div class="fontsize-minimum">A</div>
                        <div class="fontsize-maximum">A</div>
                        <div class="fontsize-default" :style="`left: calc(${100 / maxFontSizeType}% - 6px)`">{{ $t('Default') }}</div>
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
                        @range-change="fontSize = $event"
                    />
                </div>
            </f7-block>
        </f7-block>
    </f7-page>
</template>

<script>
import { mapStores } from 'pinia';
import { useSettingsStore } from '@/stores/setting.js';
import { useUserStore } from '@/stores/user.js';

import fontConstants from '@/consts/font.js';
import { getCurrentUnixTime, getDay, getDayOfWeekName } from '@/lib/datetime.js';
import { setAppFontSize, getFontSizePreviewClassName } from '@/lib/ui.mobile.js';

export default {
    props: [
        'f7router'
    ],
    data() {
        const settingsStore = useSettingsStore();

        return {
            currentTime: getCurrentUnixTime(),
            fontSize: settingsStore.appSettings.fontSize
        }
    },
    computed: {
        ...mapStores(useSettingsStore, useUserStore),
        minFontSizeType() {
            return 0;
        },
        maxFontSizeType() {
            return fontConstants.allFontSizeArray.length - 1;
        },
        fontSizePreviewClassName() {
            return getFontSizePreviewClassName(this.fontSize);
        },
        currentLongYearMonth() {
            return this.$locale.formatUnixTimeToLongYearMonth(this.userStore, this.currentTime);
        },
        currentDayOfMonth() {
            return getDay(this.currentTime);
        },
        currentDayOfWeek() {
            return this.$locale.getWeekdayShortName(getDayOfWeekName(this.currentTime));
        },
        currentShortTime() {
            return this.$locale.formatUnixTimeToShortTime(this.userStore, this.currentTime);
        }
    },
    methods: {
        setFontSize() {
            const router = this.f7router;

            if (this.fontSize !== this.settingsStore.appSettings.fontSize) {
                this.settingsStore.setFontSize(this.fontSize);
                setAppFontSize(this.fontSize);
            }

            router.back();
        },
        getFontSizeName() {
            return '';
        },
        getDisplayAmount(value) {
            return this.$locale.formatAmountWithCurrency(this.settingsStore, this.userStore, value);
        }
    }
}
</script>

<style>
.fontsize-panel {
    position: fixed;
    width: 100%;
    bottom: 0;
}

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
