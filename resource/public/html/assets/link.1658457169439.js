var l=Object.defineProperty;var n=Object.getOwnPropertySymbols;var f=Object.prototype.hasOwnProperty,m=Object.prototype.propertyIsEnumerable;var i=(e,a,t)=>a in e?l(e,a,{enumerable:!0,configurable:!0,writable:!0,value:t}):e[a]=t,u=(e,a)=>{for(var t in a||(a={}))f.call(a,t)&&i(e,t,a[t]);if(n)for(var t of n(a))m.call(a,t)&&i(e,t,a[t]);return e};import{j as p,aw as h,aB as k,$ as g,x as d,b as v,ac as R,l as w,m as _,U as c,J as $,k as y}from"./vue.1658457169439.js";import{_ as L,u as M}from"./index.1658457169439.js";const B=p({name:"layoutLinkView",setup(){const e=M(),{themeConfig:a}=h(e),t=k(),s=g({currentRouteMeta:{isLink:"",title:""}}),o=d(()=>{let{isTagsview:r}=a.value;return r?"115px":"80px"});return v(()=>t.path,()=>{s.currentRouteMeta=t.meta},{immediate:!0}),u({setLinkHeight:o},R(s))}}),C=["href"];function T(e,a,t,s,o,r){return y(),w("div",{class:"layout-view-bg-white flex layout-view-link",style:$({height:`calc(100vh - ${e.setLinkHeight}`})},[_("a",{href:e.currentRouteMeta.isLink,target:"_blank",rel:"opener",class:"flex-margin"},c(e.$t(e.currentRouteMeta.title))+"\uFF1A"+c(e.currentRouteMeta.isLink),9,C)],4)}var H=L(B,[["render",T]]);export{H as default};
