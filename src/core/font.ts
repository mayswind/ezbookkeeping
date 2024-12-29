export class FontSize {
    private static readonly allInstances: FontSize[] = [];

    public static readonly Small = new FontSize(0, 'font-size-small');
    public static readonly Default = new FontSize(1, 'font-size-default');
    public static readonly Large = new FontSize(2, 'font-size-large');
    public static readonly XLarge = new FontSize(3, 'font-size-x-large');
    public static readonly XXLarge = new FontSize(4, 'font-size-xx-large');
    public static readonly XXXLarge = new FontSize(5, 'font-size-xxx-large');
    public static readonly XXXXLarge = new FontSize(6, 'font-size-xxxx-large');

    public static readonly MinimumFontSize = FontSize.Small;
    public static readonly MaximumFontSize = FontSize.XXXXLarge;

    public readonly type: number;
    public readonly className: string;

    private constructor(type: number, className: string) {
        this.type = type;
        this.className = className;

        FontSize.allInstances.push(this);
    }

    public static values(): FontSize[] {
        return FontSize.allInstances;
    }
}

export const FONT_SIZE_PREVIEW_CLASSNAME_PREFIX: string = 'preview-';
