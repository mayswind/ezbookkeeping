import services from '@/lib/services.js';

export function getMobileUrlQrCodePath() {
    return services.generateQrCodeUrl('mobile_url');
}
