<script setup lang="ts">
interface Props {
	name: string;
}

defineProps<Props>();
</script>

<template>
	<div class="flex flex-col">
		<p>
			{{ name }}
		</p>
		<label for="select">数量</label>
		<select id="select">
			<option value="option1">Option 1</option>
			<option value="option2">Option 2</option>
			<option value="option3">Option 3</option>
		</select>
		<button class="">
			カゴに入れる
		</button>
	</div>
</template>
