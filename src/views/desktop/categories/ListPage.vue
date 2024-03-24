<template>
    <v-row class="match-height">
        <v-col cols="12">
            <v-card>
                <v-layout>
                    <v-navigation-drawer ref="navbar" :permanent="alwaysShowNav" v-model="showNav">
                        <div class="mx-6 my-4">
                            <btn-vertical-group :disabled="loading" :buttons="[
                                { name: $t('Expense'), value: allCategoryTypes.Expense },
                                { name: $t('Income'), value: allCategoryTypes.Income },
                                { name: $t('Transfer'), value: allCategoryTypes.Transfer }
                            ]" v-model="activeCategoryType" @update:model-value="switchAllPrimaryCategories" />
                        </div>
                        <v-divider />
                        <v-tabs show-arrows class="my-4" direction="vertical"
                                :disabled="loading" v-model="primaryCategoryId">
                            <v-tab class="tab-text-truncate" value="0" @click="switchAllPrimaryCategories">
                                <span class="text-truncate">{{ $t('Primary Categories') }}</span>
                            </v-tab>
                            <template :key="category.id" v-for="category in primaryCategories">
                                <v-tab class="tab-text-truncate" :value="category.id" v-if="!category.hidden"
                                       @click="switchPrimaryCategory(category)">
                                    <span class="text-truncate">{{ category.name }}</span>
                                </v-tab>
                            </template>
                            <template v-if="loading && (!primaryCategories || primaryCategories.length < 1)">
                                <v-skeleton-loader class="skeleton-no-margin mx-5 mt-4 mb-3" type="text"
                                                   :key="itemIdx" :loading="true" v-for="itemIdx in [ 1, 2, 3, 4, 5, 6, 7, 8, 9, 10 ]"></v-skeleton-loader>
                            </template>
                        </v-tabs>
                    </v-navigation-drawer>
                    <v-main>
                        <v-window class="d-flex flex-grow-1 disable-tab-transition w-100-window-container" v-model="activeTab">
                            <v-window-item value="categoryPage">
                                <v-card variant="flat" :min-height="cardMinHeight">
                                    <template #title>
                                        <div class="title-and-toolbar d-flex align-center">
                                            <v-btn class="mr-3 d-md-none" density="compact" color="default" variant="plain"
                                                   :ripple="false" :icon="true" @click="showNav = !showNav">
                                                <v-icon :icon="icons.menu" size="24" />
                                            </v-btn>
                                            <span>{{ $t('Transaction Categories') }}</span>
                                            <v-btn class="ml-3" color="default" variant="outlined"
                                                   :disabled="loading || updating" @click="add">{{ $t('Add') }}</v-btn>
                                            <v-btn class="ml-3" color="primary" variant="tonal"
                                                   :disabled="loading || updating" @click="saveSortResult"
                                                   v-if="displayOrderModified">{{ $t('Save Display Order') }}</v-btn>
                                            <v-btn density="compact" color="default" variant="text"
                                                   class="ml-2" :icon="true" :disabled="loading || updating"
                                                   v-if="!loading" @click="reload">
                                                <v-icon :icon="icons.refresh" size="24" />
                                                <v-tooltip activator="parent">{{ $t('Refresh') }}</v-tooltip>
                                            </v-btn>
                                            <v-progress-circular indeterminate size="20" class="ml-3" v-if="loading"></v-progress-circular>
                                            <v-spacer/>
                                            <v-btn density="comfortable" color="default" variant="text" class="ml-2"
                                                   :disabled="loading || updating" :icon="true">
                                                <v-icon :icon="icons.more" />
                                                <v-menu activator="parent">
                                                    <v-list>
                                                        <v-list-item :prepend-icon="icons.show"
                                                                     :title="$t('Show Hidden Transaction Category')"
                                                                     v-if="!showHidden" @click="showHidden = true"></v-list-item>
                                                        <v-list-item :prepend-icon="icons.hide"
                                                                     :title="$t('Hide Hidden Transaction Category')"
                                                                     v-if="showHidden" @click="showHidden = false"></v-list-item>
                                                    </v-list>
                                                </v-menu>
                                            </v-btn>
                                        </div>
                                    </template>

                                    <v-table class="transaction-category-table table-striped" :hover="!loading">
                                        <thead>
                                        <tr>
                                            <th class="text-uppercase">
                                                <div class="d-flex align-center">
                                                    <span>{{ $t('Category Name') }}</span>
                                                    <v-spacer/>
                                                    <span>{{ $t('Operation') }}</span>
                                                </div>
                                            </th>
                                        </tr>
                                        </thead>

                                        <tbody v-if="loading && noAvailableCategory">
                                        <tr :key="itemIdx" v-for="itemIdx in [ 1, 2, 3, 4, 5, 6, 7, 8, 9, 10 ]">
                                            <td class="px-0">
                                                <v-skeleton-loader type="text" :loading="true"></v-skeleton-loader>
                                            </td>
                                        </tr>
                                        </tbody>

                                        <tbody v-if="!loading && noAvailableCategory">
                                        <tr>
                                            <td>
                                                <div class="d-flex align-center">
                                                    <span>{{ $t('No available category') }}</span>
                                                    <v-btn class="ml-3" color="default" variant="outlined"
                                                           @click="showPresetDialog = true"
                                                           v-if="hasSubCategories && noCategory">
                                                        {{ $t('Add Default Categories') }}
                                                    </v-btn>
                                                </div>
                                            </td>
                                        </tr>
                                        </tbody>

                                        <draggable-list tag="tbody"
                                                        item-key="id"
                                                        handle=".drag-handle"
                                                        ghost-class="dragging-item"
                                                        :disabled="noAvailableCategory"
                                                        v-model="categories"
                                                        @change="onMove">
                                            <template #item="{ element }">
                                                <tr class="transaction-category-table-row text-sm" v-if="showHidden || !element.hidden">
                                                    <td>
                                                        <div class="d-flex align-center">
                                                            <div class="d-flex align-center" :class="{ 'cursor-pointer': isCategorySupportSwitch(element) }"
                                                                 @click="switchPrimaryCategory(element)">
                                                                <ItemIcon icon-type="category"
                                                                          :icon-id="element.icon" :color="element.color"
                                                                          :hidden-status="element.hidden" />
                                                                <div class="d-flex flex-column py-2">
                                                                    <span class="ml-2">{{ element.name }}</span>
                                                                    <span class="transaction-category-comment ml-2">{{ element.comment }}</span>
                                                                </div>
                                                            </div>

                                                            <v-spacer/>

                                                            <v-btn class="px-2 ml-2" color="default"
                                                                   density="comfortable" variant="text"
                                                                   :class="{ 'd-none': loading, 'hover-display': !loading }"
                                                                   :prepend-icon="element.hidden ? icons.show : icons.hide"
                                                                   :loading="categoryHiding[element.id]"
                                                                   :disabled="loading || updating"
                                                                   @click="hide(element, !element.hidden)">
                                                                {{ element.hidden ? $t('Show') : $t('Hide') }}
                                                            </v-btn>
                                                            <v-btn class="px-2" color="default"
                                                                   density="comfortable" variant="text"
                                                                   :class="{ 'd-none': loading, 'hover-display': !loading }"
                                                                   :prepend-icon="icons.edit"
                                                                   :disabled="loading || updating"
                                                                   @click="edit(element)">
                                                                {{ $t('Edit') }}
                                                            </v-btn>
                                                            <v-btn class="px-2" color="default"
                                                                   density="comfortable" variant="text"
                                                                   :class="{ 'd-none': loading, 'hover-display': !loading }"
                                                                   :prepend-icon="icons.remove"
                                                                   :loading="categoryRemoving[element.id]"
                                                                   :disabled="loading || updating"
                                                                   @click="remove(element)">
                                                                {{ $t('Delete') }}
                                                            </v-btn>
                                                            <span class="ml-2">
                                                                <v-icon :class="!loading && !updating && availableCategoryCount > 1 ? 'drag-handle' : 'disabled'"
                                                                        :icon="icons.drag"/>
                                                                <v-tooltip activator="parent" v-if="!loading && !updating && availableCategoryCount > 1">{{ $t('Drag and Drop to Change Order') }}</v-tooltip>
                                                            </span>
                                                        </div>
                                                    </td>
                                                </tr>
                                            </template>
                                        </draggable-list>
                                    </v-table>
                                </v-card>
                            </v-window-item>
                        </v-window>
                    </v-main>
                </v-layout>
            </v-card>
        </v-col>
    </v-row>

    <edit-dialog ref="editDialog" :persistent="true" />
    <preset-dialog :category-type="activeCategoryType" v-model:show="showPresetDialog"
                            @category:saved="presetCategorySaved" />

    <confirm-dialog ref="confirmDialog"/>
    <snack-bar ref="snackbar" />
