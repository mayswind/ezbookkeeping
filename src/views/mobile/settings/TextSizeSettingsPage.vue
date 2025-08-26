<template>
    <f7-page>
        <f7-navbar>
            <f7-nav-left :back-link="tt('Back')"></f7-nav-left>
            <f7-nav-title :title="tt('Text Size')"></f7-nav-title>
            <f7-nav-right>
                <f7-link :text="tt('Done')" @click="setFontSize"></f7-link>
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
                                            <span class="text-income">{{ `+${formatAmountToLocalizedNumeralsWithCurrency(12345)}` }}</span>
                                            <span class="text-expense">{{ `-${formatAmountToLocalizedNumeralsWithCurrency(67890)}` }}</span>
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
                                                        <span>{{ tt('Category Name') }}</span>
                                                    </div>
                                                </div>
                                                <div class="item-after">
                                                    <div class="transaction-amount">
                                                        <span>{{ formatAmountToLocalizedNumeralsWithCurrency(12345) }}</span>
                                                    </div>
                                                </div>
                                            </div>
                                            <div class="item-text">
                                                <div class="transaction-description">
                                                    <span>{{ tt('Description') }}</span>
                                                </div>
                                            </div>
                                            <div class="item-footer">
                                                <div class="transaction-tags">
                                                    <f7-chip media-text-color="var(--f7-chip-text-color)"
                                                             class="transaction-tag" :text="tt('Tag Title')">
                                                        <template #media>
                                                            <f7-icon f7="number"></f7-icon>
                                                        </template>
                                                    </f7-chip>
                                                </div>
                                                <div class="transaction-footer">
                                                    <span>{{ currentShortTime }}</span>
                                                    <span>·</span>
                                                    <span>{{ tt('Account Name') }}</span>
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
                        <div class="fontsize-default"
                             :style="textDirection === TextDirection.LTR ? `left: calc(${100 / FontSize.MaximumFontSize.type}% - 6px)` :  `right: calc(${100 / FontSize.MaximumFontSize.type}% - 6px)`">
                            {{ tt('Default') }}
                        </div>
                    </div>
                    <f7-range
                        :min="FontSize.MinimumFontSize.type"
                        :max="FontSize.MaximumFontSize.type"
                        :step="1"
                        :scale="true"
                        :scale-steps="FontSize.MaximumFontSize.type"
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

<script setup lang="ts">
import { ref, computed } from 'vue';
import type { Router } from 'framework7/types';

import { useI18n } from '@/locales/helpers.ts';

import { useSettingsStore } from '@/stores/setting.ts';

import { TextDirection } from '@/core/text.ts';
import { FontSize } from '@/core/font.ts';
import { parseDateTimeFromUnixTime, getCurrentUnixTime } from '@/lib/datetime.ts';
import { setAppFontSize, getFontSizePreviewClassName } from '@/lib/ui/mobile.ts';

const props = defineProps<{
    f7router: Router.Router;
}>();

const {
    tt,
    getCurrentLanguageTextDirection,
    getWeekdayShortName,
    formatUnixTimeToLongYearMonth,
    formatUnixTimeToShortTime,
    formatUnixTimeToDayOfMonth,
    formatAmountToLocalizedNumeralsWithCurrency
} = useI18n();

const settingsStore = useSettingsStore();

const currentUnixTime = ref<number>(getCurrentUnixTime());
const fontSize = ref<number>(settingsStore.appSettings.fontSize);

const textDirection = computed<string>(() => getCurrentLanguageTextDirection());
const fontSizePreviewClassName = computed<string>(() => getFontSizePreviewClassName(fontSize.value));
const currentLongYearMonth = computed<string>(() => formatUnixTimeToLongYearMonth(currentUnixTime.value));
const currentDayOfMonth = computed<string>(() => formatUnixTimeToDayOfMonth(currentUnixTime.value));
const currentDayOfWeek = computed<string>(() => getWeekdayShortName(parseDateTimeFromUnixTime(currentUnixTime.value).getWeekDay()));
const currentShortTime = computed<string>(() => formatUnixTimeToShortTime(currentUnixTime.value));

function getFontSizeName(): string {
    return '';
}

function setFontSize(): void {
    const router = props.f7router;

    if (fontSize.value !== settingsStore.appSettings.fontSize) {
        settingsStore.setFontSize(fontSize.value);
        setAppFontSize(fontSize.value);
    }

    router.back();
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
