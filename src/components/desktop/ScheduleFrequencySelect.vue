<template>
    <v-select
        persistent-placeholder
        :readonly="readonly"
        :disabled="disabled"
        :label="label"
        :menu-props="{ contentClass: 'schedule-frequency-select-menu' }"
        v-model="frequencyType"
        v-model:menu="menuState"
        @update:menu="onMenuStateChanged"
    >
        <template #selection>
            <span class="text-truncate cursor-pointer">{{ displayFrequency }}</span>
        </template>

        <template #no-data>
            <div ref="dropdownMenu" class="schedule-frequency-container">
                <div class="schedule-frequency-type-container">
                    <v-list>
                        <v-list-item :class="{ 'v-list-item--active text-primary': type.type === frequencyType }"
                                     :key="type.type" :title="type.displayName"
                                     v-for="type in allTransactionScheduledFrequencyTypes"
                                     @click="frequencyType = type.type">
                        </v-list-item>
                    </v-list>
                </div>
                <div class="schedule-frequency-value-container">
                    <v-list v-if="frequencyType === ScheduledTemplateFrequencyType.Disabled.type">
                        <v-list-item :title="tt('None')"></v-list-item>
                    </v-list>
                    <v-list select-strategy="classic" v-model:selected="frequencyValue"
                            v-else-if="frequencyType === ScheduledTemplateFrequencyType.Weekly.type">
                        <v-list-item :key="weekDay.type" :value="weekDay.type" :title="weekDay.displayName"
                                     :class="{ 'frequency-value-selected v-list-item--active text-primary': isFrequencyValueSelected(weekDay.type) }"
                                     v-for="weekDay in allWeekDays">
                            <template #prepend="{ isActive }">
                                <v-checkbox density="compact" class="me-1" :model-value="isActive"
                                            @update:model-value="updateFrequencyValue(weekDay.type, $event)"></v-checkbox>
                            </template>
                        </v-list-item>
                    </v-list>
                    <v-list select-strategy="classic" v-model:selected="frequencyValue"
                            v-else-if="frequencyType === ScheduledTemplateFrequencyType.Monthly.type">
                        <v-list-item :key="monthDay.day" :value="monthDay.day" :title="monthDay.displayName"
                                     :class="{ 'frequency-value-selected v-list-item--active text-primary': isFrequencyValueSelected(monthDay.day) }"
                                     v-for="monthDay in allAvailableMonthDays">
                            <template #prepend="{ isActive }">
                                <v-checkbox density="compact" class="me-1" :model-value="isActive"
                                            @update:model-value="updateFrequencyValue(monthDay.day, $event)"></v-checkbox>
                            </template>
                        </v-list-item>
                    </v-list>
                </div>
            </div>
        </template>
    </v-select>
</template>

<script setup lang="ts">
import { ref, computed, nextTick, useTemplateRef } from 'vue';

import { useI18n } from '@/locales/helpers.ts';
import { type CommonScheduleFrequencySelectionProps, useScheduleFrequencySelectionBase } from '@/components/base/ScheduleFrequencySelectionBase.ts';

import { useUserStore } from '@/stores/user.ts';

import { type WeekDayValue } from '@/core/datetime.ts';
import { ScheduledTemplateFrequencyType } from '@/core/template.ts';
import { sortNumbersArray } from '@/lib/common.ts';
import { scrollToSelectedItem } from '@/lib/ui/common.ts';

const props = defineProps<CommonScheduleFrequencySelectionProps>();
const emit = defineEmits<{
    (e: 'update:type', value: number): void;
    (e: 'update:modelValue', value: string): void;
}>();

const { tt, getMultiMonthdayShortNames, getMultiWeekdayLongNames } = useI18n();
const { allTransactionScheduledFrequencyTypes, allWeekDays, allAvailableMonthDays, getFrequencyValues } = useScheduleFrequencySelectionBase();

const userStore = useUserStore();

const dropdownMenu = useTemplateRef<HTMLElement>('dropdownMenu');

const menuState = ref<boolean>(false);

const firstDayOfWeek = computed<WeekDayValue>(() => userStore.currentUserFirstDayOfWeek);

const frequencyType = computed<number>({
    get: () => props.type,
    set: (value: number) => {
        if (props.type !== value) {
            emit('update:type', value);

            if (value === ScheduledTemplateFrequencyType.Weekly.type) {
                frequencyValue.value = [firstDayOfWeek.value];
            } else if (value === ScheduledTemplateFrequencyType.Monthly.type) {
                frequencyValue.value = [1];
            } else {
                frequencyValue.value = [];
            }
        }
    }
});

const frequencyValue = computed<number[]>({
    get: () => getFrequencyValues(props.modelValue),
    set: (value: number[]) => {
        emit('update:modelValue', sortNumbersArray(value).join(','));
    }
});

const displayFrequency = computed<string>(() => {
    if (frequencyType.value === ScheduledTemplateFrequencyType.Disabled.type) {
        return tt('Disabled');
    } else if (frequencyType.value === ScheduledTemplateFrequencyType.Weekly.type) {
        if (frequencyValue.value.length) {
            return tt('format.misc.everyMultiDaysOfWeek', {
                days: getMultiWeekdayLongNames(frequencyValue.value, firstDayOfWeek.value)
            });
        } else {
            return tt('Weekly');
        }
    } else if (frequencyType.value === ScheduledTemplateFrequencyType.Monthly.type) {
        if (frequencyValue.value.length) {
            return tt('format.misc.everyMultiDaysOfMonth', {
                days: getMultiMonthdayShortNames(frequencyValue.value)
            });
        } else {
            return tt('Monthly');
        }
    } else {
        return '';
    }
});

function updateFrequencyValue(value: number, selected: boolean | null): void {
    const currentFrequencyValues = frequencyValue.value;
    const newFrequencyValues: number[] = [];

    for (const currentFrequencyValue of currentFrequencyValues) {
        if (currentFrequencyValue !== value || selected) {
            newFrequencyValues.push(currentFrequencyValue);
        }
    }

    if (selected) {
        newFrequencyValues.push(value);
    }

    frequencyValue.value = sortNumbersArray(newFrequencyValues);
}

function isFrequencyValueSelected(currentValue: number): boolean {
    return frequencyValue.value.indexOf(currentValue) >= 0;
}

function onMenuStateChanged(state: boolean): void {
    if (state) {
        nextTick(() => {
            if (dropdownMenu.value && dropdownMenu.value.parentElement) {
                scrollToSelectedItem(dropdownMenu.value.parentElement, '.schedule-frequency-value-container', '.schedule-frequency-value-container', '.frequency-value-selected');
            }
        });
    }
}
</script>

<style>
.schedule-frequency-select-menu {
    max-height: inherit !important;
}

.schedule-frequency-select-menu > .v-list {
    padding: 0;
}

.schedule-frequency-select-menu .schedule-frequency-container {
    width: 100%;
    display: flex;
}

.schedule-frequency-select-menu .schedule-frequency-type-container,
.schedule-frequency-select-menu .schedule-frequency-value-container {
    width: 100%;
    max-height: 310px;
    overflow-y: scroll;
}
</style>
