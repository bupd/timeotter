import { defineConfig } from 'astro/config';
import starlight from '@astrojs/starlight';

export default defineConfig({
  site: 'https://timeotter.bupd.xyz',
  integrations: [
    starlight({
      title: 'TimeOtter',
      description: 'Calendar-driven task execution for your system',
      logo: {
        src: './src/assets/otter-logo.png',
        replacesTitle: false,
      },
      social: {
        github: 'https://github.com/bupd/timeotter',
      },
      head: [
        {
          tag: 'script',
          attrs: {
            async: true,
            src: 'https://www.googletagmanager.com/gtag/js?id=G-L2HCTNQ0WM',
          },
        },
        {
          tag: 'script',
          content: `
            window.dataLayer = window.dataLayer || [];
            function gtag(){dataLayer.push(arguments);}
            gtag('js', new Date());
            gtag('config', 'G-L2HCTNQ0WM');
          `,
        },
      ],
      sidebar: [
        {
          label: 'Getting Started',
          items: [
            { label: 'Introduction', link: '/getting-started/' },
            { label: 'Installation', link: '/installation/' },
          ],
        },
        {
          label: 'Setup',
          items: [
            { label: 'OAuth Setup', link: '/oauth-setup/' },
            { label: 'Configuration', link: '/configuration/' },
            { label: 'Cron Setup', link: '/cron-setup/' },
          ],
        },
        {
          label: 'Reference',
          items: [
            { label: 'Troubleshooting', link: '/troubleshooting/' },
          ],
        },
      ],
      customCss: ['./src/styles/custom.css'],
    }),
  ],
});
