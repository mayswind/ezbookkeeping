<template>
    <f7-page @page:afterin="onPageAfterIn">
        <f7-navbar>
            <f7-nav-left :back-link="tt('Back')" v-if="!displayOrderModified"></f7-nav-left>
            <f7-nav-left v-else-if="displayOrderModified">
                <f7-link icon-f7="xmark" :class="{ 'disabled': displayOrderSaving }" @click="cancelSort"></f7-link>
            </f7-nav-left>
            <f7-nav-title :title="tt('Transaction Tag Groups')"></f7-nav-title>
            <f7-nav-right class="navbar-compact-icons">
                <f7-link icon-f7="checkmark_alt" :class="{ 'disabled': displayOrderSaving || !displayOrderModified }" @click="saveSortResult"></f7-link>
            </f7-nav-right>
        </f7-navbar>

        <f7-list strong inset dividers class="tag-group-item-list margin-top skeleton-text" v-if="loading">
            <f7-list-item :key="itemIdx" v-for="itemIdx in [ 1, 2, 3 ]">
                <template #title>
                    <div class="display-flex">
                        <div class="transaction-tag-group-list-item-content list-item-valign-middle padding-inline-start-half">Tag Group Name</div>
                    </div>
                </template>
            </f7-list-item>
        </f7-list>

        <f7-list strong inset dividers class="tag-group-item-list margin-top" v-if="!loading && tagGroups.length < 1">
            <f7-list-item :title="tt('No available tag group')"></f7-list-item>
        </f7-list>

        <f7-list strong inset dividers sortable sortable-enabled
                 class="tag-group-item-list margin-top"
                 @sortable:sort="onSort"
                 v-if="!loading">
            <f7-list-item :id="getTagGroupDomId(tagGroup)"
                          :key="tagGroup.id"
                          v-for="tagGroup in tagGroups">
                <template #title>
                    <div class="display-flex">
                        <div class="transaction-tag-group-list-item-content list-item-valign-middle padding-inline-start-half">{{ tagGroup.name }}</div>
                    </div>
                </template>
            </f7-list-item>
        </f7-list>
    </f7-page>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';
import type { Router } from 'framework7/types';

import { useI18n } from '@/locales/helpers.ts';
import { useI18nUIComponents, showLoading, hideLoading } from '@/lib/ui/mobile.ts';

import { useTransactionTagsStore } from '@/stores/transactionTag.ts';

import { TransactionTagGroup } from '@/models/transaction_tag_group.ts';

const props = defineProps<{
    f7router: Router.Router;
}>();

const { tt } = useI18n();
const { showToast, routeBackOnError } = useI18nUIComponents();

const transactionTagsStore = useTransactionTagsStore();

const loading = ref<boolean>(true);
const loadingError = ref<unknown | null>(null);
const displayOrderModified = ref<boolean>(false);
const displayOrderSaving = ref<boolean>(false);

const tagGroups = computed<TransactionTagGroup[]>(() => transactionTagsStore.allTransactionTagGroups);

function getTagGroupDomId(tagGroup: TransactionTagGroup): string {
    return 'tagGroup_' + tagGroup.id;
}

function parseTagGroupIdFromDomId(domId: string): string | null {
    if (!domId || domId.indexOf('tagGroup_') !== 0) {
        return null;
    }

    return domId.substring(9); // tagGroup_
}

function init(): void {
    loading.value = true;

    transactionTagsStore.loadAllTagGroups({
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
}

function saveSortResult(): void {
    if (!displayOrderModified.value) {
        return;
    }

    displayOrderSaving.value = true;
    showLoading();

    transactionTagsStore.updateTagGroupDisplayOrders().then(() => {
        displayOrderSaving.value = false;
        hideLoading();

        displayOrderModified.value = false;
    }).catch(error => {
        displayOrderSaving.value = false;
        hideLoading();

        if (!error.processed) {
            showToast(error.message || error);
        }
    });
}

function cancelSort(): void {
    if (!displayOrderModified.value) {
        return;
    }

    displayOrderSaving.value = true;
    showLoading();

    transactionTagsStore.loadAllTagGroups({
        force: false
    }).then(() => {
        displayOrderSaving.value = false;
        hideLoading();

        displayOrderModified.value = false;
    }).catch(error => {
        displayOrderSaving.value = false;
        hideLoading();

        if (!error.processed) {
            showToast(error.message || error);
        }
    });
}

function onSort(event: { el: { id: string }, from: number, to: number }): void {
    if (!event || !event.el || !event.el.id) {
        showToast('Unable to move tag group');
        return;
    }

    const id = parseTagGroupIdFromDomId(event.el.id);

    if (!id) {
        showToast('Unable to move tag group');
        return;
    }

    transactionTagsStore.changeTagGroupDisplayOrder({
        tagGroupId: id,
        from: event.from,
        to: event.to
    }).then(() => {
        displayOrderModified.value = true;
    }).catch(error => {
        showToast(error.message || error);
    });
}

function onPageAfterIn(): void {
    if (transactionTagsStore.transactionTagGroupListStateInvalid && !loading.value) {
        transactionTagsStore.loadAllTagGroups({}).catch(error => {
            if (!error.processed) {
                showToast(error.message || error);
            }
        });
    }

    routeBackOnError(props.f7router, loadingError);
}

init();
</script>

<style>
.transaction-tag-group-list-item-content {
    overflow: hidden;
    text-overflow: ellipsis;
}
</style>
