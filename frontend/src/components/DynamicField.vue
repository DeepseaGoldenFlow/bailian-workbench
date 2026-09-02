<template>
  <div class="dynamic-field" :class="{ 'dynamic-field--switch': field.type === 'bool' }">
    <v-textarea v-if="isLongText" :model-value="modelValue" :label="label" :placeholder="field.placeholder" :hint="field.description" persistent-hint rows="field.type==='string' ? 4 : 3" variant="outlined" @update:model-value="$emit('update:modelValue',$event)" />
    <v-select v-else-if="field.type==='select'" :model-value="modelValue" :items="field.options" item-title="label" item-value="value" :label="label" :hint="field.description" persistent-hint variant="outlined" @update:model-value="$emit('update:modelValue',$event)" />
    <v-switch v-else-if="field.type==='bool'" :model-value="modelValue" :label="label" :hint="field.description" persistent-hint color="primary" inset @update:model-value="$emit('update:modelValue',$event)" />
    <v-text-field v-else :model-value="modelValue" :type="isNumber ? 'number' : 'text'" :label="label" :min="field.min" :max="field.max" :step="field.step" :placeholder="field.placeholder" :hint="field.description" persistent-hint variant="outlined" @update:model-value="$emit('update:modelValue',$event)" />
  </div>
</template>

<script setup>
import { computed } from 'vue'
const props=defineProps({field:{type:Object,required:true},modelValue:{default:''}})
defineEmits(['update:modelValue'])
const label=computed(()=>`${props.field.label}${props.field.required?' *':''}`)
const isLongText=computed(()=>props.field.type==='media_list'||props.field.type==='json'||(props.field.type==='string'&&['prompt','negative_prompt'].includes(props.field.name)))
const isNumber=computed(()=>['int','float'].includes(props.field.type))
</script>

<style scoped>
.dynamic-field { margin-bottom:10px; }.dynamic-field :deep(.v-field) { border-radius:14px; }.dynamic-field :deep(.v-messages) { min-height:17px; font-size:10px; }.dynamic-field--switch { margin-top:-2px; margin-bottom:2px; }
</style>
