import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { Item } from 'src/app/shared/model';
import { environment } from 'src/environments/environment';

export interface CreateOrderDto {
  quantity: number;
  itemId: string;
  price: number;
}

@Injectable({
  providedIn: 'root',
})
export class OrderService {
  constructor(private http: HttpClient) {}

  createOrder(data: CreateOrderDto): Observable<any> {
    return this.http.post(`${environment.apiUrl}/customer/order`, data)
  }

  getItems() : Observable<Item[]> {
    return this.http.get<Item[]>(`${environment.apiUrl}/customer/items`);
  }
}
