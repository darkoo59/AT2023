import { Component } from "@angular/core";
import { FormControl, FormGroup, Validators } from "@angular/forms";
import { EMPTY, Observable, Subject, catchError, switchMap, tap } from "rxjs";
import { BalanceService } from "src/app/core/services/balance.service";
import { NotificationService } from "src/app/core/services/notification.service";

@Component({
  selector: 'app-balance',
  templateUrl: './balance.component.html',
  styleUrls: ['./balance.component.scss'],
})
export class BalanceComponent {
  form: FormGroup = new FormGroup({
    quantity: new FormControl(null, [Validators.required, Validators.min(1)]),
  });
  updateBalance$: Subject<number> = new Subject<number>().pipe(
    switchMap((data: number) => this.balanceService.patchBalance({ balance: data }).pipe(
      switchMap(() => {
        this.form.reset()
        return this.balanceService.fetchBalance()
      }),
      tap(res =>
        this.notificationService.showSuccess(res.status)
      ),
      catchError(() => EMPTY),
    ))
  ) as Subject<number>;
  balance$: Observable<number | null> = this.balanceService.data$;
  fetchBalance$: Observable<any> = this.balanceService.fetchBalance();

  constructor(private balanceService: BalanceService, private notificationService: NotificationService) { }

  onSubmit() {
    if (this.form.invalid) return;
    const raw = this.form.getRawValue();
    this.updateBalance$.next(raw.quantity)
  }
}