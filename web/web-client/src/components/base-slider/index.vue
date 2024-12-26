<template>
  <div :class="vertical ? 'slider-vertical' : 'slider-line'">
    <div ref="sliderRef" class="slider-bar" @click="clickSlider($event)">
      <div class="slider-played" :style="vertical ? `height: ${activePercentage}%` : `width: ${activePercentage}%`">
        <div ref="blockRef" class="slider-block" tabindex="0"></div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { watch, onMounted, ref, onBeforeUnmount } from 'vue';

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

// 点击滑动条
const clickSlider = (e: MouseEvent) => {
  let activeSize: number;
  let percentage: number;
  if (props.vertical) {
    activeSize = sliderRef.value!.clientHeight - (e.clientY - sliderRef.value!.getBoundingClientRect().top);
    percentage = Math.round((activeSize / sliderRef.value!.clientHeight) * 100) / 100;
  } else {
    activeSize = e.clientX - sliderRef.value!.getBoundingClientRect().left;
    percentage = Math.round((activeSize / sliderRef.value!.clientWidth) * 100) / 100;
  }

  if (props.step) {
    percentage = Math.round(percentage * 10) / 10
  }

  updateSlider(percentage);
}

// 更新滑块位置和值
const updateSlider = (percentage: number) => {
  percentage = Math.max(0, percentage);
  percentage = Math.min(percentage, 1);

  activePercentage.value = percentage * 100;
  emit('changeValue', Math.round((props.max - props.min) * percentage * 100) / 100);
}

// 滑动滑动条
const slidingSlider = () => {
  blockRef.value!.onmousedown = startSliding;
  blockRef.value!.ontouchstart = startSliding;
  blockRef.value!.onkeydown = handleKeydown;
}

// 开始滑动
const startSliding = (e: Event) => {
  e.preventDefault();
  document.onmousemove = document.ontouchmove = throttle(sliderValueChange, 16);
  document.onmouseup = document.ontouchend = stopSliding;
}

// 停止滑动
const stopSliding = () => {
  document.onmousemove = document.ontouchmove = document.onmouseup = document.ontouchend = null;
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
    percentage = Math.round((activeSize / sliderRef.value!.clientHeight) * 100) / 100;
  } else {
    activeSize = (clientX - sliderRef.value!.getBoundingClientRect().left);
    percentage = Math.round((activeSize / sliderRef.value!.clientWidth) * 100) / 100;
  }

  if (props.step) {
    percentage = Math.round(percentage * 10) / 10
  }

  updateSlider(percentage);
}

// 处理键盘事件
const handleKeydown = (e: KeyboardEvent) => {
  let step = props.step ? 0.1 : 0.01;
  if (e.key === 'ArrowUp' || e.key === 'ArrowRight') {
    updateSlider(Math.min(activePercentage.value / 100 + step, 1));
  } else if (e.key === 'ArrowDown' || e.key === 'ArrowLeft') {
    updateSlider(Math.max(activePercentage.value / 100 - step, 0));
  }
}

// 节流函数
const throttle = (func: Function, limit: number) => {
  let lastFunc: NodeJS.Timeout;
  let lastRan: number;
  return function (...args: any) {
    if (!lastRan) {
      func(...args);
      lastRan = Date.now();
    } else {
      clearTimeout(lastFunc);
      lastFunc = setTimeout(function () {
        if ((Date.now() - lastRan) >= limit) {
          func(...args);
          lastRan = Date.now();
        }
      }, limit - (Date.now() - lastRan));
    }
  }
}

watch(() => props.value, (newValue) => {
  activePercentage.value = Math.round((newValue / (props.max - props.min)) * 10000) / 100;
});

onMounted(() => {
  slidingSlider();
  activePercentage.value = Math.round((props.value / (props.max - props.min)) * 100);
});

onBeforeUnmount(() => {
  stopSliding();
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
