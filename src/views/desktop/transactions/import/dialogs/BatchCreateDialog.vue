<template>
    <v-dialog width="600" :persistent="submitting || !!selectedNames.length" v-model="showState">
        <v-card class="pa-2 pa-sm-4 pa-md-4">
            <template #title>
                <div class="d-flex align-center justify-center">
                    <div class="d-flex w-100 align-center justify-center">
                        <h4 class="text-h4" v-if="type === 'expenseCategory'">{{ tt('Create Nonexistent Expense Categories') }}</h4>
                        <h4 class="text-h4" v-if="type === 'incomeCategory'">{{ tt('Create Nonexistent Income Categories') }}</h4>
                        <h4 class="text-h4" v-if="type === 'transferCategory'">{{ tt('Create Nonexistent Transfer Categories') }}</h4>
                        <h4 class="text-h4" v-if="type === 'tag'">{{ tt('Create Nonexistent Transaction Tags') }}</h4>
                    </div>
                    <v-btn density="comfortable" color="default" variant="text" class="ms-2"
                           :disabled="submitting || !invalidItems || !invalidItems.length" :icon="true">
                        <v-icon :icon="mdiDotsVertical" />
                        <v-menu activator="parent">
                            <v-list>
                                <v-list-item :prepend-icon="mdiSelectAll"
                                             :title="tt('Select All')"
                                             :disabled="!invalidItems || !invalidItems.length"
                                             @click="selectAllItems"></v-list-item>
                                <v-list-item :prepend-icon="mdiSelect"
                                             :title="tt('Select None')"
                                             :disabled="!invalidItems || !invalidItems.length"
                                             @click="selectNoneItems"></v-list-item>
                                <v-list-item :prepend-icon="mdiSelectInverse"
                                             :title="tt('Invert Selection')"
                                             :disabled="!invalidItems || !invalidItems.length"
                                             @click="selectInvertItems"></v-list-item>
                            </v-list>
                        </v-menu>
                    </v-btn>
                </div>
            </template>
            <v-card-text class="my-md-4 w-100 d-flex justify-center">
                <v-row>
                    <v-col cols="12">
                        <v-list class="py-0" density="compact" select-strategy="classic"
                                :disabled="submitting" v-model:selected="selectedNames">
                            <v-list-item class="py-0"
                                         :key="item.value" :value="item.name" :title="item.name"
                                         v-for="item in invalidItems">
                                <template #prepend="{ isActive }">
                                    <v-list-item-action start>
                                        <v-checkbox-btn :model-value="isActive"
                                                        @update:model-value="updateSelectedNames(item.name, $event)"></v-checkbox-btn>
                                    </v-list-item-action>
                                </template>
                            </v-list-item>
                        </v-list>
                    </v-col>
                </v-row>
            </v-card-text>
            <v-card-text class="overflow-y-visible">
                <div class="w-100 d-flex justify-center gap-4">
                    <v-btn :disabled="submitting || !selectedNames || !selectedNames.length" @click="confirm">
                        {{ tt('OK') }}
                        <v-progress-circular indeterminate size="22" class="ms-2" v-if="submitting"></v-progress-circular>
                    </v-btn>
                    <v-btn color="secondary" variant="tonal" :disabled="submitting" @click="cancel">{{ tt('Cancel') }}</v-btn>
                </div>
            </v-card-text>
        </v-card>
    </v-dialog>

    <snack-bar ref="snackbar" />
</template>

<script setup lang="ts">
import SnackBar from '@/components/desktop/SnackBar.vue';

import { ref, useTemplateRef } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import { useTransactionCategoriesStore } from '@/stores/transactionCategory.ts';
import { useTransactionTagsStore } from '@/stores/transactionTag.ts';

import { type NameValue, values } from '@/core/base.ts';
import { CategoryType } from '@/core/category.ts';
import { AUTOMATICALLY_CREATED_CATEGORY_ICON_ID } from '@/consts/icon.ts';
import { DEFAULT_CATEGORY_COLOR } from '@/consts/color.ts';

import { type TransactionCategoryCreateRequest, type TransactionCategoryCreateWithSubCategories, TransactionCategory } from '@/models/transaction_category.ts';
import { type TransactionTagCreateRequest, TransactionTag } from '@/models/transaction_tag.ts';

import { isDefined, arrayItemToObjectField } from '@/lib/common.ts';

import {
    mdiSelectAll,
    mdiSelect,
    mdiSelectInverse,
    mdiDotsVertical
} from '@mdi/js';

export type BatchCreateDialogDataType = 'expenseCategory' | 'incomeCategory' | 'transferCategory' | 'tag';

type SnackBarType = InstanceType<typeof SnackBar>;

interface BatchCreateDialogResponse {
    sourceTargetMap: Record<string, string>;
}

const { tt } = useI18n();

const transactionCategoriesStore = useTransactionCategoriesStore();
const transactionTagsStore = useTransactionTagsStore();

const snackbar = useTemplateRef<SnackBarType>('snackbar');

const showState = ref<boolean>(false);
const submitting = ref<boolean>(false);
const type = ref<BatchCreateDialogDataType | ''>('');
const invalidItems = ref<NameValue[] | undefined>([]);
const selectedNames = ref<string[]>([]);

let resolveFunc: ((response: BatchCreateDialogResponse) => void) | null = null;
let rejectFunc: ((reason?: unknown) => void) | null = null;

function updateSelectedNames(value: string, selected: boolean | null): void {
    const newSelectedNames: string[] = [];

    for (const name of selectedNames.value) {
        if (name !== value || selected) {
            newSelectedNames.push(name);
        }
    }

    if (selected) {
        newSelectedNames.push(value);
    }

    selectedNames.value = newSelectedNames;
}

