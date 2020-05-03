import { TestBed } from '@angular/core/testing';

import { ExampleGraphqlService } from './example-graphql.service';

describe('ExampleGraphqlService', () => {
  let service: ExampleGraphqlService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(ExampleGraphqlService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
