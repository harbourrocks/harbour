import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { RouterIconComponent } from './router-icon.component';

describe('RouterIconComponent', () => {
  let component: RouterIconComponent;
  let fixture: ComponentFixture<RouterIconComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ RouterIconComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(RouterIconComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
