const allTemplateTypes = {
    Normal: 1,
    Schedule: 2,
};

const allTemplateScheduledFrequencyTypes = {
    Disabled: {
        type: 0,
        name: 'Disabled'
    },
    Weekly: {
        type: 1,
        name: 'Weekly'
    },
    Monthly: {
        type: 2,
        name: 'Monthly'
    }
};

export default {
    allTemplateTypes: allTemplateTypes,
    allTemplateScheduledFrequencyTypes: allTemplateScheduledFrequencyTypes,
}
