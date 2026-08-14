// [PLUGIN:rules] Rules engine feature - TypeScript models.
// Mirrors the backend DTOs in pkg/models/rule*.go.

export type RuleTriggerType =
    | 'description_is'
    | 'description_contains'
    | 'amount_is'
    | 'amount_less'
    | 'amount_more'
    | 'source_account_is'
    | 'destination_account_is'
    | 'category_is'
    | 'has_no_category'
    | 'has_any_category';

export type RuleActionType =
    | 'set_category'
    | 'clear_category'
    | 'add_tag'
    | 'remove_tag'
    | 'remove_all_tags'
    | 'set_description'
    | 'append_to_description'
    | 'prepend_to_description'
    | 'set_amount'
    | 'set_source_account';

export interface RuleTriggerInfoResponse {
    readonly id: string;
    readonly triggerType: RuleTriggerType;
    readonly triggerValue: string;
    readonly displayOrder: number;
    readonly prohibited: boolean;
    readonly stopProcessing: boolean;
}

export interface RuleActionInfoResponse {
    readonly id: string;
    readonly actionType: RuleActionType;
    readonly actionValue: string;
    readonly displayOrder: number;
    readonly stopProcessing: boolean;
}

export interface RuleInfoResponse {
    readonly id: string;
    readonly ruleGroupId: string;
    readonly name: string;
    readonly comment: string;
    readonly displayOrder: number;
    readonly active: boolean;
    readonly strict: boolean;
    readonly stopProcessing: boolean;
    readonly applyOnCreate: boolean;
    readonly applyOnUpdate: boolean;
    readonly triggers: RuleTriggerInfoResponse[];
    readonly actions: RuleActionInfoResponse[];
}

export interface RuleGroupInfoResponse {
    readonly id: string;
    readonly name: string;
    readonly comment: string;
    readonly displayOrder: number;
    readonly active: boolean;
    readonly stopProcessing: boolean;
}

// --- Editable model classes (used by the editor UI) ---

export class RuleTrigger implements RuleTriggerInfoResponse {
    public id: string;
    public triggerType: RuleTriggerType;
    public triggerValue: string;
    public displayOrder: number;
    public prohibited: boolean;
    public stopProcessing: boolean;

    public constructor(triggerType: RuleTriggerType = 'description_contains', triggerValue: string = '') {
        this.id = '';
        this.triggerType = triggerType;
        this.triggerValue = triggerValue;
        this.displayOrder = 0;
        this.prohibited = false;
        this.stopProcessing = false;
    }

    public static of(resp: RuleTriggerInfoResponse): RuleTrigger {
        const t = new RuleTrigger(resp.triggerType, resp.triggerValue);
        t.id = resp.id;
        t.displayOrder = resp.displayOrder;
        t.prohibited = resp.prohibited;
        t.stopProcessing = resp.stopProcessing;
        return t;
    }
}

export class RuleAction implements RuleActionInfoResponse {
    public id: string;
    public actionType: RuleActionType;
    public actionValue: string;
    public displayOrder: number;
    public stopProcessing: boolean;

    public constructor(actionType: RuleActionType = 'set_category', actionValue: string = '') {
        this.id = '';
        this.actionType = actionType;
        this.actionValue = actionValue;
        this.displayOrder = 0;
        this.stopProcessing = false;
    }

    public static of(resp: RuleActionInfoResponse): RuleAction {
        const a = new RuleAction(resp.actionType, resp.actionValue);
        a.id = resp.id;
        a.displayOrder = resp.displayOrder;
        a.stopProcessing = resp.stopProcessing;
        return a;
    }
}

export class Rule implements RuleInfoResponse {
    public id: string;
    public ruleGroupId: string;
    public name: string;
    public comment: string;
    public displayOrder: number;
    public active: boolean;
    public strict: boolean;
    public stopProcessing: boolean;
    public applyOnCreate: boolean;
    public applyOnUpdate: boolean;
    public triggers: RuleTrigger[];
    public actions: RuleAction[];

    public constructor() {
        this.id = '';
        this.ruleGroupId = '';
        this.name = '';
        this.comment = '';
        this.displayOrder = 0;
        this.active = true;
        this.strict = false; // ANY by default (more intuitive for users)
        this.stopProcessing = false;
        this.applyOnCreate = true;
        this.applyOnUpdate = false;
        this.triggers = [new RuleTrigger()];
        this.actions = [new RuleAction()];
    }

