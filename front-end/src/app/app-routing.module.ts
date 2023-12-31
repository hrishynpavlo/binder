import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { UserFeedComponent } from './user-feed/user-feed.component';
import { LoginPageComponent } from './login-page/login-page.component';
import { authGuard } from './services/auth-guard.guard';
import { RegistrationPageComponent } from './registration-page/registration-page.component';

const routes: Routes = [
  { path: 'login', component: LoginPageComponent },
  { path: 'register', component: RegistrationPageComponent },
  { path: 'me', component: UserFeedComponent, canActivate: [authGuard] },
  { path: '**', redirectTo: "/me"}
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
