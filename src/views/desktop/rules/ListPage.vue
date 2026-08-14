<template>
    <v-row class="match-height">
        <v-col cols="12">
            <v-card>
                <v-card-item>
                    <div class="d-flex align-center justify-space-between">
                        <div class="text-h6">{{ tt('Rules') }}</div>
                        <div>
                            <v-btn variant="text" :prepend-icon="mdiFolderPlus" @click="openGroupDialog()">{{ tt('New Group') }}</v-btn>
                            <v-btn color="primary" variant="tonal" :prepend-icon="mdiPlus" :disabled="groups.length === 0" @click="openRuleDialog()">{{ tt('New Rule') }}</v-btn>
                        </div>
                    </div>
                </v-card-item>

                <v-divider />

                <v-card-text>
                    <div v-if="loading" class="text-center pa-4">
                        <v-progress-circular indeterminate />
                    </div>

                    <div v-else-if="rules.length === 0 && groups.length === 0" class="text-center pa-8 text-medium-emphasis">
                        <v-icon :icon="mdiAutoFix" size="48" />
                        <p class="mt-2">{{ tt('No rules yet. Create a group, then a rule to auto-categorize transactions.') }}</p>
                    </div>

                    <div v-else>
                        <div v-for="group in groups" :key="group.id" class="mb-4">
                            <div class="d-flex align-center justify-space-between mb-1">
                                <div class="text-subtitle-1 font-weight-bold">
                                    {{ group.name }}
                                    <v-chip v-if="!group.active" size="x-small" color="grey">{{ tt('Inactive') }}</v-chip>
                                </div>
                                <div>
                                    <v-btn variant="text" size="small" :icon="mdiPencil" @click="openGroupDialog(group)"></v-btn>
                                    <v-btn variant="text" size="small" color="error" :icon="mdiDelete" @click="confirmDeleteGroup(group)"></v-btn>
                                </div>
                            </div>

                            <v-data-table
                                :headers="ruleHeaders"
                                :items="rulesOfGroup(group.id)"
                                :no-data-text="tt('No rules in this group.')"
                                hide-default-footer
                                :items-per-page="-1"
                                density="comfortable">
                                <template v-slot:item.name="{ item }">
                                    <v-icon v-if="item.active" :icon="mdiCheckCircleOutline" color="success" size="small" class="me-1" />
                                    <v-icon v-else :icon="mdiMinusCircleOutline" color="grey" size="small" class="me-1" />
                                    <span>{{ item.name }}</span>
                                </template>
                                <template v-slot:item.summary="{ item }">
                                    <span class="text-medium-emphasis">{{ summarizeTriggers(item) }} → {{ summarizeActions(item) }}</span>
                                </template>
                                <template v-slot:item.actions="{ item }">
                                    <v-btn variant="text" size="small" :icon="mdiPencil" @click="openRuleDialog(item)"></v-btn>
                                    <v-btn variant="text" size="small" color="error" :icon="mdiDelete" @click="confirmDeleteRule(item)"></v-btn>
                                </template>
                            </v-data-table>
                        </div>
                    </div>
                </v-card-text>
            </v-card>
        </v-col>

        <snack-bar ref="snackbar" />
    </v-row>

    <!-- Group dialog -->
    <v-dialog v-model="groupDialog.show" max-width="480" persistent>
        <v-card>
            <v-card-title>{{ groupDialog.isEdit ? tt('Edit Group') : tt('New Group') }}</v-card-title>
            <v-card-text>
                <v-text-field v-model="groupDialog.name" :label="tt('Group Name')" density="comfortable" class="mb-2" />
                <v-text-field v-model="groupDialog.comment" :label="tt('Comment')" density="comfortable" class="mb-2" />
                <v-switch v-model="groupDialog.active" :label="tt('Active')" density="comfortable" hide-details />
                <v-switch v-model="groupDialog.stopProcessing" :label="tt('Stop processing after this group')" density="comfortable" hide-details />
            </v-card-text>
            <v-card-actions>
                <v-spacer />
                <v-btn variant="text" @click="groupDialog.show = false">{{ tt('Cancel') }}</v-btn>
                <v-btn color="primary" variant="tonal" :disabled="!groupDialog.name.trim()" @click="saveGroupDialog">{{ tt('Save') }}</v-btn>
            </v-card-actions>
        </v-card>
    </v-dialog>

    <!-- Rule dialog -->
    <v-dialog v-model="ruleDialog.show" max-width="720" persistent scrollable>
        <v-card>
            <v-card-title>{{ ruleDialog.isEdit ? tt('Edit Rule') : tt('New Rule') }}</v-card-title>
            <v-card-text>
                <v-text-field v-model="ruleDialog.name" :label="tt('Rule Name')" density="comfortable" class="mb-2" />
                <v-select v-model="ruleDialog.ruleGroupId" :items="groups" item-title="name" item-value="id" :label="tt('Rule Group')" density="comfortable" class="mb-2" />
                <div class="d-flex ga-2 mb-2">
                    <v-switch v-model="ruleDialog.active" :label="tt('Active')" density="comfortable" hide-details />
                    <v-switch v-model="ruleDialog.applyOnCreate" :label="tt('On Create')" density="comfortable" hide-details />
                </div>
                <v-segmented-btn v-model="ruleDialog.strictModel" class="mb-3">
                    <v-btn :value="false">{{ tt('Match ANY') }}</v-btn>
                    <v-btn :value="true">{{ tt('Match ALL') }}</v-btn>
                </v-segmented-btn>

                <div class="text-subtitle-2 mb-1">{{ tt('When transaction matches') }}</div>
                <div v-for="(trigger, idx) in ruleDialog.triggers" :key="'t' + idx" class="d-flex align-center ga-2 mb-2">
                    <v-select v-model="trigger.triggerType" :items="triggerTypeItems" item-title="label" item-value="value" density="compact" style="max-width: 260px" hide-details />
                    <v-text-field v-if="needsTriggerValue(trigger.triggerType)" v-model="trigger.triggerValue" density="compact" hide-details />
                    <v-btn variant="text" size="small" color="error" :icon="mdiDelete" @click="ruleDialog.triggers.splice(idx, 1)" :disabled="ruleDialog.triggers.length <= 1" />
                </div>
                <v-btn variant="text" color="primary" size="small" :prepend-icon="mdiPlus" @click="ruleDialog.triggers.push({ triggerType: 'description_contains', triggerValue: '', prohibited: false, stopProcessing: false })">
                    {{ tt('Add Trigger') }}
                </v-btn>

                <div class="text-subtitle-2 mt-3 mb-1">{{ tt('Then apply') }}</div>
                <div v-for="(action, idx) in ruleDialog.actions" :key="'a' + idx" class="d-flex align-center ga-2 mb-2">
                    <v-select v-model="action.actionType" :items="actionTypeItems" item-title="label" item-value="value" density="compact" style="max-width: 260px" hide-details />
                    <v-text-field v-if="needsActionValue(action.actionType)" v-model="action.actionValue" density="compact" hide-details />
                    <v-btn variant="text" size="small" color="error" :icon="mdiDelete" @click="ruleDialog.actions.splice(idx, 1)" :disabled="ruleDialog.actions.length <= 1" />
                </div>
                <v-btn variant="text" color="primary" size="small" :prepend-icon="mdiPlus" @click="ruleDialog.actions.push({ actionType: 'set_category', actionValue: '', stopProcessing: false })">
                    {{ tt('Add Action') }}
                </v-btn>
            </v-card-text>
            <v-card-actions>
                <v-spacer />
                <v-btn variant="text" @click="ruleDialog.show = false">{{ tt('Cancel') }}</v-btn>
                <v-btn color="primary" variant="tonal" :disabled="!canSaveRule" @click="saveRuleDialog">{{ tt('Save') }}</v-btn>
            </v-card-actions>
        </v-card>
    </v-dialog>
