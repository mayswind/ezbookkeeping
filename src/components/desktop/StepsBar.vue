<template>
    <div class="d-flex" :style="`min-width: ${minWidth}px`" v-if="minWidth"></div>
    <v-slide-group class="slide-group-with-stepper" :class="{ 'hidden-xs': !alwaysHorizontal }" show-arrows>
        <v-slide-group-item :key="idx" v-for="(step, idx) in steps">
            <div class="mx-1"
                 :class="{ 'slide-group-step-active': isStepActive(step), 'slide-group-step-completed': isStepCompleted(idx), 'cursor-pointer': isClickable }"
                 @click="changeStep(step)">
                <div class="d-flex align-center gap-x-2">
                    <div class="d-flex align-center gap-2">
                        <div class="d-flex align-center justify-center" style="block-size: 24px; inline-size: 24px;">
                            <div class="slide-group-stepper-indicator"></div>
                        </div>
                        <span class="text-headline-medium step-number">{{ getDisplayStep(idx + 1) }}</span>
                    </div>
                    <div style="line-height: 0;">
                        <div class="text-body-medium font-weight-medium">{{ step.title }}</div>
                        <div class="text-body-small text-medium-emphasis">{{ step.subTitle }}</div>
                    </div>
                    <div class="slide-group-stepper-link-line" v-if="idx < steps.length - 1"></div>
                </div>
            </div>
        </v-slide-group-item>
    </v-slide-group>
    <v-slide-group class="slide-group-with-stepper hidden-sm-and-up" direction="vertical" v-if="!alwaysHorizontal">
        <v-slide-group-item :key="idx" v-for="(step, idx) in steps">
            <div class="mx-1 mb-3"
                 :class="{ 'slide-group-step-active': isStepActive(step), 'slide-group-step-completed': isStepCompleted(idx), 'cursor-pointer': isClickable }"
                 @click="changeStep(step)">
                <div class="d-flex align-center gap-x-2">
                    <div class="d-flex align-center gap-2">
                        <div class="d-flex align-center justify-center" style="block-size: 24px; inline-size: 24px;">
                            <div class="slide-group-stepper-indicator"></div>
                        </div>
                        <span class="text-headline-medium step-number">{{ getDisplayStep(idx + 1) }}</span>
                    </div>
                    <div style="line-height: 0;">
                        <div class="text-body-medium font-weight-medium">{{ step.title }}</div>
                        <div class="text-body-small text-medium-emphasis">{{ step.subTitle }}</div>
                    </div>
                </div>
            </div>
        </v-slide-group-item>
    </v-slide-group>
</template>

<script setup lang="ts">
import { computed } from 'vue';

import { useI18n } from '@/locales/helpers.ts';

import { itemAndIndex } from '@/core/base.ts';
import { NumeralSystem } from '@/core/numeral.ts';

export interface StepBarItem {
    name: string;
    title: string;
    subTitle: string;
}

const props = defineProps<{
    steps: StepBarItem[];
    currentStep: string;
    clickable?: string | boolean;
    minWidth: string | number;
    alwaysHorizontal?: boolean;
}>();

const emit = defineEmits<{
    (e: 'step:change', stepName: string): void;
}>();

const { getCurrentNumeralSystemType } = useI18n();

const numeralSystem = computed<NumeralSystem>(() => getCurrentNumeralSystemType());
const isClickable = computed<boolean>(() => props.clickable !== 'false' && props.clickable !== false);

function getDisplayStep(index: number): string {
    return numeralSystem.value.replaceWesternArabicDigitsToLocalizedDigits(index.toString().padStart(2, NumeralSystem.WesternArabicNumerals.digitZero));
}

function changeStep(step: StepBarItem): void {
    if (isClickable.value) {
        emit('step:change', step.name);
    }
}

function isStepActive(step: StepBarItem): boolean {
    return props.currentStep === step.name;
}

function isStepCompleted(stepIndex: number): boolean {
    for (const [step, index] of itemAndIndex(props.steps)) {
        if (step.name === props.currentStep) {
            return stepIndex < index;
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
