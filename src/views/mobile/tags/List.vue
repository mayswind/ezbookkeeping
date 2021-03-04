<template>
    <f7-page :ptr="!sortable && !hasEditingTag" @ptr:refresh="reload" @page:afterin="onPageAfterIn">
        <f7-navbar>
            <f7-nav-left :back-link="$t('Back')"></f7-nav-left>
            <f7-nav-title :title="$t('Transaction Tags')"></f7-nav-title>
            <f7-nav-right class="navbar-compact-icons">
                <f7-link icon-f7="ellipsis" v-if="!sortable && !hasEditingTag && tags.length" @click="showMoreActionSheet = true"></f7-link>
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
                <f7-list v-if="noAvailableTag && !newTag">
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
                                     v-if="editingTag.id !== tag.id">
                                    {{ tag.name }}
                                </div>
                                <f7-input class="list-title-input padding-left-half"
                                          type="text"
                                          :placeholder="$t('Tag Title')"
                                          :value="editingTag.name"
                                          v-else-if="editingTag.id === tag.id"
                                          @input="editingTag.name = $event.target.value"
                                          @keyup.enter.native="save(tag)">
                                </f7-input>
                            </div>
                        </f7-block>
                        <f7-button slot="after"
                                   :class="{ 'no-padding': true, 'disabled': !isTagModified(tag) }"
                                   raised fill
                                   icon-f7="checkmark_alt"
                                   color="blue"
                                   v-if="editingTag.id === tag.id"
                                   @click="save(editingTag)">
                        </f7-button>
                        <f7-button slot="after"
                                   class="no-padding margin-left-half"
                                   raised fill
                                   icon-f7="xmark"
                                   color="gray"
                                   v-if="editingTag.id === tag.id"
                                   @click="cancelSave(editingTag)">
                        </f7-button>
                        <f7-swipeout-actions left v-if="sortable && editingTag.id !== tag.id">
                            <f7-swipeout-button :color="tag.hidden ? 'blue' : 'gray'" class="padding-left padding-right"
                                                overswipe close @click="hide(tag, !tag.hidden)">
                                <f7-icon :f7="tag.hidden ? 'eye' : 'eye_slash'"></f7-icon>
                            </f7-swipeout-button>
                        </f7-swipeout-actions>
                        <f7-swipeout-actions right v-if="!sortable && editingTag.id !== tag.id">
                            <f7-swipeout-button color="orange" close :text="$t('Edit')" @click="edit(tag)"></f7-swipeout-button>
                            <f7-swipeout-button color="red" class="padding-left padding-right" @click="remove(tag, false)">
                                <f7-icon f7="trash"></f7-icon>
                            </f7-swipeout-button>
                        </f7-swipeout-actions>
                    </f7-list-item>

                    <f7-list-item v-if="newTag">
                        <f7-block slot="title" class="no-padding">
                            <div class="display-flex">
                                <f7-icon slot="media" f7="number"></f7-icon>
                                <f7-input class="list-title-input padding-left-half"
                                          type="text"
                                          :placeholder="$t('Tag Title')"
                                          :value="newTag.name"
                                          @input="newTag.name = $event.target.value"
                                          @keyup.enter.native="save(newTag)">
                                </f7-input>
                            </div>
                        </f7-block>
                        <f7-button slot="after"
                                   :class="{ 'no-padding': true, 'disabled': !isTagModified(newTag) }"
                                   raised fill
                                   icon-f7="checkmark_alt"
                                   color="blue"
                                   @click="save(newTag)">
                        </f7-button>
                        <f7-button slot="after"
                                   class="no-padding margin-left-half"
                                   raised fill
                                   icon-f7="xmark"
                                   color="gray"
                                   @click="cancelSave(newTag)">
                        </f7-button>
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
                <f7-actions-button color="red" @click="remove(tagToDelete, true)">{{ $t('Delete') }}</f7-actions-button>
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
            newTag: null,
            editingTag: {
                id: '',
                name: ''
            },
            loading: true,
            loadingError: null,
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
        tags() {
            return this.$store.state.allTransactionTags;
        },
        noAvailableTag() {
            for (let i = 0; i < this.tags.length; i++) {
                if (this.showHidden || !this.tags[i].hidden) {
                    return false;
                }
            }

            return true;
        },
        hasEditingTag() {
            return this.newTag || (this.editingTag.id && this.editingTag.id !== '');
        }
    },
    created() {
        const self = this;

        self.loading = true;

        self.$store.dispatch('loadAllTags', {
            force: false
        }).then(() => {
            self.loading = false;
        }).catch(error => {
            if (error.processed) {
                self.loading = false;
            } else {
                self.loadingError = error;
                self.$toast(error.message || error);
            }
        });
    },
    methods: {
        onPageAfterIn() {
            this.$routeBackOnError('loadingError');
        },
        reload(done) {
            if (this.sortable || this.hasEditingTag) {
                done();
                return;
            }

            const self = this;

            self.$store.dispatch('loadAllTags', {
                force: true
            }).then(() => {
                if (done) {
                    done();
                }
            }).catch(error => {
                if (done) {
                    done();
                }

                if (!error.processed) {
                    self.$toast(error.message || error);
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
            const self = this;

            if (!event || !event.el || !event.el.id || event.el.id.indexOf('tag_') !== 0) {
                self.$toast('Unable to move tag');
                return;
            }

            const id = event.el.id.substr(4); // tag_

            self.$store.dispatch('changeTagDisplayOrder', {
                tagId: id,
                from: event.from,
                to: event.to
            }).then(() => {
                self.displayOrderModified = true;
            }).catch(error => {
                self.$toast(error.message || error);
            });
        },
        saveSortResult() {
            const self = this;

            if (!self.displayOrderModified) {
                self.showHidden = false;
                self.sortable = false;
                return;
            }

            self.displayOrderSaving = true;
            self.$showLoading();

            self.$store.dispatch('updateTagDisplayOrders').then(() => {
                self.displayOrderSaving = false;
                self.$hideLoading();

                self.showHidden = false;
                self.sortable = false;
                self.displayOrderModified = false;
            }).catch(error => {
                self.displayOrderSaving = false;
                self.$hideLoading();

                if (!error.processed) {
                    self.$toast(error.message || error);
                }
            });
        },
        add() {
            this.newTag = {
                name: ''
            };
        },
        edit(tag) {
            this.editingTag.id = tag.id;
            this.editingTag.name = tag.name;
        },
        save(tag) {
            const self = this;

            self.$showLoading();

            self.$store.dispatch('saveTag', {
                tag: tag
            }).then(() => {
                self.$hideLoading();

                if (tag.id) {
                    this.editingTag.id = '';
                    this.editingTag.name = '';
                } else {
                    this.newTag = null;
                }
            }).catch(error => {
                self.$hideLoading();

                if (!error.processed) {
                    self.$toast(error.message || error);
                }
            });
        },
        cancelSave(tag) {
            if (tag.id) {
                this.editingTag.id = '';
                this.editingTag.name = '';
            } else {
                this.newTag = null;
            }
        },
        isTagModified(tag) {
            if (tag.id) {
                return this.editingTag.name !== '' && this.editingTag.name !== tag.name;
            } else {
                return tag.name !== '';
            }
        },
        hide(tag, hidden) {
            const self = this;

            self.$showLoading();

            self.$store.dispatch('hideTag', {
                tag: tag,
                hidden: hidden
            }).then(() => {
                self.$hideLoading();
            }).catch(error => {
                self.$hideLoading();

                if (!error.processed) {
                    self.$toast(error.message || error);
                }
            });
        },
        remove(tag, confirm) {
            const self = this;
            const app = self.$f7;
            const $$ = app.$;

            if (!tag) {
                self.$alert('An error has occurred');
                return;
            }

            if (!confirm) {
                self.tagToDelete = tag;
                self.showDeleteActionSheet = true;
                return;
            }

            self.showDeleteActionSheet = false;
            self.tagToDelete = null;
            self.$showLoading();

            self.$store.dispatch('deleteTag', {
                tag: tag,
                beforeResolve: (done) => {
                    app.swipeout.delete($$(`#${self.$options.filters.tagDomId(tag)}`), () => {
                        done();
                    });
                }
            }).then(() => {
                self.$hideLoading();
            }).catch(error => {
                self.$hideLoading();

                if (!error.processed) {
                    self.$toast(error.message || error);
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
