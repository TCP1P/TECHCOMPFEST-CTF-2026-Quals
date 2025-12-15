import type { APIContext } from 'astro';
import { apiBackend } from '../../consts';
import { getSessionCookie } from '../../middleware';
import { genBase64, genHex } from '../../utils';


export const prerender = false;

const jsonHeaders = { 'Content-Type': 'application/json' };

export async function GET() {
	return new Response(JSON.stringify({ message: 'Method Not Allowed', id:genHex() }), {
		status: 405,
		headers: jsonHeaders,
	});
}

export async function POST({ request }: APIContext) {
	let body: unknown;
	try {
		body = await request.json();
	} catch {
		return new Response(JSON.stringify({ message: 'Invalid JSON payload.' , id:genBase64()}), {
			status: 400,
			headers: jsonHeaders,
		});
	}

	if (!body || typeof body !== 'object') {
		return new Response(JSON.stringify({ message: 'Payload must be an object.', id:genBase64() }), {
			status: 400,
			headers: jsonHeaders,
		});
	}

	const payload = body as Record<string, unknown>;
	const content =
		typeof payload.content === 'string' ? payload.content.trim() : '';
	const rawNoteId = payload.noteId;
	const noteId =
		typeof rawNoteId === 'number'
			? rawNoteId
			: typeof rawNoteId === 'string'
				? Number(rawNoteId)
				: NaN;

	if (!content || Number.isNaN(noteId)) {
		return new Response(JSON.stringify({ message: 'Both content and noteId are required.', id:genBase64() }), {
			status: 400,
			headers: jsonHeaders,
		});
	}

	const cookieToken = getSessionCookie(request.headers.get('cookie'));
	const incomingAuth = request.headers.get('authorization');
	const bearer =
		typeof incomingAuth === 'string' && incomingAuth.trim()
			? incomingAuth.trim()
			: cookieToken
				? `Bearer ${cookieToken}`
				: '';
	const upstreamHeaders: Record<string, string> = { ...jsonHeaders };
	if (bearer) upstreamHeaders.Authorization = bearer.startsWith('Bearer') ? bearer : `Bearer ${bearer}`;

	try {
		const upstreamResponse = await fetch(`${apiBackend}/comments/new`, {
			method: 'POST',
			headers: upstreamHeaders,
			body: JSON.stringify({ content, noteId }),
		});

		const raw = await upstreamResponse.text();
		let payload: unknown = null;
		if (raw) {
			try {
				payload = JSON.parse(raw);
			} catch {
				payload = { message: raw };
			}
		}

		if (!upstreamResponse.ok) {
			return new Response(
				JSON.stringify(
					payload && typeof payload === 'object'
						? payload
						: { message: 'Unable to add comment.' , id:genBase64()},
				),
				{
					status: upstreamResponse.status,
					headers: jsonHeaders,
				},
			);
		}

		return new Response(JSON.stringify({ success: true, data: payload ?? null }), {
			status: 200,
			headers: jsonHeaders,
		});
	} catch (error) {
		console.error('Comment proxy failed:', error);
		return new Response(JSON.stringify({ message: 'Comment service unavailable.' ,id:genBase64()}), {
			status: 502,
			headers: jsonHeaders,
		});
	}
}
