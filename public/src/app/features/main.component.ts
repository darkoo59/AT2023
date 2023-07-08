import { Component } from '@angular/core';
import { AuthService } from '../core/services/auth.service';
import { ThemeService } from '../core/services/theme.service';
import { Router } from '@angular/router';
import { EMPTY, Subject, catchError, switchMap } from 'rxjs';
import { UserService } from '../core/services/user.service';

@Component({
  selector: 'app-main',
  templateUrl: './main.component.html',
  styleUrls: ['./main.component.scss'],
})
export class MainComponent {
  theme$ = this.themeService.data$;
  user$ = this.userService.data$;
  isLogged$ = this.authService.isLogged$;

  logout$: Subject<any> = new Subject().pipe(
    switchMap(() => this.authService.logout().pipe(
      catchError(() => EMPTY)
    ))
  ) as Subject<any>;

  constructor(
    private themeService: ThemeService,
    private authService: AuthService,
    private userService: UserService,
    private router: Router
  ) {}

  toggleTheme(theme: 'dark' | 'light'): void {
    this.themeService.setData = theme;
    localStorage.setItem('theme', theme + '');
  }

  logout(): void {
    this.logout$.next(0)
  }

  login(): void {
    this.router.navigate(['/auth'])
  }
}
