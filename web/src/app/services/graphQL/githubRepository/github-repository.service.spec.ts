import { TestBed } from '@angular/core/testing';

import { GithubRepositoryService } from './github-repository.service';

describe('GithubRepositoryService', () => {
  let service: GithubRepositoryService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(GithubRepositoryService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
