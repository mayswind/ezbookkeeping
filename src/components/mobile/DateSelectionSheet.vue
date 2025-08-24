<template>
    <f7-sheet swipe-to-close swipe-handler=".swipe-handler" class="date-selection-sheet" style="height:auto"
              :opened="show" @sheet:open="onSheetOpen" @sheet:closed="onSheetClosed">
        <f7-toolbar>
            <div class="swipe-handler"></div>
            <div class="left">
                <f7-link :text="tt('Clear')" @click="clear"></f7-link>
            </div>
            <div class="right">
                <f7-link :text="tt('Done')" @click="confirm"></f7-link>
            </div>
        </f7-toolbar>
        <f7-page-content>
            <div class="block block-outline no-margin no-padding">
                <vue-date-picker inline auto-apply
                                 month-name-format="long"
                                 model-type="yyyy-MM-dd"
                                 six-weeks="center"
                                 class="justify-content-center"
                                 :enable-time-picker="false"
                                 :clearable="true"
                                 :dark="isDarkMode"
                                 :week-start="firstDayOfWeek"
                                 :year-range="yearRange"
                                 :day-names="dayNames"
                                 :year-first="isYearFirst"
                                 v-model="dateTime">
                    <template #month="{ text }">
                        {{ getMonthShortName(text) }}
                    </template>
                    <template #month-overlay-value="{ text }">
                        {{ getMonthShortName(text) }}
                    </template>
                </vue-date-picker>
            </div>
        </f7-page-content>
    </f7-sheet>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import { useEnvironmentsStore } from '@/stores/environment.ts';
import { useUserStore } from '@/stores/user.ts';

import { type TextualYearMonthDay, type WeekDayValue } from '@/core/datetime.ts';
import { arrangeArrayWithNewStartIndex } from '@/lib/common.ts';
import { getAllowedYearRange } from '@/lib/datetime.ts';

const props = defineProps<{
    modelValue?: TextualYearMonthDay;
    show: boolean;
}>();

const emit = defineEmits<{
    (e: 'update:modelValue', value: TextualYearMonthDay): void;
    (e: 'update:show', value: boolean): void;
}>();

const { tt, getAllMinWeekdayNames, getMonthShortName, isLongDateMonthAfterYear } = useI18n();

const environmentsStore = useEnvironmentsStore();
const userStore = useUserStore();

const yearRange = ref<number[]>(getAllowedYearRange());
const dateTime = ref<string>('');

const isDarkMode = computed<boolean>(() => environmentsStore.framework7DarkMode || false);
const firstDayOfWeek = computed<WeekDayValue>(() => userStore.currentUserFirstDayOfWeek);
const dayNames = computed<string[]>(() => arrangeArrayWithNewStartIndex(getAllMinWeekdayNames(), firstDayOfWeek.value));
const isYearFirst = computed<boolean>(() => isLongDateMonthAfterYear());

function clear(): void {
    dateTime.value = '';
    confirm();
}

function confirm(): void {
    emit('update:modelValue', dateTime.value as TextualYearMonthDay);
    emit('update:show', false);
}

function onSheetOpen(): void {
    dateTime.value = props.modelValue ?? '';
}

function onSheetClosed(): void {
    emit('update:show', false);
}
</script>

<style>
.date-selection-sheet .dp__menu {
    border: 0;
}
</style>
