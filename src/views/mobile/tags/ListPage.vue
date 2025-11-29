<template>
    <f7-page :ptr="!sortable && !hasEditingTag" @ptr:refresh="reload" @page:afterin="onPageAfterIn">
        <f7-navbar>
            <f7-nav-left :back-link="tt('Back')"></f7-nav-left>
            <f7-nav-title :title="tt('Transaction Tags')"></f7-nav-title>
            <f7-nav-right class="navbar-compact-icons">
                <f7-link :class="{ 'disabled': hasEditingTag || !tags.length }" icon-f7="ellipsis" v-if="!sortable" @click="showMoreActionSheet = true"></f7-link>
                <f7-link :class="{ 'disabled': hasEditingTag }" icon-f7="plus" v-if="!sortable" @click="add"></f7-link>
                <f7-link :text="tt('Done')" :class="{ 'disabled': displayOrderSaving || hasEditingTag }" v-else-if="sortable" @click="saveSortResult"></f7-link>
            </f7-nav-right>
        </f7-navbar>

        <f7-list strong inset dividers class="tag-item-list margin-top skeleton-text" v-if="loading">
            <f7-list-item :key="itemIdx" v-for="itemIdx in [ 1, 2, 3 ]">
                <template #media>
                    <f7-icon class="transaction-tag-icon" f7="number"></f7-icon>
                </template>
                <template #title>
                    <div class="display-flex">
                        <div class="transaction-tag-list-item-content list-item-valign-middle padding-inline-start-half">Tag Name</div>
                    </div>
                </template>
            </f7-list-item>
        </f7-list>

        <f7-list strong inset dividers class="tag-item-list margin-top" v-if="!loading && noAvailableTag && !newTag">
            <f7-list-item :title="tt('No available tag')"></f7-list-item>
        </f7-list>

        <f7-list strong inset dividers sortable class="tag-item-list margin-top"
                 :sortable-enabled="sortable" @sortable:sort="onSort"
                 v-if="!loading">
            <f7-list-item swipeout
                          :class="{ 'actual-first-child': tag.id === firstShowingId, 'actual-last-child': tag.id === lastShowingId && !newTag, 'editing-list-item': editingTag.id === tag.id }"
                          :id="getTagDomId(tag)"
                          :key="tag.id"
                          v-for="tag in tags"
                          v-show="showHidden || !tag.hidden"
                          @taphold="setSortable()">
                <template #media>
                    <f7-icon class="transaction-tag-icon" f7="number">
                        <f7-badge color="gray" class="right-bottom-icon" v-if="tag.hidden">
                            <f7-icon f7="eye_slash_fill"></f7-icon>
                        </f7-badge>
                    </f7-icon>
                </template>
                <template #title>
                    <div class="display-flex">
                        <div class="transaction-tag-list-item-content list-item-valign-middle padding-inline-start-half"
                             v-if="editingTag.id !== tag.id">
                            {{ tag.name }}
                        </div>
                        <f7-input class="list-title-input padding-inline-start-half"
                                  type="text"
                                  :placeholder="tt('Tag Title')"
                                  v-else-if="editingTag.id === tag.id"
                                  v-model:value="editingTag.name"
                                  @keyup.enter="save(editingTag)">
                        </f7-input>
                    </div>
                </template>
                <template #after>
                    <f7-button :class="{ 'no-padding': true, 'disabled': !isTagModified(tag) }"
                               raised fill
                               icon-f7="checkmark_alt"
                               color="blue"
                               v-if="editingTag.id === tag.id"
                               @click="save(editingTag)">
                    </f7-button>
                    <f7-button class="no-padding margin-inline-start-half"
                               raised fill
                               icon-f7="xmark"
                               color="gray"
                               v-if="editingTag.id === tag.id"
                               @click="cancelSave(editingTag)">
                    </f7-button>
                </template>
                <f7-swipeout-actions :left="textDirection === TextDirection.LTR"
                                     :right="textDirection === TextDirection.RTL"
                                     v-if="sortable && editingTag.id !== tag.id">
                    <f7-swipeout-button :color="tag.hidden ? 'blue' : 'gray'" class="padding-horizontal"
                                        overswipe close @click="hide(tag, !tag.hidden)">
                        <f7-icon :f7="tag.hidden ? 'eye' : 'eye_slash'"></f7-icon>
                    </f7-swipeout-button>
                </f7-swipeout-actions>
                <f7-swipeout-actions :left="textDirection === TextDirection.RTL"
                                     :right="textDirection === TextDirection.LTR"
                                     v-if="!sortable && editingTag.id !== tag.id">
                    <f7-swipeout-button color="orange" close :text="tt('Edit')" @click="edit(tag)"></f7-swipeout-button>
                    <f7-swipeout-button color="red" class="padding-horizontal" @click="remove(tag, false)">
                        <f7-icon f7="trash"></f7-icon>
                    </f7-swipeout-button>
                </f7-swipeout-actions>
            </f7-list-item>

            <f7-list-item class="editing-list-item" v-if="newTag">
                <template #media>
                    <f7-icon class="transaction-tag-icon" f7="number"></f7-icon>
                </template>
                <template #title>
                    <div class="display-flex">
                        <f7-input class="list-title-input padding-inline-start-half"
                                  type="text"
                                  :placeholder="tt('Tag Title')"
                                  v-model:value="newTag.name"
                                  @keyup.enter="save(newTag)">
                        </f7-input>
                    </div>
                </template>
                <template #after>
                    <f7-button :class="{ 'no-padding': true, 'disabled': !isTagModified(newTag) }"
                               raised fill
                               icon-f7="checkmark_alt"
                               color="blue"
                               @click="save(newTag)">
                    </f7-button>
                    <f7-button class="no-padding margin-inline-start-half"
                               raised fill
                               icon-f7="xmark"
                               color="gray"
                               @click="cancelSave(newTag)">
                    </f7-button>
                </template>
            </f7-list-item>
        </f7-list>

        <f7-actions close-by-outside-click close-on-escape :opened="showMoreActionSheet" @actions:closed="showMoreActionSheet = false">
            <f7-actions-group>
                <f7-actions-button @click="setSortable()">{{ tt('Sort') }}</f7-actions-button>
                <f7-actions-button v-if="!showHidden" @click="showHidden = true">{{ tt('Show Hidden Transaction Tags') }}</f7-actions-button>
                <f7-actions-button v-if="showHidden" @click="showHidden = false">{{ tt('Hide Hidden Transaction Tags') }}</f7-actions-button>
            </f7-actions-group>
            <f7-actions-group>
                <f7-actions-button bold close>{{ tt('Cancel') }}</f7-actions-button>
            </f7-actions-group>
        </f7-actions>

        <f7-actions close-by-outside-click close-on-escape :opened="showDeleteActionSheet" @actions:closed="showDeleteActionSheet = false">
            <f7-actions-group>
                <f7-actions-label>{{ tt('Are you sure you want to delete this tag?') }}</f7-actions-label>
                <f7-actions-button color="red" @click="remove(tagToDelete, true)">{{ tt('Delete') }}</f7-actions-button>
            </f7-actions-group>
            <f7-actions-group>
                <f7-actions-button bold close>{{ tt('Cancel') }}</f7-actions-button>
            </f7-actions-group>
        </f7-actions>
    </f7-page>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';
