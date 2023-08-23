import { inject } from '@angular/core';
import { CanActivateFn, Router } from '@angular/router';
import { BinderApiService } from './binder-api-service';



export const authGuard: CanActivateFn = (route, state) => {
  const router = inject(Router);
  const api = inject(BinderApiService);

  if(!api.getToken()){
    router.navigate(['/login']);
    return false;
  }
  return true;
};
