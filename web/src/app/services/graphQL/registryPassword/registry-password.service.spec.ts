import { TestBed } from '@angular/core/testing';

import { RegistryPasswordService } from './registry-password.service';

describe('RegistryPasswordService', () => {
  let service: RegistryPasswordService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(RegistryPasswordService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
