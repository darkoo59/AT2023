import { Component } from '@angular/core';
import { AuthService } from '../services/auth.service';

export interface NavRoute {
  path: string;
  title: string;
  protected?: boolean;
}

@Component({
  selector: 'app-nav',
  templateUrl: './navigation.component.html',
  styleUrls: ['./navigation.component.scss'],
})
export class NavigationComponent {
  isLogged$ = this.authService.isLogged$;
  routes: NavRoute[] = [
    {
      path: 'home',
      title: 'Home',
    },
    {
      path: 'balance',
      title: 'Balance',
      protected: true
    },
    {
      path: 'order',
      title: 'Order',
      protected: true
    }
  ];

  constructor(private authService: AuthService) {}
}
