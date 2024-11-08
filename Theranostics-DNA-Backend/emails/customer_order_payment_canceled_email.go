// emails/customer_order_payment_canceled_email.go

package emails

import (
	"fmt"
)

func PaymentFailedEmail(firstName, lastName, productName string) string {
	bodyContent := fmt.Sprintf(`
	<table width='100%%' cellspacing='0' cellpadding='0'>
		<tbody>
			<tr>
				<td height='30'></td>
			</tr>
			<tr>
				<td style="color: #000; font-size: 28px; font-weight: 700; text-align: center;">Order Confirmation</td>
			</tr>
			<tr>
				<td height='20'></td>
			</tr>
			<tr>
				<td style="">Dear %s %s,</td>
			</tr>
			<tr>
				<td height='20'></td>
			</tr>
			<tr>
				<td style="">Thank you for your order! Here are your order details:</td>
			</tr>
			<tr>
				<td height='20'></td>
			</tr>
			<tr>
				<td>
					<strong>Product:</strong> %s<br>
				</td>
			</tr>
			<tr>
				<td height='20'></td>
			</tr>
			<tr>
				<td style="text-align: center;">
					Your Payment is failed.
				</td>
			</tr>
			<tr>
				<td height='20'></td>
			</tr>
			<tr>
				<td style="text-align: center;">
					If you have any questions, please contact us at <a href="mailto:support@example.com" style="color:#75AC71; text-decoration: underline;">support@example.com</a>
				</td>
			</tr>
		</tbody>
	</table>
	`, firstName, lastName, productName)

	return CommonEmailTemplate(bodyContent)
}
