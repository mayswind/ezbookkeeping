<template>
    <f7-sheet swipe-to-close swipe-handler=".swipe-handler"
              :class="heightClass" :opened="show"
              @sheet:open="onSheetOpen" @sheet:closed="onSheetClosed">
        <f7-toolbar>
            <div class="swipe-handler"></div>
            <div class="left">
                <f7-link sheet-close :text="$t('Cancel')"></f7-link>
            </div>
            <div class="right">
                <f7-link :text="$t('Done')" v-if="allTags && allTags.length && !noAvailableTag" @click="save"></f7-link>
                <f7-link :class="{'disabled': newTag}"
                         :text="$t('Add')" v-if="!allTags || !allTags.length || noAvailableTag" @click="addNewTag"></f7-link>
            </div>
        </f7-toolbar>
        <f7-page-content>
            <f7-list class="no-margin-top no-margin-bottom" v-if="(!allTags || !allTags.length || noAvailableTag) && !newTag">
                <f7-list-item :title="$t('No available tag')"></f7-list-item>
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
                <f7-list-item :title="$t('Add new tag')"
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
                                      :placeholder="$t('Tag Title')"
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

<script>
import { mapStores } from 'pinia';
import { useTransactionTagsStore } from '@/stores/transactionTag.js';

import { copyArrayTo } from '@/lib/common.js';
import { scrollToSelectedItem } from '@/lib/ui.mobile.js';

export default {
    props: [
        'modelValue',
        'allowAddNewTag',
        'show'
    ],
    emits: [
        'update:modelValue',
        'update:show'
    ],
    data() {
        const self = this;
        const transactionTagsStore = useTransactionTagsStore();

        return {
            heightClass: self.getHeightClass(transactionTagsStore.allTransactionTags),
            selectedItemIds: copyArrayTo(self.modelValue, []),
            newTag: null
        }
    },
    computed: {
        ...mapStores(useTransactionTagsStore),
        allTags() {
            return this.transactionTagsStore.allTransactionTags;
        },
        noAvailableTag() {
            for (let i = 0; i < this.allTags.length; i++) {
                if (!this.allTags[i].hidden) {
                    return false;
                }
            }

            return true;
        }
    },
    methods: {
        save() {
            this.$emit('update:modelValue', this.selectedItemIds);
            this.$emit('update:show', false);
        },
        onSheetOpen(event) {
            this.selectedItemIds = copyArrayTo(this.modelValue, []);
            this.newTag = null;
            scrollToSelectedItem(event.$el, '.page-content', 'li.list-item-selected');
        },
        onSheetClosed() {
            this.$emit('update:show', false);
            this.heightClass = this.getHeightClass(this.allTags);
        },
        changeTagSelection(e) {
            const tagId = e.target.value;

            if (e.target.checked) {
                for (let i = 0; i < this.selectedItemIds.length; i++) {
                    if (this.selectedItemIds[i] === tagId) {
                        return;
                    }
                }

                this.selectedItemIds.push(tagId);
            } else {
                for (let i = 0; i < this.selectedItemIds.length; i++) {
                    if (this.selectedItemIds[i] === tagId) {
                        this.selectedItemIds.splice(i, 1);
                        break;
                    }
                }
            }
        },
        addNewTag() {
            this.newTag = {
                name: ''
            };
        },
        saveNewTag() {
            const self = this;

            self.$showLoading();

            self.transactionTagsStore.saveTag({
                tag: self.newTag
            }).then(tag => {
                self.$hideLoading();
                self.newTag = null;

                if (tag && tag.id) {
                    self.selectedItemIds.push(tag.id);
                }
            }).catch(error => {
                self.$hideLoading();

                if (!error.processed) {
                    self.$toast(error.message || error);
                }
            });
        },
        cancelSaveNewTag() {
            this.newTag = null;
        },
        isChecked(itemId) {
            for (let i = 0; i < this.selectedItemIds.length; i++) {
                if (this.selectedItemIds[i] === itemId) {
                    return true;
                }
            }

            return false;
        },
        getHeightClass(allTags) {
            if (allTags && allTags.length > 10) {
                return 'tag-selection-huge-sheet';
            } else if (allTags && allTags.length > 6) {
                return 'tag-selection-large-sheet';
            } else {
                return '';
            }
        }
    }
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
