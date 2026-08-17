<template>
    <div class="projection-table-container">
        <table class="projection-table">
            <thead>
                <tr>
                    <th class="projection-row-label">{{ tt('Category') }}</th>
                    <th class="projection-amount" :class="{ 'projection-month-projected': month.kind !== ProjectionMonthKind.Actual }"
                        :key="month.yearMonth" v-for="month in data.months">
                        <span>{{ getMonthDisplayName(month) }}</span>
                        <v-tooltip activator="parent" v-if="month.kind !== ProjectionMonthKind.Actual">
                            {{ month.kind === ProjectionMonthKind.Mixed ? tt('Partially projected from scheduled transactions') : tt('Projected from scheduled transactions') }}
                        </v-tooltip>
                    </th>
                    <th class="projection-amount projection-total-column">{{ tt('Total') }}</th>
                </tr>
            </thead>

            <tbody :key="section.id" v-for="section in [data.income, data.expense]">
                <tr class="projection-section-row" @click="toggleSection(section.id)">
                    <th class="projection-row-label">
                        <v-icon size="20" :icon="expandedSections[section.id] ? mdiChevronDown : mdiChevronRight"/>
                        <span>{{ section.categoryType === CategoryType.Income ? tt('Income') : tt('Expense') }}</span>
                    </th>
                    <td :key="month.yearMonth" v-for="month in data.months"></td>
                    <td class="projection-total-column"></td>
                </tr>

                <template v-if="expandedSections[section.id]">
                    <template :key="category.id" v-for="category in section.categories">
                        <tr class="projection-category-row" @click="toggleCategory(category.id)">
                            <th class="projection-row-label">
                                <v-icon size="18" :icon="isCategoryExpanded(category.id) ? mdiChevronDown : mdiChevronRight"/>
                                <span>{{ category.name }}</span>
                            </th>
                            <td :key="month.yearMonth" v-for="month in data.months"></td>
                            <td class="projection-total-column"></td>
                        </tr>

                        <template v-if="isCategoryExpanded(category.id)">
                            <tr class="projection-subcategory-row" :key="subCategory.id"
                                v-for="subCategory in category.subCategories">
                                <th class="projection-row-label">{{ subCategory.name }}</th>
                                <td class="projection-amount" :key="index" v-for="(amount, index) in subCategory.monthlyAmounts">
                                    {{ formatAmount(amount) }}
                                </td>
                                <td class="projection-amount projection-total-column">{{ formatAmount(subCategory.totalAmount) }}</td>
                            </tr>
                        </template>

                        <tr class="projection-subtotal-row">
                            <th class="projection-row-label">{{ tt('Subtotal') }} {{ category.name }}</th>
                            <td class="projection-amount" :key="index" v-for="(amount, index) in category.monthlyAmounts">
                                {{ formatAmount(amount) }}
                            </td>
                            <td class="projection-amount projection-total-column">{{ formatAmount(category.totalAmount) }}</td>
                        </tr>
                    </template>
                </template>

                <tr class="projection-section-total-row">
                    <th class="projection-row-label">
                        {{ section.categoryType === CategoryType.Income ? tt('Total Income') : tt('Total Expense') }}
                    </th>
                    <td class="projection-amount" :key="index" v-for="(amount, index) in section.monthlyAmounts">
                        {{ formatAmount(amount) }}
                    </td>
                    <td class="projection-amount projection-total-column">{{ formatAmount(section.totalAmount) }}</td>
                </tr>
            </tbody>

            <tfoot>
                <tr class="projection-net-row">
                    <th class="projection-row-label">{{ tt('Net') }}</th>
                    <td class="projection-amount" :class="getAmountClass(amount)" :key="index"
                        v-for="(amount, index) in data.net.monthlyAmounts">
                        {{ formatAmount(amount) }}
                    </td>
                    <td class="projection-amount projection-total-column" :class="getAmountClass(data.net.totalAmount)">
                        {{ formatAmount(data.net.totalAmount) }}
                    </td>
                </tr>
                <tr class="projection-accumulated-row">
                    <th class="projection-row-label">{{ tt('Accumulated') }}</th>
                    <td class="projection-amount" :class="getAmountClass(amount)" :key="index"
                        v-for="(amount, index) in data.accumulated.monthlyAmounts">
                        {{ formatAmount(amount) }}
                    </td>
                    <td class="projection-amount projection-total-column" :class="getAmountClass(data.accumulated.totalAmount)">
                        {{ formatAmount(data.accumulated.totalAmount) }}
                    </td>
                </tr>
            </tfoot>
        </table>

        <p class="projection-empty text-body-2" v-if="!data.months.length">{{ tt('No data') }}</p>

        <v-alert class="mt-4" density="compact" type="warning" variant="tonal" v-if="data.hasUnconvertedAmounts">
            {{ tt('Some amounts cannot be converted to the default currency, so the totals are incomplete') }}
        </v-alert>
    </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import { CategoryType } from '@/core/category.ts';
