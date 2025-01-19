<template>
    <f7-sheet swipe-to-close swipe-handler=".swipe-handler"
              :class="heightClass" :opened="show"
              @sheet:open="onSheetOpen" @sheet:closed="onSheetClosed">
        <f7-toolbar>
            <div class="swipe-handler"></div>
            <div class="left">
                <f7-link sheet-close :text="tt('Cancel')"></f7-link>
            </div>
            <div class="right">
                <f7-link :text="tt('Done')" v-if="allTags && allTags.length && !noAvailableTag" @click="save"></f7-link>
                <f7-link :class="{'disabled': newTag}"
                         :text="tt('Add')" v-if="!allTags || !allTags.length || noAvailableTag" @click="addNewTag"></f7-link>
            </div>
        </f7-toolbar>
        <f7-page-content>
            <f7-list class="no-margin-top no-margin-bottom" v-if="(!allTags || !allTags.length || noAvailableTag) && !newTag">
                <f7-list-item :title="tt('No available tag')"></f7-list-item>
            </f7-list>
            <f7-list dividers class="no-margin-top no-margin-bottom tag-selection-list" v-else-if="(allTags && allTags.length && !noAvailableTag) || newTag">
                <f7-list-item checkbox
                              :class="isChecked(tag.id) ? 'list-item-selected' : ''"
                              :value="tag.id"
                              :checked="isChecked(tag.id)"
                              :key="tag.id"
                              v-for="tag in allTags"
                              v-show="!tag.hidden || isChecked(tag.id)"
                              @change="changeTagSelection">
                    <template #title>
                        <f7-block class="no-padding no-margin">
                            <div class="display-flex">
                                <f7-icon f7="number"></f7-icon>
                                <div class="tag-selection-list-item list-item-valign-middle padding-left-half">
                                    {{ tag.name }}
                                </div>
                            </div>
                        </f7-block>
                    </template>
                </f7-list-item>
                <f7-list-item :title="tt('Add new tag')"
                              v-if="allowAddNewTag && !newTag"
                              @click="addNewTag()">
                </f7-list-item>
                <f7-list-item checkbox indeterminate disabled v-if="allowAddNewTag && newTag">
                    <template #media>
                        <f7-icon f7="number"></f7-icon>
                    </template>
                    <template #title>
                        <div class="display-flex">
                            <f7-input class="list-title-input padding-left-half"
                                      type="text"
                                      :placeholder="tt('Tag Title')"
                                      v-model:value="newTag.name"
                                      @keyup.enter="saveNewTag()">
                            </f7-input>
                        </div>
                    </template>
                    <template #after>
                        <f7-button class="no-padding"
                                   raised fill
                                   icon-f7="checkmark_alt"
                                   color="blue"
                                   @click="saveNewTag()">
                        </f7-button>
                        <f7-button class="no-padding margin-left-half"
                                   raised fill
                                   icon-f7="xmark"
                                   color="gray"
                                   @click="cancelSaveNewTag()">
                        </f7-button>
                    </template>
                </f7-list-item>
            </f7-list>
        </f7-page-content>
    </f7-sheet>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';

import { useI18n } from '@/locales/helpers.ts';
import { useI18nUIComponents, showLoading, hideLoading } from '@/lib/ui/mobile.ts';

import { TransactionTag } from '@/models/transaction_tag.ts';
import { useTransactionTagsStore } from '@/stores/transactionTag.ts';

import { copyArrayTo } from '@/lib/common.ts';
import { type Framework7Dom, scrollToSelectedItem } from '@/lib/ui/mobile.ts';

const props = defineProps<{
    modelValue: string[];
    allowAddNewTag?: boolean;
    show: boolean;
}>();

const emit = defineEmits<{
    (e: 'update:modelValue', value: string[]): void;
    (e: 'update:show', value: boolean): void;
}>();

const { tt } = useI18n();
const { showToast } = useI18nUIComponents();

const transactionTagsStore = useTransactionTagsStore();

const selectedItemIds = ref<string[]>(copyArrayTo(props.modelValue, []));
const newTag = ref<TransactionTag | null>(null);
const heightClass = ref<string>(getHeightClass());

const allTags = computed<TransactionTag[]>(() => transactionTagsStore.allTransactionTags);

const noAvailableTag = computed<boolean>(() => {
    if (transactionTagsStore.allTransactionTags) {
        for (let i = 0; i < transactionTagsStore.allTransactionTags.length; i++) {
            if (!transactionTagsStore.allTransactionTags[i].hidden) {
                return false;
            }
        }
    }

    return true;
});

function getHeightClass(): string {
    if (transactionTagsStore.allTransactionTags && transactionTagsStore.allTransactionTags.length > 8) {
        return 'tag-selection-huge-sheet';
    } else if (transactionTagsStore.allTransactionTags && transactionTagsStore.allTransactionTags.length > 4) {
        return 'tag-selection-large-sheet';
    } else {
        return '';
    }
}

function isChecked(itemId: string): boolean {
    return selectedItemIds.value.indexOf(itemId) >= 0;
}

function changeTagSelection(e: Event): void {
    const target = e.target as HTMLInputElement;
    const tagId = target.value;
    const index = selectedItemIds.value.indexOf(tagId);

    if (target.checked) {
        if (index < 0) {
            selectedItemIds.value.push(tagId);
        }
    } else {
        if (index >= 0) {
            selectedItemIds.value.splice(index, 1);
        }
    }
}

function save(): void {
    emit('update:modelValue', selectedItemIds.value);
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
            selectedItemIds.value.push(tag.id);
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

function onSheetOpen(event: { $el: Framework7Dom }): void {
    selectedItemIds.value = copyArrayTo(props.modelValue, []);
    newTag.value = null;
    scrollToSelectedItem(event.$el, '.page-content', 'li.list-item-selected');
}

function onSheetClosed(): void {
    emit('update:show', false);
}
</script>

<style>
@media (min-height: 630px) {
    .tag-selection-large-sheet {
        height: 310px;
    }

    .tag-selection-huge-sheet {
        height: 400px;
    }
}

.tag-selection-list.list .item-media + .item-inner {
    margin-left: 0;
}

.tag-selection-list-item {
    overflow: hidden;
    text-overflow: ellipsis;
}
</style>
