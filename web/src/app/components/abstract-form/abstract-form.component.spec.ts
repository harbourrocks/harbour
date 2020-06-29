import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { AbstractFormComponent } from './abstract-form.component';

describe('AbstractFormComponent', () => {
  let component: AbstractFormComponent;
  let fixture: ComponentFixture<AbstractFormComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ AbstractFormComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(AbstractFormComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
