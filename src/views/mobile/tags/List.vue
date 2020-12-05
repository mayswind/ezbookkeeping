<template>
    <f7-page :ptr="!sortable && !hasEditingTag" @ptr:refresh="reload">
        <f7-navbar>
            <f7-nav-left :back-link="$t('Back')"></f7-nav-left>
            <f7-nav-title :title="$t('Transaction Tags')"></f7-nav-title>
            <f7-nav-right class="navbar-compact-icons">
                <f7-link icon-f7="ellipsis" v-if="!sortable && !hasEditingTag && this.tags.length" @click="showMoreActionSheet = true"></f7-link>
                <f7-link icon-f7="plus" v-if="!sortable && !hasEditingTag" @click="add"></f7-link>
                <f7-link :text="$t('Done')" :class="{ 'disabled': displayOrderSaving || hasEditingTag }" v-else-if="sortable" @click="saveSortResult"></f7-link>
            </f7-nav-right>
        </f7-navbar>

        <f7-card class="skeleton-text" v-if="loading">
            <f7-card-content class="no-safe-areas" :padding="false">
                <f7-list>
                    <f7-list-item>
                        <f7-block slot="title" class="no-padding">
                            <div class="display-flex">
                                <f7-icon slot="media" f7="number"></f7-icon>
                                <div class="list-item-valign-middle padding-left-half">Tag Name</div>
                            </div>
                        </f7-block>
                    </f7-list-item>
                    <f7-list-item>
                        <f7-block slot="title" class="no-padding">
                            <div class="display-flex">
                                <f7-icon slot="media" f7="number"></f7-icon>
                                <div class="list-item-valign-middle padding-left-half">Tag Name 2</div>
                            </div>
                        </f7-block>
                    </f7-list-item>
                    <f7-list-item>
                        <f7-block slot="title" class="no-padding">
                            <div class="display-flex">
                                <f7-icon slot="media" f7="number"></f7-icon>
                                <div class="list-item-valign-middle padding-left-half">Tag Name 3</div>
                            </div>
                        </f7-block>
                    </f7-list-item>
                </f7-list>
            </f7-card-content>
        </f7-card>

        <f7-card v-else-if="!loading">
            <f7-card-content class="no-safe-areas" :padding="false">
                <f7-list v-if="noAvailableTag">
                    <f7-list-item :title="$t('No available tag')"></f7-list-item>
                </f7-list>

                <f7-list sortable :sortable-enabled="sortable" @sortable:sort="onSort">
                    <f7-list-item v-for="tag in tags"
                                  :key="tag.id"
                                  :id="tag | tagDomId"
                                  v-show="showHidden || !tag.hidden"
                                  swipeout @taphold.native="setSortable()">
                        <f7-block slot="title" class="no-padding">
                            <div class="display-flex">
                                <f7-icon slot="media" f7="number">
                                    <f7-badge color="gray" class="right-bottom-icon" v-if="tag.hidden">
                                        <f7-icon f7="eye_slash_fill"></f7-icon>
                                    </f7-badge>
                                </f7-icon>

                                <div class="list-item-valign-middle padding-left-half"
                                     v-if="!tag.editing">
                                    {{ tag.name }}
                                </div>
                                <f7-input class="list-title-input padding-left-half"
                                          type="text"
                                          :placeholder="$t('Tag Title')"
                                          :value="tag.newName"
                                          v-else-if="tag.editing"
                                          @input="tag.newName = $event.target.value"
                                          @keyup.enter.native="save(tag)">
                                </f7-input>
                            </div>
                        </f7-block>
                        <f7-button slot="after"
                                   :class="{ 'no-padding': true, 'disabled': !isTagModified(tag) }"
                                   raised fill
                                   icon-f7="checkmark_alt"
                                   color="blue"
                                   v-if="tag.editing"
                                   @click="save(tag)">
                        </f7-button>
                        <f7-button slot="after"
                                   class="no-padding margin-left-half"
                                   raised fill
                                   icon-f7="xmark"
                                   color="gray"
                                   v-if="tag.editing"
                                   @click="cancelSave(tag)">
                        </f7-button>
                        <f7-swipeout-actions left v-if="sortable && !tag.editing">
                            <f7-swipeout-button :color="tag.hidden ? 'blue' : 'gray'" class="padding-left padding-right"
                                                overswipe close @click="hide(tag, !tag.hidden)">
                                <f7-icon :f7="tag.hidden ? 'eye' : 'eye_slash'"></f7-icon>
                            </f7-swipeout-button>
                        </f7-swipeout-actions>
                        <f7-swipeout-actions right v-if="!sortable && !tag.editing">
                            <f7-swipeout-button color="orange" close :text="$t('Edit')" @click="edit(tag)"></f7-swipeout-button>
                            <f7-swipeout-button color="red" class="padding-left padding-right" @click="remove(tag)">
                                <f7-icon f7="trash"></f7-icon>
                            </f7-swipeout-button>
                        </f7-swipeout-actions>
                    </f7-list-item>
                </f7-list>
            </f7-card-content>
        </f7-card>

        <f7-actions close-by-outside-click close-on-escape :opened="showMoreActionSheet" @actions:closed="showMoreActionSheet = false">
            <f7-actions-group>
                <f7-actions-button @click="setSortable()">{{ $t('Sort') }}</f7-actions-button>
            </f7-actions-group>
            <f7-actions-group>
                <f7-actions-button bold close>{{ $t('Cancel') }}</f7-actions-button>
            </f7-actions-group>
        </f7-actions>

        <f7-actions close-by-outside-click close-on-escape :opened="showDeleteActionSheet" @actions:closed="showDeleteActionSheet = false">
            <f7-actions-group>
                <f7-actions-label>{{ $t('Are you sure you want to delete this tag?') }}</f7-actions-label>
                <f7-actions-button color="red" @click="remove(tagToDelete)">{{ $t('Delete') }}</f7-actions-button>
            </f7-actions-group>
            <f7-actions-group>
                <f7-actions-button bold close>{{ $t('Cancel') }}</f7-actions-button>
            </f7-actions-group>
        </f7-actions>
    </f7-page>
