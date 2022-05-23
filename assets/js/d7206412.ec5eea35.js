"use strict";(self.webpackChunkwebsite=self.webpackChunkwebsite||[]).push([[5442],{3905:function(e,t,n){n.d(t,{Zo:function(){return u},kt:function(){return d}});var r=n(7294);function o(e,t,n){return t in e?Object.defineProperty(e,t,{value:n,enumerable:!0,configurable:!0,writable:!0}):e[t]=n,e}function a(e,t){var n=Object.keys(e);if(Object.getOwnPropertySymbols){var r=Object.getOwnPropertySymbols(e);t&&(r=r.filter((function(t){return Object.getOwnPropertyDescriptor(e,t).enumerable}))),n.push.apply(n,r)}return n}function c(e){for(var t=1;t<arguments.length;t++){var n=null!=arguments[t]?arguments[t]:{};t%2?a(Object(n),!0).forEach((function(t){o(e,t,n[t])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(n)):a(Object(n)).forEach((function(t){Object.defineProperty(e,t,Object.getOwnPropertyDescriptor(n,t))}))}return e}function i(e,t){if(null==e)return{};var n,r,o=function(e,t){if(null==e)return{};var n,r,o={},a=Object.keys(e);for(r=0;r<a.length;r++)n=a[r],t.indexOf(n)>=0||(o[n]=e[n]);return o}(e,t);if(Object.getOwnPropertySymbols){var a=Object.getOwnPropertySymbols(e);for(r=0;r<a.length;r++)n=a[r],t.indexOf(n)>=0||Object.prototype.propertyIsEnumerable.call(e,n)&&(o[n]=e[n])}return o}var p=r.createContext({}),s=function(e){var t=r.useContext(p),n=t;return e&&(n="function"==typeof e?e(t):c(c({},t),e)),n},u=function(e){var t=s(e.components);return r.createElement(p.Provider,{value:t},e.children)},l={inlineCode:"code",wrapper:function(e){var t=e.children;return r.createElement(r.Fragment,{},t)}},h=r.forwardRef((function(e,t){var n=e.components,o=e.mdxType,a=e.originalType,p=e.parentName,u=i(e,["components","mdxType","originalType","parentName"]),h=s(n),d=o,f=h["".concat(p,".").concat(d)]||h[d]||l[d]||a;return n?r.createElement(f,c(c({ref:t},u),{},{components:n})):r.createElement(f,c({ref:t},u))}));function d(e,t){var n=arguments,o=t&&t.mdxType;if("string"==typeof e||o){var a=n.length,c=new Array(a);c[0]=h;var i={};for(var p in t)hasOwnProperty.call(t,p)&&(i[p]=t[p]);i.originalType=e,i.mdxType="string"==typeof e?e:o,c[1]=i;for(var s=2;s<a;s++)c[s]=n[s];return r.createElement.apply(null,c)}return r.createElement.apply(null,n)}h.displayName="MDXCreateElement"},3984:function(e,t,n){n.r(t),n.d(t,{assets:function(){return u},contentTitle:function(){return p},default:function(){return d},frontMatter:function(){return i},metadata:function(){return s},toc:function(){return l}});var r=n(7462),o=n(3366),a=(n(7294),n(3905)),c=["components"],i={},p="HTTP/S",s={unversionedId:"under-the-hood/cache-backends/honeypot/http",id:"under-the-hood/cache-backends/honeypot/http",title:"HTTP/S",description:"\ud83d\udfe2 Supported",source:"@site/docs/under-the-hood/cache-backends/honeypot/http.md",sourceDirName:"under-the-hood/cache-backends/honeypot",slug:"/under-the-hood/cache-backends/honeypot/http",permalink:"/under-the-hood/cache-backends/honeypot/http",draft:!1,tags:[],version:"current",frontMatter:{},sidebar:"defaultSidebar",previous:{title:"Filesystem",permalink:"/under-the-hood/cache-backends/honeypot/filesystem"},next:{title:"Database",permalink:"/category/database"}},u={},l=[{value:"Sample Filesystem Cache Backend Configuration",id:"sample-filesystem-cache-backend-configuration",level:2}],h={toc:l};function d(e){var t=e.components,n=(0,o.Z)(e,c);return(0,a.kt)("wrapper",(0,r.Z)({},h,n,{components:t,mdxType:"MDXLayout"}),(0,a.kt)("h1",{id:"https"},"HTTP/S"),(0,a.kt)("p",null,(0,a.kt)("strong",{parentName:"p"},"\ud83d\udfe2 Supported")),(0,a.kt)("p",null,"The ",(0,a.kt)("inlineCode",{parentName:"p"},"http")," and ",(0,a.kt)("inlineCode",{parentName:"p"},"https")," cache backends use jsonschemas stored at remote HTTP paths to back the in-memory schema cache."),(0,a.kt)("h2",{id:"sample-filesystem-cache-backend-configuration"},"Sample Filesystem Cache Backend Configuration"),(0,a.kt)("pre",null,(0,a.kt)("code",{parentName:"pre"},"schemaCache:\n  backend:\n    type: https                     # The backend type\n    host: registry.silverton.io     # The schema host\n    path: /some/path/somewhere      # The path to consider as root\n")))}d.isMDXComponent=!0}}]);