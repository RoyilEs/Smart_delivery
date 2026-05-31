package email_ser

import (
	"Smart_delivery_locker/global"
	"bytes"
	"crypto/tls"
	"fmt"
	"html/template"
	"net/smtp"
	"time"
)

type EmailService struct{}

// EmailRequest 邮件请求结构
type EmailRequest struct {
	ToEmail      string                 // 收件人邮箱
	ToName       string                 // 收件人名称
	Subject      string                 // 邮件主题
	Template     string                 // 邮件模板名称
	TemplateData map[string]interface{} // 模板数据
}

func (s *EmailService) SendPickupNotification(toEmail, toName, pickupCode, cabinetCode, grilleId, expireAt string) error {
	// 参数验证
	if toEmail == "" {
		return fmt.Errorf("收件人邮箱不能为空")
	}

	if pickupCode == "" {
		return fmt.Errorf("取件码不能为空")
	}

	// 计算剩余时间
	var expireTime time.Time
	var err error

	if expireAt != "" {
		expireTime, err = time.Parse(time.RFC3339, expireAt)
		if err != nil {
			global.Log.Errorf("解析过期时间失败: %v, 使用默认时间", err)
			expireTime = time.Now().Add(72 * time.Hour)
		}
	} else {
		global.Log.Warn("过期时间为空，使用默认72小时")
		expireTime = time.Now().Add(72 * time.Hour)
	}

	remainingHours := int(time.Until(expireTime).Hours())
	if remainingHours < 0 {
		remainingHours = 0
	}

	if toName == "" {
		toName = "用户"
	}

	templateData := map[string]interface{}{
		"ToName":         toName,
		"PickupCode":     pickupCode,
		"CabinetCode":    cabinetCode,
		"GrilleId":       grilleId,
		"ExpireAt":       expireAt,
		"RemainingHours": remainingHours,
		"CurrentYear":    time.Now().Year(),
	}

	req := &EmailRequest{
		ToEmail:      toEmail,
		ToName:       toName,
		Subject:      fmt.Sprintf("【智能快递柜】您有一个包裹已到达，取件码：%s", pickupCode),
		Template:     "pickup_notification",
		TemplateData: templateData,
	}

	return s.SendEmail(req)
}

// SendEmail 通用邮件发送方法
func (s *EmailService) SendEmail(req *EmailRequest) error {
	// 加载并解析模板
	tmpl, err := s.getTemplate(req.Template)
	if err != nil {
		global.Log.Errorf("加载邮件模板失败: %v", err)
		return err
	}

	// 渲染模板
	var body bytes.Buffer
	if err := tmpl.Execute(&body, req.TemplateData); err != nil {
		global.Log.Errorf("渲染邮件模板失败: %v", err)
		return err
	}

	// 发送邮件
	return s.sendVia163SMTP(req.ToEmail, req.Subject, body.String(), req.ToName)
}

// sendVia163SMTP 通过163邮箱SMTP发送邮件（支持SSL）
func (s *EmailService) sendVia163SMTP(toEmail, subject, body, toName string) error {
	// 获取配置
	host := global.Config.FromMail.Host
	port := global.Config.FromMail.Port
	from := global.Config.FromMail.From
	password := global.Config.FromMail.Password

	// 163邮箱使用SSL连接，需要配置TLS
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}

	// 连接到SMTP服务器
	conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%d", host, port), tlsConfig)
	if err != nil {
		global.Log.Errorf("连接SMTP服务器失败: %v", err)
		return err
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, host)
	if err != nil {
		global.Log.Errorf("创建SMTP客户端失败: %v", err)
		return err
	}
	defer client.Close()

	// 认证
	auth := smtp.PlainAuth("", from, password, host)
	if err = client.Auth(auth); err != nil {
		global.Log.Errorf("SMTP认证失败: %v", err)
		return err
	}

	// 设置发件人
	if err = client.Mail(from); err != nil {
		global.Log.Errorf("设置发件人失败: %v", err)
		return err
	}

	// 设置收件人
	if err = client.Rcpt(toEmail); err != nil {
		global.Log.Errorf("设置收件人失败: %v", err)
		return err
	}

	// 获取写入流
	w, err := client.Data()
	if err != nil {
		global.Log.Errorf("获取数据写入流失败: %v", err)
		return err
	}
	defer w.Close()

	// 构建邮件头部
	fromName := "智能快递柜"
	headers := make(map[string]string)
	headers["From"] = fmt.Sprintf("%s <%s>", fromName, from)
	headers["To"] = fmt.Sprintf("%s <%s>", toName, toEmail)
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=UTF-8"

	// 写入邮件内容
	for k, v := range headers {
		fmt.Fprintf(w, "%s: %s\r\n", k, v)
	}
	fmt.Fprintf(w, "\r\n%s", body)

	// 发送
	if err = w.Close(); err != nil {
		global.Log.Errorf("发送邮件失败: %v", err)
		return err
	}

	// 发送QUIT命令
	if err = client.Quit(); err != nil {
		global.Log.Warnf("关闭SMTP连接失败: %v", err)
	}

	global.Log.Infof("邮件发送成功: %s -> %s", from, toEmail)
	return nil
}

