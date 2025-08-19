<template>
    <f7-sheet swipe-to-close swipe-handler=".swipe-handler"
              style="height: auto" :opened="show" @sheet:open="onSheetOpen" @sheet:closed="onSheetClosed">
        <f7-toolbar>
            <div class="swipe-handler"></div>
            <div class="left">
                <f7-link sheet-close :text="tt('Cancel')"></f7-link>
            </div>
            <div class="right">
                <f7-link :text="tt('Done')" @click="save"></f7-link>
            </div>
        </f7-toolbar>
        <f7-page-content>
            <div class="grid grid-cols-2 grid-gap">
                <div>
                    <div class="schedule-frequency-type-container">
                        <f7-list dividers class="schedule-frequency-type-list no-margin-vertical">
                            <f7-list-item link="#" no-chevron
                                          :key="type.type"
                                          :title="type.displayName"
                                          v-for="type in allTransactionScheduledFrequencyTypes"
                                          @click="changeFrequencyType(type.type)">
                                <template #after>
                                    <f7-icon class="list-item-showing icon-with-direction" f7="chevron_right" v-if="currentFrequencyType === type.type"></f7-icon>
                                </template>
                            </f7-list-item>
                        </f7-list>
                    </div>
                </div>
                <div>
                    <div class="schedule-frequency-value-container">
                        <f7-list dividers class="schedule-frequency-value-list no-margin-vertical"
                                 v-if="currentFrequencyType === ScheduledTemplateFrequencyType.Disabled.type">
                            <f7-list-item :title="tt('None')"></f7-list-item>
                        </f7-list>
                        <f7-list dividers class="schedule-frequency-value-list no-margin-vertical"
                                 v-if="currentFrequencyType === ScheduledTemplateFrequencyType.Weekly.type">
                            <f7-list-item checkbox
                                          :class="isChecked(weekDay.type) ? 'list-item-selected' : ''"
                                          :key="weekDay.type"
                                          :value="weekDay.type"
                                          :checked="isChecked(weekDay.type)"
                                          :title="weekDay.displayName"
                                          v-for="weekDay in allWeekDays"
                                          @change="changeFrequencyValue">
                            </f7-list-item>
                        </f7-list>
                        <f7-list dividers class="schedule-frequency-value-list no-margin-vertical"
                                 v-if="currentFrequencyType === ScheduledTemplateFrequencyType.Monthly.type">
                            <f7-list-item checkbox
                                          :class="isChecked(monthDay.day) ? 'list-item-selected' : ''"
                                          :key="monthDay.day"
                                          :value="monthDay.day"
                                          :checked="isChecked(monthDay.day)"
                                          :title="monthDay.displayName"
                                          v-for="monthDay in allAvailableMonthDays"
                                          @change="changeFrequencyValue">
                            </f7-list-item>
                        </f7-list>
                    </div>
                </div>
            </div>
        </f7-page-content>
    </f7-sheet>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';

import { useI18n } from '@/locales/helpers.ts';
import { type CommonScheduleFrequencySelectionProps, useScheduleFrequencySelectionBase } from '@/components/base/ScheduleFrequencySelectionBase.ts';

import { useUserStore } from '@/stores/user.ts';

import { type WeekDayValue } from '@/core/datetime.ts';
import { ScheduledTemplateFrequencyType } from '@/core/template.ts';
import { sortNumbersArray } from '@/lib/common.ts';
import { type Framework7Dom, scrollToSelectedItem } from '@/lib/ui/mobile.ts';

interface MobileScheduleFrequencySelectionProps extends CommonScheduleFrequencySelectionProps {
    show: boolean;
}

const props = defineProps<MobileScheduleFrequencySelectionProps>();
const emit = defineEmits<{
    (e: 'update:type', value: number): void;
    (e: 'update:modelValue', value: string): void;
    (e: 'update:show', value: boolean): void;
}>();

const { tt } = useI18n();
const { allTransactionScheduledFrequencyTypes, allWeekDays, allAvailableMonthDays, getFrequencyValues } = useScheduleFrequencySelectionBase();

const userStore = useUserStore();

const currentFrequencyType = ref<number>(props.type);
const currentFrequencyValue = ref<number[]>(getFrequencyValues(props.modelValue));

const firstDayOfWeek = computed<WeekDayValue>(() => userStore.currentUserFirstDayOfWeek);

function isChecked(value: number): boolean {
    return currentFrequencyValue.value.indexOf(value) >= 0;
}

function changeFrequencyType(value: number): void {
    if (currentFrequencyType.value !== value) {
        currentFrequencyType.value = value;

        if (value === ScheduledTemplateFrequencyType.Weekly.type) {
            currentFrequencyValue.value = [firstDayOfWeek.value];
        } else if (value === ScheduledTemplateFrequencyType.Monthly.type) {
            currentFrequencyValue.value = [1];
        } else {
            currentFrequencyValue.value = [];
        }
    }
}

function changeFrequencyValue(e: Event): void {
    const value = parseInt((e.target as HTMLInputElement).value);

    if ((e.target as HTMLInputElement).checked) {
        for (let i = 0; i < currentFrequencyValue.value.length; i++) {
            if (currentFrequencyValue.value[i] === value) {
                return;
            }
        }

        currentFrequencyValue.value.push(value);
    } else {
        for (let i = 0; i < currentFrequencyValue.value.length; i++) {
            if (currentFrequencyValue.value[i] === value) {
                currentFrequencyValue.value.splice(i, 1);
                break;
            }
        }
    }
}

function save(): void {
    emit('update:type', currentFrequencyType.value);
    emit('update:modelValue', sortNumbersArray(currentFrequencyValue.value).join(','));
    emit('update:show', false);
}

function close(): void {
    emit('update:show', false);
}

function onSheetOpen(event: { $el: Framework7Dom }): void {
    currentFrequencyType.value = props.type;
    currentFrequencyValue.value = getFrequencyValues(props.modelValue);
    scrollToSelectedItem(event.$el, '.schedule-frequency-value-container', 'li.list-item-selected');
}

function onSheetClosed(): void {
    close();
}
</script>

<style>
.schedule-frequency-type-container, .schedule-frequency-value-container {
    height: 260px;
    overflow-y: auto;
}

.schedule-frequency-type-list.list .item-inner {
    padding-inline-end: 6px;
}

.schedule-frequency-value-list-list.list .item-content {
    padding-inline-start: 0;
}
</style>
