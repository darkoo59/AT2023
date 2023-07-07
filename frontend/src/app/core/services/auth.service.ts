import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable, map, switchMap, take, tap } from 'rxjs';
import { environment } from 'src/environments/environment';
import { GenericDataService } from './generic-data.service';
import { User } from 'src/app/shared/model';

export interface LoginDto {
  email: string;
  password: string;
}

export interface RegisterDto {
  email: string;
  password: string;
}

@Injectable({
  providedIn: 'root',
})
export class AuthService extends GenericDataService<User> {
  constructor(private http: HttpClient) {
    super();
    const json: string | null = localStorage.getItem('user');
    if (json != null) {
      this.setData = JSON.parse(json);
    }
  }

  refreshCookie(): Observable<any> {
    return this.http.get(`${environment.apiUrl}/auth/refresh`, {
      withCredentials: true,
    });
  }

  login(data: LoginDto): Observable<any> {
    return this.addErrorReader(
      this.http
        .post(`${environment.apiUrl}/auth/login`, data, {
          withCredentials: true,
        })
        .pipe(
          switchMap((res: any) =>
            this.data$.pipe(
              take(1),
              map((old: User | null) => {
                const newUser = { ...old, token: res.access_token };
                this.setData = newUser;
                return newUser;
              }),
              tap((user) => localStorage.setItem('user', JSON.stringify(user)))
            )
          )
        )
    );
  }

  registerUser(data: RegisterDto): Observable<any> {
    return this.addErrorReader(
      this.http.post(`${environment.apiUrl}/auth/register`, data)
    );
  }

  logout(): Observable<any> {
    return this.addErrorReader(
      this.http
        .get(`${environment.apiUrl}/auth/logout`, { withCredentials: true })
        .pipe(
          tap(() => {
            this.setData = null;
            localStorage.removeItem('user');
          })
        )
    );
  }
}
