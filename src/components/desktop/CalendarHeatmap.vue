<template>
    <div class="calendar-heatmap">
        <div class="calendar-heatmap-header" v-if="showTitle">
            <h3>{{ tt('Transaction Activity - Past 365 Days') }}</h3>
        </div>

        <div class="calendar-heatmap-container">
            <!-- CSS Grid based heatmap -->
            <div class="heatmap-grid">
                <!-- Month labels -->
                <div
                    v-for="position in gridData.monthPositions"
                    :key="`month-${position.month}`"
                    class="month-label"
                    :style="{
                        gridRow: 1,
                        gridColumn: position.column
                    }"
                >
                    {{ tt(position.name) }}
                </div>

                <!-- Weekday labels -->
                <div
                    v-for="(day, index) in weekLabels"
                    :key="`weekday-${index}`"
                    class="weekday-label"
                    :class="{ 'show': index % 2 === 1 }"
                    :style="{
                        gridRow: index + 2,
                        gridColumn: 1
                    }"
                    v-if="showWeekdays"
                >
                    {{ tt(day) }}
                </div>

                <!-- Heatmap days -->
                <div
                    v-for="item in gridData.gridItems"
                    :key="item.date"
                    class="heatmap-day"
                    :class="{
                        'has-data': item.value > 0,
                        'clickable': enableClick,
                        'future-date': isFutureDate(item.date)
                    }"
                    :style="{
                        gridRow: item.gridRow,
                        gridColumn: item.gridColumn,
                        backgroundColor: getColorForLevel(item.level, item.isIncome)
                    }"
                    @click="enableClick ? clickDay(item) : null"
                >
                    <v-tooltip location="top" :text="getTooltipText(item)">
                        <template v-slot:activator="{ props }">
                            <div v-bind="props" class="day-content"></div>
                        </template>
                    </v-tooltip>
                </div>
            </div>

            <!-- Legend -->
            <div class="heatmap-legend">
                <!-- Color scheme explanation -->
                <div class="legend-explanation">
                    <div class="legend-item">
                        <div class="legend-color-sample" :style="{ backgroundColor: getColorForLevel(2, true) }"></div>
                        <span class="legend-text">{{ tt('Income') }}</span>
                    </div>
                    <div class="legend-item">
                        <div class="legend-color-sample" :style="{ backgroundColor: getColorForLevel(2, false) }"></div>
                        <span class="legend-text">{{ tt('Expense') }}</span>
                    </div>
                </div>

                <!-- Intensity scale -->
                <div class="legend-scale">
                    <span class="legend-label">{{ tt('Less') }}</span>
                    <div class="legend-colors">
                        <div
                            v-for="level in [0, 1, 2, 3, 4]"
                            :key="level"
                            class="legend-color"
                            :style="{ backgroundColor: getColorForLevel(level) }"
                        ></div>
                    </div>
                    <span class="legend-label">{{ tt('More') }}</span>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { useI18n } from '@/locales/helpers.ts';
import { type CommonCalendarHeatmapProps, type CalendarHeatmapDataItem, useCalendarHeatmapBase } from '@/components/base/CalendarHeatmapBase.ts';

interface DesktopCalendarHeatmapProps extends CommonCalendarHeatmapProps {
    showTitle?: boolean;
    enableClick?: boolean;
}

const props = withDefaults(defineProps<DesktopCalendarHeatmapProps>(), {
    showWeekdays: true,
    showMonthLabels: true,
    showTitle: true,
    enableClick: false
});

const emit = defineEmits<{
    (e: 'dayClick', day: CalendarHeatmapDataItem): void;
}>();

const { tt } = useI18n();

const {
    gridData,
    weekLabels,
    getColorForLevel,
    getTooltipText
} = useCalendarHeatmapBase(props);

function clickDay(day: CalendarHeatmapDataItem) {
    if (props.enableClick) {
        emit('dayClick', day);
    }
}

function isFutureDate(dateStr: string): boolean {
    const date = new Date(dateStr);
    const today = new Date();
    today.setHours(23, 59, 59, 999);
    return date > today;
}
</script>

<style scoped>
.calendar-heatmap {
    width: 100%;
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', 'Helvetica', 'Arial', sans-serif;
    font-size: 12px;
}

.calendar-heatmap-header h3 {
    margin: 0 0 16px 0;
    font-size: 1.25rem;
    font-weight: 500;
    line-height: 1.6;
    letter-spacing: 0.0125em;
    color: rgba(var(--v-theme-on-surface), var(--v-high-emphasis-opacity));
}

.calendar-heatmap-container {
    overflow-x: auto;
}

.heatmap-grid {
    display: grid;
    grid-template-columns: repeat(54, 1fr); /* ~53 weeks + 1 for labels */
    grid-template-rows: repeat(9, 1fr); /* 1 for months + 7 for weekdays + 1 spare */
    gap: min(2px, 0.3vw);
    padding: 16px 0;
    width: 100%;
    height: 200px; /* Fixed height to maintain aspect ratio */
}

/* Mobile responsiveness */
@media (max-width: 768px) {
    .heatmap-grid {
        height: 150px;
        gap: 1px;
        padding: 12px 0;
    }

    .month-label {
        font-size: 8px;
    }

    .weekday-label {
        font-size: 7px;
    }

    .heatmap-legend {
        flex-direction: column;
        gap: 8px;
        align-items: flex-start;
    }

    .legend-explanation {
        order: 2;
    }

    .legend-scale {
        order: 1;
    }
}

.month-label {
    font-size: 10px;
    color: var(--v-theme-on-surface-variant);
    white-space: nowrap;
    text-align: left;
    line-height: 11px;
}

.weekday-label {
    font-size: 9px;
    color: var(--v-theme-on-surface-variant);
    line-height: 11px;
    text-align: right;
    padding-right: 4px;
    opacity: 0;
}

.weekday-label.show {
    opacity: 1;
}

.heatmap-day {
    width: 100%;
    height: 100%;
    min-width: 8px;
    min-height: 8px;
    border-radius: 2px;
    cursor: default;
    transition: all 0.2s ease;
    border: 1px solid rgba(27, 31, 35, 0.06);
}

.heatmap-day.clickable {
    cursor: pointer;
}

.heatmap-day.clickable:hover {
    border-color: rgba(27, 31, 35, 0.15);
    transform: scale(1.1);
}

.heatmap-day.future-date {
    opacity: 0.3;
    cursor: not-allowed !important;
}

.heatmap-day.future-date.clickable {
    cursor: not-allowed !important;
}

.heatmap-day.future-date:hover {
    transform: none !important;
    border-color: rgba(27, 31, 35, 0.06) !important;
}

.day-content {
    width: 100%;
    height: 100%;
    border-radius: inherit;
}

.heatmap-legend {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-top: 8px;
    font-size: 11px;
    color: var(--v-theme-on-surface-variant);
}

.legend-explanation {
    display: flex;
    gap: 12px;
}

.legend-item {
    display: flex;
    align-items: center;
    gap: 4px;
}

.legend-color-sample {
    width: 12px;
    height: 12px;
    border-radius: 2px;
    border: 1px solid rgba(27, 31, 35, 0.06);
}

.legend-text {
    font-size: 11px;
}

.legend-scale {
    display: flex;
    align-items: center;
    gap: 4px;
}

.legend-colors {
    display: flex;
    gap: 2px;
}

.legend-color {
    width: 12px;
    height: 12px;
    border-radius: 2px;
    border: 1px solid rgba(27, 31, 35, 0.06);
}

.legend-label {
    font-size: 10px;
}
</style>