import { Component, OnInit } from '@angular/core'
import { CommonModule } from '@angular/common'
import { UserStoreService } from '../../../../state/user.service'
import { NzFormModule } from 'ng-zorro-antd/form'
import { NzInputModule } from 'ng-zorro-antd/input'
import { FormsModule } from '@angular/forms'
import { NzModalModule } from 'ng-zorro-antd/modal'
import { NzMessageService } from 'ng-zorro-antd/message'
import { postApiWithAuth } from '../../../../tool/httpRequest-auth'
import { NzIconModule } from 'ng-zorro-antd/icon'
import { NzUploadFile, NzUploadModule } from 'ng-zorro-antd/upload'
import { Observable, Observer } from 'rxjs'
import { getBase64, uploadImgToBase64 } from '../../../../tool/imageUpload'
import { NzTableModule } from 'ng-zorro-antd/table'
import moment from 'moment'
@Component({
    // selector: 'app-footer',
    imports: [CommonModule, NzFormModule, NzInputModule, FormsModule, NzModalModule, NzIconModule, NzUploadModule, NzTableModule],
    templateUrl: './user-info.component.html',
    styleUrl: './user-info.component.css',
})
export class UserInfoComponent implements OnInit {
    constructor(private userService: UserStoreService, private message: NzMessageService) {}

    userInfo: any = { 
        _id: '',
        username: '',
        accessToken: '',
        avatarBase64: '',
        deptId: 0,
        roleIds: [],
        roleLists: [],
        email: '', 
        loginRecords: [],
        departments: []
    }

    resetPwForm: any = {
        newPassword: '',
        againNewPassword: ''
    }

    passwordVisible = false

    resetPwDialog: boolean = false

    ngOnInit() {
        this.userService.user$.subscribe(user => {
            console.log(user)
            this.userInfo = user
            
        })
    } 

    openResetPWDialog() {
        this.resetPwDialog = true
    }

    closeResetPWDialog() {
        this.resetPwDialog = false
    }

    async submitForm() {
        if (this.resetPwForm.newPassword === this.resetPwForm.againNewPassword) {
            const res = await postApiWithAuth('/sys/user/user-self/update-password', {   username: this.userInfo.username, newPassword: this.resetPwForm.newPassword })

            if (res.ModifiedCount === 1) {
                this.message.info('Reset password successfully!')
                this.closeResetPWDialog()
            } else {
                this.message.error(res.msg)
            }
            
        } else {
            this.message.error('Two passwords are inconsistent!')
        }
    }

    // upload

    loading = false
    avatarUrl?: string = ''
    previewImage: string = ''

    beforeUpload = (file: NzUploadFile, _fileList: NzUploadFile[]): Observable<boolean> =>
        new Observable((observer: Observer<boolean>) => {
          const isJpgOrPng = file.type === 'image/jpeg' || file.type === 'image/png'
          if (!isJpgOrPng) {
            this.message.error('You can only upload JPG or PNGfile!')
            observer.complete()
            return
          }
          const isLt2M = file.size! / 1024 / 1024 < 3
          if (!isLt2M) {
            this.message.error('Image must smaller than 3MB!')
            observer.complete()
            return
          }

          observer.next(isJpgOrPng && isLt2M)
          observer.complete()
    })
    

    async handleChange(info: { file: any }) {
        const response: any = await uploadImgToBase64(info.file.originFileObj)
        this.avatarUrl = response.data
    }


    async handleUpdateAvatar() {
        const response = await postApiWithAuth('/sys/user/user-self/update-avatar', { photo: this.avatarUrl, username: this.userInfo.username })

        if (response.result.ModifiedCount === 1) {
            this.message.success('Update success!')
            this.userService.loadUserInfo()
            this.avatarUrl = ''
        } else {
            this.message.error('Ooops! Try again!')
        }
    }

    dateFormat(data: string) {
        return data ? moment(new Date(data)).format('DD-MM-YYYY HH:MM') : null
    }
}