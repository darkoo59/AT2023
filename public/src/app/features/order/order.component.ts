import { Component, OnInit } from '@angular/core';
import { FormControl, FormGroup, Validators } from '@angular/forms';
import { EMPTY, Subject, catchError, switchMap, tap } from 'rxjs';
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
export class OrderComponent implements OnInit {
  form: FormGroup = new FormGroup({
    item: new FormControl(null, Validators.required),
    quantity: new FormControl(null, [Validators.required, Validators.min(1)]),
  });
  itemList: Item[] = [];
  createOrder$: Subject<CreateOrderDto> = new Subject<CreateOrderDto>().pipe(
    switchMap((data: CreateOrderDto) =>
      this.orderService.createOrder(data).pipe(
        tap(res => {
          this.form.reset()
          this.notificationService.showSuccess(res.status);
        }),
        catchError(() => EMPTY))
    )
  ) as Subject<CreateOrderDto>;

  constructor(private orderService: OrderService, private notificationService: NotificationService) {}

  ngOnInit(): void {
    for (let i = 1; i <= 10; i++) {
      this.itemList.push({
        id: i + '',
        name: 'item_' + i,
      });
    }
  }

  onSubmit() {
    if (this.form.invalid) return;
    const raw = this.form.getRawValue();
    this.createOrder$.next({
      quantity: raw.quantity,
      itemId: raw.item[0],
    });
  }
}
