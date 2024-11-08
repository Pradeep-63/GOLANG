// emails/password_update_email.go
package emails

import (
	"fmt"
)

// PasswordUpdatedEmail constructs the HTML body for the email.
func PasswordUpdatedEmail(firstName, lastName, email, password, appUrl string) string {
	bodyContent := fmt.Sprintf(`
	<table width='100%%' cellspacing='0' cellpadding='0'>
	<tbody>
		<tr>
			<td height='30'></td>
		</tr>
		<tr>
			<td style="color: #000; font-size: 28px; font-weight: 700; text-align: center;">Your Password has been Updated</td>
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
			<td style=" ">This is a confirmation that your password has been successfully updated.</td>
		</tr>
		<tr>
			<td height='20'></td>
		</tr>
		<tr>
			<td style=" ">If you did not initiate this change, please reset your password immediately to protect your account.</td>
		</tr>
		<tr>
			<td height='20'></td>
		</tr>

		<tr>
			<td style=" ">Your temporary password is: <strong>%s</strong>. Please change it after your first login.</td>
		</tr>
		<tr>
			<td height='20'></td>
		</tr>
		<tr>
			<td style=" text-align: center;"><a href="%s" style="background:#75AC71; font-weight: 700; color:#fff; padding: 15px 20px; border-radius: 6px; border:none; cursor: pointer; text-decoration: none; display: inline-block;">Log In</a></td>
		</tr>
		<tr>
			<td height='20'></td>
		</tr>
		<tr>
			<td style=" text-align: center;">
				If you have any questions, please email us at <a href="mailto:official.labs@theranostic.com" style="color:#75AC71; text-decoration: underline;">official.labs@theranostic.com</a> or visit our FAQ's on our website.
			</td>
		</tr>
		<tr>
			<td height='20'></td>
		</tr>
	</tbody>
</table>
	`, firstName, lastName, password, appUrl)

	return CommonEmailTemplate(bodyContent)
}
