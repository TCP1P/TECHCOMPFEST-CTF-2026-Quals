import type { APIContext } from 'astro';
import { getSessionCookie } from '../../../middleware';

export const prerender = false;

const CONTROL_CHECK_PATH = '/api/control-check';
const CONTROL_ENDPOINT = 'notes/create-note';
const SUCCESS_REDIRECT = '/notes';

export async function POST({ request }: APIContext) {
	const formData = await request.formData();
	const title = formData.get('title')?.toString().trim();
	const description = formData.get('description')?.toString().trim();
	const content = formData.get('content')?.toString().trim();

	if (!title || !description || !content) {
		return new Response(JSON.stringify({ message: 'Missing required fields.' }), {
			status: 400,
			headers: { 'Content-Type': 'application/json' },
		});
	}

    const cookieValue = getSessionCookie(request.headers.get('cookie'));



	const controlData = new FormData();
    controlData.set('token', cookieValue ?? '');
	controlData.set('endpoint', CONTROL_ENDPOINT);
	controlData.set('title', title);
	controlData.set('description', description);
	controlData.set('content', content);

	let controlResponse: Response;
	try {
		controlResponse = await fetch(new URL(CONTROL_CHECK_PATH, request.url), {
			method: 'POST',
			body: controlData,
		});
	} catch (error) {
		console.error('Control check request failed:', error);
		return new Response(JSON.stringify({ message: 'Service unavailable.' }), {
			status: 503,
			headers: { 'Content-Type': 'application/json' },
		});
	}

	const raw = await controlResponse.text();
	let controlPayload: Record<string, unknown> | null = null;
	if (raw) {
		try {
			controlPayload = JSON.parse(raw);
		} catch {
			controlPayload = { message: raw };
		}
	}

	if (!controlResponse.ok) {
		return new Response(JSON.stringify(controlPayload ?? { message: 'Create failed.' }), {
			status: controlResponse.status,
			headers: { 'Content-Type': 'application/json' },
		});
	}

	const redirectTarget =
		typeof controlPayload?.redirect === 'string' && controlPayload.redirect.trim()
			? controlPayload.redirect.toString()
			: SUCCESS_REDIRECT;

	return new Response(null, {
		status: 303,
		headers: { Location: redirectTarget },
	});
}