</template>

<script setup lang="ts">
import { ref, computed, reactive, onMounted } from 'vue';
import { mdiPlus, mdiFolderPlus, mdiPencil, mdiDelete, mdiCheckCircleOutline, mdiMinusCircleOutline, mdiAutoFix } from '@mdi/js';

import { useI18n } from '@/locales/helpers.ts';
import { useRulesStore } from '@/stores/rule.ts';

import SnackBar from '@/components/desktop/SnackBar.vue';

import { Rule, RuleGroup, RuleTrigger, RuleAction, type RuleTriggerType, type RuleActionType } from '@/models/rule.ts';

type SnackBarType = InstanceType<typeof SnackBar>;

const { tt } = useI18n();
const rulesStore = useRulesStore();

const snackbar = ref<SnackBarType | null>(null);

const loading = ref<boolean>(true);

const groups = computed<RuleGroup[]>(() => rulesStore.allRuleGroups);
const rules = computed<Rule[]>(() => rulesStore.allRules);

const ruleHeaders = computed(() => [
    { title: tt('Name'), key: 'name', sortable: true },
    { title: tt('Summary'), key: 'summary', sortable: false },
    { title: tt('Actions'), key: 'actions', sortable: false, align: 'end' as const }
]);

const triggerTypeItems: { value: RuleTriggerType; label: string }[] = [
    { value: 'description_contains', label: tt('Description contains') },
    { value: 'description_is', label: tt('Description is') },
    { value: 'amount_is', label: tt('Amount is') },
    { value: 'amount_less', label: tt('Amount is less than') },
    { value: 'amount_more', label: tt('Amount is more than') },
    { value: 'source_account_is', label: tt('Source account is') },
    { value: 'destination_account_is', label: tt('Destination account is') },
    { value: 'category_is', label: tt('Category is') },
    { value: 'has_no_category', label: tt('Has no category') },
    { value: 'has_any_category', label: tt('Has any category') }
];