import type { Router } from 'framework7/types';

import { useI18n } from '@/locales/helpers.ts';
import { useI18nUIComponents, showLoading, hideLoading, onSwipeoutDeleted } from '@/lib/ui/mobile.ts';

import { useTransactionTagsStore } from '@/stores/transactionTag.ts';

import { TextDirection } from '@/core/text.ts';
import { TransactionTag } from '@/models/transaction_tag.ts';

import {
    isNoAvailableTag,
    getFirstShowingId,
    getLastShowingId
} from '@/lib/tag.ts';

const props = defineProps<{
    f7router: Router.Router;
}>();

const { tt, getCurrentLanguageTextDirection } = useI18n();
const { showAlert, showToast, routeBackOnError } = useI18nUIComponents();

const transactionTagsStore = useTransactionTagsStore();

const newTag = ref<TransactionTag | null>(null);
const editingTag = ref<TransactionTag>(TransactionTag.createNewTag());
const loading = ref<boolean>(true);
const loadingError = ref<unknown | null>(null);
const showHidden = ref<boolean>(false);
const sortable = ref<boolean>(false);
const tagToDelete = ref<TransactionTag | null>(null);
const showMoreActionSheet = ref<boolean>(false);
const showDeleteActionSheet = ref<boolean>(false);
const displayOrderModified = ref<boolean>(false);
const displayOrderSaving = ref<boolean>(false);

const textDirection = computed<TextDirection>(() => getCurrentLanguageTextDirection());
const tags = computed<TransactionTag[]>(() => transactionTagsStore.allTransactionTags);
const firstShowingId = computed<string | null>(() => getFirstShowingId(tags.value, showHidden.value));
const lastShowingId = computed<string | null>(() => getLastShowingId(tags.value, showHidden.value));
const noAvailableTag = computed<boolean>(() => isNoAvailableTag(tags.value, showHidden.value));
const hasEditingTag = computed<boolean>(() => !!(newTag.value || (editingTag.value.id && editingTag.value.id !== '')));

