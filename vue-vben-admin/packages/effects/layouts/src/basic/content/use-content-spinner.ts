import { computed, onUnmounted, ref } from 'vue';
import { useRouter } from 'vue-router';

import { preferences } from '@vben/preferences';

function useContentSpinner() {
  const spinning = ref(false);
  const startTime = ref(0);
  let endTimer: null | ReturnType<typeof setTimeout> = null;
  const router = useRouter();
  const minShowTime = 500; // 最小显示时间
  const enableLoading = computed(() => preferences.transition.loading);

  // 结束加载动画
  const onEnd = () => {
    if (!enableLoading.value) {
      return;
    }
    if (endTimer) {
      clearTimeout(endTimer);
      endTimer = null;
    }
    const processTime = performance.now() - startTime.value;
    if (processTime < minShowTime) {
      endTimer = setTimeout(() => {
        spinning.value = false;
        endTimer = null;
      }, minShowTime - processTime);
    } else {
      spinning.value = false;
    }
  };

  // 路由前置守卫
  router.beforeEach((to) => {
    if (to.meta.loaded || !enableLoading.value || to.meta.iframeSrc) {
      return true;
    }
    startTime.value = performance.now();
    spinning.value = true;
    return true;
  });

  // 路由后置守卫
  router.afterEach((to) => {
    if (to.meta.loaded || !enableLoading.value || to.meta.iframeSrc) {
      return true;
    }
    onEnd();
    return true;
  });

  onUnmounted(() => {
    if (endTimer) {
      clearTimeout(endTimer);
      endTimer = null;
    }
  });

  return { spinning };
}

export { useContentSpinner };
