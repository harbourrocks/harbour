import { Component, OnInit, Input } from '@angular/core';
import { GroupModel } from 'src/app/models/group.model';
import { faGitlab } from '@fortawesome/free-brands-svg-icons';

@Component({
  selector: 'app-groups',
  templateUrl: './groups.component.html',
  styleUrls: ['./groups.component.scss']
})
export class GroupsComponent implements OnInit {
  @Input() groupData: GroupModel;

  public isHovered:boolean;

  constructor() {}
  

  ngOnInit(): void {
  }

  onHovered(event){
    this.isHovered = !this.isHovered;
  }

}
