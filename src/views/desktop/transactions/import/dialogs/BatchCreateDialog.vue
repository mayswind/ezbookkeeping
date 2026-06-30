<template>
    <v-dialog width="600" :persistent="submitting || !!selectedNames.length" v-model="showState">
        <v-card class="pa-sm-1 pa-md-2">
            <template #title>
                <div class="d-flex flex-wrap">
                    <h4 class="text-h4 text-wrap" v-if="type === 'expenseCategory'">{{ tt('Create Nonexistent Expense Categories') }}</h4>
                    <h4 class="text-h4 text-wrap" v-if="type === 'incomeCategory'">{{ tt('Create Nonexistent Income Categories') }}</h4>
                    <h4 class="text-h4 text-wrap" v-if="type === 'transferCategory'">{{ tt('Create Nonexistent Transfer Categories') }}</h4>
                    <h4 class="text-h4 text-wrap" v-if="type === 'account'">{{ tt('Create Nonexistent Accounts') }}</h4>
                    <h4 class="text-h4 text-wrap" v-if="type === 'tag'">{{ tt('Create Nonexistent Transaction Tags') }}</h4>
                    <v-spacer/>
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
            <v-card-text class="d-flex flex-column flex-md-row flex-grow-1 overflow-y-auto">
                <v-row>
                    <v-col cols="12" class="px-0">
                        <v-list class="py-0" density="compact" select-strategy="classic"
                                :disabled="submitting" v-model:selected="selectedNames">
                            <v-list-item class="mx-1 px-2 py-0"
                                         :key="item.value" :value="item.value" :title="item.name"
                                         v-for="item in invalidItems">
                                <template #prepend="{ isActive }">
                                    <v-list-item-action start>
                                        <v-checkbox-btn :model-value="isActive"
                                                        @update:model-value="updateSelectedNames(item.value, $event)"></v-checkbox-btn>
                                    </v-list-item-action>
                                </template>
                            </v-list-item>
                        </v-list>
                    </v-col>
                </v-row>
            </v-card-text>
            <v-card-text>
                <div class="w-100 d-flex justify-center flex-wrap mt-sm-1 mt-md-2 gap-4">
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

import { useUserStore } from '@/stores/user.ts';
import { useAccountsStore } from '@/stores/account.ts';
import { useTransactionCategoriesStore } from '@/stores/transactionCategory.ts';
import { useTransactionTagsStore } from '@/stores/transactionTag.ts';

import { type NameValue, keys } from '@/core/base.ts';
import { CategoryType } from '@/core/category.ts';
import { AccountCategory } from '@/core/account.ts';
import { AUTOMATICALLY_CREATED_CATEGORY_ICON_ID } from '@/consts/icon.ts';
import { DEFAULT_CATEGORY_COLOR } from '@/consts/color.ts';
import { DEFAULT_TAG_GROUP_ID } from '@/consts/tag.ts';

import { TransactionCategory } from '@/models/transaction_category.ts';
import { type TransactionTagCreateRequest, TransactionTag } from '@/models/transaction_tag.ts';
import { Account } from '@/models/account.ts';
import { parseImportCategoryCompositeKey, parseImportAccountCompositeKey } from '@/models/imported_transaction.ts';

import { isDefined, arrayItemToObjectField } from '@/lib/common.ts';
import { getCurrentUnixTime } from '@/lib/datetime.ts';

import {
    mdiSelectAll,
    mdiSelect,
    mdiSelectInverse,
    mdiDotsVertical
} from '@mdi/js';

export type BatchCreateDialogDataType = 'expenseCategory' | 'incomeCategory' | 'transferCategory' | 'account' | 'tag';

type SnackBarType = InstanceType<typeof SnackBar>;

interface BatchCreateDialogResponse {
    sourceTargetMap: Record<string, string>;
}

const { tt } = useI18n();

const userStore = useUserStore();
const accountsStore = useAccountsStore();
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

