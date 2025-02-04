import { ref, computed } from 'vue';
import { defineStore } from 'pinia';

import type { BeforeResolveFunction } from '@/core/base.ts';

import { TransactionType } from '@/core/transaction.ts';

import {
    type TransactionTemplateInfoResponse,
    type TransactionTemplateNewDisplayOrderRequest,
    TransactionTemplate
} from '@/models/transaction_template.ts';

import { isDefined, isObject, isArray, isEquals } from '@/lib/common.ts';

import logger from '@/lib/logger.ts';
import services, { type ApiResponsePromise } from '@/lib/services.ts';

export const useTransactionTemplatesStore = defineStore('transactionTemplates', () =>{
    const allTransactionTemplates = ref<Record<number, TransactionTemplate[]>>({});
    const allTransactionTemplatesMap = ref<Record<number, Record<string, TransactionTemplate>>>({});
    const transactionTemplateListStatesInvalid = ref<Record<number, boolean>>({});

    const allVisibleTemplates = computed<Record<number, TransactionTemplate[]>>(() => {
        const allVisibleTemplates: Record<number, TransactionTemplate[]> = {};

        for (const templateType in allTransactionTemplates.value) {
            if (!Object.prototype.hasOwnProperty.call(allTransactionTemplates.value, templateType)) {
                continue;
            }

            const allTemplates = allTransactionTemplates.value[templateType];
            const visibleTemplates: TransactionTemplate[] = [];

            for (let i = 0; i < allTemplates.length; i++) {
                const template = allTemplates[i];

                if (!template.hidden) {
                    visibleTemplates.push(template);
                }
            }

            allVisibleTemplates[templateType] = visibleTemplates;
        }

        return allVisibleTemplates;
    });

    const allAvailableTemplatesCount = computed<Record<number, number>>(() => {
        const allAvailableTemplateCounts: Record<number, number> = {};

        for (const templateType in allTransactionTemplates.value) {
            if (!Object.prototype.hasOwnProperty.call(allTransactionTemplates.value, templateType)) {
                continue;
            }

            allAvailableTemplateCounts[templateType] = allTransactionTemplates.value[templateType].length;
        }

        return allAvailableTemplateCounts;
    });

    const allVisibleTemplatesCount = computed<Record<number, number>>(() => {
        const allVisibleTemplateCounts: Record<number, number> = {};

        for (const templateType in allVisibleTemplates.value) {
            if (!Object.prototype.hasOwnProperty.call(allVisibleTemplates.value, templateType)) {
                continue;
            }

            allVisibleTemplateCounts[templateType] = allVisibleTemplates.value[templateType].length;
        }

        return allVisibleTemplateCounts;
    });

    function loadTransactionTemplateList(templateType: number, templates: TransactionTemplate[]): void {
        allTransactionTemplates.value[templateType] = templates;
        allTransactionTemplatesMap.value[templateType] = {};

        for (let i = 0; i < templates.length; i++) {
            const template = templates[i];
            allTransactionTemplatesMap.value[templateType][template.id] = template;
        }
    }

    function addTemplateToTransactionTemplateList(templateType: number, template: TransactionTemplate): void {
        const templates = allTransactionTemplates.value[templateType];
        const templateMap = allTransactionTemplatesMap.value[templateType];

        if (isArray(templates)) {
            templates.push(template);
        }

        if (isObject(templateMap)) {
            templateMap[template.id] = template;
        }
    }

    function updateTemplateInTransactionTemplateList(templateType: number, template: TransactionTemplate): void {
        const templates = allTransactionTemplates.value[templateType];
        const templateMap = allTransactionTemplatesMap.value[templateType];

        if (isArray(templates)) {
            for (let i = 0; i < templates.length; i++) {
                if (templates[i].id === template.id) {
                    templates.splice(i, 1, template);
                    break;
                }
            }
        }

        if (isObject(templateMap)) {
            templateMap[template.id] = template;
        }
    }

    function updateTemplateDisplayOrderInTransactionTemplateList(templateType: number, { from, to }: { from: number, to: number }): void {
        const templates = allTransactionTemplates.value[templateType];

        if (isArray(templates)) {
            templates.splice(to, 0, templates.splice(from, 1)[0]);
        }
    }

    function updateTemplateVisibilityInTransactionTemplateList(templateType: number, { template, hidden }: { template: TransactionTemplate, hidden: boolean }): void {
        const templateMap = allTransactionTemplatesMap.value[templateType];

        if (isObject(templateMap)) {
            if (templateMap[template.id]) {
                templateMap[template.id].hidden = hidden;
            }
        }
    }

    function removeTemplateFromTransactionTemplateList(templateType: number, template: TransactionTemplate): void {
        const templates = allTransactionTemplates.value[templateType];
        const templateMap = allTransactionTemplatesMap.value[templateType];

        if (isArray(templates)) {
            for (let i = 0; i < templates.length; i++) {
                if (templates[i].id === template.id) {
                    templates.splice(i, 1);
                    break;
                }
            }
        }

        if (isObject(templateMap)) {
            if (templateMap[template.id]) {
                delete templateMap[template.id];
            }
        }
    }

    function updateTransactionTemplateListInvalidState(templateType: number, invalidState: boolean): void {
        transactionTemplateListStatesInvalid.value[templateType] = invalidState;
    }

    function resetTransactionTemplates(): void {
        allTransactionTemplates.value = {};
        allTransactionTemplatesMap.value = {};
        transactionTemplateListStatesInvalid.value = {};
    }

    function loadAllTemplates({ templateType, force }: { templateType: number, force?: boolean }): Promise<TransactionTemplate[]> {
        if (!force && isDefined(transactionTemplateListStatesInvalid.value[templateType]) && !transactionTemplateListStatesInvalid.value[templateType]) {
            return new Promise((resolve) => {
                resolve(allTransactionTemplates.value[templateType] || []);
            });
        }

        return new Promise((resolve, reject) => {
            services.getAllTransactionTemplates({ templateType }).then(response => {
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    reject({ message: 'Unable to retrieve template list' });
                    return;
                }

                if (!isDefined(transactionTemplateListStatesInvalid.value[templateType]) || transactionTemplateListStatesInvalid.value[templateType]) {
                    updateTransactionTemplateListInvalidState(templateType, false);
                }

                const templates = TransactionTemplate.ofManyTemplates(data.result);

                if (force && data.result && isEquals(allTransactionTemplates.value[templateType], templates)) {
                    reject({ message: 'Template list is up to date', isUpToDate: true });
                    return;
                }

                loadTransactionTemplateList(templateType, templates);

                resolve(templates);
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
    }

    function getTemplate({ templateId }: { templateId: string }): Promise<TransactionTemplate> {
        return new Promise((resolve, reject) => {
            services.getTransactionTemplate({
                id: templateId
            }).then(response => {
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    reject({ message: 'Unable to retrieve template' });
                    return;
                }

                const template = TransactionTemplate.ofTemplate(data.result);

                resolve(template);
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
    }

    function saveTemplateContent({ template, isEdit, clientSessionId }: { template: TransactionTemplate, isEdit: boolean, clientSessionId: string }): Promise<TransactionTemplate> {
        return new Promise((resolve, reject) => {
            let promise: ApiResponsePromise<TransactionTemplateInfoResponse>;

            if (template.type !== TransactionType.Expense &&
                template.type !== TransactionType.Income &&
                template.type !== TransactionType.Transfer) {
                reject({ message: 'An error occurred' });
                return;
            }

            if (!isEdit) {
                promise = services.addTransactionTemplate(template.toTemplateCreateRequest(clientSessionId));
            } else {
                promise = services.modifyTransactionTemplate(template.toTemplateModifyRequest());
            }

            promise.then(response => {
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    if (!isEdit) {
                        reject({ message: 'Unable to add template' });
                    } else {
                        reject({ message: 'Unable to save template' });
                    }
                    return;
                }

                const template = TransactionTemplate.ofTemplate(data.result);

                if (!isEdit) {
                    addTemplateToTransactionTemplateList(template.templateType, template);
                } else {
                    updateTemplateInTransactionTemplateList(template.templateType, template);
                }

                resolve(template);
            }).catch(error => {
                logger.error('failed to save template', error);

                if (error.response && error.response.data && error.response.data.errorMessage) {
                    reject({ error: error.response.data });
                } else if (!error.processed) {
                    if (!isEdit) {
                        reject({ message: 'Unable to add template' });
                    } else {
                        reject({ message: 'Unable to save template' });
                    }
                } else {
                    reject(error);
                }
            });
        });
    }

    function changeTemplateDisplayOrder({ templateType, templateId, from, to }: { templateType: number, templateId: string, from: number, to: number }): Promise<void> {
        return new Promise((resolve, reject) => {
            let template: TransactionTemplate | null = null;

            if (!isArray(allTransactionTemplates.value[templateType])) {
                reject({ message: 'Unable to move template' });
                return;
            }

            for (let i = 0; i < allTransactionTemplates.value[templateType].length; i++) {
                if (allTransactionTemplates.value[templateType][i].id === templateId) {
                    template = allTransactionTemplates.value[templateType][i];
                    break;
                }
            }

            if (!template || !allTransactionTemplates.value[templateType][to]) {
                reject({ message: 'Unable to move template' });
                return;
            }

            if (isDefined(transactionTemplateListStatesInvalid.value[templateType]) && !transactionTemplateListStatesInvalid.value[templateType]) {
                updateTransactionTemplateListInvalidState(templateType, true);
            }

            updateTemplateDisplayOrderInTransactionTemplateList(templateType, { from, to });

            resolve();
        });
    }

    function updateTemplateDisplayOrders({ templateType }: { templateType: number }): Promise<boolean> {
        const newDisplayOrders: TransactionTemplateNewDisplayOrderRequest[] = [];

        if (isArray(allTransactionTemplates.value[templateType])) {
            for (let i = 0; i < allTransactionTemplates.value[templateType].length; i++) {
                newDisplayOrders.push({
                    id: allTransactionTemplates.value[templateType][i].id,
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

                if (!isDefined(transactionTemplateListStatesInvalid.value[templateType]) || transactionTemplateListStatesInvalid.value[templateType]) {
                    updateTransactionTemplateListInvalidState(templateType, false);
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
    }

    function hideTemplate({ template, hidden }: { template: TransactionTemplate, hidden: boolean }): Promise<boolean> {
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

                updateTemplateVisibilityInTransactionTemplateList(template.templateType, {
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
    }

    function deleteTemplate({ template, beforeResolve }: { template: TransactionTemplate, beforeResolve?: BeforeResolveFunction }): Promise<boolean> {
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
                        removeTemplateFromTransactionTemplateList(template.templateType, template);
                    });
                } else {
                    removeTemplateFromTransactionTemplateList(template.templateType, template);
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

    return {
        // states
        allTransactionTemplates,
        allTransactionTemplatesMap,
        transactionTemplateListStatesInvalid,
        // computed states
        allVisibleTemplates,
        allAvailableTemplatesCount,
        allVisibleTemplatesCount,
        // functions
        updateTransactionTemplateListInvalidState,
        resetTransactionTemplates,
        loadAllTemplates,
        getTemplate,
        saveTemplateContent,
        changeTemplateDisplayOrder,
        updateTemplateDisplayOrders,
        hideTemplate,
        deleteTemplate
    };
});
