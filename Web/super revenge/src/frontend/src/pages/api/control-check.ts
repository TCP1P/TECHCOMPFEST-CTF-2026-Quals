import type { APIContext } from 'astro';
import { apiBackend } from '../../consts';
import { genBase64, genHex, isPrivateIp } from '../../utils';

export const prerender = false;

const jsonHeaders = { 'Content-Type': 'application/json' };


export async function GET() {
	return new Response(JSON.stringify({ message: 'Method Not Allowed', id: genHex() }), {
		status: 405,
		headers: jsonHeaders,
	});
}


export async function POST({ request, clientAddress }: APIContext) {
	const ip = clientAddress;

	if (!isPrivateIp(ip)) {
		console.warn(`Forbidden control-check attempt from IP: ${ip}`);
		return new Response(
			JSON.stringify({ message: 'Forbidden', id: genBase64() }),
			{
				status: 403,
				headers: jsonHeaders,
			},
		);
	}



	const contentType = request.headers.get('content-type') ?? '';
	if (!contentType.includes('multipart/form-data')) {
		return new Response(JSON.stringify({ message: 'Unsupported content type.', id:genBase64() }), {
			status: 415,
			headers: { 'Content-Type': 'application/json' },
		});
	}

	const payload: Record<string, unknown> = {};

	const formData = await request.formData();
	formData.forEach((value, key) => {
		payload[key] = value.toString();
	});

	const apiEndpoint = payload.endpoint?.toString();
	if (!apiEndpoint) {
		return new Response(JSON.stringify({ message: 'Missing endpoint field.', id:genBase64() }), {
			status: 400,
			headers: { 'Content-Type': 'application/json' },
		});
	}
	const headersBack: Record<string, string> = { 'Content-Type': 'application/json' };
	const token = payload.token?.toString();
	if (token) {
		headersBack['Authorization'] = `Bearer ${token}`;
	}

	try {
		const backendResponse = await fetch(`${apiBackend}/${apiEndpoint}`, {
			method: 'POST',
			headers: headersBack,

			body: JSON.stringify(payload),
		});



		const raw = await backendResponse.text();
		let backendPayload: Record<string, unknown> | null = null;

		if (raw) {
			try {
				backendPayload = JSON.parse(raw);
			} catch {
				backendPayload = { message: raw };
			}
		}

		if (!backendResponse.ok) {
			return new Response(JSON.stringify(backendPayload ?? { message: 'Control check failed.', id:genBase64() }), {
				status: backendResponse.status,
				headers: { 'Content-Type': 'application/json' },
			});
		}

		return new Response(
			JSON.stringify({ success: true, data: backendPayload ?? null, redirect: payload.redirect?.toString() ?? null }),
			{
				status: 200,
				headers: { 'Content-Type': 'application/json' },
			},
		);
	} catch (error) {
		console.error('Error during signup relay:', error);
		return new Response(JSON.stringify({ message: 'Signup service unavailable.', id:genBase64() }), {
			status: 502,
			headers: { 'Content-Type': 'application/json' },
		});
	}
}
