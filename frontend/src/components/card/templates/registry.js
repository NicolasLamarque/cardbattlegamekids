import ClassicTemplate from './ClassicTemplate.vue'

export const templates = {
  classic: ClassicTemplate,
}

export function resolveTemplate(name) {
  return templates[name] ?? templates.classic
}