const actionTypeItems: { value: RuleActionType; label: string }[] = [
    { value: 'set_category', label: tt('Set category to') },
    { value: 'clear_category', label: tt('Clear category') },
    { value: 'add_tag', label: tt('Add tag') },
    { value: 'remove_tag', label: tt('Remove tag') },
    { value: 'remove_all_tags', label: tt('Remove all tags') },
    { value: 'set_description', label: tt('Set description to') },
    { value: 'append_to_description', label: tt('Append to description') },
    { value: 'prepend_to_description', label: tt('Prepend to description') },
    { value: 'set_amount', label: tt('Set amount to') },
    { value: 'set_source_account', label: tt('Set source account to') }
];

function needsTriggerValue(type: RuleTriggerType): boolean {
    return type !== 'has_no_category' && type !== 'has_any_category';
}

function needsActionValue(type: RuleActionType): boolean {
    return type !== 'clear_category' && type !== 'remove_all_tags';
}

function rulesOfGroup(groupId: string): Rule[] {
    return rules.value.filter(r => r.ruleGroupId === groupId);
}

function summarizeTriggers(rule: Rule): string {
    const mode = rule.strict ? tt('ALL') : tt('ANY');
    return `${mode}: ${rule.triggers.map(t => t.triggerType.replace(/_/g, ' ')).join(', ')}`;
}

function summarizeActions(rule: Rule): string {
    return rule.actions.map(a => a.actionType.replace(/_/g, ' ')).join(', ');
}

function loadData(): void {
    loading.value = true;
    Promise.all([rulesStore.loadRuleGroups(), rulesStore.loadRules()]).catch(error => {
        if (!error.processed) {
            snackbar.value?.showError(error);
        }
    }).finally(() => {
        loading.value = false;
    });
}

// --- Group dialog ---
const groupDialog = reactive({
    show: false,
    isEdit: false,
    id: '',
    name: '',
    comment: '',
    active: true,
    stopProcessing: false
});

function openGroupDialog(group?: RuleGroup): void {
    if (group) {
        groupDialog.isEdit = true;
        groupDialog.id = group.id;
        groupDialog.name = group.name;
        groupDialog.comment = group.comment;
        groupDialog.active = group.active;
        groupDialog.stopProcessing = group.stopProcessing;
    } else {
        groupDialog.isEdit = false;
        groupDialog.id = '';
        groupDialog.name = '';
        groupDialog.comment = '';
        groupDialog.active = true;
        groupDialog.stopProcessing = false;
    }
    groupDialog.show = true;
}

function saveGroupDialog(): void {
    const group = new RuleGroup();
    group.id = groupDialog.id;
    group.name = groupDialog.name;
    group.comment = groupDialog.comment;
    group.active = groupDialog.active;
    group.stopProcessing = groupDialog.stopProcessing;

    rulesStore.saveRuleGroup({ group, isEdit: groupDialog.isEdit, clientSessionId: `${Date.now()}` }).then(() => {
        groupDialog.show = false;
        snackbar.value?.showMessage(groupDialog.isEdit ? 'Group saved' : 'Group created');
        loadData();
    }).catch(error => {
        if (!error.processed) {
            snackbar.value?.showError(error);
        }
    });
}

