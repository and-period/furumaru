<script setup lang="ts">
interface Props {
	name: string;
}

defineProps<Props>();
</script>

<template>
	<div>
		{{ name }}
	</div>
</template>
