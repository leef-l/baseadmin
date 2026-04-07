/**
 * 加载js文件
 * @param src js文件地址
 */
function loadScript(src: string) {
  return new Promise<void>((resolve, reject) => {
    if (typeof document === 'undefined' || !document.head) {
      reject(new Error('Document is not available.'));
      return;
    }
    if (document.querySelector(`script[src="${src}"]`)) {
      // 如果已经加载过，直接 resolve
      return resolve();
    }
    const script = document.createElement('script');
    script.src = src;
    const handleLoad = () => {
      script.removeEventListener('load', handleLoad);
      script.removeEventListener('error', handleError);
      resolve();
    };
    const handleError = () => {
      script.removeEventListener('load', handleLoad);
      script.removeEventListener('error', handleError);
      reject(new Error(`Failed to load script: ${src}`));
    };
    script.addEventListener('load', handleLoad, { once: true });
    script.addEventListener('error', handleError, { once: true });
    document.head.append(script);
  });
}

export { loadScript };
