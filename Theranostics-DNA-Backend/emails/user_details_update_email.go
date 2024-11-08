// emails/user_details_update_email.go
package emails

import (
	"fmt"
)

// UserDetailsUpdatedEmail constructs the HTML body for the UserDetailsUpdatedEmail .
func UserDetailsUpdatedEmail(firstName, lastName, email, appUrl string) string {
	bodyContent := fmt.Sprintf(`
	<table width='100%%' cellspacing='0' cellpadding='0'>
	<tbody>
		<tr>
			<td height='30'></td>
		</tr>
		<tr>
			<td style="color: #000; font-size: 28px; font-weight: 700; text-align: center;">Your Profile Details Updated</td>
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
			<td style=" ">We wanted to let you know that your profile details have been successfully updated. 
			</td>
		</tr>
		<tr>
			<td height='20'></td>
		</tr>
	
		<tr>
		<td style=" ">Here’s a summary of the changes made:
		</td>
	</tr>
	<tr>
		<td height='20'></td>
	</tr>
	<tr>
	<td style=" ">First Name:<strong>%s</strong>
	</td>
</tr>
<tr>
	<td height='20'></td>
</tr>
<tr>
	<td style=" ">Last Name:<strong>%s</strong>
	</td>
</tr>
<tr>
	<td height='20'></td>
</tr>
<tr>
	<td style=" ">Email:<strong>%s</strong>
	</td>
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
	`, firstName, lastName, firstName, lastName, email, appUrl)

	return CommonEmailTemplate(bodyContent)
}
