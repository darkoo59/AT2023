import {
  HttpEvent,
  HttpHandler,
  HttpInterceptor,
  HttpRequest,
} from '@angular/common/http';
import { Injectable } from '@angular/core';
import { catchError, Observable, switchMap, take, throwError } from 'rxjs';
import { AuthService } from 'src/app/core/services/auth.service';
import { User } from '../model';

@Injectable()
export class AuthInterceptor implements HttpInterceptor {
  constructor(private authService: AuthService) {}

  intercept(
    req: HttpRequest<any>,
    next: HttpHandler
  ): Observable<HttpEvent<any>> {
    if (req.url.indexOf('/auth/refresh') > -1) {
      return this.authService.data$.pipe(
        switchMap((user: User | null) => {
          take(1)
          req = this.addToken(req, user?.token);
          return next.handle(req);
        })
      );
    }

    return this.authService.data$.pipe(
      take(1),
      switchMap((user: User | null) => {
        req = this.addToken(req, user?.token);
        return next.handle(req).pipe(
          catchError((error) => {
            if (error.status == 401) {
              return this.handle401Error(req, next);
            } else {
              return throwError(() => error);
            }
          })
        );
      })
    );
  }

  private handle401Error(req: HttpRequest<any>, next: HttpHandler) {
    return this.authService.refreshCookie().pipe(
      take(1),
      switchMap((res: any) =>
        this.authService.data$.pipe(
          take(1),
          switchMap((user: User | null) => {
            this.authService.setData = { ...user, token: res.access_token };
            return next.handle(req);
          })
        )
      ),
      catchError(() => this.authService.logout())
    );
  }

  private addToken(
    req: HttpRequest<any>,
    token: string | undefined
  ): HttpRequest<any> {
    if (token) {
      return req.clone({
        setHeaders: {
          Authorization: `Bearer ${token}`,
        },
      });
    }
    return req;
  }
}
