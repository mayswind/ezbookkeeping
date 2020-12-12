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
                <f7-list form>
                    <f7-list-input
                        type="text"
                        clear-button
                        :label="$t('Category Name')"
                        :placeholder="$t('Your category name')"
                        :value="category.name"
                        @input="category.name = $event.target.value"
                    ></f7-list-input>

                    <f7-list-item :header="$t('Category Icon')" key="singleTypeCategoryIconSelection" link="#"
                                  @click="category.showIconSelectionSheet = true">
                        <f7-icon slot="after" :icon="category.icon | categoryIcon" :style="{ color: '#' + category.color }"></f7-icon>
                        <IconSelectionSheet :all-icon-infos="allCategoryIcons"
                                            :show.sync="category.showIconSelectionSheet"
                                            :color="category.color"
                                            v-model="category.icon"
                        ></IconSelectionSheet>
                    </f7-list-item>

                    <f7-list-item :header="$t('Category Color')" key="singleTypeCategoryColorSelection" link="#"
                                  @click="category.showColorSelectionSheet = true">
                        <f7-icon slot="after" f7="app_fill" :style="{ color: '#' + category.color }"></f7-icon>
                        <ColorSelectionSheet :all-color-infos="allCategoryColors"
                                             :show.sync="category.showColorSelectionSheet"
                                             v-model="category.color"
                        ></ColorSelectionSheet>
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
                visible: true,
                showIconSelectionSheet: false,
                showColorSelectionSheet: false
            },
            submitting: false
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
        allCategoryIcons() {
            return this.$constants.icons.allCategoryIcons;
        },
        allCategoryColors() {
            return this.$constants.colors.allCategoryColors;
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
