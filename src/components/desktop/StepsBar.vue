<template>
    <div class="d-flex" :style="`min-width: ${minWidth}px`" v-if="minWidth"></div>
    <v-slide-group class="slide-group-with-stepper mb-10 hidden-xs" show-arrows>
        <v-slide-group-item :key="idx" v-for="(step, idx) in steps">
            <div class="mx-1"
                 :class="{ 'slide-group-step-active': isStepActive(step), 'slide-group-step-completed': isStepCompleted(idx), 'cursor-pointer': isClickable }"
                 @click="changeStep(step)">
                <div class="d-flex align-center gap-x-2">
                    <div class="d-flex align-center gap-2">
                        <div class="d-flex align-center justify-center" style="block-size: 24px; inline-size: 24px;">
                            <div class="slide-group-stepper-indicator"></div>
                        </div>
                        <h4 class="text-h4 step-number">{{ `0${idx + 1}` }}</h4>
                    </div>
                    <div style="line-height: 0;">
                        <h6 class="text-sm font-weight-medium step-title">{{ step.title }}</h6>
                        <span class="text-xs step-subtitle">{{ step.subTitle }}</span>
                    </div>
                    <div class="slide-group-stepper-link-line" v-if="idx < steps.length - 1"></div>
                </div>
            </div>
        </v-slide-group-item>
    </v-slide-group>
    <v-slide-group class="slide-group-with-stepper mb-3 hidden-sm-and-up" direction="vertical">
        <v-slide-group-item :key="idx" v-for="(step, idx) in steps">
            <div class="mx-1 mb-3"
                 :class="{ 'slide-group-step-active': isStepActive(step), 'slide-group-step-completed': isStepCompleted(idx), 'cursor-pointer': isClickable }"
                 @click="changeStep(step)">
                <div class="d-flex align-center gap-x-2">
                    <div class="d-flex align-center gap-2">
                        <div class="d-flex align-center justify-center" style="block-size: 24px; inline-size: 24px;">
                            <div class="slide-group-stepper-indicator"></div>
                        </div>
                        <h4 class="text-h4 step-number">{{ `0${idx + 1}` }}</h4>
                    </div>
                    <div style="line-height: 0;">
                        <h6 class="text-sm font-weight-medium step-title">{{ step.title }}</h6>
                        <span class="text-xs step-subtitle">{{ step.subTitle }}</span>
                    </div>
                </div>
            </div>
        </v-slide-group-item>
    </v-slide-group>
</template>

<script setup lang="ts">
import { computed } from 'vue';

interface StepBarItem {
    name: string;
    title: string;
    subTitle: string;
}

const props = defineProps<{
    steps: StepBarItem[]
    currentStep: string
    clickable?: string | boolean
    minWidth: string | number
}>();

const emit = defineEmits<{
    (e: 'step:change', stepName: string): void
}>();

const isClickable = computed<boolean>(() => {
    return props.clickable !== 'false' && props.clickable !== false;
});

function changeStep(step: StepBarItem): void {
    if (isClickable.value) {
        emit('step:change', step.name);
    }
}

function isStepActive(step: StepBarItem): boolean {
    return props.currentStep === step.name;
}

function isStepCompleted(stepIndex: number): boolean {
    for (let i = 0; i < props.steps.length; i++) {
        if (props.steps[i].name === props.currentStep) {
            return stepIndex < i;
        }
    }

    return false;
}
</script>

<style>
.slide-group-with-stepper .v-slide-group__content .slide-group-stepper-link-line {
    background-color: rgb(var(--v-theme-primary));
    border-radius: 0.1875rem;
    block-size: .1875rem;
    inline-size: 3.75rem;
    opacity: var(--v-activated-opacity);
}

.slide-group-with-stepper .v-slide-group__content .slide-group-stepper-indicator {
    background-color: rgb(var(--v-theme-surface));
    border: 0.3125rem solid rgb(var(--v-theme-primary));
    border-radius: 50%;
    block-size: 1.25rem;
    inline-size: 1.25rem;
    opacity: var(--v-activated-opacity);
}

.slide-group-with-stepper .v-slide-group__content .slide-group-step-completed .slide-group-stepper-indicator,
.slide-group-with-stepper .v-slide-group__content .slide-group-step-active .slide-group-stepper-indicator,
.slide-group-with-stepper .v-slide-group__content .slide-group-step-completed .slide-group-stepper-link-line {
    opacity: 1;
}
</style>
