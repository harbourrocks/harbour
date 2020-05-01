import { TestBed } from '@angular/core/testing';

import { OidcCallbackGuard } from './oidc-callback.guard';

describe('OidcCallbackGuard', () => {
  let guard: OidcCallbackGuard;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    guard = TestBed.inject(OidcCallbackGuard);
  });

  it('should be created', () => {
    expect(guard).toBeTruthy();
  });
});
