import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { MainComponent } from './main.component';
import { AuthGuard, GuardType } from '../core/guards/auth.guard';

const routes: Routes = [
  {
    path: '',
    component: MainComponent,
    children: [
      {
        path: 'home',
        loadChildren: () =>
          import('./home/home.module').then((m) => m.HomeModule),
      },
      {
        path: 'balance',
        loadChildren: () =>
          import('./balance/balance.module').then((m) => m.BalanceModule),
        canActivate: [AuthGuard],
      },
      {
        path: 'order',
        loadChildren: () =>
          import('./order/order.module').then((m) => m.OrderModule),
        canActivate: [AuthGuard],
      },
      { path: '', pathMatch: 'full', redirectTo: '/home' },
      {
        path: '**',
        loadChildren: () =>
          import('../shared/page-not-found/page-not-found.module').then(
            (m) => m.NotFoundModule
          ),
      },
    ],
  },
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule],
})
export class MainRoutingModule {}