function buildBatchCreateCategoryResponse(createdCategories: Record<number, TransactionCategory[]>): BatchCreateDialogResponse {
    const displayNameSourceItemMap: Record<string, string> = {};
    const sourceTargetMap: Record<string, string> = {};

    for (const item of (invalidItems.value || [])) {
        displayNameSourceItemMap[item.name] = item.value;
    }

    for (const categories of values(createdCategories)) {
        for (const category of categories) {
            if (!category.subCategories || category.subCategories.length < 1) {
                continue;
            }

            for (const subCategory of category.subCategories) {
                const sourceItem = displayNameSourceItemMap[subCategory.name];

                if (!isDefined(sourceItem)) {
                    continue;
                }

                sourceTargetMap[sourceItem] = subCategory.id;
            }
        }
    }

    const response: BatchCreateDialogResponse = {
        sourceTargetMap: sourceTargetMap
    };

    return response;
}

function buildBatchCreateTagResponse(createdTags: TransactionTag[]): BatchCreateDialogResponse {
    const displayNameSourceItemMap: Record<string, string> = {};
    const sourceTargetMap: Record<string, string> = {};

    for (const item of (invalidItems.value || [])) {
        displayNameSourceItemMap[item.name] = item.value;
    }

    for (const tag of createdTags) {
        const sourceItem = displayNameSourceItemMap[tag.name];

        if (!isDefined(sourceItem)) {
            continue;
        }

        sourceTargetMap[sourceItem] = tag.id;
    }

    const response: BatchCreateDialogResponse = {
        sourceTargetMap: sourceTargetMap
    };

    return response;
}

function open(options: { type: BatchCreateDialogDataType, invalidItems?: NameValue[] }): Promise<BatchCreateDialogResponse> {
    type.value = options.type;
    invalidItems.value = options.invalidItems;
    selectedNames.value = [];

    showState.value = true;

    return new Promise((resolve, reject) => {
        resolveFunc = resolve;
        rejectFunc = reject;
    });
}

function selectAllItems(): void {
    selectedNames.value = (invalidItems.value || []).map(item => item.name);
}

function selectNoneItems(): void {
    selectedNames.value = [];
}

function selectInvertItems(): void {
    const currentSelectedNames: Record<string, boolean> = arrayItemToObjectField(selectedNames.value, true);
    selectedNames.value = [];

    for (const item of (invalidItems.value || [])) {
        if (!currentSelectedNames[item.name]) {
            selectedNames.value.push(item.name);
        }
    }
}

function confirm(): void {
    if (type.value === 'expenseCategory' || type.value === 'incomeCategory' || type.value === 'transferCategory') {
        submitting.value = true;

        let categoryType: CategoryType = CategoryType.Expense;
        let primaryCategoryName = '';

        if (type.value === 'expenseCategory') {
            categoryType = CategoryType.Expense;
            primaryCategoryName = tt('Default Expense Category');
        } else if (type.value === 'incomeCategory') {
            categoryType = CategoryType.Income;
            primaryCategoryName = tt('Default Income Category');
        } else if (type.value === 'transferCategory') {
            categoryType = CategoryType.Transfer;
            primaryCategoryName = tt('Default Transfer Category');
        }

        const subCategories: TransactionCategoryCreateRequest[] = [];

        for (const item of selectedNames.value) {
            const category: TransactionCategory = TransactionCategory.createNewCategory(categoryType);
            category.name = item;
            category.icon = AUTOMATICALLY_CREATED_CATEGORY_ICON_ID;
            subCategories.push(category.toCreateRequest(''));
        }

        const submitCategories: TransactionCategoryCreateWithSubCategories[] = [{
            name: primaryCategoryName,
            type: categoryType,
            icon: AUTOMATICALLY_CREATED_CATEGORY_ICON_ID,
            color: DEFAULT_CATEGORY_COLOR,
            subCategories: subCategories
        }];

        transactionCategoriesStore.addCategories({
            categories: submitCategories
        }).then(response => {
            transactionCategoriesStore.loadAllCategories({ force: false }).then(() => {
                submitting.value = false;
                showState.value = false;

                resolveFunc?.(buildBatchCreateCategoryResponse(response));
            }).catch(error => {
                submitting.value = false;

                if (!error.processed) {
                    snackbar.value?.showError(error);
                }
            });
        }).catch(error => {
            submitting.value = false;

            if (!error.processed) {
                snackbar.value?.showError(error);
            }
        });
    } else if (type.value === 'tag') {
        submitting.value = true;

        const submitTags: TransactionTagCreateRequest[] = [];

        for (const item of selectedNames.value) {
            const tag: TransactionTag = TransactionTag.createNewTag(item);
            submitTags.push(tag.toCreateRequest());
        }

        transactionTagsStore.addTags({
            tags: submitTags,
            skipExists: true
        }).then(response => {
            transactionTagsStore.loadAllTags({ force: false }).then(() => {
                submitting.value = false;
                showState.value = false;

                resolveFunc?.(buildBatchCreateTagResponse(response));
            }).catch(error => {
                submitting.value = false;

                if (!error.processed) {
                    snackbar.value?.showError(error);
                }
            });
        }).catch(error => {
            submitting.value = false;

            if (!error.processed) {
                snackbar.value?.showError(error);
            }
        });
    }
}

function cancel(): void {
    rejectFunc?.();
    showState.value = false;
}

defineExpose({
    open
});
</script>
