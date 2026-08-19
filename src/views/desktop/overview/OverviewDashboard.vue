<template>
    <div ref="grid" class="overview-dashboard"
         :class="{ 'overview-dashboard-editing': editing, 'overview-dashboard-empty': !layout.widgets.length }"
         :style="gridStyle">
        <div class="overview-dashboard-drop-placeholder" :style="getGridPositionStyle(draggingWidget)"
             v-if="draggingWidget"></div>
        <div class="overview-dashboard-item"
             :class="{ 'overview-dashboard-item-editing': editing, 'overview-dashboard-item-dragging': draggingWidget?.id === widget.id }"
             :style="getWidgetStyle(widget)" :key="widget.id" v-for="widget in sortedWidgets">
            <overview-widget class="overview-dashboard-widget" :widget="widget" :loading="loading" :editing="editing"
                             @refresh="$emit('refresh')" />
            <template v-if="editing">
                <div class="overview-dashboard-title-drag-area" :aria-label="tt('Move')"
                     @pointerdown="startPointerAction($event, widget, 'move')"></div>
                <div class="overview-dashboard-editor-toolbar">
                    <v-btn density="comfortable" color="default" variant="text" class="ma-2" :icon="true"
                           :aria-label="tt('More')">
                        <v-icon :icon="mdiDotsVertical" />
                        <v-tooltip activator="parent">{{ tt('More') }}</v-tooltip>
                        <v-menu activator="parent">
                            <v-list>
                                <template v-if="DESKTOP_OVERVIEW_WIDGET_DEFINITIONS[widget.type].supportsSettings">
                                    <v-list-item :prepend-icon="mdiCogOutline" :title="tt('Settings')"
                                                 @click="$emit('configure', widget)" />
                                    <v-divider class="my-2" />
                                </template>
                                <v-list-item :prepend-icon="mdiDeleteOutline" :title="tt('Delete')"
                                             @click="$emit('remove', widget.id)" />
                            </v-list>
                        </v-menu>
                    </v-btn>
                </div>
                <button class="overview-dashboard-resize-handle" :aria-label="tt('Resize')"
                        @pointerdown="startPointerAction($event, widget, 'resize')">
                    <v-icon :icon="mdiResizeBottomRight" size="32" />
                </button>
            </template>
        </div>

        <div class="d-flex flex-column align-center justify-center ga-3" v-if="!layout.widgets.length">
            <v-icon :icon="mdiWidgetsOutline" size="64" color="secondary" />
            <span class="text-title-medium">{{ tt('No widgets') }}</span>
            <v-btn color="primary" variant="tonal" @click="$emit('add')" v-if="editing">{{ tt('Add Widget') }}</v-btn>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref, computed, onBeforeUnmount } from 'vue';

import OverviewWidget from './OverviewWidget.vue';

import { useI18n } from '@/locales/helpers.ts';

import {
    type DesktopOverviewWidgetDefinition,
    type DesktopOverviewLayout,
    type DesktopOverviewWidgetLayout
} from '@/core/overview_layout.ts';
import {
    DESKTOP_OVERVIEW_LAYOUT_COLUMNS,
    DESKTOP_OVERVIEW_WIDGET_DEFINITIONS
} from '@/consts/overview_layout.ts';

import {
    resolveOverviewWidgetCollisions,
    compactOverviewWidgets
} from '@/lib/overview_layout.ts';

import {
    mdiCogOutline,
    mdiDeleteOutline,
    mdiDotsVertical,
    mdiResizeBottomRight,
    mdiWidgetsOutline
} from '@mdi/js';

interface PointerActionState {
    pointerId: number;
    action: 'move' | 'resize';
    widget: DesktopOverviewWidgetLayout;
    startX: number;
    startY: number;
    startRect: {
        left: number;
        top: number;
        width: number;
        height: number;
    };
}

interface DraggingPreview {
    id: string;
    left: number;
    top: number;
    width: number;
    height: number;
}

const props = defineProps<{
    layout: DesktopOverviewLayout;
    loading: boolean;
    editing?: boolean
}>();

const emit = defineEmits<{
    (e: 'update:layout', value: DesktopOverviewLayout): void;
    (e: 'configure', value: DesktopOverviewWidgetLayout): void;
    (e: 'add'): void;
    (e: 'remove', value: string): void;
    (e: 'refresh'): void;
}>();

const { tt } = useI18n();

const ROW_HEIGHT: number = 60;
const GAP: number = 16;

const grid = ref<HTMLElement | null>(null);
const pointerAction = ref<PointerActionState | null>(null);
const draggingPreview = ref<DraggingPreview | null>(null);

const rowCount = computed<number>(() => Math.max(1, ...props.layout.widgets.map(widget => widget.y + widget.h)));
const sortedWidgets = computed<DesktopOverviewWidgetLayout[]>(() => [...props.layout.widgets].sort((a, b) => a.y - b.y || a.x - b.x || a.id.localeCompare(b.id)));
const draggingWidget = computed<DesktopOverviewWidgetLayout | null>(() => draggingPreview.value ? props.layout.widgets.find(widget => widget.id === draggingPreview.value?.id) || null : null);

