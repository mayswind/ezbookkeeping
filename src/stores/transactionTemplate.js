import { defineStore } from 'pinia';

import { TransactionType } from '@/core/transaction.ts';
import { TemplateType } from '@/core/template.ts';
import { isDefined, isObject, isArray, isEquals } from '@/lib/common.ts';
import services from '@/lib/services.js';
import logger from '@/lib/logger.ts';

function loadTransactionTemplateList(state, templateType, templates) {
    state.allTransactionTemplates[templateType] = templates;
    state.allTransactionTemplatesMap[templateType] = {};

    for (let i = 0; i < templates.length; i++) {
        const template = templates[i];
        state.allTransactionTemplatesMap[templateType][template.id] = template;
    }
}

function addTemplateToTransactionTemplateList(state, templateType, template) {
    if (isArray(state.allTransactionTemplates[templateType])) {
        state.allTransactionTemplates[templateType].push(template);
    }

    if (isObject(state.allTransactionTemplatesMap[templateType])) {
        state.allTransactionTemplatesMap[templateType][template.id] = template;
    }
}

function updateTemplateInTransactionTemplateList(state, templateType, template) {
    if (isArray(state.allTransactionTemplates[templateType])) {
        for (let i = 0; i < state.allTransactionTemplates[templateType].length; i++) {
            if (state.allTransactionTemplates[templateType][i].id === template.id) {
                state.allTransactionTemplates[templateType].splice(i, 1, template);
                break;
            }
        }
    }

    if (isObject(state.allTransactionTemplatesMap[templateType])) {
        state.allTransactionTemplatesMap[templateType][template.id] = template;
    }
}

function updateTemplateDisplayOrderInTransactionTemplateList(state, templateType, { from, to }) {
    if (isArray(state.allTransactionTemplates[templateType])) {
        state.allTransactionTemplates[templateType].splice(to, 0, state.allTransactionTemplates[templateType].splice(from, 1)[0]);
    }
}

function updateTemplateVisibilityInTransactionTemplateList(state, templateType, { template, hidden }) {
    if (isObject(state.allTransactionTemplatesMap[templateType])) {
        if (state.allTransactionTemplatesMap[templateType][template.id]) {
            state.allTransactionTemplatesMap[templateType][template.id].hidden = hidden;
        }
    }
}

function removeTemplateFromTransactionTemplateList(state, templateType, template) {
    if (isArray(state.allTransactionTemplates[templateType])) {
        for (let i = 0; i < state.allTransactionTemplates[templateType].length; i++) {
            if (state.allTransactionTemplates[templateType][i].id === template.id) {
                state.allTransactionTemplates[templateType].splice(i, 1);
                break;
            }
        }
    }

    if (isObject(state.allTransactionTemplatesMap[templateType])) {
        if (state.allTransactionTemplatesMap[templateType][template.id]) {
            delete state.allTransactionTemplatesMap[templateType][template.id];
        }
    }
}

