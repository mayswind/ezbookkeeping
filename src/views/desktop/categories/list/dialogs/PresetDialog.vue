<template>
    <v-dialog width="800" :persistent="submitting" v-model="showState">
        <v-card class="pa-2 pa-sm-4 pa-md-8">
            <template #title>
                <div class="d-flex align-center justify-center">
                    <h4 class="text-h4">{{ tt('Default Categories') }}</h4>
                </div>
            </template>
            <v-card-text class="preset-transaction-categories mt-sm-2 mt-md-4 pt-0">
                <template :key="categoryType" v-for="(categories, categoryType) in allPresetCategories">
                    <div class="d-flex align-center mb-1">
                        <h4>{{ getCategoryTypeName(categoryType) }}</h4>
                        <v-spacer/>
                        <language-select-button :disabled="submitting"
                                                :use-model-value="true" v-model="currentLocale" />
                    </div>

                    <v-expansion-panels class="border rounded mb-2" variant="accordion" multiple :disabled="submitting">
                        <v-expansion-panel :key="idx" v-for="(category, idx) in categories">
                            <v-expansion-panel-title class="py-0">
                                <ItemIcon icon-type="category" :icon-id="category.icon" :color="category.color"></ItemIcon>
                                <span class="ms-3">{{ category.name }}</span>
                            </v-expansion-panel-title>
                            <v-expansion-panel-text v-if="category.subCategories.length">
                                <v-list rounded density="comfortable" class="pa-0">
                                    <template :key="subIdx"
                                              v-for="(subCategory, subIdx) in category.subCategories">
                                        <v-list-item>
                                            <template #prepend>
                                                <ItemIcon icon-type="category" :icon-id="subCategory.icon" :color="subCategory.color"></ItemIcon>
                                            </template>
                                            <span class="ms-3">{{ subCategory.name }}</span>
                                        </v-list-item>
                                        <v-divider v-if="subIdx !== category.subCategories.length - 1"/>
                                    </template>
                                </v-list>
                            </v-expansion-panel-text>
                        </v-expansion-panel>
                    </v-expansion-panels>
                </template>
            </v-card-text>
            <v-card-text class="overflow-y-visible">
                <div class="w-100 d-flex justify-center mt-2 mt-sm-4 mt-md-6 gap-4">
                    <v-btn :disabled="submitting" @click="save">
                        {{ tt('Save') }}
                        <v-progress-circular indeterminate size="22" class="ms-2" v-if="submitting"></v-progress-circular>
                    </v-btn>
                    <v-btn color="secondary" density="default" variant="tonal"
                           :disabled="submitting" @click="showState = false">{{ tt('Cancel') }}</v-btn>
                </div>
            </v-card-text>
        </v-card>
    </v-dialog>

    <snack-bar ref="snackbar" />
</template>

<script setup lang="ts">
import SnackBar from '@/components/desktop/SnackBar.vue';

import { ref, computed, useTemplateRef } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import { useTransactionCategoriesStore } from '@/stores/transactionCategory.ts';

import type { PartialRecord } from '@/core/base.ts';
import { type LocalizedPresetCategory, CategoryType } from '@/core/category.ts';
import { categorizedArrayToPlainArray } from '@/lib/common.ts';
import { localizedPresetCategoriesToTransactionCategoryCreateWithSubCategories } from '@/lib/category.ts';

type SnackBarType = InstanceType<typeof SnackBar>;

const props = defineProps<{
    categoryType: CategoryType;
    show: boolean;
}>();

const emit = defineEmits<{
    (e: 'update:show', value: boolean): void;
    (e: 'category:saved', event: { message: string }): void;
}>();

const { tt, getCurrentLanguageTag, getAllTransactionDefaultCategories } = useI18n();

const transactionCategoriesStore = useTransactionCategoriesStore();

const snackbar = useTemplateRef<SnackBarType>('snackbar');

const currentLocale = ref<string>(getCurrentLanguageTag());
const submitting = ref<boolean>(false);

const allPresetCategories = computed<PartialRecord<CategoryType, LocalizedPresetCategory[]>>(() => getAllTransactionDefaultCategories(props.categoryType, currentLocale.value));

const showState = computed<boolean>({
    get: () => props.show,
    set: (value) => emit('update:show', value)
});

function getCategoryTypeName(categoryType: CategoryType): string {
    switch (categoryType) {
        case CategoryType.Income:
            return tt('Income Categories');
        case CategoryType.Expense:
            return tt('Expense Categories');
        case CategoryType.Transfer:
            return tt('Transfer Categories');
        default:
            return tt('Transaction Categories');
    }
}

function save(): void {
    submitting.value = true;

    const presetCategoriesArray = categorizedArrayToPlainArray(allPresetCategories.value);
    const submitCategories = localizedPresetCategoriesToTransactionCategoryCreateWithSubCategories(presetCategoriesArray);

    transactionCategoriesStore.addPresetCategories({
        categories: submitCategories
    }).then(() => {
        submitting.value = false;
        showState.value = false;

        emit('category:saved', {
            message: 'You have added preset categories'
        });
    }).catch(error => {
        submitting.value = false;

        if (!error.processed) {
            snackbar.value?.showError(error);
        }
    });
}
</script>

<style>
.preset-transaction-categories .v-expansion-panel-text__wrapper {
    padding: 0 0 0 0;
    padding-inline-start: 20px;
}
</style>
