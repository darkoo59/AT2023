import { Injectable } from '@angular/core';
import { GenericDataService } from 'src/app/shared/generic-data.service';

export interface LoginDto {
  email: string;
  password: string;
}

@Injectable()
export class AuthLoadingService extends GenericDataService<boolean> {}
