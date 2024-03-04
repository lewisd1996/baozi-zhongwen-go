// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.543
package components

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

import "github.com/lewisd1996/baozi-zhongwen/internal/view/icons"
import "strings"

func Navigation(userId, route string) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, templ_7745c5c3_W io.Writer) (templ_7745c5c3_Err error) {
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templ_7745c5c3_W.(*bytes.Buffer)
		if !templ_7745c5c3_IsBuffer {
			templ_7745c5c3_Buffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templ_7745c5c3_Buffer)
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		templ_7745c5c3_Err = DesktopNav(userId, route).Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = MobileNav().Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}

func DesktopNav(userId, route string) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, templ_7745c5c3_W io.Writer) (templ_7745c5c3_Err error) {
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templ_7745c5c3_W.(*bytes.Buffer)
		if !templ_7745c5c3_IsBuffer {
			templ_7745c5c3_Buffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templ_7745c5c3_Buffer)
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var2 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var2 == nil {
			templ_7745c5c3_Var2 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"hidden lg:fixed lg:inset-y-0 lg:z-50 lg:flex lg:w-52 lg:flex-col\"><div class=\"flex grow pt-6 flex-col gap-y-5 overflow-y-auto border-r border-slate-200 px-6\"><a class=\"flex items-center\" href=\"/\"><svg class=\"w-9 h-9\" viewBox=\"0 0 1024 1024\" class=\"icon\" version=\"1.1\" xmlns=\"http://www.w3.org/2000/svg\" fill=\"currentColor\"><path d=\"M619.789474 862.315789h-215.578948C232.448 862.315789 107.789474 744.340211 107.789474 581.820632c0-153.842526 110.187789-242.876632 198.736842-314.421895C373.409684 213.369263 431.157895 166.723368 431.157895 107.789474h53.894737c0 84.668632-70.251789 141.419789-144.653474 201.539368C252.550737 380.308211 161.684211 453.712842 161.684211 581.820632 161.684211 715.237053 261.416421 808.421053 404.210526 808.421053h215.578948c142.794105 0 242.526316-93.184 242.526315-226.600421 0-128.107789-90.866526-201.512421-178.714947-272.49179C609.199158 249.209263 538.947368 192.458105 538.947368 107.789474h53.894737c0 58.933895 57.775158 105.579789 124.631579 159.609263C805.995789 338.944 916.210526 427.978105 916.210526 581.820632 916.210526 744.340211 791.552 862.315789 619.789474 862.315789z m2.829473-392.165052l-80.842105-161.684211 48.208842-24.117894 80.842105 161.68421-48.208842 24.117895z m-221.237894 0l-48.208842-24.117895 80.842105-161.68421 48.208842 24.117894-80.842105 161.684211z\" fill=\"#231F20\"></path></svg></a><nav class=\"flex flex-1 flex-col\"><ul role=\"list\" class=\"flex flex-1 flex-col gap-y-7\"><li><ul role=\"list\" class=\"space-y-1\"><li>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var3 = []any{templ.KV("bg-slate-200 text-slate-900", route == "/"), "text-slate-400 group flex items-center gap-x-3 rounded-md p-2 text-sm leading-6 font-semibold"}
		templ_7745c5c3_Err = templ.RenderCSSItems(ctx, templ_7745c5c3_Buffer, templ_7745c5c3_Var3...)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<a href=\"/\" class=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ.CSSClasses(templ_7745c5c3_Var3).String()))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if route == "/" {
			templ_7745c5c3_Err = icons.HomeIcon("w-5 h-5 text-slate-900").Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		} else {
			templ_7745c5c3_Err = icons.HomeIcon("w-5 h-5 text-slate-400").Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("Home</a></li><li>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var4 = []any{templ.KV("bg-slate-200 text-slate-900", strings.Contains(route, "/decks")), "text-slate-400 group flex items-center gap-x-3 rounded-md p-2 text-sm leading-6 font-semibold"}
		templ_7745c5c3_Err = templ.RenderCSSItems(ctx, templ_7745c5c3_Buffer, templ_7745c5c3_Var4...)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<a href=\"/decks\" class=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ.CSSClasses(templ_7745c5c3_Var4).String()))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if strings.Contains(route, "/decks") {
			templ_7745c5c3_Err = icons.DecksIcon("w-5 h-5 text-slate-900").Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		} else {
			templ_7745c5c3_Err = icons.DecksIcon("w-5 h-5 text-slate-400").Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("Decks</a></li></ul></li><li class=\"-mx-6 mt-auto\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = ProfileMenu(userId).Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</li></ul></nav></div></div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}

func MobileNav() templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, templ_7745c5c3_W io.Writer) (templ_7745c5c3_Err error) {
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templ_7745c5c3_W.(*bytes.Buffer)
		if !templ_7745c5c3_IsBuffer {
			templ_7745c5c3_Buffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templ_7745c5c3_Buffer)
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var5 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var5 == nil {
			templ_7745c5c3_Var5 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"relative z-50 lg:hidden\" :class=\"{ &#39;block&#39;: open, &#39;hidden&#39;: !open }\" role=\"dialog\" aria-modal=\"true\"><div class=\"fixed inset-0 bg-slate-900/80\" @click=\"open = false\"></div><div class=\"fixed inset-0 flex\"><div class=\"relative mr-16 flex w-full max-w-xs flex-1\"><div class=\"absolute left-full top-0 flex w-16 justify-center pt-5\"><button type=\"button\" class=\"-m-2.5 p-2.5\" @click=\"open = false\"><span class=\"sr-only\">Close sidebar</span> <svg class=\"h-6 w-6 text-white\" fill=\"none\" viewBox=\"0 0 24 24\" stroke-width=\"1.5\" stroke=\"currentColor\" aria-hidden=\"true\"><path stroke-linecap=\"round\" stroke-linejoin=\"round\" d=\"M6 18L18 6M6 6l12 12\"></path></svg></button></div><div class=\"flex grow flex-col gap-y-5 overflow-y-auto bg-slate-100 px-6 pb-2\"><div class=\"flex h-16 shrink-0 items-center\"><a class=\"text-xl lg:text-3xl text-slate-900 font-bold\" href=\"/\">BAOZI <span class=\"hanzi\">中文</span></a></div><nav class=\"flex flex-1 flex-col\"><ul role=\"list\" class=\"flex flex-1 flex-col gap-y-7\"><li><ul role=\"list\" class=\"-mx-2 space-y-1\"><li><!-- Current: \"bg-slate-50 text-teal-600\", Default: \"text-slate-700 hover:text-teal-600 hover:bg-slate-50\" --><a href=\"#\" class=\"bg-teal-400/10 text-teal-300 group flex gap-x-3 rounded-md p-2 text-sm leading-6 font-semibold\"><svg class=\"h-6 w-6 shrink-0 text-teal-600\" fill=\"none\" viewBox=\"0 0 24 24\" stroke-width=\"1.5\" stroke=\"currentColor\" aria-hidden=\"true\"><path stroke-linecap=\"round\" stroke-linejoin=\"round\" d=\"M2.25 12l8.954-8.955c.44-.439 1.152-.439 1.591 0L21.75 12M4.5 9.75v10.125c0 .621.504 1.125 1.125 1.125H9.75v-4.875c0-.621.504-1.125 1.125-1.125h2.25c.621 0 1.125.504 1.125 1.125V21h4.125c.621 0 1.125-.504 1.125-1.125V9.75M8.25 21h8.25\"></path></svg> Dashboard</a></li></ul></li></ul></nav></div></div></div></div><div class=\"sticky top-0 z-40 flex items-center gap-x-6 bg-slate-100 px-4 py-4 shadow-sm sm:px-6 lg:hidden\"><button type=\"button\" class=\"-m-2.5 p-2.5 text-slate-400 lg:hidden\" @click=\"open = true\"><span class=\"sr-only\">Open sidebar</span> <svg class=\"h-6 w-6\" fill=\"none\" viewBox=\"0 0 24 24\" stroke-width=\"1.5\" stroke=\"currentColor\" aria-hidden=\"true\"><path stroke-linecap=\"round\" stroke-linejoin=\"round\" d=\"M3.75 6.75h16.5M3.75 12h16.5m-16.5 5.25h16.5\"></path></svg></button></div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}
