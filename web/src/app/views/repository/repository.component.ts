import { Component, OnInit } from '@angular/core';
import { GraphQlService } from 'src/app/services/graphql.service';
import { map } from 'rxjs/operators';
import { Observable } from 'rxjs';
import { List } from 'src/app/models/list.model';
import { Router } from '@angular/router';

@Component({
  selector: 'app-repository',
  templateUrl: './repository.component.html',
  styleUrls: ['./repository.component.scss']
})
export class RepositoryComponent implements OnInit {
  public listModel: Observable<List>;

  constructor(private graphQlService: GraphQlService, private router: Router) { }

  ngOnInit(): void {
    this.listModel = this.graphQlService.getRepositories()
      .pipe(
        map(repos => ({
          listItems: repos
            .map(repo => ({ label: repo.name, clickable: true })
            ),
          clickHandler: (listItem) => this.router.navigate(['/repository', listItem.label]),

        }
        )
        ))
  }

}