</template>

<script>
export default {
    data() {
        return {
            tags: [],
            loading: true,
            showHidden: false,
            sortable: false,
            tagToDelete: null,
            showMoreActionSheet: false,
            showDeleteActionSheet: false,
            editSaving: false,
            displayOrderModified: false,
            displayOrderSaving: false
        };
    },
    computed: {
        noAvailableTag() {
            for (let i = 0; i < this.tags.length; i++) {
                if (this.showHidden || !this.tags[i].hidden) {
                    return false;
                }
            }

            return true;
        },
        hasEditingTag() {
            for (let i = 0; i < this.tags.length; i++) {
                if (this.tags[i].editing) {
                    return true;
                }
            }

            return false;
        }
    },
    created() {
        const self = this;
        const router = self.$f7router;

        self.loading = true;

        self.$services.getAllTransactionTags().then(response => {
            const data = response.data;

            if (!data || !data.success || !data.result) {
                self.$toast('Unable to get tag list');
                router.back();
                return;
            }

            for (let i = 0; i < data.result.length; i++) {
                data.result[i].editing = false;
                data.result[i].newName = data.result[i].name;
            }

            self.tags = data.result;
            self.loading = false;
        }).catch(error => {
            self.$logger.error('failed to load tag list', error);

            if (error.response && error.response.data && error.response.data.errorMessage) {
                self.$toast({ error: error.response.data });
                router.back();
            } else if (!error.processed) {
                self.$toast('Unable to get tag list');
                router.back();
            }
        });
    },
    methods: {
        reload(done) {
            if (this.sortable || this.hasEditingTag) {
                done();
                return;
            }

            const self = this;

            self.$services.getAllTransactionTags().then(response => {
                if (done) {
                    done();
                }

                const data = response.data;

                if (!data || !data.success || !data.result) {
                    self.$toast('Unable to get tag list');
                    return;
                }

                for (let i = 0; i < data.result.length; i++) {
                    data.result[i].editing = false;
                    data.result[i].newName = data.result[i].name;
                }

                self.tags = data.result;
            }).catch(error => {
                self.$logger.error('failed to reload tag list', error);

                if (done) {
                    done();
                }

                if (error.response && error.response.data && error.response.data.errorMessage) {
                    self.$toast({ error: error.response.data });
                } else if (!error.processed) {
                    self.$toast('Unable to get category list');
                }
            });
        },
        setSortable() {
            if (this.sortable || this.hasEditingTag) {
                return;
            }

            this.showHidden = true;
            this.sortable = true;
            this.displayOrderModified = false;
        },
        onSort(event) {
            if (!event || !event.el || !event.el.id || event.el.id.indexOf('tag_') !== 0) {
                this.$toast('Unable to move tag');
                return;
            }

            const id = event.el.id.substr(4); // tag_
            let tag = null;

            for (let i = 0; i < this.tags.length; i++) {
                if (this.tags[i].id === id) {
                    tag = this.tags[i];
                    break;
                }
            }

            if (!tag || !this.tags[event.to]) {
                this.$toast('Unable to move tag');
                return;
            }

            this.tags.splice(event.to, 0, this.tags.splice(event.from, 1)[0]);

            this.displayOrderModified = true;
        },
        saveSortResult() {
            const self = this;
            const newDisplayOrders = [];

            if (!self.displayOrderModified) {
                self.showHidden = false;
                self.sortable = false;
                return;
            }

            self.displayOrderSaving = true;

            for (let i = 0; i < self.tags.length; i++) {
                newDisplayOrders.push({
                    id: self.tags[i].id,
                    displayOrder: i + 1
                });
            }

            self.$showLoading();

            self.$services.moveTransactionTag({
                newDisplayOrders: newDisplayOrders
            }).then(response => {
                self.displayOrderSaving = false;
                self.$hideLoading();

                const data = response.data;

                if (!data || !data.success || !data.result) {
                    self.$toast('Unable to move tag');
                    return;
                }

                self.showHidden = false;
                self.sortable = false;
                self.displayOrderModified = false;
            }).catch(error => {
                self.$logger.error('failed to save tags display order', error);

                self.displayOrderSaving = false;
                self.$hideLoading();

                if (error.response && error.response.data && error.response.data.errorMessage) {
                    self.$toast({ error: error.response.data });
                } else if (!error.processed) {
                    self.$toast('Unable to move tag');
                }
            });
        },
        add() {
            this.tags.push({
                id: '',
                name: '',
                newName: '',
                hidden: false,
                editing: true
            });
        },
        edit(tag) {
            tag.newName = tag.name;
            tag.editing = true;
        },
        save(tag) {
            if (tag.newName === tag.name) {
                return;
            }

            const self = this;

            self.$showLoading();

            let promise = null;

            if (tag.id) {
                promise = self.$services.modifyTransactionTag({
                    id: tag.id,
                    name: tag.newName
                });
            } else {
                promise = self.$services.addTransactionTag({
                    name: tag.newName
                });
            }

            promise.then(response => {
                self.$hideLoading();
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    self.$toast('Unable to save this tag');
                    return;
                }

                tag.id = data.result.id;
                tag.name = data.result.name;
                tag.hidden = data.result.hidden;
                tag.editing = false;
            }).catch(error => {
                self.$logger.error('failed to save tag', error);

                self.$hideLoading();

                if (error.response && error.response.data && error.response.data.errorMessage) {
                    self.$toast({ error: error.response.data });
                } else if (!error.processed) {
                    self.$toast('Unable to save this tag');
                }
            });
        },
        cancelSave(tag) {
            if (tag.id) {
                tag.newName = '';
            } else {
                for (let i = 0; i < this.tags.length; i++) {
                    if (this.tags[i] === tag) {
                        this.tags.splice(i, 1);
                        break;
                    }
                }
            }

            tag.editing = false;
        },
        isTagModified(tag) {
            return tag.newName !== tag.name;
        },
        hide(tag, hidden) {
            const self = this;

            self.$showLoading();

            self.$services.hideTransactionTag({
                id: tag.id,
                hidden: hidden
            }).then(response => {
                self.$hideLoading();
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    if (hidden) {
                        self.$toast('Unable to hide this tag');
                    } else {
                        self.$toast('Unable to unhide this tag');
                    }

                    return;
                }

                tag.hidden = hidden;
            }).catch(error => {
                self.$logger.error('failed to change tag visibility', error);

                self.$hideLoading();

                if (error.response && error.response.data && error.response.data.errorMessage) {
                    self.$toast({ error: error.response.data });
                } else if (!error.processed) {
                    if (hidden) {
                        self.$toast('Unable to hide this tag');
                    } else {
                        self.$toast('Unable to unhide this tag');
                    }
                }
            });
        },
        remove(tag) {
            const self = this;
            const app = self.$f7;
            const $$ = app.$;

            if (!tag) {
                self.$alert('An error has occurred');
                return;
            }

            if (!self.showDeleteActionSheet) {
                self.tagToDelete = tag;
                self.showDeleteActionSheet = true;
                return;
            }

            self.showDeleteActionSheet = false;
            self.tagToDelete = null;
            self.$showLoading();

            self.$services.deleteTransactionTag({
                id: tag.id
            }).then(response => {
                self.$hideLoading();
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    self.$toast('Unable to delete this tag');
                    return;
                }

                app.swipeout.delete($$(`#${self.$options.filters.tagDomId(tag)}`), () => {
                    for (let i = 0; i < self.tags.length; i++) {
                        if (self.tags[i].id === tag.id) {
                            self.tags.splice(i, 1);
                        }
                    }
                });
            }).catch(error => {
                self.$logger.error('failed to delete tag', error);

                self.$hideLoading();

                if (error.response && error.response.data && error.response.data.errorMessage) {
                    self.$toast({ error: error.response.data });
                } else if (!error.processed) {
                    self.$toast('Unable to delete this tag');
                }
            });
        }
    },
    filters: {
        tagDomId(category) {
            return 'tag_' + category.id;
        }
    }
}
</script>
