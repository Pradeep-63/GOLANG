// emails/template.go
package emails

import (
	"fmt"
	"theransticslabs/m/config"
)

func CommonEmailTemplate(bodyContent string) string {
	apiUrl := config.AppConfig.ApiUrl
	appUrl := config.AppConfig.AppUrl

	return fmt.Sprintf(`
	<!doctype html>
	<html lang="en">
	  <head>
		<meta charset="utf-8">
		<meta name="viewport" content="width=device-width, initial-scale=1">
		<title>Theranostics</title>
		<link rel="preconnect" href="https://fonts.googleapis.com">
		<link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
		<link href="https://fonts.googleapis.com/css2?family=Inria+Sans:ital,wght@0,300;0,400;0,700;1,300;1,400;1,700&display=swap" rel="stylesheet">
		<style>
			body { font-family: 'Inria Sans', sans-serif; font-weight: 400; color: #545454; line-height: 24px; font-size: 16px; }
		</style>
	  </head>
	  <body>
	<center>
		<table cellpadding="0" cellspacing="0" border="0" width="100%%" style="max-width: 600px;">
			<tr>
				<td height='30'></td>
			</tr>
			<tr>
				<td>
					<table style='border: 1px solid #E2E2E2; border-radius:16px' width='100%%' cellspacing='0' cellpadding='0' border='0' bgcolor='#ffffff'>
						<tr>
							<td>
								<table width='100%%' cellspacing='0' cellpadding='0' border='0' bgcolor="#75AC71" style='border-radius:16px 16px 0 0'>
									<tr>
										<td height='30'></td>
									</tr>
									<tr>
										<td>
											<table width='100%%' cellspacing='0' cellpadding='0'>
												<tr>
													<td width='60'></td>
													<td width='150'>
														<a href='%s'>
															<img src="%s/images/logo.png" alt="Logo" style="max-width: 100%%; height: auto;">
														</a>
													</td>
													<td></td>
													<td width='150' align='right'>
														<a href='%s'>
															<img src="%s/images/vector.png" alt="Vector" style="max-width: 100%%; height: auto;">
														</a>
													</td>
													<td width='60'></td>
												</tr>
											</table>
										</td>
									</tr>
									<tr>
										<td height='30'></td>
									</tr>
								</table>
							</td>
						</tr>
						<tr>
							<td>
								<table width='100%%' cellspacing='0' cellpadding='0'>
									<tr>
										<td width='25'></td>
										<td>
											%s
										</td>
										<td width='25'></td>
									</tr>
								</table>
							</td>
						</tr>
					</table>
				</td>
			</tr>
			<tr>
				<td height='30'></td>
			</tr>
		</table>
	</center>
	</body>
	</html>
	`, appUrl, apiUrl, appUrl, apiUrl, bodyContent)
}
