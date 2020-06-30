import { Component, OnInit } from '@angular/core';
import { FormModel } from 'src/app/models/form.model';
import {Location} from '@angular/common';

@Component({
  selector: 'app-manual-build',
  templateUrl: './manual-build.component.html',
  styleUrls: ['./manual-build.component.scss']
})
export class ManualBuildComponent implements OnInit {

  public formModel: FormModel = {
    header: "Manual Build",
    items: [
      {name: "scm", placeholder: "Source Control Repo", selections:["1", "2"]},
      {name: "commit", placeholder: "Commit Hash"},
      {name: "dockerPath", placeholder: "Dockerfile Path"},
      {name: "repository", placeholder: "Repository", selections: ["1", "2"]},
      {name: "tag", placeholder: "Image Tag"},
    ],
  }

  constructor(private _location: Location) { }

  ngOnInit(): void {
  }

  onCancel() {
    this._location.back();
  }

  onSubmit() {
    //CALL

    this._location.back();
  }

}
