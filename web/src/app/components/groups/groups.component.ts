import { Component, OnInit, Input, ViewChild, AfterViewInit } from '@angular/core';
import { GroupModel } from 'src/app/models/group.model';

@Component({
  selector: 'app-groups',
  templateUrl: './groups.component.html',
  styleUrls: ['./groups.component.scss']
})
export class GroupsComponent implements OnInit, AfterViewInit {
  @Input() groupData: GroupModel;
  @ViewChild('content') content;
  @ViewChild('expandable') expandable;

  public isHovered:boolean;
  public isCollapsed:boolean;

  public maxHeight: number;
  public blocked = false;

  constructor() {}

  ngAfterViewInit(): void {
    const el: HTMLElement = this.content.nativeElement;
    this.maxHeight = el.offsetHeight;
    console.log(this.maxHeight, el.style.maxHeight,el);
  }

  ngOnInit(): void {
  }

  onHovered(event){
    this.isHovered = !this.isHovered;
  }

  async onCollapse() {
    if(this.blocked) return;
    this.blocked = true;
    
    const el = this.content.nativeElement;
    el.style.maxHeight = `${this.isCollapsed ? 0 : this.maxHeight}px`;
    await this.sleep(10);
    this.isCollapsed = !this.isCollapsed;
    el.style.maxHeight = `${this.isCollapsed ? 0 : this.maxHeight}px`;

    if(!this.isCollapsed){
      await this.sleep(200);
      el.style.maxHeight = "";
      console.log("NOW")
    }

    this.blocked = false;
    
  }

  async sleep(ms: number){
    await new Promise(r => setTimeout(r, ms));
  }
}
