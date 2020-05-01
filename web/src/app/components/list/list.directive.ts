import { Directive, ViewContainerRef } from '@angular/core';

@Directive({
  selector: '[list-host]',
})
export class ListDirective {
  constructor(public viewContainerRef: ViewContainerRef) { }
}


