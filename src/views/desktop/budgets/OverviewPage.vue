<template>
    <v-row class="match-height">
        <v-col cols="12">
            <v-card>
                <v-card-item>
                    <div class="d-flex align-center justify-space-between">
                        <div class="d-flex align-center">
                            <v-btn variant="text" :icon="mdiChevronLeft" :disabled="loading" @click="changeMonth(-1)"></v-btn>
                            <span class="text-h6 mx-2 cursor-pointer" @click="resetToCurrentMonth">{{ formatMonth(currentMonthStartDate) }}</span>
                            <v-btn variant="text" :icon="mdiChevronRight" :disabled="loading" @click="changeMonth(1)"></v-btn>
                        </div>
                        <v-btn color="primary" variant="tonal" :prepend-icon="mdiPlus" :disabled="loading || expenseCategories.length === 0" @click="openAddDialog">
                            {{ tt('New Budget') }}
                        </v-btn>
                    </div>
                </v-card-item>

                <v-divider />

                <!-- Summary -->
                <v-card-text v-if="overview && overview.items.length > 0">
                    <v-row>
                        <v-col cols="12" md="4">
                            <v-sheet class="pa-4 rounded budget-stat-card">
                                <div class="text-caption text-medium-emphasis">{{ tt('Total Budget') }}</div>
                                <div class="text-h6">{{ displayAmount(overview.totalLimit, defaultCurrency) }}</div>
                            </v-sheet>
                        </v-col>
                        <v-col cols="12" md="4">
                            <v-sheet class="pa-4 rounded budget-stat-card">
                                <div class="text-caption text-medium-emphasis">{{ tt('Total Spent') }}</div>
                                <div class="text-h6 error--text">{{ displayAmount(overview.totalActualExpenseAmount, defaultCurrency) }}</div>
                            </v-sheet>
                        </v-col>
                        <v-col cols="12" md="4">
                            <v-sheet class="pa-4 rounded budget-stat-card">
                                <div class="text-caption text-medium-emphasis">{{ tt('Total Available') }}</div>
                                <div class="text-h6" :class="{ 'error--text': overview.totalAvailableAmount < 0 }">
                                    {{ displayAmount(overview.totalAvailableAmount, defaultCurrency) }}
                                </div>
                            </v-sheet>
                        </v-col>
                    </v-row>
                </v-card-text>

                <v-card-text>
                    <v-data-table
                        :headers="tableHeaders"
                        :items="overview ? overview.items : []"
                        :loading="loading"
                        :no-data-text="tt('No budgets for this month. Click New Budget to set one.')"
                        hide-default-footer
                        :items-per-page="-1">
                        <template v-slot:item.categoryName="{ item }">
                            <span>{{ item.categoryName || tt('(Unknown Category)') }}</span>
                        </template>
                        <template v-slot:item.progress="{ item }">
                            <v-progress-linear
                                :model-value="progressPercent(item)"
                                :color="item.availableAmount < 0 ? 'error' : (progressPercent(item) >= 90 ? 'warning' : 'success')"
                                height="8"
                                rounded />
                        </template>
                        <template v-slot:item.actualExpenseAmount="{ item }">
                            <span class="error--text font-weight-medium">{{ displayAmount(item.actualExpenseAmount, item.currency) }}</span>
                        </template>
                        <template v-slot:item.amount="{ item }">
                            <span>{{ displayAmount(item.amount, item.currency) }}</span>
                        </template>
                        <template v-slot:item.availableAmount="{ item }">
                            <span :class="{ 'error--text': item.availableAmount < 0 }">{{ displayAmount(item.availableAmount, item.currency) }}</span>
                        </template>
                        <template v-slot:item.actions="{ item }">
                            <v-btn variant="text" size="small" :icon="mdiPencil" @click="openEditDialog(item)"></v-btn>
                            <v-btn variant="text" size="small" color="error" :icon="mdiDelete" @click="confirmDelete(item)"></v-btn>
                        </template>
                    </v-data-table>
                </v-card-text>
            </v-card>
        </v-col>

        <snack-bar ref="snackbar" />
    </v-row>

    <!-- Add / Edit dialog -->
    <v-dialog v-model="dialog.show" max-width="480" persistent>
        <v-card>
            <v-card-title>{{ dialog.isEdit ? tt('Edit Budget') : tt('New Budget') }}</v-card-title>
            <v-card-text>
                <v-select
                    v-if="!dialog.isEdit"
                    v-model="dialog.categoryId"
                    :items="selectableCategories"
                    item-title="title"
                    item-value="id"
                    :label="tt('Category')"
                    :disabled="loading"
                    density="comfortable"
                    class="mb-3" />
                <v-text-field
                    v-else
                    :model-value="dialog.categoryName"
                    :label="tt('Category')"
                    disabled
                    density="comfortable"
                    class="mb-3" />
                <v-text-field
                    v-model.number="dialog.amountInput"
                    type="number"
                    step="0.01"
                    :label="tt('Budget Amount')"
                    density="comfortable"
                    class="mb-3" />
                <v-text-field
                    v-model="dialog.currency"
                    :label="tt('Currency')"
                    density="comfortable" />
            </v-card-text>
            <v-card-actions>
                <v-spacer />
                <v-btn variant="text" @click="dialog.show = false">{{ tt('Cancel') }}</v-btn>
                <v-btn color="primary" variant="tonal" :disabled="!canSaveDialog" @click="saveDialog">{{ dialog.isEdit ? tt('Save') : tt('Add') }}</v-btn>
            </v-card-actions>
        </v-card>
    </v-dialog>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { mdiChevronLeft, mdiChevronRight, mdiPlus, mdiPencil, mdiDelete } from '@mdi/js';