export const useTransactionTemplatesStore = defineStore('transactionTemplates', {
    state: () => ({
        allTransactionTemplates: {},
        allTransactionTemplatesMap: {},
        transactionTemplateListStatesInvalid: {},
    }),
    getters: {
        allVisibleTemplates(state) {
            const allVisibleTemplates = {};

            for (const templateType in state.allTransactionTemplates) {
                if (!Object.prototype.hasOwnProperty.call(state.allTransactionTemplates, templateType)) {
                    continue;
                }

                const visibleTemplates = [];

                for (let i = 0; i < state.allTransactionTemplates[templateType].length; i++) {
                    const template = state.allTransactionTemplates[templateType][i];

                    if (!template.hidden) {
                        visibleTemplates.push(template);
                    }
                }

                allVisibleTemplates[templateType] = visibleTemplates;
            }

            return allVisibleTemplates;
        },
        allAvailableTemplatesCount(state) {
            const allAvailableTemplateCounts = {};

            for (const templateType in state.allTransactionTemplates) {
                if (!Object.prototype.hasOwnProperty.call(state.allTransactionTemplates, templateType)) {
                    continue;
                }

                allAvailableTemplateCounts[templateType] = state.allTransactionTemplates[templateType].length;
            }

            return allAvailableTemplateCounts;
        },
        allVisibleTemplatesCount(state) {
            const allVisibleTemplateCounts = {};

            for (const templateType in state.allVisibleTemplates) {
                if (!Object.prototype.hasOwnProperty.call(state.allVisibleTemplates, templateType)) {
                    continue;
                }

                allVisibleTemplateCounts[templateType] = state.allVisibleTemplates[templateType].length;
            }

            return allVisibleTemplateCounts;
        }
    },
    actions: {
        updateTransactionTemplateListInvalidState(templateType, invalidState) {
            this.transactionTemplateListStatesInvalid[templateType] = invalidState;
        },
        resetTransactionTemplates() {
            this.allTransactionTemplates = {};
            this.allTransactionTemplatesMap = {};
            this.transactionTemplateListStatesInvalid = {};
        },
        loadAllTemplates({ templateType, force }) {
            const self = this;

            if (!force && isDefined(self.transactionTemplateListStatesInvalid[templateType]) && !self.transactionTemplateListStatesInvalid[templateType]) {
                return new Promise((resolve) => {
                    resolve(self.allTransactionTemplates[templateType] || []);
                });
            }

            return new Promise((resolve, reject) => {
                services.getAllTransactionTemplates({ templateType }).then(response => {
                    const data = response.data;

                    if (!data || !data.success || !data.result) {
                        reject({ message: 'Unable to retrieve template list' });
                        return;
                    }

                    if (!isDefined(self.transactionTemplateListStatesInvalid[templateType]) || self.transactionTemplateListStatesInvalid[templateType]) {
                        self.updateTransactionTemplateListInvalidState(templateType, false);
                    }

                    if (force && data.result && isEquals(self.allTransactionTemplates[templateType], data.result)) {
                        reject({ message: 'Template list is up to date' });
                        return;
                    }

                    loadTransactionTemplateList(self, templateType, data.result);

                    resolve(data.result);
                }).catch(error => {
                    if (force) {
                        logger.error('failed to force load template list', error);
                    } else {
                        logger.error('failed to load template list', error);
                    }

                    if (error.response && error.response.data && error.response.data.errorMessage) {
                        reject({ error: error.response.data });
                    } else if (!error.processed) {
                        reject({ message: 'Unable to retrieve template list' });
                    } else {
                        reject(error);
                    }
                });
            });
        },
        getTemplate({ templateId }) {
            return new Promise((resolve, reject) => {
                services.getTransactionTemplate({
                    id: templateId
                }).then(response => {
                    const data = response.data;

                    if (!data || !data.success || !data.result) {
                        reject({ message: 'Unable to retrieve template' });
                        return;
                    }

                    resolve(data.result);
                }).catch(error => {
                    logger.error('failed to load template info', error);

                    if (error.response && error.response.data && error.response.data.errorMessage) {
                        reject({ error: error.response.data });
                    } else if (!error.processed) {
                        reject({ message: 'Unable to retrieve template' });
                    } else {
                        reject(error);
                    }
                });
            });
        },
        saveTemplateContent({ template, isEdit, clientSessionId }) {
            const self = this;

            const submitTemplate = {
                templateType: template.templateType,
                name: template.name,
                type: template.type,
                sourceAccountId: template.sourceAccountId,
                sourceAmount: template.sourceAmount,
                destinationAccountId: '0',
                destinationAmount: 0,
                hideAmount: template.hideAmount,
                tagIds: template.tagIds,
                comment: template.comment
            };

            if (clientSessionId) {
                submitTemplate.clientSessionId = clientSessionId;
            }

            if (template.templateType === TemplateType.Schedule.type) {
                submitTemplate.scheduledFrequencyType = template.scheduledFrequencyType;
                submitTemplate.scheduledFrequency = template.scheduledFrequency;
                submitTemplate.utcOffset = template.utcOffset;
            }

            if (template.type === TransactionType.Expense) {
                submitTemplate.categoryId = template.expenseCategory;
            } else if (template.type === TransactionType.Income) {
                submitTemplate.categoryId = template.incomeCategory;
            } else if (template.type === TransactionType.Transfer) {
                submitTemplate.categoryId = template.transferCategory;
                submitTemplate.destinationAccountId = template.destinationAccountId;
                submitTemplate.destinationAmount = template.destinationAmount;
            } else {
                return Promise.reject('An error occurred');
            }

            if (isEdit) {
                submitTemplate.id = template.id;
            }

            return new Promise((resolve, reject) => {
                let promise = null;

                if (!submitTemplate.id) {
                    promise = services.addTransactionTemplate(submitTemplate);
                } else {
                    promise = services.modifyTransactionTemplate(submitTemplate);
                }

                promise.then(response => {
                    const data = response.data;

                    if (!data || !data.success || !data.result) {
                        if (!submitTemplate.id) {
                            reject({ message: 'Unable to add template' });
                        } else {
                            reject({ message: 'Unable to save template' });
                        }
                        return;
                    }

                    if (!submitTemplate.id) {
                        addTemplateToTransactionTemplateList(self, template.templateType, data.result);
                    } else {
                        updateTemplateInTransactionTemplateList(self, template.templateType, data.result);
                    }

                    resolve(data.result);
                }).catch(error => {
                    logger.error('failed to save template', error);

                    if (error.response && error.response.data && error.response.data.errorMessage) {
                        reject({ error: error.response.data });
                    } else if (!error.processed) {
                        if (!submitTemplate.id) {
                            reject({ message: 'Unable to add template' });
                        } else {
                            reject({ message: 'Unable to save template' });
                        }
                    } else {
                        reject(error);
                    }
                });
            });
        },
        changeTemplateDisplayOrder({ templateType, templateId, from, to }) {
            const self = this;

            return new Promise((resolve, reject) => {
                let template = null;

                if (!isArray(self.allTransactionTemplates[templateType])) {
                    reject({ message: 'Unable to move template' });
                    return;
                }

                for (let i = 0; i < self.allTransactionTemplates[templateType].length; i++) {
                    if (self.allTransactionTemplates[templateType][i].id === templateId) {
                        template = self.allTransactionTemplates[templateType][i];
                        break;
                    }
                }

                if (!template || !self.allTransactionTemplates[templateType][to]) {
                    reject({ message: 'Unable to move template' });
                    return;
                }

                if (isDefined(self.transactionTemplateListStatesInvalid[templateType]) && !self.transactionTemplateListStatesInvalid[templateType]) {
                    self.updateTransactionTemplateListInvalidState(templateType, true);
                }

                updateTemplateDisplayOrderInTransactionTemplateList(self, templateType, {
                    template: template,
                    from: from,
                    to: to
                });

                resolve();
            });
        },
        updateTemplateDisplayOrders({ templateType }) {
            const self = this;
            const newDisplayOrders = [];

            if (isArray(self.allTransactionTemplates[templateType])) {
                for (let i = 0; i < self.allTransactionTemplates[templateType].length; i++) {
                    newDisplayOrders.push({
                        id: self.allTransactionTemplates[templateType][i].id,
                        displayOrder: i + 1
                    });
                }
            }

            return new Promise((resolve, reject) => {
                services.moveTransactionTemplate({
                    newDisplayOrders: newDisplayOrders
                }).then(response => {
                    const data = response.data;

                    if (!data || !data.success || !data.result) {
                        reject({ message: 'Unable to move template' });
                        return;
                    }

                    if (!isDefined(self.transactionTemplateListStatesInvalid[templateType]) || self.transactionTemplateListStatesInvalid[templateType]) {
                        self.updateTransactionTemplateListInvalidState(templateType, false);
                    }

                    resolve(data.result);
                }).catch(error => {
                    logger.error('failed to save templates display order', error);

                    if (error.response && error.response.data && error.response.data.errorMessage) {
                        reject({ error: error.response.data });
                    } else if (!error.processed) {
                        reject({ message: 'Unable to move template' });
                    } else {
                        reject(error);
                    }
                });
            });
        },
        hideTemplate({ template, hidden }) {
            const self = this;

            return new Promise((resolve, reject) => {
                services.hideTransactionTemplate({
                    id: template.id,
                    hidden: hidden
                }).then(response => {
                    const data = response.data;

                    if (!data || !data.success || !data.result) {
                        if (hidden) {
                            reject({ message: 'Unable to hide this template' });
                        } else {
                            reject({ message: 'Unable to unhide this template' });
                        }

                        return;
                    }

                    updateTemplateVisibilityInTransactionTemplateList(self, template.templateType, {
                        template: template,
                        hidden: hidden
                    });

                    resolve(data.result);
                }).catch(error => {
                    logger.error('failed to change template visibility', error);

                    if (error.response && error.response.data && error.response.data.errorMessage) {
                        reject({ error: error.response.data });
                    } else if (!error.processed) {
                        if (hidden) {
                            reject({ message: 'Unable to hide this template' });
                        } else {
                            reject({ message: 'Unable to unhide this template' });
                        }
                    } else {
                        reject(error);
                    }
                });
            });
        },
        deleteTemplate({ template, beforeResolve }) {
            const self = this;

            return new Promise((resolve, reject) => {
                services.deleteTransactionTemplate({
                    id: template.id
                }).then(response => {
                    const data = response.data;

                    if (!data || !data.success || !data.result) {
                        reject({ message: 'Unable to delete this template' });
                        return;
                    }

                    if (beforeResolve) {
                        beforeResolve(() => {
                            removeTemplateFromTransactionTemplateList(self, template.templateType, template);
                        });
                    } else {
                        removeTemplateFromTransactionTemplateList(self, template.templateType, template);
                    }

                    resolve(data.result);
                }).catch(error => {
                    logger.error('failed to delete template', error);

                    if (error.response && error.response.data && error.response.data.errorMessage) {
                        reject({ error: error.response.data });
                    } else if (!error.processed) {
                        reject({ message: 'Unable to delete this template' });
                    } else {
                        reject(error);
                    }
                });
            });
        }
    }
});
