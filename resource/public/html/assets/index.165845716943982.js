import{_ as u,u as l,a6 as a,a7 as n}from"./index.1658457169439.js";import{j as _,aF as i,aw as c,G as d,D as m,M as f,R as p,k as y,c as v}from"./vue.1658457169439.js";const L=_({name:"layout",components:{defaults:i(()=>a(()=>import("./defaults.1658457169439.js"),["assets/defaults.1658457169439.js","assets/vue.1658457169439.js","assets/index.1658457169439.js","assets/index.1658457169439.css","assets/aside.1658457169439.js","assets/main.1658457169439.js","assets/main.1658457169439.css","assets/logo-mini.1658457169439.js","assets/sortable.esm.1658457169439.js","assets/parent.1658457169439.js"])),classic:i(()=>a(()=>import("./classic.1658457169439.js"),["assets/classic.1658457169439.js","assets/vue.1658457169439.js","assets/index.1658457169439.js","assets/index.1658457169439.css","assets/aside.1658457169439.js","assets/main.1658457169439.js","assets/main.1658457169439.css","assets/logo-mini.1658457169439.js","assets/sortable.esm.1658457169439.js","assets/parent.1658457169439.js"])),transverse:i(()=>a(()=>import("./transverse.1658457169439.js"),["assets/transverse.1658457169439.js","assets/main.1658457169439.js","assets/main.1658457169439.css","assets/vue.1658457169439.js","assets/index.1658457169439.js","assets/index.1658457169439.css","assets/logo-mini.1658457169439.js","assets/sortable.esm.1658457169439.js","assets/parent.1658457169439.js"])),columns:i(()=>a(()=>import("./columns.1658457169439.js"),["assets/columns.1658457169439.js","assets/columns.1658457169439.css","assets/vue.1658457169439.js","assets/index.1658457169439.js","assets/index.1658457169439.css","assets/aside.1658457169439.js","assets/main.1658457169439.js","assets/main.1658457169439.css","assets/logo-mini.1658457169439.js","assets/sortable.esm.1658457169439.js","assets/parent.1658457169439.js"]))},setup(){const{proxy:o}=v(),r=l(),{themeConfig:e}=c(r),t=()=>{n.get("oldLayout")||n.set("oldLayout",e.value.layout);const s=document.body.clientWidth;s<1e3?(e.value.isCollapse=!1,o.mittBus.emit("layoutMobileResize",{layout:"defaults",clientWidth:s})):o.mittBus.emit("layoutMobileResize",{layout:n.get("oldLayout")?n.get("oldLayout"):e.value.layout,clientWidth:s})};return d(()=>{t(),window.addEventListener("resize",t)}),m(()=>{window.removeEventListener("resize",t)}),{themeConfig:e}}});function E(o,r,e,t,s,C){return y(),f(p(o.themeConfig.layout))}var h=u(L,[["render",E]]);export{h as default};