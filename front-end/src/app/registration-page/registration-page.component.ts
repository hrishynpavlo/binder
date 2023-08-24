import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { AuthService } from '../services/auth.service';
import { BinderApiService } from '../services/binder-api-service';
import { FormControl, FormGroup, Validators } from '@angular/forms';

@Component({
  selector: 'app-registration-page',
  templateUrl: './registration-page.component.html',
  styleUrls: ['./registration-page.component.css']
})
export class RegistrationPageComponent implements OnInit {
  constructor(private router: Router, private authService: AuthService, private api: BinderApiService) {}

  registrationForm: FormGroup;
  latitude: number;
  longitude: number;

  ngOnInit(): void {
    this.registrationForm = new FormGroup({
      firstName: new FormControl('', [Validators.required, Validators.minLength(2)]),
      lastName: new FormControl('', [Validators.required, Validators.minLength(2)]),
      password: new FormControl('', [Validators.required, Validators.minLength(6)]),
      email: new FormControl('', [Validators.required, Validators.email]),
      displayName: new FormControl(''),
      country: new FormControl('', [Validators.required]),
      dateOfBirth: new FormControl('', [Validators.required])
    });

    navigator.geolocation.getCurrentPosition(geo => {
      console.log(geo);
      this.latitude = geo.coords.latitude,
      this.longitude = geo.coords.longitude
    }, _ => {});
  }

  onSubmit(): void {
    if(this.registrationForm.valid) {
      const firstName = this.registrationForm.controls['firstName'].value;
      const lastName = this.registrationForm.controls['lastName'].value;
      const displayName = this.registrationForm.controls['displayName'].value;

      const user = {
        FirstName: firstName,
        LastName: lastName,
        Password: this.registrationForm.controls['password'].value,
        Email: this.registrationForm.controls['email'].value,
        DisplayName: displayName ? displayName : `${firstName} ${lastName}`,
        Country: this.registrationForm.controls['country'].value,
        DateOfBirth: this.registrationForm.controls['dateOfBirth'].value,
        Latitude: this.latitude,
        Longitude: this.longitude
      };

      this.api.createUser(user).subscribe(response => {
        this.authService.setToken(response['access_token']);
        this.router.navigate(['/me']);
      });
    }
  }
}
