import type { APIContext } from 'astro';
import { genBase36, genHex } from '../../utils';


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
            return new Response(JSON.stringify({ message, id:genBase36() }), {
                status: 415,
                headers: { 'Content-Type': 'application/json' },
            });
        }
        return new Response(message, { status: 415 });
    }

    const email = payload.email?.toString().trim();
    const password = payload.password?.toString();
    

    if (!email || !password) {
        const message = 'Missing required fields.';
        if (acceptsJSON) {
            return new Response(JSON.stringify({ message, id:genBase36() }), {
                status: 400,
                headers: { 'Content-Type': 'application/json' },
            });
        }
        return new Response(message, { status: 400 });
    }


    const forwardPayload = {
        email,
        password,
        endpoint: 'auth/login',
        redirect: '/notes',
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
            console.error('Signin control check failed:', controlPayload);
            const message =
                (typeof controlPayload?.message === 'string' && controlPayload.message) ||
                'Signin failed.';
            if (acceptsJSON) {
                return new Response(JSON.stringify({ message, id:genBase36() }), {
                    status: controlResponse.status || 500,
                    headers: { 'Content-Type': 'application/json' },
                });
            }
            return new Response(message, { status: controlResponse.status || 500 });
        }

        
        if (acceptsJSON) {
            const response =  new Response(JSON.stringify({ success: true, redirect: controlPayload.redirect ?? '/signin' }), {
                status: 200,
                headers: { 'Content-Type': 'application/json' },
            });
            if (controlPayload?.data?.payload?.token) {
                console.info(`Successful signin for ${email}, setting session cookie.`);
                const cookieOptions = [
                    `Path=/`,
                    `HttpOnly`,
                    `SameSite=Lax`,
                    `Max-Age=${60 * 60 * 24 * 30}`, // 30 days
                ];
                if (new URL(request.url).protocol === 'https:') {
                    cookieOptions.push('Secure');
                }
                response.headers.append(
                    'Set-Cookie',
                    `supernote_session=${controlPayload.data.payload.token}; ${cookieOptions.join('; ')}`
                );
            }
            return response;
        }
        return Response.redirect(new URL(controlPayload.redirect ?? '/signin', request.url), 303);
    } catch (error) {
        console.error(`Error passing signin to control check: ${error}`);
        const message = 'Signin service unavailable.';
        if (acceptsJSON) {
            return new Response(JSON.stringify({ message }), {
                status: 502,
                headers: { 'Content-Type': 'application/json' },
            });
        }
        return new Response(message, { status: 502 });
    }
}
