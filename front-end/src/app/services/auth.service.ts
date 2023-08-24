import { Injectable } from '@angular/core';

@Injectable({
  providedIn: 'root'
})
export class AuthService {
  private accessToken: string | null;
  private JWT_KEY: string = "BINDER_TOKEN";

  constructor() { }

  setToken(token: string){
    this.accessToken = token;
    localStorage.setItem(this.JWT_KEY, token);
  }

  getToken(): string{
    if(!this.accessToken){
      this.accessToken = localStorage.getItem(this.JWT_KEY)!;
    }
    return this.accessToken;
  }

  resetToken(): void{
    this.accessToken = null;
  }

  logout(): void{
    this.resetToken();
    localStorage.removeItem(this.JWT_KEY);
  }
}
