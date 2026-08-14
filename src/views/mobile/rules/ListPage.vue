<template>
    <f7-page ptr @ptr:refresh="reload" @page:afterin="onPageAfterIn">
        <f7-navbar :class="{ 'disabled': loading }" :title="tt('Rules')">
            <f7-nav-right>
                <f7-link :class="{ 'disabled': loading }" icon-f7="plus" @click="goToNewRule"></f7-link>
            </f7-nav-right>
        </f7-navbar>

        <f7-list strong inset dividers class="margin-top-half skeleton-text" v-if="loading">
            <f7-list-item title="Rule Group"></f7-list-item>
            <f7-list-item title="Rule Name"></f7-list-item>
            <f7-list-item title="Rule Name"></f7-list-item>
        </f7-list>

        <f7-block v-else-if="allRules.length === 0 && allRuleGroups.length === 0" class="text-align-center">
            <f7-icon f7="wand_stars" size="48" style="color: var(--f7-list-item-footer-text-color)"></f7-icon>
            <p>{{ tt('No rules yet. Tap + to create one and auto-categorize your transactions.') }}</p>
        </f7-block>

        <template v-else>
            <f7-list strong inset dividers v-for="group in allRuleGroups" :key="group.id">
                <f7-list-item :title="group.name" group-title>
                    <div class="group-actions" slot="after">
                        <f7-link icon-f7="plus_circle" @click="goToNewRuleInGroup(group.id)"></f7-link>
                    </div>
                </f7-list-item>
                <f7-list-item
                    v-for="rule in rulesByGroup[group.id]"
                    :key="rule.id"
                    :title="rule.name"
                    :footer="ruleSummary(rule)"
                    @click="editRule(rule)">
                    <f7-icon v-if="!rule.active" slot="media" f7="minus_circle" style="color: var(--f7-list-item-footer-text-color)"></f7-icon>
                    <f7-icon v-else slot="media" f7="checkmark_seal_fill" style="color: var(--f7-color-green)"></f7-icon>
                </f7-list-item>
                <f7-list-button v-if="!rulesByGroup[group.id] || rulesByGroup[group.id]!.length === 0" color="blue" @click="goToNewRuleInGroup(group.id)">
                    {{ tt('Add Rule') }}
                </f7-list-button>
            </f7-list>
        </template>
    </f7-page>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';
import type { Router } from 'framework7/types';

import { useI18n } from '@/locales/helpers.ts';
import { useI18nUIComponents } from '@/lib/ui/mobile.ts';
import { useRulesStore } from '@/stores/rule.ts';

import type { Rule } from '@/models/rule.ts';

const props = defineProps<{
    f7router: Router.Router;
}>();

const { tt } = useI18n();
const { showToast } = useI18nUIComponents();

const rulesStore = useRulesStore();

const loading = ref<boolean>(true);

const allRuleGroups = computed(() => rulesStore.allRuleGroups);
const allRules = computed(() => rulesStore.allRules);

const rulesByGroup = computed<Record<string, Rule[]>>(() => {
    const map: Record<string, Rule[]> = {};
    for (const rule of allRules.value) {
        if (!map[rule.ruleGroupId]) {
            map[rule.ruleGroupId] = [];
        }
        map[rule.ruleGroupId]!.push(rule);
    }
    return map;
});

function onPageAfterIn(): void {
    loadData();
}

function loadData(): void {
    loading.value = true;
    Promise.all([rulesStore.loadRuleGroups(), rulesStore.loadRules()]).then(() => {
        // loaded
    }).catch(error => {
        if (!error.processed) {
            showToast(error.message || error);
        }
    }).finally(() => {
        loading.value = false;
    });
}

function reload(done?: () => void): void {
    Promise.all([rulesStore.loadRuleGroups(), rulesStore.loadRules()]).then(() => {
        done?.();
        showToast('Rules updated');
    }).catch(error => {
        done?.();
        if (!error.processed) {
            showToast(error.message || error);
        }
    });
}

function ruleSummary(rule: Rule): string {
    const triggerCount = rule.triggers.length;
    const actionCount = rule.actions.length;
    return `${triggerCount} ${tt('trigger(s)')}, ${actionCount} ${tt('action(s)')}`;
}

function editRule(rule: Rule): void {
    props.f7router.navigate('/rule/edit', {
        props: { editRuleId: rule.id }
    });
}

function goToNewRule(): void {
    if (allRuleGroups.value.length === 0) {
        showToast(tt('Create a rule group first'));
        return;
    }
    props.f7router.navigate('/rule/edit', {
        props: { presetGroupId: allRuleGroups.value[0]?.id ?? '' }
    });
}

function goToNewRuleInGroup(groupId: string): void {
    props.f7router.navigate('/rule/edit', {
        props: { presetGroupId: groupId }
    });
}
</script>

<style scoped>
.group-actions {
    display: flex;
    gap: 8px;
}
</style>
