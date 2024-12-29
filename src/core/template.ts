import type { TypeAndName } from './base.ts';

type TemplateTypeName = 'Normal' | 'Schedule';

export class TemplateType implements TypeAndName {
    private static readonly allInstances: TemplateType[] = [];
    private static readonly allInstancesByType: Record<number, TemplateType> = {};
    private static readonly allInstancesByTypeName: Record<string, TemplateType> = {};

    public static readonly Normal = new TemplateType(1, 'Normal');
    public static readonly Schedule = new TemplateType(2, 'Schedule');

    public readonly type: number;
    public readonly name: TemplateTypeName;

    private constructor(type: number, name: TemplateTypeName) {
        this.type = type;
        this.name = name;

        TemplateType.allInstances.push(this);
        TemplateType.allInstancesByType[type] = this;
        TemplateType.allInstancesByTypeName[name] = this;
    }

    public static values(): TemplateType[] {
        return TemplateType.allInstances;
    }

    public static all(): Record<TemplateTypeName, TemplateType> {
        return TemplateType.allInstancesByTypeName;
    }

    public static valueOf(type: number): TemplateType {
        return TemplateType.allInstancesByType[type];
    }
}

type ScheduledTemplateFrequencyTypeName = 'Disabled' | 'Weekly' | 'Monthly';

export class ScheduledTemplateFrequencyType implements TypeAndName {
    private static readonly allInstances: ScheduledTemplateFrequencyType[] = [];
    private static readonly allInstancesByType: Record<number, ScheduledTemplateFrequencyType> = {};
    private static readonly allInstancesByTypeName: Record<string, ScheduledTemplateFrequencyType> = {};

    public static readonly Disabled = new ScheduledTemplateFrequencyType(0, 'Disabled');
    public static readonly Weekly = new ScheduledTemplateFrequencyType(1, 'Weekly');
    public static readonly Monthly = new ScheduledTemplateFrequencyType(2, 'Monthly');

    public readonly type: number;
    public readonly name: ScheduledTemplateFrequencyTypeName;

    private constructor(type: number, name: ScheduledTemplateFrequencyTypeName) {
        this.type = type;
        this.name = name;

        ScheduledTemplateFrequencyType.allInstances.push(this);
        ScheduledTemplateFrequencyType.allInstancesByType[type] = this;
        ScheduledTemplateFrequencyType.allInstancesByTypeName[name] = this;
    }

    public static values(): ScheduledTemplateFrequencyType[] {
        return ScheduledTemplateFrequencyType.allInstances;
    }

    public static all(): Record<ScheduledTemplateFrequencyTypeName, ScheduledTemplateFrequencyType> {
        return ScheduledTemplateFrequencyType.allInstancesByTypeName;
    }

    public static valueOf(type: number): ScheduledTemplateFrequencyType {
        return ScheduledTemplateFrequencyType.allInstancesByType[type];
    }
}
