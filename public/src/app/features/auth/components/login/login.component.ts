import { Component } from '@angular/core';
import {
  UntypedFormControl,
  UntypedFormGroup,
  Validators,
} from '@angular/forms';
import { Router } from '@angular/router';
import { catchError, EMPTY, Observable, Subject, switchMap, tap } from 'rxjs';
import { AuthService, LoginDto } from 'src/app/core/services/auth.service';
import { AuthLoadingService } from '../../services/auth-loading.service';
import { NotificationService } from 'src/app/core/services/notification.service';

@Component({
  selector: 'app-login',
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.scss'],
})
export class LoginComponent {
  form: UntypedFormGroup = this.formInstance;
  $error: Observable<string | null> = this.authService.error$;

  $submit: Subject<LoginDto> = new Subject().pipe(
    switchMap((data: LoginDto | unknown) => {
      return this.authService.login(data as LoginDto).pipe(
        catchError((error: any) => {
          if (error.includes("401 Unauthorized")) {
            this.notificationService.showError(`Invalid credentials`);
          }
          this.authLoadingService.setData = false;
          return EMPTY;
        })
      );
    }),
    tap(() => {
      this.authLoadingService.setData = false;
      this.notificationService.showSuccess(`Successfully logged in`);
      this.router.navigate(['/home']);
    })
  ) as Subject<LoginDto>;

  constructor(
    private authService: AuthService,
    private notificationService: NotificationService,
    private router: Router,
    private authLoadingService: AuthLoadingService
  ) {}

  onSubmit() {
    this.authService.clearError();
    if (!this.form.valid) return;
    const data: LoginDto = this.form.getRawValue();
    this.authLoadingService.setData = true;
    this.$submit.next(data);
  }

  get formInstance(): UntypedFormGroup {
    return new UntypedFormGroup({
      email: new UntypedFormControl(null, [
        Validators.required,
        Validators.email,
      ]),
      password: new UntypedFormControl(null, [
        Validators.required,
        Validators.minLength(8),
      ]),
    });
  }
}
