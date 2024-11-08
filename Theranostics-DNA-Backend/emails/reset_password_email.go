// emails/reset_password_email.go
package emails

import (
	"fmt"
)

func ResetPasswordEmail(firstName, lastName, email, newPassword, appUrl string) string {
	bodyContent := fmt.Sprintf(`
	<table width='100%%' cellspacing='0' cellpadding='0'>
	<tbody>
		<tr>
			<td height='30'></td>
		</tr>
		<tr>
			<td style="color: #000; font-size: 28px; font-weight: 700; text-align: center;">Your password has been changed.</td>
		</tr>
		<tr>
			<td height='20'></td>
		</tr>
		<tr>
			<td style="">Hi %s %s,
			</td>
		</tr>
		<tr>
			<td height='20'></td>
		</tr>
		<tr>
			<td style=" ">You can log in using your email: <strong>%s</strong></td>
		</tr>
		<tr>
			<td height='20'></td>
		</tr>
		<tr>
			<td style=" ">Your new password is: <strong>%s</strong></td>
		</tr>
		<tr>
			<td height='20'></td>
		</tr>
		<tr>
			<td style=" text-align: center;"><a href="%s" style="background:#75AC71; font-weight: 700; color:#fff; padding: 15px 20px; border-radius: 6px; border:none; cursor: pointer; text-decoration: none; display: inline-block;">Let's Explore</a></td>
		</tr>
		<tr>
			<td height='20'></td>
		</tr>
		<tr>
			<td style=" text-align: center;">
				If you have any question,  please email us at <a href="mailto:official.labs@theranostic.com" style="color:#75AC71; text-decoration: underline;">official.labs@theranostic.com</a> or visit us our FAQ's on our website .
			</td>
		</tr>
		<tr>
			<td height='20'></td>
		</tr>
	</tbody>
</table>
	`, firstName, lastName, email, newPassword, appUrl)

	return CommonEmailTemplate(bodyContent)
}
