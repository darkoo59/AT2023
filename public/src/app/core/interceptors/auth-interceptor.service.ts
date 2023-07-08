import {
  HttpEvent,
  HttpHandler,
  HttpInterceptor,
  HttpRequest,
} from '@angular/common/http';
import { Injectable } from '@angular/core';
import { catchError, Observable, switchMap, throwError } from 'rxjs';
import { AuthService } from 'src/app/core/services/auth.service';

@Injectable()
export class AuthInterceptor implements HttpInterceptor {
  constructor(private authService: AuthService) {}

  intercept(
    req: HttpRequest<any>,
    next: HttpHandler
  ): Observable<HttpEvent<any>> {
    req = this.addToken(req, this.authService.token);
    return req.url.indexOf('/auth/refresh') > -1
      ? next.handle(req)
      : next.handle(req).pipe(
          catchError((error) => {
            if (error.status == 401 && this.authService.token !== null) {
              return this.handle401Error(req, next);
            } else {
              return throwError(() => error);
            }
          })
        );
  }

  private handle401Error(req: HttpRequest<any>, next: HttpHandler) {
    return this.authService.refreshCookie().pipe(
      switchMap((res: any) => {
        this.authService.setData = res.access_token;
        req = this.addToken(req, this.authService.token);
        return next.handle(req);
      }),
      catchError(() => this.authService.logout())
    );
  }

  private addToken(
    req: HttpRequest<any>,
    token: string | null
  ): HttpRequest<any> {
    return !token
      ? req
      : req.clone({
          setHeaders: {
            Authorization: `Bearer ${token}`,
          },
        });
  }
}
