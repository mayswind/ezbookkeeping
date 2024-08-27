<template>
    <f7-sheet swipe-to-close swipe-handler=".swipe-handler"
              :opened="show" :class="{ 'tag-selection-huge-sheet': hugeListItemRows }"
              @sheet:open="onSheetOpen" @sheet:closed="onSheetClosed">
        <f7-toolbar>
            <div class="swipe-handler"></div>
            <div class="left">
                <f7-link sheet-close :text="$t('Cancel')"></f7-link>
            </div>
            <div class="right">
                <f7-link :text="$t('Done')" @click="save"></f7-link>
            </div>
        </f7-toolbar>
        <f7-page-content>
            <f7-list class="no-margin-top no-margin-bottom" v-if="!items || !items.length || noAvailableTag">
                <f7-list-item :title="$t('No available tag')"></f7-list-item>
            </f7-list>
            <f7-list dividers class="no-margin-top no-margin-bottom" v-else-if="items && items.length && !noAvailableTag">
                <f7-list-item checkbox
                              :class="isChecked(item.id) ? 'list-item-selected' : ''"
                              :value="item.id"
                              :checked="isChecked(item.id)"
                              :key="item.id"
                              v-for="item in items"
                              v-show="!item.hidden || isChecked(item.id)"
                              @change="changeItemSelection">
                    <template #title>
                        <f7-block class="no-padding no-margin">
                            <div class="display-flex">
                                <f7-icon f7="number"></f7-icon>
                                <div class="tag-selection-list-item list-item-valign-middle padding-left-half">
                                    {{ item.name }}
                                </div>
                            </div>
                        </f7-block>
                    </template>
                </f7-list-item>
            </f7-list>
        </f7-page-content>
    </f7-sheet>
</template>

<script>
import { copyArrayTo } from '@/lib/common.js';
import { scrollToSelectedItem } from '@/lib/ui.mobile.js';

export default {
    props: [
        'modelValue',
        'items',
        'show'
    ],
    emits: [
        'update:modelValue',
        'update:show'
    ],
    data() {
        const self = this;

        return {
            selectedItemIds: copyArrayTo(self.modelValue, [])
        }
    },
    computed: {
        hugeListItemRows() {
            return this.items.length > 10;
        },
        noAvailableTag() {
            for (let i = 0; i < this.items.length; i++) {
                if (!this.items[i].hidden) {
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
            scrollToSelectedItem(event.$el, '.page-content', 'li.list-item-selected');
        },
        onSheetClosed() {
            this.$emit('update:show', false);
        },
        changeItemSelection(e) {
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
        isChecked(itemId) {
            for (let i = 0; i < this.selectedItemIds.length; i++) {
                if (this.selectedItemIds[i] === itemId) {
                    return true;
                }
            }

            return false;
        }
    }
}
</script>

<style>
@media (min-height: 630px) {
    .tag-selection-huge-sheet {
        height: 400px;
    }
}

.tag-selection-list-item {
    overflow: hidden;
    text-overflow: ellipsis;
}
</style>
