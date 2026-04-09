import { describe, expect, it } from 'vitest';

import { extractAccessCodes } from './auth-info';

describe('extractAccessCodes', () => {
  it('returns perms from auth info payload', () => {
    expect(
      extractAccessCodes({
        perms: ['system:dept:list', 'system:user:update'],
      }),
    ).toEqual(['system:dept:list', 'system:user:update']);
  });

  it('returns empty array when perms is missing', () => {
    expect(extractAccessCodes(undefined)).toEqual([]);
    expect(extractAccessCodes({})).toEqual([]);
  });
});
