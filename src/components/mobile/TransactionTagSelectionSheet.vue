<template>
    <f7-sheet ref="sheet" swipe-to-close swipe-handler=".swipe-handler"
              style="height: auto" :opened="show"
              @sheet:open="onSheetOpen" @sheet:closed="onSheetClosed">
        <f7-toolbar class="toolbar-with-swipe-handler">
            <div class="swipe-handler"></div>
            <div class="left">
                <f7-link sheet-close icon-f7="xmark"></f7-link>
            </div>
            <f7-searchbar ref="searchbar" custom-searchs
                          :value="tagSearchContent"
                          :placeholder="tt('Find tag')"
                          :disable-button="false"
                          v-if="enableFilter"
                          @input="tagSearchContent = $event.target.value"
                          @focus="onSearchBarFocus">
            </f7-searchbar>
            <div class="right">
                <f7-button round fill icon-f7="checkmark_alt" @click="save"
                           v-if="filteredTagsWithGroupHeader && filteredTagsWithGroupHeader.length > 0"></f7-button>
                <f7-link icon-f7="plus" :class="{'disabled': newTag}" @click="addNewTag"
                         v-if="!filteredTagsWithGroupHeader || filteredTagsWithGroupHeader.length < 1"></f7-link>
            </div>
        </f7-toolbar>
        <f7-page-content :class="'margin-top ' + heightClass">
            <f7-list class="no-margin-top no-margin-bottom" v-if="(!filteredTagsWithGroupHeader || filteredTagsWithGroupHeader.length < 1) && !newTag">
                <f7-list-item :title="tt('No available tag')"></f7-list-item>
            </f7-list>
            <f7-list dividers class="no-margin-top no-margin-bottom tag-selection-list" v-else-if="(filteredTagsWithGroupHeader && filteredTagsWithGroupHeader.length > 0) || newTag">
                <template :key="(tag instanceof TransactionTag) ? tag.id : (tag.type === 'subheader' ? `${tag.type}-${index}-${tag.title}` : `${tag.type}-${index}`)"
                          v-for="(tag, index) in filteredTagsWithGroupHeader">
                    <f7-list-item group-title v-if="!(tag instanceof TransactionTag) && tag.type === 'subheader'">
                        <div class="tag-selection-list-item">
                            {{ tag.title }}
                        </div>
                    </f7-list-item>
                    <template v-if="!(tag instanceof TransactionTag) && tag.type === 'addbutton'">
                        <f7-list-item link="#" no-chevron
                                      :title="tt('Add new tag')"
                                      v-if="allowAddNewTag && !newTag"
                                      @click="addNewTag()">
                        </f7-list-item>
                        <f7-list-item class="editing-list-item" checkbox indeterminate disabled v-if="allowAddNewTag && newTag">
                            <template #media>
                                <f7-icon class="transaction-tag-icon" f7="number"></f7-icon>
                            </template>
                            <template #title>
                                <div class="display-flex">
                                    <f7-input class="list-title-input padding-inline-start-half"
                                              type="text"
                                              :placeholder="tt('Tag Title')"
                                              v-model:value="newTag.name"
                                              @keyup.enter="saveNewTag()">
                                    </f7-input>
                                </div>
                            </template>
                            <template #after>
                                <f7-button :class="{ 'no-padding': true, 'disabled': !newTag || !newTag.name }"
                                           raised fill
                                           icon-f7="checkmark_alt"
                                           color="blue"
                                           @click="saveNewTag()">
                                </f7-button>
                                <f7-button class="no-padding margin-inline-start-half"
                                           raised fill
                                           icon-f7="xmark"
                                           color="gray"
                                           @click="cancelSaveNewTag()">
                                </f7-button>
                            </template>
                        </f7-list-item>
                    </template>
                    <f7-list-item checkbox
                                  :class="{ 'list-item-selected': selectedTagIds[tag.id], 'disabled': tag.hidden && !selectedTagIds[tag.id] }"
                                  :value="tag.id"
                                  :checked="selectedTagIds[tag.id]"
                                  :key="tag.id"
                                  v-else-if="tag instanceof TransactionTag"
                                  @change="changeTagSelection">
                        <template #media>
                            <f7-icon class="transaction-tag-icon" f7="number">
                                <f7-badge color="gray" class="right-bottom-icon" v-if="tag.hidden">
                                    <f7-icon f7="eye_slash_fill"></f7-icon>
                                </f7-badge>
                            </f7-icon>
                        </template>
                        <template #title>
                            <div class="display-flex">
                                <div class="tag-selection-list-item list-item-valign-middle padding-inline-start-half">
                                    {{ tag.name }}
                                </div>
                            </div>
                        </template>
                    </f7-list-item>
                </template>
            </f7-list>
        </f7-page-content>
    </f7-sheet>
