package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/urfave/cli/v3"
	"github.com/v03413/bepusdt/app/model"
	"github.com/v03413/bepusdt/app/task"
	"github.com/v03413/bepusdt/app/utils"
	"golang.org/x/crypto/bcrypt"
)

var Reset = &cli.Command{
	Name:  "reset",
	Usage: "忘记密码时，此命令可重置账号密码登录入口",
	Flags: []cli.Flag{SQLiteLogFlag},
	Before: func(ctx context.Context, c *cli.Command) (context.Context, error) {
		if err := model.Init(c.String("sqlite")); err != nil {

			return ctx, fmt.Errorf("数据库初始化失败 %w", err)
		}

		return ctx, task.Init()
	},
	Action: func(ctx context.Context, cmd *cli.Command) error {
		hash := utils.Md5String(time.Now().String())

		username := hash[8:16]
		password := hash[0:8]
		entrance := fmt.Sprintf("/%s", hash[10:20])
		encrypt, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

		model.SetK(model.AdminSecure, entrance)
		model.SetK(model.AdminUsername, username)
		model.SetK(model.AdminPassword, string(encrypt))

		fmt.Println("重置成功，对应信息如下：")
		fmt.Printf("管理员账号：%s\n管理员密码：%s\n后台管理入口：%s\n", username, password, entrance)
		fmt.Println("请妥善保存以上信息！")
		fmt.Println("-------------------------------")

		return nil
	},
}
