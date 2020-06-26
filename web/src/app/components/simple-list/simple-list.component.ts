import { Component, OnInit, Input } from '@angular/core';
import { SimpleListItem } from 'src/app/models/simple-list-item.model';
import { faTimes } from '@fortawesome/free-solid-svg-icons';

@Component({
  selector: 'app-simple-list',
  templateUrl: './simple-list.component.html',
  styleUrls: ['./simple-list.component.scss']
})
export class SimpleListComponent implements OnInit {
  @Input() listItems: Array<SimpleListItem>;
  @Input() deletable?: boolean;
  closeIcon = faTimes;

  constructor() {   }

  ngOnInit(): void {
  }

}
