import { Injectable } from '@angular/core';
import { GenericDataService } from '../../shared/generic-data.service';
import { User } from 'src/app/shared/model';

@Injectable({
  providedIn: 'root',
})
export class UserService extends GenericDataService<User> {}
