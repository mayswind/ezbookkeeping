<template>
    <f7-page :ptr="!sortable" @ptr:refresh="reload" @page:afterin="onPageAfterIn">
        <f7-navbar>
            <f7-nav-left :back-link="$t('Back')"></f7-nav-left>
            <f7-nav-title :title="templateType === allTemplateTypes.Schedule ? $t('Scheduled Transactions') : $t('Transaction Templates')"></f7-nav-title>
            <f7-nav-right class="navbar-compact-icons">
                <f7-link icon-f7="ellipsis" :class="{ 'disabled': !templates.length }" v-if="!sortable" @click="showMoreActionSheet = true"></f7-link>
                <f7-link :href="'/template/add?templateType=' + templateType" icon-f7="plus" v-if="!sortable"></f7-link>
                <f7-link :text="$t('Done')" :class="{ 'disabled': displayOrderSaving }" @click="saveSortResult" v-else-if="sortable"></f7-link>
            </f7-nav-right>
        </f7-navbar>

        <f7-list strong inset dividers class="margin-top skeleton-text" v-if="loading">
            <f7-list-item title="Template Name"
                          :key="itemIdx" v-for="itemIdx in [ 1, 2, 3 ]">
                <template #media>
                    <f7-icon f7="app_fill"></f7-icon>
                </template>
            </f7-list-item>
        </f7-list>

        <f7-list strong inset dividers class="margin-top" v-if="!loading && noAvailableTemplate">
            <f7-list-item :title="$t('No available template')"
                          :footer="$t('Once you add templates, you can long press the Add button on the home page to quickly add a new transaction')"
                          v-if="templateType === allTemplateTypes.Normal"></f7-list-item>
            <f7-list-item :title="$t('No available scheduled transactions')" v-else-if="templateType === allTemplateTypes.Schedule"></f7-list-item>
            <f7-list-item :title="$t('No available template')" v-else></f7-list-item>
        </f7-list>

        <f7-list strong inset dividers sortable class="margin-top template-list"
                 :sortable-enabled="sortable"
                 v-if="!loading"
                 @sortable:sort="onSort">
            <f7-list-item swipeout
                          :class="{ 'actual-first-child': template.id === firstShowingId, 'actual-last-child': template.id === lastShowingId }"
                          :id="getTemplateDomId(template)"
                          :title="template.name"
                          :key="template.id"
                          v-for="template in templates"
                          v-show="showHidden || !template.hidden"
                          @taphold="setSortable()">
                <template #media>
                    <f7-icon :f7="templateType === allTemplateTypes.Schedule ? 'clock' : 'doc_plaintext'">
                        <f7-badge color="gray" class="right-bottom-icon" v-if="template.hidden">
                            <f7-icon f7="eye_slash_fill"></f7-icon>
                        </f7-badge>
                    </f7-icon>
                </template>
                <f7-swipeout-actions left v-if="sortable">
                    <f7-swipeout-button :color="template.hidden ? 'blue' : 'gray'" class="padding-left padding-right"
                                        overswipe close @click="hide(template, !template.hidden)">
                        <f7-icon :f7="template.hidden ? 'eye' : 'eye_slash'"></f7-icon>
                    </f7-swipeout-button>
                </f7-swipeout-actions>
                <f7-swipeout-actions right v-if="!sortable">
                    <f7-swipeout-button color="orange" close :text="$t('Edit')" @click="edit(template)"></f7-swipeout-button>
                    <f7-swipeout-button color="red" class="padding-left padding-right" @click="remove(template, false)">
                        <f7-icon f7="trash"></f7-icon>
                    </f7-swipeout-button>
                </f7-swipeout-actions>
            </f7-list-item>
        </f7-list>

        <f7-actions close-by-outside-click close-on-escape :opened="showMoreActionSheet" @actions:closed="showMoreActionSheet = false">
            <f7-actions-group>
                <f7-actions-button @click="setSortable()">{{ $t('Sort') }}</f7-actions-button>
                <f7-actions-button v-if="!showHidden" @click="showHidden = true">{{ $t('Show Hidden Transaction Templates') }}</f7-actions-button>
                <f7-actions-button v-if="showHidden" @click="showHidden = false">{{ $t('Hide Hidden Transaction Templates') }}</f7-actions-button>
            </f7-actions-group>
            <f7-actions-group>
                <f7-actions-button bold close>{{ $t('Cancel') }}</f7-actions-button>
            </f7-actions-group>
        </f7-actions>

        <f7-actions close-by-outside-click close-on-escape :opened="showDeleteActionSheet" @actions:closed="showDeleteActionSheet = false">
            <f7-actions-group>
                <f7-actions-label>{{ $t('Are you sure you want to delete this template?') }}</f7-actions-label>
                <f7-actions-button color="red" @click="remove(templateToDelete, true)">{{ $t('Delete') }}</f7-actions-button>
            </f7-actions-group>
            <f7-actions-group>
                <f7-actions-button bold close>{{ $t('Cancel') }}</f7-actions-button>
            </f7-actions-group>
        </f7-actions>
    </f7-page>
