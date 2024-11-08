// emails/customer_order_confirmation_email.go

package emails

import (
	"fmt"
	"theransticslabs/m/config"
)

func CustomerOrderConfirmationEmail(firstName, lastName, productName string, quantity int, totalPrice float64, invoiceLink, appUrl string) string {
	apiUrl := config.AppConfig.ApiUrl

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
					<strong>Quantity:</strong> %d<br>
					<strong>Total Amount:</strong> $%.2f
				</td>
			</tr>
			<tr>
				<td height='20'></td>
			</tr>
			<tr>
				<td style="text-align: center;">
					<a href="%s/%s" style="background:#75AC71; font-weight: 700; color:#fff; padding: 15px 20px; border-radius: 6px; border:none; cursor: pointer; text-decoration: none; display: inline-block;">View Invoice</a>
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
	`, firstName, lastName, productName, quantity, totalPrice, apiUrl, invoiceLink)

	return CommonEmailTemplate(bodyContent)
}
