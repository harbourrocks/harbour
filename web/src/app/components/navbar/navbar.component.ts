import { Component, OnInit, ViewChildren, QueryList, ElementRef, ViewChild } from '@angular/core';
import { faTachometerAlt, faBoxOpen, faHatWizard, faCog } from '@fortawesome/free-solid-svg-icons';

@Component({
  selector: 'app-navbar',
  templateUrl: './navbar.component.html',
  styleUrls: ['./navbar.component.scss']
})
export class NavbarComponent implements OnInit {
  @ViewChildren('listItem') headerTexts:QueryList<ElementRef>;
  @ViewChild('underscore') underscore: ElementRef;

  underscoreEl: HTMLElement;
  activeHeadEl: HTMLElement;

  currentGroupIndex = 0;
  heightOffset = 7;

  public navbarBuild = [
    {
      routerLink: "/dashboard",
      text: "Dashboard",
    },
    {
      routerLink: "/repository",
      text: "Repository",
    },
    {
      routerLink: "/settings",
      text: "Settings",
    },

  ]

  constructor() { }

  ngOnInit(): void {
  }

  ngAfterViewInit(): void {
    this.activeHeadEl = this.headerTexts.toArray()[0].nativeElement;
    this.underscoreEl = this.underscore.nativeElement;
    this.resetUnderscore();
  }

  resetUnderscore(): void {
    const offsetTop = this.activeHeadEl.offsetTop;
    const offsetLeft = this.activeHeadEl.offsetLeft;
    const height = this.activeHeadEl.clientHeight;
    const width = this.activeHeadEl.clientWidth;
    const underscoreOffset = offsetTop + height + this.heightOffset;

    this.underscoreEl.style.left = offsetLeft + "px";
    this.underscoreEl.style.top = underscoreOffset + "px";
    this.underscoreEl.style.width = width + "px";
  }

  expandUnderscore(targetEl: HTMLElement): void {
    const currentLeft = this.underscoreEl.offsetLeft;
    const currentWidth = this.underscoreEl.offsetWidth;
    const targetLeft = targetEl.offsetLeft;
    const targetWidth = targetEl.clientWidth;

    let newLeft;
    let newWidth;

    if(targetLeft < currentLeft) {
      newLeft = targetLeft < currentLeft? targetLeft : currentLeft;
      newWidth = currentLeft - targetLeft + currentWidth;
    } else {
      newLeft = currentLeft;
      newWidth = targetLeft - currentLeft + targetWidth;
    }
    


    this.underscoreEl.style.left = newLeft + "px";
    this.underscoreEl.style.width = newWidth +"px";
  }

  onClick(headerEl: HTMLElement, index: number) {
    this.activeHeadEl = headerEl;
    this.resetUnderscore();
    this.currentGroupIndex = index;
  }

  onMouseEnter(headerText: HTMLElement) {
    this.expandUnderscore(headerText);
  }


  onMouseLeave() {
    this.resetUnderscore();
  }

}
