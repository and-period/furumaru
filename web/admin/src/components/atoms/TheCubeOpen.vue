<script lang="ts" setup>
const emit = defineEmits<{
  (e: 'click'): void
}>()

const handleClick = () => {
  emit('click')
}
</script>

<template>
  <div role="button" @click="handleClick">
    <div class="cube-top" />
    <div class="cube-x">
      <div class="cube-left" />
      <div class="cube">
        <slot />
      </div>
      <div class="cube-right" />
    </div>
    <div class="cube-bottom" />
  </div>
</template>

<style lang="scss" scoped>
$bg-color: #e7d0a9;
$border-color: #c1ae8d;

.cube-x {
  display: flex;
  align-items: center;
  justify-content: center;
}
.cube {
  background-color: $bg-color;
  border: 12px groove $border-color;
  height: 84px;
  width: 160px;
  display: inline-block;
}
.cube-top {
  border-bottom: 30px solid $bg-color;
  border-left: 4px solid transparent;
  border-right: 4px solid transparent;
  width: 160px;
  height: 30px;
  margin: auto;
}
.cube-left {
  display: inline-block;
  border-right: 30px solid $bg-color;
  border-top: 4px solid transparent;
  border-bottom: 4px solid transparent;
  width: 12px;
  height: 84px;
}
.cube-right {
  display: inline-block;
  border-left: 30px solid $bg-color;
  border-top: 4px solid transparent;
  border-bottom: 4px solid transparent;
  width: 12px;
  height: 84px;
}
.cube-bottom {
  border-top: 30px solid $bg-color;
  border-left: 4px solid transparent;
  border-right: 4px solid transparent;
  width: 160px;
  height: 30px;
  margin: auto;
}
</style>
