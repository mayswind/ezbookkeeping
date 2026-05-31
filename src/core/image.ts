import type { TypeAndName } from './base.ts';

export class ImageUploadQualityType implements TypeAndName {
    private static readonly allInstances: ImageUploadQualityType[] = [];
    private static readonly allInstancesByType: Record<number, ImageUploadQualityType> = {};

    public static readonly Original = new ImageUploadQualityType(0, 'Original', null, 1, null);
    public static readonly UHD4K = new ImageUploadQualityType(1, '4K', 3840, 0.93, 1500);
    public static readonly QHD2K = new ImageUploadQualityType(2, '2K', 2560, 0.9, 700);
    public static readonly FHD1080P = new ImageUploadQualityType(3, '1080p', 1920, 0.85, 500);
    public static readonly HD720P = new ImageUploadQualityType(4, '720p', 1280, 0.8, 300);

    public static readonly Default = ImageUploadQualityType.Original;

    public readonly type: number;
    public readonly name: string;
    public readonly maxLongSidePixels: number | null;
    public readonly quality: number;
    public readonly estimatedKiB: number | null;

    private constructor(type: number, name: string, maxLongSidePixels: number | null, quality: number, estimatedKiB: number | null) {
        this.type = type;
        this.name = name;
        this.maxLongSidePixels = maxLongSidePixels;
        this.quality = quality;
        this.estimatedKiB = estimatedKiB;

        ImageUploadQualityType.allInstances.push(this);
        ImageUploadQualityType.allInstancesByType[type] = this;
    }

    public static values(): ImageUploadQualityType[] {
        return ImageUploadQualityType.allInstances;
    }

    public static valueOf(type: number): ImageUploadQualityType | undefined {
        return ImageUploadQualityType.allInstancesByType[type];
    }
}
