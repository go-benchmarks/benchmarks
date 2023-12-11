import {sveltekit} from '@sveltejs/kit/vite';
import {defineConfig} from 'vite';
import {preprocess} from "svelte/compiler";

export default defineConfig({
    plugins: [sveltekit()],
    server: {
        fs: {
            allow: ['..'],
        }
    },
});