// getTemplate 获取邮件模板
func (s *EmailService) getTemplate(templateName string) (*template.Template, error) {
	var templateStr string

	switch templateName {
	case "pickup_notification":
		templateStr = s.getPickupNotificationTemplate()
	default:
		templateStr = s.getDefaultTemplate()
	}

	return template.New(templateName).Parse(templateStr)
}

// getPickupNotificationTemplate 取件通知模板
func (s *EmailService) getPickupNotificationTemplate() string {
	return `<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>取件通知</title>
    <style>
        body {
            font-family: 'Microsoft YaHei', Arial, sans-serif;
            line-height: 1.6;
            color: #333;
            margin: 0;
            padding: 0;
            background-color: #f5f5f5;
        }
        .container {
            max-width: 600px;
            margin: 20px auto;
            background-color: #fff;
            border-radius: 10px;
            overflow: hidden;
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
        }
        .header {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            padding: 30px;
            text-align: center;
        }
        .header h1 {
            margin: 0;
            font-size: 24px;
        }
        .content {
            padding: 30px;
        }
        .greeting {
            font-size: 18px;
            margin-bottom: 20px;
        }
        .code-box {
            background-color: #f0f7ff;
            border: 2px dashed #667eea;
            border-radius: 10px;
            padding: 20px;
            text-align: center;
            margin: 20px 0;
        }
        .code {
            font-size: 32px;
            font-weight: bold;
            color: #667eea;
            letter-spacing: 5px;
            font-family: monospace;
        }
        .info-item {
            margin: 15px 0;
            padding: 10px;
            background-color: #f9f9f9;
            border-left: 4px solid #667eea;
        }
        .warning {
            background-color: #fff3cd;
            border-left-color: #ffc107;
            color: #856404;
        }
        .footer {
            background-color: #f8f9fa;
            padding: 20px;
            text-align: center;
            font-size: 12px;
            color: #6c757d;
        }
        @media (max-width: 600px) {
            .container {
                margin: 10px;
            }
            .code {
                font-size: 24px;
            }
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🏪 智能快递柜取件通知</h1>
        </div>
        <div class="content">
            <div class="greeting">
                尊敬的 <strong>{{.ToName}}</strong>，您好！
            </div>
            <p>您的包裹已成功存入智能快递柜，请凭以下取件码前往取件。</p>
            
            <div class="code-box">
                <div style="margin-bottom: 10px;">您的取件码：</div>
                <div class="code">{{.PickupCode}}</div>
            </div>
            
            <div class="info-item">
                <strong>📍 柜子编号：</strong> {{.CabinetCode}}
            </div>
            <div class="info-item">
                <strong>📦 格口编号：</strong> {{.GrilleId}}
            </div>
            <div class="info-item warning">
                <strong>⏰ 取件时限：</strong> {{.RemainingHours}} 小时<br>
                <small>请在 {{.ExpireAt}} 前取件，超时将自动退回</small>
            </div>
            
            <p style="margin-top: 20px;">取件步骤：</p>
            <ol>
                <li>前往快递柜屏幕输入取件码</li>
                <li>或输入手机号获取所有快递</li>
                <li>根据提示打开对应格口</li>
                <li>取出包裹并关闭柜门</li>
            </ol>
        </div>
        <div class="footer">
            <p>本邮件由智能快递柜系统自动发送，请勿直接回复</p>
            <p>&copy; {{.CurrentYear}} 智能快递柜 · 让生活更便捷</p>
        </div>
    </div>
</body>
</html>`
}

// getDefaultTemplate 默认邮件模板
func (s *EmailService) getDefaultTemplate() string {
	return `<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>邮件通知</title>
</head>
<body>
    <h2>{{.Title}}</h2>
    <p>{{.Content}}</p>
    <hr>
    <p>系统自动发送，请勿回复</p>
</body>
</html>`
}
