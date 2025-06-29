import type { ApplicationCloudSetting } from '@/core/setting.ts';

export interface UserApplicationCloudSettingsUpdateRequest {
    readonly settings: ApplicationCloudSetting[];
    readonly fullUpdate: boolean;
}
