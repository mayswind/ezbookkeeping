<template>
    <v-select
        ref="frequencySelect"
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
                    <v-list ref="frequencyTypeList" @keydown.capture="onListKeydown($event, 'left')">
                        <v-list-item :class="{ 'v-list-item--active text-primary': type.type === frequencyType }"
                                     :key="type.type" :title="type.displayName"
                                     v-for="type in allTransactionScheduledFrequencyTypes"
                                     @click="frequencyType = type.type">
                        </v-list-item>
                    </v-list>
                </div>
                <div class="schedule-frequency-value-container">
                    <v-list ref="frequencyValueList" :selectable="hasFrequencyValues"
                            :select-strategy="frequencyType === ScheduledTemplateFrequencyType.EveryNDays.type ? 'single-leaf' : 'classic'"
                            v-model:selected="frequencyValue" @keydown.capture="onListKeydown($event, 'right')">
                        <v-list-item :class="{ 'frequency-value-selected v-list-item--active text-primary': hasFrequencyValues && isFrequencyValueSelected(item.type) }"
                                     :key="item.type" :value="hasFrequencyValues ? item.type : undefined" :title="item.displayName"
                                     :tabindex="hasFrequencyValues ? -2 : undefined"
                                     v-for="item in frequencyItems">
                            <template #prepend="{ isSelected }" v-if="hasFrequencyValues">
                                <svg style="margin-inline: -6px 6px" width="24" height="24" viewBox="0 0 24 24" aria-hidden="true"
                                     :class="{ 'text-primary': isSelected }">
                                    <path :d="isSelected ? mdiCheckboxMarked : mdiCheckboxBlankOutline" fill="currentColor"/>
                                </svg>
                            </template>
                        </v-list-item>
                    </v-list>
                </div>
            </div>
        </template>
    </v-select>
</template>

<script setup lang="ts">
import { VSelect } from 'vuetify/components/VSelect';
import { VList } from 'vuetify/components/VList';

import { ref, computed, nextTick, useTemplateRef } from 'vue';

import { useI18n } from '@/locales/helpers.ts';
import { type CommonScheduleFrequencySelectionProps, useScheduleFrequencySelectionBase } from '@/components/base/ScheduleFrequencySelectionBase.ts';

import { useUserStore } from '@/stores/user.ts';

import { type TypeAndDisplayName } from '@/core/base.ts';
import { type WeekDayValue } from '@/core/datetime.ts';
import { ScheduledTemplateFrequencyType } from '@/core/template.ts';

import { sortNumbersArray } from '@/lib/common.ts';
import { scrollToSelectedItem } from '@/lib/ui/common.ts';
import { type TwoColumnListColumn, setTwoColumnItemFocusWhenKeydown } from '@/lib/ui/desktop.ts';

import {
    mdiCheckboxMarked,
    mdiCheckboxBlankOutline
} from '@mdi/js';

const props = defineProps<CommonScheduleFrequencySelectionProps>();
const emit = defineEmits<{
    (e: 'update:type', value: number): void;
    (e: 'update:modelValue', value: string): void;
}>();

const {
    tt,
    getMultiMonthAndDayLongNames,
    getMultiMonthdayShortNames,
    getMultiWeekdayLongNames
} = useI18n();
const {
    allTransactionScheduledFrequencyTypes,
    allWeekDays,
    allAvailableMonthDays,
    allAvailableMonthAndDays,
    allAvailableNDays,
    getFrequencyValues
} = useScheduleFrequencySelectionBase();

const userStore = useUserStore();

const dropdownMenu = useTemplateRef<HTMLElement>('dropdownMenu');
const frequencySelect = useTemplateRef<VSelect>('frequencySelect');
const frequencyTypeList = useTemplateRef<VList>('frequencyTypeList');
const frequencyValueList = useTemplateRef<VList>('frequencyValueList');

const menuState = ref<boolean>(false);

const firstDayOfWeek = computed<WeekDayValue>(() => userStore.currentUserFirstDayOfWeek);

