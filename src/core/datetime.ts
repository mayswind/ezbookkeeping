export class MeridiemIndicator {
    private static readonly allInstances: MeridiemIndicator[] = [];

    public static readonly AM = new MeridiemIndicator(0, 'AM');
    public static readonly PM = new MeridiemIndicator(1, 'PM');

    public readonly type: number;
    public readonly value: string;

    private constructor(type: number, value: string) {
        this.type = type;
        this.value = value;

        MeridiemIndicator.allInstances.push(this);
    }

    public static values(): MeridiemIndicator[] {
        return MeridiemIndicator.allInstances;
    }
}
