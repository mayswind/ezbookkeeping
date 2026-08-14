<template>
    <f7-page>
        <f7-navbar :back-link="tt('Back')" :title="isEdit ? tt('Edit Rule') : tt('New Rule')"></f7-navbar>

        <f7-list strong inset dividers v-if="!loading">
            <f7-list-input :label="tt('Rule Name')" type="text" :placeholder="tt('Rule name')" :value="rule.name" @input="rule.name = $event.target.value"></f7-list-input>
            <f7-list-input :label="tt('Rule Group')" v-if="!isEdit || ruleGroups.length > 1">
                <select v-model="rule.ruleGroupId">
                    <option v-for="g in ruleGroups" :key="g.id" :value="g.id">{{ g.name }}</option>
                </select>
            </f7-list-input>
            <f7-list-item>
                <span>{{ tt('Active') }}</span>
                <f7-toggle slot="after" :checked="rule.active" @toggle:change="rule.active = $event"></f7-toggle>
            </f7-list-item>
        </f7-list>

        <!-- Triggers -->
        <f7-block-title v-if="!loading">{{ tt('When transaction matches') }}</f7-block-title>
        <f7-list strong inset dividers v-if="!loading">
            <f7-list-item>
                <span>{{ tt('Match') }}</span>
                <f7-segmented slot="after">
                    <f7-button :active="!rule.strict" small @click="rule.strict = false">{{ tt('Any') }}</f7-button>
                    <f7-button :active="rule.strict" small @click="rule.strict = true">{{ tt('All') }}</f7-button>
                </f7-segmented>
            </f7-list-item>
            <div v-for="(trigger, idx) in rule.triggers" :key="'t' + idx">
                <f7-list-item>
                    <div class="row-input">
                        <select class="trigger-type-select" v-model="trigger.triggerType">
                            <option v-for="t in triggerTypes" :key="t.value" :value="t.value">{{ tt(t.label) }}</option>
                        </select>
                        <f7-input v-if="needsValue(trigger.triggerType)" type="text" :placeholder="tt('Value')" :value="trigger.triggerValue" @input="trigger.triggerValue = $event.target.value"></f7-input>
                    </div>
                    <f7-link slot="after" icon-f7="minus_circle" color="red" v-if="rule.triggers.length > 1" @click="rule.triggers.splice(idx, 1)"></f7-link>
                </f7-list-item>
            </div>
            <f7-list-button color="blue" @click="rule.triggers.push(new RuleTrigger())">{{ tt('Add Trigger') }}</f7-list-button>
        </f7-list>

        <!-- Actions -->
        <f7-block-title v-if="!loading">{{ tt('Then apply') }}</f7-block-title>
        <f7-list strong inset dividers v-if="!loading">
            <div v-for="(action, idx) in rule.actions" :key="'a' + idx">
                <f7-list-item>
                    <div class="row-input">
                        <select class="action-type-select" v-model="action.actionType">
                            <option v-for="a in actionTypes" :key="a.value" :value="a.value">{{ tt(a.label) }}</option>
                        </select>
                        <f7-input v-if="needsActionValue(action.actionType)" type="text" :placeholder="tt('Value')" :value="action.actionValue" @input="action.actionValue = $event.target.value"></f7-input>
                    </div>
                    <f7-link slot="after" icon-f7="minus_circle" color="red" v-if="rule.actions.length > 1" @click="rule.actions.splice(idx, 1)"></f7-link>
                </f7-list-item>
            </div>
            <f7-list-button color="blue" @click="rule.actions.push(new RuleAction())">{{ tt('Add Action') }}</f7-list-button>
        </f7-list>

        <f7-block v-if="!loading">
            <f7-button fill color="green" :class="{ 'disabled': !canSave }" @click="save">{{ isEdit ? tt('Save') : tt('Create') }}</f7-button>
        </f7-block>
        <f7-block v-if="!loading && isEdit">
            <f7-button fill color="red" @click="confirmDelete">{{ tt('Delete Rule') }}</f7-button>
        </f7-block>
    </f7-page>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import type { Router } from 'framework7/types';

import { useI18n } from '@/locales/helpers.ts';
import { useI18nUIComponents } from '@/lib/ui/mobile.ts';
import { useRulesStore } from '@/stores/rule.ts';

import { Rule, RuleTrigger, RuleAction, type RuleTriggerType, type RuleActionType } from '@/models/rule.ts';

