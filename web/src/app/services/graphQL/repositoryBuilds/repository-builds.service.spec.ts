import { TestBed } from '@angular/core/testing';

import { RepositoryBuildsService } from './repository-builds.service';

describe('RepositoryBuildsService', () => {
  let service: RepositoryBuildsService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(RepositoryBuildsService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
