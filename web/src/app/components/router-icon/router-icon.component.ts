import { Component, OnInit, Input } from '@angular/core';

@Component({
  selector: 'app-router-icon',
  templateUrl: './router-icon.component.html',
  styleUrls: ['./router-icon.component.scss']
})
export class RouterIconComponent implements OnInit {
  @Input() routerLink!: string;
  @Input() icon!: string;

  constructor() { }

  ngOnInit(): void {
  }

}
