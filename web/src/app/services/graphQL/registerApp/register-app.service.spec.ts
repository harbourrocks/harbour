import { TestBed } from '@angular/core/testing';

import { RegisterAppService } from './register-app.service';

describe('RegisterAppService', () => {
  let service: RegisterAppService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(RegisterAppService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
