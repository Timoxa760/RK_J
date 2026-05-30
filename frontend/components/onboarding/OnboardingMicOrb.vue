<script setup lang="ts">
const props = withDefaults(
  defineProps<{
    listening?: boolean
    parsing?: boolean
    disabled?: boolean
    preview?: boolean
    label?: string
  }>(),
  {
    listening: false,
    parsing: false,
    disabled: false,
    preview: false,
    label: 'Нажмите и ответьте вслух'
  }
)

const emit = defineEmits<{
  click: []
}>()

const statusText = computed(() => {
  if (props.parsing) return 'Записываем…'
  if (props.listening) return 'Поток слушает'
  if (props.preview) return 'Голосом — быстрее'
  return props.label
})

function onClick() {
  if (props.disabled || props.parsing || props.preview) return
  emit('click')
}
</script>

<template>
  <div class="mm-onb-mic-wrap">
    <div
      v-if="!preview"
      class="mm-onb-mic-status"
      :class="{ 'opacity-80': parsing }"
    >
      <span
        class="mm-onb-mic-status__dot"
        :class="{ 'mm-onb-mic-status__dot--listen': listening && !parsing }"
      />
      <span>{{ statusText }}</span>
    </div>

    <button
      type="button"
      class="mm-onb-mic-orb-hit"
      :class="{
        'mm-onb-mic-orb-hit--preview': preview,
        'mm-onb-mic-orb-hit--listen': listening && !parsing,
        'mm-onb-mic-orb-hit--parse': parsing
      }"
      :disabled="disabled || parsing"
      :aria-pressed="listening"
      :aria-label="listening ? 'Остановить запись' : label"
      @click="onClick"
    >
      <OnboardingMicOrbVisual
        :listening="listening"
        :parsing="parsing"
        :compact="preview"
        gentle
        :ambient="preview || (!listening && !parsing)"
      />
    </button>

    <p
      v-if="preview"
      class="text-center text-xs font-medium text-[color:var(--mm-primary)]"
    >
      {{ statusText }}
    </p>
    <p
      v-else-if="!parsing && !listening"
      class="max-w-xs text-center text-sm text-[color:var(--mm-text-muted)]"
    >
      {{ label }}
    </p>
  </div>
</template>