</template>

<script>
import { mapStores } from 'pinia';
import { useTransactionTemplatesStore } from '@/stores/transactionTemplate.js';

import templateConstants from '@/consts/template.js';
import { isDefined } from '@/lib/common.js';
import { onSwipeoutDeleted } from '@/lib/ui.mobile.js';

export default {
    props: [
        'f7route',
        'f7router'
    ],
    data() {
        return {
            templateType: templateConstants.allTemplateTypes.Normal,
            loading: true,
            loadingError: null,
            showHidden: false,
            sortable: false,
            templateToDelete: null,
            showMoreActionSheet: false,
            showDeleteActionSheet: false,
            displayOrderModified: false,
            displayOrderSaving: false
        };
    },
    computed: {
        ...mapStores(useTransactionTemplatesStore),
        templates() {
            return this.transactionTemplatesStore.allTransactionTemplates[this.templateType] || [];
        },
        firstShowingId() {
            for (let i = 0; i < this.templates.length; i++) {
                if (this.showHidden || !this.templates[i].hidden) {
                    return this.templates[i].id;
                }
            }

            return null;
        },
        lastShowingId() {
            for (let i = this.templates.length - 1; i >= 0; i--) {
                if (this.showHidden || !this.templates[i].hidden) {
                    return this.templates[i].id;
                }
            }

            return null;
        },
        noAvailableTemplate() {
            for (let i = 0; i < this.templates.length; i++) {
                if (this.showHidden || !this.templates[i].hidden) {
                    return false;
                }
            }

            return true;
        },
        allTemplateTypes() {
            return templateConstants.allTemplateTypes;
        }
    },
    created() {
        const self = this;

        if (self.f7route.path === '/template/list') {
            self.templateType = templateConstants.allTemplateTypes.Normal;
        } else if (self.f7route.path === '/schedule/list') {
            self.templateType = templateConstants.allTemplateTypes.Schedule;
        }

        self.loading = true;

        self.transactionTemplatesStore.loadAllTemplates({
            templateType: self.templateType,
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
            if ((!isDefined(this.transactionTemplatesStore.transactionTemplateListStatesInvalid[this.templateType]) || this.transactionTemplatesStore.transactionTemplateListStatesInvalid[this.templateType]) && !this.loading) {
                this.reload(null);
            }

            this.$routeBackOnError(this.f7router, 'loadingError');
        },
        reload(done) {
            if (this.sortable) {
                done();
                return;
            }

            const self = this;
            const force = !!done;

            self.transactionTemplatesStore.loadAllTemplates({
                templateType: self.templateType,
                force: force
            }).then(() => {
                if (done) {
                    done();
                }

                if (force) {
                    self.$toast('Template list has been updated');
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
            if (this.sortable) {
                return;
            }

            this.showHidden = true;
            this.sortable = true;
            this.displayOrderModified = false;
        },
        onSort(event) {
            const self = this;

            if (!event || !event.el || !event.el.id) {
                self.$toast('Unable to move template');
                return;
            }

            const id = self.parseTemplateIdFromDomId(event.el.id);

            if (!id) {
                self.$toast('Unable to move template');
                return;
            }

            self.transactionTemplatesStore.changeTemplateDisplayOrder({
                templateType: self.templateType,
                templateId: id,
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

            self.transactionTemplatesStore.updateTemplateDisplayOrders({
                templateType: self.templateType
            }).then(() => {
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
        edit(template) {
            this.f7router.navigate(`/template/edit?id=${template.id}&templateType=${template.templateType}`);
        },
        hide(template, hidden) {
            const self = this;

            self.$showLoading();

            self.transactionTemplatesStore.hideTemplate({
                template: template,
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
        remove(template, confirm) {
            const self = this;

            if (!template) {
                self.$alert('An error occurred');
                return;
            }

            if (!confirm) {
                self.templateToDelete = template;
                self.showDeleteActionSheet = true;
                return;
            }

            self.showDeleteActionSheet = false;
            self.templateToDelete = null;
            self.$showLoading();

            self.transactionTemplatesStore.deleteTemplate({
                template: template,
                beforeResolve: (done) => {
                    onSwipeoutDeleted(self.getTemplateDomId(template), done);
                }
            }).then(() => {
                self.$hideLoading();
            }).catch(error => {
                self.$hideLoading();

                if (!error.processed) {
                    self.$toast(error.message || error);
                }
            });
        },
        getTemplateDomId(template) {
            return 'template_' + template.id;
        },
        parseTemplateIdFromDomId(domId) {
            if (!domId || domId.indexOf('template_') !== 0) {
                return null;
            }

            return domId.substring(9); // template_
        }
    }
};
</script>

<style>
.template-list {
    --f7-list-item-footer-font-size: var(--ebk-large-footer-font-size);
}

.template-list .item-footer {
    padding-top: 4px;
}
</style>
