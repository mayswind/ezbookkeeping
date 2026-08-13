// [PLUGIN:budget] Category budget overview store.
import { ref } from 'vue';
import { defineStore } from 'pinia';

import type { CategoryBudgetOverviewResponse, CategoryBudgetOverviewItem } from '@/models/category_budget_limit.ts';
import { CategoryBudgetLimit } from '@/models/category_budget_limit.ts';

import services from '@/lib/services.ts';
import logger from '@/lib/logger.ts';

export const useCategoryBudgetStore = defineStore('categoryBudget', () => {
    const currentOverview = ref<CategoryBudgetOverviewResponse | null>(null);
    const currentMonthStartDate = ref<number>(0);
    const loading = ref<boolean>(false);

    function loadOverview({ startDate }: { startDate: number }): Promise<CategoryBudgetOverviewResponse | null> {
        loading.value = true;
        currentMonthStartDate.value = startDate;

        return new Promise((resolve, reject) => {
            services.getCategoryBudgetOverview({ startDate: String(startDate) }).then(response => {
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    reject({ message: 'Unable to retrieve budget overview' });
                    return;
                }

                currentOverview.value = data.result;
                resolve(data.result);
            }).catch(error => {
                logger.error('failed to load budget overview', error);

                if (error.response && error.response.data && error.response.data.errorMessage) {
                    reject({ error: error.response.data });
                } else if (!error.processed) {
                    reject({ message: 'Unable to retrieve budget overview' });
                } else {
                    reject(error);
                }
            }).finally(() => {
                loading.value = false;
            });
        });
    }

    function saveBudget({ budget, isEdit, clientSessionId }: { budget: CategoryBudgetLimit, isEdit: boolean, clientSessionId: string }): Promise<CategoryBudgetLimit> {
        return new Promise((resolve, reject) => {
            let promise;

            if (!isEdit) {
                promise = services.addCategoryBudgetLimit(budget.toCreateRequest(clientSessionId));
            } else {
                promise = services.modifyCategoryBudgetLimit(budget.toModifyRequest());
            }

            promise.then(response => {
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    if (!isEdit) {
                        reject({ message: 'Unable to add budget' });
                    } else {
                        reject({ message: 'Unable to save budget' });
                    }
                    return;
                }

                resolve(CategoryBudgetLimit.of(data.result));
            }).catch(error => {
                logger.error('failed to save budget', error);

                if (error.response && error.response.data && error.response.data.errorMessage) {
                    reject({ error: error.response.data });
                } else if (!error.processed) {
                    if (!isEdit) {
                        reject({ message: 'Unable to add budget' });
                    } else {
                        reject({ message: 'Unable to save budget' });
                    }
                } else {
                    reject(error);
                }
            });
        });
    }

    function deleteBudget({ budgetId }: { budgetId: string }): Promise<boolean> {
        return new Promise((resolve, reject) => {
            services.deleteCategoryBudgetLimit({ id: budgetId }).then(response => {
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    reject({ message: 'Unable to delete this budget' });
                    return;
                }

                resolve(data.result);
            }).catch(error => {
                logger.error('failed to delete budget', error);

                if (error.response && error.response.data && error.response.data.errorMessage) {
                    reject({ error: error.response.data });
                } else if (!error.processed) {
                    reject({ message: 'Unable to delete this budget' });
                } else {
                    reject(error);
                }
            });
        });
    }

    function reset(): void {
        currentOverview.value = null;
        currentMonthStartDate.value = 0;
    }

    return {
        // states
        currentOverview,
        currentMonthStartDate,
        loading,
        // functions
        loadOverview,
        saveBudget,
        deleteBudget,
        reset
    };
});

// Helper type re-export for views
export type { CategoryBudgetOverviewResponse, CategoryBudgetOverviewItem };
