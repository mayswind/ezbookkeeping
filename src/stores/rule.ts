// [PLUGIN:rules] Rules engine store.
import { ref } from 'vue';
import { defineStore } from 'pinia';

import { Rule, RuleGroup, type RuleInfoResponse, type RuleGroupInfoResponse } from '@/models/rule.ts';

import services from '@/lib/services.ts';
import logger from '@/lib/logger.ts';

export const useRulesStore = defineStore('rules', () => {
    const allRuleGroups = ref<RuleGroup[]>([]);
    const allRules = ref<Rule[]>([]);

    function loadRuleGroups(): Promise<RuleGroup[]> {
        return new Promise((resolve, reject) => {
            services.getAllRuleGroups().then(response => {
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    reject({ message: 'Unable to retrieve rule groups' });
                    return;
                }

                allRuleGroups.value = data.result.map(RuleGroup.of);
                resolve(allRuleGroups.value);
            }).catch(error => {
                logger.error('failed to load rule groups', error);

                if (error.response && error.response.data && error.response.data.errorMessage) {
                    reject({ error: error.response.data });
                } else if (!error.processed) {
                    reject({ message: 'Unable to retrieve rule groups' });
                } else {
                    reject(error);
                }
            });
        });
    }

    function loadRules(groupId?: string): Promise<Rule[]> {
        return new Promise((resolve, reject) => {
            services.getAllRules({ groupId }).then(response => {
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    reject({ message: 'Unable to retrieve rules' });
                    return;
                }

                allRules.value = data.result.map(Rule.of);
                resolve(allRules.value);
            }).catch(error => {
                logger.error('failed to load rules', error);

                if (error.response && error.response.data && error.response.data.errorMessage) {
                    reject({ error: error.response.data });
                } else if (!error.processed) {
                    reject({ message: 'Unable to retrieve rules' });
                } else {
                    reject(error);
                }
            });
        });
    }

    function saveRuleGroup({ group, isEdit, clientSessionId }: { group: RuleGroup, isEdit: boolean, clientSessionId: string }): Promise<RuleGroup> {
        return new Promise((resolve, reject) => {
            const promise = isEdit
                ? services.modifyRuleGroup(group.toModifyRequest())
                : services.addRuleGroup(group.toCreateRequest(clientSessionId));

            promise.then(response => {
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    reject({ message: isEdit ? 'Unable to save rule group' : 'Unable to add rule group' });
                    return;
                }

                resolve(RuleGroup.of(data.result));
            }).catch(error => {
                logger.error('failed to save rule group', error);

                if (error.response && error.response.data && error.response.data.errorMessage) {
                    reject({ error: error.response.data });
                } else if (!error.processed) {
                    reject({ message: isEdit ? 'Unable to save rule group' : 'Unable to add rule group' });
                } else {
                    reject(error);
                }
            });
        });
    }

    function deleteRuleGroup({ groupId }: { groupId: string }): Promise<boolean> {
        return new Promise((resolve, reject) => {
            services.deleteRuleGroup({ id: groupId }).then(response => {
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    reject({ message: 'Unable to delete rule group' });
                    return;
                }

                resolve(data.result);
            }).catch(error => {
                logger.error('failed to delete rule group', error);

                if (error.response && error.response.data && error.response.data.errorMessage) {
                    reject({ error: error.response.data });
                } else if (!error.processed) {
                    reject({ message: 'Unable to delete rule group' });
                } else {
                    reject(error);
                }
            });
        });
    }

    function saveRule({ rule, isEdit, clientSessionId }: { rule: Rule, isEdit: boolean, clientSessionId: string }): Promise<Rule> {
        return new Promise((resolve, reject) => {
            const promise = isEdit
                ? services.modifyRule(rule.toModifyRequest())
                : services.addRule(rule.toCreateRequest(clientSessionId));

            promise.then(response => {
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    reject({ message: isEdit ? 'Unable to save rule' : 'Unable to add rule' });
                    return;
                }

                resolve(Rule.of(data.result));
            }).catch(error => {
                logger.error('failed to save rule', error);

                if (error.response && error.response.data && error.response.data.errorMessage) {
                    reject({ error: error.response.data });
                } else if (!error.processed) {
                    reject({ message: isEdit ? 'Unable to save rule' : 'Unable to add rule' });
                } else {
                    reject(error);
                }
            });
        });
    }

    function deleteRule({ ruleId }: { ruleId: string }): Promise<boolean> {
        return new Promise((resolve, reject) => {
            services.deleteRule({ id: ruleId }).then(response => {
                const data = response.data;

                if (!data || !data.success || !data.result) {
                    reject({ message: 'Unable to delete rule' });
                    return;
                }

                resolve(data.result);
            }).catch(error => {
                logger.error('failed to delete rule', error);

                if (error.response && error.response.data && error.response.data.errorMessage) {
                    reject({ error: error.response.data });
                } else if (!error.processed) {
                    reject({ message: 'Unable to delete rule' });
                } else {
                    reject(error);
                }
            });
        });
    }

    function reset(): void {
        allRuleGroups.value = [];
        allRules.value = [];
    }

    return {
        allRuleGroups,
        allRules,
        loadRuleGroups,
        loadRules,
        saveRuleGroup,
        deleteRuleGroup,
        saveRule,
        deleteRule,
        reset
    };
});

export type { RuleInfoResponse, RuleGroupInfoResponse };
