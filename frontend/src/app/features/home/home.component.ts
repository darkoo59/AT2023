import { Component } from "@angular/core";
import { EMPTY, Subject, catchError, switchMap } from "rxjs";
import { OrderService } from "src/app/core/services/order.service";

@Component({
  templateUrl: 'home.component.html'
})
export class HomeComponent {
  createOrder$: Subject<any> = new Subject<any>().pipe(
    switchMap(() => this.orderService.createOrder().pipe(
      catchError(() => EMPTY)
    ))
  ) as Subject<any>;

  constructor(private orderService: OrderService) {}
}
