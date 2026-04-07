import { getScrollbarWidth, needsScrollbar } from '@vben-core/shared/utils';

import {
  useScrollLock as _useScrollLock,
  tryOnBeforeUnmount,
  tryOnMounted,
} from '@vueuse/core';

export const SCROLL_FIXED_CLASS = `_scroll__fixed_`;

export function useScrollLock() {
  const body = typeof document === 'undefined' ? undefined : document.body;
  const isLocked = _useScrollLock(body ?? undefined);
  const scrollbarWidth = getScrollbarWidth();

  tryOnMounted(() => {
    if (!body || !needsScrollbar()) {
      return;
    }
    body.style.paddingRight = `${scrollbarWidth}px`;

    const layoutFixedNodes = document.querySelectorAll<HTMLElement>(
      `.${SCROLL_FIXED_CLASS}`,
    );
    const nodes = [...layoutFixedNodes];
    if (nodes.length > 0) {
      nodes.forEach((node) => {
        node.dataset.transition = node.style.transition;
        node.style.transition = 'none';
        node.style.paddingRight = `${scrollbarWidth}px`;
      });
    }
    isLocked.value = true;
  });

  tryOnBeforeUnmount(() => {
    if (!body || !needsScrollbar()) {
      return;
    }
    isLocked.value = false;
    const layoutFixedNodes = document.querySelectorAll<HTMLElement>(
      `.${SCROLL_FIXED_CLASS}`,
    );
    const nodes = [...layoutFixedNodes];
    if (nodes.length > 0) {
      nodes.forEach((node) => {
        node.style.paddingRight = '';
        requestAnimationFrame(() => {
          if (!node.isConnected) {
            return;
          }
          node.style.transition = node.dataset.transition || '';
        });
      });
    }
    body.style.paddingRight = '';
  });
}
