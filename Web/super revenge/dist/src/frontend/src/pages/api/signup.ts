import type { APIContext } from 'astro';

export const prerender = false;

export async function POST({ request }: APIContext) {
	const acceptsJSON = request.headers.get('accept')?.includes('application/json') ?? false;
	const contentType = request.headers.get('content-type') ?? '';
	let payload: Record<string, unknown> = {};

	if (contentType.includes('application/json')) {
		payload = await request.json();
	} else if (
		contentType.includes('multipart/form-data') ||
		contentType.includes('application/x-www-form-urlencoded')
	) {
		const formData = await request.formData();
		formData.forEach((value, key) => {
			payload[key] = value.toString();
		});
	} else {
		const message = 'Unsupported content type.';
		if (acceptsJSON) {
			return new Response(JSON.stringify({ message }), {
				status: 415,
				headers: { 'Content-Type': 'application/json' },
			});
		}
		return new Response(message, { status: 415 });
	}

	const name = payload.name?.toString().trim();
	const email = payload.email?.toString().trim();
	const password = payload.password?.toString();
	const confirmPassword = payload.confirmPassword?.toString();


	if (!name || !email || !password || !confirmPassword) {
		const message = 'Missing required fields.';
		if (acceptsJSON) {
			return new Response(JSON.stringify({ message }), {
				status: 400,
				headers: { 'Content-Type': 'application/json' },
			});
		}
		return new Response(message, { status: 400 });
	}

	if (password !== confirmPassword) {
		const message = 'Passwords do not match.';
		if (acceptsJSON) {
			return new Response(JSON.stringify({ message }), {
				status: 400,
				headers: { 'Content-Type': 'application/json' },
			});
		}
		return new Response(message, { status: 400 });
	}

	const forwardPayload = {
		endpoint: 'auth/register',
		redirect: '/signin',
		name,
		email,
		password,
		confirmPassword,
	};

	try {
		const forwardForm = new FormData();
		Object.entries(forwardPayload).forEach(([key, value]) => {
			if (typeof value === 'string') forwardForm.append(key, value);
		});

		const origin = new URL(request.url).origin;
		const controlUrl = new URL('/api/control-check', request.url);
		const controlResponse = await fetch(controlUrl, {
			method: 'POST',
			headers: { Accept: 'application/json', Origin: origin },
			body: forwardForm,
		});
		const controlPayload = await controlResponse.json().catch(() => ({}));

		if (!controlResponse.ok || !controlPayload?.success) {
			console.error('Signup control check failed:', controlPayload);
			const message =
				(typeof controlPayload?.message === 'string' && controlPayload.message) ||
				'Signup failed.';
			if (acceptsJSON) {
				return new Response(JSON.stringify({ message }), {
					status: controlResponse.status || 500,
					headers: { 'Content-Type': 'application/json' },
				});
			}
			return new Response(message, { status: controlResponse.status || 500 });
		}

		if (acceptsJSON) {
			return new Response(JSON.stringify({ success: true, data: controlPayload, redirect: controlPayload.redirect ?? '/signin' }), {
				status: 200,
				headers: { 'Content-Type': 'application/json' },
			});
		}
		return Response.redirect(new URL(controlPayload.redirect ?? '/signin', request.url), 303);
	} catch (error) {
		console.error('Error passing signup to control check:', error);
		const message = 'Signup service unavailable.';
		if (acceptsJSON) {
			return new Response(JSON.stringify({ message }), {
				status: 502,
				headers: { 'Content-Type': 'application/json' },
			});
		}
		return new Response(message, { status: 502 });
	}
}
