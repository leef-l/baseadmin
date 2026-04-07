import { describe, expect, it } from 'vitest';

import { parseUploadValue } from './shared';

describe('parseUploadValue', () => {
  it('trims blank items and removes duplicates', () => {
    expect(
      parseUploadValue(' https://a.png , ,https://b.png,https://a.png ,  '),
    ).toEqual(['https://a.png', 'https://b.png']);
  });

  it('returns empty array for blank input', () => {
    expect(parseUploadValue('   ')).toEqual([]);
    expect(parseUploadValue()).toEqual([]);
  });
});
