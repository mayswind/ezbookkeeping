<template>
    <f7-page>
        <f7-navbar :back-link="tt('Back')" :title="isEdit ? tt('Edit Budget') : tt('New Budget')"></f7-navbar>

        <f7-list strong inset dividers v-if="!loading">
            <!-- Category picker (expense categories only) -->
            <f7-list-item
                v-if="!isEdit"
                :title="tt('Category')"
                smart-select
                :smart-select-params="{ openIn: 'popup', searchbar: true, searchbarPlaceholder: tt('Search'), closeOnSelect: true }">
                <select v-model="selectedCategoryId">
                    <template v-for="primary in expenseCategories" :key="primary.id">
                        <option v-if="!primary.subCategories || primary.subCategories.length === 0" :value="primary.id">{{ primary.name }}</option>
                        <optgroup v-else :label="primary.name">
                            <option :value="primary.id">{{ primary.name }} ({{ tt('All') }})</option>
                            <option v-for="sub in primary.subCategories" :key="sub.id" :value="sub.id">{{ sub.name }}</option>
                        </optgroup>
                    </template>
                </select>
            </f7-list-item>
            <f7-list-item v-else :title="tt('Category')">
                <span>{{ selectedCategoryName }}</span>
            </f7-list-item>

            <!-- Amount input (display units, e.g. dollars) -->
            <f7-list-input
                :label="tt('Budget Amount')"
                type="number"
                step="0.01"
                inputmode="decimal"
                :placeholder="tt('Enter amount')"
                :value="amountInput"
                @input="amountInput = $event.target.value">
            </f7-list-input>

            <!-- Currency -->
            <f7-list-input
                :label="tt('Currency')"
                type="text"
                :placeholder="tt('Currency code')"
                :value="currency"
                @input="currency = $event.target.value.toUpperCase()">
            </f7-list-input>
        </f7-list>

        <f7-block v-if="!loading && isEdit">
            <f7-button fill color="red" @click="confirmDelete">{{ tt('Delete Budget') }}</f7-button>
        </f7-block>

        <f7-block v-if="!loading">
            <f7-button fill color="green" :class="{ 'disabled': !canSave }" @click="save">{{ isEdit ? tt('Save') : tt('Add') }}</f7-button>
        </f7-block>
    </f7-page>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import type { Router } from 'framework7/types';

import { useI18n } from '@/locales/helpers.ts';
import { useI18nUIComponents } from '@/lib/ui/mobile.ts';
import { useUserStore } from '@/stores/user.ts';
import { useTransactionCategoriesStore } from '@/stores/transactionCategory.ts';
import { useCategoryBudgetStore } from '@/stores/categoryBudget.ts';

import { CategoryType } from '@/core/category.ts';
import { CategoryBudgetLimit } from '@/models/category_budget_limit.ts';

const props = defineProps<{
    f7router: Router.Router;
    editBudgetId?: string;
    editCategoryId?: string;
    editAmount?: number;
    editCurrency?: string;
    monthStartDate?: number;
}>();

const { tt } = useI18n();
const { showToast, showConfirm } = useI18nUIComponents();

const userStore = useUserStore();
const categoriesStore = useTransactionCategoriesStore();
const budgetStore = useCategoryBudgetStore();

const loading = ref<boolean>(true);
const selectedCategoryId = ref<string>('');
const amountInput = ref<string>('');
const currency = ref<string>('');

const isEdit = computed(() => !!props.editBudgetId);

const expenseCategories = computed(() => categoriesStore.allTransactionCategories[CategoryType.Expense] || []);

const selectedCategoryName = computed(() => {
    if (!selectedCategoryId.value) {
        return '';
    }
    const cat = categoriesStore.allTransactionCategoriesMap[selectedCategoryId.value];
    return cat ? cat.name : '';
});

const canSave = computed(() => {
    if (!isEdit.value && !selectedCategoryId.value) {
        return false;
    }
    const amount = parseFloat(amountInput.value);
    return !isNaN(amount) && amount > 0 && currency.value.length > 0;
});

function firstOfMonth(date: Date): number {
    return new Date(date.getFullYear(), date.getMonth(), 1).getTime() / 1000;
}

function centsToDisplay(cents: number): string {
    return (cents / 100).toFixed(2);
}

function displayToCents(display: string): number {
    const value = parseFloat(display);

    if (isNaN(value)) {
        return 0;
    }

    return Math.round(value * 100);
}

function onPageInit(): void {
    loading.value = true;

    const initialCurrency = props.editCurrency || userStore.currentUserDefaultCurrency || 'USD';
    currency.value = initialCurrency;

    categoriesStore.loadAllCategories({ force: false }).then(() => {
        if (isEdit.value) {
            selectedCategoryId.value = props.editCategoryId || '';
            amountInput.value = centsToDisplay(props.editAmount || 0);
        } else {
            // Pre-select the first expense category if any.
            const cats = categoriesStore.allTransactionCategories[CategoryType.Expense];
            selectedCategoryId.value = (cats && cats.length > 0) ? cats[0]!.id : '';
        }
    }).catch(error => {
        if (!error.processed && !error.isUpToDate) {
            showToast(error.message || error);
        }

        if (isEdit.value) {
            selectedCategoryId.value = props.editCategoryId || '';
            amountInput.value = centsToDisplay(props.editAmount || 0);
        }
    }).finally(() => {
        loading.value = false;
    });
}

function save(): void {
    if (!canSave.value) {
        return;
    }

    const monthStart = props.monthStartDate || firstOfMonth(new Date());
    const amountCents = displayToCents(amountInput.value);
    const clientSessionId = `${Date.now()}`;

    const budget = props.editBudgetId
        ? CategoryBudgetLimit.of({ id: props.editBudgetId, categoryId: selectedCategoryId.value, startDate: monthStart, endDate: 0, amount: amountCents, currency: currency.value })
        : CategoryBudgetLimit.createNew(selectedCategoryId.value, monthStart, currency.value);

    budget.amount = amountCents;
    budget.currency = currency.value;

    budgetStore.saveBudget({
        budget,
        isEdit: isEdit.value,
        clientSessionId
    }).then(() => {
        showToast(isEdit.value ? 'Budget saved' : 'Budget added');
        budgetStore.loadOverview({ startDate: monthStart });
        props.f7router.back();
    }).catch(error => {
        if (!error.processed) {
            showToast(error.message || error);
        }
    });
}

function confirmDelete(): void {
    if (!props.editBudgetId) {
        return;
    }

    showConfirm('Delete this budget?', () => {
        budgetStore.deleteBudget({ budgetId: props.editBudgetId! }).then(() => {
            showToast('Budget deleted');
            budgetStore.loadOverview({ startDate: props.monthStartDate || firstOfMonth(new Date()) });
            props.f7router.back();
        }).catch(error => {
            if (!error.processed) {
                showToast(error.message || error);
            }
        });
    });
}

onMounted(onPageInit);
</script>
