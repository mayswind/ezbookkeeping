<template>
    <v-dialog width="800" :persistent="submitting" v-model="showState">
        <one-column-dialog-layout content-class="pt-0"
                                  :disabled="submitting"
                                  :title="tt('Default Categories')" :cancel-button-title="tt('Cancel')"
                                  @cancel="showState = false">
            <template #content>
                <template :key="categoryType" v-for="(categories, categoryType) in allPresetCategories">
                    <div class="d-flex align-center text-body-large my-1">
                        <span class="font-weight-bold ms-1">{{ getCategoryTypeName(parseInt(categoryType)) }}</span>
                        <v-spacer/>
                        <language-select-button :disabled="submitting"
                                                :use-model-value="true" v-model="currentLocale" />
                    </div>

                    <v-expansion-panels class="compacted-expansion-panels" variant="accordion" multiple :disabled="submitting">
                        <v-expansion-panel :key="idx" v-for="(category, idx) in categories">
                            <v-expansion-panel-title class="py-0 px-4">
                                <ItemIcon :icon-type="getCategoryIconType(category.iconType)" :icon-id="category.icon" :color="category.color"></ItemIcon>
                                <span class="text-body-medium ms-2">{{ category.name }}</span>
                            </v-expansion-panel-title>
                            <v-expansion-panel-text v-if="category.subCategories.length">
                                <v-list rounded density="comfortable" class="pa-0">
                                    <template :key="subIdx"
                                              v-for="(subCategory, subIdx) in category.subCategories">
                                        <v-list-item>
                                            <template #prepend>
                                                <ItemIcon :icon-type="getCategoryIconType(subCategory.iconType)" :icon-id="subCategory.icon" :color="subCategory.color"></ItemIcon>
                                            </template>
                                            <span class="text-body-medium ms-2">{{ subCategory.name }}</span>
                                        </v-list-item>
                                        <v-divider v-if="subIdx !== category.subCategories.length - 1"/>
                                    </template>
                                </v-list>
                            </v-expansion-panel-text>
                        </v-expansion-panel>
                    </v-expansion-panels>
                </template>
            </template>

            <template #footer>
                <v-spacer/>
                <v-btn :disabled="submitting" @click="save">
                    {{ tt('Save') }}
                    <v-progress-circular indeterminate size="22" class="ms-2" v-if="submitting"></v-progress-circular>
                </v-btn>
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

import { type LocalizedPresetCategory, CategoryType } from '@/core/category.ts';
import { categorizedArrayToPlainArray } from '@/lib/common.ts';
import { getCategoryIconType } from '@/lib/icon.ts';
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

const allPresetCategories = computed<Record<string, LocalizedPresetCategory[]>>(() => getAllTransactionDefaultCategories(props.categoryType, currentLocale.value));

const showState = computed<boolean>({
    get: () => props.show,
    set: (value) => emit('update:show', value)
});

function getCategoryTypeName(categoryType: number): string {
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
