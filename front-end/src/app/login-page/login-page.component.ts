import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { FormGroup, FormControl, Validators } from '@angular/forms';
import { BinderApiService } from '../services/binder-api-service';
import { AuthService } from '../services/auth.service';

@Component({
  selector: 'app-login-page',
  templateUrl: './login-page.component.html',
  styleUrls: ['./login-page.component.css']
})
export class LoginPageComponent implements OnInit {
  loginForm: FormGroup;

  constructor(private api: BinderApiService, private router: Router, private authService: AuthService){}

  ngOnInit(): void {
    this.loginForm = new FormGroup({
      email: new FormControl('', [Validators.required, Validators.email]),
      password: new FormControl('', [Validators.required, Validators.minLength(6)])
    })
  }

  onSubmit() {
    if(this.loginForm.valid){
      const data = {
        email: this.loginForm.controls['email'].value,
        password: this.loginForm.controls['password'].value
      }

      this.api.login(data).subscribe(loginResponse => {
        this.authService.setToken(loginResponse.body['access_token']);
        this.router.navigate(['/me']);
      });
    }
  }
}