    public static of(resp: RuleInfoResponse): Rule {
        const r = new Rule();
        r.id = resp.id;
        r.ruleGroupId = resp.ruleGroupId;
        r.name = resp.name;
        r.comment = resp.comment;
        r.displayOrder = resp.displayOrder;
        r.active = resp.active;
        r.strict = resp.strict;
        r.stopProcessing = resp.stopProcessing;
        r.applyOnCreate = resp.applyOnCreate;
        r.applyOnUpdate = resp.applyOnUpdate;
        r.triggers = (resp.triggers || []).map(RuleTrigger.of);
        r.actions = (resp.actions || []).map(RuleAction.of);
        return r;
    }

    public toCreateRequest(clientSessionId: string): RuleCreateRequest {
        return {
            ruleGroupId: this.ruleGroupId,
            name: this.name,
            comment: this.comment,
            active: this.active,
            strict: this.strict,
            stopProcessing: this.stopProcessing,
            applyOnCreate: this.applyOnCreate,
            applyOnUpdate: this.applyOnUpdate,
            triggers: this.triggers.map(t => ({
                triggerType: t.triggerType,
                triggerValue: t.triggerValue,
                prohibited: t.prohibited,
                stopProcessing: t.stopProcessing
            })),
            actions: this.actions.map(a => ({
                actionType: a.actionType,
                actionValue: a.actionValue,
                stopProcessing: a.stopProcessing
            })),
            clientSessionId
        };
    }

    public toModifyRequest(): RuleModifyRequest {
        return {
            id: this.id,
            ruleGroupId: this.ruleGroupId,
            name: this.name,
            comment: this.comment,
            active: this.active,
            strict: this.strict,
            stopProcessing: this.stopProcessing,
            applyOnCreate: this.applyOnCreate,
            applyOnUpdate: this.applyOnUpdate,
            triggers: this.triggers.map(t => ({
                triggerType: t.triggerType,
                triggerValue: t.triggerValue,
                prohibited: t.prohibited,
                stopProcessing: t.stopProcessing
            })),
            actions: this.actions.map(a => ({
                actionType: a.actionType,
                actionValue: a.actionValue,
                stopProcessing: a.stopProcessing
            }))
        };
    }
}

export class RuleGroup implements RuleGroupInfoResponse {
    public id: string;
    public name: string;
    public comment: string;
    public displayOrder: number;
    public active: boolean;
    public stopProcessing: boolean;

    public constructor() {
        this.id = '';
        this.name = '';
        this.comment = '';
        this.displayOrder = 0;
        this.active = true;
        this.stopProcessing = false;
    }

    public static of(resp: RuleGroupInfoResponse): RuleGroup {
        const g = new RuleGroup();
        g.id = resp.id;
        g.name = resp.name;
        g.comment = resp.comment;
        g.displayOrder = resp.displayOrder;
        g.active = resp.active;
        g.stopProcessing = resp.stopProcessing;
        return g;
    }

    public toCreateRequest(clientSessionId: string): RuleGroupCreateRequest {
        return {
            name: this.name,
            comment: this.comment,
            active: this.active,
            stopProcessing: this.stopProcessing,
            clientSessionId
        };
    }

    public toModifyRequest(): RuleGroupModifyRequest {
        return {
            id: this.id,
            name: this.name,
            comment: this.comment,
            active: this.active,
            stopProcessing: this.stopProcessing
        };
    }
}

// --- Request interfaces ---

export interface RuleTriggerInput {
    readonly triggerType: RuleTriggerType;
    readonly triggerValue: string;
    readonly prohibited: boolean;
    readonly stopProcessing: boolean;
}

export interface RuleActionInput {
    readonly actionType: RuleActionType;
    readonly actionValue: string;
    readonly stopProcessing: boolean;
}

export interface RuleCreateRequest {
    readonly ruleGroupId: string;
    readonly name: string;
    readonly comment: string;
    readonly active: boolean;
    readonly strict: boolean;
    readonly stopProcessing: boolean;
    readonly applyOnCreate: boolean;
    readonly applyOnUpdate: boolean;
    readonly triggers: RuleTriggerInput[];
    readonly actions: RuleActionInput[];
    readonly clientSessionId: string;
}

export interface RuleModifyRequest {
    readonly id: string;
    readonly ruleGroupId: string;
    readonly name: string;
    readonly comment: string;
    readonly active: boolean;
    readonly strict: boolean;
    readonly stopProcessing: boolean;
    readonly applyOnCreate: boolean;
    readonly applyOnUpdate: boolean;
    readonly triggers: RuleTriggerInput[];
    readonly actions: RuleActionInput[];
}

export interface RuleGroupCreateRequest {
    readonly name: string;
    readonly comment: string;
    readonly active: boolean;
    readonly stopProcessing: boolean;
    readonly clientSessionId: string;
}

export interface RuleGroupModifyRequest {
    readonly id: string;
    readonly name: string;
    readonly comment: string;
    readonly active: boolean;
    readonly stopProcessing: boolean;
}
