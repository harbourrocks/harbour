import { Component, OnInit, Input, Output, EventEmitter } from '@angular/core';
import { FormModel } from 'src/app/models/form.model';
import { FormBuilder, Form, FormGroup, Validators } from '@angular/forms';

@Component({
  selector: 'app-abstract-form',
  templateUrl: './abstract-form.component.html',
  styleUrls: ['./abstract-form.component.scss']
})
export class AbstractFormComponent implements OnInit {
  @Input() formModel: FormModel;
  @Output() submit = new EventEmitter<any>();
  @Output() cancel = new  EventEmitter<null>();

  currentSelectedField: string;
  formBuild: FormGroup;

  constructor(private fb: FormBuilder) {

  }

  ngOnInit(): void {
    let formatData = {};
    this.formModel.items.forEach(item => {
      formatData[item.name] = [null, Validators.required]
    })

    this.formBuild = this.fb.group(formatData)
  }

  selectField(name: string) {
    if (name === this.currentSelectedField) { this.currentSelectedField = undefined; return }
    this.currentSelectedField = name;
  }

  onAdd() {
      this.submit.emit(this.formBuild.value);
  } 

  onCancel() {
      this.cancel.emit();
  }

}
