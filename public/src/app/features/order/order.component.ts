import { Component, OnInit } from '@angular/core';
import { FormControl, FormGroup, Validators } from '@angular/forms';
import { EMPTY, Observable, Subject, catchError, switchMap, tap } from 'rxjs';
import { NotificationService } from 'src/app/core/services/notification.service';
import {
  CreateOrderDto,
  OrderService,
} from 'src/app/core/services/order.service';
import { Item } from 'src/app/shared/model';

@Component({
  selector: 'app-order',
  templateUrl: './order.component.html',
  styleUrls: ['./order.component.scss'],
})
export class OrderComponent {
  form: FormGroup = new FormGroup({
    item: new FormControl(null, Validators.required),
    quantity: new FormControl(null, [Validators.required, Validators.min(1)]),
  });
  getItems$: Observable<any> = this.orderService.getItems();
  itemList$: Observable<Item[] | null> = this.orderService.data$;
  createOrder$: Subject<CreateOrderDto> = new Subject<CreateOrderDto>().pipe(
    switchMap((data: CreateOrderDto) =>
      this.itemList$.pipe(
        switchMap((items: Item[] | null) => {
          data = { ...data, price: items?.find(item => item.id === data.itemId)?.price ?? 0 }
          return this.orderService.createOrder(data).pipe(
            tap(res => {
              this.form.reset()
              this.notificationService.showSuccess(res.status);
            }),
            catchError(() => EMPTY)
          )
        })
      )
    )
  ) as Subject<CreateOrderDto>;
  
  constructor(private orderService: OrderService, private notificationService: NotificationService) { }

  onSubmit() {
    if (this.form.invalid) return;
    const raw = this.form.getRawValue();
    const itemId = raw.item[0];
    this.createOrder$.next({
      quantity: raw.quantity,
      itemId: itemId
    });
  }
}