interface SelectedCategoryItem {
    // 选项 value（"一级名 + 二级名"组合键，一级可能为空）
    compositeKey: string;
    // 原始一级名（可能为空，空时回退到默认一级名）
    primaryCategoryName: string;
    // 二级名
    subCategoryName: string;
}

// 按一级分类分组创建：一级已存在则复用、不存在则创建，再在其下创建缺失的二级，
// 返回 组合键 -> 新建/复用的二级分类 id 的映射，供调用方回填交易的 categoryId
// categoryType: 分类类型（支出/收入/转账）
// defaultPrimaryCategoryName: 当原始一级名为空时使用的默认一级名
// selectedItems: 用户在对话框中选中的待创建项
async function createCategoriesHierarchically(categoryType: CategoryType, defaultPrimaryCategoryName: string, selectedItems: SelectedCategoryItem[]): Promise<Record<string, string>> {
    const sourceTargetMap: Record<string, string> = {};

    // 按"实际使用的一级名"分组（原始一级为空时回退到默认一级名）
    const groupedItems: Record<string, SelectedCategoryItem[]> = {};

    for (const item of selectedItems) {
        const primaryCategoryName = item.primaryCategoryName || defaultPrimaryCategoryName;

        if (!groupedItems[primaryCategoryName]) {
            groupedItems[primaryCategoryName] = [];
        }

        groupedItems[primaryCategoryName]!.push(item);
    }

    for (const primaryCategoryName of keys(groupedItems)) {
        // 查找已存在的同名一级分类，存在则复用，否则新建一级
        const existingPrimaryCategories = transactionCategoriesStore.allTransactionCategories[categoryType] || [];
        let parentCategory: TransactionCategory | undefined = existingPrimaryCategories.find(category => category.name === primaryCategoryName);

        if (!parentCategory) {
            const newPrimaryCategory = TransactionCategory.createNewCategory(categoryType);
            newPrimaryCategory.name = primaryCategoryName;
            newPrimaryCategory.icon = AUTOMATICALLY_CREATED_CATEGORY_ICON_ID;
            newPrimaryCategory.color = DEFAULT_CATEGORY_COLOR;
            parentCategory = await transactionCategoriesStore.saveCategory({ category: newPrimaryCategory, isEdit: false, clientSessionId: '' });
        }

        // 该一级下已存在的二级名 -> id，用于复用并避免重复创建
        const subCategoryIdByName: Record<string, string> = {};

        for (const subCategory of (parentCategory.subCategories || [])) {
            subCategoryIdByName[subCategory.name] = subCategory.id;
        }

        for (const item of groupedItems[primaryCategoryName]!) {
            const existingSubCategoryId = subCategoryIdByName[item.subCategoryName];

            if (isDefined(existingSubCategoryId)) {
                sourceTargetMap[item.compositeKey] = existingSubCategoryId;
                continue;
            }

            const newSubCategory = TransactionCategory.createNewCategory(categoryType, parentCategory.id);
            newSubCategory.name = item.subCategoryName;
            newSubCategory.icon = AUTOMATICALLY_CREATED_CATEGORY_ICON_ID;
            const createdSubCategory = await transactionCategoriesStore.saveCategory({ category: newSubCategory, isEdit: false, clientSessionId: '' });

            subCategoryIdByName[item.subCategoryName] = createdSubCategory.id;
            sourceTargetMap[item.compositeKey] = createdSubCategory.id;
        }
    }

    return sourceTargetMap;
}