function confirmDeleteGroup(group: RuleGroup): void {
    if (!confirm(tt('Delete this group and all its rules?'))) {
        return;
    }
    rulesStore.deleteRuleGroup({ groupId: group.id }).then(() => {
        snackbar.value?.showMessage('Group deleted');
        loadData();
    }).catch(error => {
        if (!error.processed) {
            snackbar.value?.showError(error);
        }
    });
}

// --- Rule dialog ---
interface DialogTrigger { triggerType: RuleTriggerType; triggerValue: string; prohibited: boolean; stopProcessing: boolean }
interface DialogAction { actionType: RuleActionType; actionValue: string; stopProcessing: boolean }

const ruleDialog = reactive({
    show: false,
    isEdit: false,
    id: '',
    name: '',
    ruleGroupId: '',
    active: true,
    applyOnCreate: true,
    strictModel: false as boolean,
    triggers: [] as DialogTrigger[],
    actions: [] as DialogAction[]
});

const canSaveRule = computed(() => {
    return ruleDialog.name.trim() !== '' && ruleDialog.ruleGroupId !== '' && ruleDialog.triggers.length > 0 && ruleDialog.actions.length > 0;
});

function openRuleDialog(rule?: Rule): void {
    if (rule) {
        ruleDialog.isEdit = true;
        ruleDialog.id = rule.id;
        ruleDialog.name = rule.name;
        ruleDialog.ruleGroupId = rule.ruleGroupId;
        ruleDialog.active = rule.active;
        ruleDialog.applyOnCreate = rule.applyOnCreate;
        ruleDialog.strictModel = rule.strict;
        ruleDialog.triggers = rule.triggers.map(t => ({ triggerType: t.triggerType, triggerValue: t.triggerValue, prohibited: t.prohibited, stopProcessing: t.stopProcessing }));
        ruleDialog.actions = rule.actions.map(a => ({ actionType: a.actionType, actionValue: a.actionValue, stopProcessing: a.stopProcessing }));
    } else {
        ruleDialog.isEdit = false;
        ruleDialog.id = '';
        ruleDialog.name = '';
        ruleDialog.ruleGroupId = groups.value[0]?.id ?? '';
        ruleDialog.active = true;
        ruleDialog.applyOnCreate = true;
        ruleDialog.strictModel = false;
        ruleDialog.triggers = [{ triggerType: 'description_contains', triggerValue: '', prohibited: false, stopProcessing: false }];
        ruleDialog.actions = [{ actionType: 'set_category', actionValue: '', stopProcessing: false }];
    }
    ruleDialog.show = true;
}

function saveRuleDialog(): void {
    if (!canSaveRule.value) {
        return;
    }

    const rule = new Rule();
    rule.id = ruleDialog.id;
    rule.name = ruleDialog.name;
    rule.ruleGroupId = ruleDialog.ruleGroupId;
    rule.active = ruleDialog.active;
    rule.applyOnCreate = ruleDialog.applyOnCreate;
    rule.applyOnUpdate = false;
    rule.stopProcessing = false;
    rule.comment = '';
    rule.strict = ruleDialog.strictModel;
    rule.triggers = ruleDialog.triggers.map(t => {
        const trigger = new RuleTrigger();
        trigger.triggerType = t.triggerType;
        trigger.triggerValue = t.triggerValue;
        trigger.prohibited = t.prohibited;
        trigger.stopProcessing = t.stopProcessing;
        return trigger;
    });
    rule.actions = ruleDialog.actions.map(a => {
        const action = new RuleAction();
        action.actionType = a.actionType;
        action.actionValue = a.actionValue;
        action.stopProcessing = a.stopProcessing;
        return action;
    });

    rulesStore.saveRule({ rule, isEdit: ruleDialog.isEdit, clientSessionId: `${Date.now()}` }).then(() => {
        ruleDialog.show = false;
        snackbar.value?.showMessage(ruleDialog.isEdit ? 'Rule saved' : 'Rule created');
        loadData();
    }).catch(error => {
        if (!error.processed) {
            snackbar.value?.showError(error);
        }
    });
}

function confirmDeleteRule(rule: Rule): void {
    if (!confirm(tt('Delete this rule?'))) {
        return;
    }
    rulesStore.deleteRule({ ruleId: rule.id }).then(() => {
        snackbar.value?.showMessage('Rule deleted');
        loadData();
    }).catch(error => {
        if (!error.processed) {
            snackbar.value?.showError(error);
        }
    });
}

onMounted(loadData);
</script>
