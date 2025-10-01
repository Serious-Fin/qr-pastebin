import { type User } from '$lib/user';
import { redirect } from '@sveltejs/kit';
import type { Actions } from './$types';

export const actions: Actions = {
    login: async ({ request, fetch, cookies }) => {
        const data = await request.formData();
        const params: User = {
            name: data.get("name") as string,
            password: data.get("password") as string
        };

        // Try sign user in
        const form = new FormData();
        form.append('name', params.name);
        form.append('password', params.password);

        let redirectTo = ""
        try {
            const response = await fetch('/api/login', {
                method: 'POST',
                body: form
            });
            const body = await response.json()

            if (!response.ok) {
                throw new Error(body.errMsg ?? "Login error")
            }
            redirectTo = body.redirectTo ?? "/"
        } catch (err) {
            if (err instanceof Error) {
                return {
                    errMsg: err.message,
                    user: {
                        name: params.name,
                        password: params.password
                    }
                }
            }
        }
        throw redirect(303, redirectTo)
    }
};