import { animate, group, keyframes, query, state, style, transition, trigger } from '@angular/animations';
import { Component, OnInit } from '@angular/core';

@Component({
  selector: 'app-auth',
  templateUrl: './auth.component.html',
  styleUrls: ['./auth.component.scss'],
  animations: [
    trigger('transitionBlock', [
      state('register', style({ transform: 'translateX(0)' })),
      state('login',    style({ transform: 'translateX(700px)' })),
      transition('register => login', [
        group([
          query('.register-text', [
            style({ opacity: 0, transform: 'translateX(100px)' }),
            animate(500, style({ opacity: 1, transform: 'none' }))
          ]),
          query('.login-text', [
            style({opacity: 1, transform: 'none'}),
            animate(500, style({ opacity: 0, transform: 'translateX(-100px)' })),
          ]),
          animate(1000, keyframes([
            style({ transform: 'scaleX(1.2) translateX(500px)',   offset: 0.35 }),
            style({ transform: 'scaleX(1.05) translateX(650px)',  offset: 0.6 }),
            style({ transform: 'scaleX(1) translateX(700px)',     offset: 1 })
          ]))
        ])
      ]),
      transition('login => register', [
        group([
          query('.login-text', [
            style({ opacity: 0, transform: 'translateX(-100px)' }),
            animate(600, style({ opacity: 1, transform: 'none' }))
          ]),
          query('.register-text', [
            style({ opacity: 1, transform: 'none' }),
            animate(600, style({ opacity: 0, transform: 'translateX(100px)' })),
          ]),
          animate(1000, keyframes([
            style({ transform: 'scaleX(1.25) translateX(200px)',  offset: 0.35 }),
            style({ transform: 'scaleX(1.05) translateX(17px)',   offset: 0.6 }),
            style({ transform: 'scaleX(1) translateX(0px)',       offset: 1 })
          ]))
        ])
      ])
    ]),
    trigger('loginBlock', [
      state('login',    style({ transform: 'translateX(0)',       opacity:1 })),
      state('register', style({ transform: 'translateX(-2000px)', opacity:0 })),
      transition('login => register', animate(400, 
        style({ transform: 'translateX(100px)', opacity: 0 })
      )),
      transition('register => login', animate(800, keyframes([
        style({ transform: 'translateX(100px)', offset: 0 }),
        style({ opacity: 0, offset: 0.3 }),
        style({ transform: 'translateX(0px)', offset:1, opacity: 1 })
      ])))
    ]),
    trigger('registerBlock', [
      state('register', style({ transform: 'translateX(0)',       opacity:1 })),
      state('login',    style({ transform: 'translateX(-2000px)',  opacity:0 })),
      transition('register => login', animate(400, 
        style({ transform: 'translateX(-100px)', opacity:0 })
      )),
      transition('login => register', animate(800, keyframes([
        style({ transform: 'translateX(-100px)', offset: 0 }),
        style({ opacity: 0, offset: 0.3 }),
        style({ transform: 'translateX(0px)', offset:1, opacity: 1 })
      ])))
    ]),
    trigger('registerMessage', [
      state('register', style({ transform: 'translateX(420px)', opacity:0 })),
      state('login', style({ transform: 'translateX(0px)', opacity:1 })),
      transition('login <=> register', animate(600))
    ]),
    trigger('loginMessage', [
      state('login', style({ transform: 'translateX(-420px)', opacity:0 })),
      state('register', style({ transform: 'translateX(0px)', opacity:1 })),
      transition('login <=> register', animate(600))
    ])
  ]
})
export class AuthComponent implements OnInit {
  state = 'login';

  constructor() { }
  ngOnInit(): void {}

  switchState(){
    if(this.state == 'login') this.state = 'register';
    else this.state = 'login';
  }
}
