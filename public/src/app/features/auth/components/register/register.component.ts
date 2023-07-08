import { Component } from '@angular/core';
import {
  UntypedFormControl,
  UntypedFormGroup,
  Validators,
} from '@angular/forms';
import { Router } from '@angular/router';
import { AuthService, RegisterDto } from 'src/app/core/services/auth.service';
import { AuthLoadingService } from '../../services/auth-loading.service';
import { Subject, switchMap, catchError, EMPTY, tap, Observable } from 'rxjs';
import { NotificationService } from 'src/app/core/services/notification.service';

@Component({
  selector: 'app-register',
  templateUrl: './register.component.html',
  styleUrls: ['./register.component.scss'],
})
export class RegisterComponent {
  form: UntypedFormGroup = this.formInstance;
  $error: Observable<string | null> = this.authService.error$;

  $submit: Subject<RegisterDto> = new Subject().pipe(
    switchMap((data: RegisterDto | unknown) => {
      return this.authService.registerUser(data as RegisterDto).pipe(
        catchError((error: any) => {
          this.authLoadingService.setData = false;
          return EMPTY;
        })
      );
    }),
    tap(() => {
      this.authLoadingService.setData = false;
      this.notificationService.showSuccess(`Successfully registered`);
      this.router.navigate(['/home']);
    })
  ) as Subject<RegisterDto>;

  constructor(
    private authService: AuthService,
    private notificationService: NotificationService,
    private router: Router,
    private authLoadingService: AuthLoadingService
  ) {}

  onSubmit() {
    this.authService.clearError();
    if (!this.form.valid) return;
    const data: RegisterDto = this.form.getRawValue();
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
