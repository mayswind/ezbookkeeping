<template>
    <f7-sheet swipe-to-close swipe-handler=".swipe-handler"
              style="height: auto" :opened="show" @sheet:open="onSheetOpen" @sheet:closed="onSheetClosed">
        <f7-toolbar>
            <div class="swipe-handler"></div>
            <div class="left">
                <f7-link sheet-close :text="$t('Cancel')"></f7-link>
            </div>
            <div class="right">
                <f7-link :text="$t('Done')" @click="save"></f7-link>
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
                                    <f7-icon class="list-item-showing" f7="chevron_right" v-if="currentFrequencyType === type.type"></f7-icon>
                                </template>
                            </f7-list-item>
                        </f7-list>
                    </div>
                </div>
                <div>
                    <div class="schedule-frequency-value-container">
                        <f7-list dividers class="schedule-frequency-value-list no-margin-vertical"
                                 v-if="currentFrequencyType === allTemplateScheduledFrequencyTypes.Disabled.type">
                            <f7-list-item :title="$t('None')"></f7-list-item>
                        </f7-list>
                        <f7-list dividers class="schedule-frequency-value-list no-margin-vertical"
                                 v-if="currentFrequencyType === allTemplateScheduledFrequencyTypes.Weekly.type">
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
                                 v-if="currentFrequencyType === allTemplateScheduledFrequencyTypes.Monthly.type">
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

<script>
import { mapStores } from 'pinia';
import { useUserStore } from '@/stores/user.js';

import { ScheduledTemplateFrequencyType } from '@/core/template.ts';
import { sortNumbersArray } from '@/lib/common.ts';
import { scrollToSelectedItem } from '@/lib/ui/mobile.js';

export default {
    props: [
        'type',
        'modelValue',
        'disabled',
        'readonly',
        'label',
        'show'
    ],
    emits: [
        'update:type',
        'update:modelValue',
        'update:show'
    ],
    data() {
        const self = this;

        return {
            currentFrequencyType: self.type,
            currentFrequencyValue: self.getFrequencyValues(self.modelValue)
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
        }
    },
    methods: {
        onSheetOpen(event) {
            this.currentFrequencyType = this.type;
            this.currentFrequencyValue = this.getFrequencyValues(this.modelValue);
            scrollToSelectedItem(event.$el, '.schedule-frequency-value-container', 'li.list-item-selected');
        },
        onSheetClosed() {
            this.close();
        },
        changeFrequencyType(value) {
            if (this.currentFrequencyType !== value) {
                this.currentFrequencyType = value;

                if (value === ScheduledTemplateFrequencyType.Weekly.type) {
                    this.currentFrequencyValue = [this.firstDayOfWeek];
                } else if (value === ScheduledTemplateFrequencyType.Monthly.type) {
                    this.currentFrequencyValue = [1];
                } else {
                    this.currentFrequencyValue = [];
                }
            }
        },
        changeFrequencyValue(e) {
            const value = parseInt(e.target.value);

            if (e.target.checked) {
                for (let i = 0; i < this.currentFrequencyValue.length; i++) {
                    if (this.currentFrequencyValue[i] === value) {
                        return;
                    }
                }

                this.currentFrequencyValue.push(value);
            } else {
                for (let i = 0; i < this.currentFrequencyValue.length; i++) {
                    if (this.currentFrequencyValue[i] === value) {
                        this.currentFrequencyValue.splice(i, 1);
                        break;
                    }
                }
            }
        },
        save() {
            this.$emit('update:type', this.currentFrequencyType);
            this.$emit('update:modelValue', sortNumbersArray(this.currentFrequencyValue).join(','));
            this.$emit('update:show', false);
        },
        close() {
            this.$emit('update:show', false);
        },
        isChecked(value) {
            for (let i = 0; i < this.currentFrequencyValue.length; i++) {
                if (this.currentFrequencyValue[i] === value) {
                    return true;
                }
            }

            return false;
        },
        getFrequencyValues(value) {
            const values = value.split(',');
            const ret = [];

            for (let i = 0; i < values.length; i++) {
                if (values[i]) {
                    ret.push(parseInt(values[i]));
                }
            }

            return sortNumbersArray(ret);
        }
    }
}
</script>

<style>
.schedule-frequency-type-container, .schedule-frequency-value-container {
    height: 260px;
    overflow-y: auto;
}

.schedule-frequency-type-list.list .item-inner {
    padding-right: 6px;
}

.schedule-frequency-value-list-list.list .item-content {
    padding-left: 0;
}
</style>
