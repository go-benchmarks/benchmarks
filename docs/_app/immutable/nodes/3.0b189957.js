import{s as i,n as r,e as m}from"../chunks/scheduler.e108d1fd.js";import{S as c,i as p,m as f,n as u,a as d,o as l,f as _}from"../chunks/index.359c4e26.js";import{p as h}from"../chunks/stores.b240cf77.js";function b(s){let t=s[0].params.benchmark+"",e;return{c(){e=f(t)},l(a){e=u(a,t)},m(a,n){d(a,e,n)},p(a,[n]){n&1&&t!==(t=a[0].params.benchmark+"")&&l(e,t)},i:r,o:r,d(a){a&&_(e)}}}function g(s,t,e){let a;m(s,h,o=>e(0,a=o));let{data:n}=t;return s.$$set=o=>{"data"in o&&e(1,n=o.data)},[a,n]}class q extends c{constructor(t){super(),p(this,t,g,b,i,{data:1})}}export{q as component};