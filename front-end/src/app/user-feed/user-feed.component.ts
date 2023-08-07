import { Component, OnInit } from '@angular/core';
import { BinderApiService } from './binder-api-service';

@Component({
  selector: 'app-user-feed',
  templateUrl: './user-feed.component.html',
  styleUrls: ['./user-feed.component.css']
})
export class UserFeedComponent implements OnInit {
  userId: number;
  user: any;
  feed: any[];

  constructor(private api: BinderApiService){}

  ngOnInit(): void {
    this.userId = 52;
    this.api.getFeed(this.userId).subscribe(response => {
      this.feed = response;
    });
    this.api.getData(this.userId).subscribe(response => {
      this.user = response;
    })
  }

  like(who: number, whom: number){
    console.log(who, whom)
  }
}
