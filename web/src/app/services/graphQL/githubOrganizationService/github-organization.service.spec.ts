import { TestBed } from '@angular/core/testing';

import { GithubOrganizationService } from './github-organization.service';

describe('GithubOrganizationService', () => {
  let service: GithubOrganizationService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(GithubOrganizationService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
