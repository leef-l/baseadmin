export function normalizeRoleIds(input: unknown): string[] {
  if (!Array.isArray(input)) {
    return [];
  }

  const values = new Set<string>();
  const queue = [...input];

  while (queue.length > 0) {
    const current = queue.shift();
    if (Array.isArray(current)) {
      queue.push(...current);
      continue;
    }

    if (
      current &&
      typeof current === 'object' &&
      'value' in current &&
      typeof (current as { value?: unknown }).value !== 'undefined'
    ) {
      queue.push((current as { value?: unknown }).value);
      continue;
    }

    if (
      typeof current === 'bigint' ||
      typeof current === 'number' ||
      typeof current === 'string'
    ) {
      values.add(String(current));
    }
  }

  return [...values];
}
