import { Injectable } from '@angular/core';
import { HttpClient, HttpResponse } from '@angular/common/http';
import { Observable } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class BinderApiService {

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

  createUser(user: any): Observable<any>{
    return this.http.post(`http://thebinderapp.com/api/user`, user);
  }
}