import tinymceCore from 'tinymce';

if (typeof window !== 'undefined') {
  (window as any).tinymce = tinymceCore;
}

export default tinymceCore;
