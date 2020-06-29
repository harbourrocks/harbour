import { Component, OnInit, Input } from '@angular/core';
import { ListItem } from 'src/app/models/list-item.model';
import { List } from 'src/app/models/list.model';

@Component({
  selector: 'app-list',
  templateUrl: './list.component.html',
  styleUrls: ['./list.component.scss']
})
export class ListComponent implements OnInit {
  @Input() listModel: List;

  constructor() { }

  ngOnInit(): void {
  }

  onClick(listItem: ListItem) {
    if (this.listModel.clickHandler)
      this.listModel.clickHandler(listItem)
  }

}