import type { BigDecimal } from '@/core/numeral.ts';
import { type ProjectionMonth, type ProjectionTableData, ProjectionMonthKind } from '@/core/projection.ts';

import { getYearMonthFirstUnixTime, parseDateTimeFromUnixTime } from '@/lib/datetime.ts';

import {
    mdiChevronDown,
    mdiChevronRight
} from '@mdi/js';

defineProps<{
    data: ProjectionTableData;
}>();

const {
    tt,
    formatAmountToLocalizedNumeralsWithCurrency,
    formatDateTimeToGregorianLikeShortYearMonth
} = useI18n();

// Both levels of the accordion are independent, and their state belongs to the view rather than to
// the store: it must survive a reload of the data but not outlive the page.
const expandedSections = ref<Record<string, boolean>>({
    income: true,
    expense: true
});
const collapsedCategories = ref<Record<string, boolean>>({});

function toggleSection(sectionId: string): void {
    expandedSections.value[sectionId] = !expandedSections.value[sectionId];
}

// categories start expanded, so the state tracks the collapsed ones instead
function isCategoryExpanded(categoryId: string): boolean {
    return !collapsedCategories.value[categoryId];
}

function toggleCategory(categoryId: string): void {
    collapsedCategories.value[categoryId] = !collapsedCategories.value[categoryId];
}

function getMonthDisplayName(month: ProjectionMonth): string {
    const unixTime = getYearMonthFirstUnixTime({ year: month.year, month1base: month.month });
    return formatDateTimeToGregorianLikeShortYearMonth(parseDateTimeFromUnixTime(unixTime));
}

// every amount of a projection is already converted to the default currency of the user, which is
// what the formatter falls back to when no currency is given
function formatAmount(amount: BigDecimal): string {
    return formatAmountToLocalizedNumeralsWithCurrency(amount);
}

function getAmountClass(amount: BigDecimal): string {
    if (amount.isNegative()) {
        return 'text-expense';
    } else if (amount.isPositive()) {
        return 'text-income';
    }

    return '';
}
</script>

<style>
.projection-table-container {
    width: 100%;
    overflow-x: auto;
}

.projection-table {
    width: 100%;
    border-collapse: collapse;
    white-space: nowrap;
}

.projection-table th,
.projection-table td {
    padding: 0.35rem 0.75rem;
    font-size: 0.875rem;
}

.projection-table thead th {
    border-bottom: 1px solid rgba(var(--v-border-color), var(--v-border-opacity));
    font-weight: 500;
}

/* the label column stays in place while the months scroll sideways */
.projection-table .projection-row-label {
    position: sticky;
    left: 0;
    z-index: 1;
    text-align: start;
    font-weight: inherit;
    background: rgb(var(--v-theme-surface));
}

.projection-table .projection-amount {
    text-align: end;
    font-variant-numeric: tabular-nums;
}

.projection-table .projection-total-column {
    border-inline-start: 1px solid rgba(var(--v-border-color), var(--v-border-opacity));
    font-weight: 500;
}

/* months that are wholly or partly simulated are not comparable with the real ones */
.projection-table .projection-month-projected {
    font-style: italic;
    opacity: 0.85;
}

.projection-table .projection-section-row,
.projection-table .projection-category-row {
    cursor: pointer;
    user-select: none;
}

.projection-table .projection-section-row > .projection-row-label {
    padding-top: 0.75rem;
    text-transform: uppercase;
    font-weight: 600;
    letter-spacing: 0.02em;
}

.projection-table .projection-category-row > .projection-row-label {
    padding-inline-start: 1.25rem;
    font-weight: 500;
}

.projection-table .projection-subcategory-row > .projection-row-label {
    padding-inline-start: 3.25rem;
}

.projection-table .projection-subtotal-row {
    font-style: italic;
    opacity: 0.85;
}

.projection-table .projection-subtotal-row > .projection-row-label {
    padding-inline-start: 2.25rem;
}

.projection-table .projection-section-total-row,
.projection-table .projection-net-row,
.projection-table .projection-accumulated-row {
    font-weight: 600;
}

.projection-table .projection-section-total-row > * {
    border-top: 1px solid rgba(var(--v-border-color), var(--v-border-opacity));
}

.projection-table tfoot > tr > * {
    border-top: 2px solid rgba(var(--v-border-color), var(--v-border-opacity));
}

.projection-empty {
    padding: 1rem 0;
    opacity: 0.7;
}
</style>
