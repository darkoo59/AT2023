import { Component } from '@angular/core';
import { Role } from 'src/app/shared/model';
import { AuthService } from '../services/auth.service';

export interface NavRoute {
  path: string;
  title: string;
}

@Component({
  selector: 'app-nav',
  templateUrl: './navigation.component.html',
  styleUrls: ['./navigation.component.scss'],
})
export class NavigationComponent {
  user$ = this.authService.data$;
  routes: NavRoute[] = [
    {
      path: 'home',
      title: 'Home',
    },
  ];

  defaultRoutes: NavRoute[] = [];

  adminRoutes: NavRoute[] = [
    {
      path: 'admin',
      title: 'Admin Panel',
    },
  ];

  constructor(private authService: AuthService) {}

  hasRole(roles: Role[] | undefined, role?: string | null): boolean {
    if (!roles) return false;
    if (roles && !role) return true;
    for (let r of roles) {
      if (r.name == role) return true;
    }
    return false;
  }
}
