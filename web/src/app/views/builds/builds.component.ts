import { Component, OnInit } from '@angular/core';
import { SimpleListItem } from 'src/app/models/simple-list-item.model';

@Component({
  selector: 'app-builds',
  templateUrl: './builds.component.html',
  styleUrls: ['./builds.component.scss']
})
export class BuildsComponent implements OnInit {
  navigationLabels = ["Images","Builds"];

  imageList: SimpleListItem[] = [
    {label: "latest", afterLabel: "2.0.0"},
    {label: "latest", afterLabel: "2.0.0"},
    {label: "latest", afterLabel: "2.0.0"},
  ]

  currentPageIndex = 0;

  constructor() { }

  ngOnInit(): void {
  }
  pageChange(newIndex) {
    this.currentPageIndex = newIndex;
  }

}
