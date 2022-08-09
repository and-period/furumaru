<template>
  <div>
    <div class="cube" role="button" @click="handleClick">
      <div class="cube-item" />
      <div class="cube-item" />
      <div class="cube-item text-center">
        <img height="45px" src="/furulogo.png" alt="ふるまるロゴ" />
      </div>
      <div class="cube-item text-center"></div>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent } from '@vue/composition-api'

export default defineComponent({
  setup(_, { emit }) {
    const handleClick = () => {
      emit('click')
    }

    return {
      handleClick,
    }
  },
})
</script>

<style lang="scss" scoped>
$bg-color: #e7d0a9;
$tape-bg-color: #b89f89;
$border-color: #c1ae8d;

.cube {
  display: block;
  position: relative;
  margin: auto;
  padding: 40px;
  width: 120px;
  transform-style: preserve-3d;
  transform: rotateX(-25deg) rotateY(-30deg);
}

.cube .cube-item {
  position: absolute;
  left: 0;
  right: 0;
  box-sizing: border-box;
  border: 1px solid $border-color;
  width: 100%;
  height: 60%;
  background-color: $bg-color;
}

// 上
.cube .cube-item:nth-child(1) {
  transform: translate3d(0, -40px, 0) rotateX(-90deg);
  height: 80px;
}

// ガムテープ部分
.cube .cube-item:nth-child(2) {
  transform: translate3d(0, -10px, 0) rotateX(-90deg);
  height: 20px;
  background-color: $tape-bg-color;
}

// 左
.cube .cube-item:nth-child(3) {
  transform: translate3d(0, 0, 40px);
}

// 右
.cube .cube-item:nth-child(4) {
  transform: translate3d(80px, 0, 0) rotateY(90deg);
  width: 80px;
}
</style>
