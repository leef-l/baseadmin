import { describe, expect, it } from 'vitest';

import { normalizeRoleIds } from './role-ids';

describe('normalizeRoleIds', () => {
  it('flattens nested values and removes duplicates', () => {
    expect(
      normalizeRoleIds([
        '1',
        2,
        ['3', { value: '4' }, [{ value: 2 }]],
        { value: ['5', '1'] },
      ]),
    ).toEqual(['1', '2', '3', '4', '5']);
  });

  it('returns empty array for unsupported input', () => {
    expect(normalizeRoleIds(null)).toEqual([]);
    expect(normalizeRoleIds('1')).toEqual([]);
  });
});
