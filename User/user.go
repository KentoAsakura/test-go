package User

import (
    "image/png"
    "os"
    "github.com/boombuler/barcode"
    "github.com/boombuler/barcode/qr"
    "fmt"
    "gorm.io/gorm"
    "gorm.io/driver/sqlite"
)

var DB *gorm.DB

func Init() {
    var err error
    DB, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
    if err != nil {
        panic("failed to connect database")
    }
    // データベースのマイグレーション
    DB.AutoMigrate(&User{})
    DB.AutoMigrate(&UserInfoMation{})
}

type UserInfoMation struct {
    gorm.Model
    UserID        uint   // ユーザーテーブルの主キー
    PhoneNumber   string
    QRCODE_Number string
    Attend        bool
}

type User struct {
    gorm.Model
    UserName     string
    PhoneNumber  string
    UserInfo     UserInfoMation
}

func CreateQRCode(phoneNumber string) (string) {
    qrCode, _ := qr.Encode(phoneNumber, qr.M, qr.Auto)
    qrCode, _ = barcode.Scale(qrCode, 200, 200)
    fileName := fmt.Sprintf("QRCode/%s_qrcode.png", phoneNumber)
    file, _ := os.Create(fileName)
    defer file.Close()

    png.Encode(file, qrCode)
    return phoneNumber + "_qrcode.png"
}

func CreateUser(username, phoneNumber string) error {
    user := &User{UserName: username, PhoneNumber: phoneNumber}
    var userInfo UserInfoMation
    userInfo.PhoneNumber = user.PhoneNumber
    userInfo.QRCODE_Number = CreateQRCode(user.PhoneNumber)
    user.UserInfo = userInfo
    if err := DB.Create(user).Error; err != nil {
        return err
    }
    return nil
}
