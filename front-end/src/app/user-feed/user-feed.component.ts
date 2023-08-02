import { Component, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { BinderApiService } from './binder-api-service';

@Component({
  selector: 'app-user-feed',
  templateUrl: './user-feed.component.html',
  styleUrls: ['./user-feed.component.css']
})
export class UserFeedComponent implements OnInit {
  userId: number;
  user: any;

  constructor(private route: ActivatedRoute, private api: BinderApiService){}

  ngOnInit(): void {
    this.userId = +this.route.snapshot.paramMap.get('id')!;
    this.api.getData(this.userId).subscribe(response => {
      this.user = response;
    });
  }
}
