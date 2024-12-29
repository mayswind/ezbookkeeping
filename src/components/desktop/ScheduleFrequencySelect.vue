<template>
    <v-select
        persistent-placeholder
        :readonly="readonly"
        :disabled="disabled"
        :label="label"
        :menu-props="{ 'content-class': 'schedule-frequency-select-menu' }"
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
                    <v-list v-if="frequencyType === allTemplateScheduledFrequencyTypes.Disabled.type">
                        <v-list-item :title="$t('None')"></v-list-item>
                    </v-list>
                    <v-list select-strategy="classic" v-model:selected="frequencyValue"
                            v-else-if="frequencyType === allTemplateScheduledFrequencyTypes.Weekly.type">
                        <v-list-item :key="weekDay.type" :value="weekDay.type" :title="weekDay.displayName"
                                     :class="{ 'frequency-value-selected v-list-item--active text-primary': isFrequencyValueSelected(weekDay.type) }"
                                     v-for="weekDay in allWeekDays">
                            <template #prepend="{ isActive }">
                                <v-checkbox density="compact" class="mr-1" :model-value="isActive"></v-checkbox>
                            </template>
                        </v-list-item>
                    </v-list>
                    <v-list select-strategy="classic" v-model:selected="frequencyValue"
                            v-else-if="frequencyType === allTemplateScheduledFrequencyTypes.Monthly.type">
                        <v-list-item :key="monthDay.day" :value="monthDay.day" :title="monthDay.displayName"
                                     :class="{ 'frequency-value-selected v-list-item--active text-primary': isFrequencyValueSelected(monthDay.day) }"
                                     v-for="monthDay in allAvailableMonthDays">
                            <template #prepend="{ isActive }">
                                <v-checkbox density="compact" class="mr-1" :model-value="isActive"></v-checkbox>
                            </template>
                        </v-list-item>
                    </v-list>
                </div>
            </div>
        </template>
    </v-select>
</template>

<script>
import { mapStores } from 'pinia';
import { useUserStore } from '@/stores/user.js';

import { ScheduledTemplateFrequencyType } from '@/core/template.ts';
import { sortNumbersArray } from '@/lib/common.ts';
import { scrollToSelectedItem } from '@/lib/ui/desktop.js';

export default {
    props: [
        'type',
        'modelValue',
        'disabled',
        'readonly',
        'label'
    ],
    emits: [
        'update:type',
        'update:modelValue'
    ],
    data() {
        return {
            menuState: false
        }
    },
    computed: {
        ...mapStores(useUserStore),
        allTransactionScheduledFrequencyTypes() {
            return this.$locale.getAllTransactionScheduledFrequencyTypes();
        },
        allTemplateScheduledFrequencyTypes() {
            return ScheduledTemplateFrequencyType.all();
        },
        allWeekDays() {
            return this.$locale.getAllWeekDays(this.firstDayOfWeek);
        },
        allAvailableMonthDays() {
            const allAvailableDays = [];

            for (let i = 1; i <= 28; i++) {
                allAvailableDays.push({
                    day: i,
                    displayName: this.$locale.getMonthdayShortName(i),
                });
            }

            return allAvailableDays;
        },
        firstDayOfWeek() {
            return this.userStore.currentUserFirstDayOfWeek;
        },
        frequencyType: {
            get: function () {
                return this.type;
            },
            set: function (value) {
                if (this.type !== value) {
                    this.$emit('update:type', value);

                    if (value === ScheduledTemplateFrequencyType.Weekly.type) {
                        this.frequencyValue = [this.firstDayOfWeek];
                    } else if (value === ScheduledTemplateFrequencyType.Monthly.type) {
                        this.frequencyValue = [1];
                    } else {
                        this.frequencyValue = [];
                    }
                }
            }
        },
        frequencyValue: {
            get: function () {
                const values = this.modelValue.split(',');
                const ret = [];

                for (let i = 0; i < values.length; i++) {
                    if (values[i]) {
                        ret.push(parseInt(values[i]));
                    }
                }

                return sortNumbersArray(ret);
            },
            set: function (value) {
                this.$emit('update:modelValue', sortNumbersArray(value).join(','));
            }
        },
        displayFrequency() {
            if (this.type === ScheduledTemplateFrequencyType.Disabled.type) {
                return this.$t('Disabled');
            } else if (this.type === ScheduledTemplateFrequencyType.Weekly.type) {
                if (this.frequencyValue.length) {
                    return this.$t('format.misc.everyMultiDaysOfWeek', {
                        days: this.$locale.getMultiWeekdayLongNames(this.frequencyValue, this.firstDayOfWeek)
                    });
                } else {
                    return this.$t('Weekly');
                }
            } else if (this.type === ScheduledTemplateFrequencyType.Monthly.type) {
                if (this.frequencyValue.length) {
                    return this.$t('format.misc.everyMultiDaysOfMonth', {
                        days: this.$locale.getMultiMonthdayShortNames(this.frequencyValue)
                    });
                } else {
                    return this.$t('Monthly');
                }
            } else {
                return '';
            }
        }
    },
    methods: {
        onMenuStateChanged(state) {
            const self = this;

            if (state) {
                self.$nextTick(() => {
                    if (self.$refs.dropdownMenu && self.$refs.dropdownMenu.parentElement) {
                        scrollToSelectedItem(self.$refs.dropdownMenu.parentElement, '.schedule-frequency-value-container', '.frequency-value-selected');
                    }
                });
            }
        },
        isFrequencyValueSelected(value) {
            for (let i = 0; i < this.frequencyValue.length; i++) {
                if (this.frequencyValue[i] === value) {
                    return true;
                }
            }

            return false;
        }
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
