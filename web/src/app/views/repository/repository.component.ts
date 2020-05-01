import { Component, OnInit } from '@angular/core';
import { GroupModel } from 'src/app/models/group.model';
import { faShip } from '@fortawesome/free-solid-svg-icons';

@Component({
  selector: 'app-repository',
  templateUrl: './repository.component.html',
  styleUrls: ['./repository.component.scss']
})
export class RepositoryComponent implements OnInit {
  public groupArr: Array<GroupModel>;

  constructor() {
    this.groupArr = [
      {
        smallColorbox: true,
        title: "Harbour",
        icon: faShip,
        collapsable: true,
        listItems: [
          {
            text: "harbour-scm",
          },
          {
            text: "harbour-build",
          },
          {
            text: "harbour-iam",
          },
        ]
      },
      {
        smallColorbox: true,
        title: "Sample Project",
        icon: faShip,
        collapsable: true,
        listItems: [
          {
            text: "sample-frontend",
          },
          {
            text: "sample-backend",
          }
        ]
      },

    ]
  }

  ngOnInit(): void {
  }

}
