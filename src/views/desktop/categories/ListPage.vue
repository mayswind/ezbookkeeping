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
                            ]" v-model="activeCategoryType" @update:modelValue="switchActiveCategoryType" />
                        </div>
                        <v-divider />
                        <v-tabs show-arrows class="my-4" direction="vertical"
                                :disabled="loading" v-model="primaryCategoryId">
                            <v-tab class="tab-text-truncate" value="0" @click="primaryCategoryId = '0'">
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
                                            <v-progress-circular indeterminate size="24" class="ml-2" v-if="loading"></v-progress-circular>
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
                                                                <span class="ml-2">{{ element.name }}</span>
                                                            </div>

                                                            <v-spacer/>

                                                            <v-btn class="hover-display px-2" color="default"
                                                                   density="comfortable" variant="text"
                                                                   :prepend-icon="icons.edit"
                                                                   :disabled="loading || updating"
                                                                   @click="edit(element)">
                                                                {{ $t('Edit') }}
                                                            </v-btn>
                                                            <v-btn class="hover-display px-2 ml-2" color="default"
                                                                   density="comfortable" variant="text"
                                                                   :prepend-icon="element.hidden ? icons.show : icons.hide"
                                                                   :loading="categoryHiding[element.id]"
                                                                   :disabled="loading || updating"
                                                                   @click="hide(element, !element.hidden)">
                                                                {{ element.hidden ? $t('Show') : $t('Hide') }}
                                                            </v-btn>
                                                            <v-btn class="hover-display px-2 ml-2" color="default"
                                                                   density="comfortable" variant="text"
                                                                   :prepend-icon="icons.remove"
                                                                   :loading="categoryRemoving[element.id]"
                                                                   :disabled="loading || updating"
                                                                   @click="remove(element)">
                                                                {{ $t('Delete') }}
                                                            </v-btn>
                                                            <span>
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

    <preset-category-dialog :category-type="activeCategoryType" v-model:show="showPresetDialog"
                            @category:saved="presetCategorySaved" />

    <confirm-dialog ref="confirmDialog"/>
    <snack-bar ref="snackbar" />
</template>

<script>
import PresetCategoryDialog from './list/dialogs/PresetCategoryDialog.vue';

import { useDisplay } from 'vuetify';

import { mapStores } from 'pinia';
import { useTransactionCategoriesStore } from '@/stores/transactionCategory.js';

import categoryConstants from '@/consts/category.js';
import { getOuterHeight } from '@/lib/ui.desktop.js';

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
        PresetCategoryDialog
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

        },
        edit() {

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
        switchActiveCategoryType() {
            this.primaryCategoryId = '0';
            this.updateCardMinHeight();
        },
        isCategorySupportSwitch(category) {
            if (!category || category.hidden) {
                return false;
            }

            return !category.parentId || category.parentId === '' || category.parentId === '0';
        },
        switchPrimaryCategory(category) {
            if (!category || category.hidden) {
                return;
            }

            if (!category.parentId || category.parentId === '' || category.parentId === '0') {
                this.primaryCategoryId = category.id;
            }
        },
        updateCardMinHeight() {
            const self = this

            self.$nextTick(() => {
                if (self.$refs.navbar && self.$refs.navbar.$el && self.$refs.navbar.$el.nextElementSibling) {
                    let navbarHeight = getOuterHeight(self.$refs.navbar.$el.nextElementSibling);

                    if (navbarHeight > self.cardMinHeight) {
                        self.cardMinHeight = navbarHeight;
                    }
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
</style>
