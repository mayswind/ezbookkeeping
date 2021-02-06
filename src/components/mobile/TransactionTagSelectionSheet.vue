<template>
    <f7-sheet :class="{ 'tag-selection-huge-sheet': hugeListItemRows }" :opened="show" @sheet:open="onSheetOpen" @sheet:closed="onSheetClosed">
        <f7-toolbar>
            <div class="left">
                <f7-link sheet-close :text="$t('Cancel')"></f7-link>
            </div>
            <div class="right">
                <f7-link :text="$t('Done')" @click="save"></f7-link>
            </div>
        </f7-toolbar>
        <f7-page-content>
            <f7-list no-hairlines class="no-margin-top no-margin-bottom" v-if="!items || !items.length || noAvailableTag">
                <f7-list-item :title="$t('No available tag')"></f7-list-item>
            </f7-list>
            <f7-list no-hairlines class="no-margin-top no-margin-bottom" v-else-if="items && items.length && !noAvailableTag">
                <f7-list-item checkbox v-for="item in items"
                              v-show="!item.hidden"
                              :key="item.id"
                              :value="item.id"
                              :checked="item.id | isChecked(selectedItemIds)"
                              @change="changeItemSelection">
                    <f7-block slot="title" class="no-padding no-margin">
                        <div class="display-flex">
                            <f7-icon slot="media" f7="number"></f7-icon>
                            <div class="list-item-valign-middle padding-left-half">
                                {{ item.name }}
                            </div>
                        </div>
                    </f7-block>
                </f7-list-item>
            </f7-list>
        </f7-page-content>
    </f7-sheet>
</template>

<script>
export default {
    props: [
        'value',
        'items',
        'show'
    ],
    data() {
        const self = this;

        return {
            selectedItemIds: self.$utilities.copyArrayTo(self.value, [])
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
            this.$emit('input', this.selectedItemIds);
            this.$emit('update:show', false);
        },
        onSheetOpen() {
            this.selectedItemIds = this.$utilities.copyArrayTo(this.value, []);
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
        }
    },
    filters: {
        isChecked(itemId, selectedItemIds) {
            for (let i = 0; i < selectedItemIds.length; i++) {
                if (selectedItemIds[i] === itemId) {
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
</style>
self