// 逐个创建不存在的账户：账户分类默认用"现金"，币种取自 CSV（缺失时回退用户默认币种），
// 返回 账户名 -> 新建账户 id 的映射，供调用方回填交易的 sourceAccountId / destinationAccountId
// selectedKeys: 用户选中的"账户名 + 币种"组合键列表
async function createNonexistentAccounts(selectedKeys: string[]): Promise<Record<string, string>> {
    const sourceTargetMap: Record<string, string> = {};
    const balanceTime = getCurrentUnixTime();

    for (const compositeKey of selectedKeys) {
        const parsed = parseImportAccountCompositeKey(compositeKey);
        const accountName = parsed.accountName;
        const accountCurrency = parsed.accountCurrency || userStore.currentUserDefaultCurrency;

        // 已存在同名账户则直接复用，避免重复创建
        const existingAccount = accountsStore.allAccounts.find(account => account.name === accountName);

        if (existingAccount) {
            sourceTargetMap[accountName] = existingAccount.id;
            continue;
        }

        const newAccount = Account.createNewAccount(AccountCategory.Cash, accountCurrency, balanceTime);
        newAccount.name = accountName;

        const createdAccount = await accountsStore.saveAccount({ account: newAccount, subAccounts: [], isEdit: false, clientSessionId: '' });
        sourceTargetMap[accountName] = createdAccount.id;
    }

    return sourceTargetMap;
}

function buildBatchCreateTagResponse(createdTags: TransactionTag[]): BatchCreateDialogResponse {
    const createdTagIdMap: Record<string, string> = {};
    const sourceTargetMap: Record<string, string> = {};

    // 新建标签名 -> id（标签的 value 即为标签名本身）
    for (const tag of createdTags) {
        createdTagIdMap[tag.name] = tag.id;
    }

    for (const item of (invalidItems.value || [])) {
        const targetItem = createdTagIdMap[item.value];

        if (!isDefined(targetItem)) {
            continue;
        }

        sourceTargetMap[item.value] = targetItem;
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
    selectedNames.value = (invalidItems.value || []).map(item => item.value);
}

function selectNoneItems(): void {
    selectedNames.value = [];
}

function selectInvertItems(): void {
    const currentSelectedNames: Record<string, boolean> = arrayItemToObjectField(selectedNames.value, true);
    selectedNames.value = [];

    for (const item of (invalidItems.value || [])) {
        if (!currentSelectedNames[item.value]) {
            selectedNames.value.push(item.value);
        }
    }
}

function confirm(): void {
    if (type.value === 'expenseCategory' || type.value === 'incomeCategory' || type.value === 'transferCategory') {
        submitting.value = true;

        let categoryType: CategoryType = CategoryType.Expense;
        let defaultPrimaryCategoryName = '';

        if (type.value === 'expenseCategory') {
            categoryType = CategoryType.Expense;
            defaultPrimaryCategoryName = tt('Default Expense Category');
        } else if (type.value === 'incomeCategory') {
            categoryType = CategoryType.Income;
            defaultPrimaryCategoryName = tt('Default Income Category');
        } else if (type.value === 'transferCategory') {
            categoryType = CategoryType.Transfer;
            defaultPrimaryCategoryName = tt('Default Transfer Category');
        }

        // 解析选中项：每项携带原始一级名（可能为空）与二级名
        const selectedItems: SelectedCategoryItem[] = selectedNames.value.map(compositeKey => {
            const parsed = parseImportCategoryCompositeKey(compositeKey);
            return {
                compositeKey: compositeKey,
                primaryCategoryName: parsed.primaryCategoryName,
                subCategoryName: parsed.subCategoryName
            };
        });

        createCategoriesHierarchically(categoryType, defaultPrimaryCategoryName, selectedItems).then(sourceTargetMap => {
            transactionCategoriesStore.loadAllCategories({ force: false }).then(() => {
                submitting.value = false;
                showState.value = false;

                resolveFunc?.({ sourceTargetMap: sourceTargetMap });
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
    } else if (type.value === 'account') {
        submitting.value = true;

        createNonexistentAccounts(selectedNames.value).then(sourceTargetMap => {
            submitting.value = false;
            showState.value = false;

            resolveFunc?.({ sourceTargetMap: sourceTargetMap });
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
