
const form = document.querySelector('form[data-signup-form]');
const toastStack = document.querySelector('.toast-stack');

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
            name: formData.get('name')?.toString().trim() ?? '',
            email: formData.get('email')?.toString().trim() ?? '',
            password: formData.get('password')?.toString() ?? '',
            confirmPassword: formData.get('confirmPassword')?.toString() ?? '',
            urlBack: window.location.href,
        };

        if (payload.password !== payload.confirmPassword) {
            showToast('Access keys do not match.', 'error');
            submitButton?.removeAttribute('disabled');
            return;
        }

        try {
            const response = await fetch('/api/signup', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    Accept: 'application/json',
                },
                body: JSON.stringify(payload),
            });

            const data = await response.json().catch(() => ({}));
            if (!response.ok || !data.success) {
                showToast(data?.message ?? 'Signup failed.', 'error');
            } else {
                showToast('Clearance activated. Redirecting...', 'success');
                const redirectTo =
                    typeof data?.redirect === 'string' && data.redirect.length > 0
                        ? data.redirect
                        : '/signin';
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