const hasFrequencyValues = computed<boolean>(() => frequencyType.value !== ScheduledTemplateFrequencyType.Disabled.type && frequencyType.value !== ScheduledTemplateFrequencyType.Daily.type);

const frequencyItems = computed<TypeAndDisplayName[]>(() => {
    switch (frequencyType.value) {
        case ScheduledTemplateFrequencyType.Disabled.type:
            return [{ type: 0, displayName: tt('None') }];
        case ScheduledTemplateFrequencyType.Daily.type:
            return [{ type: 0, displayName: tt('Daily') }];
        case ScheduledTemplateFrequencyType.EveryNDays.type:
            return allAvailableNDays.value;
        case ScheduledTemplateFrequencyType.Weekly.type:
            return allWeekDays.value;
        case ScheduledTemplateFrequencyType.Monthly.type:
            return allAvailableMonthDays.value;
        case ScheduledTemplateFrequencyType.Yearly.type:
            return allAvailableMonthAndDays.value;
        default:
            return [];
    }
});

const frequencyType = computed<number>({
    get: () => props.type,
    set: (value: number) => {
        if (props.type !== value) {
            emit('update:type', value);

            if (value === ScheduledTemplateFrequencyType.Daily.type) {
                frequencyValue.value = [0];
            } else if (value === ScheduledTemplateFrequencyType.Weekly.type) {
                frequencyValue.value = [firstDayOfWeek.value];
            } else if (value === ScheduledTemplateFrequencyType.Monthly.type) {
                frequencyValue.value = [1];
            } else if (value === ScheduledTemplateFrequencyType.Yearly.type) {
                frequencyValue.value = [101];
            } else if (value === ScheduledTemplateFrequencyType.EveryNDays.type) {
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

function onListKeydown(event: KeyboardEvent, column: TwoColumnListColumn): void {
    setTwoColumnItemFocusWhenKeydown(event, column, {
        getList: column => {
            if (column === 'left') {
                return frequencyTypeList.value?.$el;
            } else if (column === 'right') {
                return frequencyValueList.value?.$el;
            } else {
                return null;
            }
        },
        getSelectedIndex: column => {
            if (column === 'left') {
                return allTransactionScheduledFrequencyTypes.value.findIndex(item => item.type === frequencyType.value);
            } else if (column === 'right') {
                return frequencyItems.value.findIndex(item => isFrequencyValueSelected(item.type));
            } else {
                return -1;
            }
        },
        onSelect: (column, index) => {
            const item = (column === 'left' ? allTransactionScheduledFrequencyTypes.value : frequencyItems.value)[index];

            if (!item) {
                return;
            }

            if (column === 'left') {
                frequencyType.value = item.type;
            } else if (hasFrequencyValues.value) {
                const value = item.type;

                if (!isFrequencyValueSelected(value)) {
                    if (frequencyType.value === ScheduledTemplateFrequencyType.EveryNDays.type) {
                        frequencyValue.value = [value];
                    } else {
                        frequencyValue.value = [...frequencyValue.value, value];
                    }
                }
            }
        },
        onLeave: () => {
            menuState.value = false;
            frequencySelect.value?.focus();
        }
    });
}

const displayFrequency = computed<string>(() => {
    if (frequencyType.value === ScheduledTemplateFrequencyType.Disabled.type) {
        return tt('Disabled');
    } else if (frequencyType.value === ScheduledTemplateFrequencyType.Daily.type) {
        return tt('Daily');
    } else if (frequencyType.value === ScheduledTemplateFrequencyType.EveryNDays.type) {
        if (frequencyValue.value.length) {
            return tt('format.misc.everyNDays', {
                n: frequencyValue.value[0]
            });
        } else {
            return tt('Every N Days');
        }
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
    } else if (frequencyType.value === ScheduledTemplateFrequencyType.Yearly.type) {
        if (frequencyValue.value.length) {
            return tt('format.misc.everyMultiDaysOfYear', {
                days: getMultiMonthAndDayLongNames(frequencyValue.value)
            });
        } else {
            return tt('Yearly');
        }
    } else {
        return '';
    }
});

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
