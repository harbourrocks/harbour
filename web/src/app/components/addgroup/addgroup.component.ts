import { Component, OnInit } from '@angular/core';

@Component({
  selector: 'app-addgroup',
  templateUrl: './addgroup.component.html',
  styleUrls: ['./addgroup.component.scss']
})
export class AddgroupComponent implements OnInit {

  constructor() { }

  ngOnInit(): void {
  }

  addGroup() {
    console.log("Add group here");
  }
}
