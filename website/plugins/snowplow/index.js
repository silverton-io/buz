const path = require('path');

module.exports = function (context) {
    return {
      name: 'docusaurus-plugin-analytics',
      getClientModules() {
        return [path.resolve(__dirname, "./snowplow")]
      },
      injectHtmlTags: () => {
        return {
          headTags: [
            {
              tagName: 'script',
              attributes: {
                async: false,
                src: '/js/snowplow.js',
              },
            },
          ],
        }
      },
    }
  }