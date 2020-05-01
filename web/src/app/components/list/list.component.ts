import { Component, OnInit, Input, HostListener, ViewChild, ComponentFactoryResolver, Output, EventEmitter } from '@angular/core';
import { ListModel } from 'src/app/models/list.model';
import { ListDirective } from './list.directive';
import { ListContentModel } from 'src/app/models/list-item';
import { NavbarComponent } from '../navbar/navbar.component';

@Component({
  selector: 'app-list',
  templateUrl: './list.component.html',
  styleUrls: ['./list.component.scss']
})
export class ListComponent implements OnInit {
  public isHovered = false;
  @Input() listData: ListModel;
  @Input() groupHovered: Boolean;
  @Input() groupSmallColorBlock: Boolean;
  @Output() hovered = new EventEmitter<boolean>();

  @ViewChild('expandable') expandable: { nativeElement: any; }; 
  @ViewChild('colorblock') colorblock;
  @ViewChild(ListDirective, {static: true}) listHost: ListDirective;

  constructor(private componentFactoryResolver: ComponentFactoryResolver) { 
  }

  ngOnInit(): void {
  }

  loadComponent() {
    const componentFactory = this.componentFactoryResolver.resolveComponentFactory(this.listData.content.component);

    const viewContainerRef = this.listHost.viewContainerRef;
    

    const componentRef = viewContainerRef.createComponent(componentFactory);
    (<ListContentModel>componentRef.instance).data = this.listData.content.data;
  }

  unloadComponent() {
    const viewContainerRef = this.listHost.viewContainerRef;
    viewContainerRef.clear();
  }


  @HostListener('mouseenter') onMouseEnter() {
    if(!this.listData.content) return;
    this.loadComponent();
    this.isHovered = true;
    this.updateHeight();
    this.sendHoveredEvent();
  }

  @HostListener('mouseleave') onMouseLeave() {
    if(!this.listData.content) return;
    this.unloadComponent();
    this.isHovered =false;
    this.updateHeight();
    this.sendHoveredEvent();
  }

  sendHoveredEvent() {
    this.hovered.emit(!this.groupHovered);
    console.log('ev')
  }

  updateHeight(delay = 0) {
    const el = this.expandable.nativeElement;

    setTimeout(() => {

      const prevHeight = el.style.height;
      el.style.height = 'auto';
      const newHeight = el.scrollHeight + 'px';
      el.style.height = prevHeight;

      setTimeout(() => {
        el.style.height = newHeight;
      }, 50);
    }, delay);
  }

  changeBarWidth(hovered: boolean) {
    const el = this.colorblock.nativeElement;

    setTimeout(() => {

      const prevWidth = el.style.width;
      el.style.width = 'auto';
      const newWidth = el.scrollWidth + 'px';
      el.style.width = prevWidth;

      setTimeout(() => {
        el.style.height = newWidth;
      }, 50);
    }, 0);

  }
}
