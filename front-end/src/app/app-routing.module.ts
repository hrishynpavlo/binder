import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { UserFeedComponent } from './user-feed/user-feed.component';

const routes: Routes = [
  { path: 'me', component: UserFeedComponent },
  { path: '**', redirectTo: "/me"}
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
