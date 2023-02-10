import adapter from '@sveltejs/adapter-vercel';

const config = {
    kit: {
        adapter: adapter({
            runtime: 'edge',
            regions: ['fra1'],
        }),
    },
};

export default config;