</template>

<script>
import EditDialog from './list/dialogs/EditDialog.vue';
import PresetDialog from './list/dialogs/PresetDialog.vue';

import { useDisplay } from 'vuetify';

import { mapStores } from 'pinia';
import { useTransactionCategoriesStore } from '@/stores/transactionCategory.js';

import categoryConstants from '@/consts/category.js';
import { getNavSideBarOuterHeight } from '@/lib/ui.desktop.js';

import {
    mdiRefresh,
    mdiMenu,
    mdiPencilOutline,
    mdiEyeOffOutline,
    mdiEyeOutline,
    mdiDeleteOutline,
    mdiDrag,
    mdiDotsVertical
} from '@mdi/js';

export default {
    components: {
        EditDialog,
        PresetDialog
    },
    data() {
        const { mdAndUp } = useDisplay();

        return {
            activeCategoryType: categoryConstants.allCategoryTypes.Expense,
            activeTab: 'categoryPage',
            primaryCategoryId: '0',
            loading: true,
            updating: false,
            categoryHiding: {},
            categoryRemoving: {},
            displayOrderModified: false,
            cardMinHeight: 680,
            alwaysShowNav: mdAndUp.value,
            showNav: mdAndUp.value,
            showHidden: false,
            showPresetDialog: false,
            icons: {
                refresh: mdiRefresh,
                menu: mdiMenu,
                edit: mdiPencilOutline,
                show: mdiEyeOutline,
                hide: mdiEyeOffOutline,
                remove: mdiDeleteOutline,
                drag: mdiDrag,
                more: mdiDotsVertical
            }
        };
    },
    computed: {
        ...mapStores(useTransactionCategoriesStore),
        allCategoryTypes() {
            return categoryConstants.allCategoryTypes;
        },
        primaryCategories() {
            if (!this.transactionCategoriesStore.allTransactionCategories || !this.transactionCategoriesStore.allTransactionCategories[this.activeCategoryType]) {
                return [];
            }

            return this.transactionCategoriesStore.allTransactionCategories[this.activeCategoryType];
        },
        secondaryCategories() {
            if (!this.transactionCategoriesStore.allTransactionCategoriesMap || !this.transactionCategoriesStore.allTransactionCategoriesMap[this.primaryCategoryId]) {
                return [];
            }

            return this.transactionCategoriesStore.allTransactionCategoriesMap[this.primaryCategoryId].subCategories;
        },
        hasSubCategories() {
            return !this.primaryCategoryId || this.primaryCategoryId === '' || this.primaryCategoryId === '0';
        },
        categories() {
            if (this.hasSubCategories) {
                return this.primaryCategories;
            } else {
                return this.secondaryCategories;
            }
        },
        noAvailableCategory() {
            for (let i = 0; i < this.categories.length; i++) {
                if (this.showHidden || !this.categories[i].hidden) {
                    return false;
                }
            }

            return true;
        },
        noCategory() {
            return this.categories.length < 1;
        },
        availableCategoryCount() {
            let count = 0;

            for (let i = 0; i < this.categories.length; i++) {
                if (this.showHidden || !this.categories[i].hidden) {
                    count++;
                }
            }

            return count;
        }
    },
    created() {
        this.reload(false);
    },
    setup() {
        const display = useDisplay();

        return {
            display: display
        };
    },
    watch: {
        'display.mdAndUp.value': function (newValue) {
            this.alwaysShowNav = newValue;

            if (!this.showNav) {
                this.showNav = newValue;
            }
        }
    },
    methods: {
        reload(force) {
            const self = this;

            self.loading = true;

            self.transactionCategoriesStore.loadAllCategories({
                force: force
            }).then(() => {
                self.loading = false;
                self.displayOrderModified = false;

                if (force) {
                    self.$refs.snackbar.showMessage('Category list has been updated');
                }

                self.updateCardMinHeight();
            }).catch(error => {
                self.loading = false;

                if (!error.processed) {
                    self.$refs.snackbar.showError(error);
                }
            });
        },
        onMove(event) {
            if (!event || !event.moved) {
                return;
            }

            const self = this;
            const moveEvent = event.moved;

            if (!moveEvent.element || !moveEvent.element.id) {
                self.$refs.snackbar.showMessage('Unable to move category');
                return;
            }

            self.transactionCategoriesStore.changeCategoryDisplayOrder({
                categoryId: moveEvent.element.id,
                from: moveEvent.oldIndex,
                to: moveEvent.newIndex
            }).then(() => {
                self.displayOrderModified = true;
            }).catch(error => {
                self.$refs.snackbar.showError(error);
            });
        },
        saveSortResult() {
            const self = this;

            if (!self.displayOrderModified) {
                return;
            }

            self.loading = true;

            self.transactionCategoriesStore.updateCategoryDisplayOrders({
                type: self.activeCategoryType,
                parentId: self.primaryCategoryId,
            }).then(() => {
                self.loading = false;
                self.displayOrderModified = false;
            }).catch(error => {
                self.loading = false;

                if (!error.processed) {
                    self.$refs.snackbar.showError(error);
                }
            });
        },
        add() {
            const self = this;

            self.$refs.editDialog.open({
                type: self.activeCategoryType,
                parentId: self.primaryCategoryId
            }).then(result => {
                if (result && result.message) {
                    self.$refs.snackbar.showMessage(result.message);
                }

                self.updateCardMinHeight();
            }).catch(error => {
                if (error) {
                    self.$refs.snackbar.showError(error);
                }
            });
        },
        edit(category) {
            const self = this;

            self.$refs.editDialog.open({
                id: category.id,
                currentCategory: category
            }).then(result => {
                if (result && result.message) {
                    self.$refs.snackbar.showMessage(result.message);
                }

                self.updateCardMinHeight();
            }).catch(error => {
                if (error) {
                    self.$refs.snackbar.showError(error);
                }
            });
        },
        hide(category, hidden) {
            const self = this;

            self.updating = true;
            self.categoryHiding[category.id] = true;

            self.transactionCategoriesStore.hideCategory({
                category: category,
                hidden: hidden
            }).then(() => {
                self.updating = false;
                self.categoryHiding[category.id] = false;

                self.updateCardMinHeight();
            }).catch(error => {
                self.updating = false;
                self.categoryHiding[category.id] = false;

                if (!error.processed) {
                    self.$refs.snackbar.showError(error);
                }
            });
        },
        remove(category) {
            const self = this;

            self.$refs.confirmDialog.open('Are you sure you want to delete this category?').then(() => {
                self.updating = true;
                self.categoryRemoving[category.id] = true;

                self.transactionCategoriesStore.deleteCategory({
                    category: category
                }).then(() => {
                    self.updating = false;
                    self.categoryRemoving[category.id] = false;

                    self.updateCardMinHeight();
                }).catch(error => {
                    self.updating = false;
                    self.categoryRemoving[category.id] = false;

                    if (!error.processed) {
                        self.$refs.snackbar.showError(error);
                    }
                });
            });
        },
        presetCategorySaved(e) {
            if (e && e.message) {
                this.$refs.snackbar.showMessage(e.message);
                this.reload(false);
            }
        },
        isCategorySupportSwitch(category) {
            if (!category || category.hidden) {
                return false;
            }

            return !category.parentId || category.parentId === '' || category.parentId === '0';
        },
        switchAllPrimaryCategories() {
            this.primaryCategoryId = '0';
            this.updateCardMinHeight();
        },
        switchPrimaryCategory(category) {
            if (!category || category.hidden) {
                return;
            }

            if (!category.parentId || category.parentId === '' || category.parentId === '0') {
                this.primaryCategoryId = category.id;
            }

            this.updateCardMinHeight();
        },
        updateCardMinHeight() {
            const self = this

            self.$nextTick(() => {
                if (self.$refs.navbar && self.$refs.navbar.$el && self.$refs.navbar.$el.nextElementSibling) {
                    let navbarHeight = getNavSideBarOuterHeight(self.$refs.navbar.$el.nextElementSibling);
                    self.cardMinHeight = Math.max(navbarHeight, 680);
                }
            });
        }
    }
}
</script>

<style>
.transaction-category-table tr.transaction-category-table-row .hover-display {
    display: none;
}

.transaction-category-table tr.transaction-category-table-row:hover .hover-display {
    display: grid;
}

.transaction-category-table .transaction-category-comment {
    font-size: 0.8rem;
    color: rgba(var(--v-theme-on-background), var(--v-medium-emphasis-opacity)) !important;
}
</style>
