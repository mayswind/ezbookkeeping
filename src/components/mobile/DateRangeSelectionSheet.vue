<template>
    <f7-sheet style="height:auto" :opened="show"
              @sheet:open="onSheetOpen" @sheet:closed="onSheetClosed">
        <f7-page-content>
            <div class="display-flex padding justify-content-space-between align-items-center">
                <div style="font-size: 18px" v-if="title"><b>{{ title }}</b></div>
            </div>
            <div class="padding-horizontal padding-bottom">
                <p class="no-margin-top margin-bottom-half" v-if="hint">{{ hint }}</p>
                <slot></slot>
                <f7-list no-hairlines inline-labels class="no-margin-top margin-bottom">
                    <f7-list-input
                        :label="$t('Begin Time')"
                        type="datetime-local"
                        class="date-range-sheet-time-item"
                        :value="currentMinDate"
                        @input="currentMinDate = $event.target.value"
                    >
                    </f7-list-input>

                    <f7-list-input
                        :label="$t('End Time')"
                        type="datetime-local"
                        class="date-range-sheet-time-item"
                        :value="currentMaxDate"
                        @input="currentMaxDate = $event.target.value"
                    >
                    </f7-list-input>
                </f7-list>
                <f7-button large fill
                           :class="{ 'disabled': !currentMinDate || !currentMaxDate }"
                           :text="$t('Continue')"
                           @click="confirm">
                </f7-button>
                <div class="margin-top text-align-center">
                    <f7-link  @click="cancel" :text="$t('Cancel')"></f7-link>
                </div>
            </div>
        </f7-page-content>
    </f7-sheet>
</template>

<script>
export default {
    props: [
        'minTime',
        'maxTime',
        'title',
        'hint',
        'show'
    ],
    data() {
        const self = this;
        let minDate = self.$utilities.getTodayFirstUnixTime();
        let maxDate = self.$utilities.getCurrentUnixTime();

        if (self.minTime) {
            minDate = self.minTime;
        }

        if (self.maxTime) {
            maxDate = self.maxTime;
        }

        return {
            currentMinDate: self.$utilities.formatUnixTime(minDate, 'YYYY-MM-DDTHH:mm'),
            currentMaxDate: self.$utilities.formatUnixTime(maxDate, 'YYYY-MM-DDTHH:mm')
        }
    },
    watch: {
        'currentMinDate': function (newValue) {
            if (!newValue) {
                this.currentMinDate = this.$utilities.formatUnixTime(this.$utilities.getCurrentUnixTime(), 'YYYY-MM-DDTHH:mm');
            }
        },
        'currentMaxDate': function (newValue) {
            if (!newValue) {
                this.currentMaxDate = this.$utilities.formatUnixTime(this.$utilities.getCurrentUnixTime(), 'YYYY-MM-DDTHH:mm');
            }
        }
    },
    methods: {
        onSheetOpen() {
            if (this.minTime) {
                this.currentMinDate = this.$utilities.formatUnixTime(this.minTime, 'YYYY-MM-DDTHH:mm');
            }

            if (this.maxTime) {
                this.currentMaxDate = this.$utilities.formatUnixTime(this.maxTime, 'YYYY-MM-DDTHH:mm');
            }
        },
        onSheetClosed() {
            this.$emit('update:show', false);
        },
        confirm() {
            if (!this.currentMinDate || !this.currentMaxDate) {
                return;
            }

            const minUnixTime = this.$utilities.getMinuteFirstUnixTime(this.currentMinDate);
            const maxUnixTime = this.$utilities.getMinuteLastUnixTime(this.currentMaxDate);

            if (minUnixTime < 0 || maxUnixTime < 0) {
                this.$toast('Date is too early');
                return;
            }

            this.$emit('dateRange:change', minUnixTime, maxUnixTime);
        },
        cancel() {
            this.$emit('update:show', false);
        }
    }
}
</script>

<style>
.date-range-sheet-time-item input[type="datetime-local"] {
    max-width: inherit;
}

.list .date-range-sheet-time-item > .item-content {
    padding-left: 0;
}
</style>
