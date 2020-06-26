import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { VerticalNavbarComponent } from './vertical-navbar.component';

describe('VerticalNavbarComponent', () => {
  let component: VerticalNavbarComponent;
  let fixture: ComponentFixture<VerticalNavbarComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ VerticalNavbarComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(VerticalNavbarComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
