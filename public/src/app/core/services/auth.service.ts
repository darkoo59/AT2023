import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable, map, tap } from 'rxjs';
import { environment } from 'src/environments/environment';
import { GenericDataService } from '../../shared/generic-data.service';
import { Router } from '@angular/router';

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
export class AuthService extends GenericDataService<string> {
  token: string | null = null;
  isLogged$: Observable<boolean> = this.data$.pipe(map((token) => !!token));
  
  constructor(private http: HttpClient, private router: Router) {
    super();
    const token: string | null = localStorage.getItem('access_token');
    this.setData = token;
  }

  override set setData(data: string | null) {
    super.setData = data;
    this.token = data;
    data !== null
      ? localStorage.setItem('access_token', data)
      : localStorage.removeItem('access_token');
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
        .pipe(tap((res: any) => (this.setData = res.access_token ?? null)))
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
        .pipe(tap(() => {
          this.setData = null
          this.router.navigate(['/home'])
        }))
    );
  }
}
