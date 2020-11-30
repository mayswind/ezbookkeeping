<template>
    <f7-page>
        <f7-navbar>
            <f7-nav-left :back-link="$t('Back')"></f7-nav-left>
            <f7-nav-title :title="$t(title)"></f7-nav-title>
            <f7-nav-right>
                <f7-link :class="{ 'disabled': inputIsEmpty || submitting }" :text="$t(saveButtonTitle)" @click="save"></f7-link>
            </f7-nav-right>
        </f7-navbar>

        <f7-card class="skeleton-text" v-if="loading">
            <f7-card-content class="no-safe-areas" :padding="false">
                <f7-list>
                    <f7-list-input label="Category Name" placeholder="Your category name"></f7-list-input>
                    <f7-list-item header="Category Icon" after="Icon"></f7-list-item>
                    <f7-list-item header="Category Color" after="Color"></f7-list-item>
                    <f7-list-input type="textarea" label="Description" placeholder="Your category description (optional)"></f7-list-input>
                </f7-list>
            </f7-card-content>
        </f7-card>

        <f7-card v-else-if="!loading">
            <f7-card-content class="no-safe-areas" :padding="false">
                <f7-list>
                    <f7-list-input
                        type="text"
                        clear-button
                        :label="$t('Category Name')"
                        :placeholder="$t('Your category name')"
                        :value="category.name"
                        @input="category.name = $event.target.value"
                    ></f7-list-input>

                    <f7-list-item :header="$t('Category Icon')" key="singleTypeCategoryIconSelection" link="#"
                                  @click="showIconSelectionSheet(category)">
                        <f7-icon slot="after" :icon="category.icon | categoryIcon" :style="{ color: '#' + category.color }"></f7-icon>
                    </f7-list-item>

                    <f7-list-item :header="$t('Category Color')" key="singleTypeCategoryColorSelection" link="#"
                                  @click="showColorSelectionSheet(category)">
                        <f7-icon slot="after" f7="app_fill" :style="{ color: '#' + category.color }"></f7-icon>
                    </f7-list-item>

                    <f7-list-input
                        type="textarea"
                        :label="$t('Description')"
                        :placeholder="$t('Your category description (optional)')"
                        :value="category.comment"
                        @input="category.comment = $event.target.value"
                    ></f7-list-input>

                    <f7-list-item :header="$t('Visible')" v-if="editCategoryId">
                        <f7-toggle :checked="category.visible" @toggle:change="category.visible = $event"></f7-toggle>
                    </f7-list-item>
                </f7-list>
            </f7-card-content>
        </f7-card>

        <f7-sheet class="category-icon-sheet" :opened="showIconSelection" @sheet:closed="hideIconSelectionSheet">
            <f7-toolbar>
                <div class="left"></div>
                <div class="right">
                    <f7-link sheet-close :text="$t('Done')"></f7-link>
                </div>
            </f7-toolbar>
            <f7-page-content>
                <f7-block class="margin-vertical">
                    <f7-row class="padding-vertical-half padding-horizontal-half" v-for="(row, idx) in allCategoryIconRows" :key="idx">
                        <f7-col class="text-align-center" v-for="categoryIcon in row" :key="categoryIcon.id">
                            <f7-icon :icon="categoryIcon.icon" :style="{ color: '#' + (categoryChoosingIcon ? categoryChoosingIcon.color : '000000') }" @click.native="setSelectedIcon(categoryIcon)">
                                <f7-badge color="default" class="right-bottom-icon" v-if="categoryChoosingIcon && categoryChoosingIcon.icon === categoryIcon.id">
                                    <f7-icon f7="checkmark_alt"></f7-icon>
                                </f7-badge>
                            </f7-icon>
                        </f7-col>
                        <f7-col v-for="idx in (iconCountPerRow - row.length)" :key="idx"></f7-col>
                    </f7-row>
                </f7-block>
            </f7-page-content>
        </f7-sheet>

        <f7-sheet :opened="showColorSelection" @sheet:closed="hideColorSelectionSheet">
            <f7-toolbar>
                <div class="left"></div>
                <div class="right">
                    <f7-link sheet-close :text="$t('Done')"></f7-link>
                </div>
            </f7-toolbar>
            <f7-page-content>
                <f7-block class="margin-vertical">
                    <f7-row class="padding-vertical padding-horizontal-half" v-for="(row, idx) in allCategoryColorRows" :key="idx">
                        <f7-col class="text-align-center" v-for="categoryColor in row" :key="categoryColor.color">
                            <f7-icon f7="app_fill" :style="{ color: '#' + categoryColor.color }" @click.native="setSelectedColor(categoryColor.color)">
                                <f7-badge color="default" class="right-bottom-icon" v-if="categoryChoosingColor && categoryChoosingColor.color === categoryColor.color">
                                    <f7-icon f7="checkmark_alt"></f7-icon>
                                </f7-badge>
                            </f7-icon>
                        </f7-col>
                        <f7-col v-for="idx in (iconCountPerRow - row.length)" :key="idx"></f7-col>
                    </f7-row>
                </f7-block>
            </f7-page-content>
        </f7-sheet>
    </f7-page>
