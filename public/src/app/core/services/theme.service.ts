import { Injectable } from '@angular/core';
import { GenericDataService } from '../../shared/generic-data.service';

@Injectable({
  providedIn: 'root',
})
export class ThemeService
  extends GenericDataService<'dark' | 'light'>
{
  constructor() {
    super()
    this.setData = localStorage.getItem('theme') === 'dark' ? 'dark' : 'light';
  }
}