import { useI18n } from '@/locales/helpers.ts';
import { useUserStore } from '@/stores/user.ts';
import { useTransactionCategoriesStore } from '@/stores/transactionCategory.ts';
import { useCategoryBudgetStore } from '@/stores/categoryBudget.ts';
import { parseBigDecimal } from '@/lib/numeral.ts';

import { CategoryType } from '@/core/category.ts';
import SnackBar from '@/components/desktop/SnackBar.vue';

import type { CategoryBudgetOverviewItem } from '@/models/category_budget_limit.ts';
import { CategoryBudgetLimit } from '@/models/category_budget_limit.ts';

type SnackBarType = InstanceType<typeof SnackBar>;

const { tt, formatAmountToLocalizedNumeralsWithCurrency } = useI18n();
const userStore = useUserStore();
const categoriesStore = useTransactionCategoriesStore();
const budgetStore = useCategoryBudgetStore();

const snackbar = ref<SnackBarType | null>(null);

const loading = ref<boolean>(true);
const currentMonthStartDate = ref<number>(firstOfMonth(new Date()));

function firstOfMonth(date: Date): number {
    return Math.floor(new Date(date.getFullYear(), date.getMonth(), 1).getTime() / 1000);
}
function addMonthsToDate(unixSeconds: number, months: number): number {
    const d = new Date(unixSeconds * 1000);
    return Math.floor(new Date(d.getFullYear(), d.getMonth() + months, 1).getTime() / 1000);
}

const overview = computed(() => budgetStore.currentOverview);
const defaultCurrency = computed(() => userStore.currentUserDefaultCurrency || 'USD');

const expenseCategories = computed(() => categoriesStore.allTransactionCategories[CategoryType.Expense] || []);

const selectableCategories = computed(() => {
    const result: { id: string; title: string }[] = [];

    for (const primary of expenseCategories.value) {
        if (!primary.subCategories || primary.subCategories.length === 0) {
            result.push({ id: primary.id, title: primary.name });
        } else {
            result.push({ id: primary.id, title: `${primary.name} (${tt('All')})` });

            for (const sub of primary.subCategories) {
                result.push({ id: sub.id, title: `${primary.name} / ${sub.name}` });
            }
        }
    }

    return result;
});

const tableHeaders = computed(() => [
    { title: tt('Category'), key: 'categoryName', sortable: true },
    { title: tt('Progress'), key: 'progress', sortable: false },
    { title: tt('Spent'), key: 'actualExpenseAmount', sortable: true, align: 'end' as const },
    { title: tt('Budget'), key: 'amount', sortable: true, align: 'end' as const },
    { title: tt('Available'), key: 'availableAmount', sortable: true, align: 'end' as const },
    { title: tt('Actions'), key: 'actions', sortable: false, align: 'end' as const }
]);