const props = defineProps<{
    f7router: Router.Router;
    editRuleId?: string;
    presetGroupId?: string;
}>();

const { tt } = useI18n();
const { showToast, showConfirm } = useI18nUIComponents();

const rulesStore = useRulesStore();

const loading = ref<boolean>(true);
const rule = ref<Rule>(new Rule());

const isEdit = computed(() => !!props.editRuleId);
const ruleGroups = computed(() => rulesStore.allRuleGroups);

const triggerTypes: { value: RuleTriggerType; label: string }[] = [
    { value: 'description_contains', label: 'Description contains' },
    { value: 'description_is', label: 'Description is' },
    { value: 'amount_is', label: 'Amount is' },
    { value: 'amount_less', label: 'Amount is less than' },
    { value: 'amount_more', label: 'Amount is more than' },
    { value: 'source_account_is', label: 'Source account is' },
    { value: 'destination_account_is', label: 'Destination account is' },
    { value: 'category_is', label: 'Category is' },
    { value: 'has_no_category', label: 'Has no category' },
    { value: 'has_any_category', label: 'Has any category' }
];

const actionTypes: { value: RuleActionType; label: string }[] = [
    { value: 'set_category', label: 'Set category to' },
    { value: 'clear_category', label: 'Clear category' },
    { value: 'add_tag', label: 'Add tag' },
    { value: 'remove_tag', label: 'Remove tag' },
    { value: 'remove_all_tags', label: 'Remove all tags' },
    { value: 'set_description', label: 'Set description to' },
    { value: 'append_to_description', label: 'Append to description' },
    { value: 'prepend_to_description', label: 'Prepend to description' },
    { value: 'set_amount', label: 'Set amount to' },
    { value: 'set_source_account', label: 'Set source account to' }
];

function needsValue(type: RuleTriggerType): boolean {
    return type !== 'has_no_category' && type !== 'has_any_category';
}

function needsActionValue(type: RuleActionType): boolean {
    return type !== 'clear_category' && type !== 'remove_all_tags';
}

const canSave = computed(() => {
    return rule.value.name.trim() !== '' &&
        rule.value.ruleGroupId !== '' &&
        rule.value.triggers.length > 0 &&
        rule.value.actions.length > 0;
});

function init(): void {
    loading.value = true;

    const chain: Promise<void> = rulesStore.loadRuleGroups().then(() => {
        if (props.editRuleId) {
            return rulesStore.loadRules().then(() => {
                const existing = rulesStore.allRules.find(r => r.id === props.editRuleId);
                if (existing) {
                    rule.value = existing;
                }
            });
        }

        rule.value = new Rule();
        rule.value.ruleGroupId = props.presetGroupId || (rulesStore.allRuleGroups[0]?.id ?? '');
        return Promise.resolve();
    });

    chain.catch(error => {
        if (!error.processed) {
            showToast(error.message || error);
        }
    }).finally(() => {
        loading.value = false;
    });
}

function save(): void {
    if (!canSave.value) {
        return;
    }

    // Validate: triggers/actions that need a value must have one
    for (const t of rule.value.triggers) {
        if (needsValue(t.triggerType) && t.triggerValue.trim() === '') {
            showToast(tt('Please fill in all trigger values'));
            return;
        }
    }
    for (const a of rule.value.actions) {
        if (needsActionValue(a.actionType) && a.actionValue.trim() === '') {
            showToast(tt('Please fill in all action values'));
            return;
        }
    }

    const clientSessionId = `${Date.now()}`;

    rulesStore.saveRule({
        rule: rule.value,
        isEdit: isEdit.value,
        clientSessionId
    }).then(() => {
        showToast(isEdit.value ? 'Rule saved' : 'Rule created');
        rulesStore.loadRules();
        props.f7router.back();
    }).catch(error => {
        if (!error.processed) {
            showToast(error.message || error);
        }
    });
}

function confirmDelete(): void {
    showConfirm(tt('Delete this rule?'), () => {
        rulesStore.deleteRule({ ruleId: rule.value.id }).then(() => {
            showToast('Rule deleted');
            rulesStore.loadRules();
            props.f7router.back();
        }).catch(error => {
            if (!error.processed) {
                showToast(error.message || error);
            }
        });
    });
}

onMounted(init);
</script>

<style scoped>
.row-input {
    display: flex;
    flex-direction: column;
    gap: 8px;
    width: 100%;
    padding-right: 12px;
}
.trigger-type-select,
.action-type-select {
    width: 100%;
}
</style>
