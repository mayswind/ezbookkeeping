import { describe, expect, it, vi } from 'vitest';

import { TemplateType } from '@/core/template.ts';
import { TransactionTemplate, type TransactionTemplateInfoResponse } from '@/models/transaction_template.ts';

vi.mock('@/lib/datetime.ts', () => ({
    getBrowserTimezoneName: () => 'Europe/Berlin'
}));

function createTemplate(timeZone?: string): TransactionTemplate {
    const response: TransactionTemplateInfoResponse = {
        id: '1', timeSequenceId: '0', templateType: TemplateType.Schedule.type,
        name: 'Monthly bill', type: 3, categoryId: '11', time: 0,
        timeZone: timeZone, utcOffset: 120, sourceAccountId: '1', destinationAccountId: '0',
        sourceAmount: 1250, destinationAmount: 0, hideAmount: false, tagIds: [], comment: '',
        editable: true, displayOrder: 0, hidden: false,
        scheduledFrequencyType: 2, scheduledFrequency: '1', scheduledAt: 1320
    };

    return TransactionTemplate.ofTemplate(response);
}

describe('Scheduled template time zone', () => {
    it('should retain the named zone when loading and saving a template', () => {
        const template = createTemplate('Europe/Berlin');

        expect(template.timeZone).toBe('Europe/Berlin');
        expect(template.toTemplateCreateRequest('session').timeZone).toBe('Europe/Berlin');
        expect(template.toTemplateModifyRequest().timeZone).toBe('Europe/Berlin');
    });

    it('should retain the named zone when copying a template', () => {
        const template = createTemplate();
        template.fillFrom(createTemplate('Europe/Berlin'));

        expect(template.timeZone).toBe('Europe/Berlin');
        expect(template.toTemplateModifyRequest().timeZone).toBe('Europe/Berlin');
    });

    it('should send the browser zone for a schedule that follows the browser', () => {
        const template = createTemplate('');

        expect(template.toTemplateCreateRequest('session').timeZone).toBe('Europe/Berlin');
        expect(template.toTemplateModifyRequest().timeZone).toBe('Europe/Berlin');
    });

    it('should preserve the fixed offset of a legacy template', () => {
        const template = createTemplate();

        expect(template.timeZone).toBeUndefined();
        expect(template.toTemplateModifyRequest().timeZone).toBeUndefined();
        expect(template.toTemplateModifyRequest().utcOffset).toBe(120);
    });

    it('should omit the named zone from a normal template', () => {
        const template = createTemplate('Europe/Berlin');
        template.templateType = TemplateType.Normal.type;

        expect(template.toTemplateCreateRequest('session').timeZone).toBeUndefined();
        expect(template.toTemplateModifyRequest().timeZone).toBeUndefined();
    });
});
