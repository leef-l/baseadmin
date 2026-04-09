import { readFileSync } from 'node:fs';
import { resolve } from 'node:path';

import { describe, expect, it } from 'vitest';

import { supportedBaseAdminComponents } from './supported-components';

describe('supportedBaseAdminComponents', () => {
  it('matches the shared baseadmin scope contract', () => {
    const contractPath = resolve(process.cwd(), '../contracts/baseadmin-scope.json');
    const contract = JSON.parse(readFileSync(contractPath, 'utf8')) as {
      supportedComponents: string[];
    };

    expect([...supportedBaseAdminComponents]).toEqual(contract.supportedComponents);
  });
});
