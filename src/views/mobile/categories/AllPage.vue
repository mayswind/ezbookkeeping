<template>
    <f7-page ptr @ptr:refresh="reload" @page:afterin="onPageAfterIn">
        <f7-navbar :title="tt('Transaction Categories')" :back-link="tt('Back')"></f7-navbar>

        <f7-list strong inset dividers class="margin-top skeleton-text" v-if="loading">
            <f7-list-item title="Expense" link="#"></f7-list-item>
            <f7-list-item title="Income" link="#"></f7-list-item>
            <f7-list-item title="Transfer" link="#"></f7-list-item>
        </f7-list>

        <f7-list strong inset dividers class="margin-top" v-else-if="!loading">
            <f7-list-item :title="tt('Expense')" link="/category/list?type=2"></f7-list-item>
            <f7-list-item :title="tt('Income')" link="/category/list?type=1"></f7-list-item>
            <f7-list-item :title="tt('Transfer')" link="/category/list?type=3"></f7-list-item>
        </f7-list>
    </f7-page>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import type { Router } from 'framework7/types';

import { useI18n } from '@/locales/helpers.ts';
import { useI18nUIComponents } from '@/lib/ui/mobile.ts';

import { useTransactionCategoriesStore } from '@/stores/transactionCategory.ts';

const props = defineProps<{
    f7router: Router.Router;
}>();

const { tt } = useI18n();
const { showToast, routeBackOnError } = useI18nUIComponents();

const transactionCategoriesStore = useTransactionCategoriesStore();

const loading = ref<boolean>(true);
const loadingError = ref<unknown | null>(null);

function onPageAfterIn(): void {
    routeBackOnError(props.f7router, loadingError);
}

function reload(done?: () => void): void {
    const force = !!done;

    transactionCategoriesStore.loadAllCategories({
        force: force
    }).then(() => {
        done?.();

        if (force) {
            showToast('Category list has been updated');
        }
    }).catch(error => {
        done?.();

        if (!error.processed) {
            showToast(error.message || error);
        }
    });
}

transactionCategoriesStore.loadAllCategories({
    force: false
}).then(() => {
    loading.value = false;
}).catch(error => {
    if (error.processed) {
        loading.value = false;
    } else {
        loadingError.value = error;
        showToast(error.message || error);
    }
});
</script>
