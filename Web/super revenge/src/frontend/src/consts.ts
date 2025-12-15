// Place any global data in this file.
// You can import this data from anywhere in your site by using the `import` keyword.

export const SITE_TITLE = 'Super Notes';
export const SITE_DESCRIPTION = 'Welcome to Super Notes!, where you can store and manage your notes efficiently.';

const backendBase =
    import.meta.env.API_BASE_URL ??
    import.meta.env.PUBLIC_API_BASE_URL ??
    'http://localhost:8080';

export const apiBackend = backendBase + '/api/v1';

console.log('API Backend URL:', apiBackend);