</template>

<script>
export default {
    data() {
        const self = this;
        const query = self.$f7route.query;

        return {
            editCategoryId: null,
            loading: false,
            category: {
                type: query.type,
                name: '',
                parentId: query.parentId,
                icon: self.$constants.icons.defaultCategoryIconId,
                color: self.$constants.colors.defaultCategoryColor,
                comment: '',
                visible: true
            },
            iconCountPerRow: 7,
            categoryChoosingIcon: null,
            categoryChoosingColor: null,
            submitting: false,
            showIconSelection: false,
            showColorSelection: false
        };
    },
    computed: {
        title() {
            if (!this.editCategoryId) {
                if (this.category.parentId === '0') {
                    return 'Add Primary Category';
                } else {
                    return 'Add Secondary Category';
                }
            } else {
                return 'Edit Category';
            }
        },
        saveButtonTitle() {
            if (!this.editCategoryId) {
                return 'Add';
            } else {
                return 'Save';
            }
        },
        allCategoryIconRows() {
            const allCategoryIcons = this.$constants.icons.allCategoryIcons;
            const ret = [];
            let rowCount = 0;

            for (let categoryIconId in allCategoryIcons) {
                if (!Object.prototype.hasOwnProperty.call(allCategoryIcons, categoryIconId)) {
                    continue;
                }

                const categoryIcon = allCategoryIcons[categoryIconId];

                if (!ret[rowCount]) {
                    ret[rowCount] = [];
                } else if (ret[rowCount] && ret[rowCount].length >= this.iconCountPerRow) {
                    rowCount++;
                    ret[rowCount] = [];
                }

                ret[rowCount].push({
                    id: categoryIconId,
                    icon: categoryIcon.icon
                });
            }

            return ret;
        },
        allCategoryColorRows() {
            const allCategoryColors = this.$constants.colors.allCategoryColors;
            const ret = [];
            let rowCount = -1;

            for (let i = 0; i < allCategoryColors.length; i++) {
                if (i % this.iconCountPerRow === 0) {
                    ret[++rowCount] = [];
                }

                ret[rowCount].push({
                    color: allCategoryColors[i]
                });
            }

            return ret;
        },
        inputIsEmpty() {
            return !!this.inputEmptyProblemMessage;
        },
        inputEmptyProblemMessage() {
            if (!this.category.name) {
                return 'Category name cannot be empty';
            } else {
                return null;
            }
        }
    },
    created() {
        const self = this;
        const query = self.$f7route.query;
        const router = self.$f7router;

        if (!query.id && !query.parentId) {
            self.$toast('Parameter Invalid');
            router.back();
            return;
        }

        if (query.id) {
            self.loading = true;

            self.editCategoryId = query.id;
            self.$services.getTransactionCategory({
                id: self.editCategoryId
            }).then(response => {
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    self.$toast('Unable to get category');
                    router.back();
                    return;
                }

                const category = data.result;
                self.category.id = category.id;
                self.category.type = category.type.toString();
                self.category.parentId = category.type.parentId;
                self.category.name = category.name;
                self.category.icon = category.icon;
                self.category.color = category.color;
                self.category.comment = category.comment;
                self.category.visible = !category.hidden;

                self.loading = false;
            }).catch(error => {
                self.$logger.error('failed to load category info', error);

                if (error.response && error.response.data && error.response.data.errorMessage) {
                    self.$toast({ error: error.response.data });
                    router.back();
                } else if (!error.processed) {
                    self.$toast('Unable to get category');
                    router.back();
                }
            });
        } else if (query.parentId) {
            if (query.type !== '1' && query.type !== '2' && query.type !== '3') {
                self.$toast('Parameter Invalid');
                router.back();
                return;
            }

            self.loading = false;
        }
    },
    methods: {
        showIconSelectionSheet(category) {
            this.categoryChoosingIcon = category;
            this.showIconSelection = true;
        },
        setSelectedIcon(categoryIcon) {
            if (!this.categoryChoosingIcon) {
                return;
            }

            this.categoryChoosingIcon.icon = categoryIcon.id;
            this.categoryChoosingIcon = null;
            this.showIconSelection = false;
        },
        hideIconSelectionSheet() {
            this.categoryChoosingIcon = null;
            this.showIconSelection = false;
        },
        showColorSelectionSheet(category) {
            this.categoryChoosingColor = category;
            this.showColorSelection = true;
        },
        setSelectedColor(color) {
            if (!this.categoryChoosingColor) {
                return;
            }

            this.categoryChoosingColor.color = color;
            this.categoryChoosingColor = null;
            this.showColorSelection = false;
        },
        hideColorSelectionSheet() {
            this.categoryChoosingColor = null;
            this.showColorSelection = false;
        },
        save() {
            const self = this;
            const router = self.$f7router;

            const problemMessage = self.inputEmptyProblemMessage;

            if (problemMessage) {
                self.$alert(problemMessage);
                return;
            }

            self.submitting = true;
            self.$showLoading(() => self.submitting);

            const submitCategory = {
                type: parseInt(self.category.type),
                name: self.category.name,
                parentId: self.category.parentId,
                icon: self.category.icon,
                color: self.category.color,
                comment: self.category.comment
            };

            let promise = null;

            if (!self.editCategoryId) {
                promise = self.$services.addTransactionCategory(submitCategory);
            } else {
                submitCategory.id = self.category.id;
                submitCategory.hidden = !self.category.visible;
                promise = self.$services.modifyTransactionCategory(submitCategory);
            }

            promise.then(response => {
                self.submitting = false;
                self.$hideLoading();
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    if (!self.editCategoryId) {
                        self.$toast('Unable to add category');
                    } else {
                        self.$toast('Unable to save category');
                    }
                    return;
                }

                if (!self.editCategoryId) {
                    self.$toast('You have added a new category');
                } else {
                    self.$toast('You have saved this category');
                }

                router.back();
            }).catch(error => {
                self.$logger.error('failed to save category', error);

                self.submitting = false;
                self.$hideLoading();

                if (error.response && error.response.data && error.response.data.errorMessage) {
                    self.$toast({ error: error.response.data });
                } else if (!error.processed) {
                    if (!self.editCategoryId) {
                        self.$toast('Unable to add category');
                    } else {
                        self.$toast('Unable to save category');
                    }
                }
            });
        }
    }
}
</script>

<style>
@media (min-height: 630px) {
    .category-icon-sheet {
        height: 400px;
    }
}
</style>
