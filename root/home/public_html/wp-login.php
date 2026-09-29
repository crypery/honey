<!DOCTYPE html>
	<html lang="en-US">
	<head>
	<meta http-equiv="Content-Type" content="text/html; charset=UTF-8" />
	<title>WordPress</title>
	<meta name='robots' content='max-image-preview:large, noindex, noarchive' />
	<meta name="viewport" content="width=device-width" />
	<link rel='stylesheet' id='bootstrap-css' href='/wp-includes/css/bootstrap.css' type='text/css' media='all' />
	<link rel='stylesheet' id='bootstrap-theme-css' href='/wp-includes/css/bootstrap-theme.css' type='text/css' media='all' />
	<script type='text/javascript' src='/wp-includes/js/jquery.js' id='jquery-js'></script>
	<script type='text/javascript' src='/wp-includes/js/bootstrap.js' id='bootstrap-js'></script>
	<style type="text/css">
		#login label { display: block; margin-bottom: 5px; font-weight: bold; }
		#login .forgetmenot label { display: inline; font-weight: normal; }
		#login #logo { margin: 0; font-size: 24px; }
	</style>
	</head>
	<body>
	<h1 class="sr-only">Log In</h1>
	<div id="login">
		<div class="container">
			<div class="row">
				<div class="col-sm-6 col-sm-offset-3 col-md-6 col-md-offset-3">
					<div class="panel panel-default" style="margin-top: 40px;">
						<div class="panel-heading text-center">
							<h1 id="logo"><a href="/"><span class="glyphicon glyphicon-fire"></span> WordPress</a></h1>
						</div>
						<div class="panel-body">
							<form name="loginform" id="loginform" action="wp-login.php" method="get">
								<p>
									<label for="user_login">Username</label>
									<input type="text" name="log" id="user_login" class="input form-control" value="" size="20" autocapitalize="off" autocomplete="username" required="required" />
								</p>

								<p>
									<label for="user_pass">Password</label>
									<input type="password" name="pwd" id="user_pass" class="input password-input form-control" value="" size="20" autocomplete="current-password" spellcheck="false" required="required" />
								</p>
								<p class="submit text-right">
									<input type="submit" name="wp-submit" id="wp-submit" class="button button-primary button-large btn btn-primary" value="Log In" />
									<input type="hidden" name="redirect_to" value="http://wp/wp-admin/" />
									<input type="hidden" name="testcookie" value="1" />
								</p>
							</form>
						</div>
					</div>
					<p id="backtoblog" class="text-center"><a href="/">&larr; Back to Main</a></p>
				</div>
			</div>
		</div>
	</div>
	</body>
	</html>
