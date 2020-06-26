import { Component, OnInit, Input, ViewChildren, QueryList, ElementRef, ViewChild, Output, EventEmitter } from '@angular/core';

@Component({
  selector: 'app-vertical-navbar',
  templateUrl: './vertical-navbar.component.html',
  styleUrls: ['./vertical-navbar.component.scss']
})
export class VerticalNavbarComponent implements OnInit {
  @Input() labels: Array<string>;
  @Output() navigationChange = new EventEmitter<number>();

  @ViewChildren('listItem') headerTexts:QueryList<ElementRef>;
  @ViewChild('underscore') underscore: ElementRef;

  underscoreEl: HTMLElement;
  activeHeadEl: HTMLElement;

  currentGroupIndex = 0;
  distance = 20;
  initalHeight: number;

  constructor() { }

  ngOnInit(): void {
  }

  ngAfterViewInit(): void {
    this.activeHeadEl = this.headerTexts.toArray()[0].nativeElement;
    this.underscoreEl = this.underscore.nativeElement;
    
    this.initalHeight = this.underscoreEl.offsetHeight;
    this.resetUnderscore();
  }

  resetUnderscore(): void {

    const offsetTop = this.activeHeadEl.offsetTop; 
    const offsetLeft = this.activeHeadEl.offsetLeft;
    const underscoreOffset = offsetLeft + 20;

    this.underscoreEl.style.left = "-" + underscoreOffset + "px";
    this.underscoreEl.style.top = offsetTop - this.initalHeight/4 + "px";
    this.underscoreEl.style.height = this.initalHeight + "px";
  }

  expandUnderscore(targetEl: HTMLElement): void {
    const currentTop = this.underscoreEl.offsetTop;
    const currentHeight = this.underscoreEl.offsetHeight;
    const targetTop = targetEl.offsetTop;
    const targetHeight = targetEl.clientHeight;

    let newTop;
    let newHeight;

    if(targetTop < currentTop) {
      newTop = targetTop < currentTop? targetTop : currentTop;
      newHeight = currentTop - targetTop + currentHeight;
    } else {
      newTop = currentTop;
      newHeight = targetTop - currentTop + targetHeight;
    }
    


    this.underscoreEl.style.top = newTop + "px";
    this.underscoreEl.style.height = newHeight +"px";
  }

  onClick(headerEl: HTMLElement, index: number) {
    this.activeHeadEl = headerEl;
    this.resetUnderscore();
    this.currentGroupIndex = index;
    this.navigationChange.emit(index);
  }

  onMouseEnter(headerText: HTMLElement) {
    this.expandUnderscore(headerText);
  }


  onMouseLeave() {
    this.resetUnderscore();
  }


}
