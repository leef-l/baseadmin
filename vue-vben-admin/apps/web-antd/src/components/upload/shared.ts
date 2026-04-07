export function parseUploadValue(value?: string) {
  if (!value) {
    return [];
  }
  const seen = new Set<string>();
  const items: string[] = [];
  for (const part of value.split(',')) {
    const normalized = part.trim();
    if (!normalized || seen.has(normalized)) {
      continue;
    }
    seen.add(normalized);
    items.push(normalized);
  }
  return items;
}
