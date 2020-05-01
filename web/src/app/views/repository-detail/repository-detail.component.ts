import { Component, OnInit, ViewChild, ViewChildren, QueryList, AfterViewInit, ElementRef } from '@angular/core';
import { GroupModel } from 'src/app/models/group.model';
import { ListDetailsModel } from 'src/app/models/list-detail.model';
import { GroupsComponent } from 'src/app/components/groups/groups.component';

@Component({
  selector: 'app-repository-detail',
  templateUrl: './repository-detail.component.html',
  styleUrls: ['./repository-detail.component.scss']
})
export class RepositoryDetailComponent implements OnInit, AfterViewInit {
  @ViewChildren('headerText') headerTexts:QueryList<ElementRef>;
  @ViewChild('underscore') underscore: ElementRef;
  @ViewChild('groupWrapper') groupWrapper: ElementRef;
  
  groupEl: HTMLElement;
  underscoreEl: HTMLElement;
  activeHeadEl: HTMLElement;
  currentGroupIndex = 0;
  heightOffset = 7;

  //INPUT IN FUTURE
  public repoDetails: Array<ListDetailsModel>;

  constructor() { 
    const groupProperty:GroupModel = {
      smallColorbox:true,
      listItems:[
        {text: "2.0.0-release1"},
        {text: "1.3.1"},
        {text: "1.2.4"},
        {text: "1.2.3"},
        {text: "1.2.2"},
        {text: "1.2.1"},
        {text: "1.1.4"},
        {text: "1.1.3"},
        {text: "1.1.2"},
        {text: "1.1.1"},
        {text: "1.0.1"}
      ]
    }
    const groupImage:GroupModel = {
      listItems:[
        {text: "2.0.0-release1"},
        {text: "1.3.1"},
      ]
    }
    const groupBuilds:GroupModel = {
      listItems:[
        {text: "2.0.0-release1"},
      ]
    }

    this.repoDetails = [
      {title: "Properties", group: groupProperty},
      {title: "Images", group: groupImage},
      {title: "Builds", group: groupBuilds}
    ]
  }

  ngOnInit(): void {
  }

  ngAfterViewInit(): void {
    this.groupEl = this.groupWrapper.nativeElement;
    this.activeHeadEl = this.headerTexts.toArray()[0].nativeElement;
    this.underscoreEl = this.underscore.nativeElement;
    this.resetUnderscore();

    const offsetTop = this.groupEl.offsetTop;
    const windowHeight = window.innerHeight;
    const calculatedHeight = windowHeight - offsetTop;
    console.log(calculatedHeight)

    this.groupEl.style.height = calculatedHeight + "px";
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
  

  getGroupData() {
    return this.repoDetails[this.currentGroupIndex].group;
  }

}
