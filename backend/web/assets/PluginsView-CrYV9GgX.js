import{G as $,h,g as b,d as F,i as B,r as v,j as x,o as t,c as l,u as a,k,l as w,F as g,m as y,S as C,a as c,t as m,_ as P}from"./index-VTW-DXlx.js";async function E(){const i=await $("/plugins");return h("Failed to get plugins",i),i}async function I(i,r){const u=await b(`/plugins/install?name=${i}&version=${r}`);h("Failed to install plugin",u)}async function G(i){const r=await b(`/plugins/uninstall?name=${i}`);h("Failed to uninstall plugin",r)}const M={key:2},R={class:"row pb-4"},U={class:"col-2 text-center"},j=["href"],D={class:"col-6"},L={class:"col-2"},N=["onChange"],W=["value"],q={class:"col-2"},z=["onClick"],A=["onClick"],J=F({__name:"PluginsView",setup(i){const r=B(),u=v(new Map);let d=v();function S(s,n,e){for(const o of s)if(o.name===n){for(const _ of o.versions)if(_.name===e)return _.installed}return!1}function f(){E().then(s=>{d.value=s;for(const n of d.value)n.versions.length<=0||u.value.set(n.name,n.versions[0].name)})}function p(s){const n=u.value.get(s);if(n===void 0)throw new Error(`selected version for ${s} is undefined`);return n}function V(s,n){const e=n.target;u.value.set(s,e.value)}return x(()=>{f()}),(s,n)=>(t(),l(g,null,[a(r).state===a(C).PluginInstalling?(t(),k(P,{key:0,message:"installing plugin"})):a(r).state===a(C).PluginUninstalling?(t(),k(P,{key:1,message:"uninstalling plugin"})):w("",!0),a(d)!==void 0?(t(),l("table",M,[(t(!0),l(g,null,y(a(d),e=>(t(),l("tr",R,[c("td",U,[c("a",{href:e.url},m(e.name),9,j)]),c("td",D,m(e.description),1),c("td",L,[c("select",{class:"w-100 form-select",onChange:o=>V(e.name,o)},[(t(!0),l(g,null,y(e.versions,o=>(t(),l("option",{value:o.name},m(o.name),9,W))),256))],40,N)]),c("td",q,[S(a(d),e.name,p(e.name))?(t(),l("button",{key:0,class:"btn btn-outline-info w-100",onClick:o=>a(G)(e.name).then(_=>f())}," Uninstall ",8,z)):(t(),l("button",{key:1,class:"btn btn-outline-info w-100",onClick:o=>a(I)(e.name,p(e.name)).then(_=>f())}," Install ",8,A))])]))),256))])):w("",!0)],64))}});export{J as default};
