// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.771
package auth_views

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

func SignUpForm() templ.Component {
	return templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
		templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templruntime.GetBuffer(templ_7745c5c3_W)
		if !templ_7745c5c3_IsBuffer {
			defer func() {
				templ_7745c5c3_BufErr := templruntime.ReleaseBuffer(templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err == nil {
					templ_7745c5c3_Err = templ_7745c5c3_BufErr
				}
			}()
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<section class=\"section sign-in\"><form hx-post=\"/sign-up\" hx-swap=\"outerHTML\" class=\"container is-max-desktop\"><h2 class=\"subtitle is-4\">Sign Up</h2><hr><div class=\"field\"><label class=\"label\" for=\"username\">Username</label><div class=\"control\"><input class=\"input\" id=\"usernameInput\" name=\"username\" type=\"text\" required></div></div><div class=\"field\"><label class=\"label\" for=\"password\">Password</label><div class=\"control\"><input class=\"input\" id=\"passwordInput\" name=\"password\" type=\"password\" required></div></div><div class=\"field\"><label class=\"label\" for=\"password\">Confirm Password</label><div class=\"control\"><input class=\"input\" id=\"repeatPasswordInput\" name=\"repeat_password\" type=\"password\" required></div></div><div class=\"field\"><div class=\"control\"><input class=\"button is-success w-100\" type=\"submit\" value=\"Log In\"></div></div><div class=\"field\"><div class=\"control\"><a href=\"\" hx-get=\"/sign-in\" hx-swap=\"outerHTML\" hx-target=\".sign-in\" hx-push-url=\"true\">Already a user? Log In</a></div></div></form></section>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return templ_7745c5c3_Err
	})
}

var _ = templruntime.GeneratedTemplate
