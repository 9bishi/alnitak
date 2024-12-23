<template>
  <div :class="vertical ? 'slider-vertical' : 'slider-line'">
    <div ref="sliderRef" class="slider-bar" @click="clickSlider($event)">
      <div class="slider-played" :style="vertical ? `height: ${activePercentage}%` : `width: ${activePercentage}%`">
        <div ref="blockRef" class="slider-block"></div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { watch, onMounted, ref, nextTick, onBeforeUnmount, computed } from 'vue';

const emit = defineEmits(["changeValue"]);
const props = withDefaults(defineProps<{
  value: number
  max?: number
  min?: number
  vertical?: boolean
  step?: boolean
}>(), {
  value: 0,
  max: 100,
  min: 0,
  vertical: false,
  step: false
});

const activePercentage = ref(props.value);
const blockRef = ref<HTMLElement | null>(null);
const sliderRef = ref<HTMLElement | null>(null);

// 步进处理：根据step计算步进
const calculateStep = (percentage: number) => {
  if (props.step) {
    const stepValue = (props.max - props.min) / 10;  // 假设步进值为 1/10
    return Math.round(percentage * 100) / 100 * stepValue;
  }
  return percentage;
}

// 点击滑动条
const clickSlider = (e: MouseEvent) => {
  let activeSize: number;
  let percentage: number;
  if (props.vertical) {
    activeSize = sliderRef.value!.clientHeight - (e.clientY - sliderRef.value!.getBoundingClientRect().top);
    percentage = (activeSize / sliderRef.value!.clientHeight);
  } else {
    activeSize = e.clientX - sliderRef.value!.getBoundingClientRect().left;
    percentage = (activeSize / sliderRef.value!.clientWidth);
  }

  percentage = Math.max(0, Math.min(percentage, 1));

  activePercentage.value = percentage * 100;
  emit('changeValue', Math.round((props.max - props.min) * calculateStep(percentage) + props.min));
}

// 滑动滑动条
const slidingSlider = () => {
  // PC端
  blockRef.value!.onmousedown = function () {
    document.onmousemove = function (e) {
      sliderValueChange(e);
    };
    document.onmouseup = function () {
      document.onmousemove = document.onmouseup = null;
    };
  };

  // 移动端
  blockRef.value!.ontouchstart = function () {
    document.ontouchmove = function (e) {
      sliderValueChange(e);
    };
    document.ontouchend = function () {
      document.ontouchmove = document.ontouchend = null;
    };
  };
}

// 滑动条值改变
const sliderValueChange = (e: MouseEvent | TouchEvent) => {
  let activeSize: number;
  let percentage: number;
  let clientX, clientY: number;

  if (e instanceof MouseEvent) {
    clientX = e.clientX;
    clientY = e.clientY;
  } else {
    clientX = e.changedTouches[0].clientX;
    clientY = e.changedTouches[0].clientY;
  }

  if (props.vertical) {
    activeSize = sliderRef.value!.clientHeight - (clientY - sliderRef.value!.getBoundingClientRect().top);
    percentage = (activeSize / sliderRef.value!.clientHeight);
  } else {
    activeSize = (clientX - sliderRef.value!.getBoundingClientRect().left);
    percentage = (activeSize / sliderRef.value!.clientWidth);
  }

  percentage = Math.max(0, Math.min(percentage, 1));

  activePercentage.value = percentage * 100;
  emit('changeValue', Math.round((props.max - props.min) * calculateStep(percentage) + props.min));
}

// 监听value的变化，更新滑动条位置
watch(() => props.value, (newValue) => {
  activePercentage.value = Math.round(((newValue - props.min) / (props.max - props.min)) * 100);
});

onMounted(() => {
  slidingSlider();
  activePercentage.value = Math.round(((props.value - props.min) / (props.max - props.min)) * 100);
});

onBeforeUnmount(() => {
  // 清除 touch 和 mouse 事件
  document.ontouchmove = document.ontouchend = null;
  document.onmousemove = document.onmouseup = null;
})
</script>

<style lang="scss" scoped>
.slider-line {
  margin-left: 6px;
  width: calc(100% - 12px);

  .slider-bar {
    position: relative;
    width: 100%;
    background: rgba(107, 107, 107, 0.6);
    height: 4px;
    cursor: pointer;
  }

  .slider-played {
    position: absolute;
    height: 4px;
    background-color: var(--primary-color);

    .slider-block {
      position: absolute;
      top: -6px;
      right: -6px;
      width: 12px;
      height: 12px;
      border-radius: 50%;
      background: #fff;
      border: 2px solid var(--primary-color);
      transition: 0.2s all;
      user-select: none;
    }
  }
}

.slider-vertical {
  height: calc(100% - 24px);
  margin: 12px 0;

  .slider-bar {
    position: relative;
    width: 4px;
    height: 100%;
    cursor: pointer;
    background: rgba(107, 107, 107, 0.6);
  }

  .slider-played {
    position: absolute;
    width: 100%;
    bottom: 0;

    .slider-block {
      position: absolute;
      top: -4px;
      right: -4px;
      width: 8px;
      height: 8px;
      border-radius: 50%;
      border: 2px solid;
      background: #fff;
      transition: 0.2s all;
      user-select: none;
    }
  }
}
</style>
