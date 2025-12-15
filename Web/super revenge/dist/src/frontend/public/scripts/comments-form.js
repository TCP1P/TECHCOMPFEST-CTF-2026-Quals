const FORM_SELECTOR = '[data-add-note-form]';
const STATUS_SELECTOR = '[data-status]';

function attachHandlers() {
	document.querySelectorAll(FORM_SELECTOR).forEach((form) => {
		if (form.dataset.commentsBound === 'true') return;
		form.dataset.commentsBound = 'true';

		form.addEventListener('submit', async (event) => {
			event.preventDefault();

			const status = form.querySelector(STATUS_SELECTOR);
			const textarea = form.querySelector('textarea[name="content"]');
			const action = form.dataset.action ?? '';
			const noteIdValue = Number(form.dataset.noteId);

			if (!action) {
				if (status) status.textContent = 'Missing submit action.';
				return;
			}
			if (!Number.isFinite(noteIdValue)) {
				if (status) status.textContent = 'Invalid note reference.';
				return;
			}

			const content = textarea?.value?.trim() ?? '';
			if (!content) {
				if (status) status.textContent = 'Content is required.';
				return;
			}

			if (status) status.textContent = 'Saving...';

			try {
				const response = await fetch(action, {
					method: 'POST',
					headers: { 'Content-Type': 'application/json' },
					body: JSON.stringify({ noteId: noteIdValue, content }),
				});

				if (!response.ok) {
					const payload = await response.json().catch(() => ({}));
					throw new Error(
						typeof payload.message === 'string' ? payload.message : 'Failed to add comment.',
					);
				}

				form.reset();
				if (status) status.textContent = 'Comment added successfully.';
			} catch (error) {
				const message = error instanceof Error ? error.message : 'Could not add comment.';
				if (status) status.textContent = message;
			}
		});
	});
}

if (document.readyState === 'loading') {
	document.addEventListener('DOMContentLoaded', attachHandlers, { once: true });
} else {
	attachHandlers();
}