const dialog = ref({
    show: false,
    isEdit: false,
    id: '',
    categoryId: '',
    categoryName: '',
    amountInput: 0 as number,
    currency: ''
});

const canSaveDialog = computed(() => {
    if (!dialog.value.isEdit && !dialog.value.categoryId) {
        return false;
    }
    return dialog.value.amountInput > 0 && dialog.value.currency.length > 0;
});

function loadCurrentOverview(): void {
    loading.value = true;
    budgetStore.loadOverview({ startDate: currentMonthStartDate.value }).catch(error => {
        if (!error.processed) {
            snackbar.value?.showError(error);
        }
    }).finally(() => {
        loading.value = false;
    });
}

function changeMonth(delta: number): void {
    currentMonthStartDate.value = addMonthsToDate(currentMonthStartDate.value, delta);
    loadCurrentOverview();
}

function resetToCurrentMonth(): void {
    currentMonthStartDate.value = firstOfMonth(new Date());
    loadCurrentOverview();
}

function progressPercent(item: CategoryBudgetOverviewItem): number {
    if (item.amount <= 0) {
        return 0;
    }
    const pct = (item.actualExpenseAmount / item.amount) * 100;
    return Math.min(100, Math.max(0, pct));
}

function formatMonth(unixSeconds: number): string {
    const d = new Date(unixSeconds * 1000);
    const monthNames = ['January', 'February', 'March', 'April', 'May', 'June', 'July', 'August', 'September', 'October', 'November', 'December'];
    return `${tt(`datetime.${monthNames[d.getMonth()]}.long`)} ${d.getFullYear()}`;
}

function displayAmount(cents: number, currency: string): string {
    return formatAmountToLocalizedNumeralsWithCurrency(parseBigDecimal(cents), currency);
}

function openAddDialog(): void {
    dialog.value = {
        show: true,
        isEdit: false,
        id: '',
        categoryId: '',
        categoryName: '',
        amountInput: 0,
        currency: defaultCurrency.value
    };
}

function openEditDialog(item: CategoryBudgetOverviewItem): void {
    dialog.value = {
        show: true,
        isEdit: true,
        id: item.id,
        categoryId: item.categoryId,
        categoryName: item.categoryName,
        amountInput: item.amount / 100,
        currency: item.currency
    };
}

function saveDialog(): void {
    if (!canSaveDialog.value) {
        return;
    }

    const amountCents = Math.round(dialog.value.amountInput * 100);
    const budget = dialog.value.isEdit
        ? CategoryBudgetLimit.of({ id: dialog.value.id, categoryId: dialog.value.categoryId, startDate: currentMonthStartDate.value, endDate: 0, amount: amountCents, currency: dialog.value.currency })
        : CategoryBudgetLimit.createNew(dialog.value.categoryId, currentMonthStartDate.value, dialog.value.currency);

    budget.amount = amountCents;
    budget.currency = dialog.value.currency;

    budgetStore.saveBudget({
        budget,
        isEdit: dialog.value.isEdit,
        clientSessionId: `${Date.now()}`
    }).then(() => {
        dialog.value.show = false;
        snackbar.value?.showMessage(dialog.value.isEdit ? 'Budget saved' : 'Budget added');
        loadCurrentOverview();
    }).catch(error => {
        if (!error.processed) {
            snackbar.value?.showError(error);
        }
    });
}

function confirmDelete(item: CategoryBudgetOverviewItem): void {
    if (!confirm(tt('Delete this budget?'))) {
        return;
    }

    budgetStore.deleteBudget({ budgetId: item.id }).then(() => {
        snackbar.value?.showMessage('Budget deleted');
        loadCurrentOverview();
    }).catch(error => {
        if (!error.processed) {
            snackbar.value?.showError(error);
        }
    });
}

onMounted(() => {
    categoriesStore.loadAllCategories({ force: false }).catch(() => {
        // Categories may already be loaded; ignore errors here.
    });
    loadCurrentOverview();
});
</script>

<style scoped>
.budget-stat-card {
    background: rgba(var(--v-theme-on-surface), 0.04);
}
.cursor-pointer {
    cursor: pointer;
}
</style>