const gridStyle = computed<Record<string, string>>(() => ({
    '--overview-dashboard-row-height': `${ROW_HEIGHT}px`,
    '--overview-dashboard-gap': `${GAP}px`,
    minHeight: props.layout.widgets.length ? `${rowCount.value * ROW_HEIGHT + (rowCount.value - 1) * GAP}px` : '360px'
}));

function getGridPositionStyle(widget: DesktopOverviewWidgetLayout): Record<string, string> {
    return {
        gridColumn: `${widget.x + 1} / span ${widget.w}`,
        gridRow: `${widget.y + 1} / span ${widget.h}`
    };
}

function getWidgetStyle(widget: DesktopOverviewWidgetLayout): Record<string, string> {
    const style: Record<string, string> = getGridPositionStyle(widget);
    const preview: DraggingPreview | null = draggingPreview.value;

    if (preview?.id === widget.id) {
        style['position'] = 'fixed';
        style['left'] = `${preview.left}px`;
        style['top'] = `${preview.top}px`;
        style['width'] = `${preview.width}px`;
        style['height'] = `${preview.height}px`;
    }

    return style;
}

function startPointerAction(event: PointerEvent, widget: DesktopOverviewWidgetLayout, action: 'move' | 'resize'): void {
    if (!props.editing || !grid.value) {
        return;
    }

    event.preventDefault();

    if (!event.currentTarget || !(event.currentTarget instanceof HTMLElement)) {
        return;
    }

    const target: HTMLElement = event.currentTarget;
    target.setPointerCapture?.(event.pointerId);

    const itemElement = target.closest('.overview-dashboard-item') as HTMLElement | null;
    const itemRect = itemElement?.getBoundingClientRect();

    if (!itemRect) {
        return;
    }

    pointerAction.value = {
        pointerId: event.pointerId,
        action,
        widget: { ...widget, settings: { ...widget.settings } },
        startX: event.clientX,
        startY: event.clientY,
        startRect: {
            left: itemRect.left,
            top: itemRect.top,
            width: itemRect.width,
            height: itemRect.height
        }
    };

    if (action === 'move') {
        draggingPreview.value = {
            id: widget.id,
            left: itemRect.left,
            top: itemRect.top,
            width: itemRect.width,
            height: itemRect.height
        };
    }

    window.addEventListener('pointermove', handlePointerMove);
    window.addEventListener('pointerup', finishPointerAction);
    window.addEventListener('pointercancel', finishPointerAction);
}

function handlePointerMove(event: PointerEvent): void {
    const state = pointerAction.value;

    if (!state || event.pointerId !== state.pointerId || !grid.value) {
        return;
    }

    const columnWidth: number = (grid.value.clientWidth - GAP * (DESKTOP_OVERVIEW_LAYOUT_COLUMNS - 1)) / DESKTOP_OVERVIEW_LAYOUT_COLUMNS;
    const deltaColumns: number = Math.round((event.clientX - state.startX) / (columnWidth + GAP));
    const deltaRows: number = Math.round((event.clientY - state.startY) / (ROW_HEIGHT + GAP));
    const definition: DesktopOverviewWidgetDefinition = DESKTOP_OVERVIEW_WIDGET_DEFINITIONS[state.widget.type];
    const nextWidget = { ...state.widget, settings: { ...state.widget.settings } };

    if (state.action === 'move') {
        draggingPreview.value = {
            id: state.widget.id,
            left: state.startRect.left + event.clientX - state.startX,
            top: state.startRect.top + event.clientY - state.startY,
            width: state.startRect.width,
            height: state.startRect.height
        };
        nextWidget.x = Math.min(DESKTOP_OVERVIEW_LAYOUT_COLUMNS - nextWidget.w, Math.max(0, state.widget.x + deltaColumns));
        nextWidget.y = Math.max(0, state.widget.y + deltaRows);
    } else {
        nextWidget.w = Math.min(definition.maxWidth ?? DESKTOP_OVERVIEW_LAYOUT_COLUMNS, Math.max(definition.minWidth, state.widget.w + deltaColumns));
        nextWidget.w = Math.min(nextWidget.w, DESKTOP_OVERVIEW_LAYOUT_COLUMNS - nextWidget.x);
        nextWidget.h = Math.min(definition.maxHeight ?? 1000, Math.max(definition.minHeight, state.widget.h + deltaRows));
    }

    const widgets = props.layout.widgets.map(widget => widget.id === nextWidget.id ? nextWidget : widget);
    emit('update:layout', { ...props.layout, widgets: resolveOverviewWidgetCollisions(widgets, nextWidget.id) });
}

function finishPointerAction(event: PointerEvent): void {
    if (!pointerAction.value || event.pointerId !== pointerAction.value.pointerId) {
        return;
    }

    pointerAction.value = null;
    draggingPreview.value = null;

    window.removeEventListener('pointermove', handlePointerMove);
    window.removeEventListener('pointerup', finishPointerAction);
    window.removeEventListener('pointercancel', finishPointerAction);

    emit('update:layout', { ...props.layout, widgets: compactOverviewWidgets(props.layout.widgets) });
}

onBeforeUnmount(() => {
    window.removeEventListener('pointermove', handlePointerMove);
    window.removeEventListener('pointerup', finishPointerAction);
    window.removeEventListener('pointercancel', finishPointerAction);
});
</script>