</template>

<script setup lang="ts">
import { ref, useTemplateRef } from 'vue';
import type { Sheet, Searchbar } from 'framework7/types';

import { useI18n } from '@/locales/helpers.ts';
import { useI18nUIComponents, showLoading, hideLoading } from '@/lib/ui/mobile.ts';
import { type CommonTransactionTagSelectionProps, useTransactionTagSelectionBase } from '@/components/base/TransactionTagSelectionBase.ts';

import { useTransactionTagsStore } from '@/stores/transactionTag.ts';

import { TransactionTag } from '@/models/transaction_tag.ts';

import { scrollToSelectedItem } from '@/lib/ui/common.ts';
import { type Framework7Dom, scrollSheetToTop } from '@/lib/ui/mobile.ts';

interface MobileransactionTagSelectionProps extends CommonTransactionTagSelectionProps {
    enableFilter?: boolean;
    show: boolean;
}

const props = defineProps<MobileransactionTagSelectionProps>();

const emit = defineEmits<{
    (e: 'update:modelValue', value: string[]): void;
    (e: 'update:show', value: boolean): void;
}>();

const { tt } = useI18n();
const { showToast } = useI18nUIComponents();

const {
    clonedModelValue,
    tagSearchContent,
    selectedTagIds,
    filteredTagsWithGroupHeader
} = useTransactionTagSelectionBase(props, true, true);

const transactionTagsStore = useTransactionTagsStore();

const sheet = useTemplateRef<Sheet.Sheet>('sheet');
const searchbar = useTemplateRef<Searchbar.Searchbar>('searchbar');

const newTag = ref<TransactionTag | null>(null);
const heightClass = ref<string>(getHeightClass());

function getHeightClass(): string {
    if (filteredTagsWithGroupHeader.value.length > 6) {
        return 'tag-selection-huge-sheet';
    } else if (filteredTagsWithGroupHeader.value.length > 3) {
        return 'tag-selection-large-sheet';
    } else {
        return 'tag-selection-default-sheet';
    }
}

function changeTagSelection(e: Event): void {
    const target = e.target as HTMLInputElement;
    const tagId = target.value;
    const index = clonedModelValue.value.indexOf(tagId);

    if (target.checked) {
        if (index < 0) {
            clonedModelValue.value.push(tagId);
        }
    } else {
        if (index >= 0) {
            clonedModelValue.value.splice(index, 1);
        }
    }
}

function save(): void {
    emit('update:modelValue', clonedModelValue.value);
    emit('update:show', false);
}

function addNewTag(): void {
    newTag.value = TransactionTag.createNewTag();
}

function saveNewTag(): void {
    if (!newTag.value) {
        return;
    }

    showLoading();

    transactionTagsStore.saveTag({
        tag: newTag.value
    }).then(tag => {
        hideLoading();
        newTag.value = null;

        if (tag && tag.id) {
            clonedModelValue.value.push(tag.id);
        }
    }).catch(error => {
        hideLoading();

        if (!error.processed) {
            showToast(error.message || error);
        }
    });
}

function cancelSaveNewTag(): void {
    newTag.value = null;
}

function onSearchBarFocus(): void {
    scrollSheetToTop(sheet.value?.$el as HTMLElement, window.innerHeight); // $el is not Framework7 Dom
}

function onSheetOpen(event: { $el: Framework7Dom }): void {
    clonedModelValue.value = Array.from(props.modelValue);
    newTag.value = null;
    scrollToSelectedItem(event.$el[0], '.sheet-modal-inner', '.page-content', 'li.list-item-selected');
}

function onSheetClosed(): void {
    emit('update:show', false);
    tagSearchContent.value = '';
    searchbar.value?.clear();
}
</script>

<style>
.tag-selection-default-sheet {
    height: 310px;
}

@media (min-height: 630px) {
    .tag-selection-large-sheet {
        height: 370px;
    }

    .tag-selection-huge-sheet {
        height: 500px;
    }
}

@media (max-height: 629px) {
    .tag-selection-large-sheet,
    .tag-selection-huge-sheet {
        height: 320px;
    }
}

.tag-selection-list.list.list-dividers li.list-group-title:first-child,
.tag-selection-list.list.list-dividers li.list-group-title.actual-first-child {
    margin-top: 10px;
    border-radius: inherit;
}

.tag-selection-list.list .item-media + .item-inner {
    margin-inline-start: 0;
}

.tag-selection-list-item {
    overflow: hidden;
    text-overflow: ellipsis;
}
</style>
