import services from './services.js';

export function getMobileUrlQrCodePath() {
    return services.generateQrCodeUrl('mobile_url');
}
