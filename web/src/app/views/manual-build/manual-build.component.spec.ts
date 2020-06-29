import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { ManualBuildComponent } from './manual-build.component';

describe('ManualBuildComponent', () => {
  let component: ManualBuildComponent;
  let fixture: ComponentFixture<ManualBuildComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ ManualBuildComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ManualBuildComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