function isTagModified(tag: TransactionTag): boolean {
    if (tag.id) {
        return editingTag.value.name !== '' && editingTag.value.name !== tag.name;
    } else {
        return tag.name !== '';
    }
}

function getTagDomId(tag: TransactionTag): string {
    return 'tag_' + tag.id;
}

function parseTagIdFromDomId(domId: string): string | null {
    if (!domId || domId.indexOf('tag_') !== 0) {
        return null;
    }

    return domId.substring(4); // tag_
}

function init(): void {
    loading.value = true;

    transactionTagsStore.loadAllTags({
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

function reload(done?: () => void): void {
    if (sortable.value || hasEditingTag.value) {
        done?.();
        return;
    }

    const force = !!done;

    transactionTagsStore.loadAllTags({
        force: force
    }).then(() => {
        done?.();

        if (force) {
            showToast('Tag list has been updated');
        }
    }).catch(error => {
        done?.();

        if (!error.processed) {
            showToast(error.message || error);
        }
    });
}

function add(): void {
    newTag.value = TransactionTag.createNewTag();
}

function edit(tag: TransactionTag): void {
    editingTag.value.id = tag.id;
    editingTag.value.name = tag.name;
}

function save(tag: TransactionTag): void {
    showLoading();

    transactionTagsStore.saveTag({
        tag: tag
    }).then(() => {
        hideLoading();

        if (tag.id) {
            editingTag.value.id = '';
            editingTag.value.name = '';
        } else {
            newTag.value = null;
        }
    }).catch(error => {
        hideLoading();

        if (!error.processed) {
            showToast(error.message || error);
        }
    });
}

function cancelSave(tag: TransactionTag): void {
    if (tag.id) {
        editingTag.value.id = '';
        editingTag.value.name = '';
    } else {
        newTag.value = null;
    }
}

function hide(tag: TransactionTag, hidden: boolean): void {
    showLoading();

    transactionTagsStore.hideTag({
        tag: tag,
        hidden: hidden
    }).then(() => {
        hideLoading();
    }).catch(error => {
        hideLoading();

        if (!error.processed) {
            showToast(error.message || error);
        }
    });
}

function remove(tag: TransactionTag | null, confirm: boolean): void {
    if (!tag) {
        showAlert('An error occurred');
        return;
    }

    if (!confirm) {
        tagToDelete.value = tag;
        showDeleteActionSheet.value = true;
        return;
    }

    showDeleteActionSheet.value = false;
    tagToDelete.value = null;
    showLoading();

    transactionTagsStore.deleteTag({
        tag: tag,
        beforeResolve: (done) => {
            onSwipeoutDeleted(getTagDomId(tag), done);
        }
    }).then(() => {
        hideLoading();
    }).catch(error => {
        hideLoading();

        if (!error.processed) {
            showToast(error.message || error);
        }
    });
}

function setSortable(): void {
    if (sortable.value || hasEditingTag.value) {
        return;
    }

    showHidden.value = true;
    sortable.value = true;
    displayOrderModified.value = false;
}

function saveSortResult(): void {
    if (!displayOrderModified.value) {
        showHidden.value = false;
        sortable.value = false;
        return;
    }

    displayOrderSaving.value = true;
    showLoading();

    transactionTagsStore.updateTagDisplayOrders().then(() => {
        displayOrderSaving.value = false;
        hideLoading();

        showHidden.value = false;
        sortable.value = false;
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
        showToast('Unable to move tag');
        return;
    }

    const id = parseTagIdFromDomId(event.el.id);

    if (!id) {
        showToast('Unable to move tag');
        return;
    }

    transactionTagsStore.changeTagDisplayOrder({
        tagId: id,
        from: event.from,
        to: event.to
    }).then(() => {
        displayOrderModified.value = true;
    }).catch(error => {
        showToast(error.message || error);
    });
}

function onPageAfterIn(): void {
    if (transactionTagsStore.transactionTagListStateInvalid && !loading.value) {
        reload();
    }

    routeBackOnError(props.f7router, loadingError);
}

init();
</script>

<style>
.tag-item-list.list .item-media + .item-inner {
    margin-inline-start: 5px;
}

.transaction-tag-list-item-content {
    overflow: hidden;
    text-overflow: ellipsis;
}
</style>
