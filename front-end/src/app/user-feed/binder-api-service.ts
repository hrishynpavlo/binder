import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class BinderApiService {
  constructor(private http: HttpClient) {}

  getData(id: number): Observable<any> {
    return this.http.get(`http://localhost:8080/api/user/${id}`);
  }
}