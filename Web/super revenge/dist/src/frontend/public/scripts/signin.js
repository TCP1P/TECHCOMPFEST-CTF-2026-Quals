const form = document.querySelector('form[data-signin-form]');
const toastStack = document.querySelector('.toast-stack');
const paramRedirect = new URLSearchParams(window.location.search).get('redirect');
if (!paramRedirect) {
    window.history.replaceState({}, '', window.location.pathname);
}
const redirect = paramRedirect ?? '/notes';
if (form && toastStack) {
    const showToast = (message, variant = 'success') => {
        const toast = document.createElement('div');
        toast.className = `toast toast-${variant}`;
        toast.textContent = message;
        toastStack.appendChild(toast);
        setTimeout(() => {
            toast.classList.add('hide');
            toast.addEventListener('transitionend', () => toast.remove(), { once: true });
        }, 3500);
    };

    form.addEventListener('submit', async (event) => {
        event.preventDefault();
        const submitButton = form.querySelector('button[type="submit"]');
        submitButton?.setAttribute('disabled', 'true');

        const formData = new FormData(form);
        const payload = {
            email: formData.get('email')?.toString().trim() ?? '',
            password: formData.get('password')?.toString() ?? '',
            redirect: redirect
        };

        if (!payload.email || !payload.password) {
            showToast('Signal ID and Access Key are required.', 'error');
            submitButton?.removeAttribute('disabled');
            return;
        }

        try {
            const response = await fetch('/api/signin', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    Accept: 'application/json',
                },
                body: JSON.stringify(payload),
            });

            const data = await response.json().catch(() => ({}));
            if (!response.ok || !data.success) {
                showToast(data?.message ?? 'Signin failed.', 'error');
            } else {
                showToast('Link established. Redirecting...', 'success');
                const redirectTo =
                    typeof data?.redirect === 'string' && data.redirect.length > 0
                        ? data.redirect
                        : '/notes';
                setTimeout(() => {
                    window.location.href = redirectTo;
                }, 1200);
            }
        } catch {
            showToast('Network error. Try again shortly.', 'error');
        } finally {
            submitButton?.removeAttribute('disabled');
        }
    });
}