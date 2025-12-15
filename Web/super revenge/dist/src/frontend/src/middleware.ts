import { defineMiddleware } from 'astro/middleware';

const AUTH_COOKIES = ['supernote_session', 'session', 'authToken'];

export function getSessionCookie(cookieHeader: string | null): string | null {
	if (!cookieHeader) return null;
	const match = cookieHeader
		.split(';')
		.map((part) => part.trim())
		.find((item) => AUTH_COOKIES.some((name) => item.startsWith(`${name}=`)));
	return match ? match.split('=')[1] ?? null : null;
}

async function isSessionValid(cookieValue: string, request: Request): Promise<boolean> {
	const form = new FormData();
	form.append('endpoint', 'auth/check-token');
	form.append('token', cookieValue);

	try {
		const controlUrl = new URL('/api/control-check', request.url);
		const response = await fetch(controlUrl, {
			method: 'POST',
			headers: { Accept: 'application/json' },
			body: form,
		});
		
		if (!response.ok) return false;
		const payload = await response.json().catch(() => null);
		return payload?.success === true;
	} catch (error) {
		console.info('Error during session validation:', error);

		console.info('Failed to validate session token due to network error.');
		return false;
	}
}

export const onRequest = defineMiddleware(async ({ request }, next) => {
	const url = new URL(request.url);
	console.log(`Middleware checking URL: ${JSON.stringify(url)}`);
	if (url.pathname.startsWith('/notes')) {
		const cookieValue = getSessionCookie(request.headers.get('cookie'));
		console.log(`Checking session cookie for /notes:`, cookieValue);
		if (!cookieValue || !(await isSessionValid(cookieValue, request))) {
			const redirectTarget = new URL('/signin', url.origin);
			redirectTarget.searchParams.set('redirect', url.pathname + url.search);
			return Response.redirect(redirectTarget, 303);
		}
	}

	return next();
});
