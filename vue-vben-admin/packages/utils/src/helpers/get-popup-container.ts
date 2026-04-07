/**
 * If the node is holding inside a form, return the form element,
 * otherwise return the parent node of the given element or
 * the document body if the element is not provided.
 */
export function getPopupContainer(node?: HTMLElement): HTMLElement {
  const body = typeof document === 'undefined' ? undefined : document.body;
  if (!node?.isConnected) {
    return body ?? node ?? document.createElement('div');
  }

  const form = node.closest('form');
  if (form?.isConnected) {
    return form;
  }

  const parentElement = node.parentElement;
  if (parentElement?.isConnected) {
    return parentElement;
  }

  const parentNode = node.parentNode;
  if (parentNode instanceof HTMLElement && parentNode.isConnected) {
    return parentNode;
  }

  return body ?? node;
}
