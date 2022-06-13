const path = require('path');

module.exports = function (context) {
  const { siteConfig } = context;
  const { themeConfig } = siteConfig;
  const { snowplow } = themeConfig || {};

  if (!snowplow) {
    throw new Error(`'snowplow' object in 'themeConfig' must be specified`);
  }

  const {
    collector,
    appId,
    withCredentials
  } = snowplow;

  return {
    name: 'docusaurus-plugin-snowplow',
    getClientModules() {
      return [path.resolve(__dirname, "./snowplow")]
    },
    injectHtmlTags: () => {
      return {
        headTags: [
          {
            tagName: 'script',
            innerHTML: `
            ;(function(p,l,o,w,i,n,g){if(!p[i]){p.GlobalSnowplowNamespace=p.GlobalSnowplowNamespace||[]; p.GlobalSnowplowNamespace.push(i);p[i]=function(){(p[i].q=p[i].q||[]).push(arguments) };p[i].q=p[i].q||[];n=l.createElement(o);g=l.getElementsByTagName(o)[0];n.async=1; n.src=w;g.parentNode.insertBefore(n,g)}}(window,document,"script","https://cdn.jsdelivr.net/npm/@snowplow/javascript-tracker@3.4.0/dist/sp.js","snowplow"));
            snowplow('newTracker', 'sp', '${collector}', {
                appId: '${appId}',
                platform: 'web',
                withCredentials: ${withCredentials},
                cookieSameSite: 'Lax',
                discoverRootDomain: true
            });
            `
          },
        ],
      }
    },
  }
}
