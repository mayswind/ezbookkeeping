import { ref, computed } from 'vue';
import { defineStore } from 'pinia';

import { type BeforeResolveFunction, itemAndIndex } from '@/core/base.ts';
import { TemplateType, type ScheduledTransactionStatus } from '@/core/template.ts';
import { TransactionType } from '@/core/transaction.ts';

import {
    type TransactionTemplateInfoResponse,
    type TransactionTemplateNewDisplayOrderRequest,
    type ScheduledTransactionOccurrenceResponse,
    TransactionTemplate
} from '@/models/transaction_template.ts';

import { isDefined, isEquals } from '@/lib/common.ts';

import logger from '@/lib/logger.ts';
import services, { type ApiResponsePromise } from '@/lib/services.ts';

export const useTransactionTemplatesStore = defineStore('transactionTemplates', () =>{
    let loadingPromise: Promise<TransactionTemplate[]> | null = null;

    const allNormalTransactionTemplates = ref<TransactionTemplate[]>([]);
    const allNormalTransactionTemplatesMap = ref<Record<string, TransactionTemplate>>({});
    const normalTransactionTemplateListStateInvalid = ref<boolean>(true);

    const allVisibleNormalTemplates = computed<TransactionTemplate[]>(() => allNormalTransactionTemplates.value.filter(template => !template.hidden));
    const allAvailableNormalTemplatesCount = computed<number>(() => allNormalTransactionTemplates.value.length);
    const allVisibleNormalTemplatesCount = computed<number>(() => allVisibleNormalTemplates.value.length);

    function loadNormalTransactionTemplateList(templates: TransactionTemplate[]): void {
        allNormalTransactionTemplates.value = templates;
        allNormalTransactionTemplatesMap.value = {};

        for (const template of templates) {
            allNormalTransactionTemplatesMap.value[template.id] = template;
        }
    }

    function addTemplateToNormalTransactionTemplateList(template: TransactionTemplate): void {
        allNormalTransactionTemplates.value.push(template);
        allNormalTransactionTemplatesMap.value[template.id] = template;
    }

    function updateTemplateInNormalTransactionTemplateList(currentTemplate: TransactionTemplate): void {
        for (const [template, index] of itemAndIndex(allNormalTransactionTemplates.value)) {
            if (template.id === currentTemplate.id) {
                allNormalTransactionTemplates.value.splice(index, 1, currentTemplate);
                break;
            }
        }

        allNormalTransactionTemplatesMap.value[currentTemplate.id] = currentTemplate;
    }

    function updateTemplateDisplayOrderInNormalTransactionTemplateList({ from, to }: { from: number, to: number }): void {
        allNormalTransactionTemplates.value.splice(to, 0, allNormalTransactionTemplates.value.splice(from, 1)[0] as TransactionTemplate);
    }

    function updateTemplateVisibilityInNormalTransactionTemplateList({ template, hidden }: { template: TransactionTemplate, hidden: boolean }): void {
        if (allNormalTransactionTemplatesMap.value[template.id]) {
            allNormalTransactionTemplatesMap.value[template.id]!.hidden = hidden;
        }
    }

    function removeTemplateFromNormalTransactionTemplateList(currentTemplate: TransactionTemplate): void {
        for (const [template, index] of itemAndIndex(allNormalTransactionTemplates.value)) {
            if (template.id === currentTemplate.id) {
                allNormalTransactionTemplates.value.splice(index, 1);
                break;
            }
        }

        if (allNormalTransactionTemplatesMap.value[currentTemplate.id]) {
            delete allNormalTransactionTemplatesMap.value[currentTemplate.id];
        }
    }

    function updateNormalTransactionTemplateListInvalidState(invalidState: boolean): void {
        normalTransactionTemplateListStateInvalid.value = invalidState;
    }

    function resetNormalTransactionTemplates(): void {
        loadingPromise = null;
        allNormalTransactionTemplates.value = [];
        allNormalTransactionTemplatesMap.value = {};
        normalTransactionTemplateListStateInvalid.value = true;
    }

    function loadAllNormalTemplates({ force }: { force?: boolean }): Promise<TransactionTemplate[]> {
        if (loadingPromise) {
            return loadingPromise;
        }

        if (!force && !normalTransactionTemplateListStateInvalid.value) {
            return new Promise((resolve) => {
                resolve(allNormalTransactionTemplates.value || []);
            });
        }

        const currentPromise = loadingPromise = new Promise((resolve, reject) => {
            services.getAllTransactionTemplates({
                templateType: TemplateType.Normal.type
            }).then(response => {
                if (!loadingPromise || loadingPromise !== currentPromise) {
                    logger.error('loadingPromise is invalid after retrieving normal template list');
                    reject({ message: 'An error occurred' });
                    return;
                }

                const data = response.data;

                if (!data || !data.success || !data.result) {
                    loadingPromise = null;
                    reject({ message: 'Unable to retrieve template list' });
                    return;
                }

                if (normalTransactionTemplateListStateInvalid.value) {
                    updateNormalTransactionTemplateListInvalidState(false);
                }

                const templates = TransactionTemplate.ofMultiTemplates(data.result);

                if (force && data.result && isEquals(allNormalTransactionTemplates.value, templates)) {
                    loadingPromise = null;
                    reject({ message: 'Template list is up to date', isUpToDate: true });
                    return;
                }

                loadNormalTransactionTemplateList(templates);
                loadingPromise = null;

                resolve(templates);
            }).catch(error => {
                if (loadingPromise === currentPromise) {
                    loadingPromise = null;
                }

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

        return loadingPromise;
    }

    function loadAllScheduledTransactions({ status }: { status: ScheduledTransactionStatus }): Promise<TransactionTemplate[]> {
        return new Promise((resolve, reject) => {
            services.getAllTransactionTemplates({
                templateType: TemplateType.Schedule.type,
                scheduledStatus: status
            }).then(response => {
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    reject({ message: 'Unable to retrieve template list' });
                    return;
                }

                const templates = TransactionTemplate.ofMultiTemplates(data.result);

                resolve(templates);
            }).catch(error => {
                logger.error('failed to load scheduled transaction list', error);

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

    function loadScheduledTransactionsInMonth({ year, month }: { year: number, month: number }): Promise<ScheduledTransactionOccurrenceResponse[]> {
        return services.getScheduledTransactionsInMonth({
            year: year,
            month: month
        }).then(response => {
            const data = response.data;

            if (!data || !data.success || !data.result) {
                return Promise.reject({ message: 'Unable to retrieve scheduled transaction calendar' });
            }

            return data.result;
        }).catch(error => {
            logger.error('failed to load scheduled transaction calendar', error);

            if (error.response && error.response.data && error.response.data.errorMessage) {
                return Promise.reject({ error: error.response.data });
            } else if (!error.processed) {
                return Promise.reject({ message: 'Unable to retrieve scheduled transaction calendar' });
            }

            return Promise.reject(error);
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

                if (template.templateType === TemplateType.Normal.type) {
                    if (!isEdit) {
                        addTemplateToNormalTransactionTemplateList(template);
                    } else {
                        updateTemplateInNormalTransactionTemplateList(template);
                    }
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

    function changeNormalTemplateDisplayOrder({ templateId, from, to }: { templateId: string, from: number, to: number }): Promise<void> {
        return new Promise((resolve, reject) => {
            let currentTemplate: TransactionTemplate | null = null;

            for (const template of allNormalTransactionTemplates.value) {
                if (template.id === templateId) {
                    currentTemplate = template;
                    break;
                }
            }

            if (!currentTemplate || !allNormalTransactionTemplates.value[to]) {
                reject({ message: 'Unable to move template' });
                return;
            }

            if (!normalTransactionTemplateListStateInvalid.value) {
                updateNormalTransactionTemplateListInvalidState(true);
            }

            updateTemplateDisplayOrderInNormalTransactionTemplateList({ from, to });

            resolve();
        });
    }

    function updateNormalTemplateDisplayOrders({ newDisplayOrders: specifiedDisplayOrders }: { newDisplayOrders?: TransactionTemplateNewDisplayOrderRequest[] } = {}): Promise<boolean> {
        const newDisplayOrders: TransactionTemplateNewDisplayOrderRequest[] = specifiedDisplayOrders || [];

        if (!isDefined(specifiedDisplayOrders)) {
            for (const [template, index] of itemAndIndex(allNormalTransactionTemplates.value)) {
                newDisplayOrders.push({
                    id: template.id,
                    displayOrder: index + 1
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

                if (isDefined(specifiedDisplayOrders)) {
                    resolve(data.result);
                    return;
                }

                loadAllNormalTemplates({
                    force: false
                }).finally(() => {
                    if (normalTransactionTemplateListStateInvalid.value) {
                        updateNormalTransactionTemplateListInvalidState(false);
                    }

                    resolve(data.result);
                });
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

                if (template.templateType === TemplateType.Normal.type) {
                    updateTemplateVisibilityInNormalTransactionTemplateList({
                        template: template,
                        hidden: hidden
                    });
                }

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

                if (template.templateType !== TemplateType.Normal.type) {
                    if (beforeResolve) {
                        beforeResolve(() => resolve(data.result));
                    } else {
                        resolve(data.result);
                    }
                    return;
                }

                if (beforeResolve) {
                    beforeResolve(() => {
                        removeTemplateFromNormalTransactionTemplateList(template);
                    });
                } else {
                    removeTemplateFromNormalTransactionTemplateList(template);
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
        allNormalTransactionTemplates,
        allNormalTransactionTemplatesMap,
        normalTransactionTemplateListStateInvalid,
        // computed states
        allVisibleNormalTemplates,
        allAvailableNormalTemplatesCount,
        allVisibleNormalTemplatesCount,
        // functions
        updateNormalTransactionTemplateListInvalidState,
        resetNormalTransactionTemplates,
        loadAllNormalTemplates,
        loadAllScheduledTransactions,
        loadScheduledTransactionsInMonth,
        getTemplate,
        saveTemplateContent,
        changeNormalTemplateDisplayOrder,
        updateNormalTemplateDisplayOrders,
        hideTemplate,
        deleteTemplate
    };
});
