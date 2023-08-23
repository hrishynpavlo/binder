import { Injectable } from '@angular/core';
import { HttpClient, HttpResponse } from '@angular/common/http';
import { Observable } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class BinderApiService {
  private accessToken: string;

  constructor(private http: HttpClient) {}

  getData(id: number): Observable<any> {
    return this.http.get(`http://thebinderapp.com/api/user/${id}`);
  }

  getFeed(): Observable<any>{
    return this.http.get(`http://thebinderapp.com/api/feed`)
  }

  login(data: any): Observable<HttpResponse<any>>{
    return this.http.post(`http://thebinderapp.com/api/login`, data, { observe: 'response' })
  }

  setToken(token: string){
    this.accessToken = token;
  }

  getToken(): string{
    return this.accessToken;
  }
}