import { Component, OnInit } from '@angular/core'
import { CommonModule } from '@angular/common'
import { Router } from '@angular/router'
import { NglModule } from 'ng-lightning'
import { FormsModule } from '@angular/forms'
import { UserStoreService } from '../../../../state/user.service'
import { HttpService } from '../../../../tool/HttpService'
import { callPostApi } from '../../../../tool/call-http'
import { getApiFromOutside, postApi } from '../../../../tool/httpRequest-public'
import { NzMessageService } from 'ng-zorro-antd/message'
import { postApiWithAuth } from '../../../../tool/httpRequest-auth'


@Component({
    // selector: 'app-footer',
    standalone: true,
    imports: [CommonModule, NglModule, FormsModule],
    templateUrl: './login.component.html',
    styleUrl: './login.component.scss',
})
export class LoginComponent implements OnInit{
    constructor(
        private userStoreService: UserStoreService, 
        private router: Router,
        private httpService: HttpService,
        private message: NzMessageService
    ) {}
    ngOnInit(): void {
        this.getIpAdress()
    }

    loginForm: any = {
        username: '',
        password: '',
        ipAddress: ''
    }

    standalone = {
        standalone: true
    }

    async login() {

        const res = await postApiWithAuth('/auth/login', this.loginForm)

        if (!res.msg) {
            this.userStoreService.setAccessToken(res.accessToken)
            const checkpoint = await this.userStoreService.isAuthenticated()
            if (checkpoint) {
                this.router.navigate(['/'])
            }
        } else {
            this.message.error(res.msg)
        }
        
       
    }

    cleanForm() {
        this.loginForm = {
            username: '',
            password: '',
            ipAddress: ''
        }
    }

    async getIpAdress() {
        const res = await getApiFromOutside('https://geolocation-db.com/json/')
        this.loginForm.ipAddress = res.IPv4
    }
}