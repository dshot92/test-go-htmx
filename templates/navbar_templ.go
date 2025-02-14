// Code generated by templ - DO NOT EDIT.

// templ: version: v0.3.833
package templates

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

func Navbar() templ.Component {
	return templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
		templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
		if templ_7745c5c3_CtxErr := ctx.Err(); templ_7745c5c3_CtxErr != nil {
			return templ_7745c5c3_CtxErr
		}
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
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 1, "<nav class=\"navbar\"><div class=\"container navbar-content\"><a href=\"/\" class=\"site-title\">3D Model Viewer</a> <button id=\"theme-toggle\" type=\"button\" class=\"theme-toggle\"><svg id=\"theme-toggle-dark-icon\" class=\"hidden\" width=\"20\" height=\"20\" fill=\"currentColor\" viewBox=\"0 0 20 20\" xmlns=\"http://www.w3.org/2000/svg\"><path d=\"M17.293 13.293A8 8 0 016.707 2.707a8.001 8.001 0 1010.586 10.586z\"></path></svg> <svg id=\"theme-toggle-light-icon\" class=\"hidden\" width=\"20\" height=\"20\" fill=\"currentColor\" viewBox=\"0 0 20 20\" xmlns=\"http://www.w3.org/2000/svg\"><path d=\"M10 2a1 1 0 011 1v1a1 1 0 11-2 0V3a1 1 0 011-1zm4 8a4 4 0 11-8 0 4 4 0 018 0zm-.464 4.95l.707.707a1 1 0 001.414-1.414l-.707-.707a1 1 0 00-1.414 1.414zm2.12-10.607a1 1 0 010 1.414l-.706.707a1 1 0 11-1.414-1.414l.707-.707a1 1 0 011.414 0zM17 11a1 1 0 100-2h-1a1 1 0 100 2h1zm-7 4a1 1 0 011 1v1a1 1 0 11-2 0v-1a1 1 0 011-1zM5.05 6.464A1 1 0 106.465 5.05l-.708-.707a1 1 0 00-1.414 1.414l.707.707zm1.414 8.486l-.707.707a1 1 0 01-1.414-1.414l.707-.707a1 1 0 011.414 1.414zM4 11a1 1 0 100-2H3a1 1 0 000 2h1z\"></path></svg></button></div></nav><script>\n        // Wait for DOM to be fully loaded\n        document.addEventListener('DOMContentLoaded', function() {\n            // Initialize theme\n            if (localStorage.getItem('color-theme') === 'dark' || (!('color-theme' in localStorage) && window.matchMedia('(prefers-color-scheme: dark)').matches)) {\n                document.documentElement.classList.add('dark');\n            } else {\n                document.documentElement.classList.remove('dark');\n            }\n\n            const themeToggleDarkIcon = document.getElementById('theme-toggle-dark-icon');\n            const themeToggleLightIcon = document.getElementById('theme-toggle-light-icon');\n            const themeToggleBtn = document.getElementById('theme-toggle');\n\n            // Set initial icon state\n            if (document.documentElement.classList.contains('dark')) {\n                themeToggleLightIcon.classList.remove('hidden');\n                themeToggleDarkIcon.classList.add('hidden');\n            } else {\n                themeToggleLightIcon.classList.add('hidden');\n                themeToggleDarkIcon.classList.remove('hidden');\n            }\n\n            // Handle theme toggle\n            themeToggleBtn.addEventListener('click', function() {\n                // Toggle icons\n                themeToggleDarkIcon.classList.toggle('hidden');\n                themeToggleLightIcon.classList.toggle('hidden');\n\n                // Toggle theme\n                if (document.documentElement.classList.contains('dark')) {\n                    document.documentElement.classList.remove('dark');\n                    localStorage.setItem('color-theme', 'light');\n                } else {\n                    document.documentElement.classList.add('dark');\n                    localStorage.setItem('color-theme', 'dark');\n                }\n            });\n        });\n    </script>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return nil
	})
}

var _ = templruntime.GeneratedTemplate
