import services from './services.ts';

export function getMobileUrlQrCodePath(): string {
    return services.generateQrCodeUrl('mobile_url');
}
