import { page } from '$app/stores';
import { onMount } from 'svelte';

export function rand(min, max) {
    return Math.floor(Math.random() * (max - min + 1)) + min;
}

// Deprecated
export function lazyLoader(total, getCount, fetchPage) {
    let app;
    let loading = false;
    let _page = 2;
    async function checkScroll() {
        if (
            app.scrollTop + app.offsetHeight >= app.scrollHeight - app.offsetHeight * 0.25 &&
            getCount() < total &&
            !loading
        ) {
            loading = true;
            await fetchPage(_page++);
            loading = false;
            checkScroll();
        }
    }
    if (getCount() < total) {
        onMount(() => {
            app = document.querySelector('.app');
            app.addEventListener('scroll', checkScroll);
            const unsubscribe = page.subscribe(() => {
                _page = 2;
                checkScroll();
            });
            checkScroll();
            return () => {
                app.removeEventListener('scroll', checkScroll);
                unsubscribe();
            };
        });
    }
}

export function newLazyLoader({ getTotal, getCount, loadPage }) {
    onMount(() => {
        const app = document.querySelector('.app');
        let _page = 1;
        async function checkScroll() {
            if (
                (getTotal() == null || getCount() < getTotal()) &&
                app.scrollTop + app.offsetHeight >= app.scrollHeight - app.offsetHeight * 0.25
            ) {
                loadPage(_page++);
            }
        }
        app.addEventListener('scroll', checkScroll);
        const unsubscribe = page.subscribe(() => {
            _page = 1;
            loadPage(_page++);
        });
        loadPage(_page++);
        return () => {
            app.removeEventListener('scroll', checkScroll);
            unsubscribe();
        };
    });
}
