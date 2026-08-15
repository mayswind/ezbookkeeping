<template>
    <v-dialog width="600" :persistent="submitting || !!selectedNames.length" v-model="showState">
        <one-column-dialog-layout content-class="pa-0"
                                  :title="tt(title)" :cancel-button-title="tt('Cancel')"
                                  :disabled="submitting"
                                  @cancel="cancel">
            <template #toolbar>
                <v-btn class="ms-2" density="comfortable" variant="outlined"
                       :disabled="submitting || !selectedNames || !selectedNames.length" @click="confirm">
                    {{ tt('OK') }}
                    <v-progress-circular indeterminate size="22" class="ms-2" v-if="submitting"></v-progress-circular>
                </v-btn>

                <v-btn density="compact" color="default" variant="text" class="ms-2"
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
            </template>

            <template #content>
                <v-list class="mx-3 d-flex flex-column gap-1" density="comfortable" select-strategy="classic"
                        :disabled="submitting" v-model:selected="selectedNames">
                    <v-list-item class="py-0" density="compact"
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
            </template>
        </one-column-dialog-layout>
    </v-dialog>

    <snack-bar ref="snackbar" />
</template>

<script setup lang="ts">
import SnackBar from '@/components/desktop/SnackBar.vue';

import { ref, computed, useTemplateRef } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import { useTransactionCategoriesStore } from '@/stores/transactionCategory.ts';
import { useTransactionTagsStore } from '@/stores/transactionTag.ts';

import { type NameValue, values } from '@/core/base.ts';
import { IconType } from '@/core/icon.ts';
import { CategoryType } from '@/core/category.ts';
import { AUTOMATICALLY_CREATED_CATEGORY_ICON_ID } from '@/consts/icon.ts';
import { DEFAULT_CATEGORY_COLOR } from '@/consts/color.ts';
import { DEFAULT_TAG_GROUP_ID } from '@/consts/tag.ts';

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

let resolveFunc: ((response: BatchCreateDialogResponse) => void) | null = null;
let rejectFunc: ((reason?: unknown) => void) | null = null;

const showState = ref<boolean>(false);
const submitting = ref<boolean>(false);
const type = ref<BatchCreateDialogDataType | ''>('');
const invalidItems = ref<NameValue[] | undefined>([]);
const selectedNames = ref<string[]>([]);

const title = computed<string>(() => {
    if (type.value === 'expenseCategory') {
        return tt('Create Nonexistent Expense Categories');
    } else if (type.value === 'incomeCategory') {
        return tt('Create Nonexistent Income Categories');
    } else if (type.value === 'transferCategory') {
        return tt('Create Nonexistent Transfer Categories');
    } else if (type.value === 'tag') {
        return tt('Create Nonexistent Transaction Tags');
    }

    return '';
});

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
            iconType: IconType.System,
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
            const tag: TransactionTag = TransactionTag.createNewTag(item, DEFAULT_TAG_GROUP_ID);
            submitTags.push(tag.toCreateRequest());
        }

        transactionTagsStore.addTags({
            tags: submitTags,
            groupId: DEFAULT_TAG_GROUP_ID,
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
