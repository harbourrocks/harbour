import { Component, OnInit, Input, Output, EventEmitter } from '@angular/core';
import { ListItem } from 'src/app/models/list-item.model';


@Component({
  selector: 'app-list-item',
  templateUrl: './list-item.component.html',
  styleUrls: ['./list-item.component.scss']
})
export class ListItemComponent implements OnInit {
  @Input() listItem: ListItem;
  @Output() onClick = new EventEmitter<ListItem>();

  constructor() { }

  ngOnInit(): void {
  }
}
