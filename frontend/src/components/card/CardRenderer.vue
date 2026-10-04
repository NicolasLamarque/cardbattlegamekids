<script setup>
import { computed } from 'vue'
import { resolveTemplate } from './templates/registry.js'

const props = defineProps({
  character: { type: Object, required: true },
  stats: { type: Object, required: true },
  template: { type: String, default: 'classic' },
  accentColor: { type: String, default: '' },
})

const TemplateComponent = computed(() => resolveTemplate(props.template))
</script>

<template>
  <component
    :is="TemplateComponent"
    :character="character"
    :stats="stats"
    :accent-color="accentColor"
  />
</template>
